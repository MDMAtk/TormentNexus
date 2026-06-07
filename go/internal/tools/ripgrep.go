package tools
import ("context"; "fmt"; "os/exec")
func HandleRipgrep(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	pattern, _ := getString(args, "pattern"); if pattern == "" { return err("pattern is required") }
	path, _ := getString(args, "path"); if path == "" { path = "." }
	cmd := exec.CommandContext(ctx, "rg", "--json", pattern, path)
	output, e := cmd.CombinedOutput()
	if e != nil { if exitErr, okVal := e.(*exec.ExitError); okVal && exitErr.ExitCode() == 1 { return ok("[]") }; return err(fmt.Sprintf("ripgrep failed: %v", e)) }
	return ok(string(output))
}
