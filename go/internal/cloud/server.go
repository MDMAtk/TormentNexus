package cloud

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// CloudServer represents the cloud server
type CloudServer struct {
	am          *AccountManager
	provisioner *Provisioner
	backup      *BackupManager
	mcp         *MCPStreamableHTTP
	db          *sql.DB
}

// NewCloudServer creates a new cloud server
func NewCloudServer(db *sql.DB, backupDir string) *CloudServer {
	am := NewAccountManager(db)
	provisioner := NewProvisioner()
	backup := NewBackupManager(backupDir)
	mcp := NewMCPStreamableHTTP(am, provisioner)

	return &CloudServer{
		am:          am,
		provisioner: provisioner,
		backup:      backup,
		mcp:         mcp,
		db:          db,
	}
}

// Initialize initializes the cloud server
func (s *CloudServer) Initialize() error {
	// Initialize account manager
	if err := s.am.Initialize(); err != nil {
		return fmt.Errorf("initialize account manager: %w", err)
	}

	// Initialize provisioner
	if err := s.provisioner.Initialize(); err != nil {
		return fmt.Errorf("initialize provisioner: %w", err)
	}

	// Initialize backup manager
	if err := s.backup.Initialize(); err != nil {
		return fmt.Errorf("initialize backup manager: %w", err)
	}

	return nil
}

// RegisterRoutes registers all cloud routes
func (s *CloudServer) RegisterRoutes(mux *http.ServeMux) {
	// Account routes
	mux.HandleFunc("/api/cloud/auth/register", s.handleRegister)
	mux.HandleFunc("/api/cloud/auth/login", s.handleLogin)
	mux.HandleFunc("/api/cloud/auth/logout", s.handleLogout)
	mux.HandleFunc("/api/cloud/auth/me", s.handleMe)
	mux.HandleFunc("/api/cloud/auth/api-key", s.handleAPIKey)

	// Account management routes
	mux.HandleFunc("/api/cloud/account", s.handleAccount)
	mux.HandleFunc("/api/cloud/account/plan", s.handlePlan)

	// Container routes
	mux.HandleFunc("/api/cloud/container", s.handleContainer)
	mux.HandleFunc("/api/cloud/container/start", s.handleContainerStart)
	mux.HandleFunc("/api/cloud/container/stop", s.handleContainerStop)
	mux.HandleFunc("/api/cloud/container/destroy", s.handleContainerDestroy)

	// Backup routes
	mux.HandleFunc("/api/cloud/backup", s.handleBackup)
	mux.HandleFunc("/api/cloud/backup/restore", s.handleBackupRestore)
	mux.HandleFunc("/api/cloud/backup/list", s.handleBackupList)

	// MCP routes
	s.mcp.RegisterRoutes(mux)

	// Stripe routes
	s.RegisterStripeRoutes(mux)

	// Health check
	mux.HandleFunc("/api/cloud/health", s.handleHealth)
}

// handleRegister handles user registration
func (s *CloudServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	account, err := s.am.CreateAccount(req.Email, req.Name, req.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// Provision container
	ctx := r.Context()
	container, err := s.provisioner.Provision(ctx, account.ID, string(account.Plan))
	if err != nil {
		log.Printf("[Cloud] Error provisioning container: %v", err)
	} else {
		s.am.UpdateContainerID(account.ID, container.ID)
		account.ContainerID = container.ID
	}

	// Schedule automatic backups
	s.backup.ScheduleBackup(account.ID, account.ContainerID, 24*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"account":   account,
		"container": container,
	})
}

// handleLogin handles user login
func (s *CloudServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	account, err := s.am.Authenticate(req.Email, req.Password)
	if err != nil {
		http.Error(w, `{"error": "invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// Create session token
	token := account.APIKey

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"account": account,
		"token":   token,
	})
}

// handleLogout handles user logout
func (s *CloudServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you would invalidate the session
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// handleMe handles getting the current user
func (s *CloudServer) handleMe(w http.ResponseWriter, r *http.Request) {
	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"account": account,
	})
}

// handleAPIKey handles API key regeneration
func (s *CloudServer) handleAPIKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Generate new API key
	apiKey, err := generateAPIKey()
	if err != nil {
		http.Error(w, `{"error": "failed to generate api key"}`, http.StatusInternalServerError)
		return
	}

	// Update account
	_, err = s.db.Exec("UPDATE accounts SET api_key = ?, updated_at = ? WHERE id = ?", apiKey, time.Now(), account.ID)
	if err != nil {
		http.Error(w, `{"error": "failed to update api key"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"api_key": apiKey,
	})
}

// handleAccount handles account updates
func (s *CloudServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"account": account,
		})
	case http.MethodPut:
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
			return
		}

		_, err = s.db.Exec("UPDATE accounts SET name = ?, updated_at = ? WHERE id = ?", req.Name, time.Now(), account.ID)
		if err != nil {
			http.Error(w, `{"error": "failed to update account"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
		})
	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handlePlan handles plan updates
func (s *CloudServer) handlePlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		Plan string `json:"plan"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	// Update plan
	if err := s.am.UpdatePlan(account.ID, Plan(req.Plan)); err != nil {
		http.Error(w, `{"error": "failed to update plan"}`, http.StatusInternalServerError)
		return
	}

	// Update container resources
	container, err := s.provisioner.GetContainerByAccount(account.ID)
	if err == nil {
		container.MemoryLimit = s.provisioner.getMemoryLimit(req.Plan)
		container.CPULimit = s.provisioner.getCPULimit(req.Plan)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"plan":    req.Plan,
	})
}

// handleContainer handles container operations
func (s *CloudServer) handleContainer(w http.ResponseWriter, r *http.Request) {
	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	container, err := s.provisioner.GetContainerByAccount(account.ID)
	if err != nil {
		http.Error(w, `{"error": "container not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"container": container,
	})
}

// handleContainerStart handles container start
func (s *CloudServer) handleContainerStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	container, err := s.provisioner.GetContainerByAccount(account.ID)
	if err != nil {
		http.Error(w, `{"error": "container not found"}`, http.StatusNotFound)
		return
	}

	if err := s.provisioner.Start(r.Context(), container.ID); err != nil {
		http.Error(w, `{"error": "failed to start container"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// handleContainerStop handles container stop
func (s *CloudServer) handleContainerStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	container, err := s.provisioner.GetContainerByAccount(account.ID)
	if err != nil {
		http.Error(w, `{"error": "container not found"}`, http.StatusNotFound)
		return
	}

	if err := s.provisioner.Stop(r.Context(), container.ID); err != nil {
		http.Error(w, `{"error": "failed to stop container"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// handleContainerDestroy handles container destruction
func (s *CloudServer) handleContainerDestroy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	container, err := s.provisioner.GetContainerByAccount(account.ID)
	if err != nil {
		http.Error(w, `{"error": "container not found"}`, http.StatusNotFound)
		return
	}

	if err := s.provisioner.Destroy(r.Context(), container.ID); err != nil {
		http.Error(w, `{"error": "failed to destroy container"}`, http.StatusInternalServerError)
		return
	}

	// Update account
	s.am.UpdateContainerID(account.ID, "")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// handleBackup handles backup creation
func (s *CloudServer) handleBackup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	container, err := s.provisioner.GetContainerByAccount(account.ID)
	if err != nil {
		http.Error(w, `{"error": "container not found"}`, http.StatusNotFound)
		return
	}

	backupFile, err := s.backup.CreateBackup(account.ID, container.ID)
	if err != nil {
		http.Error(w, `{"error": "failed to create backup"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"backup_file": backupFile,
	})
}

// handleBackupRestore handles backup restoration
func (s *CloudServer) handleBackupRestore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	container, err := s.provisioner.GetContainerByAccount(account.ID)
	if err != nil {
		http.Error(w, `{"error": "container not found"}`, http.StatusNotFound)
		return
	}

	var req struct {
		BackupFile string `json:"backup_file"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	if err := s.backup.RestoreBackup(req.BackupFile, container.ID); err != nil {
		http.Error(w, `{"error": "failed to restore backup"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

// handleBackupList handles listing backups
func (s *CloudServer) handleBackupList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	account, err := s.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	backups, err := s.backup.ListBackups(account.ID)
	if err != nil {
		http.Error(w, `{"error": "failed to list backups"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"backups": backups,
	})
}

// handleHealth handles health check
func (s *CloudServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

// authenticate authenticates a request
func (s *CloudServer) authenticate(r *http.Request) (*Account, error) {
	// Check API key header
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != "" {
		return s.am.GetByAPIKey(apiKey)
	}

	// Check Authorization header
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		apiKey = auth[7:]
		return s.am.GetByAPIKey(apiKey)
	}

	// Check query parameter
	apiKey = r.URL.Query().Get("api_key")
	if apiKey != "" {
		return s.am.GetByAPIKey(apiKey)
	}

	return nil, fmt.Errorf("no authentication provided")
}
