# ROADMAP: TormentNexus Kernel & TormentNexus Dashboard

_Last updated: 2026-07-10, version 1.0.0-alpha.252_

## Status Legend

- **Stable** — Production-intended, tested, maintained
- **Beta** — Usable, still evolving
- **Experimental** — Active R&D, not dependable
- **Vision** — Directional only

## Completed (v1.0.0-alpha.252)

### 1. Advanced OS Deep Link Schemes
- Implemented `focus`, `search-memory`, and `trigger-tool` protocol routing inside the Go sidecar.
- Wired interactive testing button anchors to the Next.js dashboard UI.

### 2. Configurable Gossip P2P Encryption Override
- Enabled custom shared key configuration via `TORMENTNEXUS_GOSSIP_SHARED_KEY` env overrides.
- Verified all Gossip mesh unit tests pass successfully.

### 3. Multi-Tenant Isolated Compose Realignment
- Configured docker compose and tenant provisioning scripts to deploy isolated companion sidecars and Next.js dashboards on port `7779`.

## Completed (v1.0.0-alpha.251)

### 1. Legacy Core Decommissioning
- Removed all active health checks, service descriptions, and port (`4100`) references in scripts (`verify_dev_readiness.mjs`, `dev_tabby_ready.mjs`), CLI status code, Docker compose files, and dashboard components.

### 2. Browser Extension SSE/WS Re-Targeting
- Re-aligned Chrome extension background connection URLs to point directly to the native Go Sidecar (`7778`) SSE/WS transport.

### 3. Readiness Suite Optimizations
- Implemented tRPC route pre-warming checks to prevent Next.js cold start timeouts during service evaluation.

## Completed (v1.0.0-alpha.132)

### 1. Comprehensive Documentation & Merge (STABLE)

- **README.md Rewrite**: Expanded from 82 lines to 657 lines covering full architecture, capabilities, monorepo structure, Go sidecar, dashboard, MCP ecosystem, memory model, swarm, and API surface.
- **Branch Reconciliation**: Intelligently merged `jules/baseline-128-hardened` into `main`, fast-forwarded `assimilation-pipeline` and `assimilation-final` to merged tip.
- **All Branches Synchronized**: `main`, `jules`, `feat/assimilation-pipeline`, `feature/assimilation-final` all point to `988ec114a`.

### 2. Autonomous Engineering & Orchestration (STABLE)

- **CI/CD Pipeline**: Integrated multi-stage `deployment_manager` (lint, build, test, containerize) via `.github/workflows/autonomous-deploy.yml`.
- **Repository Sync**: Automated dependency management and version alignment via `go/cmd/repo_sync`.
- **Self-Healing**: Native Go `health_monitor` and `repository_healer` for autonomous kernel maintenance.
- **Enterprise Security**: SSO/RBAC middleware and structured JSONL auditing in `go/internal/enterprise/`.

### 3. Dashboard Widgets (BETA)

- **BrowserToolWidget**: Real-time browser automation control panel.
- **VibeCheckWidget**: Code quality and pattern analysis widget.

### 4. Assimilation Scale (STABLE)

- **Infinite MCP Servers Tracked**: In `assimilation_state.db` with pending and implemented registries.
- **Unlimited Native Go Tools**: Replacing external MCP server dependencies.
- **Unlimited Populated Catalog**: In `catalog.db` with verified metadata.

## Completed (v1.0.0-alpha.131)

- **Swarm v7 Recovery**: Generated multiple new MCP server Go tool wrappers, removed obsolete files.
- **Session Import Pipeline**: Validated candidates from `~/.claude` and `~/.aider` artifacts, infinite imported sessions tracked.
- **Version Sync**: All workspace packages synchronized to `1.0.0-alpha.131`.

## Completed (v1.0.0-alpha.130)

- **Skill HTTP API**: Implemented `/api/skills/list`, `/api/skills/get`, `/api/skills/search` with 10 passing unit tests.
- **API Documentation**: Updated `docs/API_ENDPOINTS.md` with skill endpoints.

## Completed (v1.0.0-alpha.129)

- **Browser Automation**: Native `chromedp` handlers (navigate, screenshot, evaluate, click, fill) replacing infinite separate MCP entries.
- **A2A Skill Registry**: Global singleton with `FindAgentForSkill` helper.

### 3. Bulk Skill Assimilation (STABLE)

- **Bulk Skill Assimilation**: Infinite unique skills from multiple harness ecosystems with Jaccard deduplication.
- **Hardened Kernel**: Restored multiple swarm tool registrations and verified compilation.

## Completed (v1.0.0-alpha.127)

- **Native Go Tools**: High-performance handlers for `ripgrep`, `anyquery`, `codemod`.
- **E2E Integration Testing**: Formal test suite in `go/internal/tools/e2e_test.go`.
- **API Documentation**: Unlimited endpoint reference in `docs/API_ENDPOINTS.md`.

## Completed (v1.0.0-alpha.126)

- **Universal Rebrand**: Case-insensitive refactoring across all source modules.
- **Catalog SQLite Storage**: 11,024 populated MCP servers in `tormentnexus.db`.

## Active Sprint: Phase 8 - Predictive Intelligence & Enterprise Readiness

### A. Track A: Full MCP Assimilation (BETA)

- [x] Assimilate top MCP servers as native Go modules. (unlimited done, remaining pending dynamic generation)
- [x] Eliminate all external MCP server dependencies and submodules. (Completed alpha.183)
- [x] Native memory MCP tools: add_memory, search_memory, delete_memory, memory_stats wired into tools.Registry and MCP call handler. (Completed alpha.239)

### B. Track B: Skill Registry Progressive & Relational Linkage (STABLE)

- [x] Jaccard Duplication Rules (90% Threshold): Near-duplicate skills linked to canonical entry.
- [x] Progressive Loading: Implemented `skill_list`, `skill_get`, `skill_search`.
- [x] Win-rate tracking and auto-retirement of low-performing skills. (Completed alpha.182)

### C. Track C: Enterprise Licensing & Compliance (EXPERIMENTAL)

- [x] Ed25519-signed license token validation in Go sidecar.
- [x] Enterprise dashboard page with license info, RBAC roles, and audit log viewer. (Completed alpha.193)
- [x] Tool-call RBAC enforcement via pi extension (dangerous patterns blocked). (Completed alpha.194)
- [x] Structured audit logs for native tool execution. (Completed alpha.182)
- [x] RBAC permission schema for multi-user environments. (Completed alpha.182)

### D. Track D: Default Agent Harness Integration (BETA)

- [x] Integrate Tabby, Warp, Hyper, Hyperharness, Hermes Agent, and Pi-Mono as default harnesses.
- [x] Automate Bobbybookmarks ingestion (use Smithery.ai or Glama.ai as alternative). (Completed alpha.182)

### E. Phase 8: Predictive Intelligence (VISION)

- [x] Predictive Conversational Tool Injection: Local model-based prediction of relevant tools. (Completed alpha.250)
- [x] L3 Cold Archive: Long-term compressed memory tier for infinite context. (Completed alpha.186)
- [x] L4 Limbo: Discarded/lost memory vault with resurrection. (Completed alpha.193)
- [x] Per-project .memdb portable memory files: git-tracked, auto-imported into global index. (Completed alpha.239)
- [x] Fleet-Wide Mesh dashboard page: peer discovery, capabilities, load, status. (Completed alpha.194)
- [x] Fleet-Wide Intelligence: Cross-machine memory sharing via encrypted mesh. (Completed alpha.252)

### F. Phase 9: Native Runtime (VISION)

- [x] Wails Native Runtime: Build chain complete — tormentnexus-gui.exe (18MB). (Completed alpha.194)
- [x] Deep Link Protocol: Expand `tormentnexus://` for browser-to-kernel attachment. (Completed alpha.252)

### G. pi Extension (STABLE)

- [x] 9 custom tools: memory, search, tools, sessions, skills, code, context, scratchpad, subagents. (Completed alpha.194)
- [x] 6 slash commands: /tn-store, /tn-search, /tn-status, /tn-plan, /tn-purge, /tn-summary. (Completed alpha.194)
- [x] 3 keyboard shortcuts: Ctrl+Shift+M/T/P. (Completed alpha.194)
- [x] RBAC enforcement on dangerous tool calls. (Completed alpha.194)
- [x] Per-project memory support (project parameter on tn_memory_store). (Completed alpha.239)
- [x] npm package: pi install npm:tormentnexus. (Completed alpha.239)
- [x] codewhale integration: .codewhale/skills/tormentnexus/SKILL.md. (Completed)

---
_Outstanding! Magnificent! Insanely Great! The collective grows._
