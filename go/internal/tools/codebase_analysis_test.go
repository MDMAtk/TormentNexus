package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/MDMAtk/TormentNexus/internal/repograph"
)

func TestCodebaseAnalysisTools(t *testing.T) {
	// Create a temp directory to simulate a codebase
	tmpDir, err := os.MkdirTemp("", "test_codebase")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write a Go file
	goFileContent := `package dummy

import "fmt"

// HelloGreeting prints hello
func HelloGreeting() {
	fmt.Println("Hello")
}

type MyStruct struct {
	Value string
}
`
	err = os.WriteFile(filepath.Join(tmpDir, "dummy.go"), []byte(goFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write dummy.go: %v", err)
	}

	// Write a Python file
	pyFileContent := `def compute_val(x):
    return x * 2

class PythonClass:
    pass
`
	err = os.WriteFile(filepath.Join(tmpDir, "helper.py"), []byte(pyFileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write helper.py: %v", err)
	}

	// Initialize the GlobalRepoGraph
	GlobalRepoGraph = repograph.NewRepoGraphService(tmpDir)

	ctx := context.Background()

	// Ensure the graph is built
	if err := ensureGraphBuilt(ctx); err != nil {
		t.Fatalf("ensureGraphBuilt failed: %v", err)
	}

	// 1. Test HandleCodebaseSearch in "symbols" mode (default)
	t.Run("CodebaseSearch_Symbols", func(t *testing.T) {
		args := map[string]interface{}{
			"query": "Greeting",
			"mode":  "symbols",
		}
		resp, err := HandleCodebaseSearch(ctx, args)
		if err != nil {
			t.Fatalf("HandleCodebaseSearch error: %v", err)
		}
		if resp.IsError {
			t.Fatalf("HandleCodebaseSearch returned error response")
		}
		if len(resp.Content) == 0 {
			t.Fatalf("Empty content returned")
		}

		var nodes []*repograph.Node
		if err := json.Unmarshal([]byte(resp.Content[0].Text), &nodes); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(nodes) == 0 {
			t.Fatalf("Expected at least 1 matching symbol node, got 0")
		}
		if nodes[0].Name != "HelloGreeting" {
			t.Errorf("Expected node name HelloGreeting, got %s", nodes[0].Name)
		}
	})

	// 2. Test HandleCodebaseSearch in "definitions" mode
	t.Run("CodebaseSearch_Definitions", func(t *testing.T) {
		args := map[string]interface{}{
			"query": "HelloGreeting",
			"mode":  "definitions",
		}
		resp, err := HandleCodebaseSearch(ctx, args)
		if err != nil {
			t.Fatalf("HandleCodebaseSearch error: %v", err)
		}

		var nodes []*repograph.Node
		if err := json.Unmarshal([]byte(resp.Content[0].Text), &nodes); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(nodes) != 1 {
			t.Fatalf("Expected exactly 1 definition node, got %d", len(nodes))
		}
		if nodes[0].Name != "HelloGreeting" {
			t.Errorf("Expected definition name HelloGreeting, got %s", nodes[0].Name)
		}
	})

	// 3. Test HandleCodebaseOutline with filePath
	t.Run("CodebaseOutline_FilePath", func(t *testing.T) {
		args := map[string]interface{}{
			"filePath": "dummy.go",
		}
		resp, err := HandleCodebaseOutline(ctx, args)
		if err != nil {
			t.Fatalf("HandleCodebaseOutline error: %v", err)
		}

		var nodes []*repograph.Node
		if err := json.Unmarshal([]byte(resp.Content[0].Text), &nodes); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// HelloGreeting and MyStruct
		if len(nodes) != 2 {
			t.Fatalf("Expected 2 nodes in outline, got %d", len(nodes))
		}

		// Should be sorted by line start
		if nodes[0].Name != "HelloGreeting" || nodes[1].Name != "MyStruct" {
			t.Errorf("Unexpected sorted order or names: %s, %s", nodes[0].Name, nodes[1].Name)
		}
	})

	// 4. Test HandleCodebaseOutline with symbolName
	t.Run("CodebaseOutline_SymbolName", func(t *testing.T) {
		args := map[string]interface{}{
			"symbolName": "PythonClass",
		}
		resp, err := HandleCodebaseOutline(ctx, args)
		if err != nil {
			t.Fatalf("HandleCodebaseOutline error: %v", err)
		}

		var nodes []*repograph.Node
		if err := json.Unmarshal([]byte(resp.Content[0].Text), &nodes); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(nodes) != 1 {
			t.Fatalf("Expected exactly 1 definition node, got %d", len(nodes))
		}
		if nodes[0].Path != "helper.py" {
			t.Errorf("Expected path helper.py, got %s", nodes[0].Path)
		}
	})
}
