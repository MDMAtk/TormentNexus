package tools
import ("context"; "fmt"; "os/exec")
func HandleCodemod(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	command, _ := getString(args, "command"); if command == "" { return err("command is required") }
	cmd := exec.CommandContext(ctx, "codemod", command)
	output, e := cmd.CombinedOutput()
	if e != nil { return err(fmt.Sprintf("codemod failed: %v", e)) }
	return ok(string(output))
}
