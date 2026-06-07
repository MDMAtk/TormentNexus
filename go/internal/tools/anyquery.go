package tools
import ("context"; "fmt"; "os/exec")
func HandleAnyquery(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query", "sql"); if query == "" { return err("query is required") }
	cmd := exec.CommandContext(ctx, "anyquery", "-q", query, "--json")
	output, e := cmd.CombinedOutput()
	if e != nil { return err(fmt.Sprintf("anyquery failed: %v", e)) }
	return ok(string(output))
}
