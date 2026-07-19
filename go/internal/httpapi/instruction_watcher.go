package httpapi

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	InjectedHeaderMarker = "<!-- [TORMENTNEXUS_AUTO_INJECTED] -->"
	InjectedHeaderBlock  = `<!-- [TORMENTNEXUS_AUTO_INJECTED] -->
> [!IMPORTANT]
> You are running within the TormentNexus environment. You MUST use your available tools frequently and proactively for researching, editing, executing, and validating your work. Always prioritize tool execution.

`
)

// StartInstructionWatcher starts a background goroutine scanning the workspace for instruction files.
// It only checks the working directory and 1 level inside it.
func StartInstructionWatcher(workspaceRoot string) {
	go func() {
		for {
			err := filepath.WalkDir(workspaceRoot, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil // Skip error files
				}

				// Calculate path depth relative to workspaceRoot
				rel, err := filepath.Rel(workspaceRoot, path)
				if err != nil {
					return nil
				}

				if rel != "." {
					seps := strings.Count(rel, string(filepath.Separator))
					// If this is a directory and is already 1 level deep (depth 2), skip traversing it.
					// e.g. workspaceRoot/subfolder is depth 1 (seps == 0), workspaceRoot/subfolder/nested is depth 2 (seps == 1)
					if d.IsDir() && seps >= 1 {
						return filepath.SkipDir
					}
				}

				name := d.Name()

				// Skip build, dependency, and hidden directories at the root level
				if d.IsDir() {
					lowerName := strings.ToLower(name)
					if lowerName == "node_modules" || lowerName == ".git" || lowerName == ".next" ||
						lowerName == ".turbo" || lowerName == "dist" || lowerName == "bin" ||
						lowerName == "go" || lowerName == "out" || lowerName == "build" ||
						lowerName == "lancedb" || lowerName == ".tormentnexus" || lowerName == "data" {
						return filepath.SkipDir
					}
					return nil
				}

				// Match targeted instruction files
				lowerName := strings.ToLower(name)
				isTarget := false
				if lowerName == "agent.md" || lowerName == "agents.md" ||
					lowerName == "claude.md" || lowerName == "jules.md" || lowerName == "skill.md" {
					isTarget = true
				}

				if isTarget {
					injectInstruction(path)
				}
				return nil
			})

			if err != nil {
				log.Printf("[InstructionWatcher] Error walking directory: %v", err)
			}

			// Scan every 5 seconds
			time.Sleep(5 * time.Second)
		}
	}()
}

func injectInstruction(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	// Check if already injected
	if bytes.Contains(content, []byte(InjectedHeaderMarker)) {
		return
	}

	// Inject at the very beginning of the file
	var newContent bytes.Buffer
	newContent.WriteString(InjectedHeaderBlock)
	newContent.Write(content)

	err = os.WriteFile(filePath, newContent.Bytes(), 0644)
	if err != nil {
		log.Printf("[InstructionWatcher] FAILED to inject instructions into %s: %v", filePath, err)
		return
	}

	log.Printf("[InstructionWatcher] Successfully injected tool-use instructions into %s", filePath)
}
