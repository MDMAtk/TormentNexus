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

var batch5API = &http.Client{Timeout: 15 * time.Second}

// ═══════════════════════════════════════════════════════════════════
// 35. ITUNES SEARCH — Search music, movies, podcasts (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleiTunesSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	media, _ := getString(args, "media")
	limit, _ := getInt(args, "limit", 5)
	if q == "" {
		return err("query is required")
	}
	if media == "" {
		media = "music"
	}
	u := fmt.Sprintf("https://itunes.apple.com/search?term=%s&media=%s&limit=%d", url.QueryEscape(q), media, limit)
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return err("iTunes: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		ResultCount int                      `json:"resultCount"`
		Results     []map[string]interface{} `json:"results"`
	}
	json.Unmarshal(body, &data)
	if data.ResultCount == 0 {
		return ok(fmt.Sprintf("No %s results for %q on iTunes", media, q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("iTunes Search — %d results for %q [%s]:\n\n", data.ResultCount, q, media))
	for i, r := range data.Results {
		if i >= limit {
			break
		}
		name, _ := r["trackName"].(string)
		if name == "" {
			name, _ = r["collectionName"].(string)
		}
		artist, _ := r["artistName"].(string)
		kind, _ := r["kind"].(string)
		price, _ := r["trackPrice"].(float64)
		currency, _ := r["currency"].(string)
		genre, _ := r["primaryGenreName"].(string)
		url, _ := r["trackViewUrl"].(string)

		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, name))
		sb.WriteString(fmt.Sprintf("   Artist: %s | Genre: %s | Type: %s\n", artist, genre, kind))
		if price > 0 {
			sb.WriteString(fmt.Sprintf("   Price: %.2f %s\n", price, currency))
		}
		if url != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", url))
		}
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 36. IPIFY — Get public IP address (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleMyIP(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch5API.Get("https://api.ipify.org?format=json")
	if apiErr != nil {
		return ok("Your IP: [unavailable]")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		IP string `json:"ip"`
	}
	json.Unmarshal(body, &data)
	if data.IP == "" {
		return ok("Your IP: [unavailable]")
	}
	return ok(fmt.Sprintf("Your public IP address: %s", data.IP))
}

// ═══════════════════════════════════════════════════════════════════
// 37. CHUCK NORRIS — Random jokes (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleChuckNorrisJoke(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	category, _ := getString(args, "category")
	u := "https://api.chucknorris.io/jokes/random"
	if category != "" {
		u += "?category=" + category
	}
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return ok("Chuck Norris doesn't use APIs. APIs use Chuck Norris.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Value string `json:"value"`
		URL   string `json:"url"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("💪 Chuck Norris: %s\n%s", data.Value, data.URL))
}

// HandleChuckCategories lists Chuck Norris joke categories.
func HandleChuckCategories(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch5API.Get("https://api.chucknorris.io/jokes/categories")
	if apiErr != nil {
		return err("ChuckNorris: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var cats []string
	json.Unmarshal(body, &cats)
	return ok(fmt.Sprintf("Chuck Norris joke categories:\n%s", strings.Join(cats, ", ")))
}

// ═══════════════════════════════════════════════════════════════════
// 38. KANYE WEST — Random Kanye quotes (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleKanyeQuote(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch5API.Get("https://api.kanye.rest")
	if apiErr != nil {
		return ok("\"The only disability in life is a bad attitude.\" — Scott Hamilton")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Quote string `json:"quote"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("\"%s\" — Kanye West", data.Quote))
}

// ═══════════════════════════════════════════════════════════════════
// 39. BACON IPSUM — Generate filler text (free)
// ═══════════════════════════════════════════════════════════════════

func HandleLoremIpsum(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	kind, _ := getString(args, "type")
	paras, _ := getInt(args, "paragraphs", 3)
	sentences, _ := getInt(args, "sentences", 0)
	if kind == "" {
		kind = "meat-and-filler"
	}
	u := fmt.Sprintf("https://baconipsum.com/api/?type=%s&paras=%d", kind, paras)
	if sentences > 0 {
		u += fmt.Sprintf("&sentences=%d", sentences)
	}
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return ok("Bacon ipsum dolor amet pastrami brisket jerky...")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var paragraphs []string
	json.Unmarshal(body, &paragraphs)
	return ok(fmt.Sprintf("📝 Generated Text:\n\n%s", strings.Join(paragraphs, "\n\n")))
}

// ═══════════════════════════════════════════════════════════════════
// 40. CORPORATE BS GENERATOR (free)
// ═══════════════════════════════════════════════════════════════════

func HandleCorporateBS(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch5API.Get("https://corporatebs-generator.sameerkumar.website/")
	if apiErr != nil {
		return ok("Synergize your core competencies to drive best-in-class ROI.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Phrase string `json:"phrase"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("🏢 Corporate BS: %s", data.Phrase))
}

// ═══════════════════════════════════════════════════════════════════
// 41. ZEN QUOTES — Inspirational quotes (free)
// ═══════════════════════════════════════════════════════════════════

func HandleZenQuote(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch5API.Get("https://zenquotes.io/api/random")
	if apiErr != nil {
		return ok("\"The journey of a thousand miles begins with a single step.\" — Lao Tzu")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data []struct {
		Q string `json:"q"`
		A string `json:"a"`
	}
	json.Unmarshal(body, &data)
	if len(data) == 0 || data[0].Q == "" {
		return ok("\"Be the change you wish to see in the world.\" — Gandhi")
	}
	return ok(fmt.Sprintf("\"%s\" — %s", data[0].Q, data[0].A))
}

// ═══════════════════════════════════════════════════════════════════
// 42. GITHUB TRENDING — Top starred repos (free)
// ═══════════════════════════════════════════════════════════════════

func HandleGitHubTrending(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lang, _ := getString(args, "language")
	since, _ := getString(args, "since")
	n, _ := getInt(args, "limit", 10)
	q := "stars:>1000"
	if since == "weekly" {
		q = "created:>7+days+stars:>100"
	} else if since == "daily" {
		q = "created:>1+day+stars:>10"
	}
	if lang != "" {
		q += "+language:" + lang
	}
	u := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=%d", url.QueryEscape(q), n)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "TormentNexus/1.0")

	resp, apiErr := batch5API.Do(req)
	if apiErr != nil {
		return err("GitHub: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Items []struct {
			FullName    string `json:"full_name"`
			Description string `json:"description"`
			Stars       int    `json:"stargazers_count"`
			Forks       int    `json:"forks_count"`
			Language    string `json:"language"`
			HTMLURL     string `json:"html_url"`
		} `json:"items"`
	}
	json.Unmarshal(body, &data)
	if len(data.Items) == 0 {
		return ok("No trending repos found")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🌟 GitHub Trending (%s):\n\n", since))
	for i, r := range data.Items {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, r.FullName))
		if r.Description != "" {
			d := r.Description
			if len(d) > 100 {
				d = d[:100] + "..."
			}
			sb.WriteString(fmt.Sprintf("   %s\n", d))
		}
		sb.WriteString(fmt.Sprintf("   ⭐ %d | 🍴 %d", r.Stars, r.Forks))
		if r.Language != "" {
			sb.WriteString(fmt.Sprintf(" | %s", r.Language))
		}
		sb.WriteString(fmt.Sprintf("\n   %s\n", r.HTMLURL))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 43. OPEN LIBRARY AUTHOR — Author info and works (free)
// ═══════════════════════════════════════════════════════════════════

func HandleAuthorLookup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	authorKey, _ := getString(args, "authorKey")
	name, _ := getString(args, "name")
	if authorKey == "" && name == "" {
		return err("authorKey or name is required")
	}
	var u string
	if authorKey != "" {
		u = fmt.Sprintf("https://openlibrary.org/authors/%s.json", authorKey)
	} else {
		u = fmt.Sprintf("https://openlibrary.org/search/authors.json?q=%s", url.QueryEscape(name))
	}
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if name != "" {
		var data struct {
			Docs []map[string]interface{} `json:"docs"`
		}
		json.Unmarshal(body, &data)
		if len(data.Docs) == 0 {
			return ok(fmt.Sprintf("No authors found for %q", name))
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Authors matching \"%s\":\n\n", name))
		for i, d := range data.Docs {
			if i >= 5 {
				break
			}
			n, _ := d["name"].(string)
			wk, _ := d["key"].(string)
			birth, _ := d["birth_date"].(string)
			works, _ := d["work_count"].(float64)
			sb.WriteString(fmt.Sprintf("%d. %s (born: %s, %d works)\n", i+1, n, birth, int(works)))
			sb.WriteString(fmt.Sprintf("   https://openlibrary.org%s\n", wk))
		}
		return ok(sb.String())
	}

	var data struct {
		Name         string `json:"name"`
		BirthDate    string `json:"birth_date"`
		DeathDate    string `json:"death_date"`
		Bio          string `json:"bio"`
		PersonalName string `json:"personal_name"`
		FullerName   string `json:"fuller_name"`
		Photos       []int  `json:"photos"`
	}
	json.Unmarshal(body, &data)
	bio := data.Bio
	if len(bio) > 500 {
		bio = bio[:500] + "..."
	}
	return ok(fmt.Sprintf("Author: %s\nFull name: %s\nBorn: %s | Died: %s\n\nBio: %s",
		data.Name, data.PersonalName, data.BirthDate, data.DeathDate, bio))
}

// ═══════════════════════════════════════════════════════════════════
// 44. SPACEX — Latest launches and info (free)
// ═══════════════════════════════════════════════════════════════════

func HandleSpaceXLaunch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	latest := getBool(args, "latest")
	next := getBool(args, "next")
	flight, _ := getInt(args, "flightNumber", 0)
	var u string
	if latest {
		u = "https://api.spacexdata.com/v5/launches/latest"
	} else if next {
		u = "https://api.spacexdata.com/v5/launches/next"
	} else if flight > 0 {
		u = fmt.Sprintf("https://api.spacexdata.com/v5/launches/%d", flight)
	} else {
		u = "https://api.spacexdata.com/v5/launches/latest"
	}
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return ok("SpaceX: service unavailable. Try again later.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Name         string `json:"name"`
		FlightNumber int    `json:"flight_number"`
		DateUTC      string `json:"date_utc"`
		Success      *bool  `json:"success"`
		Details      string `json:"details"`
		Links        struct {
			Webcast string `json:"webcast"`
			Article string `json:"article"`
		} `json:"links"`
		Failures []struct {
			Reason string `json:"reason"`
		} `json:"failures"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SpaceX — Launch #%d: %s\n", data.FlightNumber, data.Name))
	sb.WriteString(fmt.Sprintf("Date: %s\n", data.DateUTC[:10]))
	if data.Success != nil {
		if *data.Success {
			sb.WriteString("Status: ✅ Success\n")
		} else {
			sb.WriteString("Status: ❌ Failed\n")
			for _, f := range data.Failures {
				sb.WriteString(fmt.Sprintf("  Reason: %s\n", f.Reason))
			}
		}
	}
	if data.Details != "" {
		d := data.Details
		if len(d) > 500 {
			d = d[:500] + "..."
		}
		sb.WriteString(fmt.Sprintf("\n%s\n", d))
	}
	if data.Links.Webcast != "" {
		sb.WriteString(fmt.Sprintf("\nWebcast: %s\n", data.Links.Webcast))
	}
	if data.Links.Article != "" {
		sb.WriteString(fmt.Sprintf("Article: %s\n", data.Links.Article))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 45. PROGRAMMING QUOTES — Developer quotes (free)
// ═══════════════════════════════════════════════════════════════════

func HandleDevQuote(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch5API.Get("https://programming-quotesapi.vercel.app/api/random")
	if apiErr != nil {
		return ok("\"Any fool can write code that a computer can understand. Good programmers write code that humans can understand.\" — Martin Fowler")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Quote  string `json:"quote"`
		Author string `json:"author"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("💻 \"%s\" — %s", data.Quote, data.Author))
}

// ═══════════════════════════════════════════════════════════════════
// 46. OPEN LIBRARY WORK — Book/work details (free)
// ═══════════════════════════════════════════════════════════════════

func HandleWorkDetails(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	workKey, _ := getString(args, "workKey")
	if workKey == "" {
		return err("workKey is required (e.g. 'OL12345W')")
	}
	u := fmt.Sprintf("https://openlibrary.org/works/%s.json", workKey)
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Title  string `json:"title"`
		Author struct {
			Key string `json:"key"`
		} `json:"author"`
		FirstPublishYear int         `json:"first_publish_year"`
		Description      interface{} `json:"description"`
		Subjects         []string    `json:"subjects"`
		Excerpts         []struct {
			Text string `json:"text"`
		} `json:"excerpts"`
	}
	json.Unmarshal(body, &data)
	desc := ""
	if d, ok := data.Description.(string); ok {
		desc = d
	} else if m, ok := data.Description.(map[string]interface{}); ok {
		desc, _ = m["value"].(string)
	}
	if len(desc) > 500 {
		desc = desc[:500] + "..."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Work: %s\n", data.Title))
	sb.WriteString(fmt.Sprintf("First published: %d\n", data.FirstPublishYear))
	sb.WriteString(fmt.Sprintf("Author key: %s\n", data.Author.Key))
	if desc != "" {
		sb.WriteString(fmt.Sprintf("\nDescription:\n%s\n", desc))
	}
	if len(data.Subjects) > 0 {
		subj := data.Subjects
		if len(subj) > 10 {
			subj = subj[:10]
		}
		sb.WriteString(fmt.Sprintf("\nSubjects: %s\n", strings.Join(subj, ", ")))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 47. IP DATA — Combined IP + location (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleIPData(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	ip, _ := getString(args, "ip")
	u := "http://ip-api.com/json/"
	if ip != "" {
		u += ip
	}
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return ok("IP data: service unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Query       string  `json:"query"`
		Status      string  `json:"status"`
		Country     string  `json:"country"`
		CountryCode string  `json:"countryCode"`
		Region      string  `json:"regionName"`
		City        string  `json:"city"`
		Zip         string  `json:"zip"`
		Lat         float64 `json:"lat"`
		Lon         float64 `json:"lon"`
		Timezone    string  `json:"timezone"`
		ISP         string  `json:"isp"`
		Org         string  `json:"org"`
		As          string  `json:"as"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("IP Geolocation — %s\nStatus: %s\nLocation: %s, %s, %s %s\nCoords: %.4f, %.4f\nISP: %s\nOrg: %s\nASN: %s\nTimezone: %s",
		data.Query, data.Status, data.City, data.Region, data.Country, data.Zip,
		data.Lat, data.Lon, data.ISP, data.Org, data.As, data.Timezone))
}

// ═══════════════════════════════════════════════════════════════════
// 48. COLOR API — Color information (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleColorInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	hex, _ := getString(args, "hex")
	name, _ := getString(args, "name")
	if hex == "" && name == "" {
		return err("hex or name is required (e.g. 'ff0000' for red)")
	}
	var u string
	if hex != "" {
		u = fmt.Sprintf("https://www.thecolorapi.com/id?hex=%s", strings.TrimPrefix(hex, "#"))
	} else {
		u = fmt.Sprintf("https://www.thecolorapi.com/id?name=%s", url.QueryEscape(name))
	}
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Color %s%s: service unavailable", hex, name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name struct {
			Value string `json:"value"`
		} `json:"name"`
		Hex struct {
			Value string `json:"value"`
		} `json:"hex"`
		RGB struct {
			Value string `json:"value"`
		} `json:"rgb"`
		HSL struct {
			Value string `json:"value"`
		} `json:"hsl"`
		Image struct {
			Bare string `json:"bare"`
		} `json:"image"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Color: %s (%s)\nHex: %s\nRGB: %s\nHSL: %s\nPreview: https://www.thecolorapi.com/id?format=svg&hex=%s",
		data.Name.Value, hex, data.Hex.Value, data.RGB.Value, data.HSL.Value, strings.TrimPrefix(data.Hex.Value, "#")))
}

// ═══════════════════════════════════════════════════════════════════
// 49. DATE NAMEDAY — Nameday calendar (free)
// ═══════════════════════════════════════════════════════════════════

func HandleNameday(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	country, _ := getString(args, "country")
	d, _ := getString(args, "date")
	if country == "" {
		country = "us"
	}
	var u string
	if d != "" {
		u = fmt.Sprintf("https://nameday.abalin.net/api/V1/getdate?country=%s&day=%s&month=%s", country,
			strings.Split(d, "-")[2], strings.Split(d, "-")[1])
	} else {
		u = fmt.Sprintf("https://nameday.abalin.net/api/V1/today?country=%s", country)
	}
	resp, apiErr := batch5API.Get(u)
	if apiErr != nil {
		return ok("Nameday: service unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Nameday map[string]string `json:"nameday"`
		Date    string            `json:"date"`
	}
	json.Unmarshal(body, &data)
	if name, exists := data.Nameday[country]; exists {
		return ok(fmt.Sprintf("Nameday for %s (%s): %s", data.Date, country, name))
	}
	// Try 'us' fallback
	if name, exists := data.Nameday["us"]; exists {
		return ok(fmt.Sprintf("Nameday for %s: %s", data.Date, name))
	}
	return ok(fmt.Sprintf("Nameday for %s: service data", data.Date))
}

// ── HELPERS ───────────────────────────────────────────────────────
