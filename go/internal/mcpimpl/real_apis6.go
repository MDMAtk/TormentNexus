package mcpimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════════════════
// Additional utility handlers (batch6 extras)
// ═══════════════════════════════════════════════════════════════════

func HandleTrendingBooks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	period, _ := getString(args, "period")
	if period == "" {
		period = "daily"
	}
	u := fmt.Sprintf("https://openlibrary.org/trending/%s.json?limit=10", period)
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Works []struct {
			Title  string `json:"title"`
			Author []struct {
				Name string `json:"name"`
			} `json:"authors"`
			Year int `json:"first_publish_year"`
		} `json:"works"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Trending (%s):\n", period))
	for _, w := range data.Works {
		a := ""
		if len(w.Author) > 0 {
			a = " by " + w.Author[0].Name
		}
		sb.WriteString(fmt.Sprintf("%s%s (%d)\n", w.Title, a, w.Year))
	}
	return ok(sb.String())
}

func HandleISSLocation(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch6API.Get("http://api.open-notify.org/iss-now.json")
	if apiErr != nil {
		return ok("ISS location unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ISSPosition struct {
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"iss_position"`
		Timestamp int64 `json:"timestamp"`
	}
	json.Unmarshal(body, &data)
	t := time.Unix(data.Timestamp, 0)
	return ok(fmt.Sprintf("ISS Location\nLat: %s Lon: %s\nTime: %s", data.ISSPosition.Latitude, data.ISSPosition.Longitude, t.Format("2006-01-02 15:04:05 UTC")))
}

func HandleYesNo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch6API.Get("https://yesno.wtf/api")
	if apiErr != nil {
		return ok("Yes")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Answer string `json:"answer"`
		Image  string `json:"image"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Answer: %s\n%s", strings.ToUpper(data.Answer), data.Image))
}

func HandleGenderPredict(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required")
	}
	u := fmt.Sprintf("https://api.genderize.io?name=%s", url.QueryEscape(name))
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return err("Genderize: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name        string  `json:"name"`
		Gender      string  `json:"gender"`
		Probability float64 `json:"probability"`
		Count       int     `json:"count"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Gender for %q: %s (%.0f%% based on %d)", data.Name, data.Gender, data.Probability*100, data.Count))
}

func HandleActivitySuggest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	u := "https://boredapi.com/api/activity"
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return ok("Try learning something new today!")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Activity string  `json:"activity"`
		Type     string  `json:"type"`
		Price    float64 `json:"price"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Activity: %s\nType: %s | Price: %.2f", data.Activity, data.Type, data.Price))
}

func HandleExchangeRate(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	from, _ := getString(args, "from")
	to, _ := getString(args, "to")
	if from == "" {
		return err("from currency required (e.g. USD, EUR, GBP)")
	}
	if to == "" {
		to = "USD"
	}
	u := fmt.Sprintf("https://api.frankfurter.app/latest?from=%s&to=%s", from, to)
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Rate %s->%s unavailable", from, to))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Date  string             `json:"date"`
		Rates map[string]float64 `json:"rates"`
	}
	json.Unmarshal(body, &data)
	if rate, exists := data.Rates[to]; exists {
		return ok(fmt.Sprintf("1 %s = %.4f %s (as of %s)", from, rate, to, data.Date))
	}
	return ok(fmt.Sprintf("Rate %s->%s not found", from, to))
}

func HandleMarsPhotos(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	rover, _ := getString(args, "rover")
	sol, _ := getInt(args, "sol", 1000)
	if rover == "" {
		rover = "curiosity"
	}
	u := fmt.Sprintf("https://api.nasa.gov/mars-photos/api/v1/rovers/%s/photos?sol=%d&api_key=DEMO_KEY", rover, sol)
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Mars %s sol %d unavailable", rover, sol))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Photos []struct {
			ID        int    `json:"id"`
			ImgSrc    string `json:"img_src"`
			EarthDate string `json:"earth_date"`
		} `json:"photos"`
	}
	json.Unmarshal(body, &data)
	if len(data.Photos) == 0 {
		return ok(fmt.Sprintf("No photos from %s on sol %d", rover, sol))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Mars Rover %s (sol %d): %d photos\n\n", rover, sol, len(data.Photos)))
	for i, p := range data.Photos {
		if i >= 3 {
			break
		}
		sb.WriteString(fmt.Sprintf("Photo %d (%s)\n  %s\n", p.ID, p.EarthDate, p.ImgSrc))
	}
	return ok(sb.String())
}

func HandleQuoteOfTheDay(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch6API.Get("https://favqs.com/api/qotd")
	if apiErr != nil {
		return ok("Quote of the day unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Quote struct {
			Body   string `json:"body"`
			Author string `json:"author"`
		} `json:"quote"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Quote of the Day:\n\n\"%s\"\n— %s", data.Quote.Body, data.Quote.Author))
}
