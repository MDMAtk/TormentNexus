package tools

import "context"

func HandleListTools(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok(map[string]interface{}{
}
		"tools": []string{"tool1", "tool2"},
	})

func HandleGetTool(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		return err("missing tool name")
}

	return ok(map[string]interface{}{
}
		"name":        name,
		"description": "A sample MCP tool",
	})
}