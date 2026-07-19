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

var batch12 = &http.Client{Timeout: 15 * time.Second}

func HandleUVIndex(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	v := getFloat(args, "latitude"); if v == 0 { v = 20.0 }; lat := v
	v2 := getFloat(args, "longitude"); if v2 == 0 { v2 = -155.0 }; lon := v2
	u := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&daily=uv_index_max,uv_index_clear_sky_max&timezone=auto&forecast_days=7", lat, lon)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("UV index at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Daily struct {
			Time []string  `json:"time"`
			UV   []float64 `json:"uv_index_max"`
		} `json:"daily"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("UV Index Forecast at %.2f, %.2f:\n\n", lat, lon))
	for i, t := range data.Daily.Time {
		uv := data.Daily.UV[i]
		risk := "Low"
		if uv >= 3 { risk = "Moderate" }
		if uv >= 6 { risk = "High" }
		if uv >= 8 { risk = "Very High" }
		if uv >= 11 { risk = "Extreme" }
		sb.WriteString(fmt.Sprintf("%s: UV %.1f (%s)\n", t, uv, risk))
	}
	return ok(sb.String())
}

func HandleClimateData(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	v := getFloat(args, "latitude"); if v == 0 { v = 40.71 }; lat := v
	v2 := getFloat(args, "longitude"); if v2 == 0 { v2 = -74.01 }; lon := v2
	start, _ := getInt(args, "startYear", 2000)
	end, _ := getInt(args, "endYear", 2020)
	u := fmt.Sprintf("https://climate-api.open-meteo.com/v1/climate?latitude=%.2f&longitude=%.2f&start_year=%d&end_year=%d&daily=temperature_2m_max,temperature_2m_min,precipitation_sum", lat, lon, start, end)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Climate data at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Daily struct {
			Time  []string  `json:"time"`
			TMax  []float64 `json:"temperature_2m_max"`
			TMin  []float64 `json:"temperature_2m_min"`
			Precip []float64 `json:"precipitation_sum"`
		} `json:"daily"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Climate Data at %.2f, %.2f (%d-%d):\n\n", lat, lon, start, end))
	for i, t := range data.Daily.Time {
		sb.WriteString(fmt.Sprintf("%s: %.1f - %.1f°C, %.0fmm precip\n", t, data.Daily.TMin[i], data.Daily.TMax[i], data.Daily.Precip[i]))
	}
	return ok(sb.String())
}

func HandleHistoricalWeather(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	v := getFloat(args, "latitude"); if v == 0 { v = 52.52 }; lat := v
	v2 := getFloat(args, "longitude"); if v2 == 0 { v2 = 13.41 }; lon := v2
	start, _ := getString(args, "startDate")
	end, _ := getString(args, "endDate")
	if start == "" {
		start = "2024-01-01"
	}
	if end == "" {
		end = "2024-01-07"
	}
	u := fmt.Sprintf("https://archive-api.open-meteo.com/v1/archive?latitude=%.2f&longitude=%.2f&start_date=%s&end_date=%s&daily=temperature_2m_max,temperature_2m_min,precipitation_sum,sunshine_duration", lat, lon, start, end)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Historical weather at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Daily struct {
			Time   []string  `json:"time"`
			TMax   []float64 `json:"temperature_2m_max"`
			TMin   []float64 `json:"temperature_2m_min"`
			Precip []float64 `json:"precipitation_sum"`
			Sun    []float64 `json:"sunshine_duration"`
		} `json:"daily"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Historical Weather at %.2f, %.2f (%s to %s):\n\n", lat, lon, start, end))
	for i, t := range data.Daily.Time {
		sunMin := data.Daily.Sun[i] / 60
		sb.WriteString(fmt.Sprintf("%s: %.1f-%.1f°C, %.0fmm, %.0f min sun\n", t, data.Daily.TMin[i], data.Daily.TMax[i], data.Daily.Precip[i], sunMin))
	}
	return ok(sb.String())
}

func HandleForecastExtended(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	v := getFloat(args, "latitude"); if v == 0 { v = 35.68 }; lat := v
	v2 := getFloat(args, "longitude"); if v2 == 0 { v2 = 139.69 }; lon := v2
	days, _ := getInt(args, "days", 7)
	u := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.2f&longitude=%.2f&daily=temperature_2m_max,temperature_2m_min,precipitation_sum,precipitation_probability_max,wind_speed_10m_max&forecast_days=%d", lat, lon, days)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Forecast at %.2f, %.2f unavailable", lat, lon))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Daily struct {
			Time       []string  `json:"time"`
			TMax       []float64 `json:"temperature_2m_max"`
			TMin       []float64 `json:"temperature_2m_min"`
			Precip     []float64 `json:"precipitation_sum"`
			PrecipProb []float64 `json:"precipitation_probability_max"`
			Wind       []float64 `json:"wind_speed_10m_max"`
		} `json:"daily"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Weather Forecast at %.2f, %.2f (%d days):\n\n", lat, lon, days))
	for i, t := range data.Daily.Time {
		sb.WriteString(fmt.Sprintf("%s: %.1f-%.1f°C, %.0fmm (%.0f%%), wind %.0f km/h\n", t, data.Daily.TMin[i], data.Daily.TMax[i], data.Daily.Precip[i], data.Daily.PrecipProb[i], data.Daily.Wind[i]))
	}
	return ok(sb.String())
}

func HandleUserTodos(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	userID, _ := getInt(args, "userId", 1)
	u := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d/todos", userID)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return err("JSONPlaceholder: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var todos []struct {
		Title    string `json:"title"`
		Completed bool  `json:"completed"`
	}
	json.Unmarshal(body, &todos)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Todos for user %d (%d items):\n\n", userID, len(todos)))
	for _, t := range todos {
		icon := "☐"
		if t.Completed {
			icon = "☑"
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", icon, t.Title))
	}
	return ok(sb.String())
}

func HandleUserAlbums(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	userID, _ := getInt(args, "userId", 1)
	u := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d/albums", userID)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return err("JSONPlaceholder: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var albums []struct {
		Title string `json:"title"`
	}
	json.Unmarshal(body, &albums)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Albums for user %d:\n\n", userID))
	for _, a := range albums {
		sb.WriteString(fmt.Sprintf("%s\n", a.Title))
	}
	return ok(sb.String())
}

func HandlePostComments(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	postID, _ := getInt(args, "postId", 1)
	u := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d/comments", postID)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return err("JSONPlaceholder: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var comments []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Body  string `json:"body"`
	}
	json.Unmarshal(body, &comments)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Comments for post %d (%d total):\n\n", postID, len(comments)))
	for _, c := range comments {
		sb.WriteString(fmt.Sprintf("%s (%s)\n  %s\n\n", c.Name, c.Email, c.Body))
	}
	return ok(sb.String())
}

func HandleJSONUser(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 1)
	u := fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", id)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return err("JSONPlaceholder: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Website  string `json:"website"`
		Address  struct {
			Street  string `json:"street"`
			Suite   string `json:"suite"`
			City    string `json:"city"`
			Zipcode string `json:"zipcode"`
		} `json:"address"`
		Company struct {
			Name string `json:"name"`
		} `json:"company"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("User #%d: %s (%s)\nEmail: %s\nPhone: %s\nWebsite: %s\nAddress: %s %s, %s %s\nCompany: %s",
		id, data.Name, data.Username, data.Email, data.Phone, data.Website,
		data.Address.Street, data.Address.Suite, data.Address.City, data.Address.Zipcode, data.Company.Name))
}

func HandleTVPerson(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 1)
	u := fmt.Sprintf("https://api.tvmaze.com/people/%d", id)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name      string `json:"name"`
		Birthday  string `json:"birthday"`
		Deathday  string `json:"deathday"`
		Gender    string `json:"gender"`
		Country   struct {
			Name string `json:"name"`
		} `json:"country"`
	}
	json.Unmarshal(body, &data)
	d := ""
	if data.Deathday != "" {
		d = " - " + data.Deathday
	}
	return ok(fmt.Sprintf("Person: %s\nBorn: %s%s\nGender: %s\nCountry: %s",
		data.Name, data.Birthday, d, data.Gender, data.Country.Name))
}

func HandleShowCast(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	showID, _ := getInt(args, "showId", 1)
	u := fmt.Sprintf("https://api.tvmaze.com/shows/%d/cast", showID)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return err("TVMaze: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var cast []struct {
		Person struct {
			Name string `json:"name"`
		} `json:"person"`
		Character struct {
			Name string `json:"name"`
		} `json:"character"`
	}
	json.Unmarshal(body, &cast)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Cast (%d members):\n\n", len(cast)))
	for _, c := range cast {
		sb.WriteString(fmt.Sprintf("%s as %s\n", c.Person.Name, c.Character.Name))
	}
	return ok(sb.String())
}

func HandleWorldBankCountry(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	code, _ := getString(args, "code")
	if code == "" {
		code = "US"
	}
	u := fmt.Sprintf("https://api.worldbank.org/v2/country/%s?format=json", strings.ToUpper(code))
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("World Bank info for %s unavailable", code))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw []interface{}
	json.Unmarshal(body, &raw)
	if len(raw) < 2 {
		return ok(fmt.Sprintf("No data for %s", code))
	}
	entries, _ := raw[1].([]interface{})
	if len(entries) == 0 {
		return ok(fmt.Sprintf("No data for %s", code))
	}
	entry, _ := entries[0].(map[string]interface{})
	name, _ := entry["name"].(string)
	region, _ := entry["region"].(map[string]interface{})
	regName, _ := region["value"].(string)
	income, _ := entry["incomeLevel"].(map[string]interface{})
	incName, _ := income["value"].(string)
	capital, _ := entry["capitalCity"].(string)
	lat, _ := entry["latitude"].(string)
	lon, _ := entry["longitude"].(string)
	return ok(fmt.Sprintf("Country: %s (%s)\nRegion: %s\nIncome: %s\nCapital: %s\nCoords: %s, %s",
		name, code, regName, incName, capital, lat, lon))
}

func HandleAuthorBio(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	authorKey, _ := getString(args, "authorKey")
	if authorKey == "" {
		return err("authorKey is required (e.g. 'OL28618A')")
	}
	u := fmt.Sprintf("https://openlibrary.org/authors/%s.json", authorKey)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name       string `json:"name"`
		BirthDate  string `json:"birth_date"`
		DeathDate  string `json:"death_date"`
		Bio        string `json:"bio"`
		Wikipedia  string `json:"wikipedia"`
		Personal   string `json:"personal_name"`
	}
	json.Unmarshal(body, &data)
	bio := data.Bio
	if len(bio) > 500 {
		bio = bio[:500] + "..."
	}
	return ok(fmt.Sprintf("Author: %s (%s)\nBorn: %s | Died: %s\nWikipedia: %s\n\n%s",
		data.Name, data.Personal, data.BirthDate, data.DeathDate, data.Wikipedia, bio))
}

func HandleBookCover(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	olid, _ := getString(args, "olid")
	isbn, _ := getString(args, "isbn")
	size, _ := getString(args, "size")
	if olid == "" && isbn == "" {
		return err("olid or isbn is required")
	}
	if size == "" {
		size = "L"
	}
	var url string
	if isbn != "" {
		url = fmt.Sprintf("https://covers.openlibrary.org/b/isbn/%s-%s.jpg", isbn, size)
	} else {
		url = fmt.Sprintf("https://covers.openlibrary.org/b/olid/%s-%s.jpg", olid, size)
	}
	return ok(fmt.Sprintf("Book Cover: %s\n(Sizes: S=small, M=medium, L=large)", url))
}

func HandleOpenLibraryUser(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	username, _ := getString(args, "username")
	if username == "" {
		username = "george08"
	}
	u := fmt.Sprintf("https://openlibrary.org/people/%s.json", username)
	resp, apiErr := batch12.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Open Library user %q not found", username))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		DisplayName string `json:"displayname"`
		Created     string `json:"created"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("User: %s\nJoined: %s\nUsername: %s", data.DisplayName, data.Created[:10], username))
}
