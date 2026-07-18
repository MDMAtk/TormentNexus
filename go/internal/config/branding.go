package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Edition represents the product edition
type Edition string

const (
	EditionHyperNexus Edition = "hypernexus" // Open source, free edition
)

// BrandingConfig holds all branding-related configuration
type BrandingConfig struct {
	// Edition determines the product branding and behavior
	Edition Edition `json:"edition"`

	// ProductName is the display name (HyperNexus or HyperNexus)
	ProductName string `json:"product_name"`

	// CompanyName is the company/organization name
	CompanyName string `json:"company_name"`

	// TrayTooltip is the system tray tooltip text
	TrayTooltip string `json:"tray_tooltip"`

	// DashboardTitle is the title shown in the dashboard
	DashboardTitle string `json:"dashboard_title"`

	// CloudEndpoint is the Streamable HTTP endpoint for corporate mode
	// Empty string means local mode
	CloudEndpoint string `json:"cloud_endpoint,omitempty"`

	// CloudAuth is the authentication token for cloud connection
	CloudAuth string `json:"cloud_auth,omitempty"`

	// ConfigDir is the configuration directory name
	ConfigDir string `json:"config_dir"`

	// RegistryKey is the Windows registry key name
	RegistryKey string `json:"registry_key"`
}

var (
	currentConfig *BrandingConfig
	configOnce    sync.Once
	configMutex   sync.RWMutex
)

// GetBranding returns the current branding configuration
func GetBranding() *BrandingConfig {
	configOnce.Do(func() {
		currentConfig = loadBrandingConfig()
	})
	configMutex.RLock()
	defer configMutex.RUnlock()
	return currentConfig
}

// SetBranding updates the branding configuration (for testing or runtime switching)
func SetBranding(cfg *BrandingConfig) {
	configMutex.Lock()
	defer configMutex.Unlock()
	currentConfig = cfg
}

// IsCorporateMode returns true if running in HyperNexus corporate mode
func IsCorporateMode() bool {
	return GetBranding().Edition == EditionHyperNexus
}

// IsCloudConnected returns true if connected to cloud endpoint
func IsCloudConnected() bool {
	return GetBranding().CloudEndpoint != ""
}

// loadBrandingConfig loads branding configuration from file or detects from environment
func loadBrandingConfig() *BrandingConfig {
	// Try to load from config file
	cfg, err := loadFromFile()
	if err == nil {
		return cfg
	}

	// Try to detect from environment variable
	if edition := os.Getenv("HN_EDITION"); edition != "" {
		if edition == "hypernexus" || edition == "corporate" {
			return HyperNexusCorporateBranding()
		}
	}

	// Try to detect from config directory name
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		hyperNexusDir := filepath.Join(homeDir, ".hypernexus")
		if _, err := os.Stat(hyperNexusDir); err == nil {
			return HyperNexusBranding()
		}
	}

	// Default to HyperNexus (open source)
	return HyperNexusBranding()
}

// HyperNexusBranding returns branding for the open source edition
func HyperNexusBranding() *BrandingConfig {
	return &BrandingConfig{
		Edition:        EditionHyperNexus,
		ProductName:    "HyperNexus",
		CompanyName:    "HyperNexus Team",
		TrayTooltip:    "HyperNexus (Running)",
		DashboardTitle: "HyperNexus Dashboard",
		ConfigDir:      ".hypernexus",
		RegistryKey:    "HyperNexus",
	}
}

// HyperNexusCorporateBranding returns branding for the corporate edition
func HyperNexusCorporateBranding() *BrandingConfig {
	return &BrandingConfig{
		Edition:        EditionHyperNexus,
		ProductName:    "HyperNexus",
		CompanyName:    "HyperNexus Corp",
		TrayTooltip:    "HyperNexus (Running)",
		DashboardTitle: "HyperNexus Dashboard",
		ConfigDir:      ".hypernexus",
		RegistryKey:    "HyperNexus",
		// Cloud endpoint can be configured later
		CloudEndpoint: os.Getenv("HN_CLOUD_ENDPOINT"),
		CloudAuth:     os.Getenv("HN_CLOUD_AUTH"),
	}
}

// loadFromFile loads branding configuration from the config file
func loadFromFile() (*BrandingConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Check for HyperNexus config
	hyperNexusConfig := filepath.Join(homeDir, ".hypernexus", "branding.json")
	if _, err := os.Stat(hyperNexusConfig); err == nil {
		return loadFromPath(hyperNexusConfig)
	}

	return nil, os.ErrNotExist
}

// loadFromPath loads branding configuration from a specific file path
func loadFromPath(path string) (*BrandingConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg BrandingConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveToFile saves the current branding configuration to a file
func SaveToFile(cfg *BrandingConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, cfg.ConfigDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "branding.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
