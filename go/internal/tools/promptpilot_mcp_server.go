package tools

import (
	"context"
)

func HandleGetPrompt(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ :=getString(args, "name")
	if name == "" {
		name = "World"
	}
	return success(map[string]interface{}{
}
		"result": "Hello " + name + "! Here is your prompt.",
	})
}