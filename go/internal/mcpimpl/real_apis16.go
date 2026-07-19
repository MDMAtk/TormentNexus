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

var batch16 = &http.Client{Timeout: 15 * time.Second}

func HandleNASADailyEarth(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	u := "https://epic.gsfc.nasa.gov/api/natural"
	resp, apiErr := batch16.Get(u)
	if apiErr != nil {
		return ok("EPIC unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var items []struct {
		Date     string `json:"date"`
		Caption  string `json:"caption"`
		Image    string `json:"image"`
		Centroid struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"centroid_coordinates"`
	}
	json.Unmarshal(body, &items)
	if len(items) == 0 {
		return ok("No EPIC images available")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("NASA EPIC — Daily Earth Images (%d):\n\n", len(items)))
	for _, item := range items[:5] {
		date := item.Date[:10]
		sb.WriteString(fmt.Sprintf("%s — %s\n", date, item.Caption))
		sb.WriteString(fmt.Sprintf("  Center: %.2f, %.2f\n", item.Centroid.Lat, item.Centroid.Lon))
		sb.WriteString(fmt.Sprintf("  https://epic.gsfc.nasa.gov/archive/natural/%s/png/%s.png\n", strings.ReplaceAll(date, "-", "/"), item.Image))
	}
	return ok(sb.String())
}

func HandleMarsWeather(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	u := "https://api.nasa.gov/insight_weather/?api_key=DEMO_KEY&feedtype=json&ver=1.0"
	resp, apiErr := batch16.Get(u)
	if apiErr != nil {
		return ok("Mars weather unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw map[string]interface{}
	json.Unmarshal(body, &raw)
	var sb strings.Builder
	sb.WriteString("Mars Weather — InSight Lander (Elysium Planitia):\n\n")
	for key, val := range raw {
		sol, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		if _, err := strconv.Atoi(key); err != nil {
			continue
		}
		season, _ := sol["Season"].(string)
		temp, _ := sol["AT"].(map[string]interface{})
		press, _ := sol["PRE"].(map[string]interface{})
		wind, _ := sol["HWS"].(map[string]interface{})
		avgTemp, _ := temp["av"].(float64)
		avgPress, _ := press["av"].(float64)
		avgWind, _ := wind["av"].(float64)
		sb.WriteString(fmt.Sprintf("Sol %s [%s]: %.1f°C, %.1f Pa, %.1f m/s wind\n", key, season, avgTemp, avgPress, avgWind))
	}
	return ok(sb.String())
}

func HandlePokemonMove(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/move/%s", strings.ToLower(q))
	resp, apiErr := batch16.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Move %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name     string `json:"name"`
		Power    int    `json:"power"`
		PP       int    `json:"pp"`
		Accuracy int    `json:"accuracy"`
		Priority int    `json:"priority"`
		Type     struct {
			Name string `json:"name"`
		} `json:"type"`
		Damage struct {
			Name string `json:"name"`
		} `json:"damage_class"`
		Effect string `json:"effect"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Move: %s (%s)\nType: %s | Power: %d | PP: %d | Accuracy: %d%% | Priority: %d",
		data.Name, data.Damage.Name, data.Type.Name, data.Power, data.PP, data.Accuracy, data.Priority))
}

func HandlePokemonItem(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	id, _ := getInt(args, "id", 0)
	if name == "" && id == 0 {
		return err("name or id is required")
	}
	q := name
	if q == "" {
		q = fmt.Sprintf("%d", id)
	}
	u := fmt.Sprintf("https://pokeapi.co/api/v2/item/%s", strings.ToLower(q))
	resp, apiErr := batch16.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Item %q not found", q))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name     string `json:"name"`
		Cost     int    `json:"cost"`
		Category struct {
			Name string `json:"name"`
		} `json:"category"`
		Effect []struct {
			Text string `json:"text"`
		} `json:"effect_entries"`
		Flavor []struct {
			Text string `json:"text"`
		} `json:"flavor_text_entries"`
	}
	json.Unmarshal(body, &data)
	effect := ""
	if len(data.Effect) > 0 {
		effect = data.Effect[0].Text
	}
	if effect == "" && len(data.Flavor) > 0 {
		effect = data.Flavor[0].Text
	}
	if len(effect) > 300 {
		effect = effect[:300] + "..."
	}
	return ok(fmt.Sprintf("Item: %s\nCost: %d¥ | Category: %s\n\n%s", data.Name, data.Cost, data.Category.Name, effect))
}

func HandleSolarRadiation(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude")
	if lat == 0 {
		lat = 40.0
	}
	lon := getFloat(args, "longitude")
	if lon == 0 {
		lon = -105.0
	}
	u := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&daily=shortwave_radiation_sum,direct_radiation_sum&timezone=auto&forecast_days=3", lat, lon)
	resp, apiErr := batch16.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Solar data at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Daily struct {
			Time  []string  `json:"time"`
			Short []float64 `json:"shortwave_radiation_sum"`
		} `json:"daily"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Solar Radiation Forecast at %.2f, %.2f:\n\n", lat, lon))
	for i, t := range data.Daily.Time {
		sb.WriteString(fmt.Sprintf("%s: %.0f MJ/m²\n", t, data.Daily.Short[i]))
	}
	return ok(sb.String())
}

func HandleEvapotranspiration(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lat := getFloat(args, "latitude")
	if lat == 0 {
		lat = 40.0
	}
	lon := getFloat(args, "longitude")
	if lon == 0 {
		lon = -105.0
	}
	u := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&daily=et0_fao_evapotranspiration&timezone=auto&forecast_days=3", lat, lon)
	resp, apiErr := batch16.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("ET₀ data at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Daily struct {
			Time []string  `json:"time"`
			ET0  []float64 `json:"et0_fao_evapotranspiration"`
		} `json:"daily"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Reference Evapotranspiration (ET₀) at %.2f, %.2f:\n\n", lat, lon))
	for i, t := range data.Daily.Time {
		sb.WriteString(fmt.Sprintf("%s: %.1f mm\n", t, data.Daily.ET0[i]))
	}
	return ok(sb.String())
}
