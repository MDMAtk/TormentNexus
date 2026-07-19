package mcpimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var batch20 = &http.Client{Timeout: 15 * time.Second}

func HandleBerryFirmness(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id required")
	}
	q := name
	if q == "" {
		q = strconv.Itoa(id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/berry-firmness/%s", strings.ToLower(q))
	resp, apiErr := batch20.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Firmness %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Berries []struct {
			Name string `json:"name"`
		} `json:"berries"`
	}
	json.Unmarshal(body, &data)
	var berries []string
	for _, b := range data.Berries {
		berries = append(berries, b.Name)
	}
	return ok(fmt.Sprintf("Berry Firmness: %s (%d berries)\n%s", data.Name, len(berries), strings.Join(berries, ", ")))
}

func HandleBerryFlavor(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id required")
	}
	q := name
	if q == "" {
		q = strconv.Itoa(id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/berry-flavor/%s", strings.ToLower(q))
	resp, apiErr := batch20.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Flavor %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Contest struct {
			Name string `json:"name"`
		} `json:"contest_type"`
		Berries []struct {
			Berry struct {
				Name string `json:"name"`
			} `json:"berry"`
			Potency int `json:"potency"`
		} `json:"berries"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Berry Flavor: %s (Contest: %s)\n\nBerries:\n", data.Name, data.Contest.Name))
	for _, b := range data.Berries {
		sb.WriteString(fmt.Sprintf("%s (potency: %d)\n", b.Berry.Name, b.Potency))
	}
	return ok(sb.String())
}

func HandlePokemonCharacteristic(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 1)
	u := fmt.Sprintf("https://pokeapi.co/api/v2/characteristic/%d", id)
	resp, apiErr := batch20.Get(u)
	if apiErr != nil {
		return err("Characteristic: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ID     int   `json:"id"`
		Modulo int   `json:"gene_modulo"`
		Values []int `json:"possible_values"`
		Stat   struct {
			Name string `json:"name"`
		} `json:"highest_stat"`
	}
	json.Unmarshal(body, &data)
	var vals []string
	for _, v := range data.Values {
		vals = append(vals, strconv.Itoa(v))
	}
	return ok(fmt.Sprintf("Characteristic #%d: %s (%s mod %d)\nPossible values: %s", data.ID, data.Stat.Name, data.Stat.Name, data.Modulo, strings.Join(vals, ", ")))
}

func HandlePokemonColor(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required (e.g. 'red', 'blue', 'green')")
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-color/%s", strings.ToLower(name))
	resp, apiErr := batch20.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Color %q not found", name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Species []struct {
			Name string `json:"name"`
		} `json:"pokemon_species"`
	}
	json.Unmarshal(body, &data)
	var species []string
	for _, s := range data.Species {
		species = append(species, s.Name)
	}
	return ok(fmt.Sprintf("Color: %s (%d species)\n\n%s\n\n(Use query param for full list)", data.Name, len(species), strings.Join(species[:15], ", ")))
}

func HandlePokemonHabitat(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required (e.g. 'cave', 'forest', 'sea')")
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-habitat/%s", strings.ToLower(name))
	resp, apiErr := batch20.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Habitat %q not found", name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Species []struct {
			Name string `json:"name"`
		} `json:"pokemon_species"`
	}
	json.Unmarshal(body, &data)
	var species []string
	for _, s := range data.Species {
		species = append(species, s.Name)
	}
	return ok(fmt.Sprintf("Habitat: %s (%d species)\n\n%s\n\n(Use query param for full list)", data.Name, len(species), strings.Join(species[:10], ", ")))
}

func HandlePokemonShape(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required (e.g. 'quadruped', 'bipedal', 'fish')")
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-shape/%s", strings.ToLower(name))
	resp, apiErr := batch20.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Shape %q not found", name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Species []struct {
			Name string `json:"name"`
		} `json:"pokemon_species"`
	}
	json.Unmarshal(body, &data)
	var species []string
	for _, s := range data.Species {
		species = append(species, s.Name)
	}
	return ok(fmt.Sprintf("Shape: %s (%d species)\n\n%s\n\n(Use query param for full list)", data.Name, len(species), strings.Join(species[:10], ", ")))
}

func HandleTVUpdates(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	since, _ := getString(args, "since")
	if since == "" {
		since = "day"
	}
	u := fmt.Sprintf("https://api.tvmaze.com/updates/shows?since=%s", since)
	resp, apiErr := batch20.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw map[string]int64
	json.Unmarshal(body, &raw)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TV Show Updates (%d total):\n\n", len(raw)))
	count := 0
	for showID, timestamp := range raw {
		if count >= 20 {
			break
		}
		t := time.Unix(timestamp, 0)
		u2 := fmt.Sprintf("https://api.tvmaze.com/shows/%s", showID)
		resp2, _ := batch20.Get(u2)
		if resp2 != nil {
			defer resp2.Body.Close()
			body2, _ := io.ReadAll(resp2.Body)
			var show struct {
				Name string `json:"name"`
			}
			json.Unmarshal(body2, &show)
			sb.WriteString(fmt.Sprintf("%s: %s\n", show.Name, t.Format("15:04:05")))
		} else {
			sb.WriteString(fmt.Sprintf("Show #%s: %s\n", showID, t.Format("15:04:05")))
		}
		count++
	}
	return ok(sb.String())
}

func HandleOpenLibraryAvailability(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	workID, _ := getString(args, "workId")
	if workID == "" {
		return err("workId is required (e.g. OL45804W)")
	}
	u := fmt.Sprintf("https://openlibrary.org/works/%s/availability.json", workID)
	resp, apiErr := batch20.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Availability for %s unavailable", workID))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw map[string]interface{}
	json.Unmarshal(body, &raw)
	return ok(fmt.Sprintf("Availability for %s:\n%v", workID, raw))
}
