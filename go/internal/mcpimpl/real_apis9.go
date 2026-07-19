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

var batch9 = &http.Client{Timeout: 15 * time.Second}

func HandleMealRandom(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch9.Get("https://www.themealdb.com/api/json/v1/1/random.php")
	if apiErr != nil {
		return ok("Recipe unavailable. Try making scrambled eggs!")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Meals []map[string]interface{} `json:"meals"`
	}
	json.Unmarshal(body, &data)
	if len(data.Meals) == 0 {
		return ok("No random meal found")
	}
	m := data.Meals[0]
	name, _ := m["strMeal"].(string)
	cat, _ := m["strCategory"].(string)
	area, _ := m["strArea"].(string)
	instructions, _ := m["strInstructions"].(string)
	thumb, _ := m["strMealThumb"].(string)
	if len(instructions) > 500 {
		instructions = instructions[:500] + "..."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Random Recipe: %s\n", name))
	sb.WriteString(fmt.Sprintf("Category: %s | Cuisine: %s\n\n", cat, area))
	sb.WriteString(fmt.Sprintf("Instructions:\n%s\n\n", instructions))
	sb.WriteString(fmt.Sprintf("📷 %s\n", thumb))
	return ok(sb.String())
}

func HandleMealCategories(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch9.Get("https://www.themealdb.com/api/json/v1/1/categories.php")
	if apiErr != nil {
		return err("MealDB: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Categories []struct {
			Name        string `json:"strCategory"`
			Description string `json:"strCategoryDescription"`
		} `json:"categories"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Meal Categories (%d):\n\n", len(data.Categories)))
	for _, c := range data.Categories {
		d := c.Description
		if len(d) > 120 {
			d = d[:120] + "..."
		}
		sb.WriteString(fmt.Sprintf("%s\n  %s\n\n", c.Name, d))
	}
	return ok(sb.String())
}

func HandleMealAreas(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch9.Get("https://www.themealdb.com/api/json/v1/1/list.php?a=list")
	if apiErr != nil {
		return err("MealDB: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Meals []struct {
			Area string `json:"strArea"`
		} `json:"meals"`
	}
	json.Unmarshal(body, &data)
	var areas []string
	for _, m := range data.Meals {
		areas = append(areas, m.Area)
	}
	return ok(fmt.Sprintf("Cuisine Areas (%d):\n%s", len(areas), strings.Join(areas, ", ")))
}

func HandleCocktailSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	ingredient, _ := getString(args, "ingredient")
	random := getBool(args, "random")
	var u string
	if random {
		u = "https://www.thecocktaildb.com/api/json/v1/1/random.php"
	} else if name != "" {
		u = fmt.Sprintf("https://www.thecocktaildb.com/api/json/v1/1/search.php?s=%s", url.QueryEscape(name))
	} else if ingredient != "" {
		u = fmt.Sprintf("https://www.thecocktaildb.com/api/json/v1/1/filter.php?i=%s", url.QueryEscape(ingredient))
	} else {
		u = "https://www.thecocktaildb.com/api/json/v1/1/random.php"
	}
	resp, apiErr := batch9.Get(u)
	if apiErr != nil {
		return ok("CocktailDB unavailable. Try a classic: Mojito!")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Drinks []map[string]interface{} `json:"drinks"`
	}
	json.Unmarshal(body, &data)
	if len(data.Drinks) == 0 {
		return ok("No cocktails found")
	}
	var sb strings.Builder
	for _, d := range data.Drinks {
		name, _ := d["strDrink"].(string)
		cat, _ := d["strCategory"].(string)
		glass, _ := d["strGlass"].(string)
		alcoholic, _ := d["strAlcoholic"].(string)
		instructions, _ := d["strInstructions"].(string)
		thumb, _ := d["strDrinkThumb"].(string)
		sb.WriteString(fmt.Sprintf("%s\n", name))
		sb.WriteString(fmt.Sprintf("Category: %s | Glass: %s | %s\n", cat, glass, alcoholic))
		if len(instructions) > 300 {
			instructions = instructions[:300] + "..."
		}
		sb.WriteString(fmt.Sprintf("Instructions: %s\n", instructions))
		// Ingredients
		for i := 1; i <= 15; i++ {
			ing := d[fmt.Sprintf("strIngredient%d", i)]
			meas := d[fmt.Sprintf("strMeasure%d", i)]
			if ing == nil || ing.(string) == "" {
				break
			}
			sb.WriteString(fmt.Sprintf("  • %s %s\n", meas, ing))
		}
		sb.WriteString(fmt.Sprintf("%s\n", thumb))
	}
	return ok(sb.String())
}

func HandleTopAnime(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	n, _ := getInt(args, "limit", 5)
	u := fmt.Sprintf("https://api.jikan.moe/v4/top/anime?limit=%d", n)
	resp, apiErr := batch9.Get(u)
	if apiErr != nil {
		return ok("Top anime: service unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Data []struct {
			Title    string  `json:"title"`
			Score    float64 `json:"score"`
			Type     string  `json:"type"`
			Episodes int     `json:"episodes"`
			Status   string  `json:"status"`
			URL      string  `json:"url"`
		} `json:"data"`
	}
	json.Unmarshal(body, &data)
	if len(data.Data) == 0 {
		return ok("No anime found")
	}
	var sb strings.Builder
	sb.WriteString("Top Anime:\n\n")
	for _, a := range data.Data {
		sb.WriteString(fmt.Sprintf("%s\n", a.Title))
		sb.WriteString(fmt.Sprintf("Score: %.2f | Type: %s | Episodes: %d | Status: %s\n", a.Score, a.Type, a.Episodes, a.Status))
		sb.WriteString(fmt.Sprintf("%s\n\n", a.URL))
	}
	return ok(sb.String())
}

func HandleLanguageList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch9.Get("https://api.languagetool.org/v2/languages")
	if apiErr != nil {
		return ok("Languages: service unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var langs []struct {
		Name     string `json:"name"`
		Code     string `json:"code"`
		LongCode string `json:"longCode"`
	}
	json.Unmarshal(body, &langs)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("LanguageTool — %d supported languages:\n\n", len(langs)))
	for _, l := range langs[:20] {
		sb.WriteString(fmt.Sprintf("%-20s %s\n", l.Name, l.Code))
	}
	return ok(sb.String())
}

func HandleIP2Location(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	ip, _ := getString(args, "ip")
	if ip == "" {
		ip = "8.8.8.8"
	}
	u := fmt.Sprintf("https://api.ip2location.io/?ip=%s&format=json", ip)
	resp, apiErr := batch9.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("IP2Location for %s unavailable", ip))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		IP          string  `json:"ip"`
		CountryCode string  `json:"country_code"`
		CountryName string  `json:"country_name"`
		RegionName  string  `json:"region_name"`
		CityName    string  `json:"city_name"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		ZipCode     string  `json:"zip_code"`
		TimeZone    string  `json:"time_zone"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("IP2Location — %s\nCountry: %s (%s)\nRegion: %s | City: %s\nCoords: %.4f, %.4f\nZip: %s | TZ: %s",
		data.IP, data.CountryName, data.CountryCode, data.RegionName, data.CityName,
		data.Latitude, data.Longitude, data.ZipCode, data.TimeZone))
}

func HandleMealByArea(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	area, _ := getString(args, "area")
	if area == "" {
		return err("area is required (e.g. Canadian, Italian, Japanese)")
	}
	u := fmt.Sprintf("https://www.themealdb.com/api/json/v1/1/filter.php?a=%s", url.QueryEscape(area))
	resp, apiErr := batch9.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("No recipes found for %q", area))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Meals []struct {
			Name string `json:"strMeal"`
		} `json:"meals"`
	}
	json.Unmarshal(body, &data)
	if len(data.Meals) == 0 {
		return ok(fmt.Sprintf("No %s recipes found", area))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Recipes from %s (%d):\n\n", area, len(data.Meals)))
	for _, m := range data.Meals {
		sb.WriteString(fmt.Sprintf("%s\n", m.Name))
	}
	return ok(sb.String())
}

func HandleCocktailByIngredient(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	ingredient, _ := getString(args, "ingredient")
	if ingredient == "" {
		return err("ingredient is required (e.g. gin, vodka, rum)")
	}
	u := fmt.Sprintf("https://www.thecocktaildb.com/api/json/v1/1/filter.php?i=%s", url.QueryEscape(ingredient))
	resp, apiErr := batch9.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("No cocktails with %q", ingredient))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Drinks []struct {
			Name string `json:"strDrink"`
		} `json:"drinks"`
	}
	json.Unmarshal(body, &data)
	if len(data.Drinks) == 0 {
		return ok(fmt.Sprintf("No cocktails with %q", ingredient))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Cocktails with %s (%d):\n\n", ingredient, len(data.Drinks)))
	for _, d := range data.Drinks {
		sb.WriteString(fmt.Sprintf("🍸 %s\n", d.Name))
	}
	return ok(sb.String())
}

func HandleOpenLibrarySubjectsList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	u := "https://openlibrary.org/subjects.json?limit=50"
	resp, apiErr := batch9.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Subjects []struct {
			Name  string `json:"name"`
			Count int    `json:"work_count"`
		} `json:"subjects"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Open Library Subjects (%d):\n\n", len(data.Subjects)))
	for _, s := range data.Subjects {
		sb.WriteString(fmt.Sprintf("%s (%d works)\n", s.Name, s.Count))
	}
	return ok(sb.String())
}
