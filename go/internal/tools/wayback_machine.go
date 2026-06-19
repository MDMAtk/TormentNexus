var result map[string]interface{}
if e := json.Unmarshal(body, &result); e != nil {
    return ok(string(body))
}

archived, found := result["archived"].(bool)
if !found {
    return ok(string(body))
}

if archived {
    closest, _ := result["url"].(string)
    timestamp, _ := result["timestamp"].(string)
    return ok(fmt.Sprintf("URL is archived.\nClosest Snapshot: %s\nTimestamp: %s", closest, timestamp))
}

var waybackClient = http.DefaultClient
endpoint := fmt.Sprintf("https://archive.org/wayback/available?url=%s", url.QueryEscape(target))
limit, _ :=getInt(args, "limit")
if limit <= 0 {
    limit = 10
}