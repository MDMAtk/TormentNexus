package httpapi

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Account — represents a HyperNexus tenant account
type Account struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Company       string `json:"company"`
	Subdomain     string `json:"subdomain"`
	Plan          string `json:"plan"`
	Seats         int    `json:"seats"`
	Active        bool   `json:"active"`
	ProvisionedAt string `json:"provisioned_at,omitempty"`
	CreatedAt     string `json:"created_at"`
}

// LoginRequest — login payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// handleAccountRegister — POST /api/account/register
func (s *Server) handleAccountRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Company   string `json:"company"`
		Subdomain string `json:"subdomain"`
		Plan      string `json:"plan"`
		Seats     int    `json:"seats"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password required"})
		return
	}
	if req.Subdomain == "" {
		req.Subdomain = strings.ToLower(strings.ReplaceAll(req.Company, " ", "-"))
	}
	if req.Plan == "" {
		req.Plan = "basic"
	}
	if req.Seats == 0 {
		req.Seats = 1
	}

	// Hash password
	hash, salt := hashPassword(req.Password)
	id := generateID()

	now := time.Now().UTC().Format(time.RFC3339)

	// Store in SQLite
	db := s.ensureAccountDB()
	_, err := db.Exec(`INSERT INTO accounts (id, email, password_hash, salt, company, subdomain, plan, seats, active, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?)`,
		id, req.Email, hash, salt, req.Company, req.Subdomain, req.Plan, req.Seats, now)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "email already registered"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create account"})
		return
	}

	writeJSON(w, http.StatusCreated, Account{
		ID:        id,
		Email:     req.Email,
		Company:   req.Company,
		Subdomain: req.Subdomain,
		Plan:      req.Plan,
		Seats:     req.Seats,
		Active:    true,
		CreatedAt: now,
	})
}

// handleAccountLogin — POST /api/account/login
func (s *Server) handleAccountLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	db := s.ensureAccountDB()
	var id, hash, salt, company, subdomain, plan string
	var seats int
	err := db.QueryRow(`SELECT id, password_hash, salt, company, subdomain, plan, seats FROM accounts WHERE email = ? AND active = 1`,
		req.Email).Scan(&id, &hash, &salt, &company, &subdomain, &plan, &seats)
	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "database error"})
		return
	}

	if !verifyPassword(req.Password, hash, salt) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		return
	}

	// Generate session token
	token := generateID()
	db.Exec(`UPDATE accounts SET session_token = ?, last_login = ? WHERE id = ?`,
		token, time.Now().UTC().Format(time.RFC3339), id)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"token":     token,
		"account":   Account{ID: id, Email: req.Email, Company: company, Subdomain: subdomain, Plan: plan, Seats: seats, Active: true},
		"dashboard": fmt.Sprintf("https://%s.hypernexus.site", subdomain),
	})
}

// handleAccountProvision — POST /api/account/provision (called by Stripe webhook or manually)
func (s *Server) handleAccountProvision(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID string `json:"account_id"`
		Email     string `json:"email"`
		Plan      string `json:"plan"`
		Seats     int    `json:"seats"`
		Subdomain string `json:"subdomain"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	db := s.ensureAccountDB()

	// If creating from Stripe webhook, get/create account
	var company, subdomain, plan string
	var seats int
	if req.AccountID != "" {
		err := db.QueryRow(`SELECT company, subdomain, plan, seats FROM accounts WHERE id = ?`,
			req.AccountID).Scan(&company, &subdomain, &plan, &seats)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "account not found"})
			return
		}
	} else {
		// Create from Stripe checkout metadata
		subdomain = req.Subdomain
		if subdomain == "" {
			subdomain = "org-" + generateID()[:8]
		}
		plan = req.Plan
		if plan == "" {
			plan = "basic"
		}
		seats = req.Seats
		if seats == 0 {
			seats = 1
		}
	}

	// Provision containers (best-effort — succeeds even without Docker)
	provisionLog, provisionErr := provisionContainers(subdomain, plan, seats)
	provisionStatus := "provisioned"
	if provisionErr != nil {
		provisionStatus = "account_created_pending_provision"
	}

	now := time.Now().UTC().Format(time.RFC3339)

	// Update account
	if req.AccountID != "" {
		db.Exec(`UPDATE accounts SET active = 1, provisioned_at = ?, plan = COALESCE(NULLIF(?, ''), plan), seats = CASE WHEN ? > 0 THEN ? ELSE seats END WHERE id = ?`,
			now, req.Plan, req.Seats, req.Seats, req.AccountID)
	} else {
		id := generateID()
		hash, salt := hashPassword("changeme")
		db.Exec(`INSERT OR REPLACE INTO accounts (id, email, password_hash, salt, company, subdomain, plan, seats, active, provisioned_at, created_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?)`,
			id, req.Email, hash, salt, req.Email, subdomain, plan, seats, now, now)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":       provisionStatus,
		"subdomain":    subdomain,
		"dashboard":    fmt.Sprintf("https://%s.hypernexus.site", subdomain),
		"provisioning": provisionLog,
	})
}

// handleAccountStatus — GET /api/account/status?token=xxx
func (s *Server) handleAccountStatus(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "token required"})
		return
	}

	db := s.ensureAccountDB()
	var acc Account
	err := db.QueryRow(`SELECT id, email, company, subdomain, plan, seats, COALESCE(provisioned_at,''), created_at FROM accounts WHERE session_token = ? AND active = 1`,
		token).Scan(&acc.ID, &acc.Email, &acc.Company, &acc.Subdomain, &acc.Plan, &acc.Seats, &acc.ProvisionedAt, &acc.CreatedAt)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired session"})
		return
	}
	acc.Active = true

	// Get container stats
	stats := getContainerStats(acc.Subdomain)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"account": acc,
		"stats":   stats,
	})
}

// ensureAccountDB — lazy-init the accounts SQLite database
func (s *Server) ensureAccountDB() *sql.DB {
	if s.accountDB != nil {
		return s.accountDB
	}
	db, err := sql.Open("sqlite", s.cfg.AccountDBPath())
	if err != nil {
		panic("account db: " + err.Error())
	}
	db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE,
		password_hash TEXT,
		salt TEXT,
		company TEXT,
		subdomain TEXT,
		plan TEXT DEFAULT 'basic',
		seats INTEGER DEFAULT 1,
		active INTEGER DEFAULT 0,
		session_token TEXT,
		provisioned_at TEXT,
		last_login TEXT,
		created_at TEXT
	)`)
	s.accountDB = db
	return db
}

// hashPassword — bcrypt-like simple hash with salt
func hashPassword(password string) (hash, salt string) {
	b := make([]byte, 16)
	rand.Read(b)
	salt = hex.EncodeToString(b)
	h := sha256.Sum256([]byte(salt + password + "hypernexus"))
	return hex.EncodeToString(h[:]), salt
}

func verifyPassword(password, hash, salt string) bool {
	h := sha256.Sum256([]byte(salt + password + "hypernexus"))
	return hex.EncodeToString(h[:]) == hash
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// provisionContainers — shell out to the tenant provision script
func provisionContainers(subdomain, plan string, seats int) (string, error) {
	script := "/opt/tormentnexus/deploy/stripe-webhook-provisioner.sh"
	cmd := exec.Command("bash", script, subdomain, "admin@hypernexus.site", subdomain, plan)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// getContainerStats — get stats from running containers
func getContainerStats(subdomain string) map[string]interface{} {
	out, err := exec.Command("docker", "stats", "--no-stream",
		"--format", "{{.Name}},{{.CPUPerc}},{{.MemUsage}},{{.NetIO}},{{.BlockIO}}",
		fmt.Sprintf("tn-core-%s", subdomain),
	).CombinedOutput()
	if err != nil {
		return map[string]interface{}{
			"status":  "not_running",
			"message": "Container not provisioned or not running yet",
		}
	}
	parts := strings.Split(strings.TrimSpace(string(out)), ",")
	if len(parts) >= 5 {
		return map[string]interface{}{
			"status": "running",
			"name":   parts[0],
			"cpu":    parts[1],
			"memory": parts[2],
			"net_io": parts[3],
			"disk":   strings.TrimSpace(parts[4]),
		}
	}
	return map[string]interface{}{"status": "unknown"}
}
