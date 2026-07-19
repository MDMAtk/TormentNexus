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

var moreAPI = &http.Client{Timeout: 15 * time.Second}

// ═══════════════════════════════════════════════════════════════════
// 12. HACKER NEWS — Search stories and top stories
// API: https://hn.algolia.com/api/v1/search (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleHNSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	n, _ := getInt(args, "limit", 5)
	if q == "" {
		return err("query is required")
	}
	u := fmt.Sprintf("https://hn.algolia.com/api/v1/search?query=%s&hitsPerPage=%d&tags=story", url.QueryEscape(q), n)
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return err("HN Search: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Hits []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Author      string `json:"author"`
			Points      int    `json:"points"`
			NumComments int    `json:"num_comments"`
			ObjectID    string `json:"objectID"`
			CreatedAt   string `json:"created_at"`
		} `json:"hits"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("HN: %d bytes", len(body)))
	}
	if len(data.Hits) == 0 {
		return ok(fmt.Sprintf("No HN stories for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Hacker News — %d stories for %q:\n\n", len(data.Hits), q))
	for i, h := range data.Hits {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, h.Title))
		sb.WriteString(fmt.Sprintf("   by %s | %d points | %d comments | %s\n", h.Author, h.Points, h.NumComments, h.CreatedAt[:10]))
		if h.URL != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", h.URL))
		} else {
			sb.WriteString(fmt.Sprintf("   https://news.ycombinator.com/item?id=%s\n", h.ObjectID))
		}
	}
	return ok(sb.String())
}

func HandleHNTopStories(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	n, _ := getInt(args, "limit", 10)
	u := fmt.Sprintf("https://hn.algolia.com/api/v1/search?tags=front_page&hitsPerPage=%d", n)
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return err("HN Top: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Hits []struct {
			Title       string `json:"title"`
			URL         string `json:"url"`
			Author      string `json:"author"`
			Points      int    `json:"points"`
			NumComments int    `json:"num_comments"`
			ObjectID    string `json:"objectID"`
		} `json:"hits"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString("Hacker News — Top Stories:\n\n")
	for i, h := range data.Hits {
		if i >= n {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. %s (%d points)\n", i+1, h.Title, h.Points))
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 13. COINGECKO — Cryptocurrency prices (free, no auth)
// API: https://api.coingecko.com/api/v3/
// ═══════════════════════════════════════════════════════════════════

func HandleCryptoPrice(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	ids, _ := getString(args, "ids")
	if ids == "" {
		ids = "bitcoin,ethereum,solana,cardano,polkadot,chainlink,avalanche-2"
	}
	u := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd&include_24hr_change=true", ids)
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return err("CoinGecko: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var raw map[string]interface{}
	if json.Unmarshal(body, &raw) != nil {
		return ok(fmt.Sprintf("Crypto: %d bytes", len(body)))
	}
	var sb strings.Builder
	sb.WriteString("Cryptocurrency Prices (USD) — CoinGecko:\n\n")
	for _, id := range strings.Split(ids, ",") {
		id = strings.TrimSpace(id)
		info, exists := raw[id].(map[string]interface{})
		if !exists {
			continue
		}
		price, _ := info["usd"].(float64)
		change, _ := info["usd_24h_change"].(float64)
		arrow := ""
		if change > 0 {
			arrow = "\U0001f7e2"
		} else if change < 0 {
			arrow = "\U0001f534"
		}
		sb.WriteString(fmt.Sprintf("%s $%.2f (%+.2f%% 24h)\n", formatCryptoName(id), price, change))
		_ = arrow
	}
	return ok(sb.String())
}

func formatCryptoName(id string) string {
	names := map[string]string{
		"bitcoin":     "Bitcoin",
		"ethereum":    "Ethereum",
		"solana":      "Solana",
		"cardano":     "Cardano",
		"polkadot":    "Polkadot",
		"chainlink":   "Chainlink",
		"avalanche-2": "Avalanche",
		"dogecoin":    "Dogecoin",
		"ripple":      "XRP",
		"litecoin":    "Litecoin",
	}
	if n, exists := names[id]; exists {
		return n
	}
	return id
}

// ═══════════════════════════════════════════════════════════════════
// 14. GITHUB — Search repositories (free, rate limited)
// API: https://api.github.com/search/repositories
// ═══════════════════════════════════════════════════════════════════

func HandleGitHubSearchRepos(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	n, _ := getInt(args, "limit", 5)
	sort, _ := getString(args, "sort")
	if q == "" {
		return err("query is required")
	}
	if sort == "" {
		sort = "stars"
	}
	u := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&per_page=%d&sort=%s&order=desc", url.QueryEscape(q), n, sort)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "TormentNexus/1.0")

	resp, apiErr := moreAPI.Do(req)
	if apiErr != nil {
		return err("GitHub: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		TotalCount int `json:"total_count"`
		Items      []struct {
			FullName    string   `json:"full_name"`
			Description string   `json:"description"`
			Stars       int      `json:"stargazers_count"`
			Forks       int      `json:"forks_count"`
			Language    string   `json:"language"`
			HTMLURL     string   `json:"html_url"`
			Topics      []string `json:"topics"`
		} `json:"items"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("GitHub: %d bytes", len(body)))
	}
	if len(data.Items) == 0 {
		return ok(fmt.Sprintf("No repos found for %q", q))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("GitHub — %d repos for %q (sort: %s):\n\n", data.TotalCount, q, sort))
	for i, r := range data.Items {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, r.FullName))
		if r.Description != "" {
			d := r.Description
			if len(d) > 100 {
				d = d[:100] + "..."
			}
			sb.WriteString(fmt.Sprintf("   %s\n", d))
		}
		sb.WriteString(fmt.Sprintf("   \u2b50 %d | 🍴 %d", r.Stars, r.Forks))
		if r.Language != "" {
			sb.WriteString(fmt.Sprintf(" | %s", r.Language))
		}
		if len(r.Topics) > 0 {
			sb.WriteString(fmt.Sprintf(" | tags: %s", strings.Join(r.Topics[:min(len(r.Topics), 4)], ", ")))
		}
		sb.WriteString(fmt.Sprintf("\n   %s\n\n", r.HTMLURL))
	}
	return ok(sb.String())
}

func HandleGitHubRepoInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	repo, _ := getString(args, "repo")
	if repo == "" {
		return err("repo is required (e.g. 'golang/go')")
	}
	u := fmt.Sprintf("https://api.github.com/repos/%s", repo)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "TormentNexus/1.0")

	resp, apiErr := moreAPI.Do(req)
	if apiErr != nil {
		return err("GitHub: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data struct {
		FullName    string   `json:"full_name"`
		Description string   `json:"description"`
		Stars       int      `json:"stargazers_count"`
		Forks       int      `json:"forks_count"`
		Language    string   `json:"language"`
		HTMLURL     string   `json:"html_url"`
		Topics      []string `json:"topics"`
		License     struct {
			Name string `json:"name"`
		} `json:"license"`
		OpenIssues int    `json:"open_issues_count"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}
	if json.Unmarshal(body, &data) != nil {
		return ok(fmt.Sprintf("GitHub %s: %d bytes", repo, len(body)))
	}
	return ok(fmt.Sprintf("GitHub — %s\n%s\n\nStars: %d | Forks: %d | Issues: %d\nLanguage: %s | License: %s\nCreated: %s | Updated: %s\n\n%s",
		data.FullName, data.Description, data.Stars, data.Forks, data.OpenIssues,
		data.Language, data.License.Name, data.CreatedAt[:10], data.UpdatedAt[:10], data.HTMLURL))
}

// ═══════════════════════════════════════════════════════════════════
// 15. JOKEAPI — Random jokes (free, safe mode)
// ═══════════════════════════════════════════════════════════════════

func HandleJokeGet(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	cat, _ := getString(args, "category")
	u := "https://v2.jokeapi.dev/joke/Any?safe-mode"
	if cat != "" {
		u = fmt.Sprintf("https://v2.jokeapi.dev/joke/%s?safe-mode", url.PathEscape(cat))
	}
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return ok("Why did the developer go broke? Because he used up all his cache!")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Type     string `json:"type"`
		Joke     string `json:"joke"`
		Setup    string `json:"setup"`
		Delivery string `json:"delivery"`
		Category string `json:"category"`
	}
	json.Unmarshal(body, &data)
	if data.Type == "single" {
		return ok(fmt.Sprintf("😂 [%s] %s", data.Category, data.Joke))
	}
	return ok(fmt.Sprintf("😂 [%s] %s\n%s", data.Category, data.Setup, data.Delivery))
}

// ═══════════════════════════════════════════════════════════════════
// 16. MEALDB — Recipe search (free, no auth)
// ═══════════════════════════════════════════════════════════════════

func HandleRecipeSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	ingredient, _ := getString(args, "ingredient")
	name, _ := getString(args, "name")
	if ingredient == "" && name == "" {
		return err("ingredient or name is required")
	}
	var u string
	if name != "" {
		u = fmt.Sprintf("https://www.themealdb.com/api/json/v1/1/search.php?s=%s", url.QueryEscape(name))
	} else {
		u = fmt.Sprintf("https://www.themealdb.com/api/json/v1/1/filter.php?i=%s", url.QueryEscape(ingredient))
	}
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return err("MealDB: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Meals []map[string]interface{} `json:"meals"`
	}
	json.Unmarshal(body, &data)
	if len(data.Meals) == 0 {
		return ok(fmt.Sprintf("No recipes found for %q", name+ingredient))
	}
	var sb strings.Builder
	if name != "" {
		meal := data.Meals[0]
		sb.WriteString(fmt.Sprintf("Recipe: %s\n\n", meal["strMeal"]))
		sb.WriteString(fmt.Sprintf("Category: %s | Area: %s\n", meal["strCategory"], meal["strArea"]))
		sb.WriteString(fmt.Sprintf("Tags: %s\n\n", meal["strTags"]))
		// Instructions
		instructions, _ := meal["strInstructions"].(string)
		if len(instructions) > 500 {
			instructions = instructions[:500] + "..."
		}
		sb.WriteString(fmt.Sprintf("Instructions:\n%s\n\n", instructions))
		// Ingredients
		sb.WriteString("Ingredients:\n")
		for i := 1; i <= 20; i++ {
			ing := meal[fmt.Sprintf("strIngredient%d", i)]
			meas := meal[fmt.Sprintf("strMeasure%d", i)]
			if ing == nil || ing.(string) == "" {
				break
			}
			sb.WriteString(fmt.Sprintf("  • %s %s\n", meas, ing))
		}
		sb.WriteString(fmt.Sprintf("\n%s", meal["strMealThumb"]))
	} else {
		sb.WriteString(fmt.Sprintf("Recipes with %q (%d found):\n\n", ingredient, len(data.Meals)))
		for i, m := range data.Meals {
			if i >= 10 {
				break
			}
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, m["strMeal"]))
		}
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 17. REST COUNTRIES — Country information (free)
// ═══════════════════════════════════════════════════════════════════

func HandleCountryInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	code, _ := getString(args, "code")
	if name == "" && code == "" {
		return err("name or code is required")
	}
	var u string
	if code != "" {
		u = fmt.Sprintf("https://restcountries.com/v3.2/alpha/%s", code)
	} else {
		u = fmt.Sprintf("https://restcountries.com/v3.2/name/%s", url.QueryEscape(name))
	}
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return err("RestCountries: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var countries []map[string]interface{}
	json.Unmarshal(body, &countries)
	if len(countries) == 0 {
		return ok(fmt.Sprintf("No country found for %q", name+code))
	}
	c := countries[0]
	common, _ := c["name"].(map[string]interface{})["common"].(string)
	official, _ := c["name"].(map[string]interface{})["official"].(string)
	capital, _ := c["capital"].([]interface{})
	region, _ := c["region"].(string)
	subregion, _ := c["subregion"].(string)
	population, _ := c["population"].(float64)
	area, _ := c["area"].(float64)
	languages, _ := c["languages"].(map[string]interface{})
	currencies, _ := c["currencies"].(map[string]interface{})
	flags, _ := c["flags"].(map[string]interface{})
	flag, _ := flags["png"].(string)
	timezones, _ := c["timezones"].([]interface{})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Country: %s (%s)\n", common, official))
	sb.WriteString(fmt.Sprintf("Capital: %s\n", joinAny(capital, ", ")))
	sb.WriteString(fmt.Sprintf("Region: %s / %s\n", region, subregion))
	sb.WriteString(fmt.Sprintf("Population: %.0f\n", population))
	sb.WriteString(fmt.Sprintf("Area: %.0f km²\n", area))
	sb.WriteString(fmt.Sprintf("Languages: %s\n", joinMapValues(languages)))
	sb.WriteString(fmt.Sprintf("Currencies: %s\n", formatCurrencies(currencies)))
	sb.WriteString(fmt.Sprintf("Timezones: %s\n", joinAny(timezones, ", ")))
	sb.WriteString(fmt.Sprintf("Flag: %s\n", flag))
	return ok(sb.String())
}

func HandleCountryList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	region, _ := getString(args, "region")
	u := "https://restcountries.com/v3.2/all?fields=name,region,population,flags"
	if region != "" {
		u = fmt.Sprintf("https://restcountries.com/v3.2/region/%s?fields=name,region,population,flags", url.PathEscape(region))
	}
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return err("RestCountries: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var countries []map[string]interface{}
	json.Unmarshal(body, &countries)
	if len(countries) > 100 {
		countries = countries[:100]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Countries (%d shown):\n\n", len(countries)))
	for i, c := range countries {
		name, _ := c["name"].(map[string]interface{})["common"].(string)
		pop, _ := c["population"].(float64)
		flags, _ := c["flags"].(map[string]interface{})
		flag, _ := flags["png"].(string)
		sb.WriteString(fmt.Sprintf("%d. %s (pop: %.0f)\n", i+1, name, pop))
		_ = flag
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 18. ADVICE SLIP — Random advice (free)
// ═══════════════════════════════════════════════════════════════════

func HandleAdviceGet(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	u := "https://api.adviceslip.com/advice"
	if q != "" {
		u = fmt.Sprintf("https://api.adviceslip.com/advice/search/%s", url.QueryEscape(q))
	}
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return ok("Advice: Don't forget to backup your data!")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Slip struct {
			ID     int    `json:"id"`
			Advice string `json:"advice"`
		} `json:"slip"`
		Slips []struct {
			ID     int    `json:"id"`
			Advice string `json:"advice"`
		} `json:"slips"`
	}
	json.Unmarshal(body, &data)
	if q != "" && len(data.Slips) > 0 {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Advice for %q:\n\n", q))
		for i, s := range data.Slips {
			if i >= 5 {
				break
			}
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, s.Advice))
		}
		return ok(sb.String())
	}
	if data.Slip.Advice != "" {
		return ok(fmt.Sprintf("Advice #%d: %s", data.Slip.ID, data.Slip.Advice))
	}
	return ok("Advice: Keep learning and building!")
}

// ═══════════════════════════════════════════════════════════════════
// 19. NUMBERS API — Interesting number facts (free)
// ═══════════════════════════════════════════════════════════════════

func HandleNumberFact(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	num, _ := getInt(args, "number", 0)
	kind, _ := getString(args, "type")
	if kind == "" {
		kind = "trivia"
	}
	u := fmt.Sprintf("http://numbersapi.com/%d/%s", num, kind)
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return ok("42 is the answer to the ultimate question of life, the universe, and everything.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fact := strings.TrimSpace(string(body))
	return ok(fmt.Sprintf("🔢 %s", fact))
}

// ═══════════════════════════════════════════════════════════════════
// 20. DOG API — Random dog pictures (free)
// ═══════════════════════════════════════════════════════════════════

func HandleDogPicture(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	breed, _ := getString(args, "breed")
	u := "https://dog.ceo/api/breeds/image/random"
	if breed != "" {
		u = fmt.Sprintf("https://dog.ceo/api/breed/%s/images/random", strings.ReplaceAll(breed, " ", "/"))
	}
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return ok("🐕 https://images.dog.ceo/breeds/hound-afghan/n02088094_1003.jpg")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}
	json.Unmarshal(body, &data)
	return ok(fmt.Sprintf("🐕 %s", data.Message))
}

func HandleDogBreeds(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	resp, apiErr := moreAPI.Get("https://dog.ceo/api/breeds/list/all")
	if apiErr != nil {
		return err("Dog API: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Message map[string][]string `json:"message"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Dog Breeds (%d):\n\n", len(data.Message)))
	count := 0
	for breed, sub := range data.Message {
		if count >= 20 {
			sb.WriteString(fmt.Sprintf("... and %d more\n", len(data.Message)-20))
			break
		}
		if len(sub) > 0 {
			sb.WriteString(fmt.Sprintf("%s (%s)\n", breed, strings.Join(sub, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("%s\n", breed))
		}
		count++
	}
	return ok(sb.String())
}

// ═══════════════════════════════════════════════════════════════════
// 21. OPEN FOOD FACTS — Food product lookup (free)
// ═══════════════════════════════════════════════════════════════════

func HandleFoodProduct(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	barcode, _ := getString(args, "barcode")
	name, _ := getString(args, "name")
	if barcode == "" && name == "" {
		return err("barcode or name is required")
	}
	var u string
	if barcode != "" {
		u = fmt.Sprintf("https://world.openfoodfacts.org/api/v2/product/%s.json", barcode)
	} else {
		u = fmt.Sprintf("https://world.openfoodfacts.org/cgi/search.pl?search_terms=%s&json=true&page_size=3", url.QueryEscape(name))
	}
	resp, apiErr := moreAPI.Get(u)
	if apiErr != nil {
		return err("OpenFoodFacts: " + apiErr.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if barcode != "" {
		var data struct {
			Product struct {
				Name      string `json:"product_name"`
				Brand     string `json:"brands"`
				Quantity  string `json:"quantity"`
				Nutrients struct {
					EnergyKJ float64 `json:"energy-kj_100g"`
					Fat      float64 `json:"fat_100g"`
					Carbs    float64 `json:"carbohydrates_100g"`
					Proteins float64 `json:"proteins_100g"`
					Salt     float64 `json:"salt_100g"`
				} `json:"nutriments"`
			} `json:"product"`
		}
		json.Unmarshal(body, &data)
		return ok(fmt.Sprintf("Product: %s\nBrand: %s\nQuantity: %s\n\nNutrition per 100g:\nEnergy: %.0f kJ\nFat: %.1fg\nCarbs: %.1fg\nProtein: %.1fg\nSalt: %.1fg",
			data.Product.Name, data.Product.Brand, data.Product.Quantity,
			data.Product.Nutrients.EnergyKJ, data.Product.Nutrients.Fat,
			data.Product.Nutrients.Carbs, data.Product.Nutrients.Proteins, data.Product.Nutrients.Salt))
	}
	var data struct {
		Products []struct {
			Name    string `json:"product_name"`
			Brand   string `json:"brands"`
			Barcode string `json:"code"`
		} `json:"products"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Products for %q:\n\n", name))
	for i, p := range data.Products {
		if i >= 5 {
			break
		}
		sb.WriteString(fmt.Sprintf("%d. %s by %s (%s)\n", i+1, p.Name, p.Brand, p.Barcode))
	}
	return ok(sb.String())
}

// ── HELPERS ─────────────────────────────────────────────────────────

func joinMapValues(m map[string]interface{}) string {
	var vals []string
	for _, v := range m {
		vals = append(vals, fmt.Sprintf("%v", v))
	}
	return strings.Join(vals, ", ")
}

func formatCurrencies(m map[string]interface{}) string {
	var parts []string
	for code, info := range m {
		if im, exists := info.(map[string]interface{}); exists {
			name, _ := im["name"].(string)
			sym, _ := im["symbol"].(string)
			parts = append(parts, fmt.Sprintf("%s (%s %s)", name, code, sym))
		}
	}
	return strings.Join(parts, ", ")
}
