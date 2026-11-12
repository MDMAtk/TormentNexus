# TormentNexus

An AI control plane: a local-first system for coordinating multiple AI coding
agents, MCP (Model Context Protocol) tools, provider routing, session
continuity, and operator observability — one place to run tools, models,
sessions, and subagents instead of juggling them by hand.

TormentNexus is the canonical, open-source codebase. **HyperNexus is not a
separate project** — it is the rename plus cloud hosting / commercial face of
this same codebase (see HyperNexusLLC/HyperNexus).

## Status

**In active development, not production-proven.** The repo is large and
fast-moving, with a lot of scaffolding alongside working parts. Performance
and routing claims have not been validated by a spec benchmark — the latest
exploratory runs show parity-level results with average latency still in
minutes, not the near-instant routing the architecture targets. Do not take
tool/agent counts in older docs at face value; verify against the code.

Last commit: **2026-11-05** (submodule re-pin). Before that, a batch of
merges on 2026-10-24 (CI fixes, integrations fold-in, scaffold removal).

## What it is

- **MCP registry and catalog** — tracks a large set of MCP servers and tools
  in a local SQLite catalog, with native Go tool implementations for many
  common integrations (filesystem, Slack, SQLite, search, Ollama, and more).
- **Provider routing** — routes requests across model providers, with a fast
  router plus background queueing (Architecture A). Router model pinned to
  `qwen3:1.7b` (think off) as of 2026-10-31.
- **Multi-agent orchestration** — swarm launcher, task plane, session
  continuity, and operator dashboard.
- **Memory & context** — persistent memory bank, context portal, LanceDB
  vector store.
- **Go sidecar + TypeScript core** — Go for performance-sensitive tooling,
  TypeScript/Next.js for the dashboard and orchestration layer.

## What's working

- The repo builds and runs locally (Go sidecar + TS core + dashboard).
- MCP catalog and native Go tools exist and are exercised in local runs.
- AGPL-3.0 license adopted 2026-10-23 (open core; commercial licensing via
  HyperNexus).

## What's not / known gaps

- **Live proof is not closed.** No spec benchmark has validated the routing
  or latency claims. Re-spec pins (queueing, latency measurement, predictor
  commit, real tokenizer, deployment machine) are still required before a
  run counts as proof.
- Harness plugins for the six agent CLIs (Claude, Cursor, Codex, Gemini,
  Windsurf, Copilot) are built by the author on his own machine — they are
  not in this repo.
- The repo carries a lot of dead weight: stale docs, duplicate session
  archives, `.db` files, and abandoned scaffolding. Cleanup is ongoing.
- No paying customers. Outreach to prospective design partners is held until
  live proof and billing verification exist.

## Layout (top level)

```
TormentNexus/
├── main.go, go/            Go sidecar — native tools, kernel
├── packages/, apps/        TypeScript core and dashboard
├── mcp-servers/, mcp/      MCP server implementations and registry
├── orchestrator/           Multi-agent orchestration
├── memory-bank/            Persistent memory
├── swarm_launcher.py       Swarm / multi-agent runner
├── docs/                   Architecture and design docs
└── landing/                Marketing site source
```

## Getting started

See `DEPLOY.md` and `docs/` for setup. In short: Go toolchain + Node.js,
`start.bat` / `start.sh` launchers for the Go sidecar and TS core. Expect
rough edges — this is a working codebase, not a polished release.

## License

AGPL-3.0 (see `LICENSE`). Network use triggers the copyleft clause;
commercial licensing is available through HyperNexus.

## Links

- Commercial face / cloud hosting: https://github.com/HyperNexusLLC/HyperNexus
- Sites: https://tormentnexus.site — https://hypernexus.site
- Author: https://robertpelloni.com
