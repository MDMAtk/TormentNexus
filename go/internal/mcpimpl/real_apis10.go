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

var batch10 = &http.Client{Timeout: 15 * time.Second}

func HandleEarthquakes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	minMag := getFloat(args, "minMagnitude")
	if minMag == 0 {
		minMag = 2.5
	}
	period, _ := getString(args, "period")
	if period == "" {
		period = "day"
	}
	feed := map[string]string{"hour": "1hour", "day": "day", "week": "week"}
	suffix := feed[period]
	if suffix == "" {
		suffix = "day"
	}
	u := fmt.Sprintf("https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/%s_%s.geojson", minMagStr(minMag), suffix)
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return err("USGS: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Metadata struct {
			Count int `json:"count"`
		} `json:"metadata"`
		Features []struct {
			Properties struct {
				Mag     float64 `json:"mag"`
				Place   string  `json:"place"`
				Time    int64   `json:"time"`
				URL     string  `json:"url"`
				Detail  string  `json:"detail"`
				Tsunami int     `json:"tsunami"`
			} `json:"properties"`
			Geometry struct {
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	json.Unmarshal(body, &data)
	if len(data.Features) == 0 {
		return ok("No earthquakes found")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Earthquakes (last %s) — %d events:\n\n", period, data.Metadata.Count))
	for _, f := range data.Features {
		t := time.UnixMilli(f.Properties.Time)
		sb.WriteString(fmt.Sprintf("M%.1f — %s\n", f.Properties.Mag, f.Properties.Place))
		sb.WriteString(fmt.Sprintf("  %s | Depth: %.0f km\n", t.Format("2006-01-02 15:04 UTC"), f.Geometry.Coordinates[2]))
		if f.Properties.Tsunami > 0 {
			sb.WriteString("  ⚠️ TSUNAMI WARNING\n")
		}
	}
	return ok(sb.String())
}

func minMagStr(mag float64) string {
	if mag >= 4.5 {
		return "4.5"
	} else if mag >= 2.5 {
		return "2.5"
	} else if mag >= 1.0 {
		return "1.0"
	}
	return "all"
}

func HandleSunTimes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude")
	if lat == 0 {
		lat = 38.9072
	}
	lon := getFloat(args, "longitude")
	if lon == 0 {
		lon = -77.0369
	}
	u := fmt.Sprintf("https://api.sunrisesunset.io/json?lat=%.4f&lng=%.4f", lat, lon)
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Sun times at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Results struct {
			Sunrise   string `json:"sunrise"`
			Sunset    string `json:"sunset"`
			Dawn      string `json:"dawn"`
			Dusk      string `json:"dusk"`
			DayLength string `json:"day_length"`
			SolarNoon string `json:"solar_noon"`
			TimeZone  string `json:"timezone"`
		} `json:"results"`
	}
	json.Unmarshal(body, &data)
	r := data.Results
	return ok(fmt.Sprintf("Sun Times at %.4f, %.4f [%s]:\n\nSunrise: %s\nSunset: %s\nDawn: %s\nDusk: %s\nDay Length: %s\nSolar Noon: %s",
		lat, lon, r.TimeZone, r.Sunrise, r.Sunset, r.Dawn, r.Dusk, r.DayLength, r.SolarNoon))
}

func HandleAirQuality(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude")
	if lat == 0 {
		lat = 52.52
	}
	lon := getFloat(args, "longitude")
	if lon == 0 {
		lon = 13.41
	}
	u := fmt.Sprintf("https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%.2f&longitude=%.2f&current=european_aqi,us_aqi,pm2_5,pm10,nitrogen_dioxide,ozone,sulphur_dioxide",
		lat, lon)
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Air quality at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Current struct {
			EuropeanAQI     float64 `json:"european_aqi"`
			USAQI           float64 `json:"us_aqi"`
			PM25            float64 `json:"pm2_5"`
			PM10            float64 `json:"pm10"`
			NitrogenDioxide float64 `json:"nitrogen_dioxide"`
			Ozone           float64 `json:"ozone"`
			SulphurDioxide  float64 `json:"sulphur_dioxide"`
		} `json:"current"`
	}
	json.Unmarshal(body, &data)
	c := data.Current
	aqiDesc := func(v float64) string {
		if v <= 20 {
			return "Good"
		}
		if v <= 40 {
			return "Fair"
		}
		if v <= 60 {
			return "Moderate"
		}
		if v <= 80 {
			return "Poor"
		}
		return "Very Poor"
	}
	return ok(fmt.Sprintf("Air Quality at %.2f, %.2f:\n\nEU AQI: %.0f (%s)\nUS AQI: %.0f\nPM2.5: %.1f µg/m³\nPM10: %.1f µg/m³\nNO₂: %.1f µg/m³\nO₃: %.1f µg/m³\nSO₂: %.1f µg/m³",
		lat, lon, c.EuropeanAQI, aqiDesc(c.EuropeanAQI), c.USAQI, c.PM25, c.PM10, c.NitrogenDioxide, c.Ozone, c.SulphurDioxide))
}

func HandleMusicSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	artist, _ := getString(args, "artist")
	id, _ := getString(args, "id")
	if artist == "" && id == "" {
		return err("artist name or MusicBrainz ID is required")
	}
	var u string
	if id != "" {
		u = fmt.Sprintf("https://musicbrainz.org/ws/2/artist/%s?fmt=json&inc=tags+genres", id)
	} else {
		u = fmt.Sprintf("https://musicbrainz.org/ws/2/artist?query=%s&fmt=json&limit=5", url.QueryEscape(artist))
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("User-Agent", "TormentNexus/1.0 ( tornexus@dev )")
	req.Header.Set("Accept", "application/json")
	resp, apiErr := batch10.Do(req)
	if apiErr != nil {
		return ok(fmt.Sprintf("MusicBrainz: service unavailable"))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if id != "" {
		var data struct {
			Name    string `json:"name"`
			Type    string `json:"type"`
			Country string `json:"country"`
			Tags    []struct {
				Name string `json:"name"`
			} `json:"tags"`
		}
		json.Unmarshal(body, &data)
		var tags []string
		for _, t := range data.Tags {
			tags = append(tags, t.Name)
		}
		return ok(fmt.Sprintf("Artist: %s\nType: %s | Country: %s\nTags: %s",
			data.Name, data.Type, data.Country, strings.Join(tags, ", ")))
	}
	var data struct {
		Artists []struct {
			Name    string `json:"name"`
			Type    string `json:"type"`
			Country string `json:"country"`
			ID      string `json:"id"`
		} `json:"artists"`
	}
	json.Unmarshal(body, &data)
	if len(data.Artists) == 0 {
		return ok(fmt.Sprintf("No artists found for %q", artist))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("MusicBrainz — artists matching %q:\n\n", artist))
	for _, a := range data.Artists {
		sb.WriteString(fmt.Sprintf("%s (%s, %s)\n  ID: %s\n", a.Name, a.Type, a.Country, a.ID))
	}
	return ok(sb.String())
}

func HandleTVShowByID(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 1)
	u := fmt.Sprintf("https://api.tvmaze.com/shows/%d", id)
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
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
	}
	json.Unmarshal(body, &data)
	summary := stripHTMLTags(data.Summary)
	if len(summary) > 300 {
		summary = summary[:300] + "..."
	}
	return ok(fmt.Sprintf("TV Show: %s\nType: %s | Language: %s\nGenres: %s\nStatus: %s | Rating: %.1f\n%s\n%s",
		data.Name, data.Type, data.Language, strings.Join(data.Genres, ", "), data.Status, data.Rating.Average, summary, data.URL))
}

func HandleNASANotifications(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	u := "https://api.nasa.gov/DONKI/notifications?api_key=DEMO_KEY"
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return ok("NASA DONKI: service unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var notifications []struct {
		MessageType string `json:"messageType"`
		Message     string `json:"message"`
		MessageURL  string `json:"messageURL"`
	}
	json.Unmarshal(body, &notifications)
	if len(notifications) == 0 {
		return ok("No space weather notifications")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("NASA Space Weather — %d notifications:\n\n", len(notifications)))
	for _, n := range notifications[:10] {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", n.MessageType, n.MessageURL))
	}
	return ok(sb.String())
}

func HandleOpenLibraryWork(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	olid, _ := getString(args, "olid")
	if olid == "" {
		return err("olid is required (e.g. 'OL45804W')")
	}
	u := fmt.Sprintf("https://openlibrary.org/works/%s.json", olid)
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Title   string `json:"title"`
		Key     string `json:"key"`
		Authors []struct {
			Author struct {
				Key string `json:"key"`
			} `json:"author"`
		} `json:"authors"`
		Description interface{} `json:"description"`
		Subjects    []string    `json:"subjects"`
	}
	json.Unmarshal(body, &data)
	desc := ""
	if d, ok := data.Description.(string); ok {
		desc = d
	} else if m, ok := data.Description.(map[string]interface{}); ok {
		desc, _ = m["value"].(string)
	}
	if len(desc) > 400 {
		desc = desc[:400] + "..."
	}
	var authors []string
	for _, a := range data.Authors {
		authors = append(authors, a.Author.Key)
	}
	return ok(fmt.Sprintf("Work: %s (%s)\nAuthors: %s\n\n%s\n\nSubjects: %s",
		data.Title, data.Key, strings.Join(authors, ", "), desc, strings.Join(data.Subjects[:min(len(data.Subjects), 8)], ", ")))
}

func HandleBeerSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	abvMin := getFloat(args, "abvMin")
	u := "https://api.punkapi.com/v2/beers?page=1&per_page=5"
	if name != "" {
		u += "&beer_name=" + url.QueryEscape(name)
	}
	if abvMin > 0 {
		u += fmt.Sprintf("&abv_gt=%.0f", abvMin)
	}
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return ok("Beer data unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var beers []struct {
		Name        string   `json:"name"`
		Tagline     string   `json:"tagline"`
		ABV         float64  `json:"abv"`
		IBU         float64  `json:"ibu"`
		EBC         float64  `json:"ebc"`
		Description string   `json:"description"`
		FoodPairing []string `json:"food_pairing"`
	}
	json.Unmarshal(body, &beers)
	if len(beers) == 0 {
		return ok("No beers found")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Beers (%d found):\n\n", len(beers)))
	for _, b := range beers {
		sb.WriteString(fmt.Sprintf("%s — %s\n", b.Name, b.Tagline))
		sb.WriteString(fmt.Sprintf("ABV: %.1f%% | IBU: %.0f | EBC: %.0f\n", b.ABV, b.IBU, b.EBC))
		d := b.Description
		if len(d) > 200 {
			d = d[:200] + "..."
		}
		sb.WriteString(fmt.Sprintf("%s\n", d))
		if len(b.FoodPairing) > 0 {
			sb.WriteString(fmt.Sprintf("Pairs with: %s\n", strings.Join(b.FoodPairing[:min(3, len(b.FoodPairing))], ", ")))
		}
	}
	return ok(sb.String())
}

func HandleCoinCapPrices(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	limit, _ := getInt(args, "limit", 10)
	u := fmt.Sprintf("https://api.coincap.io/v2/assets?limit=%d", limit)
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return ok("CoinCap: service unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Data []struct {
			Symbol    string `json:"symbol"`
			Name      string `json:"name"`
			PriceUSD  string `json:"priceUsd"`
			Change24  string `json:"changePercent24Hr"`
			MarketCap string `json:"marketCapUsd"`
		} `json:"data"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString("Crypto Prices (CoinCap):\n\n")
	for _, a := range data.Data {
		sb.WriteString(fmt.Sprintf("%s (%s): $%.4f (24h: %s%%)\n", a.Name, a.Symbol, parseFloatSafe(a.PriceUSD), a.Change24))
	}
	return ok(sb.String())
}

func HandleCoinCapSingle(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	symbol, _ := getString(args, "symbol")
	if symbol == "" {
		return err("symbol is required (e.g. bitcoin, ethereum)")
	}
	u := fmt.Sprintf("https://api.coincap.io/v2/assets/%s", strings.ToLower(symbol))
	resp, apiErr := batch10.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("CoinCap: %q not found", symbol))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Data struct {
			Symbol    string `json:"symbol"`
			Name      string `json:"name"`
			PriceUSD  string `json:"priceUsd"`
			Change24  string `json:"changePercent24Hr"`
			Volume    string `json:"volumeUsd24Hr"`
			Supply    string `json:"supply"`
			MarketCap string `json:"marketCapUsd"`
		} `json:"data"`
	}
	json.Unmarshal(body, &data)
	d := data.Data
	return ok(fmt.Sprintf("%s (%s)\nPrice: $%.4f\n24h Change: %s%%\nVolume: $%s\nSupply: %s\nMarket Cap: $%s",
		d.Name, d.Symbol, parseFloatSafe(d.PriceUSD), d.Change24, d.Volume, d.Supply, d.MarketCap))
}

func parseFloatSafe(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
