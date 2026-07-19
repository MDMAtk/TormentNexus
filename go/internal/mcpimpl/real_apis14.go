package mcpimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var batch14 = &http.Client{Timeout: 15 * time.Second}

func HandleNEOFeed(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	start, _ := getString(args, "startDate")
	end, _ := getString(args, "endDate")
	if start == "" {
		start = time.Now().Format("2006-01-02")
	}
	if end == "" {
		end = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}
	u := fmt.Sprintf("https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=DEMO_KEY", start, end)
	resp, apiErr := batch14.Get(u)
	if apiErr != nil {
		return ok("NEO feed unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ElementCount int `json:"element_count"`
		NearEarth    map[string][]struct {
			Name        string  `json:"name"`
			Diameter    float64 `json:"estimated_diameter_min"`
			Hazardous   bool    `json:"is_potentially_hazardous_asteroid"`
			Magnitude   float64 `json:"absolute_magnitude_h"`
		} `json:"near_earth_objects"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Near Earth Objects — %d detected (%s to %s):\n\n", data.ElementCount, start, end))
	for date, objects := range data.NearEarth {
		sb.WriteString(fmt.Sprintf("%s (%d objects):\n", date, len(objects)))
		for _, o := range objects {
			hazard := "✅"
			if o.Hazardous {
				hazard = "⚠️ POTENTIALLY HAZARDOUS"
			}
			sb.WriteString(fmt.Sprintf("  %s | Diam: %.2f m | Mag: %.1f %s\n", o.Name, o.Diameter, o.Magnitude, hazard))
		}
	}
	return ok(sb.String())
}

func HandleTVSeasons(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	showID, _ := getInt(args, "showId", 1)
	u := fmt.Sprintf("https://api.tvmaze.com/shows/%d/seasons", showID)
	resp, apiErr := batch14.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var seasons []struct {
		Number   int    `json:"number"`
		Episodes int    `json:"episodeOrder"`
		Airdate  string `json:"premiereDate"`
		EndDate  string `json:"endDate"`
		Summary  string `json:"summary"`
	}
	json.Unmarshal(body, &seasons)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Seasons for show %d (%d seasons):\n\n", showID, len(seasons)))
	for _, s := range seasons {
		sb.WriteString(fmt.Sprintf("Season %d: %d episodes (%s to %s)\n", s.Number, s.Episodes, s.Airdate, s.EndDate))
	}
	return ok(sb.String())
}

func HandleSeasonEpisodes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	seasonID, _ := getInt(args, "seasonId", 1)
	u := fmt.Sprintf("https://api.tvmaze.com/seasons/%d/episodes", seasonID)
	resp, apiErr := batch14.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var episodes []struct {
		Number   int    `json:"number"`
		Name     string `json:"name"`
		Airdate  string `json:"airdate"`
		Runtime  int    `json:"runtime"`
	}
	json.Unmarshal(body, &episodes)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Episodes for season %d (%d episodes):\n\n", seasonID, len(episodes)))
	for _, e := range episodes {
		sb.WriteString(fmt.Sprintf("E%02d: %s (%s, %d min)\n", e.Number, e.Name, e.Airdate, e.Runtime))
	}
	return ok(sb.String())
}

func HandlePokeLocation(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/location/%s", strings.ToLower(q))
	resp, apiErr := batch14.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Location %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Region  struct {
			Name string `json:"name"`
		} `json:"region"`
		Areas []struct {
			Name string `json:"name"`
		} `json:"areas"`
	}
	json.Unmarshal(body, &data)
	var areas []string
	for _, a := range data.Areas {
		areas = append(areas, a.Name)
	}
	return ok(fmt.Sprintf("Location: %s\nRegion: %s\nAreas (%d): %s", data.Name, data.Region.Name, len(areas), strings.Join(areas, ", ")))
}

func HandlePokeType(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required (e.g. 'fire', 'water', 'electric')")
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/type/%s", strings.ToLower(name))
	resp, apiErr := batch14.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Type %q not found", name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	type pokeType struct { Name string }
	var data struct {
		Name   string `json:"name"`
		Damage struct {
			Double []pokeType `json:"double_damage_to"`
			Half   []pokeType `json:"half_damage_to"`
			No     []pokeType `json:"no_damage_to"`
		} `json:"damage_relations"`
	}
	json.Unmarshal(body, &data)
	join := func(items []pokeType) string {
		var s []string
		for _, item := range items {
			s = append(s, item.Name)
		}
		return strings.Join(s, ", ")
	}
	return ok(fmt.Sprintf("Type: %s\n\nStrong against: %s\nWeak against: %s\nImmune against: %s",
		data.Name, join(data.Damage.Double), join(data.Damage.Half), join(data.Damage.No)))
}

func HandleOpenLibraryLanguages(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch14.Get("https://openlibrary.org/languages.json")
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var langs []struct {
		Name string `json:"name"`
		Key  string `json:"key"`
	}
	json.Unmarshal(body, &langs)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Open Library Languages (%d):\n\n", len(langs)))
	for _, l := range langs {
		sb.WriteString(fmt.Sprintf("%s (%s)\n", l.Name, l.Key))
	}
	return ok(sb.String())
}

func HandleExoplanets(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	n, _ := getInt(args, "limit", 5)
	u := fmt.Sprintf("https://exoplanetarchive.ipac.caltech.edu/TAP/sync?query=select+top+%d+pl_name,pl_rade,pl_orbper,pl_orbeccen,st_spectype+from+ps&format=json", n)
	resp, apiErr := batch14.Get(u)
	if apiErr != nil {
		return ok("Exoplanet data unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var planets []map[string]interface{}
	json.Unmarshal(body, &planets)
	if len(planets) == 0 {
		return ok("No exoplanet data")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Exoplanets (%d found):\n\n", len(planets)))
	for _, p := range planets {
		name, _ := p["pl_name"].(string)
		rade, _ := p["pl_rade"].(float64)
		orbper, _ := p["pl_orbper"].(float64)
		spec, _ := p["st_spectype"].(string)
		sb.WriteString(fmt.Sprintf("%s: R=%.2f Earth radii, P=%.1f days, Host: %s\n", name, rade, orbper, spec))
	}
	return ok(sb.String())
}

func HandleOpenLibraryBookFormats(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	bookKey, _ := getString(args, "bookKey")
	if bookKey == "" {
		return err("bookKey is required (e.g. 'OL7353617M')")
	}
	u := fmt.Sprintf("https://openlibrary.org/books/%s.json", bookKey)
	resp, apiErr := batch14.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Title       string `json:"title"`
		Publishers  []string `json:"publishers"`
		PublishDate string `json:"publish_date"`
		Pages       int    `json:"number_of_pages"`
		Formats     map[string]interface{} `json:"formats"`
		Identifiers map[string]interface{} `json:"identifiers"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Book: %s\n", data.Title))
	sb.WriteString(fmt.Sprintf("Published: %s by %s\n", data.PublishDate, strings.Join(data.Publishers, ", ")))
	if data.Pages > 0 {
		sb.WriteString(fmt.Sprintf("Pages: %d\n", data.Pages))
	}
	for name, val := range data.Identifiers {
		sb.WriteString(fmt.Sprintf("%s: %v\n", name, val))
	}
	for name := range data.Formats {
		sb.WriteString(fmt.Sprintf("Format: %s\n", name))
	}
	return ok(sb.String())
}
