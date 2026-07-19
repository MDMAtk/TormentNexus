package tools

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

func HandleAvGetInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	now := time.Now().Format(time.RFC3339)
	return ok(fmt.Sprintf("Av tool running. Current time: %s", now))
}

func HandleAvProcessFile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	filepath, _ :=getString(args, "filepath")
	if filepath == "" {
		return err("filepath is required")
}

	info, statErr := os.Stat(filepath)
	if statErr != nil {
		return err(statErr.Error())
}

	size := info.Size()
	modTime := info.ModTime().Format(time.RFC3339)
	return ok(fmt.Sprintf("File: %s, Size: %d bytes, Modified: %s", filepath, size, modTime))
}

func HandleAvConvertString(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	input, _ :=getString(args, "input")
	operation, _ :=getString(args, "operation")
	var result string
	switch operation {
	case "uppercase":
		result = strings.ToUpper(input)
	case "lowercase":
		result = strings.ToLower(input)
	default:
		return err("unsupported operation, use 'uppercase' or 'lowercase'")
}

	return ok(result)
}