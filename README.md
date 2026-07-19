# TormentNexus — Universal AI Control Plane

> **One-command install:** `npx @tormentnexus/install`
>
> Local-first, open-source AI control plane with persistent multi-tier memory, MCP tool orchestration across 20,000+ servers, multi-agent swarm coordination, and corporate-grade security — all running on your machine.
>
> **📄 [Read the Technical Whitepaper](https://tormentnexus.site/tormentnexus_whitepaper.html)** | **💬 [Join Discord](https://discord.gg/Hj9P3GbVxR)**

## Installation

### Quick Start (Windows)

```bash
scripts\install.bat
```

### Quick Start (macOS / Linux)

```bash
python3 scripts/install-client-support.py
```

### GUI Installer

```bash
python3 scripts/install-gui.py
```

### What Gets Installed

The installer detects and configures TormentNexus for **38+ AI coding agents**.
Each client receives every integration it supports:

| Client | SKILL | MCP | CMD | HOOK | EXT | AGENT |
|--------|-------|-----|-----|------|-----|-------|
| Claude Code | x | x | x | x | - | - |
| Cursor | x | x | x | x | x | - |
| Windsurf | x | x | x | x | x | - |
| Gemini | x | x | x | - | x | x |
| Codex CLI | x | x | x | - | - | - |
| Grok Build | x | x | x | x | x | - |
| Antigravity | x | x | - | - | x | x |
| Aider | - | x | - | - | - | - |
| CodeWhale | x | x | x | x | x | - |
| Goose | x | x | x | x | x | - |
| Pi | x | x | x | - | x | - |
| Cline / Roo | x | x | x | x | - | - |
| + 26 more | x | x | varies | varies | varies | varies |

### Corporate On-Premises Installation

For enterprises running TormentNexus behind a firewall:

```bash
python3 scripts/install-client-support.py --mode corporate
scripts\install.bat corporate
```

Corporate mode adds:

- **SSO/OIDC** configuration for enterprise identity providers
- **RBAC roles** with admin/user/auditor defaults
- **Audit logging** to local SQLite with daily rotation
- **Multi-tenant isolation** for on-premises tenant separation
- **Enterprise license** key setup

### Server Deployment

```bash
GOOS=linux GOARCH=amd64 go build -o tn-kernel ./cmd/tormentnexus
pm2 start tn-kernel --name tn-kernel -- serve -port 8090 -host 127.0.0.1
```

---

# TormentNexus: The Cognitive Kernel — Universal Open-Source AI Control Plane for Multi-Agent Workflows, MCP Tools & Context-Aware Memory

![Version](https://img.shields.io/badge/version-1.0.0--b1-blue)
![Build](https://img.shields.io/badge/build-passing-brightgreen)
![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go)
![TypeScript](https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript)
![Next.js](https://img.shields.io/badge/Next.js-16-000000?logo=next.js)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)
![License](https://img.shields.io/badge/license-MIT%2FOSS-orange)

> **TormentNexus** is a local-first, open-source universal AI control plane and coding harness engineered to give your LLM agents infinite persistent memory, tool orchestration, and autonomous infrastructure assimilation—all running locally on your machine or in production.

---

## 🏢 Corporate & Enterprise Deployments (HyperNexus)

For professional, corporate, and enterprise deployments, visit **[HyperNexus](https://hypernexus.site)**.
**HyperNexus** is the official licensed, certified, authoritative, cloud-hosted corporate installation of the TormentNexus kernel. It provides:

- **Tenant Isolation**: Automated provisioning scripts to spin up per-tenant containers with isolated volumes and resource limits.
- **Enterprise-Grade Security**: Full SSO support and fine-grained Role-Based Access Control (RBAC) middleware.
- **Air-Gapped & VPC Support**: Production-ready deployment configurations with daily audit logs for SIEM integration.
- **Stripe Billing Integration**: Built-in webhook configurations and subscription seats management.

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
- [Roadmap](ROADMAP.md)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

---

## What It Does

TormentNexus is a **decision system and universal bridge** — not just an aggregator. It runs locally as a modular monolith that unifies the chaotic landscape of AI tools, models, and agents into a single, coherent operating system for AI-driven development.

### Current Capabilities (v1.0.0-b1)

| Capability | Status | Details |
|------------|--------|---------|
| **MCP Registry** | Stable | Infinite tracked MCP servers, unlimited populated in SQLite catalog, infinite verified servers, unlimited verified tools |
| **Native Go Tools** | Beta | Unlimited native Go tool implementations replacing external MCP servers (filesystem, Slack, SQLite, DuckDuckGo, Ollama, TTS, Vercel, NWS, DexPaprika, Firecrawl, Exa, arXiv, Semantic Scholar, Mem0, Alpaca, Alpha Vantage, Hugging Face, Semgrep, Octagon, Browser Automation, ChromaDB, Basic Memory, MindsDB, Serena, AST-grep, PAL, Thoughtbox, and more) |
| **Progressive Tool Routing** | Stable | Semantic vector search + BM25 ranking injects only the most relevant tools into LLM context windows |
| **Dual-Tier Memory** | Stable | L1 (session scratchpad) + L2 (semantic SQLite vault) with heat-score lifecycle and autonomous context harvesting |
| **LLM Waterfall** | Stable | Cascading failover: NVIDIA NIM → OpenRouter → Local LM Studio / Ollama with 429/5xx handling |
| **Multi-Agent Swarm** | Beta | A2A protocol coordination, role rotation (Planner→Implementer→Tester→Critic), consensus engine |
| **Autonomous Healer** | Stable | `Diagnose → Fix → Verify → Retry` loop with native code execution (tsc, vitest, go test) and L2 vault persistence |
| **Browser Automation** | Beta | Native chromedp handlers: navigate, screenshot, evaluate, click, fill, get HTML |
| **Skill Registry** | Stable | Unlimited assimilated skills from multiple harness ecosystems (Aider, Agent, CCS, Hermes, Pi, etc.) with Jaccard deduplication |
| **CodeWhale Extension** | Stable | Native Rust extension with full Pi extension parity — infinite MCP tools, unlimited custom tools, infinite slash commands, unlimited shortcuts, L2 memory hooks, RBAC, @memory:key expansion, per-turn context harvesting, session compaction |
| **Dashboard** | Stable | Next.js 16 + React 19 + Tailwind CSS 4 with real-time telemetry, knowledge graph, healer view, swarm visualizer |
| **npm Packages** | Published | `@tormentnexus/install`, `@tormentnexus/core`, `@tormentnexus/cli`, `@tormentnexus/openhands`, `@tormentnexus/vscode`, `@tormentnexus/cursor` |
| **tRPC Bridge** | Stable | Type-safe API layer (port 7778) connecting UI to Go sidecar |
| **Session Import** | Beta | Automatic ingestion of Claude, Aider, and other harness session artifacts |
| **Enterprise Licensing** | Experimental | Ed25519-signed license token validation with offline verification |
| **Supervisor Nudge** | Stable | Autonomous Windows UI automation to maintain development momentum across AI chat surfaces |
| **Deep Link Protocol** | Beta | `tormentnexus://attach?session=ID` and `tormentnexus://create?cliType=aider` URI handling |
| **Provider Metrics** | Stable | Real-time telemetry tracking for all LLM providers with latency and cost analysis |

---

## 🚀 One-Command Install

```bash
npx @tormentnexus/install
```

That's it. Detects and configures TormentNexus for **38+ AI coding agents** including Claude, Gemini, Cursor, Windsurf, Aider, CodeWhale, and more. Each gets MCP config, skills, hooks, and commands.

```bash
# CLI tools
npm install -g @tormentnexus/cli
tn search "your topic"
tn status

# OpenHands integration
npx @tormentnexus/openhands

# Core kernel
npx @tormentnexus/core
```

---

## Architecture

TormentNexus is a **high-performance Go modular monolith** with a **TypeScript/Next.js frontend**, operating as a local-first control plane.

```
┌─────────────────────────────────────────────────────────────┐
│  OPERATOR LAYER                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐       │
│  │  Web Dash   │  │  CLI (TS)   │  │  VS Code    │       │
│  │  Port 7779  │  │  tormentnexus│  │  Extension  │       │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘       │
│         │                │                │                 │
│         └────────────────┴────────────────┘                 │
│                          │                                │
│  ┌───────────────────────┴───────────────────────┐        │
│  │  GO SIDECAR (Port 7778) — The Authoritative Kernel    │
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
| Go Sidecar | 7778 | Authoritative native kernel (HTTP API + tRPC) |
| Next.js Dashboard | 7779 | Web observation deck |

---

## Core Pillars

### 1. Progressive MCP Tool Routing & Parity

Models should never be overwhelmed with a 50,000-token tool dump. TormentNexus employs a multi-layered, progressive disclosure system:

- **Semantic Search**: Local vector embeddings match the active prompt against a global MCP directory of infinite servers.
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
│  └─ tormentnexus-extension/              # Tormentnexus extension
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
├─ catalog.db                    # unlimited populated MCP server catalog
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

The Go sidecar is the **authoritative execution kernel** of TormentNexus. It is a high-performance modular monolith with dozens of internal packages and unlimited native tool implementations.

### Why Go?

- **Performance**: Single-binary deployment, zero runtime dependencies, sub-millisecond API latency
- **Reliability**: Native goroutine-based concurrency, robust error handling, memory safety
- **Portability**: Cross-compilation to any platform; runs on Windows, macOS, Linux, and WSL
- **MCP Assimilation**: Replacing fragile Node.js/Python MCP servers with compiled, type-safe Go handlers

### Native Go Tool Categories

| Category | Example Integrations |
|----------|---------------------|
| Web Search & Scraping | DuckDuckGo, Firecrawl, Exa |
| Academic & Research | arXiv, Semantic Scholar, PubMed, Open Library |
| Finance & Markets | Alpaca, Alpha Vantage, CoinGecko, Yahoo Finance |
| Science & Space | NASA (APOD, NEO, Mars Weather, EPIC), Open-Meteo |
| Entertainment | PokéAPI, TVMaze, JokeAPI |
| Reference | World Bank, Nobel Prize, Open Library |
| AI & LLM | Ollama, Hugging Face, MindsDB |
| Security | Semgrep, AST-grep |
| Communication | Slack |
| Cloud & DevOps | Vercel, Filesystem |
| Memory | Basic Memory, ChromaDB, Mem0 |

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

- **20,000+** total tracked MCP servers (`assimilation_state.db`)
- **Thousands** verified individual tools
- **Continuous native Go reimplementation** eliminating external dependencies

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
# 1. One-command setup (recommended)
npx @tormentnexus/install

# 2. Or clone and build from source
git clone https://github.com/MDMAtk/TormentNexus.git
cd tormentnexus
pnpm install && pnpm rebuild better-sqlite3
cd go && go build -buildvcs=false ./cmd/tormentnexus && cd ..
pnpm run dev
```

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
| [VISION.md](docs/VISION.md) | North star and philosophical pillars |
| [ROADMAP.md](ROADMAP.md) | Active development roadmap and milestones |
| [CHANGELOG.md](CHANGELOG.md) | Detailed version history |
| [MEMORY.md](docs/MEMORY.md) | Accumulated multi-agent insights and gotchas |
| [HANDOFF.md](docs/HANDOFF.md) | Session handoff protocol for agent continuity |
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

## Community

### pi Extension (npm)

```bash
pi install npm:tormentnexus
```

Connects pi to the TormentNexus control plane — persistent L2 memory, tool search, session import, skill registry, code search, subagent orchestration, enterprise RBAC. Full source at `packages/tormentnexus/`.

**CodeWhale now has equivalent functionality through a native Rust extension** — see below.

### CodeWhale Integration

CodeWhale has **full Pi extension parity** through two integration layers:

#### Layer 1: Native Rust Extension (tn-extension)

A native Rust crate at `crates/tn-extension` compiled into the CodeWhale binary. Hooks fire automatically on every lifecycle event:

| Hook | What it does |
|------|-------------|
| **SessionStart** | Logs session to TN L2 memory |
| **BeforeAgentStart** | Injects TN system prompt + searches L2 for relevant context per turn |
| **ToolCall** | Logs tool + args; checks 6 dangerous patterns against RBAC |
| **ToolResult** | Auto-stores substantial results from key tools to L2 |
| **TurnEnd** | Logs tool usage summary per turn |
| **Input** | Expands `@memory:key` inline with L2 content |
| **UserBash** | Audit-logs shell commands to TN enterprise audit |
| **ModelSelect** | Tracks model changes to L2 |
| **SessionCompact** | Preserves memory across compaction |

Also registers: 9 custom tools, 6 slash commands, 3 keyboard shortcuts, auto-registers the MCP server.

#### Layer 2: CodeWhale Skill (SKILL.md)

Tracked at `.codewhale/plugins/tormentnexus/skills/SKILL.md`, auto-installed to `~/.codewhale/skills/tormentnexus/`. Contains 49 MCP tools, REST API reference, slash command handling, best practices, security notes.

#### Installation

```bash
scripts\install_codewhale.bat
```

Verify:

```bash
codewhale mcp connect tormentnexus
codewhale mcp tools | findstr mcp_tormentnexus | find /c /v ""
curl http://127.0.0.1:7778/api/health
```

Build from `~/codewhale-source`:

```bash
cargo build --release -p codewhale-cli
```

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
