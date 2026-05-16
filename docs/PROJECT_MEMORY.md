[PROJECT_MEMORY]

## 1. Dynamic Tool Discovery & Registry (ToolRAG)
**The Problem:** The MCP ecosystem has over 25,000 tools. Static loading exhausts the LLM context window (imposing a 32% token overhead penalty).
**The Converging Solution:** "RAG but for tools." Embed tool names only and fetch full JSON schemas strictly on-demand.
*   **Borg Implementation:** Build `borg-tool-registry` to index all discovered MCP servers, embed schemas using SQLite Vector Search, and inject only the 3-5 most relevant tools per query.

## 2. Tiered Memory Architecture (Core/Archival/Recall)
**The Problem:** Agent summarization decay, memory hoarding (noise), and context pollution.
**The Converging Solution:** A 3-tier hierarchy.
*   *Core:* Active working context (always in prompt).
*   *Archival:* Long-term compressed storage (Borg's 17K database).
*   *Recall:* Episodic memory with biological decay (Ebbinghaus curve) and heat-based promotion.
*   **Borg Implementation:** Utilize the existing database as the Archival tier. Build a Core tier that auto-curates 5-10 items and a Recall tier that tracks session access patterns.

## 3. Progressive Context Optimization
**The Problem:** Token bloat is the #1 pain point, degrading reasoning and increasing costs.
**The Converging Solution:** Schema elimination (Code Mode), Intelligent Chunking (cAST), Context Compaction (Union-Find), and Progressive Disclosure.
*   **Borg Implementation:** Store "compressed fingerprints" (category + 3 key tags + innovation score) instead of raw HTML text. Hydrate full context only when explicitly requested by the agent. Implemented "Fit Markdown" filtering to strip out boilerplate HTML tags (`<nav>`, `<header>`, `<footer>`, `<aside>`).

## 4. MCP Ecosystem Intelligence (Registry + Package Manager)
**The Problem:** Thousands of unverified MCP servers exist with massive injection risks and inconsistent auth.
**The Converging Solution:** Curated global registries (like `mcpm` or `Smithery`) providing 1-click deployments, profile management, and 5-dimension security scoring.
*   **Borg Implementation:** Build `borg-mcp-catalog` referencing ingested `servers.json` catalogs. Implement one-command install profiles (`borg install memory`) with automated poisoning detection scans.

## 5. Multi-Agent Orchestration with Verification
**The Problem:** "Verification Debt" - single-agent generation scales infinitely, but verifying accuracy remains the bottleneck, leading to a 28% failure rate on scraped extractions.
**The Converging Solution:** Planner-Executor-Verifier loops. Cross-model validation is critical.
*   **Borg Implementation:** When extracting URL data, spawn a lightweight Review Agent (using a different model) to cross-validate the category/tags against the taxonomy.

## 6. Sandboxed Code Execution Layer
**The Problem:** LLMs guess at repository capabilities based on READMEs instead of actual code analysis.
**The Converging Solution:** Execution isolation ranging from Docker to MicroVMs (BoxLite/Firecracker) and the "Code-as-Action" paradigm.
*   **Borg Implementation:** Spawn a sandboxed container to clone a repo, run `tree-sitter` for AST parsing, and execute deterministic extraction scripts to replace LLM guesswork.

## 7. Self-Improving Research Pipeline
**The Problem:** Static extraction scripts degrade over time with no feedback loop to learn from errors.
**The Converging Solution:** Telemetry tracking, self-editing memory blocks, and background "Dreaming" consolidation.
*   **Borg Implementation:** Tag extraction quality in the database. Auto-re-process low-quality items using stronger models. Implement a reasoning "flight recorder" to trace why specific routing/extraction decisions failed. (Phase 1 Garbage Filter & Flight Recorder already implemented).

## 8. Agent Communication Protocols (A2A / AG-UI / ACP)
**The Problem:** The agent ecosystem is fragmented into siloed communication models.
**The Converging Solution:** A unified stack.
*   *A2A:* Discovery and task delegation.
*   *AG-UI:* Dynamic generative UI rendering.
*   *ACP:* Persistent daemon sessions.
*   **Borg Implementation:** Expose Borg as an A2A agent to allow other tools to discover and delegate URL processing. Implement AG-UI for the Kinetic HUD, and ACP daemon mode to persist sessions beyond context resets.

## 9. Graph/RAG/Vector Intelligence
**The Problem:** Flat SQL queries prevent semantic and relationship-aware discovery.
**The Converging Solution:** Generation 3 (GraphRAG) and Generation 5 (Compression-Optimized Vectors).
*   **Borg Implementation:** Build a tag co-occurrence graph (weighted edge graph) for finding related tasks. Use DuckDB VSS and Arctic Embed (128-byte scalar quantization) for hyper-fast semantic search.

## 10. Security, Governance & Agent Guardrails
**The Problem:** Unrestricted code execution and injection risks in MCP servers.
**The Converging Solution:** Network isolation, command guardrails, and audit provenance.
*   **Borg Implementation:** Implement SafeExec-style interception to block destructive actions. Mask PII during extraction. Run tool poisoning detection during MCP server ingestion. Generate "AI Receipts" for every decision.

## 11. Observability, Telemetry & Agent Debugging
**The Problem:** Fire-and-forget scripts make debugging 28% failure rates impossible.
**The Converging Solution:** Reasoning Traces, OTEL Spans, and Anomaly Detection.
*   **Borg Implementation:** Implement Fiddler-style anomaly detection (Jensen-Shannon Divergence) to flag statistically abnormal reasoning outputs. Add Reticle-style session recording to capture full JSON-RPC payloads for replay debugging.

## 12. Model Intelligence & Cost Optimization
**The Problem:** Single-model routing is inefficient; 10x cost variations exist.
**The Converging Solution:** Complexity-based tiered routing.
*   **Borg Implementation:** The `WaterfallRouter` implements tiered routing: Simple URLs (Reddit, YCombinator) route to a local 1.2B SLM, and complex tasks (GitHub repos) route to frontier models. Track model performance ELO and auto-demote degrading models.

## 13. Browser Automation & Web Interaction
**The Problem:** Naive HTTP fetching fails on SPAs, captchas, and auth walls.
**The Converging Solution:** Accessibility Tree extraction (Generation 2) and Zero-Copy Vision (Generation 3).
*   **Borg Implementation:** Use Accessibility Trees for SPA/Docs to achieve 90% token reduction vs raw DOM. Fall back to Zero-Copy POSIX shared memory vision for complex diagrams.

## 14. CI/CD & Git-Integrated Workflows
**The Problem:** Moving from code generation to code *verification* as the primary bottleneck.
**The Converging Solution:** Agent-native git workflows, worktree isolation, and review-as-quality-gates.
*   **Borg Implementation:** Version the database schema using git-commit-like operations. Implement closed-loop self-healing (if extraction quality is poor, automatically checkout a new worktree and re-process with a different model).

### 1. Identity & Vision: The AI Hypervisor
The project has evolved from its origin as "Borg" into a dual-brand architectural vision:
- **Nexus:** The underlying coordination kernel or "AI Hypervisor." It manages active memory, tool routing, and orchestration.
- **HyperCode:** The flagship, autonomous developer-facing coding product powered by the Nexus kernel.

The "AI Hypervisor" model treats AI models as compute resources and tools as peripheral drivers, with Nexus acting as the management layer that optimizes model selection, context management, and execution loops.

### 2. Current State (v1.0.0-alpha.56)
The project is currently in the transition between **Phase 1 (Active Memory)** and **Phase 2 (Autonomy Loop)**.
- **Phase 1 Status:** Complete. The foundational memory substrate is production-ready.
- **Phase 2 Status:** Initiated. Focus has shifted to self-healing reactors and the "execute-fix-verify-retry" autonomous loop.

### 3. Core Architectural Patterns
- **Kernel/Control Plane Split:**
    - **Kernel:** Deterministic execution, memory, and routing (being migrated toward `go/` and `@nexus/kernel`).
    - **Control Plane:** Dashboards, session management, and operator UI (`apps/web`, `packages/core`).
- **Active Memory Substrate:**
    - **Heat-Based Tiering:** Entries have a `heat_score` (0-100). Utility increases heat; time causes exponential decay (24h half-life).
    - **Feedback Loops:** Tool success/failure directly modifies the heat of the context used to achieve that outcome.
- **Provider Routing:**
    - Uses a waterfall fallback system. If one model/provider quota is exhausted, it automatically falls back to the next best available resource.
- **Progressive Tool Disclosure:**
    - Instead of flooding context with all tools, Nexus uses semantic ranking to disclose only relevant tools based on the active goal.

### 4. Monorepo Structure & Module Roles
- **`packages/core`:** The central hub ("Brain") of the TypeScript control plane. It hosts tRPC routers, session logic, and bridges to the Go sidecar.
- **`packages/memory`:** The implementation layer for LanceDB and vector-based storage.
- **`go/`:** The Go Sidecar (Port 4300). Currently serves as a high-performance state authority and BM25 ranking engine, mirroring and bridging TypeScript services.
- **`apps/web`:** The primary operator dashboard for managing sessions and visualizing the knowledge graph.
- **`packages/tools`:** Contains functional tool implementations (Read, Write, Shell, etc.) shared across CLI and Web surfaces.

### 5. Technical Decisions & Constraints
- **Shell Hardening:** `child_process.exec` is strictly prohibited. All command execution must use `spawn` with tokenized argument arrays and `shell: false`.
- **Environment:** Standardized on Node 24 and Go 1.24.3. Port 443 is restricted; local caches/binaries must be used for dependency management.
- **Version Authority:** Versioning is synchronized globally. The current baseline is `1.0.0-alpha.56`.

### 6. Roadmap: The Autonomy Path
The next immediate milestones involve:
1.  **The Healer Loop:** Implementing the full `execute-fix-verify-retry` autonomous cycle within the `HealerReactor`.
2.  **Fleet Management:** Extending Nexus to manage multiple concurrent "HyperCode" sessions with shared organizational memory.
3.  **Assimilation:** Systematically migrating high-performance logic (ranking, sync, memory) from TypeScript into the native Go kernel.

### 7. Governance & Intelligence
- **Supervisor/Council Pattern:** The system is designed to run under the supervision of an "Architect" or "Council" of models that verify plans before implementation.
- **Passive Harvesting:** The system automatically extracts facts and patterns from agent traffic to populate its L2 "Vault" memory without manual operator intervention.

---
*Last updated: Session v1.0.0-alpha.56*
# AI Hypervisor (Nexus) - Comprehensive Architectural Memory

This document summarizes the foundational architecture, established patterns, and strategic decisions of the project as of version **1.0.0-alpha.56**.

## 1. Strategic Identity: Nexus & HyperCode
The project has successfully pivoted from "Borg" to a dual-brand infrastructure model:
*   **Nexus (The Kernel/Hypervisor):** The underlying coordination runtime and "AI Hypervisor." It treats LLMs as "guest operating systems" and manages the low-level memory, routing, and execution buses.
*   **HyperCode (The Product):** The user-facing, local-first autonomous coding environment powered by the Nexus kernel.

## 2. Active Tiered Memory Substrate (Implemented Phase 1)
*   **Heat Scoring (0-100):** Every memory entry tracks utility. Heat increases on access and decays exponentially (24-hour half-life).
*   **Tiered Hierarchy:** L1 (Working Memory, heat > 80) is promoted to context; L2 (Vault) is for semantic recall.
*   **Tool-Outcome Feedback:** `MemoryManager.recordToolOutcome()` boosts the heat of successful patterns and demotes failures.

## 3. "Kernel / Control Plane" Topology
*   **/kernel**: Deterministic brain (runtime, memory, router).
*   **/control-plane**: Observer layer (UI, Telemetry).

## 4. State Authority & The Sidecar Pattern
*   **Go Sidecar (Port 4300):** State authority and BM25 ranking.
*   **TS Bridge (Port 4100):** Primary control-plane bridge and tRPC host.

## 5. Intelligence Management: Progressive Disclosure
*   **Decision System:** Ranked discovery and LRU eviction ensures only 3-5 tools/skills are in the active working set.

## 6. Hardened Execution & Security
*   **Standard:** Tokenized argument arrays with `shell: false` for all command executions.
*   **Parity Principle:** 1:1 behavioral and schema parity for tools expected by proprietary models (e.g., Claude Code).

---
### Meta-Protocol for Future Sessions
1.  **Truth over Fiction:** Dashboards must reflect real state.
2.  **Autonomous Momentum:** Proceed through Phase 2 (Autonomy/Self-Healing).
3.  **Documentation Sync:** Every version bump syncs all meta files and manifests.
