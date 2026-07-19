package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MDMAtk/TormentNexus/internal/controlplane"
	"github.com/MDMAtk/TormentNexus/internal/memorystore"
)

// TN-native memory tool handlers.
// These use GlobalVectorStore directly — same backend as the pi extension's tn_memory_store tool
// and the /api/memory/add HTTP endpoint.

// HandleAddMemory stores a memory in the L2 vault with full structured fields.
// Supports: content, tags, category, importance, session_id.
func HandleAddMemory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	content, _ := getString(args, "content")
	if content == "" {
		return err("content is required")
	}

	if GlobalVectorStore == nil {
		return err("vector store not initialized")
	}

	tags, _ := getString(args, "tags")
	category, _ := getString(args, "category")
	importance, _ := getFloat(args, "importance")
	sessionID, _ := getString(args, "session_id")

	if category == "" {
		category = "general"
	}
	if importance <= 0 {
		importance = 0.5
	}
	if sessionID == "" {
		sessionID = "manual"
	}

	// Wrap content with metadata JSON (matching pi extension format)
	wrapped := map[string]interface{}{
		"content":   content,
		"tags":      parseTags(tags),
		"category":  category,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	wrappedJSON, _ := json.Marshal(wrapped)

	record := controlplane.L2VaultRecord{
		ID:             fmt.Sprintf("mem-%s-%d", strings.ReplaceAll(category, " ", "_"), time.Now().UnixNano()),
		SessionID:      sessionID,
		Type:           controlplane.MemoryLongTerm,
		Kind:           "fact",
		Category:       category,
		Tags:           tags,
		Content:        string(wrappedJSON),
		Importance:     importance,
		HeatScore:      50.0,
		LastAccessedAt: time.Now(),
		CreatedAt:      time.Now(),
	}

	if storeErr := GlobalVectorStore.Commit(ctx, record); storeErr != nil {
		return err(fmt.Sprintf("store failed: %v", storeErr))
	}

	return ok(fmt.Sprintf("Memory stored (category: %s, tags: %s)", category, tags))
}

// HandleSearchMemory searches the L2 vault by keyword, with optional tag/category filter.
// Supports: query, tag, category, limit, offset, includeCold.
func HandleSearchMemory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query")
	if query == "" {
		// Return recent memories instead of failing
		query = ""
	}

	if GlobalVectorStore == nil {
		return err("vector store not initialized")
	}

	limit, _ := getInt(args, "limit")
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, _ := getInt(args, "offset")
	if offset < 0 {
		offset = 0
	}
	tagFilter, _ := getString(args, "tag")
	categoryFilter, _ := getString(args, "category")
	includeCold, _ := getBool(args, "include_cold")

	// Use FTS search via the memory store
	db := GlobalVectorStore.DB()
	ftsSearcher, ftsErr := memorystore.NewFTSMemorySearch(db)
	if ftsErr != nil {
		return err(fmt.Sprintf("search init failed: %v", ftsErr))
	}

	results, searchErr := ftsSearcher.Search(ctx, query, includeCold, limit, offset)
	if searchErr != nil {
		return err(fmt.Sprintf("search failed: %v", searchErr))
	}

	// Apply tag/category filters post-hoc
	filtered := make([]memorystore.FTSMemorySearchResult, 0)
	for _, r := range results.Results {
		if tagFilter != "" && !strings.Contains(r.Record.Tags, tagFilter) {
			continue
		}
		if categoryFilter != "" && r.Record.Category != categoryFilter {
			continue
		}
		filtered = append(filtered, r)
	}

	// Return formatted results
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d memories (showing %d)\n", results.Total, len(filtered)))
	for i, r := range filtered {
		rec := r.Record
		// Try to extract content from wrapped JSON
		displayContent := rec.Content
		var wrapped map[string]interface{}
		if err := json.Unmarshal([]byte(rec.Content), &wrapped); err == nil {
			if c, ok := wrapped["content"].(string); ok {
				displayContent = c
			}
		}
		contentStr := displayContent
		if len(contentStr) > 120 {
			contentStr = contentStr[:120] + "..."
		}
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, rec.Category, contentStr))
		if rec.Tags != "" {
			sb.WriteString(fmt.Sprintf("   Tags: %s\n", rec.Tags))
		}
	}

	return ok(sb.String())
}

// HandleDeleteMemory removes a memory from the L2 vault by ID.
func HandleDeleteMemory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getString(args, "id")
	if id == "" {
		return err("id is required")
	}

	if GlobalVectorStore == nil {
		return err("vector store not initialized")
	}

	db := GlobalVectorStore.DB()
	_, delErr := db.ExecContext(ctx, `DELETE FROM l2_vault WHERE id = ?`, id)
	if delErr != nil {
		return err(fmt.Sprintf("delete failed: %v", delErr))
	}
	_, _ = db.ExecContext(ctx, `DELETE FROM l2_memory_fts WHERE memory_id = ?`, id)
	_, _ = db.ExecContext(ctx, `DELETE FROM vec_l2_vault WHERE id = ?`, id)

	return ok(fmt.Sprintf("Memory %s deleted", id))
}

// HandleMemoryStats returns counts across all memory tiers.
func HandleMemoryStats(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	if GlobalVectorStore == nil {
		return err("vector store not initialized")
	}

	db := GlobalVectorStore.DB()
	var vaultCount int
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM l2_vault`).Scan(&vaultCount)

	var ftsCount int
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM l2_memory_fts`).Scan(&ftsCount)

	var srCount int
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM spaced_repetition_metadata`).Scan(&srCount)

	return ok(fmt.Sprintf("L2 Vault: %d | FTS indexed: %d | Spaced repetition: %d", vaultCount, ftsCount, srCount))
}

func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return nil
	}
	parts := strings.Split(tagsStr, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			result = append(result, t)
		}
	}
	return result
}
