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

var batch17 = &http.Client{Timeout: 15 * time.Second}

func HandleTVShowAKAs(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	showID, _ := getInt(args, "showId", 1)
	u := fmt.Sprintf("https://api.tvmaze.com/shows/%d/akas", showID)
	resp, apiErr := batch17.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var akas []struct {
		Name    string `json:"name"`
		Country struct {
			Name string `json:"name"`
		} `json:"country"`
	}
	json.Unmarshal(body, &akas)
	if len(akas) == 0 {
		return ok("No alternative names")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Alternative names for show %d (%d):\n\n", showID, len(akas)))
	for _, a := range akas {
		country := ""
		if a.Country.Name != "" {
			country = " [" + a.Country.Name + "]"
		}
		sb.WriteString(fmt.Sprintf("%s%s\n", a.Name, country))
	}
	return ok(sb.String())
}

func HandleTMMachine(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 1)
	u := fmt.Sprintf("https://pokeapi.co/api/v2/machine/%d", id)
	resp, apiErr := batch17.Get(u)
	if apiErr != nil {
		return err("Machine: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ID   int `json:"id"`
		Item struct {
			Name string `json:"name"`
		} `json:"item"`
		Move struct {
			Name string `json:"name"`
		} `json:"move"`
		Group struct {
			Name string `json:"name"`
		} `json:"version_group"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("TM/HM #%d: %s — teaches %s (Gen: %s)", data.ID, data.Item.Name, data.Move.Name, data.Group.Name))
}

func HandleEggGroup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/egg-group/%s", strings.ToLower(q))
	resp, apiErr := batch17.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Egg group %q not found", q))
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
	return ok(fmt.Sprintf("Egg Group: %s (%d species)\n\n%s\n\n(Use query param for full list)", data.Name, len(species), strings.Join(species[:10], ", ")))
}

func HandlePokemonNature(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required (e.g. 'adamant', 'modest')")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/nature/%s", strings.ToLower(q))
	resp, apiErr := batch17.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Nature %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name      string `json:"name"`
		Increased struct {
			Name string `json:"name"`
		} `json:"increased_stat"`
		Decreased struct {
			Name string `json:"name"`
		} `json:"decreased_stat"`
		Likes struct {
			Name string `json:"name"`
		} `json:"likes_flavor"`
		Hates struct {
			Name string `json:"name"`
		} `json:"hates_flavor"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Nature: %s\n↑ %s | ↓ %s\nLikes: %s | Hates: %s", data.Name, data.Increased.Name, data.Decreased.Name, data.Likes.Name, data.Hates.Name))
}

func HandlePokemonGender(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required (e.g. 'male', 'female', 'genderless')")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/gender/%s", strings.ToLower(q))
	resp, apiErr := batch17.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Gender %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Species []struct {
			Species struct {
				Name string `json:"name"`
			} `json:"pokemon_species"`
			Rate int `json:"rate"`
		} `json:"pokemon_species_details"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Gender: %s (%d species)\n\n", data.Name, len(data.Species)))
	for _, s := range data.Species[:10] {
		sb.WriteString(fmt.Sprintf("%s (rate: %d)\n", s.Species.Name, s.Rate))
	}
	return ok(sb.String())
}

func HandlePokemonStat(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required (e.g. 'hp', 'attack', 'defense')")
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/stat/%s", strings.ToLower(name))
	resp, apiErr := batch17.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Stat %q not found", name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name         string `json:"name"`
		IsBattleOnly bool   `json:"is_battle_only"`
		Moves        []struct {
			Move struct {
				Name string `json:"name"`
			} `json:"move"`
		} `json:"affecting_moves"`
		Items []struct {
			Item struct {
				Name string `json:"name"`
			} `json:"item"`
		} `json:"affecting_items"`
	}
	json.Unmarshal(body, &data)
	var moves []string
	for _, m := range data.Moves {
		moves = append(moves, m.Move.Name)
	}
	return ok(fmt.Sprintf("Stat: %s\nBattle Only: %v\nAffecting Moves (%d): %s", data.Name, data.IsBattleOnly, len(moves), strings.Join(moves, ", ")))
}

func HandleNASATechport(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	// limit removed - API always returns all
	u := fmt.Sprintf("https://api.nasa.gov/techport/api/projects?api_key=DEMO_KEY")
	resp, apiErr := batch17.Get(u)
	if apiErr != nil {
		return ok("NASA Techport unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Projects []struct {
			ID     int    `json:"id"`
			Title  string `json:"title"`
			Status string `json:"status"`
		} `json:"projects"`
		Total int `json:"totalCount"`
	}
	json.Unmarshal(body, &data)
	if len(data.Projects) == 0 {
		return ok("No projects found")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("NASA Technology Projects — %d total:\n\n", data.Total))
	for _, p := range data.Projects {
		if len(sb.String()) > 1500 {
			break
		}
		sb.WriteString(fmt.Sprintf("%d: %s [%s]\n", p.ID, p.Title, p.Status))
	}
	return ok(sb.String())
}
