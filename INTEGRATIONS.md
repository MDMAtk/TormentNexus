# TormentNexus Integrations

Every harness/editor integration for TormentNexus, in one place.
TormentNexus is the canonical repo: integrations live here, and the
HyperNexus repo is only a rename of TormentNexus plus cloud hosting — not
a separate codebase.

Status legend: **present** = code is in this tree; **verified** = exercised
here; **unverified** = not run live. Nothing below is claimed tested unless
it says so.

## Shipped in this tree

| Integration | Location | Status |
|---|---|---|
| VSCode extension | `apps/vscode/` (`tormentnexus-vscode-extension`, v1.0.0-alpha.157) | Present; runtime unverified here |
| OpenCode plugin | `integrations/opencode/` | Folded in 2026-10-24 from the HyperNexus tree; parses + loads under node; live OpenCode run NOT tested |
| MimoCode / MiMo Desktop addon | `integrations/mimo/` | Folded in 2026-10-24 from the HyperNexus tree; `install.js --dry-run` exercised; live install NOT run; `tools/dist/` NOT built (run `build-tools.mjs` first) |
| Chrome extension | `apps/tormentnexus-extension/` (`tormentnexus-extension`, v1.0.0-alpha.157; manifest, pages, packages, e2e) | Present; referenced in repo docs; runtime unverified here |
| borg extension | `apps/borg-extension/` | Scaffold only — template packages (`dev-utils`, `env`, `hmr`, `shared`, `storage`, `tsconfig`, `vite-config`), no manifest or extension code; present but unwired |
| CodeWhale / Pi | Archived in the HyperNexus tree (`archive/crates/foundation-pi/`) | Reference only — NOT copied. It is a Rust crate bridging the `codewhale_foundation_pi::extension` API into the CodeWhale engine, a different engine from TormentNexus's Go/TS tree, so it does not belong here. The OpenCode plugin header claims "parity with Pi extension"; this crate is the surviving Pi artifact. |

## Details

### OpenCode (`integrations/opencode/`)

- `plugins/tormentnexus/index.js` — the plugin (copy to
  `~/.opencode/plugins/tormentnexus/`).
- `tormentnexus-opencode.js` — standalone variant of the same logic
  (`opencode plugin install ./tormentnexus-opencode.js`).
- The two files are near-duplicates upstream (same logic, different
  formatting); both are kept so neither install path breaks.
- 12 MCP tools, 6 slash commands, L2 session priming, `@memory:key`
  expansion, audit logging, status bar widget, RBAC enforcement (per the
  plugin header).
- Talks HTTP to `http://127.0.0.1:7778` by default (`HYPERNEXUS_URL` env
  override). 7778 is a port this tree itself uses for its Go sidecar/API
  (`go/cmd/tormentnexus/mcp_server.go`, `packages/cli` start); whether the
  plugin's `/api/*` routes are served there was not live-tested.
- On activation it writes `~/.opencode/mcp.json` registering a
  `tormentnexus` MCP server entry that spawns `tormentnexus mcp` (stdio).
- Fold-in adaptations: spawned binary `hypernexus` → `tormentnexus`,
  plugin name `hypernexus` → `tormentnexus`, install-comment paths
  corrected. `HYPERNEXUS_URL` and "HyperNexus" description wording kept
  (commercial face; full rename is a separate decision).

### MimoCode / MiMo Desktop (`integrations/mimo/`)

- Skills, 16 desktop tools (source in `tools/src/`, built to `tools/dist/`
  by `build-tools.mjs`), 6 commands (`/hn-store`, `/hn-search`,
  `/hn-status`, `/hn-tools`, `/hn-session`, `/hn-remember`), hooks,
  orchestrator agent, MCP entry merged into `mimocode.jsonc`.
- Install: `node integrations/mimo/install.js [--dry-run]`.
- Fold-in adaptations: comment paths `addons/mimo/` → `integrations/mimo/`;
  `findBinary()` prefers TormentNexus candidates with HyperNexus fallback.
  Skill/tool/command names (`hypernexus_*`, `/hn-*`) and the
  `HYPERNEXUS_WORKSPACE_ROOT` env var kept (commercial face; full rename is
  a separate decision).

## Layout decision (2026-10-24)

Folded integrations live under root `integrations/`, not `apps/`: the root
`pnpm-workspace.yaml` globs `apps/*`, `packages/*`, `mcp-servers/*`, and
these are externally-installed plugin packages, not monorepo build units.
Root-level `integrations/` matches the repo's existing flat top-level
convention (`mcp-servers/`, `cli-tools/`, etc.) and keeps them out of the
workspace build graph.

## NOT built here — Robert's task

Native plugins for **Claude Code, Cursor, Codex, Gemini CLI, Windsurf,
and Copilot** do not exist in this tree. Robert is building those himself
on his machine (he can run full Windsurf etc. there; the sandbox cannot).
No manifests for those six were found anywhere in the TormentNexus or
HyperNexus trees — the website's "native plugins for most harnesses"
phrasing is overstated for those six until his plugins land. (An earlier
audit claim of "only VSCode" was corrected 2026-10-24: VSCode + OpenCode +
MimoCode + Chrome + borg-scaffold is the honest inventory.)

## Open questions for Robert

- Merge `muse/fold-in-2026-10-24` (branch-only; nothing pushed or merged
  by this fold).
- The HyperNexus rename + cloud-hosting step (separate from this fold).
- Full branding rename (`hypernexus_*` tool/skill names,
  `HYPERNEXUS_URL` / `HYPERNEXUS_WORKSPACE_ROOT` env vars) vs keeping
  HyperNexus as the commercial name on these surfaces.
- Whether the OpenCode/Mimo integrations' live behavior gets proven
  against a running kernel before any customer pilot uses them.
- borg-extension: finish it, or drop the scaffold.
