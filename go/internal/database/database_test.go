package database

import (
	"os"
	"testing"
)

func TestRewriteQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Placeholders",
			input:    "SELECT * FROM users WHERE name = ? AND email = ?",
			expected: "SELECT * FROM users WHERE name = $1 AND email = $2",
		},
		{
			name:     "Literal String Placeholders",
			input:    "SELECT * FROM users WHERE name = 'John ?' AND email = ?",
			expected: "SELECT * FROM users WHERE name = 'John ?' AND email = $1",
		},
		{
			name:     "Datetime replace",
			input:    "INSERT INTO logs (time) VALUES (datetime('now'))",
			expected: "INSERT INTO logs (time) VALUES (CURRENT_TIMESTAMP)",
		},
		{
			name:     "FTS Table rewrite",
			input:    "CREATE VIRTUAL TABLE IF NOT EXISTS items USING fts5(id, content)",
			expected: "CREATE TABLE IF NOT EXISTS items (id, content)",
		},
		{
			name:     "Insert Ignore rewrite",
			input:    "INSERT OR IGNORE INTO items (id, val) VALUES (?, ?)",
			expected: "INSERT INTO items (id, val) VALUES ($1, $2) ON CONFLICT DO NOTHING",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := RewriteQuery(tc.input)
			if actual != tc.expected {
				t.Errorf("Expected: %q\nGot: %q", tc.expected, actual)
			}
		})
	}
}

func TestOpenFallback(t *testing.T) {
	os.Unsetenv("DATABASE_URL")
	db, err := Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE test (id INT)")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
}
