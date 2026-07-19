package sessionimport

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/glebarez/go-sqlite"

	"github.com/MDMAtk/TormentNexus/internal/database")

type ImportedSession struct {
	ID                string `json:"id"`
	SourceTool        string `json:"sourceTool"`
	SourcePath        string `json:"sourcePath"`
	ExternalSessionID string `json:"externalSessionId"`
	Title             string `json:"title"`
	SessionFormat     string `json:"sessionFormat"`
	Transcript        string `json:"transcript"`
}

func ImportSession(dbPath string, session ImportedSession) error {
	db, err := database.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open tormentnexus.db: %w", err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UnixMilli()

	// Calculate a unique transcript_hash
	hashInput := session.Transcript + "_" + session.ID + "_" + session.SourcePath
	hasher := sha256.New()
	hasher.Write([]byte(hashInput))
	tHash := hex.EncodeToString(hasher.Sum(nil))

	// Check if already exists based on SourcePath
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM imported_sessions WHERE source_path = ?)", session.SourcePath).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		_, err = tx.Exec(`
			UPDATE imported_sessions 
			SET transcript = ?, updated_at = ?, transcript_hash = ?
			WHERE source_path = ?
		`, session.Transcript, now, tHash, session.SourcePath)
	} else {
		newUUID := uuid.New().String()
		_, err = tx.Exec(`
			INSERT INTO imported_sessions (
				uuid, source_tool, source_path, external_session_id, title, 
				session_format, transcript, discovered_at, imported_at, created_at, updated_at,
				transcript_hash, normalized_session
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, newUUID, session.SourceTool, session.SourcePath, session.ExternalSessionID, session.Title,
			session.SessionFormat, session.Transcript, now, now, now, now, tHash, "{}")
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}
