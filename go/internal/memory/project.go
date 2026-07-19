package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ProjectDB struct {
	WorkspaceRoot string
}

func NewProjectDB(workspaceRoot string) *ProjectDB {
	return &ProjectDB{WorkspaceRoot: workspaceRoot}
}

func (db *ProjectDB) SyncMemDB() (map[string]any, error) {
	memdbPath := filepath.Join(db.WorkspaceRoot, ".memdb")

	// Create it if it doesn't exist
	if _, err := os.Stat(memdbPath); os.IsNotExist(err) {
		err := os.WriteFile(memdbPath, []byte("{}"), 0644)
		if err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(memdbPath)
	if err != nil {
		return nil, err
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		parsed = map[string]any{}
	}

	return parsed, nil
}
