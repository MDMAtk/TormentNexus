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

var batch6API = &http.Client{Timeout: 15 * time.Second}

// ═══════════════════════════════════════════════════════════════════
// 50. NAGER.DATE — Public holidays by year/country (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandlePublicHolidays(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	year, _ := getInt(args, "year", 0)
	country, _ := getString(args, "country")
	if year == 0 {
		year = time.Now().Year()
	}
	if country == "" {
		country = "US"
	}
	u := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%d/%s", year, strings.ToUpper(country))
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return err("Nager.Date: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var holidays []struct {
		Date        string   `json:"date"`
		LocalName   string   `json:"localName"`
		Name        string   `json:"name"`
		CountryCode string   `json:"countryCode"`
		Global      bool     `json:"global"`
		Types       []string `json:"types"`
	}
	json.Unmarshal(body, &holidays)
	if len(holidays) == 0 {
		return ok(fmt.Sprintf("No holidays found for %s in %d", country, year))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Public Holidays — %s %d (%d total):\n\n", strings.ToUpper(country), year, len(holidays)))
	for _, h := range holidays {
		sb.WriteString(fmt.Sprintf("%s: %s\n", h.Date, h.LocalName))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 51. CAT FACTS — Random cat facts (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleCatFact(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	n, _ := getInt(args, "count", 1)
	var u string
	if n > 1 {
		u = fmt.Sprintf("https://catfact.ninja/facts?limit=%d", n)
	} else {
		u = "https://catfact.ninja/fact"
	}
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return ok("Cats have 32 muscles in each ear. 🐱")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if n > 1 {
		var data struct {
			Data []struct {
				Fact string `json:"fact"`
			} `json:"data"`
		}
		json.Unmarshal(body, &data)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("🐱 Cat Facts (%d):\n\n", len(data.Data)))
		for i, d := range data.Data {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, d.Fact))
		}
		return ok(sb.String())
	}
	var data struct {
		Fact string `json:"fact"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("🐱 Cat Fact: %s", data.Fact))
}

// ═══════════════════════════════════════════════════════════════════
// 52. SPACEFLIGHT NEWS — Space news articles (free)
// ═══════════════════════════════════════════════════════════════════

func HandleSpaceNews(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	n, _ := getInt(args, "limit", 5)
	u := fmt.Sprintf("https://api.spaceflightnewsapi.net/v4/articles/?limit=%d&ordering=-published_at", n)
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return err("Spaceflight News: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Count   int `json:"count"`
		Results []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Summary     string `json:"summary"`
			PublishedAt string `json:"published_at"`
			NewsSite    string `json:"news_site"`
		} `json:"results"`
	}
	json.Unmarshal(body, &data)
	if len(data.Results) == 0 {
		return ok("No space news available")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🚀 Space News (%d articles, showing %d):\n\n", data.Count, len(data.Results)))
	for i, a := range data.Results {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, a.Title))
		s := a.Summary
		if len(s) > 200 {
			s = s[:200] + "..."
		}
		sb.WriteString(fmt.Sprintf("   %s — %s (%s)\n", s, a.NewsSite, a.PublishedAt[:10]))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 53. OPEN NOTIFY — People currently in space (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandlePeopleInSpace(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch6API.Get("http://api.open-notify.org/astros.json")
	if apiErr != nil {
		return ok("🌍 People in space: [service unavailable]")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Number int `json:"number"`
		People []struct {
			Name  string `json:"name"`
			Craft string `json:"craft"`
		} `json:"people"`
	}
	json.Unmarshal(body, &data)
	if data.Number == 0 {
		return ok("🌍 No people currently in space.")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🌍 People in Space: %d\n\n", data.Number))
	for _, p := range data.People {
		sb.WriteString(fmt.Sprintf("  👨‍🚀 %s (%s)\n", p.Name, p.Craft))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 54. BIBLE API — Random Bible verses (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleBibleVerse(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	reference, _ := getString(args, "reference")
	var u string
	if reference != "" {
		u = fmt.Sprintf("https://bible-api.com/%s", url.QueryEscape(reference))
	} else {
		u = "https://bible-api.com/?random=verse"
	}
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return ok("\"For I know the plans I have for you, declares the LORD.\" — Jeremiah 29:11")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Reference string `json:"reference"`
		Text      string `json:"text"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("📖 %s\n\n%s", data.Reference, strings.TrimSpace(data.Text)))
}

// ═══════════════════════════════════════════════════════════════════
// 55. DATAMUSE — Synonyms, antonyms, rhymes (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleThesaurus(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	word, _ := getString(args, "word")
	rel, _ := getString(args, "relation")
	if word == "" {
		return err("word is required")
	}
	if rel == "" {
		rel = "syn"
	}
	u := fmt.Sprintf("https://api.datamuse.com/words?rel_%s=%s&max=10", rel, url.QueryEscape(word))
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return err("Datamuse: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var words []struct {
		Word  string `json:"word"`
		Score int    `json:"score"`
	}
	json.Unmarshal(body, &words)
	if len(words) == 0 {
		return ok(fmt.Sprintf("No %s found for \"%s\"", rel, word))
	}
	var sb strings.Builder
	relName := map[string]string{"syn": "Synonyms", "ant": "Antonyms", "rhy": "Rhymes", "trg": "Related"}[rel]
	sb.WriteString(fmt.Sprintf("📚 %s for \"%s\":\n\n", relName, word))
	for _, w := range words {
		sb.WriteString(fmt.Sprintf("  %s (score: %d)\n", w.Word, w.Score))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 56. EMOJI HUB — Random emoji with meaning (free)
// ═══════════════════════════════════════════════════════════════════

func HandleRandomEmoji(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch6API.Get("https://emojihub.yurace.pro/api/random")
	if apiErr != nil {
		return ok("Emoji: 😊 (service unavailable)")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Name     string   `json:"name"`
		Category string   `json:"category"`
		Group    string   `json:"group"`
		HTMLCode []string `json:"htmlCode"`
	}
	json.Unmarshal(body, &data)
	emoji := ""
	if len(data.HTMLCode) > 0 {
		emoji = data.HTMLCode[0]
	}
	return ok(fmt.Sprintf("Emoji: %s\nName: %s\nCategory: %s\nGroup: %s", emoji, data.Name, data.Category, data.Group))
}

// ═══════════════════════════════════════════════════════════════════
// 57. USELESS FACTS — Random interesting facts (free)
// ═══════════════════════════════════════════════════════════════════

func HandleUselessFact(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch6API.Get("https://uselessfacts.jsph.pl/random.json?language=en")
	if apiErr != nil {
		return ok("Did you know? The Eiffel Tower can be 15 cm taller during the summer.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Text   string `json:"text"`
		Source string `json:"source"`
	}
	json.Unmarshal(body, &data)
	extra := ""
	if data.Source != "" {
		extra = fmt.Sprintf(" (source: %s)", data.Source)
	}
	return ok(fmt.Sprintf("💡 Random Fact:\n%s%s", data.Text, extra))
}

// ═══════════════════════════════════════════════════════════════════
// 58. GEEK JOKES — Programming/geek humor (free)
// ═══════════════════════════════════════════════════════════════════

func HandleGeekJoke(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := batch6API.Get("https://geek-jokes.sameerkumar.website/api?format=json")
	if apiErr != nil {
		return ok("Why do programmers prefer dark mode? Because light attracts bugs.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Joke string `json:"joke"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("🤓 Geek Joke: %s", data.Joke))
}

// ═══════════════════════════════════════════════════════════════════
// 59. OPEN BREWERY DB — Brewery search (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleBrewerySearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	city, _ := getString(args, "city")
	state, _ := getString(args, "state")
	name, _ := getString(args, "name")
	n, _ := getInt(args, "limit", 5)

	var u string
	if name != "" {
		u = fmt.Sprintf("https://api.openbrewerydb.org/v1/breweries?by_name=%s&per_page=%d", url.QueryEscape(name), n)
	} else if city != "" && state != "" {
		u = fmt.Sprintf("https://api.openbrewerydb.org/v1/breweries?by_city=%s&by_state=%s&per_page=%d", url.QueryEscape(city), url.QueryEscape(state), n)
	} else if city != "" {
		u = fmt.Sprintf("https://api.openbrewerydb.org/v1/breweries?by_city=%s&per_page=%d", url.QueryEscape(city), n)
	} else {
		u = fmt.Sprintf("https://api.openbrewerydb.org/v1/breweries?per_page=%d", n)
	}

	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return err("Open Brewery DB: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var breweries []struct {
		Name        string `json:"name"`
		BreweryType string `json:"brewery_type"`
		Street      string `json:"street"`
		City        string `json:"city"`
		State       string `json:"state"`
		WebsiteURL  string `json:"website_url"`
		Phone       string `json:"phone"`
	}
	json.Unmarshal(body, &breweries)
	if len(breweries) == 0 {
		return ok("No breweries found")
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🍺 Breweries (%d found):\n\n", len(breweries)))
	for _, b := range breweries {
		sb.WriteString(fmt.Sprintf("  %s (%s)\n", b.Name, b.BreweryType))
		sb.WriteString(fmt.Sprintf("  %s, %s, %s\n", b.Street, b.City, b.State))
		if b.WebsiteURL != "" {
			sb.WriteString(fmt.Sprintf("  %s\n", b.WebsiteURL))
		}
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 60. D&D 5e API — Monster and spell data (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleDnDMonster(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return err("name is required (e.g. 'aboleth', 'dragon', 'goblin')")
	}
	u := fmt.Sprintf("https://www.dnd5eapi.co/api/monsters/%s", strings.ToLower(strings.ReplaceAll(name, " ", "-")))
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("D&D Monster %q not found. Try: aboleth, ancient-red-dragon, goblin", name))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Name      string `json:"name"`
		Size      string `json:"size"`
		Type      string `json:"type"`
		Alignment string `json:"alignment"`
		AC        []struct {
			Value int `json:"value"`
		} `json:"armor_class"`
		HP              int               `json:"hit_points"`
		Speed           map[string]string `json:"speed"`
		ChallengeRating float64           `json:"challenge_rating"`
		XP              int               `json:"xp"`
	}
	json.Unmarshal(body, &data)
	ac := 0
	if len(data.AC) > 0 {
		ac = data.AC[0].Value
	}
	speed := ""
	for k, v := range data.Speed {
		speed += fmt.Sprintf("%s %s ", k, v)
	}
	return ok(fmt.Sprintf("🐉 D&D Monster: %s\nSize: %s | Type: %s | Alignment: %s\nAC: %d | HP: %d | Speed: %s\nCR: %.1f (%d XP)\n",
		data.Name, data.Size, data.Type, data.Alignment, ac, data.HP, speed, data.ChallengeRating, data.XP))
}

// ═══════════════════════════════════════════════════════════════════
// 61. POKEMON TCG — Trading card search (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandlePokemonTCG(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	set, _ := getString(args, "set")
	n, _ := getInt(args, "limit", 5)
	filter := ""
	if name != "" {
		filter = fmt.Sprintf("name:%s", name)
	}
	if set != "" {
		if filter != "" {
			filter += " "
		}
		filter += fmt.Sprintf("set.name:%s", url.QueryEscape(set))
	}
	if filter == "" {
		filter = "name:pikachu"
	}
	u := fmt.Sprintf("https://api.pokemontcg.io/v2/cards?q=%s&pageSize=%d", url.QueryEscape(filter), n)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("User-Agent", "TormentNexus/1.0")

	resp, apiErr := batch6API.Do(req)
	if apiErr != nil {
		return ok("Pokémon TCG: service unavailable")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Data []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
			Set  struct {
				Name string `json:"name"`
			} `json:"set"`
			Images struct {
				Small string `json:"small"`
			} `json:"images"`
		} `json:"data"`
		TotalCount int `json:"totalCount"`
	}
	json.Unmarshal(body, &data)
	if len(data.Data) == 0 {
		return ok(fmt.Sprintf("No Pokémon cards found for %s", name))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🃏 Pokémon TCG — %d cards:\n\n", data.TotalCount))
	for _, c := range data.Data {
		sb.WriteString(fmt.Sprintf("  %s (%s) — set: %s\n", c.Name, c.ID, c.Set.Name))
		sb.WriteString(fmt.Sprintf("  %s\n", c.Images.Small))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 62. FAKER API — Generate fake companies, persons, products (free)
// ═══════════════════════════════════════════════════════════════════

func HandleFakerData(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	kind, _ := getString(args, "type")
	n, _ := getInt(args, "count", 2)
	if kind == "" {
		kind = "companies"
	}
	u := fmt.Sprintf("https://fakerapi.it/api/v1/%s?_quantity=%d", kind, n)
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return err("FakerAPI: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Status string        `json:"status"`
		Code   int           `json:"code"`
		Total  int           `json:"total"`
		Data   []interface{} `json:"data"`
	}
	json.Unmarshal(body, &data)
	if len(data.Data) == 0 {
		return ok(fmt.Sprintf("Faker: no %s generated", kind))
	}
	pretty, _ := json.MarshalIndent(data.Data, "", "  ")
	return ok(fmt.Sprintf("📝 Faker — %d %s:\n\n%s", len(data.Data), kind, string(pretty)))
}

// ═══════════════════════════════════════════════════════════════════
// 63. OPEN LIBRARY EDITION — Book edition by ISBN (free)
// ═══════════════════════════════════════════════════════════════════

func HandleBookEdition(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	isbn, _ := getString(args, "isbn")
	if isbn == "" {
		return err("isbn is required (10 or 13 digit ISBN)")
	}
	u := fmt.Sprintf("https://openlibrary.org/isbn/%s.json", isbn)
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return err("Open Library: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		Title         string   `json:"title"`
		Publishers    []string `json:"publishers"`
		PublishDate   string   `json:"publish_date"`
		NumberOfPages int      `json:"number_of_pages"`
		Authors       []struct {
			Key string `json:"key"`
		} `json:"authors"`
		Covers []int `json:"covers"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Book Edition — ISBN %s\n", isbn))
	sb.WriteString(fmt.Sprintf("Title: %s\n", data.Title))
	if len(data.Publishers) > 0 {
		sb.WriteString(fmt.Sprintf("Publisher: %s\n", strings.Join(data.Publishers, ", ")))
	}
	sb.WriteString(fmt.Sprintf("Published: %s\n", data.PublishDate))
	if data.NumberOfPages > 0 {
		sb.WriteString(fmt.Sprintf("Pages: %d\n", data.NumberOfPages))
	}
	for i := range data.Authors {
		sb.WriteString(fmt.Sprintf("Author %d: %s\n", i+1, data.Authors[i].Key))
	}
	if len(data.Covers) > 0 {
		sb.WriteString(fmt.Sprintf("Cover ID: %d", data.Covers[0]))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 64. DICTIONARY API — Word definitions (free)
// ═══════════════════════════════════════════════════════════════════

func HandleWordDefinition(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	word, _ := getString(args, "word")
	if word == "" {
		return err("word is required")
	}
	u := fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", url.PathEscape(word))
	resp, apiErr := batch6API.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("Definition for \"%s\": [service unavailable]", word))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var entries []struct {
		Word     string `json:"word"`
		Phonetic string `json:"phonetic"`
		Meanings []struct {
			PartOfSpeech string `json:"partOfSpeech"`
			Definitions  []struct {
				Definition string   `json:"definition"`
				Example    string   `json:"example"`
				Synonyms   []string `json:"synonyms"`
			} `json:"definitions"`
		} `json:"meanings"`
	}
	json.Unmarshal(body, &entries)
	if len(entries) == 0 {
		return ok(fmt.Sprintf("No definition found for \"%s\"", word))
	}
	e := entries[0]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📖 %s (%s)\n", e.Word, e.Phonetic))
	for _, m := range e.Meanings {
		for _, d := range m.Definitions {
			sb.WriteString(fmt.Sprintf("\n[%s] %s\n", m.PartOfSpeech, d.Definition))
			if d.Example != "" {
				sb.WriteString(fmt.Sprintf("   \"%s\"\n", d.Example))
			}
			if len(d.Synonyms) > 0 {
				s := d.Synonyms
				if len(s) > 5 {
					s = s[:5]
				}
				sb.WriteString(fmt.Sprintf("   Synonyms: %s\n", strings.Join(s, ", ")))
			}
			break
		}
	}
	return ok(sb.String())
}
