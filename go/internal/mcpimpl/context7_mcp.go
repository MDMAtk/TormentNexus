package mcpimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Context7 MCP server — up-to-date code documentation for LLMs.
// Original: https://github.com/upstash/context7 (58K stars)
// API: https://context7.com/api

var context7HTTP = &http.Client{Timeout: 15 * time.Second}

// HandleSearchLibraries searches libraries based on a query.
func HandleSearchLibraries(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query")
	if query == "" {
		return err("query is required")
	}

	apiURL := fmt.Sprintf("https://context7.com/api/search?q=%s", url.QueryEscape(query))
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("request build: %v", reqErr))
	}
	req.Header.Set("User-Agent", "TormentNexus/1.0")

	resp, doErr := context7HTTP.Do(req)
	if doErr != nil {
		return err(fmt.Sprintf("search failed: %v", doErr))
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("read response: %v", readErr))
	}

	var results []map[string]interface{}
	if parseErr := json.Unmarshal(body, &results); parseErr != nil {
		return err(fmt.Sprintf("parse response: %v", parseErr))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d libraries for \"%s\":\n\n", len(results), query))
	for i, lib := range results {
		if i >= 10 {
			break
		}
		name, _ := lib["name"].(string)
		desc, _ := lib["description"].(string)
		stars, _ := lib["stars"].(float64)
		url, _ := lib["url"].(string)
		sb.WriteString(fmt.Sprintf("%d. %s", i+1, name))
		if stars > 0 {
			sb.WriteString(fmt.Sprintf(" (★ %.0f)", stars))
		}
		sb.WriteString("\n")
		if desc != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", truncateStr(desc, 150)))
		}
		if url != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", url))
		}
		sb.WriteString("\n")
	}

	return ok(sb.String())
}

// HandleFetchLibraryContext fetches documentation context for a specific library.
func HandleFetchLibraryContext(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	library, _ := getString(args, "library")
	if library == "" {
		return err("library is required (e.g. 'next.js', 'react', 'express')")
	}
	query, _ := getString(args, "query")
	if query == "" {
		query = "getting started"
	}

	apiURL := fmt.Sprintf("https://context7.com/api/context?library=%s&query=%s",
		url.QueryEscape(library), url.QueryEscape(query))
	req, reqErr := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if reqErr != nil {
		return err(fmt.Sprintf("request build: %v", reqErr))
	}
	req.Header.Set("User-Agent", "TormentNexus/1.0")

	resp, doErr := context7HTTP.Do(req)
	if doErr != nil {
		return err(fmt.Sprintf("context fetch failed: %v", doErr))
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(fmt.Sprintf("read response: %v", readErr))
	}

	var result map[string]interface{}
	if parseErr := json.Unmarshal(body, &result); parseErr != nil {
		return err(fmt.Sprintf("parse response: %v", parseErr))
	}

	content, _ := result["content"].(string)
	source, _ := result["source"].(string)
	if content == "" {
		return ok(fmt.Sprintf("No documentation found for \"%s\" about \"%s\".", library, query))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📚 %s — %s\n\n", library, query))
	if source != "" {
		sb.WriteString(fmt.Sprintf("Source: %s\n\n", source))
	}
	sb.WriteString(content)

	return ok(sb.String())
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
