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

var batch15 = &http.Client{Timeout: 15 * time.Second}

func HandleTVMazeCrew(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	showID, _ := getInt(args, "showId", 1)
	u := fmt.Sprintf("https://api.tvmaze.com/shows/%d/crew", showID)
	resp, apiErr := batch15.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var crew []struct {
		Type   string `json:"type"`
		Person struct {
			Name string `json:"name"`
		} `json:"person"`
	}
	json.Unmarshal(body, &crew)
	if len(crew) == 0 {
		return ok("No crew data")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Crew for show %d (%d members):\n\n", showID, len(crew)))
	for _, c := range crew {
		sb.WriteString(fmt.Sprintf("%s — %s\n", c.Person.Name, c.Type))
	}
	return ok(sb.String())
}

func HandleEvolutionChain(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 1)
	u := fmt.Sprintf("https://pokeapi.co/api/v2/evolution-chain/%d", id)
	resp, apiErr := batch15.Get(u)
	if apiErr != nil {
		return err("Evolution: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Chain struct {
			Species struct {
				Name string `json:"name"`
			} `json:"species"`
			EvolvesTo []struct {
				Species struct {
					Name string `json:"name"`
				} `json:"species"`
				EvolvesTo []struct {
					Species struct {
						Name string `json:"name"`
					} `json:"species"`
				} `json:"evolves_to"`
			} `json:"evolves_to"`
		} `json:"chain"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Evolution Chain %d:\n\n%s", id, data.Chain.Species.Name))
	if len(data.Chain.EvolvesTo) > 0 {
		for _, e := range data.Chain.EvolvesTo {
			sb.WriteString(fmt.Sprintf("\n  ↓ evolves to\n  %s", e.Species.Name))
			if len(e.EvolvesTo) > 0 {
				for _, e2 := range e.EvolvesTo {
					sb.WriteString(fmt.Sprintf("\n    ↓ evolves to\n    %s", e2.Species.Name))
				}
			}
		}
	}
	return ok(sb.String())
}

func HandlePokemonSpecies(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%s", strings.ToLower(q))
	resp, apiErr := batch15.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Species %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name        string `json:"name"`
		CaptureRate int    `json:"capture_rate"`
		Color       struct {
			Name string `json:"name"`
		} `json:"color"`
		EggGroups []struct {
			Name string `json:"name"`
		} `json:"egg_groups"`
		Habitat struct {
			Name string `json:"name"`
		} `json:"habitat"`
		BaseHappiness int `json:"base_happiness"`
		GrowthRate    struct {
			Name string `json:"name"`
		} `json:"growth_rate"`
		Shape struct {
			Name string `json:"name"`
		} `json:"shape"`
		IsLegendary bool `json:"is_legendary"`
		IsMythical  bool `json:"is_mythical"`
	}
	json.Unmarshal(body, &data)
	var eggs []string
	for _, g := range data.EggGroups {
		eggs = append(eggs, g.Name)
	}
	legendary := ""
	if data.IsLegendary {
		legendary = " ⭐ LEGENDARY"
	}
	if data.IsMythical {
		legendary = " 🌟 MYTHICAL"
	}
	return ok(fmt.Sprintf("Species: %s%s\nColor: %s | Shape: %s\nCapture Rate: %d/%d | Happiness: %d\nGrowth Rate: %s\nHabitat: %s\nEgg Groups: %s",
		data.Name, legendary, data.Color.Name, data.Shape.Name, data.CaptureRate, 255, data.BaseHappiness,
		data.GrowthRate.Name, data.Habitat.Name, strings.Join(eggs, ", ")))
}

func HandleWorldBankCountries(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	region, _ := getString(args, "region")
	n, _ := getInt(args, "limit", 20)
	u := "https://api.worldbank.org/v2/country?format=json&per_page=" + fmt.Sprintf("%d", n)
	if region != "" {
		u += "&region=" + region
	}
	resp, apiErr := batch15.Get(u)
	if apiErr != nil {
		return err("World Bank: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw []interface{}
	json.Unmarshal(body, &raw)
	if len(raw) < 2 {
		return ok("No country data")
	}
	entries, _ := raw[1].([]interface{})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("World Bank Countries (%d):\n\n", len(entries)))
	for _, e := range entries {
		entry, _ := e.(map[string]interface{})
		if entry == nil {
			continue
		}
		name, _ := entry["name"].(string)
		code, _ := entry["iso2Code"].(string)
		region, _ := entry["region"].(map[string]interface{})
		regName, _ := region["value"].(string)
		sb.WriteString(fmt.Sprintf("%s (%s) — %s\n", name, code, regName))
	}
	return ok(sb.String())
}

func HandleWorldBankIndicators(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	n, _ := getInt(args, "limit", 20)
	u := fmt.Sprintf("https://api.worldbank.org/v2/indicator?format=json&per_page=%d", n)
	resp, apiErr := batch15.Get(u)
	if apiErr != nil {
		return err("World Bank: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw []interface{}
	json.Unmarshal(body, &raw)
	if len(raw) < 2 {
		return ok("No indicators")
	}
	entries, _ := raw[1].([]interface{})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("World Bank Indicators (%d):\n\n", len(entries)))
	for _, e := range entries {
		entry, _ := e.(map[string]interface{})
		if entry == nil {
			continue
		}
		id, _ := entry["id"].(string)
		name, _ := entry["name"].(string)
		sb.WriteString(fmt.Sprintf("%s — %s\n", id, name))
	}
	return ok(sb.String())
}

func HandleNASAAstronomy(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	date, _ := getString(args, "date")
	count, _ := getInt(args, "count", 1)
	u := fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY&count=%d", count)
	if date != "" {
		u = fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY&date=%s", date)
	}
	resp, apiErr := batch15.Get(u)
	if apiErr != nil {
		return ok("APOD unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if count > 1 && date == "" {
		var items []struct {
			Title       string `json:"title"`
			Date        string `json:"date"`
			Explanation string `json:"explanation"`
			URL         string `json:"url"`
		}
		json.Unmarshal(body, &items)
		var sb strings.Builder
		for _, item := range items {
			sb.WriteString(fmt.Sprintf("APOD: %s (%s)\n%s\n%s\n\n", item.Title, item.Date, item.Explanation[:min(len(item.Explanation), 200)], item.URL))
		}
		return ok(sb.String())
	}
	var data struct {
		Title       string `json:"title"`
		Date        string `json:"date"`
		Explanation string `json:"explanation"`
		URL         string `json:"url"`
		HdURL       string `json:"hdurl"`
		Copyright   string `json:"copyright"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("APOD — %s (%s)\n© %s\n\n%s\n\n%s\n%s", data.Title, data.Date, data.Copyright, data.Explanation, data.URL, data.HdURL))
}

func HandlePokemonGrowth(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required (e.g. 'slow', 'medium', 'fast')")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/growth-rate/%s", strings.ToLower(q))
	resp, apiErr := batch15.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Growth rate %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Formula string `json:"formula"`
		Levels  []struct {
			Level      int `json:"level"`
			Experience int `json:"experience"`
		} `json:"levels"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Growth Rate: %s (%s)\n\n", data.Name, data.Formula))
	for _, l := range data.Levels {
		sb.WriteString(fmt.Sprintf("Level %d: %d XP\n", l.Level, l.Experience))
	}
	return ok(sb.String())
}
