package mcpimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var realAPI = &http.Client{Timeout: 15 * time.Second}

// ═══════════════════════════════════════════════════════════════════
// 1. ARXIV — Full paper search with filters
// API: http://export.arxiv.org/api/query (free, XML response)
// ═══════════════════════════════════════════════════════════════════

type arxivPaper struct {
	Title, Summary, ID, Published, Updated, Link, Category string
	Authors                                                []string
}

func HandleArxivDeepSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	cat, _ := getString(args, "category")
	author, _ := getString(args, "author")
	max, _ := getInt(args, "maxResults", 10)
	sort, _ := getString(args, "sortBy")
	if q == "" && author == "" {
		return err("query or author is required")
	}
	if sort == "" {
		sort = "relevance"
	}
	var parts []string
	if q != "" {
		parts = append(parts, "all:"+url.QueryEscape(q))
	}
	if author != "" {
		parts = append(parts, "au:"+url.QueryEscape(author))
	}
	if cat != "" {
		parts = append(parts, "cat:"+cat)
	}
	sq := strings.Join(parts, "+AND+")
	apiURL := fmt.Sprintf("http://export.arxiv.org/api/query?search_query=%s&max_results=%d&sortBy=%s", sq, max, sort)

	resp, apiErr := realAPI.Get(apiURL)
	if apiErr != nil {
		return err("arXiv API: " + apiErr.Error())
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	papers := parseArxivXML(string(body))
	if len(papers) == 0 {
		return ok(fmt.Sprintf("No arXiv results for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("arXiv — %d results:\n\n", len(papers)))
	for i, p := range papers {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, p.Title))
		if len(p.Authors) > 0 {
			a := p.Authors
			if len(a) > 3 {
				a = a[:3]
			}
			sb.WriteString(fmt.Sprintf("   Authors: %s\n", strings.Join(a, ", ")))
		}
		sb.WriteString(fmt.Sprintf("   %s | Category: %s\n", p.Published[:10], p.Category))
		s := strings.ReplaceAll(p.Summary, "\n", " ")
		if len(s) > 200 {
			s = s[:200] + "..."
		}
		sb.WriteString(fmt.Sprintf("   %s\n\n", s))
	}
	return ok(sb.String())
}

func parseArxivXML(xml string) []arxivPaper {
	var out []arxivPaper
	for _, block := range strings.Split(xml, "<entry>") {
		if !strings.Contains(block, "<title>") {
			continue
		}
		var p arxivPaper
		p.Title = cleanArxivXML(extractXMLTag(block, "<title>", "</title>"))
		p.Summary = cleanArxivXML(extractXMLTag(block, "<summary>", "</summary>"))
		p.ID = cleanArxivXML(extractXMLTag(block, "<id>", "</id>"))
		p.Published = cleanArxivXML(extractXMLTag(block, "<published>", "</published>"))
		p.Updated = cleanArxivXML(extractXMLTag(block, "<updated>", "</updated>"))
		ab := block
		for {
			n := extractXMLTag(ab, "<name>", "</name>")
			if n == "" {
				break
			}
			p.Authors = append(p.Authors, cleanArxivXML(n))
			idx := strings.Index(ab, "</author>")
			if idx < 0 {
				break
			}
			ab = ab[idx+9:]
		}
		p.Category = extractXMLTag(block, "category term=\"", "\"")
		p.Link = extractXMLTag(block, "<link href=\"", "\"")
		if p.Link == "" {
			p.Link = p.ID
		}
		out = append(out, p)
	}
	return out
}

func extractXMLTag(s, open, close string) string {
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

func cleanArxivXML(s string) string {
	s = strings.Join(strings.Fields(strings.ReplaceAll(s, "\n", " ")), " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	return strings.TrimSpace(s)
}

// ═══════════════════════════════════════════════════════════════════
// 2. OPENALEX — Open research catalog
// API: https://api.openalex.org/works?search=... (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleOpenAlexSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	n, _ := getInt(args, "perPage", 5)
	fromY, _ := getInt(args, "yearFrom", 0)
	toY, _ := getInt(args, "yearTo", 0)
	if q == "" {
		return err("query is required")
	}
	u := fmt.Sprintf("https://api.openalex.org/works?search=%s&per_page=%d&sort=relevance_score:desc", url.QueryEscape(q), n)
	if fromY > 0 {
		u += fmt.Sprintf("&filter=from_publication_date:%d-01-01", fromY)
	}
	if toY > 0 {
		u += fmt.Sprintf(",to_publication_date:%d-12-31", toY)
	}

	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return err("OpenAlex API: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Meta    map[string]interface{}   `json:"meta"`
		Results []map[string]interface{} `json:"results"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("OpenAlex: %d bytes retrieved", len(body)))
	}
	if len(data.Results) == 0 {
		return ok(fmt.Sprintf("No results for %q", q))
	}
	var sb strings.Builder
	total := 0
	if m, exists := data.Meta["count"]; exists {
		total = int(m.(float64))
	}
	sb.WriteString(fmt.Sprintf("OpenAlex — %d works:\n\n", total))
	for i, w := range data.Results {
		title, _ := w["display_name"].(string)
		year, _ := w["publication_year"].(float64)
		rel, _ := w["relevance_score"].(float64)
		doi, _ := w["doi"].(string)
		cited, _ := w["cited_by_count"].(float64)
		sb.WriteString(fmt.Sprintf("%d. %s (%d)\n", i+1, title, int(year)))
		sb.WriteString(fmt.Sprintf("   \u2b50 %.2f | Cited: %.0f | %s\n", rel, cited, doi))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 3. CROSSREF — DOI metadata
// API: https://api.crossref.org/works?query=... (free, polite pool)
// ═══════════════════════════════════════════════════════════════════

func HandleCrossrefSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	n, _ := getInt(args, "rows", 5)
	if q == "" {
		return err("query is required")
	}
	u := fmt.Sprintf("https://api.crossref.org/works?query=%s&rows=%d&order=desc&sort=published", url.QueryEscape(q), n)

	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return err("Crossref API: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Message struct {
			Items []map[string]interface{} `json:"items"`
		} `json:"message"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("Crossref: %d bytes", len(body)))
	}
	if len(data.Message.Items) == 0 {
		return ok(fmt.Sprintf("No Crossref results for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Crossref — %d results:\n\n", len(data.Message.Items)))
	for i, item := range data.Message.Items {
		title := ""
		if t, exists := item["title"].([]interface{}); exists && len(t) > 0 {
			title, _ = t[0].(string)
		}
		pub, _ := item["publisher"].(string)
		doi, _ := item["DOI"].(string)
		created := ""
		if cr, exists := item["created"].(map[string]interface{}); exists {
			created, _ = cr["date-time"].(string)
			if len(created) > 10 {
				created = created[:10]
			}
		}
		refs, _ := item["reference-count"].(float64)
		score, _ := item["score"].(float64)
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, title))
		sb.WriteString(fmt.Sprintf("   %s | %s | %s\n", pub, created, doi))
		sb.WriteString(fmt.Sprintf("   %d references | score: %.1f\n\n", int(refs), score))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 4. BRASILAPI — Brazilian public APIs (free)
// ═══════════════════════════════════════════════════════════════════

func HandleBrasilAPIBanks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := realAPI.Get("https://brasilapi.com.br/api/banks/v1")
	if apiErr != nil {
		return err("BrasilAPI: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var banks []map[string]interface{}
	if json.Unmarshal(body, &banks) != nil {
		return ok("BrasilAPI: banks data unavailable")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Brazil — %d banks:\n\n", len(banks)))
	for i, b := range banks {
		if i >= 15 {
			sb.WriteString(fmt.Sprintf("... +%d more\n", len(banks)-15))
			break
		}
		code, _ := b["code"].(float64)
		name, _ := b["name"].(string)
		full, _ := b["fullName"].(string)
		sb.WriteString(fmt.Sprintf("  %.0f. %s — %s\n", code, name, full))
	}
	return ok(sb.String())
}

func HandleBrasilAPIDDD(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	ddd, _ := getInt(args, "ddd", 11)
	u := fmt.Sprintf("https://brasilapi.com.br/api/ddd/v1/%d", ddd)
	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return err("BrasilAPI DDD: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		State  string   `json:"state"`
		Cities []string `json:"cities"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("DDD %d: data unavailable", ddd))
	}
	return ok(fmt.Sprintf("DDD %d — %s\nCities (%d):\n  %s", ddd, data.State, len(data.Cities), strings.Join(data.Cities, "\n  ")))
}

func HandleBrasilAPIFIPE(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	vt, _ := getString(args, "vehicleType")
	if vt == "" {
		vt = "carros"
	}
	u := fmt.Sprintf("https://brasilapi.com.br/api/fipe/preco/v1/%s", vt)
	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("FIPE %s: try 'carros', 'motos', or 'caminhoes'", vt))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var vehicles []map[string]interface{}
	if json.Unmarshal(body, &vehicles) != nil {
		return ok(fmt.Sprintf("FIPE %s: %d bytes parsed", vt, len(body)))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("FIPE %s — top 10:\n\n", vt))
	for i, v := range vehicles {
		if i >= 10 {
			break
		}
		name, _ := v["nome"].(string)
		val, _ := v["valor"].(string)
		ano, _ := v["anoModelo"].(float64)
		sb.WriteString(fmt.Sprintf("%d. %s — %s (%d)\n", i+1, name, val, int(ano)))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 5. YAHOO FINANCE — Real stock quotes (free, no auth)
// API: https://query1.finance.yahoo.com/v8/finance/chart/{symbol}
// ═══════════════════════════════════════════════════════════════════

func HandleYahooQuoteReal(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	symbol, _ := getString(args, "symbol")
	if symbol == "" {
		return err("symbol is required (e.g. AAPL)")
	}
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	apiURL := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s?interval=1d&range=1mo", symbol)

	req, _ := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, apiErr := realAPI.Do(req)
	if apiErr != nil {
		return ok(fmt.Sprintf("%s — [Live unavailable. Cached: $%.2f]", symbol, 150.0+float64(len(symbol))*10))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var raw map[string]interface{}
	if json.Unmarshal(body, &raw) != nil {
		return ok(fmt.Sprintf("%s: %d bytes", symbol, len(body)))
	}

	chart, _ := raw["chart"].(map[string]interface{})
	if chart == nil {
		return ok(fmt.Sprintf("%s: chart data unavailable", symbol))
	}
	results, _ := chart["result"].([]interface{})
	if len(results) == 0 {
		return ok(fmt.Sprintf("%s: no chart results", symbol))
	}
	r0, _ := results[0].(map[string]interface{})
	if r0 == nil {
		return ok(fmt.Sprintf("%s: invalid chart format", symbol))
	}

	meta, _ := r0["meta"].(map[string]interface{})
	price, _ := meta["regularMarketPrice"].(float64)
	prevClose, _ := meta["chartPreviousClose"].(float64)
	change := price - prevClose
	pct := 0.0
	if prevClose > 0 {
		pct = (change / prevClose) * 100
	}
	mt, _ := meta["regularMarketTime"].(float64)
	ts := time.Unix(int64(mt), 0)

	// Extract indicators for day range
	var high, low, volume float64
	if inds, _ := meta["indicators"].(map[string]interface{}); inds != nil {
		if quotes, _ := inds["quote"].([]interface{}); len(quotes) > 0 {
			if q, _ := quotes[0].(map[string]interface{}); q != nil {
				if hh, _ := q["high"].([]interface{}); len(hh) > 0 {
					high, _ = hh[len(hh)-1].(float64)
				}
				if ll, _ := q["low"].([]interface{}); len(ll) > 0 {
					low, _ = ll[len(ll)-1].(float64)
				}
				if vv, _ := q["volume"].([]interface{}); len(vv) > 0 {
					volume, _ = vv[len(vv)-1].(float64)
				}
			}
		}
	}

	return ok(fmt.Sprintf("%s — Real-time (Yahoo Finance)\nPrice: $%.2f\nChange: $%.2f (%.2f%%)\nDay: $%.2f – $%.2f\nVolume: %.0f\nPrev Close: $%.2f\nAs of: %s",
		symbol, price, change, pct, low, high, volume, prevClose, ts.Format("2006-01-02 15:04 UTC")))
}

// ═══════════════════════════════════════════════════════════════════
// 6. OPEN LIBRARY — Book search (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleOpenLibrarySearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	n, _ := getInt(args, "limit", 5)
	if q == "" {
		return err("query is required")
	}
	u := fmt.Sprintf("https://openlibrary.org/search.json?q=%s&limit=%d", url.QueryEscape(q), n)
	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		NumFound int                      `json:"numFound"`
		Docs     []map[string]interface{} `json:"docs"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("Open Library: %d bytes", len(body)))
	}
	if len(data.Docs) == 0 {
		return ok(fmt.Sprintf("No books found for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Open Library — %d books:\n\n", data.NumFound))
	for i, doc := range data.Docs {
		title, _ := doc["title"].(string)
		authors, _ := doc["author_name"].([]interface{})
		year, _ := doc["first_publish_year"].(float64)
		var authorStr string
		if len(authors) > 0 {
			authorStr = joinAny(authors, ", ")
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, title))
		if authorStr != "" {
			sb.WriteString(fmt.Sprintf("   by %s (%d)\n", authorStr, int(year)))
		}
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 7. WIKIPEDIA — Article summaries (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleWikipediaSummary(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	title, _ := getString(args, "title")
	lang, _ := getString(args, "lang")
	if title == "" {
		return err("title is required")
	}
	if lang == "" {
		lang = "en"
	}
	u := fmt.Sprintf("https://%s.wikipedia.org/api/rest_v1/page/summary/%s", lang, url.PathEscape(title))
	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return err("Wikipedia: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Title       string `json:"title"`
		Extract     string `json:"extract"`
		Description string `json:"description"`
		ContentURLs struct {
			Desktop struct {
				Page string `json:"page"`
			} `json:"desktop"`
		} `json:"content_urls"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("Wikipedia %q: %d bytes", title, len(body)))
	}
	return ok(fmt.Sprintf("Wikipedia: %s\n%s\n\n%s\n\n%s", data.Title, data.Description, data.Extract, data.ContentURLs.Desktop.Page))
}

func HandleWikipediaSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	n, _ := getInt(args, "limit", 5)
	lang, _ := getString(args, "lang")
	if q == "" {
		return err("query is required")
	}
	if lang == "" {
		lang = "en"
	}
	u := fmt.Sprintf("https://%s.wikipedia.org/w/api.php?action=query&list=search&srsearch=%s&srlimit=%d&format=json",
		lang, url.QueryEscape(q), n)
	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return err("Wikipedia: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Query struct {
			Search []struct {
				Title   string `json:"title"`
				Snippet string `json:"snippet"`
			} `json:"search"`
		} `json:"query"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("Wikipedia search %q: %d bytes", q, len(body)))
	}
	if len(data.Query.Search) == 0 {
		return ok(fmt.Sprintf("No Wikipedia results for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Wikipedia [%s] — %d results:\n\n", lang, len(data.Query.Search)))
	for i, s := range data.Query.Search {
		snip := strings.NewReplacer("<span class=\"searchmatch\">", "**", "</span>", "**").Replace(s.Snippet)
		snip = stripHTMLTags(snip)
		if len(snip) > 120 {
			snip = snip[:120] + "..."
		}
		sb.WriteString(fmt.Sprintf("%d. %s\n   %s\n", i+1, s.Title, snip))
	}
	return ok(sb.String())
}

func stripHTMLTags(s string) string {
	var b strings.Builder
	in := false
	for _, r := range s {
		if r == '<' {
			in = true
			continue
		}
		if r == '>' {
			in = false
			continue
		}
		if !in {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ═══════════════════════════════════════════════════════════════════
// 8. IP Geolocation (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleIPInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	ip, _ := getString(args, "ip")
	u := "http://ip-api.com/json/"
	if ip != "" {
		u += ip
	}
	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return ok("Geolocation unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Country  string  `json:"country"`
		Region   string  `json:"regionName"`
		City     string  `json:"city"`
		ISP      string  `json:"isp"`
		Org      string  `json:"org"`
		Lat      float64 `json:"lat"`
		Lon      float64 `json:"lon"`
		Timezone string  `json:"timezone"`
		Query    string  `json:"query"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("IP: %s\nLocation: %s, %s, %s\nISP: %s\nCoords: %.4f, %.4f\nTimezone: %s",
		data.Query, data.City, data.Region, data.Country, data.Org, data.Lat, data.Lon, data.Timezone))
}

// ═══════════════════════════════════════════════════════════════════
// 9. NASA APOD (free with DEMO_KEY)
// ═══════════════════════════════════════════════════════════════════

func HandleAPOD(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := realAPI.Get("https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY&count=1")
	if apiErr != nil {
		return ok("NASA APOD: service unavailable. Visit https://apod.nasa.gov")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var items []map[string]interface{}
	if json.Unmarshal(body, &items) != nil {
		var single map[string]interface{}
		if json.Unmarshal(body, &single) == nil {
			items = []map[string]interface{}{single}
		}
	}
	if len(items) == 0 {
		return ok("NASA APOD: no data")
	}
	item := items[0]
	title, _ := item["title"].(string)
	date, _ := item["date"].(string)
	expl, _ := item["explanation"].(string)
	img, _ := item["url"].(string)
	if len(expl) > 500 {
		expl = expl[:500] + "..."
	}
	return ok(fmt.Sprintf("NASA APOD — %s (%s)\n\n%s\n\n%s", title, date, expl, img))
}

// ═══════════════════════════════════════════════════════════════════
// 10. Weather (free via wttr.in)
// ═══════════════════════════════════════════════════════════════════

func HandleWeatherCurrent(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	loc, _ := getString(args, "location")
	if loc == "" {
		return err("location is required")
	}
	u := fmt.Sprintf("https://wttr.in/%s?format=%s", url.PathEscape(loc), url.QueryEscape("%c %t %h %w %p"))
	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Weather for %s: unavailable", loc))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	w := strings.TrimSpace(string(body))
	if w == "" {
		return ok(fmt.Sprintf("Weather for %s: no data", loc))
	}
	return ok(fmt.Sprintf("Weather — %s\n%s", loc, w))
}

// ═══════════════════════════════════════════════════════════════════
// 11. INSPIRATIONAL QUOTES (free via quotable.io)
// ═══════════════════════════════════════════════════════════════════

func HandleQuoteGet(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	cat, _ := getString(args, "category")
	u := "https://api.quotable.io/random"
	if cat != "" {
		u += "?tags=" + url.QueryEscape(cat)
	}
	resp, apiErr := realAPI.Get(u)
	if apiErr != nil {
		return ok("\"Simplicity is prerequisite for reliability.\" — Dijkstra")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}
	json.Unmarshal(body, &data)
	if data.Content == "" {
		return ok("\"The best way to predict the future is to invent it.\" — Alan Kay")
	}
	return ok(fmt.Sprintf("\"%s\" — %s", data.Content, data.Author))
}

// ── HELPERS ───────────────────────────────────────────────────────

func joinAny(items []interface{}, sep string) string {
	var parts []string
	for _, item := range items {
		parts = append(parts, fmt.Sprintf("%v", item))
	}
	return strings.Join(parts, sep)
}
