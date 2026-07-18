package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Host          string
	Port          int
	ConfigDir     string
	MainConfigDir string
	WorkspaceRoot string
}

func Default() Config {
	return Config{
		Host:          "127.0.0.1",
		Port:          7778,
		ConfigDir:     DefaultConfigDir(),
		MainConfigDir: DefaultMainConfigDir(),
		WorkspaceRoot: DefaultWorkspaceRoot(),
	}
}

func DefaultConfigDir() string {
	// Check new env var first, then old one for backward compat
	if configured := os.Getenv("HYPERNEXUS_CONFIG_DIR"); configured != "" {
		return configured
	}
	if configured := os.Getenv("HYPERNEXUS_GO_CONFIG_DIR"); configured != "" {
		return configured
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".hypernexus"
	}

	// Migrate from old .hypernexus-go if it exists
	oldDir := filepath.Join(homeDir, ".hypernexus-go")
	newDir := filepath.Join(homeDir, ".hypernexus")
	if oldInfo, oldErr := os.Stat(oldDir); oldErr == nil && oldInfo.IsDir() {
		if newInfo, newErr := os.Stat(newDir); newErr != nil || !newInfo.IsDir() {
			// Old dir exists, new doesn't — migrate
			if err := os.Rename(oldDir, newDir); err == nil {
				fmt.Printf("[Config] Migrated %s → %s\n", oldDir, newDir)
				return newDir
			}
			// Rename failed, try copying
			fmt.Printf("[Config] Migration rename failed: %v, using old dir\n", err)
			return oldDir
		}
	}

	return newDir
}

func DefaultMainConfigDir() string {
	if configured := os.Getenv("HYPERNEXUS_MAIN_CONFIG_DIR"); configured != "" {
		return configured
	}

	// Fallback to legacy env var for backward compatibility
	if legacy := os.Getenv("HYPERNEXUS_CONFIG_DIR"); legacy != "" {
		return legacy
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".hypernexus"
	}

	return filepath.Join(homeDir, ".hypernexus")
}

func DefaultWorkspaceRoot() string {
	if configured := os.Getenv("HYPERNEXUS_WORKSPACE_ROOT"); configured != "" {
		return configured
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return "."
	}

	return workingDir
}

func (c Config) BaseURL() string {
	return fmt.Sprintf("http://%s:%d", browserHost(c.Host), c.Port)
}

func (c Config) LockPath() string {
	return filepath.Join(c.ConfigDir, "lock")
}

func (c Config) MainLockPath() string {
	return filepath.Join(c.MainConfigDir, "lock")
}

func (c Config) ImportedInstructionsPath() string {
	return filepath.Join(c.WorkspaceRoot, ".hypernexus", "imported_sessions", "docs", "auto-imported-agent-instructions.md")
}

// AccountDBPath returns the path to the accounts SQLite database.
func (c Config) AccountDBPath() string {
	return filepath.Join(c.ConfigDir, "accounts.db")
}

func browserHost(host string) string {
	switch host {
	case "", "0.0.0.0", "::", "[::]":
		return "127.0.0.1"
	default:
		return host
	}
}
