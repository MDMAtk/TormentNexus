package catalogingestor

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/glebarez/go-sqlite"

	"github.com/MDMAtk/TormentNexus/internal/database")

func SyncRegisteredToolsToCatalog(workspaceRoot string, toolNames []string) error {
	dbPath := filepath.Join(workspaceRoot, "catalog.db")
	db, err := database.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open catalog.db: %w", err)
	}
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA busy_timeout=5000")
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	stmt, err := tx.Prepare(`
		INSERT INTO published_mcp_servers (
			uuid, canonical_id, display_name, description, tags, categories, transport, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(canonical_id) DO UPDATE SET
			display_name = excluded.display_name,
			description = excluded.description,
			updated_at = excluded.updated_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, tool := range toolNames {
		uid := uuid.New().String()
		canonicalID := "native-tool:" + tool
		displayName := "Go Native: " + tool
		desc := fmt.Sprintf("Built-in, high-performance Go-native TormentNexus tool: %s", tool)
		tags := "[\"native\", \"always-on\", \"performance\"]"
		categories := "[\"core\", \"utility\"]"
		transport := "in-process"
		status := "active"

		_, err := stmt.Exec(
			uid, canonicalID, displayName, desc, tags, categories, transport, status, now, now,
		)
		if err != nil {
			return fmt.Errorf("failed to sync tool %s: %w", tool, err)
		}
	}

	return tx.Commit()
}
