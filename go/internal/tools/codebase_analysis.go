package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/MDMAtk/TormentNexus/internal/repograph"
)

// ensureGraphBuilt ensures that the global repository graph is populated.
func ensureGraphBuilt(ctx context.Context) error {
	if GlobalRepoGraph == nil {
		return fmt.Errorf("GlobalRepoGraph is not initialized")
	}
	if GlobalRepoGraph.GetGraph() == nil {
		_, err := GlobalRepoGraph.Build(ctx)
		if err != nil {
			return fmt.Errorf("failed to build repository graph: %w", err)
		}
	}
	return nil
}

// HandleCodebaseSearch searches the codebase symbol/import/definition graph.
func HandleCodebaseSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, okVal := getString(args, "query")
	if !okVal || query == "" {
		return err("missing required parameter 'query'")
	}

	mode, _ := getString(args, "mode")
	if mode == "" {
		mode = "symbols"
	}

	limit, okLimit := getInt(args, "limit")
	if !okLimit || limit <= 0 {
		limit = 20
	}

	if errBuild := ensureGraphBuilt(ctx); errBuild != nil {
		return err(errBuild.Error())
	}

	var data interface{}
	switch strings.ToLower(mode) {
	case "definitions":
		data = GlobalRepoGraph.FindDefinitions(query)
	case "references":
		data = GlobalRepoGraph.FindReferences(query)
	case "symbols":
		data = GlobalRepoGraph.SearchSymbols(query, limit)
	default:
		return err(fmt.Sprintf("invalid search mode '%s'. Allowed: symbols, definitions, references", mode))
	}

	responseBytes, e := json.Marshal(data)
	if e != nil {
		return err(fmt.Sprintf("failed to marshal search results: %s", e.Error()))
	}

	return ok(string(responseBytes))
}

// HandleCodebaseOutline outlines symbols defined in a file or details a specific symbol.
func HandleCodebaseOutline(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	filePath, hasPath := getString(args, "filePath")
	symbolName, hasSymbol := getString(args, "symbolName")

	if !hasPath && !hasSymbol {
		return err("must provide either 'filePath' or 'symbolName'")
	}

	if errBuild := ensureGraphBuilt(ctx); errBuild != nil {
		return err(errBuild.Error())
	}

	if hasPath {
		// Clean filePath path separators for consistency
		cleanPath := filepathToSlash(filePath)
		graph := GlobalRepoGraph.GetGraph()
		if graph == nil {
			return err("graph is nil")
		}

		var fileNodes []*repograph.Node
		for _, node := range graph.Nodes {
			// Check path matches (ignoring node type NodeFile to list definition structures)
			if filepathToSlash(node.Path) == cleanPath && node.Type != repograph.NodeFile && node.Type != repograph.NodeImport {
				fileNodes = append(fileNodes, node)
			}
		}

		// Sort by line start to outline in sequential order
		sort.Slice(fileNodes, func(i, j int) bool {
			return fileNodes[i].LineStart < fileNodes[j].LineStart
		})

		responseBytes, e := json.Marshal(fileNodes)
		if e != nil {
			return err(fmt.Sprintf("failed to marshal file outline: %s", e.Error()))
		}
		return ok(string(responseBytes))
	}

	// Lookup specific symbol definition
	defs := GlobalRepoGraph.FindDefinitions(symbolName)
	responseBytes, e := json.Marshal(defs)
	if e != nil {
		return err(fmt.Sprintf("failed to marshal symbol definitions: %s", e.Error()))
	}
	return ok(string(responseBytes))
}

// Helper to convert windows path slashes to forward slashes to match repo graph structure.
func filepathToSlash(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
