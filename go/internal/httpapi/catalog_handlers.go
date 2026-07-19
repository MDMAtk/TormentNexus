package httpapi

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// CatalogResult represents a single catalog entry
type CatalogResult struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Description string   `json:"description,omitempty"`
	Source      string   `json:"source"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags,omitempty"`
}

// CatalogSearchResponse is the response from catalog search
type CatalogSearchResponse struct {
	Query   string          `json:"query"`
	Total   int             `json:"total"`
	Results []CatalogResult `json:"results"`
	Sources map[string]int  `json:"sources"`
}

// handleBacklogSearch handles GET /api/catalog/search?q=...&category=...&limit=...
func (s *Server) handleBacklogSearch(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	category := strings.TrimSpace(r.URL.Query().Get("category"))
	limit := 50

	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if limit > 200 {
		limit = 200
	}
	if limit < 1 {
		limit = 50
	}

	// Try both catalog.db locations
	dbPaths := []string{
		filepath.Join(s.cfg.WorkspaceRoot, "catalog.db"),
		filepath.Join(s.cfg.WorkspaceRoot, "catalog.db"),
	}

	var db *sql.DB
	var err error
	for _, p := range dbPaths {
		if _, statErr := os.Stat(p); statErr == nil {
			db, err = sql.Open("sqlite", p)
			if err == nil {
				break
			}
		}
	}

	if db == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"error": "catalog database not available",
		})
		return
	}
	defer db.Close()

	// Build query
	var results []CatalogResult
	sources := make(map[string]int)
	total := 0

	if query == "" && category == "" {
		// Return top entries by source
		rows, qerr := db.Query(`
			SELECT title, url, description, source 
			FROM links_backlog 
			ORDER BY source, title 
			LIMIT ?
		`, limit)
		if qerr != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": qerr.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var r CatalogResult
			if err := rows.Scan(&r.Title, &r.URL, &r.Description, &r.Source); err == nil {
				r.Category = categorizeSource(r.Source)
				results = append(results, r)
				sources[r.Source]++
				total++
			}
		}
	} else if query != "" {
		// Search by title and description
		searchTerm := "%" + strings.ToLower(query) + "%"
		rows, qerr := db.Query(`
			SELECT title, url, description, source 
			FROM links_backlog 
			WHERE lower(title) LIKE ? OR lower(description) LIKE ?
			ORDER BY 
				CASE WHEN lower(title) LIKE ? THEN 0 ELSE 1 END,
				title
			LIMIT ?
		`, searchTerm, searchTerm, searchTerm, limit)
		if qerr != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": qerr.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var r CatalogResult
			if err := rows.Scan(&r.Title, &r.URL, &r.Description, &r.Source); err == nil {
				r.Category = categorizeSource(r.Source)
				if category == "" || r.Category == category {
					results = append(results, r)
					sources[r.Source]++
					total++
				}
			}
		}
	} else {
		// Filter by category only
		catSources := sourcesForCategory(category)
		if len(catSources) == 0 {
			writeJSON(w, http.StatusOK, CatalogSearchResponse{
				Query:   query,
				Total:   0,
				Results: []CatalogResult{},
				Sources: map[string]int{},
			})
			return
		}

		placeholders := make([]string, len(catSources))
		args := make([]any, len(catSources)+1)
		for i, src := range catSources {
			placeholders[i] = "?"
			args[i] = src
		}
		args[len(catSources)] = limit

		rows, qerr := db.Query(fmt.Sprintf(`
			SELECT title, url, description, source 
			FROM links_backlog 
			WHERE source IN (%s)
			ORDER BY title
			LIMIT ?
		`, strings.Join(placeholders, ",")), args...)
		if qerr != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": qerr.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var r CatalogResult
			if err := rows.Scan(&r.Title, &r.URL, &r.Description, &r.Source); err == nil {
				r.Category = categorizeSource(r.Source)
				results = append(results, r)
				sources[r.Source]++
				total++
			}
		}
	}

	if results == nil {
		results = []CatalogResult{}
	}

	writeJSON(w, http.StatusOK, CatalogSearchResponse{
		Query:   query,
		Total:   total,
		Results: results,
		Sources: sources,
	})
}

// handleBacklogStats handles GET /api/catalog/stats
func (s *Server) handleBacklogStats(w http.ResponseWriter, r *http.Request) {
	dbPaths := []string{
		filepath.Join(s.cfg.WorkspaceRoot, "catalog.db"),
		filepath.Join(s.cfg.WorkspaceRoot, "catalog.db"),
	}

	var db *sql.DB
	var err error
	for _, p := range dbPaths {
		if _, statErr := os.Stat(p); statErr == nil {
			db, err = sql.Open("sqlite", p)
			if err == nil {
				break
			}
		}
	}

	if db == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{"error": "catalog database not available"})
		return
	}
	defer db.Close()

	// Get counts by source
	rows, err := db.Query("SELECT source, count(*) FROM links_backlog GROUP BY source ORDER BY count(*) DESC")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	defer rows.Close()

	bySource := make(map[string]int)
	byCategory := make(map[string]int)
	total := 0

	for rows.Next() {
		var source string
		var count int
		if err := rows.Scan(&source, &count); err == nil {
			bySource[source] = count
			cat := categorizeSource(source)
			byCategory[cat] += count
			total += count
		}
	}

	// Get skills count
	var skillsCount int
	skillDB, err := sql.Open("sqlite", filepath.Join(s.cfg.ConfigDir, "catalog.db"))
	if err == nil {
		skillDB.QueryRow("SELECT count(*) FROM published_skills").Scan(&skillsCount)
		skillDB.Close()
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"total":      total,
		"bySource":   bySource,
		"byCategory": byCategory,
		"skills":     skillsCount,
		"goHandlers": 5668,
	})
}

// handleBacklogCategories handles GET /api/catalog/categories
func (s *Server) handleBacklogCategories(w http.ResponseWriter, r *http.Request) {
	categories := map[string]any{
		"mcp_server": map[string]any{
			"label":       "MCP Servers",
			"description": "Model Context Protocol server implementations",
			"count":       19464,
		},
		"ai_dev_tool": map[string]any{
			"label":       "AI/Dev Tools",
			"description": "AI development tools and SDKs",
			"count":       5582,
		},
		"mcp_language": map[string]any{
			"label":       "MCP by Language",
			"description": "MCP servers by programming language",
			"count":       249,
		},
		"mcp_category": map[string]any{
			"label":       "MCP by Category",
			"description": "MCP servers by integration category",
			"count":       367,
		},
		"prompt": map[string]any{
			"label":       "Prompts & Templates",
			"description": "System prompts, prompt patterns, and templates",
			"count":       295,
		},
		"agent": map[string]any{
			"label":       "Agent Frameworks",
			"description": "AI agent frameworks and tools",
			"count":       223,
		},
		"skill": map[string]any{
			"label":       "Skills",
			"description": "Reusable skill modules for AI agents",
			"count":       5441,
		},
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"categories": categories,
	})
}

// categorizeSource maps a source string to a category
func categorizeSource(source string) string {
	mcp := map[string]bool{
		"awesome-mcp-servers": true, "awesome-mcp-fork": true, "npm": true,
		"npm-expanded": true, "npm-tool": true, "npm-mcp-sdk": true,
		"github-search": true, "github-mcp-name": true, "github-mcp-desc": true,
		"github-specific": true, "github-topic": true, "github-topic-mcp": true,
		"github-topic-mcp-proto": true, "pypi": true, "crates-io": true,
		"docker-hub": true, "glama-html": true, "glama-mock": true,
		"mcp-get": true, "multi-source-scrape": true, "awesome-fork-v2": true,
	}
	prompt := map[string]bool{
		"prompt-repo": true, "github-prompt-search": true, "fabric-patterns": true,
		"awesome-ai-system-prompts": true, "chatgpt-system-prompts": true,
		"system-prompts-leaks": true, "big-prompt-library": true,
		"meigen-ai-design-mcp": true, "system-prompts-ai-tools": true,
		"claude-code-prompts": true, "claude-code-prompts-v2": true,
		"leaked-system-prompts": true,
	}
	agent := map[string]bool{
		"agent-framework": true, "awesome-ai-agent": true,
	}

	if mcp[source] {
		return "mcp_server"
	}
	if prompt[source] {
		return "prompt"
	}
	if agent[source] {
		return "agent"
	}
	if strings.HasPrefix(source, "github-lang-") {
		return "mcp_language"
	}
	if strings.HasPrefix(source, "github-cat-") {
		return "mcp_category"
	}
	return "ai_dev_tool"
}

// sourcesForCategory returns all sources that belong to a category
func sourcesForCategory(category string) []string {
	switch category {
	case "mcp_server":
		return []string{
			"awesome-mcp-servers", "awesome-mcp-fork", "npm", "npm-expanded",
			"npm-tool", "npm-mcp-sdk", "github-search", "github-mcp-name",
			"github-mcp-desc", "github-specific", "github-topic", "github-topic-mcp",
			"github-topic-mcp-proto", "pypi", "crates-io", "docker-hub",
			"glama-html", "glama-mock", "mcp-get", "multi-source-scrape",
			"awesome-fork-v2",
		}
	case "prompt":
		return []string{
			"prompt-repo", "github-prompt-search", "fabric-patterns",
			"awesome-ai-system-prompts", "chatgpt-system-prompts",
			"system-prompts-leaks", "big-prompt-library", "meigen-ai-design-mcp",
			"system-prompts-ai-tools", "claude-code-prompts", "claude-code-prompts-v2",
			"leaked-system-prompts",
		}
	case "agent":
		return []string{"agent-framework", "awesome-ai-agent"}
	case "mcp_language":
		return []string{
			"github-lang-python", "github-lang-typescript", "github-lang-go",
			"github-lang-rust", "github-lang-java", "github-lang-ruby",
			"github-lang-swift", "github-lang-c++", "github-lang-c#",
		}
	case "mcp_category":
		return []string{
			"github-cat-database", "github-cat-filesystem", "github-cat-browser",
			"github-cat-search", "github-cat-github", "github-cat-slack",
			"github-cat-discord", "github-cat-notion", "github-cat-docker",
			"github-cat-aws", "github-cat-stripe", "github-cat-email",
			"github-cat-calendar", "github-cat-monitoring", "github-cat-testing",
			"github-cat-security", "github-cat-payment", "github-cat-crm",
			"github-cat-cms", "github-cat-ecommerce", "github-cat-social",
			"github-cat-analytics", "github-cat-logging", "github-cat-auth",
		}
	case "ai_dev_tool":
		return []string{
			"awesome-general", "npm-ai-sdk", "npm-tool-v2", "github-ai-tool",
			"npm-adjacent", "awesome-domain", "awesome-ai-v2",
		}
	default:
		return nil
	}
}
