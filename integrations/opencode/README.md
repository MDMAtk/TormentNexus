# OpenCode Integration

TormentNexus bridge for [OpenCode](https://opencode.ai): persistent memory,
MCP tools, sessions, skills, and commercial RBAC, talking to the TormentNexus
kernel over HTTP (`http://127.0.0.1:7778` by default — 7778 is a port the
TormentNexus tree itself uses for its Go sidecar/API
(`go/cmd/tormentnexus/mcp_server.go`, `packages/cli/src/commands/start.ts`),
so the default needed no change on fold-in; whether the plugin's
`/api/*` routes are served on that port was not live-tested here).

## Files

- `plugins/tormentnexus/index.js` — the OpenCode plugin. Install by copying
  this directory to `~/.opencode/plugins/tormentnexus/` (so OpenCode loads
  `~/.opencode/plugins/tormentnexus/index.js`).
- `tormentnexus-opencode.js` — standalone variant of the same plugin logic
  (same tools, commands, and HTTP bridge; formatted differently). Install
  with `opencode plugin install ./tormentnexus-opencode.js`.

The two files are near-duplicates upstream (same logic, different
formatting); both are kept here deliberately so neither install path breaks.
If they drift, `plugins/tormentnexus/index.js` is the canonical one.

## What it provides (per the plugin header)

- 12 custom MCP tools (memory, tools, sessions, skills, code, context,
  scratchpad, processes, billing)
- Session priming with L2 context injection; per-turn context harvesting
- `@memory:key` inline expansion in prompts
- Tool-call audit logging; auto-storage of tool results to L2
- 6 slash commands; live status bar widget; subagent orchestration dispatcher
- Commercial RBAC enforcement

On activation the plugin also writes `~/.opencode/mcp.json`, registering a
`tormentnexus` MCP server entry that spawns `tormentnexus mcp` (stdio).

## Fold-in adaptations (2026-10-24, from the HyperNexus tree)

- Spawned binary `hypernexus`/`hypernexus.exe` → `tormentnexus`/
  `tormentnexus.exe` (the TormentNexus tree ships `tormentnexus`, which has
  an `mcp` subcommand).
- Plugin `module.exports.name` `hypernexus` → `tormentnexus`; MCP server
  key in the written `mcp.json` stays `hypernexus`.
- Header install comments corrected: `./tormentnexus-opencode.js` and
  `~/.opencode/plugins/tormentnexus/index.js` (upstream comments named
  files/paths that did not exist).
- Kept unchanged: `HYPERNEXUS_URL` env var, port `7778`, and
  "HyperNexus" wording in descriptions (HyperNexus remains the commercial
  face; a full branding rename is a separate product decision).

## Status

UNVERIFIED at runtime: both files pass `node --check` and the plugin loads
under node (exports tools + commands), but a live run against OpenCode and a
running TormentNexus kernel was not performed here. The plugin header claims
"parity with Pi extension" — see `INTEGRATIONS.md` for the archived Pi
artifact.
