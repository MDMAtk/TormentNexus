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

var batch7 = &http.Client{Timeout: 15 * time.Second}

func HandlePicsumPhoto(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	id, _ := getInt(args, "id", 0)
	w, _ := getInt(args, "width", 400)
	h, _ := getInt(args, "height", 300)
	var u string
	if id > 0 {
		u = fmt.Sprintf("https://picsum.photos/id/%d/info", id)
	} else {
		u = "https://picsum.photos/id/0/info"
	}
	resp, apiErr := batch7.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Picsum: https://picsum.photos/%d/%d", w, h))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		ID     int    `json:"id"`
		Author string `json:"author"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		URL    string `json:"url"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("Photo #%d by %s\nOriginal: %dx%d\nURL: %s\nDownload: https://picsum.photos/%d/%d",
		data.ID, data.Author, data.Width, data.Height, data.URL, w, h))
}

func HandleIsItChristmas(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch7.Get("https://isitchristmas.com/api")
	if apiErr != nil {
		return ok("Christmas check unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Christmasess []struct {
			IsChristmas string `json:"is_christmas"`
		} `json:"christmases"`
	}
	json.Unmarshal(body, &data)
	if len(data.Christmasess) > 0 && data.Christmasess[0].IsChristmas == "yes" {
		return ok("Yes! It IS Christmas! 🎄🎅")
	}
	return ok("No, it's not Christmas yet. 🎄 Check back on December 25!")
}

func HandleBookDetails(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	key, _ := getString(args, "key")
	if key == "" {
		return err("key is required (e.g. 'OL7353617M', '/works/OL45804W')")
	}
	u := fmt.Sprintf("https://openlibrary.org%s.json", key)
	resp, apiErr := batch7.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Title      string   `json:"title"`
		Publishers []string `json:"publishers"`
		PubDate    string   `json:"publish_date"`
		Pages      int      `json:"number_of_pages"`
		Subjects   []string `json:"subjects"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Book: %s\n", data.Title))
	if len(data.Publishers) > 0 {
		sb.WriteString(fmt.Sprintf("Publisher: %s\n", strings.Join(data.Publishers, ", ")))
	}
	sb.WriteString(fmt.Sprintf("Published: %s\n", data.PubDate))
	if data.Pages > 0 {
		sb.WriteString(fmt.Sprintf("Pages: %d\n", data.Pages))
	}
	if len(data.Subjects) > 5 {
		data.Subjects = data.Subjects[:5]
	}
	if len(data.Subjects) > 0 {
		sb.WriteString(fmt.Sprintf("Subjects: %s\n", strings.Join(data.Subjects, ", ")))
	}
	return ok(sb.String())
}

func HandleOfficialJoke(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	kind, _ := getString(args, "type")
	u := "https://official-joke-api.appspot.com/random_joke"
	if kind != "" {
		u = fmt.Sprintf("https://official-joke-api.appspot.com/jokes/%s/random", kind)
	}
	resp, apiErr := batch7.Get(u)
	if apiErr != nil {
		return ok("Why did the chicken cross the road? To get to the other side!")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(body), "[") {
		var jokes []struct {
			Setup     string `json:"setup"`
			Punchline string `json:"punchline"`
			Type      string `json:"type"`
		}
		json.Unmarshal(body, &jokes)
		if len(jokes) > 0 {
			return ok(fmt.Sprintf("[%s] %s\n%s", jokes[0].Type, jokes[0].Setup, jokes[0].Punchline))
		}
	}
	var joke struct {
		Setup     string `json:"setup"`
		Punchline string `json:"punchline"`
		Type      string `json:"type"`
	}
	json.Unmarshal(body, &joke)
	return ok(fmt.Sprintf("[%s] %s\n%s", joke.Type, joke.Setup, joke.Punchline))
}

func HandleFoxPicture(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch7.Get("https://randomfox.ca/floof/")
	if apiErr != nil {
		return ok("Fox picture unavailable. 🦊")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Image string `json:"image"`
		Link  string `json:"link"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("🦊 Random Fox:\n%s\n%s", data.Image, data.Link))
}

func HandleBookSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	limit, _ := getInt(args, "limit", 5)
	if q == "" {
		return err("query is required")
	}
	u := fmt.Sprintf("https://openlibrary.org/search.json?q=%s&limit=%d", url.QueryEscape(q), limit)
	resp, apiErr := batch7.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		NumFound int `json:"numFound"`
		Docs     []struct {
			Title   string   `json:"title"`
			Author  []string `json:"author_name"`
			Year    int      `json:"first_publish_year"`
			Subject []string `json:"subject"`
		} `json:"docs"`
	}
	json.Unmarshal(body, &data)
	if len(data.Docs) == 0 {
		return ok(fmt.Sprintf("No books found for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Open Library — %d results for %q:\n\n", data.NumFound, q))
	for _, d := range data.Docs {
		sb.WriteString(fmt.Sprintf("%s", d.Title))
		if len(d.Author) > 0 {
			sb.WriteString(fmt.Sprintf(" by %s", strings.Join(d.Author, ", ")))
		}
		if d.Year > 0 {
			sb.WriteString(fmt.Sprintf(" (%d)", d.Year))
		}
		sb.WriteString("\n")
	}
	return ok(sb.String())
}

func HandleZippoLookup(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	code, _ := getString(args, "code")
	country, _ := getString(args, "country")
	if code == "" {
		return err("code is required")
	}
	if country == "" {
		country = "us"
	}
	u := fmt.Sprintf("https://api.zippopotam.us/%s/%s", strings.ToLower(country), code)
	resp, apiErr := batch7.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Zip code %q in %s: not found", code, country))
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
	sb.WriteString(fmt.Sprintf("%s (%s):\n", data.PostCode, data.Country))
	for _, p := range data.Places {
		sb.WriteString(fmt.Sprintf("%s, %s (%s, %s)\n", p.PlaceName, p.State, p.Latitude, p.Longitude))
	}
	return ok(sb.String())
}

func HandleQRCode(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	data, _ := getString(args, "data")
	size, _ := getInt(args, "size", 200)
	bg, _ := getString(args, "bg")
	fg, _ := getString(args, "fg")
	if data == "" {
		return err("data is required (text or URL to encode)")
	}
	u := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=%dx%d&data=%s", size, size, url.QueryEscape(data))
	if bg != "" {
		u += "&bgcolor=" + strings.TrimPrefix(bg, "#")
	}
	if fg != "" {
		u += "&color=" + strings.TrimPrefix(fg, "#")
	}
	return ok(fmt.Sprintf("QR Code generated:\n%s\n\nScan this QR code to access: %s", u, data))
}

func HandleNekosImage(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	kind, _ := getString(args, "type")
	if kind == "" {
		kind = "neko"
	}
	u := fmt.Sprintf("https://nekos.best/api/v2/%s", kind)
	resp, apiErr := batch7.Get(u)
	if apiErr != nil {
		return ok("Anime image unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Results []struct {
			URL    string `json:"url"`
			Artist string `json:"artist_name"`
		} `json:"results"`
	}
	json.Unmarshal(body, &data)
	if len(data.Results) == 0 {
		return ok("No image found")
	}
	r := data.Results[0]
	result := fmt.Sprintf("Anime Image (%s):\n%s", kind, r.URL)
	if r.Artist != "" {
		result += fmt.Sprintf("\nArtist: %s", r.Artist)
	}
	return ok(result)
}

func HandleHTTPCat(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	code, _ := getInt(args, "code", 200)
	return ok(fmt.Sprintf("HTTP Cat — Status %d:\nhttps://http.cat/%d.jpg\n\n(Not all status codes have cat images. Try: 100, 200, 301, 404, 500)", code, code))
}

func HandleBooksBySubject(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	subject, _ := getString(args, "subject")
	limit, _ := getInt(args, "limit", 10)
	if subject == "" {
		return err("subject is required (e.g. 'science_fiction', 'fantasy', 'history')")
	}
	u := fmt.Sprintf("https://openlibrary.org/subjects/%s.json?limit=%d", url.PathEscape(subject), limit)
	resp, apiErr := batch7.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Name  string `json:"name"`
		Count int    `json:"work_count"`
		Works []struct {
			Title string `json:"title"`
			Year  int    `json:"first_publish_year"`
		} `json:"works"`
	}
	json.Unmarshal(body, &data)
	if len(data.Works) == 0 {
		return ok(fmt.Sprintf("No books for subject %q", subject))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Subject: %s (%d books)\n\n", data.Name, data.Count))
	for _, w := range data.Works {
		sb.WriteString(w.Title)
		if w.Year > 0 {
			sb.WriteString(fmt.Sprintf(" (%d)", w.Year))
		}
		sb.WriteString("\n")
	}
	return ok(sb.String())
}
