package mcpimpl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var arxivHTTP = &http.Client{Timeout: 15 * time.Second}

// HandleSearchArxiv searches arXiv papers by keyword.
func HandleSearchArxiv(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query")
	maxResults, _ := getInt(args, "maxResults", 10)

	resp, fetchErr := arxivHTTP.Get(fmt.Sprintf("http://export.arxiv.org/api/query?search_query=all:%s&max_results=%d&sortBy=relevance", url.QueryEscape(query), maxResults))
	if fetchErr != nil {
		return err(fmt.Sprintf("arxiv query: %v", fetchErr))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	text := string(body)

	// Parse the Atom XML response into something useful
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📄 arXiv results for \"%s\":\n\n", query))

	entries := strings.Split(text, "<entry>")
	for i, entry := range entries {
		if i == 0 || i > maxResults {
			continue
		}
		title := extractXML(entry, "<title>", "</title>")
		authors := extractXML(entry, "<author>", "</author>")
		summary := extractXML(entry, "<summary>", "</summary>")
		link := extractXML(entry, "<id>", "</id>")
		published := extractXML(entry, "<published>", "</published>")

		title = cleanXML(title)
		summary = cleanXML(summary)
		// Truncate summary
		if len(summary) > 200 {
			summary = summary[:200] + "..."
		}

		sb.WriteString(fmt.Sprintf("%d. %s\n", i, title))
		if published != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", published[:10]))
		}
		if authors != "" {
			authorName := extractXML(authors, "<name>", "</name>")
			sb.WriteString(fmt.Sprintf("   %s\n", cleanXML(authorName)))
		}
		sb.WriteString(fmt.Sprintf("   %s\n", link))
		sb.WriteString(fmt.Sprintf("   %s\n\n", summary))
	}

	return ok(sb.String())
}

// HandleGetAbstract gets the abstract of a specific arXiv paper.
func HandleGetAbstract(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	arxivID, _ := getString(args, "arxivId")
	if arxivID == "" {
		return err("arxivId is required (e.g. '2303.08774')")
	}

	resp, fetchErr := arxivHTTP.Get(fmt.Sprintf("http://export.arxiv.org/api/query?id_list=%s", arxivID))
	if fetchErr != nil {
		return err(fmt.Sprintf("arxiv fetch: %v", fetchErr))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	text := string(body)

	title := cleanXML(extractXML(text, "<title>", "</title>"))
	abstract := cleanXML(extractXML(text, "<summary>", "</summary>"))
	authors := cleanXML(extractXML(extractXML(text, "<author>", "</author>"), "<name>", "</name>"))

	return ok(fmt.Sprintf("📄 %s\n\nAuthors: %s\n\n%s", title, authors, abstract))
}

func extractXML(s, open, close string) string {
	i := strings.Index(s, open)
	if i < 0 {
		return ""
	}
	i += len(open)
	j := strings.Index(s[i:], close)
	if j < 0 {
		return ""
	}
	return s[i : i+j]
}

func cleanXML(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\t", " ")
	// Collapse spaces
	parts := strings.Fields(s)
	return strings.Join(parts, " ")
}
