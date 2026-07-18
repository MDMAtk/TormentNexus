package adapters

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gitlab.com/robertpelloni/HyperNexus/tormentnexus"
)

type HyperNexusStatus struct {
	Assimilated       bool           `json:"assimilated"`
	HyperNexusCoreURL       string         `json:"hypernexusCoreUrl,omitempty"`
	MemoryContext     string         `json:"memoryContext,omitempty"`
	Provider          ProviderStatus `json:"provider"`
	MCPServerNames    []string       `json:"mcpServerNames,omitempty"`
	MCPConfigPath     string         `json:"mcpConfigPath,omitempty"`
	HyperNexusRepoPath string         `json:"hypernexusRepoPath,omitempty"`
	Warnings          []string       `json:"warnings,omitempty"`
}

type HyperNexusAdapter struct {
	hypernexusAdapter *hypernexus.Adapter
	workingDir  string
	homeDir     string
}

func NewHyperNexusAdapter(workingDir string) *HyperNexusAdapter {
	homeDir, _ := os.UserHomeDir()
	return &HyperNexusAdapter{
		hypernexusAdapter: hypernexus.NewAdapter(),
		workingDir:  workingDir,
		homeDir:     homeDir,
	}
}

func (a *HyperNexusAdapter) Status() HyperNexusStatus {
	status := HyperNexusStatus{
		Assimilated:   a.hypernexusAdapter != nil && a.hypernexusAdapter.Assimilated,
		MemoryContext: a.MemoryContext(),
		Provider:      BuildProviderStatus(),
	}
	if a.hypernexusAdapter != nil {
		status.HyperNexusCoreURL = a.hypernexusAdapter.HyperNexusCoreURL
	}
	if repoPath, ok := a.findHyperNexusRepo(); ok {
		status.HyperNexusRepoPath = repoPath
	} else {
		status.Warnings = append(status.Warnings, "adjacent hypernexus repo not found")
	}
	if configPath, names, err := a.listMCPServers(); err == nil {
		status.MCPConfigPath = configPath
		status.MCPServerNames = names
	} else {
		status.Warnings = append(status.Warnings, err.Error())
	}
	return status
}

func (a *HyperNexusAdapter) MemoryContext() string {
	if a.hypernexusAdapter == nil {
		return ""
	}
	return a.hypernexusAdapter.GetMemoryContext()
}

func (a *HyperNexusAdapter) RouteMCP(request string) string {
	if a.hypernexusAdapter == nil {
		return request
	}
	return a.hypernexusAdapter.RouteMCP(request)
}

func (a *HyperNexusAdapter) BuildSystemContext() string {
	status := a.Status()
	parts := []string{
		"[HyperNexus Adapter]",
		fmt.Sprintf("Assimilated: %t", status.Assimilated),
	}
	if status.HyperNexusCoreURL != "" {
		parts = append(parts, fmt.Sprintf("HyperNexus Core URL: %s", status.HyperNexusCoreURL))
	}
	if status.MemoryContext != "" {
		parts = append(parts, status.MemoryContext)
	}
	if len(status.Provider.Available) > 0 {
		parts = append(parts, BuildProviderContext())
	}
	if len(status.MCPServerNames) > 0 {
		parts = append(parts, fmt.Sprintf("Configured MCP servers: %s", strings.Join(status.MCPServerNames, ", ")))
	}
	if status.HyperNexusRepoPath != "" {
		parts = append(parts, fmt.Sprintf("HyperNexus repo: %s", status.HyperNexusRepoPath))
	}
	if len(status.Warnings) > 0 {
		parts = append(parts, fmt.Sprintf("Warnings: %s", strings.Join(status.Warnings, "; ")))
	}
	return strings.Join(parts, "\n")
}

func (a *HyperNexusAdapter) listMCPServers() (string, []string, error) {
	configPath, config, err := ParseMCPConfig(a.homeDir)
	if err != nil {
		return configPath, nil, fmt.Errorf("mcp config unavailable: %w", err)
	}
	names := make([]string, 0, len(config.MCPServers))
	for name := range config.MCPServers {
		names = append(names, name)
	}
	sort.Strings(names)
	return configPath, names, nil
}

func (a *HyperNexusAdapter) findHyperNexusRepo() (string, bool) {
	candidates := []string{}
	if a.workingDir != "" {
		candidates = append(candidates,
			filepath.Join(a.workingDir, "..", "hypernexus"),
			filepath.Join(a.workingDir, "../hypernexus"),
		)
	}
	if a.homeDir != "" {
		candidates = append(candidates, filepath.Join(a.homeDir, "workspace", "hypernexus"))
	}
	for _, candidate := range candidates {
		clean := filepath.Clean(candidate)
		if stat, err := os.Stat(filepath.Join(clean, "README.md")); err == nil && !stat.IsDir() {
			return clean, true
		}
	}
	return "", false
}
