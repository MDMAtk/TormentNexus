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

var batch8 = &http.Client{Timeout: 15 * time.Second}

func HandleDadJoke(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "search")
	u := "https://icanhazdadjoke.com/"
	if q != "" {
		u = fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s&limit=3", url.QueryEscape(q))
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "TormentNexus/1.0")
	resp, apiErr := batch8.Do(req)
	if apiErr != nil {
		return ok("Why don't scientists trust atoms? They make up everything!")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if q != "" {
		var data struct {
			Results []struct {
				Joke string `json:"joke"`
			} `json:"results"`
		}
		json.Unmarshal(body, &data)
		if len(data.Results) == 0 {
			return ok(fmt.Sprintf("No dad jokes found for %q", q))
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Dad Jokes for %q:\n\n", q))
		for _, r := range data.Results {
			sb.WriteString(fmt.Sprintf("%s\n\n", r.Joke))
		}
		return ok(sb.String())
	}
	var data struct {
		Joke string `json:"joke"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Dad Joke: %s", data.Joke))
}

func HandleCatPicture(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	tag, _ := getString(args, "tag")
	u := "https://cataas.com/cat?json=true"
	if tag != "" {
		u = fmt.Sprintf("https://cataas.com/cat/%s?json=true", url.PathEscape(tag))
	}
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return ok("Cat picture: https://cataas.com/cat")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		URL string `json:"url"`
	}
	json.Unmarshal(body, &data)
	catURL := data.URL
	if !strings.HasPrefix(catURL, "http") {
		catURL = "https://cataas.com" + catURL
	}
	return ok(fmt.Sprintf("🐱 %s", catURL))
}

func HandleHTTPBinDebug(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	kind, _ := getString(args, "kind")
	if kind == "" {
		kind = "ip"
	}
	u := fmt.Sprintf("https://httpbin.org/%s", kind)
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return err("HTTPBin: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	pretty, _ := json.MarshalIndent(result, "", "  ")
	return ok(fmt.Sprintf("HTTPBin [/%s]:\n%s", kind, string(pretty)))
}

func HandleUniversitySearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	country, _ := getString(args, "country")
	if name == "" && country == "" {
		return err("name or country is required")
	}
	u := "http://universities.hipolabs.com/search?"
	if name != "" {
		u += "name=" + url.QueryEscape(name)
	}
	if country != "" {
		if name != "" {
			u += "&"
		}
		u += "country=" + url.QueryEscape(country)
	}
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return err("Universities: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var unis []struct {
		Name     string   `json:"name"`
		Country  string   `json:"country"`
		WebPages []string `json:"web_pages"`
		Domains  []string `json:"domains"`
	}
	json.Unmarshal(body, &unis)
	if len(unis) == 0 {
		return ok(fmt.Sprintf("No universities found for %q in %s", name, country))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Universities (%d found):\n\n", len(unis)))
	for _, uni := range unis {
		sb.WriteString(fmt.Sprintf("%s (%s)\n", uni.Name, uni.Country))
		if len(uni.WebPages) > 0 {
			sb.WriteString(fmt.Sprintf("  %s\n", uni.WebPages[0]))
		}
	}
	return ok(sb.String())
}

func HandleAuthorWorks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	authorKey, _ := getString(args, "authorKey")
	name, _ := getString(args, "name")
	limit, _ := getInt(args, "limit", 10)
	if authorKey == "" && name == "" {
		return err("authorKey or name is required")
	}
	if authorKey == "" && name != "" {
		searchURL := fmt.Sprintf("https://openlibrary.org/search/authors.json?q=%s&limit=1", url.QueryEscape(name))
		resp, apiErr := batch8.Get(searchURL)
		if apiErr != nil {
			return err("Open Library: " + apiErr.Error())
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		var search struct {
			Docs []struct {
				Key  string `json:"key"`
				Name string `json:"name"`
			} `json:"docs"`
		}
		json.Unmarshal(body, &search)
		if len(search.Docs) == 0 {
			return ok(fmt.Sprintf("Author %q not found", name))
		}
		authorKey = search.Docs[0].Key
	}
	u := fmt.Sprintf("https://openlibrary.org%s/works.json?limit=%d", authorKey, limit)
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Entries []struct {
			Title string `json:"title"`
		} `json:"entries"`
	}
	json.Unmarshal(body, &data)
	if len(data.Entries) == 0 {
		return ok("No works found")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Works by %s (%d):\n\n", authorKey, len(data.Entries)))
	for _, e := range data.Entries {
		sb.WriteString(fmt.Sprintf("%s\n", e.Title))
	}
	return ok(sb.String())
}

func HandleBerryDetail(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/berry/%s", strings.ToLower(q))
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Berry %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name       string `json:"name"`
		Size       int    `json:"size"`
		GrowthTime int    `json:"growth_time"`
		Smoothness int    `json:"smoothness"`
		Firmness   struct {
			Name string `json:"name"`
		} `json:"firmness"`
		Item struct {
			Name string `json:"name"`
		} `json:"item"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Berry: %s\nItem: %s\nFirmness: %s\nSize: %d | Growth: %d hours | Smoothness: %d",
		data.Name, data.Item.Name, data.Firmness.Name, data.Size, data.GrowthTime, data.Smoothness))
}

func HandleBookRating(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	key, _ := getString(args, "key")
	isbn, _ := getString(args, "isbn")
	if key == "" && isbn == "" {
		return err("key (e.g. /works/OL45804W) or isbn is required")
	}
	workKey := key
	if workKey == "" && isbn != "" {
		workKey = "/isbn/" + isbn
	}
	u := fmt.Sprintf("https://openlibrary.org%s/ratings.json", workKey)
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Rating for %s unavailable", key+isbn))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Summary struct {
			Average float64 `json:"average"`
			Count   int     `json:"count"`
		} `json:"summary"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Rating for %s:\nAverage: %.2f (%d ratings)", workKey, data.Summary.Average, data.Summary.Count))
}

func HandlePokemonAbility(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required (e.g. 'overgrow', 'levitate')")
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/ability/%s", strings.ToLower(name))
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Ability %q not found", name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name   string `json:"name"`
		Gen    string `json:"generation"`
		Flavor []struct {
			Text string `json:"flavor_text"`
		} `json:"flavor_text_entries"`
	}
	json.Unmarshal(body, &data)
	flavor := ""
	if len(data.Flavor) > 0 {
		flavor = data.Flavor[0].Text
		if len(flavor) > 200 {
			flavor = flavor[:200] + "..."
		}
	}
	return ok(fmt.Sprintf("Ability: %s\nGeneration: %s\n\n%s", data.Name, data.Gen, flavor))
}

func HandleStarWarsSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	category, _ := getString(args, "category")
	search, _ := getString(args, "search")
	if category == "" {
		category = "people"
	}
	if search == "" {
		return err("search is required")
	}
	u := fmt.Sprintf("https://swapi.dev/api/%s/?search=%s", category, url.QueryEscape(search))
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return err("SWAPI: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Results []struct {
			Name      string `json:"name"`
			BirthYear string `json:"birth_year"`
			Gender    string `json:"gender"`
		} `json:"results"`
	}
	json.Unmarshal(body, &data)
	if len(data.Results) == 0 {
		return ok(fmt.Sprintf("No %s found for %q", category, search))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Star Wars %s matching %q:\n\n", category, search))
	for _, r := range data.Results {
		sb.WriteString(fmt.Sprintf("%s (%s, %s)\n", r.Name, r.BirthYear, r.Gender))
	}
	return ok(sb.String())
}

func HandleSWAPILookup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	category, _ := getString(args, "category")
	id, _ := getInt(args, "id", 1)
	if category == "" {
		category = "people"
	}
	u := fmt.Sprintf("https://swapi.dev/api/%s/%d/", category, id)
	resp, apiErr := batch8.Get(u)
	if apiErr != nil {
		return err("SWAPI: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Star Wars %s #%d\n", category, id))
	for k, v := range result {
		if k == "created" || k == "edited" || k == "url" {
			continue
		}
		switch val := v.(type) {
		case string:
			if len(val) < 120 {
				sb.WriteString(fmt.Sprintf("%s: %s\n", k, val))
			}
		case float64:
			sb.WriteString(fmt.Sprintf("%s: %.0f\n", k, val))
		case []interface{}:
			if len(val) > 0 && len(val) <= 3 {
				var items []string
				for _, item := range val {
					items = append(items, fmt.Sprintf("%v", item))
				}
				sb.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(items, ", ")))
			}
		}
	}
	return ok(sb.String())
}
