package memorystore

import (
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/MDMAtk/TormentNexus/internal/controlplane"
)

// SpacedRepetitionMetadata tracks Leitner / SuperMemo SM-2 state for memory entries
type SpacedRepetitionMetadata struct {
	MemoryID     string    `json:"memory_id"`
	Repetitions  int       `json:"repetitions"`
	Interval     int       `json:"interval"` // in days
	EaseFactor   float64   `json:"ease_factor"`
	NextReviewAt time.Time `json:"next_review_at"`
}

// InitSpacedRepetition creates the metadata table if it doesn't exist
func InitSpacedRepetition(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS spaced_repetition_metadata (
		memory_id TEXT PRIMARY KEY,
		repetitions INTEGER DEFAULT 0,
		interval INTEGER DEFAULT 0,
		ease_factor REAL DEFAULT 2.5,
		next_review_at TEXT
	);`
	_, err := db.Exec(query)
	return err
}

// ReviewMemory updates the spaced repetition metrics for a memory entry using SM-2
func ReviewMemory(db *sql.DB, memoryID string, quality int) error {
	if quality < 0 || quality > 5 {
		return errors.New("review quality must be between 0 and 5")
	}

	// Fetch current metadata or initialize default
	var meta SpacedRepetitionMetadata
	meta.MemoryID = memoryID
	meta.EaseFactor = 2.5
	meta.Interval = 0
	meta.Repetitions = 0
	meta.NextReviewAt = time.Now()

	var nextReviewStr string
	querySelect := `SELECT repetitions, interval, ease_factor, next_review_at FROM spaced_repetition_metadata WHERE memory_id = ?`
	err := db.QueryRow(querySelect, memoryID).Scan(&meta.Repetitions, &meta.Interval, &meta.EaseFactor, &nextReviewStr)
	if err == nil && nextReviewStr != "" {
		if t, err := time.Parse(time.RFC3339, nextReviewStr); err == nil {
			meta.NextReviewAt = t
		}
	}

	// SM-2 Algorithm
	if quality < 3 {
		// Failure: reset repetitions, interval to 1 day, keep ease factor
		meta.Repetitions = 0
		meta.Interval = 1
	} else {
		// Success: recalculate ease factor
		q := float64(quality)
		meta.EaseFactor = meta.EaseFactor + (0.1 - (5-q)*(0.08+(5-q)*0.02))
		if meta.EaseFactor < 1.3 {
			meta.EaseFactor = 1.3
		}

		if meta.Repetitions == 0 {
			meta.Interval = 1
		} else if meta.Repetitions == 1 {
			meta.Interval = 6
		} else {
			meta.Interval = int(math.Round(float64(meta.Interval) * meta.EaseFactor))
		}
		meta.Repetitions++
	}

	meta.NextReviewAt = time.Now().Add(time.Duration(meta.Interval) * 24 * time.Hour)

	queryUpsert := `
	INSERT INTO spaced_repetition_metadata (memory_id, repetitions, interval, ease_factor, next_review_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(memory_id) DO UPDATE SET
		repetitions = excluded.repetitions,
		interval = excluded.interval,
		ease_factor = excluded.ease_factor,
		next_review_at = excluded.next_review_at`

	_, err = db.Exec(queryUpsert, meta.MemoryID, meta.Repetitions, meta.Interval, meta.EaseFactor, meta.NextReviewAt.Format(time.RFC3339))
	return err
}

// GetDueMemories returns all memory IDs that are due or overdue for a review
func GetDueMemories(db *sql.DB) ([]string, error) {
	// First ensure we have metadata initialized for all L2 vault records.
	// To do this, we find any memory records in l2_vault that don't exist in spaced_repetition_metadata.
	// We can select memory_id from l2_vault directly.
	// Let's check what the L2 vault table name is in vector_sqlite.go.
	// Typically it's "l2_vault" or "vault". Let's run a query that selects due items.
	
	query := `
	SELECT id FROM l2_vault
	WHERE id NOT IN (SELECT memory_id FROM spaced_repetition_metadata)
	UNION
	SELECT memory_id FROM spaced_repetition_metadata
	WHERE datetime(next_review_at) <= datetime(?)
	`
	rows, err := db.Query(query, time.Now().Format(time.RFC3339))
	if err != nil {
		// Fallback if l2_vault does not exist yet
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

// ReviewMemory is a VectorStore wrapper around ReviewMemory function
func (vs *VectorStore) ReviewMemory(memoryID string, quality int) error {
	return ReviewMemory(vs.db, memoryID, quality)
}

// GetDueMemories is a VectorStore wrapper around GetDueMemories function
func (vs *VectorStore) GetDueMemories() ([]string, error) {
	return GetDueMemories(vs.db)
}

// GetDueMemoriesRecords returns the full L2VaultRecord of all memory entries that are due or overdue for review
func GetDueMemoriesRecords(db *sql.DB) ([]controlplane.L2VaultRecord, error) {
	query := `
	SELECT id, session_id, memory_type, memory_kind, category, tags, source_url, content, importance, heat_score, last_accessed_at, created_at
	FROM l2_vault
	WHERE id NOT IN (SELECT memory_id FROM spaced_repetition_metadata)
	UNION
	SELECT l.id, l.session_id, l.memory_type, l.memory_kind, l.category, l.tags, l.source_url, l.content, l.importance, l.heat_score, l.last_accessed_at, l.created_at
	FROM l2_vault l
	JOIN spaced_repetition_metadata s ON l.id = s.memory_id
	WHERE datetime(s.next_review_at) <= datetime(?)
	`
	rows, err := db.Query(query, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []controlplane.L2VaultRecord
	for rows.Next() {
		var r controlplane.L2VaultRecord
		var lastAccessedStr, createdStr string
		var memoryTypeStr string
		err := rows.Scan(
			&r.ID, &r.SessionID, &memoryTypeStr, &r.Kind, &r.Category, &r.Tags, &r.SourceURL, &r.Content, &r.Importance, &r.HeatScore, &lastAccessedStr, &createdStr,
		)
		if err != nil {
			continue
		}
		r.Type = controlplane.MemoryType(memoryTypeStr)
		if t, err := time.Parse(time.RFC3339, lastAccessedStr); err == nil {
			r.LastAccessedAt = t
		}
		if t, err := time.Parse(time.RFC3339, createdStr); err == nil {
			r.CreatedAt = t
		}
		records = append(records, r)
	}
	return records, nil
}

// GetDueMemoriesRecords is a VectorStore wrapper around GetDueMemoriesRecords function
func (vs *VectorStore) GetDueMemoriesRecords() ([]controlplane.L2VaultRecord, error) {
	return GetDueMemoriesRecords(vs.db)
}


