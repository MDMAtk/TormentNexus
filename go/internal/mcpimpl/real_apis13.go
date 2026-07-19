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

var batch13 = &http.Client{Timeout: 15 * time.Second}

func HandleNASASearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	n, _ := getInt(args, "limit", 5)
	media, _ := getString(args, "mediaType")
	if q == "" {
		return err("query is required (e.g. 'mars', 'apollo', 'nebula')")
	}
	u := fmt.Sprintf("https://images-api.nasa.gov/search?q=%s&page_size=%d", url.QueryEscape(q), n)
	if media != "" {
		u += "&media_type=" + media
	}
	resp, apiErr := batch13.Get(u)
	if apiErr != nil {
		return err("NASA Images: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Collection struct {
			Items []struct {
				Data []struct {
					Title        string `json:"title"`
					Description  string `json:"description"`
					DateCreated  string `json:"date_created"`
					Photographer string `json:"photographer"`
				} `json:"data"`
				Links []struct {
					Href string `json:"href"`
				} `json:"links"`
			} `json:"items"`
		} `json:"collection"`
	}
	json.Unmarshal(body, &data)
	items := data.Collection.Items
	if len(items) == 0 {
		return ok(fmt.Sprintf("No NASA images found for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("NASA Image Library — %d results for %q:\n\n", len(items), q))
	for _, item := range items {
		if len(item.Data) == 0 {
			continue
		}
		d := item.Data[0]
		sb.WriteString(fmt.Sprintf("%s\n", d.Title))
		desc := d.Description
		if len(desc) > 150 {
			desc = desc[:150] + "..."
		}
		sb.WriteString(fmt.Sprintf("  %s\n", desc))
		sb.WriteString(fmt.Sprintf("  Date: %s\n", d.DateCreated[:10]))
		if len(item.Links) > 0 {
			sb.WriteString(fmt.Sprintf("  %s\n", item.Links[0].Href))
		}
	}
	return ok(sb.String())
}

func HandleFloodWarnings(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude")
	if lat == 0 {
		lat = 48.85
	}
	lon := getFloat(args, "longitude")
	if lon == 0 {
		lon = 2.35
	}
	u := fmt.Sprintf("https://flood-api.open-meteo.com/v1/flood?latitude=%.2f&longitude=%.2f&daily=river_discharge&forecast_days=3", lat, lon)
	resp, apiErr := batch13.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Flood data at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Daily struct {
			Time      []string  `json:"time"`
			Discharge []float64 `json:"river_discharge"`
		} `json:"daily"`
	}
	json.Unmarshal(body, &data)
	if len(data.Daily.Time) == 0 {
		return ok("No flood data available")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Flood Warning — River Discharge at %.2f, %.2f:\n\n", lat, lon))
	for i, t := range data.Daily.Time {
		risk := "Normal"
		discharge := data.Daily.Discharge[i]
		if discharge > 1000 {
			risk = "⚠️ HIGH FLOOD RISK"
		} else if discharge > 500 {
			risk = "Moderate"
		}
		sb.WriteString(fmt.Sprintf("%s: %.0f m³/s (%s)\n", t, discharge, risk))
	}
	return ok(sb.String())
}

func HandleAirQualityForecast(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude")
	if lat == 0 {
		lat = 52.52
	}
	lon := getFloat(args, "longitude")
	if lon == 0 {
		lon = 13.41
	}
	u := fmt.Sprintf("https://air-quality-api.open-meteo.com/v1/air-quality?latitude=%.2f&longitude=%.2f&hourly=european_aqi,us_aqi,pm2_5,pm10&forecast_days=2", lat, lon)
	resp, apiErr := batch13.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Air quality forecast at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Hourly struct {
			Time  []string  `json:"time"`
			EUAQI []float64 `json:"european_aqi"`
			USAQI []float64 `json:"us_aqi"`
		} `json:"hourly"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Air Quality Forecast at %.2f, %.2f:\n\n", lat, lon))
	for i, t := range data.Hourly.Time {
		if i >= 12 {
			break
		}
		sb.WriteString(fmt.Sprintf("%s: EU AQI %.0f | US AQI %.0f\n", t[11:16], data.Hourly.EUAQI[i], data.Hourly.USAQI[i]))
	}
	return ok(sb.String())
}

func HandleNobelPrizes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	year, _ := getInt(args, "year", 0)
	category, _ := getString(args, "category")
	limit, _ := getInt(args, "limit", 5)
	u := fmt.Sprintf("https://api.nobelprize.org/2.1/nobelPrizes?limit=%d", limit)
	if year > 0 {
		u += fmt.Sprintf("&offset=%d", (time.Now().Year()-year)*6)
	}
	if category != "" {
		u += "&nobelPrizeCategory=" + category
	}
	resp, apiErr := batch13.Get(u)
	if apiErr != nil {
		return err("Nobel: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		NobelPrizes []struct {
			AwardYear string `json:"awardYear"`
			Category  string `json:"category"`
			Laureates []struct {
				Name string `json:"fullName"`
			} `json:"laureates"`
		} `json:"nobelPrizes"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString("Nobel Prizes:\n\n")
	for _, p := range data.NobelPrizes {
		sb.WriteString(fmt.Sprintf("%s — %s\n", p.AwardYear, p.Category))
		for _, l := range p.Laureates {
			sb.WriteString(fmt.Sprintf("  %s\n", l.Name))
		}
	}
	return ok(sb.String())
}

func HandleTVShowsByPage(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	page, _ := getInt(args, "page", 1)
	u := fmt.Sprintf("https://api.tvmaze.com/shows?page=%d", page)
	resp, apiErr := batch13.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("TV shows page %d unavailable", page))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var shows []struct {
		Name   string   `json:"name"`
		Type   string   `json:"type"`
		Genres []string `json:"genres"`
		Rating struct {
			Average float64 `json:"average"`
		} `json:"rating"`
		Status string `json:"status"`
	}
	json.Unmarshal(body, &shows)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("TV Shows (page %d, %d shows):\n\n", page, len(shows)))
	for _, s := range shows[:20] {
		sb.WriteString(fmt.Sprintf("%s [%s] %.1f — %s\n", s.Name, strings.Join(s.Genres, ", "), s.Rating.Average, s.Status))
	}
	return ok(sb.String())
}

func HandleBookShelves(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	workID, _ := getString(args, "workId")
	if workID == "" {
		workID = "OL45804W"
	}
	u := fmt.Sprintf("https://openlibrary.org/works/%s/bookshelves.json", workID)
	resp, apiErr := batch13.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Bookshelves for %s unavailable", workID))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Counts []struct {
			Name  string `json:"name"`
			Count int    `json:"count"`
		} `json:"counts"`
	}
	json.Unmarshal(body, &data)
	if len(data.Counts) == 0 {
		return ok("No bookshelf data")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Bookshelves for work %s:\n\n", workID))
	for _, c := range data.Counts {
		sb.WriteString(fmt.Sprintf("%s: %d\n", c.Name, c.Count))
	}
	return ok(sb.String())
}
