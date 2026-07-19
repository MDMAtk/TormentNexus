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

var batch11 = &http.Client{Timeout: 15 * time.Second}

func HandleNobelLaureates(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	limit, _ := getInt(args, "limit", 5)
	u := fmt.Sprintf("https://api.nobelprize.org/2.1/laureates?limit=%d", limit)
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return err("Nobel: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Laureates []struct {
			Name      string `json:"fullName"`
			BirthDate string `json:"birthDate"`
			DeathDate string `json:"deathDate"`
			Prizes    []struct {
				AwardYear string `json:"awardYear"`
				Category  string `json:"category"`
			} `json:"nobelPrizes"`
		} `json:"laureates"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Nobel Laureates:\n\n"))
	for _, l := range data.Laureates {
		sb.WriteString(fmt.Sprintf("%s\n", l.Name))
		for _, p := range l.Prizes {
			sb.WriteString(fmt.Sprintf("  %s Nobel Prize in %s\n", p.AwardYear, p.Category))
		}
	}
	return ok(sb.String())
}

func HandleNobelSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	category, _ := getString(args, "category")
	if name == "" && category == "" {
		return err("name or category is required")
	}
	u := "https://api.nobelprize.org/2.1/laureates?limit=10"
	if name != "" {
		u += "&query=" + url.QueryEscape(name)
	}
	if category != "" {
		u += "&nobelPrizeCategory=" + category
	}
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return err("Nobel: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Laureates []struct {
			ID   string `json:"id"`
			Name string `json:"fullName"`
		} `json:"laureates"`
	}
	json.Unmarshal(body, &data)
	if len(data.Laureates) == 0 {
		return ok(fmt.Sprintf("No Nobel laureates found for %q", name+category))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Nobel Laureates matching %q:\n\n", name))
	for _, l := range data.Laureates {
		sb.WriteString(fmt.Sprintf("%s (ID: %s)\n", l.Name, l.ID))
	}
	return ok(sb.String())
}

func HandleGeocode(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	place, _ := getString(args, "place")
	if place == "" {
		return err("place is required (city name, address)")
	}
	u := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=5", url.QueryEscape(place))
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return err("Geocoding: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Results []struct {
			Name      string  `json:"name"`
			Country   string  `json:"country"`
			Admin     string  `json:"admin1"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Elevation float64 `json:"elevation"`
			Timezone  string  `json:"timezone"`
		} `json:"results"`
	}
	json.Unmarshal(body, &data)
	if len(data.Results) == 0 {
		return ok(fmt.Sprintf("%q not found", place))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Places matching %q:\n\n", place))
	for _, r := range data.Results {
		sb.WriteString(fmt.Sprintf("%s, %s (%s)\n", r.Name, r.Country, r.Admin))
		sb.WriteString(fmt.Sprintf("  Lat: %.4f Lon: %.4f\n", r.Latitude, r.Longitude))
		sb.WriteString(fmt.Sprintf("  Elevation: %.0fm | TZ: %s\n", r.Elevation, r.Timezone))
	}
	return ok(sb.String())
}

func HandleMarineWeather(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude"); if lat == 0 { lat = 52.52 }
	lon := getFloat(args, "longitude"); if lon == 0 { lon = 13.41 }
	u := fmt.Sprintf("https://marine-api.open-meteo.com/v1/marine?latitude=%.2f&longitude=%.2f&hourly=wave_height,swell_wave_height,swell_wave_direction,swell_wave_period&forecast_days=1", lat, lon)
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Marine weather at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Hourly struct {
			Times  []string  `json:"time"`
			WaveH  []float64 `json:"wave_height"`
			SwellH []float64 `json:"swell_wave_height"`
			SwellD []float64 `json:"swell_wave_direction"`
			SwellP []float64 `json:"swell_wave_period"`
		} `json:"hourly"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Marine Weather at %.2f, %.2f:\n\n", lat, lon))
	for i, t := range data.Hourly.Times {
		if i >= 6 {
			break
		}
		h := data.Hourly.WaveH[i]
		sh := data.Hourly.SwellH[i]
		sp := data.Hourly.SwellP[i]
		sb.WriteString(fmt.Sprintf("%s: Wave %.1fm | Swell %.1fm (period: %.1fs)\n", t[11:16], h, sh, sp))
	}
	return ok(sb.String())
}

func HandleElevation(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude"); if lat == 0 { lat = 48.8567 }
	lon := getFloat(args, "longitude"); if lon == 0 { lon = 2.3510 }
	u := fmt.Sprintf("https://api.open-meteo.com/v1/elevation?latitude=%.4f&longitude=%.4f", lat, lon)
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Elevation at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Elevation []float64 `json:"elevation"`
	}
	json.Unmarshal(body, &data)
	if len(data.Elevation) == 0 {
		return ok(fmt.Sprintf("Elevation at %.2f, %.2f: N/A", lat, lon))
	}
	return ok(fmt.Sprintf("Elevation at %.2f, %.2f: %.0f meters", lat, lon, data.Elevation[0]))
}

func HandleWorldBankData(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	country, _ := getString(args, "country")
	indicator, _ := getString(args, "indicator")
	if country == "" {
		country = "US"
	}
	if indicator == "" {
		indicator = "NY.GDP.MKTP.CD"
	}
	labels := map[string]string{
		"NY.GDP.MKTP.CD": "GDP (current US$)",
		"SP.POP.TOTL": "Population, total",
		"NY.GDP.PCAP.CD": "GDP per capita (US$)",
		"NY.GDP.MKTP.KD.ZG": "GDP growth (annual %)",
		"FP.CPI.TOTL.ZG": "Inflation (CPI %)",
	}
	label := labels[indicator]
	if label == "" {
		label = indicator
	}
	u := fmt.Sprintf("https://api.worldbank.org/v2/country/%s/indicator/%s?format=json&per_page=10", country, indicator)
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("World Bank data for %s unavailable", country))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw []interface{}
	json.Unmarshal(body, &raw)
	if len(raw) < 2 {
		return ok(fmt.Sprintf("No data for %s/%s", country, indicator))
	}
	entries, _ := raw[1].([]interface{})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("World Bank — %s for %s:\n\n", label, country))
	for _, e := range entries {
		entry, _ := e.(map[string]interface{})
		if entry == nil {
			continue
		}
		year, _ := entry["year"].(string)
		val, _ := entry["value"].(interface{})
		if val == nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s: %v\n", year, val))
	}
	return ok(sb.String())
}

func HandleTVSchedule(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	country, _ := getString(args, "country")
	date, _ := getString(args, "date")
	if country == "" {
		country = "US"
	}
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	u := fmt.Sprintf("https://api.tvmaze.com/schedule?country=%s&date=%s", country, date)
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("TV schedule for %s on %s unavailable", country, date))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var entries []struct {
		Show struct {
			Name string `json:"name"`
		} `json:"show"`
		Name      string `json:"name"`
		Type    string `json:"type"`
		Runtime int    `json:"runtime"`
	}
	json.Unmarshal(body, &entries)
	if len(entries) == 0 {
		return ok(fmt.Sprintf("No TV schedule for %s on %s", country, date))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TV Schedule — %s on %s:\n\n", country, date))
	for _, e := range entries[:15] {
		sb.WriteString(fmt.Sprintf("%s — %s (%d min)\n", e.Show.Name, e.Name, e.Runtime))
	}
	return ok(sb.String())
}

func HandleTVEpisodes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	showID, _ := getInt(args, "showId", 1)
	n, _ := getInt(args, "limit", 10)
	u := fmt.Sprintf("https://api.tvmaze.com/shows/%d/episodes", showID)
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var episodes []struct {
		Season int    `json:"season"`
		Number int    `json:"number"`
		Name   string `json:"name"`
		Runtime int   `json:"runtime"`
		Airdate string `json:"airdate"`
	}
	json.Unmarshal(body, &episodes)
	if len(episodes) == 0 {
		return ok("No episodes found")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Episodes (%d total):\n\n", len(episodes)))
	for _, e := range episodes[:n] {
		sb.WriteString(fmt.Sprintf("S%02dE%02d: %s (%s, %d min)\n", e.Season, e.Number, e.Name, e.Airdate, e.Runtime))
	}
	return ok(sb.String())
}

func HandleRecentChanges(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch11.Get("https://openlibrary.org/recentchanges.json?limit=10")
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var changes []struct {
		Kind    string `json:"kind"`
		Comment string `json:"comment"`
		Author  string `json:"author"`
		IP      string `json:"ip"`
	}
	json.Unmarshal(body, &changes)
	var sb strings.Builder
	sb.WriteString("Recent Open Library Changes:\n\n")
	for _, c := range changes {
		who := c.Author
		if who == "" {
			who = c.IP
		}
		sb.WriteString(fmt.Sprintf("[%s] by %s: %s\n", c.Kind, who, c.Comment))
	}
	return ok(sb.String())
}

func HandleJokeCategories(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch11.Get("https://v2.jokeapi.dev/categories")
	if apiErr != nil {
		return err("JokeAPI: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Categories []string `json:"categories"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("JokeAPI Categories:\n\n%s", strings.Join(data.Categories, ", ")))
}

func HandleOpenLibraryByISBN(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	isbn, _ := getString(args, "isbn")
	if isbn == "" {
		return err("isbn is required")
	}
	u := fmt.Sprintf("https://openlibrary.org/search.json?q=isbn:%s", url.QueryEscape(isbn))
	resp, apiErr := batch11.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Docs []struct {
			Title    string   `json:"title"`
			Author   []string `json:"author_name"`
			Year     int      `json:"first_publish_year"`
			Edition  int      `json:"edition_count"`
		} `json:"docs"`
	}
	json.Unmarshal(body, &data)
	if len(data.Docs) == 0 {
		return ok(fmt.Sprintf("No results for ISBN %s", isbn))
	}
	d := data.Docs[0]
	return ok(fmt.Sprintf("ISBN %s: %s by %s (%d, %d editions)",
		isbn, d.Title, strings.Join(d.Author, ", "), d.Year, d.Edition))
}

// ═══════════════════════════════════════════════════════════════════
// Fix: replace getFloat usage with helpers.go's getFloat
// helpers.go getFloat(args, key) returns float64
// ═══════════════════════════════════════════════════════════════════

func getFloatOr(args map[string]interface{}, key string, def float64) float64 {
	v := getFloat(args, key)
	if v == 0 && def != 0 {
		return def
	}
	return v
}
