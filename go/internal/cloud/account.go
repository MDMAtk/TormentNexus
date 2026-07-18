package cloud

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Account represents a cloud user account
type Account struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	Plan         Plan      `json:"plan"`
	Status       Status    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLoginAt  time.Time `json:"last_login_at,omitempty"`
	APIKey       string    `json:"api_key,omitempty"`
	ContainerID  string    `json:"container_id,omitempty"`
	Metadata     string    `json:"metadata,omitempty"`
}

// Plan represents the subscription plan
type Plan string

const (
	PlanFree       Plan = "free"
	PlanStarter    Plan = "starter"
	PlanPro        Plan = "pro"
	PlanEnterprise Plan = "enterprise"
)

// Status represents account status
type Status string

const (
	StatusActive    Status = "active"
	StatusSuspended Status = "suspended"
	StatusDeleted   Status = "deleted"
)

// AccountManager handles account operations
type AccountManager struct {
	db *sql.DB
}

// NewAccountManager creates a new account manager
func NewAccountManager(db *sql.DB) *AccountManager {
	return &AccountManager{db: db}
}

// Initialize creates the accounts table
func (am *AccountManager) Initialize() error {
	query := `
	CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		plan TEXT NOT NULL DEFAULT 'free',
		status TEXT NOT NULL DEFAULT 'active',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		last_login_at DATETIME,
		api_key TEXT UNIQUE,
		container_id TEXT,
		metadata TEXT
	);
	
	CREATE INDEX IF NOT EXISTS idx_accounts_email ON accounts(email);
	CREATE INDEX IF NOT EXISTS idx_accounts_api_key ON accounts(api_key);
	CREATE INDEX IF NOT EXISTS idx_accounts_status ON accounts(status);
	`
	_, err := am.db.Exec(query)
	return err
}

// CreateAccount creates a new account
func (am *AccountManager) CreateAccount(email, name, password string) (*Account, error) {
	// Check if email already exists
	var exists bool
	err := am.db.QueryRow("SELECT EXISTS(SELECT 1 FROM accounts WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Generate API key
	apiKey, err := generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("generate api key: %w", err)
	}

	// Generate ID
	id := generateID()

	account := &Account{
		ID:           id,
		Email:        email,
		Name:         name,
		PasswordHash: string(hash),
		Plan:         PlanFree,
		Status:       StatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		APIKey:       apiKey,
	}

	query := `
	INSERT INTO accounts (id, email, name, password_hash, plan, status, created_at, updated_at, api_key)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = am.db.Exec(query,
		account.ID,
		account.Email,
		account.Name,
		account.PasswordHash,
		account.Plan,
		account.Status,
		account.CreatedAt,
		account.UpdatedAt,
		account.APIKey,
	)
	if err != nil {
		return nil, fmt.Errorf("insert account: %w", err)
	}

	return account, nil
}

// Authenticate verifies email/password and returns the account
func (am *AccountManager) Authenticate(email, password string) (*Account, error) {
	account, err := am.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if account.Status != StatusActive {
		return nil, errors.New("account is not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Update last login
	am.db.Exec("UPDATE accounts SET last_login_at = ? WHERE id = ?", time.Now(), account.ID)

	return account, nil
}

// GetByID retrieves an account by ID
func (am *AccountManager) GetByID(id string) (*Account, error) {
	account := &Account{}
	query := `
	SELECT id, email, name, password_hash, plan, status, created_at, updated_at, 
	       last_login_at, api_key, container_id, metadata
	FROM accounts WHERE id = ?
	`
	err := am.db.QueryRow(query, id).Scan(
		&account.ID,
		&account.Email,
		&account.Name,
		&account.PasswordHash,
		&account.Plan,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
		&account.APIKey,
		&account.ContainerID,
		&account.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetByEmail retrieves an account by email
func (am *AccountManager) GetByEmail(email string) (*Account, error) {
	account := &Account{}
	query := `
	SELECT id, email, name, password_hash, plan, status, created_at, updated_at, 
	       last_login_at, api_key, container_id, metadata
	FROM accounts WHERE email = ?
	`
	err := am.db.QueryRow(query, email).Scan(
		&account.ID,
		&account.Email,
		&account.Name,
		&account.PasswordHash,
		&account.Plan,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
		&account.APIKey,
		&account.ContainerID,
		&account.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetByAPIKey retrieves an account by API key
func (am *AccountManager) GetByAPIKey(apiKey string) (*Account, error) {
	account := &Account{}
	query := `
	SELECT id, email, name, password_hash, plan, status, created_at, updated_at, 
	       last_login_at, api_key, container_id, metadata
	FROM accounts WHERE api_key = ?
	`
	err := am.db.QueryRow(query, apiKey).Scan(
		&account.ID,
		&account.Email,
		&account.Name,
		&account.PasswordHash,
		&account.Plan,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLoginAt,
		&account.APIKey,
		&account.ContainerID,
		&account.Metadata,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// UpdatePlan updates the account's plan
func (am *AccountManager) UpdatePlan(id string, plan Plan) error {
	_, err := am.db.Exec(
		"UPDATE accounts SET plan = ?, updated_at = ? WHERE id = ?",
		plan, time.Now(), id,
	)
	return err
}

// UpdateContainerID updates the account's container ID
func (am *AccountManager) UpdateContainerID(id, containerID string) error {
	_, err := am.db.Exec(
		"UPDATE accounts SET container_id = ?, updated_at = ? WHERE id = ?",
		containerID, time.Now(), id,
	)
	return err
}

// UpdateStatus updates the account's status
func (am *AccountManager) UpdateStatus(id string, status Status) error {
	_, err := am.db.Exec(
		"UPDATE accounts SET status = ?, updated_at = ? WHERE id = ?",
		status, time.Now(), id,
	)
	return err
}

// ListAccounts lists all accounts with pagination
func (am *AccountManager) ListAccounts(offset, limit int) ([]*Account, error) {
	query := `
	SELECT id, email, name, password_hash, plan, status, created_at, updated_at, 
	       last_login_at, api_key, container_id, metadata
	FROM accounts
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`
	rows, err := am.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account
	for rows.Next() {
		account := &Account{}
		err := rows.Scan(
			&account.ID,
			&account.Email,
			&account.Name,
			&account.PasswordHash,
			&account.Plan,
			&account.Status,
			&account.CreatedAt,
			&account.UpdatedAt,
			&account.LastLoginAt,
			&account.APIKey,
			&account.ContainerID,
			&account.Metadata,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

// CountAccounts returns the total number of accounts
func (am *AccountManager) CountAccounts() (int, error) {
	var count int
	err := am.db.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&count)
	return count, err
}

// generateAPIKey generates a random API key
func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "hn_" + hex.EncodeToString(bytes), nil
}

// generateID generates a random ID
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
