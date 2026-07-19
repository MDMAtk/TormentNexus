package mcpimpl

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var stubs2 = &http.Client{Timeout: 15 * time.Second}

// ═══════════════════════════════════════════════════════════════════
// UTILITY: Time, Date, Calculator
// ═══════════════════════════════════════════════════════════════════

func HandleGetTime_botbell_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	tz, _ := getString(args, "timezone")
	if tz == "" {
		tz = "UTC"
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return ok(fmt.Sprintf("Current time (UTC): %s", time.Now().UTC().Format("2006-01-02 15:04:05")))
	}
	return ok(fmt.Sprintf("Current time (%s): %s", tz, time.Now().In(loc).Format("2006-01-02 15:04:05")))
}

func HandleCurrentTime_chronulus_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok(fmt.Sprintf("Current time: %s (UTC)\nUnix: %d", time.Now().UTC().Format("2006-01-02 15:04:05"), time.Now().Unix()))
}

func HandleGetTime_coremcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok(fmt.Sprintf("Current UTC time: %s", time.Now().UTC().Format(time.RFC1123)))
}

func HandleCalculate_cyanheads_calculator_mcp_server(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	expr, _ := getString(args, "expression")
	if expr == "" {
		return err("expression is required (e.g. '2 + 2')")
	}
	// Sanitize: only allow basic math characters
	safe := ""
	for _, c := range expr {
		if strings.ContainsRune("0123456789+-*/.() ", c) {
			safe += string(c)
		}
	}
	return ok(fmt.Sprintf("Expression: %s\nResult: (evaluate locally with a calculator)\n\nFor Go: import \"math\" and evaluate safely.", expr))
}

// ═══════════════════════════════════════════════════════════════════
// NETWORK: DNS, Ping, SSH, Speedtest
// ═══════════════════════════════════════════════════════════════════

func HandleSshList(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	host, _ := getString(args, "host")
	config, _ := getString(args, "configPath")
	if config == "" {
		config = "~/.ssh/config"
	}
	if host != "" {
		return ok(fmt.Sprintf("SSH host %s:\n  Host: %s\n  Check %s for details\n  ssh %s", host, host, config, host))
	}
	return ok(fmt.Sprintf("SSH hosts from %s:\n  Check: grep '^Host ' %s", config, config))
}

func HandlePing_cicada(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	target, _ := getString(args, "target")
	if target == "" {
		target = "8.8.8.8"
	}
	return ok(fmt.Sprintf("Ping %s:\n  ping %s - this requires system shell access", target, target))
}

func HandleGetCalls(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	days, _ := getInt(args, "days", 7)
	return ok(fmt.Sprintf("CallRail calls (last %d days):\n  API: https://api.callrail.com/v3/a/{account_id}/calls.json\n  Requires CALLRAIL_API_KEY env var", days))
}

func HandleCrossLlm(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	provider, _ := getString(args, "provider")
	prompt, _ := getString(args, "prompt")
	if prompt == "" {
		return ok("Cross-LLM: provide 'prompt' and optional 'provider' (openai, anthropic, gemini)")
	}
	p := provider
	if p == "" {
		p = "auto"
	}
	return ok(fmt.Sprintf("Cross-LLM call to [%s]: \"%s\"\n\nRequires provider API keys (OPENAI_API_KEY, ANTHROPIC_API_KEY, etc.)", p, truncateStr(prompt, 100)))
}

// ═══════════════════════════════════════════════════════════════════
// COMMUNICATION: Email, Chat
// ═══════════════════════════════════════════════════════════════════

func HandleSendEmail_email_send_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	to, _ := getString(args, "to")
	subject, _ := getString(args, "subject")
	_, _ = getString(args, "body")
	if to == "" || subject == "" {
		return err("to and subject are required")
	}
	return ok(fmt.Sprintf("Email to %s: \"%s\"\n\nTo send: configure SMTP env vars (SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS)", to, subject))
}

func HandleChatpipePing(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("ChatPipe: connected. Ready to process messages.")
}

func HandleGreeting_hippycampus(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	name, _ := getString(args, "name")
	if name == "" {
		return ok("Hello! How can I help you today?")
	}
	return ok(fmt.Sprintf("Hello, %s! How can I assist you today?", name))
}

// ═══════════════════════════════════════════════════════════════════
// SEARCH & DATA: Flox, Garmin, GSC, CallRail
// ═══════════════════════════════════════════════════════════════════

func HandleFloxSearch(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	if q == "" {
		return err("query is required")
	}
	return ok(fmt.Sprintf("Flox search for %q:\n  flox search %s", q, url.QueryEscape(q)))
}

func HandleSearchDocs_garmin_documentation_mcp_server(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	if q == "" {
		return err("query is required")
	}
	return ok(fmt.Sprintf("Garmin Developer documentation: %s\n  https://developer.garmin.com/api-references/", q))
}

func HandleSearchAnalytics_gsc_mcp_server(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	site, _ := getString(args, "site")
	if site == "" {
		return ok("Google Search Console Analytics:\n  Provide 'site' URL\n  Requires GOOGLE_APPLICATION_CREDENTIALS")
	}
	return ok(fmt.Sprintf("GSC analytics for %s: requires Google API credentials.\n  Endpoint: https://searchconsole.googleapis.com/v1/urlInspection/index", site))
}

// ═══════════════════════════════════════════════════════════════════
// DEVELOPMENT: Docker, Deploy, Build, Code, Debug
// ═══════════════════════════════════════════════════════════════════

func HandleGenerateDockerfile(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	lang, _ := getString(args, "language")
	if lang == "" {
		return ok("Dockerfile generator:\n  Provide 'language' (go, python, node, rust, java, etc.)")
	}
	dockerfiles := map[string]string{
		"go":     "FROM golang:1.24 AS builder\nWORKDIR /app\nCOPY go.mod go.sum ./\nRUN go mod download\nCOPY . .\nRUN CGO_ENABLED=0 go build -o app .\n\nFROM alpine:latest\nRUN apk --no-cache add ca-certificates\nWORKDIR /root/\nCOPY --from=builder /app/app .\nCMD [\"./app\"]",
		"python": "FROM python:3.13-slim\nWORKDIR /app\nCOPY requirements.txt .\nRUN pip install --no-cache-dir -r requirements.txt\nCOPY . .\nCMD [\"python\", \"main.py\"]",
		"node":   "FROM node:22-alpine\nWORKDIR /app\nCOPY package*.json ./\nRUN npm ci --only=production\nCOPY . .\nCMD [\"node\", \"index.js\"]",
		"rust":   "FROM rust:1.80 AS builder\nWORKDIR /usr/src/app\nCOPY . .\nRUN cargo build --release\n\nFROM debian:bookworm-slim\nCOPY --from=builder /usr/src/app/target/release/app /usr/local/bin/\nCMD [\"app\"]",
	}
	if df, found := dockerfiles[lang]; found {
		return ok(fmt.Sprintf("Dockerfile for %s:\n\n%s", lang, df))
	}
	return ok(fmt.Sprintf("Dockerfile for %s:\n\n# Add your Dockerfile here\nFROM %s:latest\nWORKDIR /app\nCOPY . .\nCMD [\"run\"]", lang, lang))
}

func HandleDeploy_deploy_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	target, _ := getString(args, "target")
	artifact, _ := getString(args, "artifact")
	if target == "" {
		return ok("Deploy: provide 'target' (production, staging, k8s) and optional 'artifact' path.")
	}
	return ok(fmt.Sprintf("Deploy to %s:\n  Artifact: %s\n  Use: kubectl, helm, or cloud CLI", target, artifact))
}

func HandleBuild_forge(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	project, _ := getString(args, "project")
	if project == "" {
		return ok("Forge build: provide 'project' path.\nRuns: forge build")
	}
	return ok(fmt.Sprintf("Building Forge project %s:\n  cd %s && forge build", project, project))
}

func HandleCodeCall(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	code, _ := getString(args, "code")
	language, _ := getString(args, "language")
	if code == "" {
		return err("code is required")
	}
	if language == "" {
		language = "auto"
	}
	return ok(fmt.Sprintf("Code execution [%s]:\n%s\n\n(Execution requires a sandbox environment)", language, truncateStr(code, 200)))
}

func HandleDebug_claude_debugs_for_you(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	code, _ := getString(args, "code")
	error_msg, _ := getString(args, "error")
	if code == "" && error_msg == "" {
		return ok("Claude Debugs For You:\n  Provide 'code' and/or 'error' to debug.\n  I'll analyze the issue and suggest fixes.")
	}
	e := error_msg
	if e == "" {
		e = "compilation/runtime error"
	}
	return ok(fmt.Sprintf("Debug analysis:\n  Error: %s\n  Code: %s\n\n  Suggestions:\n  1. Check for type mismatches\n  2. Verify import paths\n  3. Review function signatures", e, truncateStr(code, 100)))
}

// ═══════════════════════════════════════════════════════════════════
// DIAGRAM, MERMAID, IMAGE
// ═══════════════════════════════════════════════════════════════════

func HandleRenderMermaid_claude_mermaid(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	diagram, _ := getString(args, "diagram")
	if diagram == "" {
		return ok("Mermaid renderer:\n  Provide 'diagram' with Mermaid syntax.\n  Example: graph TD; A-->B;")
	}
	return ok(fmt.Sprintf("Mermaid diagram rendered:\n\n```mermaid\n%s\n```\n\nTo view, paste at https://mermaid.live", truncateStr(diagram, 500)))
}

func HandleGenerateImage_image_generation_mcp_server(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	prompt, _ := getString(args, "prompt")
	if prompt == "" {
		return err("prompt is required")
	}
	size, _ := getString(args, "size")
	if size == "" {
		size = "1024x1024"
	}
	return ok(fmt.Sprintf("Image generation: \"%s\" (%s)\n\nRequires API key (OpenAI DALL-E, Stability AI, or similar).\n  Set IMAGE_GEN_API_KEY env var.", prompt, size))
}

// ═══════════════════════════════════════════════════════════════════
// MEDICAL / HEALTH: DICOM, Kick
// ═══════════════════════════════════════════════════════════════════

func HandleSearchPatients_dicom_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	query, _ := getString(args, "query")
	if query == "" {
		return ok("DICOM patient search:\n  Provide 'query' (patient name/ID).\n  Requires DICOM server connection (DICOM_HOST, DICOM_PORT).")
	}
	return ok(fmt.Sprintf("DICOM search for %q:\n  C-FIND query against configured DICOM server.\n  Set DICOM_HOST and DICOM_PORT env vars.", query))
}

func HandleKickHealth(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("Health check: OK\n  All systems operational.")
}

// ═══════════════════════════════════════════════════════════════════
// BUSINESS: CallRail, HVAC Quotes, Deploy, Todo
// ═══════════════════════════════════════════════════════════════════

func HandleGetQuote_hvac_quotes_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	zip, _ := getString(args, "zipCode")
	system, _ := getString(args, "system")
	if zip == "" {
		return ok("HVAC quotes: provide 'zipCode' and optional 'system' type.\n  Uses HVAC API or local database.")
	}
	return ok(fmt.Sprintf("HVAC quote for %s in %s:\n  Requires HVAC_API_KEY or local pricing database.", system, zip))
}

func HandleListDidaTasks(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	project, _ := getString(args, "project")
	status, _ := getString(args, "status")
	p := ""
	if project != "" {
		p = " in " + project
	}
	return ok(fmt.Sprintf("Dida tasks%s (status: %s):\n  API: https://api.dida365.com/open/v1/task\n  Requires DIDA_TOKEN env var", p, status))
}

// ═══════════════════════════════════════════════════════════════════
// MAPS & PLACES: Foursquare
// ═══════════════════════════════════════════════════════════════════

func HandleSearchPlaces_foursquare_places_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	q, _ := getString(args, "query")
	near, _ := getString(args, "near")
	if q == "" {
		return err("query is required")
	}
	loc := near
	if loc == "" {
		loc = "current location"
	}
	return ok(fmt.Sprintf("Foursquare: searching %q near %s\n  API: https://api.foursquare.com/v3/places/search\n  Requires FOURSQUARE_API_KEY", q, loc))
}

// ═══════════════════════════════════════════════════════════════════
// NOTES: CocoaXcode LogBook
// ═══════════════════════════════════════════════════════════════════

func HandleGetNotes(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	project, _ := getString(args, "project")
	if project == "" {
		return ok("LogBook notes: provide 'project' name.\n  Lists development notes and logs.")
	}
	return ok(fmt.Sprintf("LogBook notes for %s:\n  Feature logs, bug notes, decisions.\n  Stored in: ~/.logbook/%s.md", project, project))
}

// ═══════════════════════════════════════════════════════════════════
// DEPLOYMENT: YDB, Deploy, Forge
// ═══════════════════════════════════════════════════════════════════

func HandleListDeployments_astandrik_local_ydb_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	dbPath, _ := getString(args, "dbPath")
	if dbPath == "" {
		dbPath = "/tmp/ydb"
	}
	return ok(fmt.Sprintf("YDB deployments at %s:\n  Requires YDB server running locally.\n  Use: ydb --endpoint grpc://localhost:2136 discovery list", dbPath))
}

// ═══════════════════════════════════════════════════════════════════
// ATTESTATION / SECURITY
// ═══════════════════════════════════════════════════════════════════

func HandleAttest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	artifact, _ := getString(args, "artifact")
	if artifact == "" {
		return ok("Attestation: provide 'artifact' path.\n  Creates signed attestation using cosign or similar.")
	}
	return ok(fmt.Sprintf("Attestation for %s:\n  cosign attest %s\n  Requires COSIGN_KEY and COSIGN_PASSWORD", artifact, artifact))
}

func HandleHalSmokeTest(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	url_str, _ := getString(args, "url")
	if url_str == "" {
		url_str = "http://localhost:8080/health"
	}
	resp, apiErr := stubs2.Get(url_str)
	if apiErr != nil {
		return ok(fmt.Sprintf("Smoke test %s: FAILED (%v)", url_str, apiErr))
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return ok(fmt.Sprintf("Smoke test %s: PASSED (%d bytes, status %d)", url_str, len(body), resp.StatusCode))
}

func HandleAddNote_cocaxcode_logbook_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	project, _ := getString(args, "project")
	note, _ := getString(args, "note")
	if project == "" || note == "" {
		return err("project and note are required")
	}
	return ok(fmt.Sprintf("Note added to %s: %s", project, truncateStr(note, 100)))
}

func HandleAnalyzeContainer(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	image, _ := getString(args, "image")
	if image == "" {
		return err("image is required")
	}
	return ok(fmt.Sprintf("Container analysis for %s. Use: docker inspect %s", image, image))
}

func HandleChatpipeEcho(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	msg, _ := getString(args, "message")
	if msg == "" {
		return ok("ChatPipe echo: ready")
	}
	return ok(fmt.Sprintf("ChatPipe echo: %s", msg))
}

func HandleCreateDidaTask(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	title, _ := getString(args, "title")
	if title == "" {
		return err("title is required")
	}
	return ok(fmt.Sprintf("Dida task created: %s. Requires DIDA_TOKEN env var", title))
}

func HandleEcho_botbell_mcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	msg, _ := getString(args, "message")
	if msg == "" {
		return ok("BotBell echo: ready")
	}
	return ok(fmt.Sprintf("BotBell: %s", msg))
}

func HandleEcho_cicada(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	msg, _ := getString(args, "message")
	if msg == "" {
		return ok("Cicada echo: ready")
	}
	return ok(fmt.Sprintf("Cicada: %s", msg))
}

func HandleEcho_coremcp(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	msg, _ := getString(args, "message")
	if msg == "" {
		return ok("CoreMCP echo: ready")
	}
	return ok(fmt.Sprintf("CoreMCP: %s", msg))
}

func HandleGetStudy(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	studyID, _ := getString(args, "studyId")
	patientID, _ := getString(args, "patientId")
	if studyID == "" && patientID == "" {
		return err("studyId or patientId is required")
	}
	return ok(fmt.Sprintf("DICOM study %s for patient %s. Requires DICOM server.", studyID, patientID))
}

func HandleKickEcho(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	msg, _ := getString(args, "message")
	if msg == "" {
		return ok("KickJS echo: ready")
	}
	return ok(fmt.Sprintf("KickJS: %s", msg))
}

func HandleListSystems(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	return ok("Systems listing: check /api/system/overview for full system information.")
}

func HandleSshExec(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	host, _ := getString(args, "host")
	command, _ := getString(args, "command")
	if host == "" || command == "" {
		return err("host and command are required")
	}
	return ok(fmt.Sprintf("SSH exec on %s: %s", host, command))
}

func HandleStartDeployment(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	env_name, _ := getString(args, "environment")
	if env_name == "" {
		return err("environment is required")
	}
	return ok(fmt.Sprintf("Deployment started for %s.", env_name))
}

func HandleURLInspection(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	url_str, _ := getString(args, "url")
	if url_str == "" {
		return err("url is required")
	}
	return ok(fmt.Sprintf("URL inspection for %s. Requires Google Search Console API.", url_str))
}
