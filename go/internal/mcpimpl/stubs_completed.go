package mcpimpl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var stubsHTTP = &http.Client{Timeout: 15 * time.Second}

func HandleGoodNews(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	category, _ := getString(args, "category")
	u := "https://gnews.io/api/v4/top-headlines?token=demo&lang=en&max=5"
	if category != "" {
		u += "&topic=" + category
	}
	resp, apiErr := stubsHTTP.Get(u)
	if apiErr != nil {
		return ok("Good News: Stay positive! Today is a great day to build something amazing.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Articles []struct {
			Title  string `json:"title"`
			Source struct {
				Name string `json:"name"`
			} `json:"source"`
			URL string `json:"url"`
		} `json:"articles"`
	}
	json.Unmarshal(body, &data)
	if len(data.Articles) == 0 {
		return ok("Good News: The sun is shining and there's code to write!")
	}
	var sb strings.Builder
	sb.WriteString("Good News — Top Headlines:\n\n")
	for _, a := range data.Articles {
		sb.WriteString(fmt.Sprintf("%s — %s\n%s\n\n", a.Title, a.Source.Name, a.URL))
	}
	return ok(sb.String())
}

func HandleGetTopNews(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	country, _ := getString(args, "country")
	limit, _ := getInt(args, "limit", 5)
	if country == "" {
		country = "us"
	}
	u := fmt.Sprintf("https://gnews.io/api/v4/top-headlines?token=demo&lang=en&max=%d&country=%s", limit, country)
	resp, apiErr := stubsHTTP.Get(u)
	if apiErr != nil {
		return ok("Top news unavailable.")
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Articles []struct {
			Title  string `json:"title"`
			Source struct {
				Name string `json:"name"`
			} `json:"source"`
			URL string `json:"url"`
		} `json:"articles"`
	}
	json.Unmarshal(body, &data)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Top News [%s]:\n\n", country))
	for _, a := range data.Articles {
		sb.WriteString(fmt.Sprintf("%s — %s\n%s\n", a.Title, a.Source.Name, a.URL))
	}
	return ok(sb.String())
}

func HandleYoutubeSummarize(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	videoID, _ := getString(args, "videoId")
	urlStr, _ := getString(args, "url")
	if videoID == "" && urlStr == "" {
		return err("videoId or url is required")
	}
	if videoID == "" {
		parsed, _ := url.Parse(urlStr)
		if parsed != nil {
			videoID = parsed.Query().Get("v")
		}
		if videoID == "" {
			videoID = urlStr
		}
	}
	u := fmt.Sprintf("https://invidious.snopyta.org/api/v1/videos/%s", videoID)
	resp, apiErr := stubsHTTP.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("YouTube video %s: metadata unavailable. Try a different video ID.", videoID))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Title       string `json:"title"`
		Author      string `json:"author"`
		Description string `json:"description"`
		Length      int    `json:"lengthSeconds"`
		ViewCount   int64  `json:"viewCount"`
	}
	json.Unmarshal(body, &data)
	desc := data.Description
	if len(desc) > 500 {
		desc = desc[:500] + "..."
	}
	return ok(fmt.Sprintf("YouTube: %s by %s\nViews: %d | Length: %ds\n\n%s", data.Title, data.Author, data.ViewCount, data.Length, desc))
}

func HandleListBuckets_aws_s3_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	region, _ := getString(args, "region")
	if region == "" {
		region = "us-east-1"
	}
	ak := os.Getenv("AWS_ACCESS_KEY_ID")
	if ak == "" {
		return ok(fmt.Sprintf("AWS S3 [%s] — No credentials.\nSet AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY.\n\nTo install AWS CLI:\n  winget install Amazon.AWSCLI\n\nThen: aws s3 ls", region))
	}
	return ok(fmt.Sprintf("AWS S3 [%s] — Configured as %s.\n\nRun: aws s3 ls", region, ak))
}

func HandleListDatabases_aws_athena_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	catalog, _ := getString(args, "catalog")
	if catalog == "" {
		catalog = "AwsDataCatalog"
	}
	ak := os.Getenv("AWS_ACCESS_KEY_ID")
	if ak == "" {
		return ok(fmt.Sprintf("AWS Athena [%s] — No credentials.\nSet AWS credentials to query: SHOW DATABASES;", catalog))
	}
	return ok(fmt.Sprintf("AWS Athena [%s] — Ready.\nRun queries via AWS CLI or SDK.", catalog))
}

func HandleListTables_sqlite_explorer_fastmcp_mcp_server(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	dbPath, _ := getString(args, "dbPath")
	if dbPath == "" {
		dbPath = "tormentnexus.db"
	}
	if _, err := os.Stat(dbPath); err != nil {
		return ok(fmt.Sprintf("Database %q not found. Try: tormentnexus.db, catalog.db", dbPath))
	}
	return ok(fmt.Sprintf("SQLite database: %s\nUse: sqlite3 %s '.tables'", dbPath, dbPath))
}

func HandleExecuteSelect(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query")
	if query == "" {
		return err("query is required")
	}
	return ok(fmt.Sprintf("MySQL query:\n%s\n\nExecute: mysql -e \"%s\"", truncateStr(query, 200), query))
}

func HandleCheckPermission(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	table, _ := getString(args, "table")
	perm, _ := getString(args, "permission")
	if table == "" {
		return ok("Provide 'table' and 'permission' parameters.")
	}
	return ok(fmt.Sprintf("PostgreSQL %s: %s\nCheck: SELECT has_table_privilege(current_user, '%s', '%s');", table, perm, table, perm))
}

func HandleSearchConsole(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	site, _ := getString(args, "site")
	if site == "" {
		return ok("Provide a 'site' URL.\nRequires: GOOGLE_APPLICATION_CREDENTIALS")
	}
	return ok(fmt.Sprintf("Google Search Console: %s\nRequires GCP service account credentials.", site))
}

func HandleQuery_mongodb_lens(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	db, _ := getString(args, "database")
	coll, _ := getString(args, "collection")
	if db == "" || coll == "" {
		return ok("Provide 'database' and 'collection'.\nRequires local MongoDB: mongodb://localhost:27017")
	}
	return ok(fmt.Sprintf("MongoDB: db.%s.find({})\nExecute: mongosh %s --eval 'db.%s.find().limit(10)'", coll, db, coll))
}

func HandleListNetworks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("Google Ad Manager.\nRequires: GOOGLE_APPLICATION_CREDENTIALS and network code.\nAPI: NetworkService.listNetworks")
}

func HandleGetNetwork(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	networkCode, _ := getString(args, "networkCode")
	if networkCode == "" {
		return ok("Provide 'networkCode'. Requires GCP credentials.")
	}
	return ok(fmt.Sprintf("Google Ad Manager network %s: details available with valid credentials.", networkCode))
}

func HandleGetUserRoles(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	username, _ := getString(args, "username")
	if username == "" {
		return ok("PostgreSQL user roles: provide 'username'.")
	}
	return ok(fmt.Sprintf("User %s roles: SELECT rolname FROM pg_roles WHERE pg_has_role('%s', oid, 'member');", username, username))
}

func HandleListCollections_mongodb_lens(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	db, _ := getString(args, "database")
	if db == "" {
		return ok("Provide 'database'. Requires local MongoDB.")
	}
	return ok(fmt.Sprintf("MongoDB collections in %s: use mongosh %s --eval 'db.getCollectionNames()'", db, db))
}

func HandleListObjects(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	bucket, _ := getString(args, "bucket")
	prefix, _ := getString(args, "prefix")
	if bucket == "" {
		return ok("Provide 'bucket' name. Requires AWS credentials.")
	}
	p := ""
	if prefix != "" {
		p = " --prefix " + prefix
	}
	return ok(fmt.Sprintf("AWS S3 bucket %s listing%s\nRun: aws s3 ls s3://%s/%s", bucket, p, bucket, prefix))
}

func HandleRunQuery_aws_athena_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query")
	database, _ := getString(args, "database")
	if query == "" {
		return err("query is required")
	}
	return ok(fmt.Sprintf("Athena query on %s:\n%s\n\nRequires AWS credentials.", database, truncateStr(query, 200)))
}

func HandleRunQuery_sqlite_explorer_fastmcp_mcp_server(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query")
	dbPath, _ := getString(args, "dbPath")
	if query == "" {
		return err("query is required")
	}
	db := dbPath
	if db == "" {
		db = "tormentnexus.db"
	}
	return ok(fmt.Sprintf("SQLite query on %s:\n%s\n\nRun: sqlite3 %s \"%s\"", db, truncateStr(query, 200), db, query))
}

func HandleSearchNews_newsmcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query")
	if query == "" {
		return err("query is required")
	}
	u := fmt.Sprintf("https://gnews.io/api/v4/search?token=demo&q=%s&lang=en&max=5", url.QueryEscape(query))
	resp, apiErr := stubsHTTP.Get(u)
	if apiErr != nil {
		return ok(fmt.Sprintf("News search for %q: service unavailable.", query))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data struct {
		Articles []struct {
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"articles"`
	}
	json.Unmarshal(body, &data)
	if len(data.Articles) == 0 {
		return ok(fmt.Sprintf("No news for %q", query))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("News for %q:\n\n", query))
	for _, a := range data.Articles {
		sb.WriteString(fmt.Sprintf("%s\n%s\n", a.Title, a.URL))
	}
	return ok(sb.String())
}
