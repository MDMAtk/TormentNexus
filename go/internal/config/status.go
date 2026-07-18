package config

import (
	"os"
	"path/filepath"
)

type PathStatus struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
}

type Status struct {
	Host                 string     `json:"host"`
	Port                 int        `json:"port"`
	BaseURL              string     `json:"baseUrl"`
	WorkspaceRoot        PathStatus `json:"workspaceRoot"`
	ConfigDir            PathStatus `json:"configDir"`
	MainConfigDir        PathStatus `json:"mainConfigDir"`
	HyperNexusConfigFile       PathStatus `json:"hypernexusConfigFile"`
	MCPConfigFile        PathStatus `json:"mcpConfigFile"`
	GoLockPath           PathStatus `json:"goLockPath"`
	MainLockPath         PathStatus `json:"mainLockPath"`
	ImportedInstructions PathStatus `json:"importedInstructions"`
	SectionedMemoryStore PathStatus `json:"sectionedMemoryStore"`
	LegacyMemoryStore    PathStatus `json:"legacyMemoryStore"`
	HyperNexusSubmodule        PathStatus `json:"hypernexusSubmodule"`
}

func Snapshot(cfg Config) Status {
	return Status{
		Host:                 cfg.Host,
		Port:                 cfg.Port,
		BaseURL:              cfg.BaseURL(),
		WorkspaceRoot:        buildPathStatus(cfg.WorkspaceRoot),
		ConfigDir:            buildPathStatus(cfg.ConfigDir),
		MainConfigDir:        buildPathStatus(cfg.MainConfigDir),
		HyperNexusConfigFile:       buildPathStatus(filepath.Join(cfg.WorkspaceRoot, "hypernexus.config.json")),
		MCPConfigFile:        buildPathStatus(filepath.Join(cfg.WorkspaceRoot, "mcp.jsonc")),
		GoLockPath:           buildPathStatus(cfg.LockPath()),
		MainLockPath:         buildPathStatus(cfg.MainLockPath()),
		ImportedInstructions: buildPathStatus(cfg.ImportedInstructionsPath()),
		SectionedMemoryStore: buildPathStatus(filepath.Join(cfg.WorkspaceRoot, ".hypernexus", "sectioned_memory.json")),
		LegacyMemoryStore:    buildPathStatus(filepath.Join(cfg.WorkspaceRoot, ".hypernexus", "claude_mem.json")),
		HyperNexusSubmodule:        buildPathStatus(filepath.Join(cfg.WorkspaceRoot, "submodules", "hypernexus")),
	}
}

func buildPathStatus(path string) PathStatus {
	_, err := os.Stat(path)
	return PathStatus{
		Path:   path,
		Exists: err == nil,
	}
}
