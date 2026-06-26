package tools

import (
    "context"
    "fmt"
    "net/http"
    "time"
)

// HandleExample1 is an example handler function.
func HandleExample1(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    val, _ :=getString(args, "key")
    if val == "" {
        return err("Key not found in args")
}

    return ok(fmt.Sprintf("Received value: %s", val))
}

// HandleExample2 is another example handler function.
func HandleExample2(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    client := http.DefaultClient
    req, e := http.NewRequest("GET", "https://example.com", nil)
    if e != nil {
        return err(e.Error())
}

    resp, e := client.Do(req)
    if e != nil {
        return err(e.Error())
}

    defer resp.Body.Close()
    return ok(fmt.Sprintf("Response status: %s", resp.Status))
}