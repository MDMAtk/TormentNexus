package memorystore

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/MDMAtk/TormentNexus/internal/controlplane"

	"github.com/MDMAtk/TormentNexus/internal/database")

// ─── Project .memdb Schema ────────────────────────────────────────────

// ProjectMemorySchema is the DDL for a per-project .memdb file.
// Minimal schema: just the fields that matter for portable project memory.
const ProjectMemorySchema = `
CREATE TABLE IF NOT EXISTS memories (
    id          TEXT PRIMARY KEY,
    session_id  TEXT NOT NULL DEFAULT '',
    kind        TEXT NOT NULL DEFAULT 'fact',
    category    TEXT NOT NULL DEFAULT 'general',
    tags        TEXT NOT NULL DEFAULT '',
    content     TEXT NOT NULL,
    importance  REAL NOT NULL DEFAULT 0.5,
    created_at  TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX IF NOT EXISTS idx_memories_category ON memories(category);
CREATE INDEX IF NOT EXISTS idx_memories_kind ON memories(kind);
`

// ProjectDB manages a per-project .memdb SQLite file.
// Each file is a portable, git-trackable collection of project-specific memories.
type ProjectDB struct {
	path       string
	db         *sql.DB
	mu         sync.RWMutex
	stmtInsert *sql.Stmt
	stmtSearch *sql.Stmt
	stmtList   *sql.Stmt
}

// OpenProjectDB opens or creates a .memdb file for a project.
func OpenProjectDB(projectPath string) (*ProjectDB, error) {
	if err := os.MkdirAll(filepath.Dir(projectPath), 0755); err != nil {
		return nil, fmt.Errorf("project db mkdir: %w", err)
	}

	db, err := database.Open("sqlite", projectPath)
	if err != nil {
		return nil, fmt.Errorf("project db open: %w", err)
	}

	// WAL disabled intentionally — .memdb files should be clean for git tracking
	if _, err := db.Exec("PRAGMA journal_mode=DELETE"); err != nil {
		db.Close()
		return nil, fmt.Errorf("project db pragma: %w", err)
	}
	if _, err := db.Exec("PRAGMA synchronous=NORMAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("project db sync: %w", err)
	}

	if _, err := db.Exec(ProjectMemorySchema); err != nil {
		db.Close()
		return nil, fmt.Errorf("project db schema: %w", err)
	}

	insertStmt, err := db.Prepare(`INSERT OR IGNORE INTO memories (id, session_id, kind, category, tags, content, importance, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("project db prepare insert: %w", err)
	}

	searchStmt, err := db.Prepare(`SELECT id, session_id, kind, category, tags, content, importance, created_at FROM memories WHERE content LIKE ? OR tags LIKE ? OR category = ? ORDER BY importance DESC LIMIT ?`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("project db prepare search: %w", err)
	}

	listStmt, err := db.Prepare(`SELECT id, session_id, kind, category, tags, content, importance, created_at FROM memories ORDER BY created_at DESC LIMIT ?`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("project db prepare list: %w", err)
	}

	return &ProjectDB{
		path:       projectPath,
		db:         db,
		stmtInsert: insertStmt,
		stmtSearch: searchStmt,
		stmtList:   listStmt,
	}, nil
}

func (p *ProjectDB) Close() error {
	if p.stmtInsert != nil {
		p.stmtInsert.Close()
	}
	if p.stmtSearch != nil {
		p.stmtSearch.Close()
	}
	if p.stmtList != nil {
		p.stmtList.Close()
	}
	return p.db.Close()
}

func (p *ProjectDB) Path() string { return p.path }

// Store writes a memory to this project's .memdb.
func (p *ProjectDB) Store(ctx context.Context, entry controlplane.L2VaultRecord) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	created := entry.CreatedAt
	if created.IsZero() {
		created = time.Now()
	}
	_, err := p.stmtInsert.ExecContext(ctx, entry.ID, entry.SessionID, entry.Kind, entry.Category, entry.Tags, entry.Content, entry.Importance, created.UTC().Format(time.RFC3339))
	return err
}

// Search queries this project's memories by keyword.
func (p *ProjectDB) Search(ctx context.Context, query string, limit int) ([]controlplane.L2VaultRecord, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	like := "%" + query + "%"
	rows, err := p.stmtSearch.QueryContext(ctx, like, like, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProjectRows(rows)
}

// List returns the most recent memories from this project.
func (p *ProjectDB) List(ctx context.Context, limit int) ([]controlplane.L2VaultRecord, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rows, err := p.stmtList.QueryContext(ctx, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProjectRows(rows)
}

// Count returns the number of memories in this project's .memdb.
func (p *ProjectDB) Count(ctx context.Context) (int, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var count int
	err := p.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM memories").Scan(&count)
	return count, err
}

func scanProjectRows(rows *sql.Rows) ([]controlplane.L2VaultRecord, error) {
	var results []controlplane.L2VaultRecord
	for rows.Next() {
		var r controlplane.L2VaultRecord
		var createdStr string
		if err := rows.Scan(&r.ID, &r.SessionID, &r.Kind, &r.Category, &r.Tags, &r.Content, &r.Importance, &createdStr); err != nil {
			continue
		}
		r.Type = controlplane.MemoryLongTerm
		if t, err := time.Parse(time.RFC3339, createdStr); err == nil {
			r.CreatedAt = t
		}
		r.LastAccessedAt = time.Now()
		r.HeatScore = 50.0
		results = append(results, r)
	}
	return results, nil
}

// ─── Workspace Scanner ────────────────────────────────────────────────

// FindProjectMemDBs scans a workspace root for all .memdb files.
// Looks in: workspace root, and first-level subdirectories.
func FindProjectMemDBs(workspaceRoot string) ([]string, error) {
	var paths []string

	// Check workspace root
	root := filepath.Join(workspaceRoot, ".memdb")
	if _, err := os.Stat(root); err == nil {
		paths = append(paths, root)
	}

	// Check first-level subdirectories
	entries, err := os.ReadDir(workspaceRoot)
	if err != nil {
		return paths, nil // workspace might not exist
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Skip hidden dirs, node_modules, .git, etc.
		name := entry.Name()
		if strings.HasPrefix(name, ".") || name == "node_modules" || name == ".git" {
			continue
		}
		memdbPath := filepath.Join(workspaceRoot, name, ".memdb")
		if _, err := os.Stat(memdbPath); err == nil {
			paths = append(paths, memdbPath)
		}
	}

	return paths, nil
}

// ImportProjectMemDB reads all memories from a .memdb file and
// commits them into the global VectorStore. Deduplicates by ID.
func ImportProjectMemDB(ctx context.Context, memdbPath string, globalVS *VectorStore) (int, error) {
	pdb, err := OpenProjectDB(memdbPath)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", memdbPath, err)
	}
	defer pdb.Close()

	records, err := pdb.List(ctx, 10000)
	if err != nil {
		return 0, fmt.Errorf("list %s: %w", memdbPath, err)
	}

	var imported int
	for _, r := range records {
		// Tag with project source if not already tagged
		projectName := filepath.Base(filepath.Dir(memdbPath))
		if projectName == "." || projectName == string(filepath.Separator) {
			projectName = "root"
		}
		if !strings.Contains(r.Tags, "project:") {
			if r.Tags != "" {
				r.Tags += ","
			}
			r.Tags += "project:" + projectName
		}

		if err := globalVS.Commit(ctx, r); err == nil {
			imported++
		}
	}

	return imported, nil
}

// SyncAllProjectMemDBs scans the workspace and imports all .memdb files.
func SyncAllProjectMemDBs(ctx context.Context, workspaceRoot string, globalVS *VectorStore) (int, int, error) {
	paths, err := FindProjectMemDBs(workspaceRoot)
	if err != nil {
		return 0, 0, err
	}

	var totalFiles, totalMemories int
	for _, p := range paths {
		count, err := ImportProjectMemDB(ctx, p, globalVS)
		if err != nil {
			fmt.Printf("[ProjectDB] Import error %s: %v\n", p, err)
			continue
		}
		if count > 0 {
			totalFiles++
			totalMemories += count
			fmt.Printf("[ProjectDB] Imported %d memories from %s\n", count, p)
		}
	}

	return totalFiles, totalMemories, nil
}

// RetroactivelySplitMemoriesById scans existing memories in the global vault
// for project tags and writes them into corresponding .memdb files.
// Returns (filesCreated, memoriesStored, errors).
func RetroactivelySplitMemoriesById(ctx context.Context, globalVS *VectorStore, workspaceRoot string) (int, int, error) {
	records, err := globalVS.GetAllVaultRecords(ctx, 100000)
	if err != nil {
		return 0, 0, fmt.Errorf("get all records: %w", err)
	}

	// Group by project
	projects := make(map[string][]controlplane.L2VaultRecord)
	for _, r := range records {
		// Extract project from tags
		projectName := extractProjectFromTags(r.Tags)
		if projectName == "" {
			// Try session_id
			if strings.HasPrefix(r.SessionID, "project:") {
				projectName = strings.TrimPrefix(r.SessionID, "project:")
			}
		}
		if projectName == "" {
			projectName = "_ungrouped"
		}
		projects[projectName] = append(projects[projectName], r)
	}

	var filesCreated, memoriesStored int
	for projectName, recs := range projects {
		if projectName == "_ungrouped" {
			continue
		}
		memdbPath := filepath.Join(workspaceRoot, projectName, ".memdb")
		if _, err := os.Stat(filepath.Dir(memdbPath)); os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(memdbPath), 0755); err != nil {
				fmt.Printf("[ProjectDB] Cannot create dir for %s: %v\n", projectName, err)
				continue
			}
		}

		pdb, err := OpenProjectDB(memdbPath)
		if err != nil {
			fmt.Printf("[ProjectDB] Cannot open %s: %v\n", memdbPath, err)
			continue
		}

		for _, r := range recs {
			if err := pdb.Store(ctx, r); err == nil {
				memoriesStored++
			}
		}
		pdb.Close()
		filesCreated++
		fmt.Printf("[ProjectDB] Split %d memories to %s\n", len(recs), memdbPath)
	}

	return filesCreated, memoriesStored, nil
}

func extractProjectFromTags(tags string) string {
	for _, part := range strings.Split(tags, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "project:") {
			return strings.TrimPrefix(part, "project:")
		}
	}
	return ""
}
