package tools

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func HandleProcessImage(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	imagePath, _ :=getString(args, "image_path")
	operation, _ :=getString(args, "operation")
	data, readErr := os.ReadFile(imagePath)
	if readErr != nil {
		return err("failed to read image: " + readErr.Error())
}

	return ok(fmt.Sprintf("Processed image %s (%d bytes) with operation %s", imagePath, len(data), operation))
}

func HandleListFormats(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	formats := []string{"JPEG", "PNG", "GIF", "BMP", "WebP"}
	return ok(strings.Join(formats, ", "))
}

func HandleApplyFilter(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	imagePath, _ :=getString(args, "image_path")
	filterName, _ :=getString(args, "filter_name")
	intensity, _ :=getInt(args, "intensity")
	_, statErr := os.Stat(imagePath)
	if statErr != nil {
		return err("image not found: " + statErr.Error())
}

	return ok(fmt.Sprintf("Applied filter %s with intensity %d to %s", filterName, intensity, imagePath))
}