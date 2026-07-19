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

var batch18 = &http.Client{Timeout: 15 * time.Second}

func HandlePokemonGeneration(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required (e.g. 'generation-i', 1)")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/generation/%s", strings.ToLower(q))
	resp, apiErr := batch18.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Generation %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		MainRegion struct {
			Name string `json:"name"`
		} `json:"main_region"`
		Pokemon []struct {
			Name string `json:"name"`
		} `json:"pokemon_species"`
		Moves []struct {
			Name string `json:"name"`
		} `json:"moves"`
		Abilities []struct {
			Name string `json:"name"`
		} `json:"abilities"`
		Types []struct {
			Name string `json:"name"`
		} `json:"types"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Generation %d: %s\nRegion: %s\nPokemon: %d | Moves: %d | Abilities: %d | Types: %d",
		data.ID, data.Name, data.MainRegion.Name, len(data.Pokemon), len(data.Moves), len(data.Abilities), len(data.Types)))
}

func HandlePokedex(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required (e.g. 'kanto', 'national')")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/pokedex/%s", strings.ToLower(q))
	resp, apiErr := batch18.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Pokedex %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		IsMain  bool   `json:"is_main_series"`
		Entries []struct {
			EntryNumber int `json:"entry_number"`
			Species     struct {
				Name string `json:"name"`
			} `json:"pokemon_species"`
		} `json:"pokemon_entries"`
		Descriptions []struct {
			Text string `json:"text"`
		} `json:"descriptions"`
	}
	json.Unmarshal(body, &data)
	desc := ""
	if len(data.Descriptions) > 0 {
		desc = data.Descriptions[0].Text
		if len(desc) > 200 {
			desc = desc[:200] + "..."
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pokedex: %s (main series: %v)\n%s\n\n", data.Name, data.IsMain, desc))
	for _, e := range data.Entries[:10] {
		sb.WriteString(fmt.Sprintf("#%03d %s\n", e.EntryNumber, e.Species.Name))
	}
	return ok(sb.String())
}

func HandlePokemonVersion(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required (e.g. 'red', 'blue', 'scarlet')")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/version/%s", strings.ToLower(q))
	resp, apiErr := batch18.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Version %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name  string `json:"name"`
		Group struct {
			Name string `json:"name"`
		} `json:"version_group"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Version: %s (%s)", data.Name, data.Group.Name))
}

func HandleVersionGroup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/version-group/%s", strings.ToLower(q))
	resp, apiErr := batch18.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Version group %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name string `json:"name"`
		Gen  struct {
			Name string `json:"name"`
		} `json:"generation"`
		Moves []struct {
			Name string `json:"name"`
		} `json:"move_learn_methods"`
		Versions []struct {
			Name string `json:"name"`
		} `json:"versions"`
	}
	json.Unmarshal(body, &data)
	var versions []string
	for _, v := range data.Versions {
		versions = append(versions, v.Name)
	}
	var methods []string
	for _, m := range data.Moves {
		methods = append(methods, m.Name)
	}
	return ok(fmt.Sprintf("Version Group: %s\nGeneration: %s\nGames: %s\nLearn Methods: %s",
		data.Name, data.Gen.Name, strings.Join(versions, ", "), strings.Join(methods, ", ")))
}

func HandlePokeResourceList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resource, _ := getString(args, "resource")
	n, _ := getInt(args, "limit", 20)
	if resource == "" {
		return err("resource is required (e.g. 'ability', 'type', 'move', 'pokemon', 'item')")
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/%s?limit=%d", resource, n)
	resp, apiErr := batch18.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Resource %q not found", resource))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Count   int `json:"count"`
		Results []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pokemon %s (%d total, showing %d):\n\n", resource, data.Count, len(data.Results)))
	for _, r := range data.Results {
		sb.WriteString(fmt.Sprintf("%s\n", r.Name))
	}
	return ok(sb.String())
}

func HandleNASAAsteroids(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	start, _ := getString(args, "startDate")
	end, _ := getString(args, "endDate")
	if start == "" {
		start = time.Now().Format("2006-01-02")
	}
	if end == "" {
		end = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	}
	u := fmt.Sprintf("https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=DEMO_KEY", start, end)
	resp, apiErr := batch18.Get(u)
	if apiErr != nil {
		return ok("NEO feed unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw map[string]interface{}
	json.Unmarshal(body, &raw)
	cnt, _ := raw["element_count"].(float64)
	neo, _ := raw["near_earth_objects"].(map[string]interface{})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("NEO — %.0f asteroids (%s to %s):\n\n", cnt, start, end))
	hazardCount := 0
	for date, objs := range neo {
		items, _ := objs.([]interface{})
		for _, obj := range items {
			o, _ := obj.(map[string]interface{})
			name, _ := o["name"].(string)
			hazard, _ := o["is_potentially_hazardous_asteroid"].(bool)
			if hazard {
				hazardCount++
				sb.WriteString(fmt.Sprintf("⚠️ %s on %s\n", name, date))
			}
		}
	}
	sb.WriteString(fmt.Sprintf("\nPotentially hazardous: %d", hazardCount))
	return ok(sb.String())
}
