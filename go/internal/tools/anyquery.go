package tools

import (
    "context"
    "fmt"
    "os/exec"
    "strings"
)

// Define runAnyquery since it's not in parity.go.
func runAnyquery(ctx context.Context, args ...string) (string, error) {
    cmd := exec.CommandContext(ctx, "anyquery", args...)
    output, e := cmd.CombinedOutput()
    if e != nil {
        return string(output), e
    }
    return string(output), nil
}

func HandleRunQuery(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    query, _ :=getString(args, "query")
    if query == "" {
        return err("query is required")
}

    result, apiErr := runAnyquery(ctx, "-q", query)
    if apiErr != nil {
        return err(apiErr.Error())
}

    return ok(result)
}

func HandleListTables(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    result, fetchErr := runAnyquery(ctx, "-q", "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
    if fetchErr != nil {
        return err(fetchErr.Error())
}

    return ok(result)
}

func HandleDescribeTable(ctx context.Context, args map[string.Interface{}) (ToolResponse, error) {
    table, _ :=getString(args, "table")
    if table == "" {
        return err("table is required")
}

    query := fmt.Sprintf("PRAGMA table_info('%s')", strings.ReplaceAll(table, "'", "''"))
    result, execErr := runAnyquery(ctx, "-q", query)
    if execErr != nil {
        return err(execErr.Error())
}

    return ok(result)
}

func HandleListIntegrations(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    result, apiErr := runAnyquery(ctx, "list")
    if apiErr != nil {
        result2, apiErr2 := runAnyquery(ctx, "plugin", "list")
        if apiErr2 != nil {
            return err(apiErr.Error() + "; also tried plugin list: " + apiErr2.Error())
}

        return ok(result2)
}

    return ok(result)
}

func HandleExplainQuery(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
    query, _ :=getString(args, "query")
    if query == "" {
        return err("query is required")
}

    explainQuery := fmt.Sprintf("EXPLAIN QUERY PLAN %s", query)
    result, parseErr := runAnyquery(ctx, "-q", explainQuery)
    if parseErr != nil {
        return err(parseErr.Error())
}

    return ok(result)
}