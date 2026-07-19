package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/MDMAtk/TormentNexus/internal/memorystore"
	"github.com/MDMAtk/TormentNexus/internal/repograph"
)

var GlobalVectorStore *memorystore.VectorStore
var GlobalRepoGraph *repograph.RepoGraphService

type ToolResponse struct {
	Content []TextContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func ok(text string) (ToolResponse, error) {
	return ToolResponse{
		Content: []TextContent{{Type: "text", Text: text}},
	}, nil
}

func err(msg string) (ToolResponse, error) {
	return ToolResponse{}, fmt.Errorf("%s", msg)
}

func getString(args map[string]interface{}, key string) (string, bool) {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

func getInt(args map[string]interface{}, key string) (int, bool) {
	if v, ok := args[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val), true
		case int:
			return val, true
		}
	}
	return 0, false
}

func getFloat(args map[string]interface{}, key string) (float64, bool) {
	v, ok := args[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	}
	return 0, false
}

func getBool(args map[string]interface{}, key string) (bool, bool) {
	if v, ok := args[key]; ok {
		if b, ok := v.(bool); ok {
			return b, true
		}
	}
	return false, false
}

type ToolHandler func(ctx context.Context, args map[string]interface{}) (ToolResponse, error)

type Registry struct {
	mu       sync.RWMutex
	handlers map[string]ToolHandler
}

func NewRegistry() *Registry {
	r := &Registry{
		handlers: make(map[string]ToolHandler),
	}
	// Built-in tools — handlers are in server.go (the only clean handler file).
	// Additional tools are served via the MCP server (mcp_server.go) and
	// mcpimpl dispatch (5,400+ generated handlers).
	r.Register("echo", HandleEcho)
	r.Register("hello_world", HandleHelloWorld)
	r.Register("codebase_search", HandleCodebaseSearch)
	r.Register("codebase_outline", HandleCodebaseOutline)
	// TN-native memory tools — same backend as /api/memory/* HTTP API and tn_memory_store pi extension
	r.Register("add_memory", HandleAddMemory)
	r.Register("search_memory", HandleSearchMemory)
	r.Register("delete_memory", HandleDeleteMemory)
	r.Register("memory_stats", HandleMemoryStats)
	return r
}

func (r *Registry) Register(name string, handler ToolHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[name] = handler
}

func (r *Registry) Execute(ctx context.Context, name string, args map[string]interface{}) (ToolResponse, error) {
	r.mu.RLock()
	handler, ok := r.handlers[name]
	r.mu.RUnlock()
	if !ok {
		return ToolResponse{}, fmt.Errorf("unknown tool: %s", name)
	}
	return handler(ctx, args)
}

func (r *Registry) HasTool(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.handlers[name]
	return ok
}

func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var names []string
	for name := range r.handlers {
		names = append(names, name)
	}
	return names




















}

var _ = json.Marshal
