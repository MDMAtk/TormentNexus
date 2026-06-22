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
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	netLogoURL = "https://ccl.northwestern.edu/netlogo/"
)

func HandleNetLogoDownload(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	version, _ :=getString(args, "version")
	if version == "" {
		return err("version parameter is required")
}

	downloadURL := fmt.Sprintf("%s%s/NetLogo%[2]s.zip", netLogoURL, version, version)
	resp, fetchErr := http.Get(downloadURL)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to download: %s", resp.Status))
}

	zipPath := filepath.Join(os.TempDir(), fmt.Sprintf("NetLogo%s.zip", version))
	zipFile, createErr := os.Create(zipPath)
	if createErr != nil {
		return err(createErr.Error())
}

	defer zipFile.Close()

	_, copyErr := io.Copy(zipFile, resp.Body)
	if copyErr != nil {
		return err(copyErr.Error())
}

	return ok(fmt.Sprintf("Downloaded NetLogo %s to %s", version, zipPath))
}

func HandleNetLogoRun(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	modelPath, _ :=getString(args, "model_path")
	if modelPath == "" {
		return err("model_path parameter is required")
}

	netLogoPath, _ :=getString(args, "netlogo_path")
	if netLogoPath == "" {
		return err("netlogo_path parameter is required")
}

	cmd := exec.CommandContext(ctx, netLogoPath, "--model", modelPath)
	output, runErr := cmd.CombinedOutput()
	if runErr != nil {
		return err(fmt.Sprintf("failed to run NetLogo: %v\nOutput: %s", runErr, string(output)))
}

	return ok(fmt.Sprintf("NetLogo model executed successfully:\n%s", string(output)))
}

func HandleNetLogoInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	modelPath, _ :=getString(args, "model_path")
	if modelPath == "" {
		return err("model_path parameter is required")
}

	modelFile, openErr := os.Open(modelPath)
	if openErr != nil {
		return err(openErr.Error())
}

	defer modelFile.Close()

	modelContent, readErr := io.ReadAll(modelFile)
	if readErr != nil {
		return err(readErr.Error())
}

	// Simple regex to extract model info (this is a simplified version)
	re := regexp.MustCompile(`to-report\\s+([a-zA-Z0-9_]+)\\s*\\(([^)]*)\\)\\s*[\\r\\n]+\\s*report\\s+([^\\r\\n]+)`)
	matches := re.FindAllStringSubmatch(string(modelContent), -1)

	if len(matches) == 0 {
		return ok("No reporter functions found in the model")
}

	var reporters []map[string]string
	for _, match := range matches {
		reporters = append(reporters, map[string]string{
			"name":    match[1],
			"inputs":  match[2],
			"returns": match[3],
		})

	info, marshalErr := json.Marshal(reporters)
	if marshalErr != nil {
		return err(marshalErr.Error())
}

	return ok(fmt.Sprintf("Model reporters:\n%s", string(info)))
}

}

func HandleNetLogoVersion(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	netLogoPath, _ :=getString(args, "netlogo_path")
	if netLogoPath == "" {
		return err("netlogo_path parameter is required")
}

	cmd := exec.CommandContext(ctx, netLogoPath, "--version")
	output, runErr := cmd.CombinedOutput()
	if runErr != nil {
		return err(fmt.Sprintf("failed to get version: %v\nOutput: %s", runErr, string(output)))
}

	version := strings.TrimSpace(string(output))
	if version == "" {
		return err("could not determine NetLogo version")
}

	return ok(fmt.Sprintf("NetLogo version: %s", version))
}

func HandleNetLogoConvert(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	modelPath, _ :=getString(args, "model_path")
	if modelPath == "" {
		return err("model_path parameter is required")
}

	netLogoPath, _ :=getString(args, "netlogo_path")
	if netLogoPath == "" {
		return err("netlogo_path parameter is required")
}

	outputPath, _ :=getString(args, "output_path")
	if outputPath == "" {
		return err("output_path parameter is required")
}

	cmd := exec.CommandContext(ctx, netLogoPath, "--convert", modelPath, outputPath)
	output, runErr := cmd.CombinedOutput()
	if runErr != nil {
		return err(fmt.Sprintf("failed to convert model: %v\nOutput: %s", runErr, string(output)))
}

	return ok(fmt.Sprintf("Model converted successfully to %s", outputPath))
}

func HandleNetLogoDocumentation(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	version, _ :=getString(args, "version")
	if version == "" {
		return err("version parameter is required")
}

	docURL := fmt.Sprintf("%s%s/documentation.html", netLogoURL, version)
	resp, fetchErr := http.Get(docURL)
	if fetchErr != nil {
		return err(fetchErr.Error())
}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err(fmt.Sprintf("failed to fetch documentation: %s", resp.Status))
}

	docContent, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return err(readErr.Error())
}

	return ok(fmt.Sprintf("NetLogo %s documentation:\n%s", version, string(docContent)))
}