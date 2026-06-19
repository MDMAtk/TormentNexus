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
		return "", fmt.Errorf("missing key: %s", key)
}

	strVal, found := val.(string)
	if !found {
		return "", fmt.Errorf("key %s is not a string", key)
}

	return strVal, nil
}

func getInt(args map[string]interface{}, key string) (int, error) {
	val, found := args[key]
	if !found {
		return 0, fmt.Errorf("missing key: %s", key)
}

	intVal, found := val.(int)
	if !found {
		return 0, fmt.Errorf("key %s is not an integer", key)
}

	return intVal, nil
}

func getBool(args map[string]interface{}, key string) (bool, error) {
	val, found := args[key]
	if !found {
		return false, fmt.Errorf("missing key: %s", key)
}

	boolVal, found := val.(bool)
	if !found {
		return false, fmt.Errorf("key %s is not a boolean", key)
}

	return boolVal, nil
}

func HandleAddNewElementShape(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// Phase 0: Gather Requirements
	var shapeName string
	var referenceImage string

	// Get shape name
	shapeName, _ :=getString(args, "shapeName")
	if e != nil {
		return err(e)
}

	// Get reference image
	referenceImage, e = getString(args, "referenceImage")
	if e != nil {
		return err(e)
}

	// Phase 1: Sketch & Validate (interactive)
	// This phase is interactive and requires user input.
	// For the purpose of this example, we'll simulate the user's response.
	// In a real implementation, you would use a UI or command line interface to gather user feedback.

	// Simulate user validation of the sketch
	userConfirmation := "yes" // This would be the user's actual response
	if userConfirmation != "yes" {
		return err(fmt.Errorf("user did not confirm sketch"))
}

	// Phase 2: Wire Everything Up (automated)
	// This phase involves updating various files and configurations.
	// For the purpose of this example, we'll simulate the updates.

}