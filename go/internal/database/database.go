package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/lib/pq"
	_ "github.com/glebarez/go-sqlite"
)

func init() {
	sql.Register("pgcompat", &CompatDriver{})
}

// Open opens a database connection. If the environment variable DATABASE_URL is set,
// it will connect to PostgreSQL using a compatibility wrapper. Otherwise, it will connect
// using the specified driver (usually "sqlite") and dataSourceName.
func Open(driverName, dataSourceName string) (*sql.DB, error) {
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return sql.Open("pgcompat", dbURL)
	}
	return sql.Open(driverName, dataSourceName)
}

type CompatDriver struct{}

func (d *CompatDriver) Open(name string) (driver.Conn, error) {
	pqDriver := &pq.Driver{}
	conn, err := pqDriver.Open(name)
	if err != nil {
		return nil, err
	}
	return &CompatConn{conn: conn}, nil
}

type CompatConn struct {
	conn driver.Conn
}

func (c *CompatConn) Prepare(query string) (driver.Stmt, error) {
	rewritten := RewriteQuery(query)
	stmt, err := c.conn.Prepare(rewritten)
	if err != nil {
		return nil, fmt.Errorf("prepare failed: %w (rewritten query: %s)", err, rewritten)
	}
	return &CompatStmt{stmt: stmt}, nil
}

func (c *CompatConn) Close() error {
	return c.conn.Close()
}

func (c *CompatConn) Begin() (driver.Tx, error) {
	return c.conn.Begin()
}

func (c *CompatConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	rewritten := RewriteQuery(query)
	if execer, ok := c.conn.(driver.ExecerContext); ok {
		return execer.ExecContext(ctx, rewritten, args)
	}
	// Fallback to Prepare
	stmt, err := c.conn.Prepare(rewritten)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	dargs := make([]driver.Value, len(args))
	for i, v := range args {
		dargs[i] = v.Value
	}
	return stmt.Exec(dargs)
}

func (c *CompatConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	rewritten := RewriteQuery(query)
	if queryer, ok := c.conn.(driver.QueryerContext); ok {
		return queryer.QueryContext(ctx, rewritten, args)
	}
	// Fallback to Prepare
	stmt, err := c.conn.Prepare(rewritten)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	dargs := make([]driver.Value, len(args))
	for i, v := range args {
		dargs[i] = v.Value
	}
	return stmt.Query(dargs)
}

type CompatStmt struct {
	stmt driver.Stmt
}

func (s *CompatStmt) Close() error {
	return s.stmt.Close()
}

func (s *CompatStmt) NumInput() int {
	return s.stmt.NumInput()
}

func (s *CompatStmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.stmt.Exec(args)
}

func (s *CompatStmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.stmt.Query(args)
}

// RewriteQuery translates SQLite syntax to PostgreSQL syntax at runtime.
func RewriteQuery(query string) string {
	// 1. Replace placeholder ? with $1, $2, ...
	var result strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	paramIndex := 1
	n := len(query)
	for i := 0; i < n; i++ {
		ch := query[i]
		if ch == '\'' && (i == 0 || query[i-1] != '\\') {
			inSingleQuote = !inSingleQuote
			result.WriteByte(ch)
		} else if ch == '"' && (i == 0 || query[i-1] != '\\') {
			inDoubleQuote = !inDoubleQuote
			result.WriteByte(ch)
		} else if ch == '?' && !inSingleQuote && !inDoubleQuote {
			result.WriteString("$" + strconv.Itoa(paramIndex))
			paramIndex++
		} else {
			result.WriteByte(ch)
		}
	}
	q := result.String()

	// 2. Translate date/time and types:
	q = regexp.MustCompile(`(?i)datetime\('now'\)`).ReplaceAllString(q, "CURRENT_TIMESTAMP")
	q = regexp.MustCompile(`(?i)datetime\("now"\)`).ReplaceAllString(q, "CURRENT_TIMESTAMP")

	// 3. SQLite VIRTUAL TABLE USING fts5 -> Standard PostgreSQL table
	if strings.Contains(strings.ToUpper(q), "USING FTS5") {
		q = regexp.MustCompile(`(?i)CREATE\s+VIRTUAL\s+TABLE`).ReplaceAllString(q, "CREATE TABLE")
		q = regexp.MustCompile(`(?i)USING\s+fts5\s*`).ReplaceAllString(q, "")
		q = regexp.MustCompile(`(?i)\bUNINDEXED\b`).ReplaceAllString(q, "")
	}

	// 4. Ignore SQLite triggers for PostgreSQL (we use standard fallback or manual replication if needed)
	if strings.Contains(strings.ToUpper(q), "CREATE TRIGGER") {
		return "/* ignored trigger */ SELECT 1"
	}

	// 5. SQLite FTS MATCH syntax to PostgreSQL ILIKE
	q = regexp.MustCompile(`(?i)\bcontent\s+MATCH\s+(\S+)`).ReplaceAllString(q, "content ILIKE $1")

	// 6. Handle INSERT OR IGNORE and INSERT OR REPLACE compatibility
	reIgnore := regexp.MustCompile(`(?i)\bINSERT\s+OR\s+IGNORE\s+INTO\b`)
	if reIgnore.MatchString(q) {
		q = reIgnore.ReplaceAllString(q, "INSERT INTO")
		if !strings.Contains(strings.ToUpper(q), "ON CONFLICT") {
			q = q + " ON CONFLICT DO NOTHING"
		}
	}

	reReplace := regexp.MustCompile(`(?i)\bINSERT\s+OR\s+REPLACE\s+INTO\b`)
	if reReplace.MatchString(q) {
		q = reReplace.ReplaceAllString(q, "INSERT INTO")
		if !strings.Contains(strings.ToUpper(q), "ON CONFLICT") {
			lowerQ := strings.ToLower(q)
			if strings.Contains(lowerQ, "l3_cold_archive") {
				q = q + " ON CONFLICT (id) DO UPDATE SET content=EXCLUDED.content, heat_score=EXCLUDED.heat_score, archived_at=EXCLUDED.archived_at"
			} else if strings.Contains(lowerQ, "tool_embeddings") {
				q = q + " ON CONFLICT (name) DO UPDATE SET description=EXCLUDED.description, embedding=EXCLUDED.embedding"
			} else {
				q = q + " ON CONFLICT (id) DO UPDATE SET content=EXCLUDED.content"
			}
		}
	}

	return q
}
