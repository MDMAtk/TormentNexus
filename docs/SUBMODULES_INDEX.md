# Submodules Index & Project Structure

_Last updated: 2026-07-02, version 1.0.0-alpha.232_

> **All legacy submodules and subprojects have been removed or archived.** 

## Repository Layout

```
tormentnexus/
├── go/                      # Go sidecar (kernel, control plane, tools)
│   ├── cmd/tormentnexus/    # Main binary
│   └── internal/            # Core Go implementations (Memory, MCP, HTTP API, CodeExec)
├── apps/
│   ├── web/                 # Next.js dashboard (port 3000)
│   ├── vscode/              # VS Code extension integration
│   └── tormentnexus-extension/ # Browser context extension
├── packages/                # Monorepo packages (TS types, React components, CLI)
├── data/                    # Database files (.db, assimilated states)
├── archive/                 # Retired submodules or legacy ports (untracked)
│   ├── go_bobbybookmarks/   # Archived bookmarks scraper
│   └── go_marketing_agent/ # Archived enterprise sales bot
└── docs/                    # Architecture, API, and LLM instructions
```
