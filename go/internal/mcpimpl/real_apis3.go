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

var batch4API = &http.Client{Timeout: 15 * time.Second}

// ═══════════════════════════════════════════════════════════════════
// 22. POKÉAPI — Pokémon data (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandlePokemonLookup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", strings.ToLower(q))
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("PokéAPI: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Height int    `json:"height"`
		Weight int    `json:"weight"`
		Types  []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
		Abilities []struct {
			Ability struct {
				Name string `json:"name"`
			} `json:"ability"`
		} `json:"abilities"`
		Stats []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		} `json:"stats"`
		Sprites struct {
			FrontDefault string `json:"front_default"`
		} `json:"sprites"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("Pokémon %q: not found", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pokémon #%d — %s\n", data.ID, strings.Title(data.Name)))
	sb.WriteString(fmt.Sprintf("Height: %.1f m | Weight: %.1f kg\n", float64(data.Height)/10, float64(data.Weight)/10))
	var types []string
	for _, t := range data.Types {
		types = append(types, t.Type.Name)
	}
	sb.WriteString(fmt.Sprintf("Types: %s\n", strings.Join(types, ", ")))
	var abilities []string
	for _, a := range data.Abilities {
		abilities = append(abilities, a.Ability.Name)
	}
	sb.WriteString(fmt.Sprintf("Abilities: %s\n", strings.Join(abilities, ", ")))
	sb.WriteString("Base Stats:\n")
	for _, s := range data.Stats {
		bar := strings.Repeat("█", s.BaseStat/10)
		sb.WriteString(fmt.Sprintf("  %s: %d %s\n", s.Stat.Name, s.BaseStat, bar))
	}
	if data.Sprites.FrontDefault != "" {
		sb.WriteString(fmt.Sprintf("\n%s", data.Sprites.FrontDefault))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 23. SWAPI — Star Wars data (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleStarWarsLookup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	category, _ := getString(args, "category")
	id, _ := getInt(args, "id", 1)
	search, _ := getString(args, "search")
	if category == "" {
		category = "people"
	}
	var u string
	if search != "" {
		u = fmt.Sprintf("https://swapi.dev/api/%s/?search=%s", category, url.QueryEscape(search))
	} else {
		u = fmt.Sprintf("https://swapi.dev/api/%s/%d/", category, id)
	}
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("SWAPI: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if search != "" {
		var data struct {
			Results []map[string]interface{} `json:"results"`
		}
		json.Unmarshal(body, &data)
		if len(data.Results) == 0 {
			return ok(fmt.Sprintf("No Star Wars %s found for %q", category, search))
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Star Wars %s — %d results:\n\n", category, len(data.Results)))
		for i, r := range data.Results {
			name, _ := r["name"].(string)
			if name == "" {
				name, _ = r["title"].(string)
			}
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, name))
			if i >= 9 {
				break
			}
		}
		return ok(sb.String())
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Star Wars — %s #%d\n", category, id))
	for k, v := range result {
		if k == "created" || k == "edited" || k == "url" {
			continue
		}
		switch val := v.(type) {
		case string:
			if len(val) < 100 {
				sb.WriteString(fmt.Sprintf("%s: %s\n", k, val))
			}
		case float64:
			sb.WriteString(fmt.Sprintf("%s: %.0f\n", k, val))
		case []interface{}:
			if len(val) > 0 && len(val) < 5 {
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

// ═══════════════════════════════════════════════════════════════════
// 24. RICK AND MORTY API — Character/location/episode data
// ═══════════════════════════════════════════════════════════════════

func HandleRickAndMorty(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	character, _ := getInt(args, "character", 1)
	name, _ := getString(args, "name")
	var u string
	if name != "" {
		u = fmt.Sprintf("https://rickandmortyapi.com/api/character/?name=%s", url.QueryEscape(name))
	} else {
		u = fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", character)
	}
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("Rick&Morty: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if name != "" {
		var data struct {
			Results []map[string]interface{} `json:"results"`
		}
		json.Unmarshal(body, &data)
		if len(data.Results) == 0 {
			return ok(fmt.Sprintf("No Rick & Morty characters for %q", name))
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Rick & Morty — %d characters:\n\n", len(data.Results)))
		for i, r := range data.Results {
			if i >= 10 {
				break
			}
			n, _ := r["name"].(string)
			status, _ := r["status"].(string)
			species, _ := r["species"].(string)
			gender, _ := r["gender"].(string)
			img, _ := r["image"].(string)
			sb.WriteString(fmt.Sprintf("%d. %s — %s %s (%s)\n", i+1, n, status, species, gender))
			_ = img
		}
		return ok(sb.String())
	}

	var data struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Species string `json:"species"`
		Gender  string `json:"gender"`
		Origin  struct {
			Name string `json:"name"`
		} `json:"origin"`
		Location struct {
			Name string `json:"name"`
		} `json:"location"`
		Episode []string `json:"episode"`
		Image   string   `json:"image"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Rick & Morty — %s\nStatus: %s | Species: %s | Gender: %s\nOrigin: %s | Location: %s\nEpisodes: %d\n%s",
		data.Name, data.Status, data.Species, data.Gender, data.Origin.Name, data.Location.Name, len(data.Episode), data.Image))
}

// ═══════════════════════════════════════════════════════════════════
// 25. TVMAZE — TV show search and episode info (free)
// ═══════════════════════════════════════════════════════════════════

func HandleTVShowSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	n, _ := getInt(args, "limit", 5)
	if q == "" {
		return err("query is required")
	}
	u := fmt.Sprintf("https://api.tvmaze.com/search/shows?q=%s", url.QueryEscape(q))
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var results []struct {
		Score float64 `json:"score"`
		Show  struct {
			ID       int      `json:"id"`
			Name     string   `json:"name"`
			Type     string   `json:"type"`
			Language string   `json:"language"`
			Genres   []string `json:"genres"`
			Status   string   `json:"status"`
			Rating   struct {
				Average float64 `json:"average"`
			} `json:"rating"`
			Summary string `json:"summary"`
			URL     string `json:"url"`
		} `json:"show"`
	}
	json.Unmarshal(body, &results)
	if len(results) == 0 {
		return ok(fmt.Sprintf("No shows found for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TV Shows — %d results for %q:\n\n", len(results), q))
	for i, r := range results {
		if i >= n {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. %s (score: %.2f)\n", i+1, r.Show.Name, r.Score))
		sb.WriteString(fmt.Sprintf("   Type: %s | Language: %s | Status: %s\n", r.Show.Type, r.Show.Language, r.Show.Status))
		if len(r.Show.Genres) > 0 {
			sb.WriteString(fmt.Sprintf("   Genres: %s\n", strings.Join(r.Show.Genres, ", ")))
		}
		if r.Show.Rating.Average > 0 {
			sb.WriteString(fmt.Sprintf("   Rating: %.1f/10\n", r.Show.Rating.Average))
		}
		summary := stripHTMLTags(r.Show.Summary)
		if len(summary) > 200 {
			summary = summary[:200] + "..."
		}
		sb.WriteString(fmt.Sprintf("   %s\n", summary))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 26. AGIFY — Age prediction from name (free)
// ═══════════════════════════════════════════════════════════════════

func HandleAgePredict(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required")
	}
	u := fmt.Sprintf("https://api.agify.io?name=%s", url.QueryEscape(name))
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("Agify: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Count int    `json:"count"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Age prediction for \"%s\": %d years old (based on %d people)", data.Name, data.Age, data.Count))
}

// ═══════════════════════════════════════════════════════════════════
// 27. NATIONALIZE — Nationality prediction from name (free)
// ═══════════════════════════════════════════════════════════════════

func HandleNationalityPredict(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required")
	}
	u := fmt.Sprintf("https://api.nationalize.io?name=%s", url.QueryEscape(name))
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("Nationalize: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name    string `json:"name"`
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Nationality prediction for \"%s\":\n", data.Name))
	for _, c := range data.Country {
		prob := c.Probability * 100
		sb.WriteString(fmt.Sprintf("  %s: %.1f%%\n", c.CountryID, prob))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 28. OPEN-METEO — Weather forecast (free, no API key)
// ═══════════════════════════════════════════════════════════════════

func HandleForecast(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude")
	if lat == 0 {
		lat = 52.52
	}
	lon := getFloat(args, "longitude")
	if lon == 0 {
		lon = 13.41
	}
	days, _ := getInt(args, "days", 3)
	u := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&daily=temperature_2m_max,temperature_2m_min,precipitation_sum,weathercode&current_weather=true&forecast_days=%d",
		lat, lon, days)
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Forecast at %.2f, %.2f: service unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		CurrentWeather struct {
			Temperature float64 `json:"temperature"`
			Windspeed   float64 `json:"windspeed"`
			Weathercode int     `json:"weathercode"`
		} `json:"current_weather"`
		Daily struct {
			Time          []string  `json:"time"`
			TempMax       []float64 `json:"temperature_2m_max"`
			TempMin       []float64 `json:"temperature_2m_min"`
			Precipitation []float64 `json:"precipitation_sum"`
			WeatherCode   []int     `json:"weathercode"`
		} `json:"daily"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Forecast at %.2f, %.2f\n\n", lat, lon))
	sb.WriteString(fmt.Sprintf("Current: %.1f°C, Wind: %.1f km/h, %s\n\n", data.CurrentWeather.Temperature, data.CurrentWeather.Windspeed, weatherCodeDesc(data.CurrentWeather.Weathercode)))
	for i := range data.Daily.Time {
		sb.WriteString(fmt.Sprintf("%s: %.1f–%.1f°C, %.0fmm precip, %s\n",
			data.Daily.Time[i], data.Daily.TempMin[i], data.Daily.TempMax[i],
			data.Daily.Precipitation[i], weatherCodeDesc(data.Daily.WeatherCode[i])))
	}
	return ok(sb.String())
}

func weatherCodeDesc(code int) string {
	desc := map[int]string{
		0: "☀️ Clear", 1: "🌤 Mostly clear", 2: "⛅ Partly cloudy", 3: "☁️ Overcast",
		45: "🌫 Foggy", 48: "🌫 Depositing rime",
		51: "🌦 Light drizzle", 53: "🌦 Moderate drizzle", 55: "🌦 Dense drizzle",
		61: "🌧 Light rain", 63: "🌧 Moderate rain", 65: "🌧 Heavy rain",
		71: "🌨 Light snow", 73: "🌨 Moderate snow", 75: "🌨 Heavy snow",
		80: "🌦 Slight rain showers", 81: "🌧 Moderate rain showers", 82: "🌧 Violent rain showers",
		95: "⛈ Thunderstorm", 96: "⛈ Thunderstorm with hail", 99: "⛈ Thunderstorm with heavy hail",
	}
	if d, exists := desc[code]; exists {
		return d
	}
	return fmt.Sprintf("Code %d", code)
}

// ═══════════════════════════════════════════════════════════════════
// 29. RANDOM USER — Generate random user profiles
// ═══════════════════════════════════════════════════════════════════

func HandleRandomUser(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	n, _ := getInt(args, "count", 1)
	gender, _ := getString(args, "gender")
	nat, _ := getString(args, "nationality")
	u := fmt.Sprintf("https://randomuser.me/api/?results=%d", n)
	if gender != "" {
		u += "&gender=" + gender
	}
	if nat != "" {
		u += "&nat=" + nat
	}
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("RandomUser: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Results []struct {
			Name struct {
				Title string `json:"title"`
				First string `json:"first"`
				Last  string `json:"last"`
			} `json:"name"`
			Gender   string `json:"gender"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Nat      string `json:"nat"`
			Location struct {
				Street struct {
					Name   string `json:"name"`
					Number int    `json:"number"`
				} `json:"street"`
				City    string `json:"city"`
				State   string `json:"state"`
				Country string `json:"country"`
			} `json:"location"`
			Picture struct {
				Large string `json:"large"`
			} `json:"picture"`
		} `json:"results"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	for i, u := range data.Results {
		sb.WriteString(fmt.Sprintf("User %d:\n", i+1))
		sb.WriteString(fmt.Sprintf("  Name: %s %s %s\n", u.Name.Title, u.Name.First, u.Name.Last))
		sb.WriteString(fmt.Sprintf("  Gender: %s | Nationality: %s\n", u.Gender, u.Nat))
		sb.WriteString(fmt.Sprintf("  Email: %s\n", u.Email))
		sb.WriteString(fmt.Sprintf("  Phone: %s\n", u.Phone))
		sb.WriteString(fmt.Sprintf("  Location: %d %s, %s, %s %s\n", u.Location.Street.Number, u.Location.Street.Name, u.Location.City, u.Location.State, u.Location.Country))
		sb.WriteString(fmt.Sprintf("  Photo: %s\n", u.Picture.Large))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 30. ZIPPO — Zip code lookup (free)
// ═══════════════════════════════════════════════════════════════════

func HandleZipLookup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	code, _ := getString(args, "code")
	country, _ := getString(args, "country")
	if code == "" {
		return err("code is required")
	}
	if country == "" {
		country = "us"
	}
	u := fmt.Sprintf("https://api.zippopotam.us/%s/%s", strings.ToLower(country), code)
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Zip code %s (%s): not found or service unavailable", code, country))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		PostCode string `json:"post code"`
		Country  string `json:"country"`
		Places   []struct {
			PlaceName string `json:"place name"`
			State     string `json:"state"`
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"places"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Zip code: %s (%s)\n", data.PostCode, data.Country))
	for _, p := range data.Places {
		sb.WriteString(fmt.Sprintf("  %s, %s (lat: %s, lon: %s)\n", p.PlaceName, p.State, p.Latitude, p.Longitude))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 31. FRUITYVICE — Fruit nutrition data (free)
// ═══════════════════════════════════════════════════════════════════

func HandleFruitInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required (e.g. apple, banana, strawberry)")
	}
	u := fmt.Sprintf("https://www.fruityvice.com/api/fruit/%s", strings.ToLower(name))
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Fruit %q: not found. Try: apple, banana, orange, strawberry, blueberry", name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name       string `json:"name"`
		Family     string `json:"family"`
		Order      string `json:"order"`
		Genus      string `json:"genus"`
		Nutritions struct {
			Calories      float64 `json:"calories"`
			Fat           float64 `json:"fat"`
			Sugar         float64 `json:"sugar"`
			Carbohydrates float64 `json:"carbohydrates"`
			Protein       float64 `json:"protein"`
		} `json:"nutritions"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Fruit: %s\nFamily: %s | Order: %s | Genus: %s\n\nNutrition per 100g:\nCalories: %.0f\nFat: %.1fg\nSugar: %.1fg\nCarbs: %.1fg\nProtein: %.1fg",
		data.Name, data.Family, data.Order, data.Genus,
		data.Nutritions.Calories, data.Nutritions.Fat, data.Nutritions.Sugar,
		data.Nutritions.Carbohydrates, data.Nutritions.Protein))
}

// ═══════════════════════════════════════════════════════════════════
// 32. CHEAPSHARK — Game deals (free)
// ═══════════════════════════════════════════════════════════════════

func HandleGameDeals(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	title, _ := getString(args, "title")
	storeID, _ := getInt(args, "storeID", 0)
	maxPrice := getFloat(args, "maxPrice")
	u := "https://www.cheapshark.com/api/1.0/deals?pageSize=10"
	if title != "" {
		u += "&title=" + url.QueryEscape(title)
	}
	if storeID > 0 {
		u += fmt.Sprintf("&storeID=%d", storeID)
	}
	if maxPrice > 0 {
		u += fmt.Sprintf("&upperPrice=%.2f", maxPrice)
	}
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("CheapShark: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var deals []struct {
		Title           string `json:"title"`
		SalePrice       string `json:"salePrice"`
		NormalPrice     string `json:"normalPrice"`
		Savings         string `json:"savings"`
		MetacriticScore string `json:"metacriticScore"`
		StoreID         int    `json:"storeID"`
		Thumb           string `json:"thumb"`
	}
	json.Unmarshal(body, &deals)
	if len(deals) == 0 {
		return ok("No game deals found")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Game Deals (%d found):\n\n", len(deals)))
	for i, d := range deals {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, d.Title))
		sb.WriteString(fmt.Sprintf("   💰 $%s → $%s (save %s%%)\n", d.NormalPrice, d.SalePrice, d.Savings))
		if d.MetacriticScore != "" {
			sb.WriteString(fmt.Sprintf("   Metacritic: %s\n", d.MetacriticScore))
		}
		_ = d.Thumb
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 33. JSONPLACEHOLDER — Sample data/TODO (free)
// ═══════════════════════════════════════════════════════════════════

func HandleJSONPlaceholder(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resource, _ := getString(args, "resource")
	id, _ := getInt(args, "id", 0)
	if resource == "" {
		resource = "todos"
	}
	u := fmt.Sprintf("https://jsonplaceholder.typicode.com/%s", resource)
	if id > 0 {
		u += fmt.Sprintf("/%d", id)
	}
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("JSONPlaceholder: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result interface{}
	json.Unmarshal(body, &result)
	pretty, _ := json.MarshalIndent(result, "", "  ")
	return ok(fmt.Sprintf("JSONPlaceholder — %s:\n\n%s", resource, string(pretty)))
}

// ═══════════════════════════════════════════════════════════════════
// 34. OPEN TRIVIA DB — Quiz questions (free)
// ═══════════════════════════════════════════════════════════════════

func HandleTriviaQuestion(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	n, _ := getInt(args, "count", 1)
	category, _ := getInt(args, "category", 0)
	difficulty, _ := getString(args, "difficulty")
	u := fmt.Sprintf("https://opentdb.com/api.php?amount=%d", n)
	if category > 0 {
		u += fmt.Sprintf("&category=%d", category)
	}
	if difficulty != "" {
		u += "&difficulty=" + difficulty
	}
	resp, apiErr := batch4API.Get(u)
	if apiErr != nil {
		return err("Open Trivia: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ResponseCode int `json:"response_code"`
		Results      []struct {
			Category   string   `json:"category"`
			Difficulty string   `json:"difficulty"`
			Question   string   `json:"question"`
			Correct    string   `json:"correct_answer"`
			Incorrect  []string `json:"incorrect_answers"`
		} `json:"results"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	for i, q := range data.Results {
		sb.WriteString(fmt.Sprintf("Q%d [%s/%s]: %s\n", i+1, q.Category, q.Difficulty, decodeHTML(q.Question)))
		sb.WriteString(fmt.Sprintf("Correct answer: %s\n\n", decodeHTML(q.Correct)))
	}
	return ok(sb.String())
}

func decodeHTML(s string) string {
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#039;", "'")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&eacute;", "é")
	s = strings.ReplaceAll(s, "&agrave;", "à")
	s = strings.ReplaceAll(s, "&ouml;", "ö")
	return s
}

// ── HELPERS ───────────────────────────────────────────────────────
