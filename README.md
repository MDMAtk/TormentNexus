🚀 TormentNexus — Universal AI Control Plane for Multi-Agent Workflows, MCP Tools & Context-Aware Memory

```text
╔══════════════════════════════════════════════════════════════════════════════╗
║                                                                              ║
║                     ██╗   ██╗███╗   ██╗██████╗ ███████╗██████╗              ║
║                     ██║   ██║████╗  ██║██╔══██╗██╔════╝██╔══██╗             ║
║                     ██║   ██║██╔██╗ ██║██║  ██║█████╗  ██████╔╝             ║
║                     ██║   ██║██║╚██╗██║██║  ██║██╔══╝  ██╔══██╗             ║
║                     ╚██████╔╝██║ ╚████║██████╔╝███████╗██║  ██║             ║
║                      ╚═════╝ ╚═╝  ╚═══╝╚═════╝ ╚══════╝╚═╝  ╚═╝             ║
║                                                                              ║
║                     ██████╗ ██████╗ ███╗   ██╗███████╗████████╗██████╗      ║
║                    ██╔════╝██╔═══██╗████╗  ██║██╔════╝╚══██╔══╝██╔══██╗     ║
║                    ██║     ██║   ██║██╔██╗ ██║███████╗   ██║   ██████╔╝     ║
║                    ██║     ██║   ██║██║╚██╗██║╚════██║   ██║   ██╔══██╗     ║
║                    ╚██████╗╚██████╔╝██║ ╚████║███████║   ██║   ██║  ██║     ║
║                     ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚═╝  ╚═╝     ║
║                                                                              ║
║                     █████╗ ██╗     ██████╗ ██╗  ██╗ █████╗                  ║
║                    ██╔══██╗██║     ██╔══██╗██║  ██║██╔══██╗                 ║
║                    ███████║██║     ██████╔╝███████║███████║                 ║
║                    ██╔══██║██║     ██╔═══╝ ██╔══██║██╔══██║                 ║
║                    ██║  ██║███████╗██║     ██║  ██║██║  ██║                 ║
║                    ╚═╝  ╚═╝╚══════╝╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝                 ║
║                                                                              ║
║                    ╔══════════════════════════════════════╗                  ║
║                    ║     🧠  PRODUCTION READY  🧠           ║                  ║
║                    ║  PERSISTENT MEMORY & 20,000+ TOOLS     ║                  ║
║                    ║  WORKS WITH 38+ AI AGENTS       ║                  ║
║                    ╚══════════════════════════════════════╝                  ║
║                                                                              ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

# TormentNexus: The Cognitive Kernel — Universal AI Control Plane for Multi-Agent Workflows, MCP Tools & Context-Aware Memory

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript)
![Next.js](https://img.shields.io/badge/Next.js-16-000000?logo=next.js)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)
![License](https://img.shields.io/badge/license-Enterprise%2FOSS-orange)

> **TormentNexus** is the ultimate local-first control plane for multi-agent workflows, Model Context Protocol (MCP) tooling, provider routing, session continuity, and operator observability. It is the substrate where a single local system seamlessly coordinates tools, models, sessions, context, subagents, and full visibility across the entire AI-driven software development stack.

---

## Table of Contents

- [What It Does](#what-it-does)
- [Architecture](#architecture)
- [Core Pillars](#core-pillars)
- [Monorepo Structure](#monorepo-structure)
- [The Go Sidecar](#the-go-sidecar)
- [The Dashboard](#the-dashboard)
- [MCP Ecosystem](#mcp-ecosystem)
- [Memory & Context](#memory--context)
- [Swarm & Multi-Agent](#swarm--multi-agent)
- [API Surface](#api-surface)
- [Quick Start](#quick-start)
- [What's Planned](#whats-planned)
- [Roadmap](#roadmap)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

---

## What It Does

TormentNexus is a **decision system and universal bridge** — not just an aggregator. It runs locally as a modular monolith that unifies the chaotic landscape of AI tools, models, and agents into a single, coherent operating system for AI-driven development.

### Current Capabilities (v1.0.0-alpha.132)

| Capability | Status | Details |
|------------|--------|---------|
| **MCP Registry** | Stable | 14,250+ tracked MCP servers, 11,024+ populated in SQLite catalog, 600+ verified servers, 11,000+ verified tools |
| **Native Go Tools** | Beta | 3,900+ native Go tool implementations replacing external MCP servers (filesystem, Slack, SQLite, DuckDuckGo, Ollama, TTS, Vercel, NWS, DexPaprika, Firecrawl, Exa, arXiv, Semantic Scholar, Mem0, Alpaca, Alpha Vantage, Hugging Face, Semgrep, Octagon, Browser Automation, ChromaDB, Basic Memory, MindsDB, Serena, AST-grep, PAL, Thoughtbox, and more) |
| **Progressive Tool Routing** | Stable | Semantic vector search + BM25 ranking injects only the most relevant tools into LLM context windows |
| **Dual-Tier Memory** | Stable | L1 (session scratchpad) + L2 (semantic SQLite vault) with heat-score lifecycle and autonomous context harvesting |
| **LLM Waterfall** | Stable | Cascading failover: NVIDIA NIM → OpenRouter → Local LM Studio / Ollama with 429/5xx handling |
| **Multi-Agent Swarm** | Beta | A2A protocol coordination, role rotation (Planner→Implementer→Tester→Critic), consensus engine |
| **Autonomous Healer** | Stable | `Diagnose → Fix → Verify → Retry` loop with native code execution (tsc, vitest, go test) and L2 vault persistence |
| **Browser Automation** | Beta | Native chromedp handlers: navigate, screenshot, evaluate, click, fill, get HTML |
| **Skill Registry** | Stable | 3,229+ assimilated skills from 7 harness ecosystems (Aider, Agent, CCS, Hermes, Pi, etc.) with Jaccard deduplication |
| **Dashboard** | Stable | Next.js 16 + React 19 + Tailwind CSS 4 with real-time telemetry, knowledge graph, healer view, swarm visualizer |
| **tRPC Bridge** | Stable | Type-safe API layer (port 4100) connecting UI to Go sidecar |
| **Session Import** | Beta | Automatic ingestion of Claude, Aider, and other harness session artifacts |
| **Enterprise Licensing** | Experimental | Ed25519-signed license token validation with offline verification |
| **Supervisor Nudge** | Stable | Autonomous Windows UI automation to maintain development momentum across AI chat surfaces |
| **Deep Link Protocol** | Beta | `tormentnexus://attach?session=ID` and `tormentnexus://create?cliType=aider` URI handling |
| **Provider Metrics** | Stable | Real-time telemetry tracking for all LLM providers with latency and cost analysis |

---

## Architecture

TormentNexus is a **high-performance Go modular monolith** with a **TypeScript/Next.js frontend**, operating as a local-first control plane.

```
┌─────────────────────────────────────────────────────────────┐
│  OPERATOR LAYER                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │  Web Dash   │  │  CLI (TS)   │  │  VS Code    │       │
│  │  Port 3000  │  │  tormentnexus│  │  Extension  │       │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘       │
│         │                │                │                 │
│         └────────────────┴────────────────┘                 │
│                          │                                │
│  ┌───────────────────────┴───────────────────────┐        │
│  │  TYPESCRIPT CONTROL PLANE (Port 4100)        │        │
│  │  tRPC routers · NativeSidecarDaemon · ResilientStream  │
│  └───────────────────────┬───────────────────────┘        │
│                          │                                │
│  ┌───────────────────────┴───────────────────────┐        │
│  │  GO SIDECAR (Port 4300) — The Authoritative Kernel    │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐     │
│  │  │ SkillStore│ │ EventBus │ │  Vault  │ │ Healer  │     │
│  │  │ (BM25)  │ │ (Swarm) │ │(sqlite) │ │(Immune) │     │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘     │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐     │
│  │  │  Router  │ │ PairOrchestrator│ │ CodeExecutor│ │ MCP Sync │     │
│  │  │(Progressive)│ │(Consensus) │ │(Sandbox) │ │(Registry)│     │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘     │
│  └──────────────────────────────────────────────────────┘
│                          │                                │
│  ┌───────────────────────┴───────────────────────┐        │
│  │  EXTERNAL MODELS & TOOLS                     │        │
│  │  OpenAI · Anthropic · Gemini · OpenRouter · Ollama  │
│  │  600+ MCP Servers · 3,900+ Native Go Tools          │
│  └──────────────────────────────────────────────────────┘
└─────────────────────────────────────────────────────────────┘
```

### Key Ports

| Service | Port | Purpose |
|---------|------|---------|
| Next.js Dashboard | 3000 | Web observation deck |
| tRPC Bridge | 4100 | TypeScript Control Plane API |
| Go Sidecar | 4300 | Authoritative native kernel |

---

## Core Pillars

### 1. Progressive MCP Tool Routing & Parity

Models should never be overwhelmed with a 50,000-token tool dump. TormentNexus employs a multi-layered, progressive disclosure system:

- **Semantic Search**: Local vector embeddings match the active prompt against a global MCP directory of 14,250+ servers.
- **The Router**: Only the top highly relevant tool schemas are injected into the active LLM context.
- **Universal Parity**: Byte-for-byte identical tool signatures for Claude Code, Codex, Gemini CLI, Cursor, and Windsurf.

### 2. Dual-Tier Memory Architecture (L1 / L2)

Context is finite; memory must be infinite.

- **L1 — Session Scratchpad**: Ephemeral, lightning-fast memory tied directly to the active session.
- **L2 — The Vault**: Permanent semantic storage in SQLite with `sqlite-vec`. Saves exact transcripts and LLM-compressed heuristics.
- **Context Harvesting**: Every session autonomously queries the L2 Vault to pull in relevant historical heuristics.
- **Heat Mechanics**: Relevance increases heat, time causes decay — biological memory modeling.

### 3. The Resilient LLM Waterfall

Uptime is non-negotiable. TormentNexus's inference client natively catches 429s (Rate Limits) and 5xx (Server Errors), seamlessly cascading the exact payload down a prioritized chain without crashing:

1. **NVIDIA NIM** / Primary APIs
2. **OpenRouter** (Secondary aggregator fallback)
3. **Local LM Studio / Ollama** (Ultimate offline fallback)

### 4. Multi-Agent Swarm & P2P Mesh

TormentNexus coordinates specialized models inside shared chatrooms via the Agent-to-Agent (A2A) protocol.

- **Role Rotation**: Models take turns acting as Planner, Implementer, Tester, and Critic.
- **Consensus & Debate**: Agents autonomously bid on tasks, share context via a neural transcript, and debate implementations until consensus is reached.
- **PairOrchestrator**: Enforces a strict `Planner → Checker → Implementer → Critic` state machine with weighted consensus scoring.

### 5. Truth Over Hype Dashboards

TormentNexus's dashboards reflect actual SQLite database rows and active Go goroutine states. No mocked UI scaffolds. Monitor telemetry, traffic inspection, working-set capacity, and LLM routing histories in real-time.

### 6. Autonomous Immune System

Every failure is an opportunity for diagnosis, remediation, and verification.

- **HealerService**: Multi-turn `Diagnose → Fix → Verify → Retry` loop using the native `CodeExecutor`.
- **L2 Vault Persistence**: All healing events and extracted facts are saved as long-term memory for fleet-wide intelligence sharing.
- **Supervisor Nudge Protocol**: Autonomously maintains development momentum by re-engaging inactive agents through professional, context-aware directives.

---

## Monorepo Structure

TormentNexus is a **pnpm monorepo** with four major layers:

```
tormentnexus/
├─ apps/                          # Operator-facing applications
│  ├─ web/                        # Next.js 16 dashboard (primary browser UI)
│  ├─ maestro/                    # Electron desktop shell
│  ├─ maestro-go/                 # Go-adjacent desktop lane (experimental)
│  ├─ vscode/                     # VS Code extension
│  ├─ tormentnexus-extension/      # Browser extension
│  ├─ cloud-orchestrator/          # Nested cloud stack (mini-monorepo)
│  └─ borg-extension/              # Borg extension
│
├─ packages/                       # Shared libraries & TypeScript control plane
│  ├─ core/                        # Main TS control plane, tRPC routers, orchestration
│  ├─ cli/                         # CLI entrypoint (`tormentnexus` command)
│  ├─ ui/                          # Shared React UI components (Radix + Tailwind)
│  ├─ ai/                          # Model/provider SDK integration layer
│  ├─ memory/                      # Memory storage, retrieval, embeddings, vector DB
│  ├─ types/                       # Shared TypeScript types
│  ├─ tools/                       # Tool definitions and helpers
│  ├─ mcp-registry/                # MCP metadata and registry surfaces
│  ├─ mcp-client/                  # MCP client integration
│  ├─ agents/                      # Agent-related logic and adapters
│  ├─ adk/                         # Agent Development Kit layer
│  ├─ search/                      # Search and indexing support
│  ├─ tormentnexus-supervisor/      # Windows supervisor bridge (UI automation)
│  ├─ browser/                     # Legacy browser support
│  ├─ browser-extension/          # Shared browser-extension package
│  ├─ enterprise/                  # Enterprise features (SSO, RBAC stubs)
│  ├─ jetbrains/                  # JetBrains IDE integration
│  ├─ zed-extension/              # Zed editor extension
│  └─ tsconfig/                   # Shared TypeScript configuration
│
├─ go/                             # Go modular monolith (authoritative kernel)
│  ├─ cmd/tormentnexus/            # Go entrypoint
│  ├─ internal/
│  │  ├─ httpapi/                  # HTTP API server (600+ endpoints, 18K+ lines)
│  │  ├─ tools/                    # 3,900+ native Go tool implementations
│  │  ├─ harnesses/                # Harness registry and management
│  │  ├─ memory/                   # L1/L2 memory manager
│  │  ├─ memorystore/             # SQLite vault, hydration, search
│  │  ├─ mcp/                     # MCP native router, progressive disclosure
│  │  ├─ orchestration/           # PairOrchestrator, A2A skill registry
│  │  ├─ healer/                 # Autonomous immune system
│  │  ├─ llm/                    # Waterfall routing, provider abstraction
│  │  ├─ eventbus/               # High-frequency resilient message broker
│  │  ├─ codeexec/               # Sandboxed code execution
│  │  ├─ config/                 # Configuration management
│  │  ├─ vault/                  # Secure persistence layer
│  │  ├─ sync/                   # MCP sync and assimilation
│  │  ├─ graph/                  # Dependency graph analysis
│  │  ├─ repograph/              # Repository graph visualization
│  │  ├─ license/                # Ed25519 license validation
│  │  ├─ supervisor/             # Supervisor automation
│  │  ├─ session/                # Session management
│  │  ├─ sessionimport/          # Session import pipeline
│  │  ├─ skillregistry/          # Skill registry with Jaccard deduplication
│  │  ├─ toolregistry/           # Native tool registry
│  │  ├─ flightrecorder/          # Audit and telemetry
│  │  ├─ metrics/               # Provider performance metrics
│  │  ├─ marketplace/            # Tool marketplace
│  │  ├─ gossip/                 # P2P mesh communication
│  │  ├─ mesh/                  # Mesh networking
│  │  ├─ process/               # Process management
│  │  ├─ workspaces/            # Workspace management
│  │  ├─ workflow/             # Workflow engine
│  │  ├─ git/                  # Git operations
│  │  ├─ gitservice/           # Git service layer
│  │  ├─ ctxharvester/         # Context harvesting
│  │  ├─ hsync/                # Harness synchronization
│  │  ├─ interop/               # Language interoperability
│  │  ├─ toon/                 # Animation/graphics utilities
│  │  ├─ ai/                   # AI integrations (Go)
│  │  ├─ buffer/               # Buffer management
│  │  ├─ cache/                # Caching layer
│  │  ├─ controlplane/         # Go control plane
│  │  ├─ buildinfo/            # Build information
│  │  ├─ lockfile/             # Lock file management
│  │  ├─ submodules/           # Submodule management
│  │  └─ providers/            # Provider integrations (Go)
│  └─ go.mod, go.sum             # Go module dependencies
│
├─ data/                           # Local knowledge assets
│  ├─ assimilation_state.db         # MCP assimilation tracking (14,250 rows)
│  ├─ bobbybookmarks/              # Bookmark ecosystems for catalog updates
│  ├─ prompt_library.db             # Prompt library (planned)
│  └─ assimilate_skills.py         # Skill ingestion script
│
├─ docs/                           # Comprehensive documentation
│  ├─ API_ENDPOINTS.md           # 600+ endpoint reference
│  ├─ PROJECT_STRUCTURE.md         # Full module diagram
│  ├─ ARCHITECTURE.md             # Architecture overview
│  ├─ UNIVERSAL_LLM_INSTRUCTIONS.md # Agent coordination rules
│  ├─ GLOBAL_LIBRARY_INDEX.md     # Global MCP library index (2.3MB)
│  ├─ BUILTIN_TOOLS_EVIDENCE_LOCK.md # Evidence lock status
│  └─ [40+ more docs]             # Guides, protocols, security FAQ
│
├─ scripts/                        # Workspace build/dev/maintenance scripts
│  ├─ build_all.mjs               # Full monorepo build
│  ├─ dev_tabby_ready.mjs         # Development launcher
│  ├─ check_release_gate.mjs      # CI release gate
│  ├─ sync_versions.mjs            # Version synchronization
│  └─ [20+ more scripts]          # Validation, indexing, pruning
│
├─ submodules/                     # External upstream assimilations
│  ├─ tormentnexus/                # External CLI harness upstream
│  └─ tormentnexus-mcp/            # External MCP reference lane
│
├─ swarm.py, swarm_v7.py          # Automated Go tool generation swarm
├─ tormentnexus.db               # Main registry (MCP servers, tools, sessions)
├─ catalog.db                    # 11,024+ populated MCP server catalog
├─ provider_metrics.db           # Provider telemetry database
├─ README.md                     # This file
├─ ROADMAP.md                    # Active development roadmap
├─ CHANGELOG.md                  # Detailed version history (400+ entries)
├─ AGENTS.md                     # Multi-agent coordination protocol
├─ VISION.md                     # North star and philosophical pillars
├─ MEMORY.md                     # Accumulated multi-agent insights
├─ HANDOFF.md                    # Session handoff protocol
├─ package.json                  # Root monorepo scripts
└─ pnpm-workspace.yaml           # Workspace boundaries
```

---

## The Go Sidecar

The Go sidecar is the **authoritative execution kernel** of TormentNexus. It is a high-performance modular monolith with 40+ internal packages and 3,900+ native tool implementations.

### Why Go?

- **Performance**: Single-binary deployment, zero runtime dependencies, sub-millisecond API latency
- **Reliability**: Native goroutine-based concurrency, robust error handling, memory safety
- **Portability**: Cross-compilation to any platform; runs on Windows, macOS, Linux, and WSL
- **MCP Assimilation**: Replacing fragile Node.js/Python MCP servers with compiled, type-safe Go handlers

### Native Go Tool Categories

| Category | Native Tools | Example Implementations |
|----------|-------------|------------------------|
| Web Search & Scraping | DuckDuckGo, Firecrawl, Exa | `ddg_search.go`, `firecrawl.go`, `exa.go` |
| Academic & Research | arXiv, Semantic Scholar | `arxiv.go`, `semantic_scholar.go` |
| Databases & Storage | SQLite, ChromaDB, Mem0 | `sqlite_mcp.go`, `chroma.go`, `mem0.go` |
| Finance & Markets | Alpaca, Alpha Vantage, DexPaprika | `alpaca.go`, `alpha_vantage.go`, `dexpaprika.go` |
| Cloud & DevOps | Vercel, Filesystem | `vercel.go`, `filesystem.go` |
| AI & LLM | Ollama, Hugging Face, MindsDB | `ollama.go`, `huggingface.go`, `mindsdb.go` |
| Security | Semgrep, AST-grep | `semgrep.go`, `ast_grep.go` |
| Code Intelligence | Serena, GitIngest, ripgrep | `serena.go`, `gitingest.go`, `ripgrep_search.go` |
| Communication | Slack | `slack.go` |
| Media | TTS, Browser Automation | `tts.go`, `browser_automation.go` |
| Provider Abstraction | PAL (multi-model routing) | `pal.go` |
| Memory | Basic Memory, Thoughtbox | `basic_memory.go`, `thoughtbox.go` |
| System | codemod, anyquery | `codemod.go`, `anyquery.go` |
| Weather | NWS (National Weather Service) | `nws_weather.go` |

### Swarm Tool Generation

The project includes an automated **swarm pipeline** (`swarm.py` / `swarm_v7.py`) that generates native Go tool implementations from MCP server specifications. This enables rapid assimilation of the MCP ecosystem at scale.

- **5 workers** with `--forever` mode for continuous generation
- **200 task limit** per run
- **Parallel batch validation** of generated tools
- **Self-healing**: Automatically removes broken implementations and regenerates

---

## The Dashboard

The **TormentNexus Dashboard** (`apps/web`) is a rich Next.js 16 / React 19 / Tailwind CSS 4 operator interface providing real-time observability and control over the entire system.

### Dashboard Pages

| Route | Purpose |
|-------|---------|
| `/` | Landing page with system overview |
| `/dashboard` | Main command center |
| `/dashboard/brain` | Knowledge graph visualization (force-graph, Mermaid) |
| `/dashboard/chronicle` | Healer / Immune System — active pathogens, L2 vault records |
| `/dashboard/blocks` | Block-based workflow builder |
| `/dashboard/claude-chrome` | Claude Chrome integration console |
| `/dashboard/claude-cloud` | Claude Cloud integration console |
| `/dashboard/copilot` | GitHub Copilot integration |
| `/dashboard/code` | Code execution and sandbox |
| `/dashboard/code/sandbox` | Secure sandbox environment |
| `/dashboard/autopilot` | Autopilot configuration and monitoring |
| `/dashboard/audit` | Audit log viewer with real-time events |
| `/dashboard/config` | System configuration editor |
| `/dashboard/architecture` | Architecture visualization |
| `/dashboard/command` | Command center with cheatsheet |
| `/dashboard/billing` | Billing and subscription management |
| `/dashboard/council` | Multi-agent council visualizer |
| `/dashboard/director` | Director chat interface |
| `/dashboard/swarm` | Swarm mission control |
| `/dashboard/symbols` | Symbol search and LSP integration |
| `/dashboard/vault` | L2 vault memory browser |
| `/dashboard/workflow` | Workflow orchestration |

### UI Components

- **Glassmorphic dark mode** design language
- **Real-time SSE streaming** from Go sidecar
- **Drag-and-drop dashboard** (`dnd-kit`)
- **Knowledge graph** with `@xyflow/react` and `react-force-graph-2d`
- **Charts and telemetry** via `recharts`
- **Animations** via `framer-motion`
- **Toast notifications** via `sonner`

---

## MCP Ecosystem

TormentNexus is the world's largest local MCP registry and the only system with **native Go assimilation** at scale.

### Registry Scale

- **14,250+** total tracked MCP servers (`assimilation_state.db`)
- **11,024+** populated with verified metadata (`catalog.db`)
- **600+** verified and registered in production (`tormentnexus.db`)
- **11,000+** verified individual tools
- **3,900+** reimplemented as native Go handlers (eliminating external dependencies)

### Progressive Disclosure Pipeline

```
User Prompt
    ↓
Vector Embedding (sqlite-vec)
    ↓
BM25 + Semantic Search (Top-K matching)
    ↓
Tool Schema Injection (only relevant schemas)
    ↓
LLM Context Window (clean, relevant, small)
```

### Assimilation Categories

1. **Developer Tools** — GitIngest, ripgrep, codemod, anyquery, AST-grep, Serena
2. **Databases** — SQLite, ChromaDB, Mem0, Basic Memory
3. **Web Search** — DuckDuckGo, Firecrawl, Exa, arXiv, Semantic Scholar
4. **Communication** — Slack
5. **Cloud** — Vercel, Filesystem
6. **AI/LLM** — Ollama, Hugging Face, MindsDB, PAL (multi-model)
7. **Finance** — Alpaca, Alpha Vantage, DexPaprika, Octagon
8. **Security** — Semgrep
9. **Weather** — NWS
10. **Browser** — Playwright, chromedp automation
11. **Sandbox** — Thoughtbox (Node VM wrapper)
12. **Media** — TTS

---

## Memory & Context

### The Hippocampus Model

```
┌─────────────────────────────────────────┐
│  L1 — Session Scratchpad (Active)       │
│  · In-memory, ephemeral                 │
│  · Lightning-fast retrieval             │
│  · Tied to active session ID              │
│  · ~4K-8K token window                    │
├─────────────────────────────────────────┤
│  L2 — The Vault (Semantic)              │
│  · SQLite + sqlite-vec embeddings       │
│  · Permanent, searchable                │
│  · Heat-score lifecycle management        │
│  · LLM-compressed heuristics              │
│  · Exact transcript preservation          │
├─────────────────────────────────────────┤
│  L3 — Cold Archive (Planned)             │
│  · Long-term cold storage               │
│  · Compressed, summarized                 │
│  · Retrieved only on explicit demand      │
└─────────────────────────────────────────┘
```

### Context Harvesting

Every session autonomously queries the L2 Vault before responding, pulling in relevant historical heuristics based on semantic similarity and heat scores. This ensures the model never starts from zero context.

### TrafficObserver

A passive fact extraction layer that monitors all system traffic (tool calls, LLM responses, errors) and automatically persists facts into the L2 Vault without explicit user action.

---

## Swarm & Multi-Agent

### A2A Protocol Coordination

TormentNexus implements the **Agent-to-Agent (A2A)** protocol for multi-agent coordination inside shared chatrooms.

### Role Rotation State Machine

```
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Planner  │───→│ Checker  │───→│Implementer│───→│  Critic  │
│ (Design) │    │ (Verify) │    │ (Build)  │    │(Validate)│
└──────────┘    └──────────┘    └──────────┘    └────┬─────┘
     ↑─────────────────────────────────────────────────┘
                    (Consensus Loop)
```

### Consensus Engine

- Weighted voting based on model confidence
- Debate rounds until threshold consensus is reached
- Neural transcript sharing for context synchronization
- Task bidding and autonomous delegation

---

## API Surface

The Go sidecar exposes **600+ HTTP endpoints** organized into:

### Core API Categories

| Category | Endpoints | Examples |
|----------|-----------|----------|
| **Health & System** | 5 | `/health`, `/version`, `/api/index` |
| **Configuration** | 8 | `/api/config/*`, `/api/config/mcp-timeout` |
| **MCP** | 10 | `/api/mcp/tools`, `/api/mcp/tools/call`, `/api/mcp/tools/search` |
| **Skills** | 12 | `/api/skills/list`, `/api/skills/search`, `/api/skills/assimilate` |
| **Memory** | 7 | `/api/memory/list`, `/api/memory/search`, `/api/memory/hydrate` |
| **Agents & Swarm** | 6 | `/api/swarm/start`, `/api/squad/spawn`, `/api/supervisor/decompose` |
| **Governance** | 5 | `/api/api-keys`, `/api/audit`, `/api/autonomy/set-level` |
| **DevOps** | 5 | `/api/git/status`, `/api/submodules`, `/api/scripts/execute` |
| **Code & Symbols** | 4 | `/api/code/exec`, `/api/graph`, `/api/lsp/find-symbol` |
| **Sessions** | 4 | `/api/sessions/imported/scan`, session import pipeline |

### Response Format

All API endpoints return a standardized envelope:

```json
{
  "success": true,
  "data": { ... }
}
```

---

## Quick Start

### Prerequisites

- **Node.js 24+**
- **Go 1.26+**
- **pnpm 10+**

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/NexusSoftMDMA/TormentNexus.git
cd tormentnexus

# 2. Install dependencies & rebuild SQLite bindings
pnpm install
pnpm rebuild better-sqlite3

# 3. Build the Go sidecar
cd go && go build -buildvcs=false ./cmd/tormentnexus && cd ..

# 4. Start the TormentNexus Control Plane
pnpm run dev
```

The Next.js dashboard will automatically open at `http://localhost:3000/dashboard` once the TypeScript Control Plane (port 4100) and Go Sidecar (port 4300) are successfully locked and humming.

### Development Scripts

| Script | Purpose |
|--------|---------|
| `pnpm run dev` | Start full development environment |
| `pnpm run dev:web` | Start only the web dashboard |
| `pnpm run build` | Full production build |
| `pnpm run test` | Run all test suites |
| `pnpm run check:release-gate` | CI release gate validation |
| `pnpm run index:sync` | Sync master MCP index |
| `pnpm run clean` | Clean all build artifacts |

---

## What's Planned

### Phase 6: Comprehensive Assimilation & Enterprise Readiness (Active)

| Track | Goal | Status |
|-------|------|--------|
| **Track A** | Assimilate top 500 MCP servers as native Go modules | In Progress (3,900+ done) |
| **Track B** | Hermes addons & prompt library migration | Experimental |
| **Track C** | Enterprise licensing (Ed25519) + SSO/RBAC | Experimental |
| **Track D** | Default agent harness integration (Tabby, Warp, Hyper, Hermes, Pi) | Beta |

### Phase 7: Session Continuity & Deep Linking (In Progress)

- **Session Import Pipeline**: Automatic ingestion of Claude, Aider, and harness artifacts (49 candidates discovered, 586 imported sessions)
- **`tormentnexus://` Protocol**: Browser-to-kernel deep linking for seamless IDE attachment
- **Wails Native Runtime**: Replacing Electron with a Go-native desktop shell

### Phase 8: Predictive Intelligence (Vision)

- **Predictive Conversational Tool Injection**: Local model-based prediction of relevant tools before the user asks
- **L3 Cold Archive**: Long-term compressed memory tier for infinite context
- **Fleet-Wide Intelligence**: Cross-machine memory sharing via encrypted mesh

### Long-Term Vision

- **The AI TormentNexus**: The operating system for AI models, abstracting all provider complexity
- **Models as Compute**: Ephemeral resources managed by allocation, fallback routing, and token budgets
- **Tools as Drivers**: MCP servers as "device drivers" for the AI OS
- **Biological Memory**: L1/L2/L3 tiers with heat-based mechanics
- **Autonomous Immune System**: Self-healing through diagnosis, remediation, and verification

---

## Roadmap

See [ROADMAP.md](ROADMAP.md) for the detailed path to v1.0.0 stable, including:

- Progressive Skill Disclosure (Context Hygiene)
- Go-native MCP sync migration
- Native UI replacement for Electron (Maestro/Go)
- Full L3 cold archive implementation
- Cross-platform binary distribution
- Fleet telemetry and cross-machine memory mesh

---

## Documentation

| Document | Purpose |
|----------|---------|
| [AGENTS.md](AGENTS.md) | Multi-agent coordination and handoff protocol |
| [VISION.md](VISION.md) | North star and philosophical pillars |
| [ROADMAP.md](ROADMAP.md) | Active development roadmap and milestones |
| [CHANGELOG.md](CHANGELOG.md) | Detailed version history (400+ entries) |
| [MEMORY.md](MEMORY.md) | Accumulated multi-agent insights and gotchas |
| [HANDOFF.md](HANDOFF.md) | Session handoff protocol for agent continuity |
| [docs/API_ENDPOINTS.md](docs/API_ENDPOINTS.md) | Complete 600+ endpoint reference |
| [docs/PROJECT_STRUCTURE.md](docs/PROJECT_STRUCTURE.md) | Full module dependency diagram |
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Architecture deep dive |
| [docs/UNIVERSAL_LLM_INSTRUCTIONS.md](docs/UNIVERSAL_LLM_INSTRUCTIONS.md) | Agent coordination rules (read first) |
| [docs/GLOBAL_LIBRARY_INDEX.md](docs/GLOBAL_LIBRARY_INDEX.md) | Global MCP library index (2.3MB) |
| [docs/LAUNCH_METRICS.md](docs/LAUNCH_METRICS.md) | Launch metrics and telemetry |
| [docs/LAUNCH_SECURITY_FAQ.md](docs/LAUNCH_SECURITY_FAQ.md) | Security FAQ |
| [docs/GO_SIDECAR_API.md](docs/GO_SIDECAR_API.md) | Go sidecar API documentation |

---

## Contributing

TormentNexus is built by a **multi-agent swarm** of specialized AI models (Gemini, Claude, GPT) coordinated through the [AGENTS.md](AGENTS.md) protocol. Human operators supervise and validate.

### Model Specializations

| Model | Strengths | Focus Areas |
|-------|-----------|-------------|
| **Gemini** | Speed, massive context, repo maintenance | Bulk refactoring, recursive scripts, context analysis |
| **Claude** | UI/UX perfection, documentation, deep features | Responsive layouts, type safety, precise documentation |
| **GPT** | Systemic architecture, distributed debugging, race conditions | Go/TS bridge contracts, DB migration, concurrency safety |

### Session Protocol

1. **Read** `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` and `AGENTS.md`
2. **Check** `VERSION` and `HANDOFF.md` for current state
3. **Run** git checks to ensure workspace cleanliness
4. **Work** autonomously unless changes are destructive or ambiguous
5. **Commit** small, incremental, verifiable changes
6. **Update** `HANDOFF.md`, `MEMORY.md`, and `CHANGELOG.md`
7. **Bump** `VERSION` and sync with `node scripts/sync-versions.mjs`
8. **Push** to `origin` and `tormentnexus-upstream`

---

## License

TormentNexus uses a **dual-licensing model**:

- **Open Source**: Core kernel and tool implementations available under a permissive open-source license
- **Enterprise**: Ed25519-signed license tokens for advanced features (SSO, RBAC, fleet management, cross-machine mesh)

See [docs/LAUNCH_SECURITY_FAQ.md](docs/LAUNCH_SECURITY_FAQ.md) for licensing details.

---

> *Praise the LORD! Keep on going! Don't ever stop! Don't stop the party!!!*
>
> The collective grows. 🚀
