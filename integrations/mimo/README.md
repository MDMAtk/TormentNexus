# HyperNexus MiMo Addons

Full-featured HyperNexus integration for **MiMoCode CLI** and **MiMo Desktop**.

## Contents

| Component | Target | Files |
|---|---|---|
| Skill | MiMoCode + MiMo Desktop + Agents | `skills/hypernexus/` (SKILL.md + locales) |
| Desktop tools (16) | MiMo Desktop `engine-config/tools/` | `tools/dist/*.js` |
| Commands (6) | MiMoCode `commands/` | `mimocode/commands/*.md` |
| Hooks | MiMoCode `hooks/` | `mimocode/hooks/config.json` |
| Agent | MiMoCode `agents/` | `mimocode/agents/hypernexus-orchestrator.md` |
| MCP config | MiMoCode `mimocode.jsonc` | merged at install time |

## Desktop Tools

`hypernexus_search_memory`, `hypernexus_memory_vector_search`, `hypernexus_memory_list`,
`hypernexus_add_fact`, `hypernexus_record_observation`, `hypernexus_knowledge_graph`,
`hypernexus_context_harvest`, `hypernexus_session_search`, `hypernexus_session_context`,
`hypernexus_mcp_search`, `hypernexus_tool_discover`, `hypernexus_mcp_call`,
`hypernexus_code_search`, `hypernexus_repomap`, `hypernexus_system_status`,
`hypernexus_billing_status`

## MiMoCode Commands

`/hn-store`, `/hn-search`, `/hn-status`, `/hn-tools`, `/hn-session`, `/hn-remember`

## Build & Install

```bash
node integrations/mimo/generate.mjs    # regenerate sources (optional)
node integrations/mimo/build-tools.mjs # bundle tools with zod (esbuild)
node integrations/mimo/install.js      # install to MiMoCode + MiMo Desktop + Agents
```

## Requirements

- TormentNexus binary (`tormentnexus mcp`) — preferred; the installer falls
  back to a `hypernexus` binary if no TormentNexus binary is found
- Node 20+
- zod (resolved from `~/.config/mimocode/node_modules/zod` or the repo root)
- esbuild (resolved from the repo root)

## Fold-in (2026-10-24, from the HyperNexus tree)

This addon was copied verbatim from the HyperNexus tree's `addons/mimo/`
into `integrations/mimo/` as part of folding all integrations into
TormentNexus (the canonical repo). Adaptations, and only these:

- Usage/comment paths `addons/mimo/` → `integrations/mimo/` (comments
  only; the installer resolves everything relative to its own directory).
- `install.js findBinary()` now prefers TormentNexus candidates
  (`~/.tormentnexus/tormentnexus`, `/usr/local/bin/tormentnexus`, PATH
  `tormentnexus`) and falls back to the HyperNexus candidates, so the MCP
  entry spawns a binary this tree ships. The entry still runs
  `<binary> mcp`; TormentNexus has an `mcp` subcommand.
- Kept unchanged: skill/tool/command names (`hypernexus_*`, `/hn-*`,
  `skills/hypernexus`), the `HYPERNEXUS_WORKSPACE_ROOT` env var, and
  "HyperNexus" wording (HyperNexus remains the commercial face; a full
  branding rename is a separate product decision).

Status: UNVERIFIED at runtime. All JS/MJS files pass `node --check` and
`install.js --dry-run` was exercised from this location (it reports planned
copies without writing). A live install was NOT run (it writes to the user's
home directory), `tools/dist/` is NOT built here (run `build-tools.mjs`
first — it needs zod), and no live MimoCode/Desktop session was tested.
