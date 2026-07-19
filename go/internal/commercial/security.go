package commercial

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// SecurityProvider defines the interface for commercial-grade security features.
type SecurityProvider interface {
	ValidateSSO(ctx context.Context, token string) (bool, error)
	Authorize(ctx context.Context, userID string, resource string, action string) (bool, error)
}

// CommercialWrapper wraps the core execution engine with commercial security.
type CommercialWrapper struct {
	provider   SecurityProvider
	configPath string
	mu         sync.RWMutex
	ssoConfig  map[string]string
	roles      []map[string]any
}

type CommercialConfig struct {
	SSO   map[string]string `json:"sso"`
	Roles []map[string]any  `json:"roles"`
}

// NewCommercialWrapper creates a new wrapper with the given provider.
func NewCommercialWrapper(provider SecurityProvider, workspaceRoot string) *CommercialWrapper {
	cfgPath := filepath.Join(workspaceRoot, ".tormentnexus", "commercial_config.json")
	ew := &CommercialWrapper{
		provider:   provider,
		configPath: cfgPath,
		ssoConfig:  make(map[string]string),
		roles:      defaultRoles(),
	}
	_ = ew.Load()
	return ew
}

func defaultRoles() []map[string]any {
	return []map[string]any{
		{"name": "admin", "description": "Full system access", "permissions": []string{"read", "write", "admin", "audit"}},
		{"name": "operator", "description": "Daily operations", "permissions": []string{"read", "write", "execute"}},
		{"name": "viewer", "description": "Read-only access", "permissions": []string{"read"}},
	}
}

// Load loads the configuration from disk.
func (ew *CommercialWrapper) Load() error {
	ew.mu.Lock()
	defer ew.mu.Unlock()

	data, err := os.ReadFile(ew.configPath)
	if err != nil {
		return err
	}

	var config CommercialConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	if config.SSO != nil {
		ew.ssoConfig = config.SSO
	}
	if config.Roles != nil {
		ew.roles = config.Roles
	}
	return nil
}

// Save saves the current configuration to disk.
func (ew *CommercialWrapper) Save() error {
	ew.mu.RLock()
	config := CommercialConfig{
		SSO:   ew.ssoConfig,
		Roles: ew.roles,
	}
	ew.mu.RUnlock()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(ew.configPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(ew.configPath, data, 0644)
}

// Info returns commercial license and security info.
func (ew *CommercialWrapper) Info() map[string]any {
	ew.mu.RLock()
	defer ew.mu.RUnlock()

	return map[string]any{
		"valid":       true,
		"licensedTo":  "TormentNexus Commercial",
		"tier":        "commercial",
		"maxNodes":    10,
		"features":    []string{"sso", "rbac", "audit", "encryption"},
		"expiresAt":   "",
		"ssoSettings": ew.ssoConfig,
	}
}

// GetRoles returns the available RBAC roles.
func (ew *CommercialWrapper) GetRoles() []map[string]any {
	ew.mu.RLock()
	defer ew.mu.RUnlock()
	return ew.roles
}

// UpdateSSO updates the SSO configuration.
func (ew *CommercialWrapper) UpdateSSO(sso map[string]string) error {
	ew.mu.Lock()
	ew.ssoConfig = sso
	ew.mu.Unlock()
	return ew.Save()
}

// UpdateRoles updates the RBAC roles.
func (ew *CommercialWrapper) UpdateRoles(roles []map[string]any) error {
	ew.mu.Lock()
	ew.roles = roles
	ew.mu.Unlock()
	return ew.Save()
}

// Middleware provides an HTTP middleware for commercial security checks.
func (ew *CommercialWrapper) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Commercial-SSO")
		if token != "" && ew.provider != nil {
			valid, err := ew.provider.ValidateSSO(r.Context(), token)
			if err != nil || !valid {
				http.Error(w, "Unauthorized: Invalid SSO token", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
