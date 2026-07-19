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

var batch19 = &http.Client{Timeout: 15 * time.Second}

func HandleTVLookup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getString(args, "externalId")
	source, _ := getString(args, "source")
	if name != "" {
		u := fmt.Sprintf("https://api.tvmaze.com/singlesearch/shows?q=%s", url.QueryEscape(name))
		resp, apiErr := batch19.Get(u)
		if apiErr != nil {
			return ok(fmt.Sprintf("Show %q not found", name))
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var show struct {
			Name   string `json:"name"`
			Type   string `json:"type"`
			Status string `json:"status"`
			Rating struct {
				Average float64 `json:"average"`
			} `json:"rating"`
			Genres []string `json:"genres"`
			URL    string   `json:"url"`
		}
		json.Unmarshal(body, &show)
		return ok(fmt.Sprintf("Show: %s\nType: %s | Status: %s\nRating: %.1f | Genres: %s\n%s",
			show.Name, show.Type, show.Status, show.Rating.Average, strings.Join(show.Genres, ", "), show.URL))
	}
	if id != "" && source != "" {
		u := fmt.Sprintf("https://api.tvmaze.com/lookup/shows?%s=%s", source, id)
		resp, apiErr := batch19.Get(u)
		if apiErr != nil {
			return ok(fmt.Sprintf("Show not found for %s=%s", source, id))
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var show struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
		json.Unmarshal(body, &show)
		return ok(fmt.Sprintf("Show: %s\n%s", show.Name, show.URL))
	}
	return err("name or externalId+source required")
}

func HandleContestEffect(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 1)
	u := fmt.Sprintf("https://pokeapi.co/api/v2/contest-effect/%d", id)
	resp, apiErr := batch19.Get(u)
	if apiErr != nil {
		return err("Contest: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ID     int `json:"id"`
		Appeal int `json:"appeal"`
		Jam    int `json:"jam"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Contest Effect #%d: Appeal %d, Jam %d", data.ID, data.Appeal, data.Jam))
}

func HandleSuperContest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 1)
	u := fmt.Sprintf("https://pokeapi.co/api/v2/super-contest-effect/%d", id)
	resp, apiErr := batch19.Get(u)
	if apiErr != nil {
		return err("Super contest: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ID     int `json:"id"`
		Appeal int `json:"appeal"`
		Moves  []struct {
			Name string `json:"name"`
		} `json:"moves"`
	}
	json.Unmarshal(body, &data)
	var moves []string
	for _, m := range data.Moves {
		moves = append(moves, m.Name)
	}
	return ok(fmt.Sprintf("Super Contest #%d: Appeal %d\nMoves (%d): %s", data.ID, data.Appeal, len(moves), strings.Join(moves, ", ")))
}

func HandleEncounterCondition(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/encounter-condition/%s", strings.ToLower(q))
	resp, apiErr := batch19.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Condition %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name   string `json:"name"`
		Values []struct {
			Name string `json:"name"`
		} `json:"values"`
	}
	json.Unmarshal(body, &data)
	var vals []string
	for _, v := range data.Values {
		vals = append(vals, v.Name)
	}
	return ok(fmt.Sprintf("Encounter Condition: %s\nValues: %s", data.Name, strings.Join(vals, ", ")))
}

func HandlePalPark(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/pal-park-area/%s", strings.ToLower(q))
	resp, apiErr := batch19.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Pal Park %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name       string `json:"name"`
		Encounters []struct {
			Species struct {
				Name string `json:"name"`
			} `json:"pokemon_species"`
			Rate int `json:"rate"`
		} `json:"pokemon_encounters"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pal Park Area: %s (%d encounters)\n\n", data.Name, len(data.Encounters)))
	for _, e := range data.Encounters[:10] {
		sb.WriteString(fmt.Sprintf("%s (rate: %d)\n", e.Species.Name, e.Rate))
	}
	return ok(sb.String())
}

func HandleOpenLibrarySearchSort(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	limit, _ := getInt(args, "limit", 5)
	sort, _ := getString(args, "sort")
	if q == "" {
		return err("query is required")
	}
	u := fmt.Sprintf("https://openlibrary.org/search.json?q=%s&limit=%d", url.QueryEscape(q), limit)
	if sort != "" {
		u += "&sort=" + sort
	}
	resp, apiErr := batch19.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Docs []struct {
			Title   string `json:"title"`
			Author  string `json:"author_name"`
			Year    int    `json:"first_publish_year"`
			Edition int    `json:"edition_count"`
		} `json:"docs"`
	}
	json.Unmarshal(body, &data)
	if len(data.Docs) == 0 {
		return ok(fmt.Sprintf("No results for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Open Library — %d results for %q (sort: %s):\n\n", len(data.Docs), q, sort))
	for _, d := range data.Docs {
		sb.WriteString(fmt.Sprintf("%s by %s (%d, %d editions)\n", d.Title, d.Author, d.Year, d.Edition))
	}
	return ok(sb.String())
}
