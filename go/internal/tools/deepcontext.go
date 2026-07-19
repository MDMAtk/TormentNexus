package tools

import (
	"context"
	"fmt"
)

func HandleSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ :=getString(args, "query")
	if query == "" {
		return err("query parameter is required")
}

	result := fmt.Sprintf("Search results for query: %s", query)
	return ok(result)
}

func HandleGet(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ :=getString(args, "id")
	if id == "" {
		return err("id parameter is required")
}

	result := fmt.Sprintf("Retrieved document with id: %s", id)
	return ok(result)
}