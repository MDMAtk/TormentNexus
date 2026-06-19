package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func ok(text string) (ToolResponse, error) {
	return ToolResponse{TextContent: text}, nil
}

func err(e error) (ToolResponse, error) {
	return ToolResponse{}, e
}

func getString(args map[string]interface{}, key string) (string, error) {
	val, found := args[key]
	if !found {
		return "", fmt.Errorf("key %s not found", key)
}

	strVal, found := val.(string)
	if !found {
		return "", fmt.Errorf("value for key %s is not a string", key)
}

	return strVal, nil
}

func getInt(args map[string]interface{}, key string) (int, error) {
	strVal, _ :=getString(args, key)
	if e != nil {
		return 0, e
	}
	intVal, e := strconv.Atoi(strVal)
	if e != nil {
		return 0, fmt.Errorf("value for key %s is not an integer: %v", key, e)
}

	return intVal, nil
}

func getBool(args map[string]interface{}, key string) (bool, error) {
	strVal, _ :=getString(args, key)
	if e != nil {
		return false, e
	}
	boolVal, e := strconv.ParseBool(strVal)
	if e != nil {
		return false, fmt.Errorf("value for key %s is not a boolean: %v", key, e)
}

	return boolVal, nil
}

func HandlePing(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("pong")
}

func HandleStart(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// Implement the logic to start MCP clients here
	// This is a placeholder for the actual implementation
	return ok("MCP clients started successfully")
}

func HandleRestart(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {

	return ok("not yet implemented")
}