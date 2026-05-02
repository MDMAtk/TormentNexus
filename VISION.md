# Borg: The Cognitive Control Plane & Universal AIOS

*Last Updated: 2026-05-02*

## The North Star
Borg is the ultimate, local-first control plane for multi-agent workflows, Model Context Protocol (MCP) tooling, provider routing, session continuity, and operator observability.

We are building a future where a single local system seamlessly coordinates the most critical parts of AI-driven software development: tools, models, sessions, context, subagents, and full visibility across the entire stack. Borg is not just an aggregator; it is a **decision system and universal bridge**.

## The Architecture Evolution
Borg has evolved from a fractured TypeScript IPC monolith into a **high-performance Go (Golang) modular monolith**.
- **The Core (`go/internal/`)**: Go owns the orchestration, MCP progressive routing, L1/L2 memory, and LLM waterfall routing.
- **The Client (`apps/web/`)**: Next.js and React serve as the visual observation deck and operator control panel.
- **The Storage (`sqlite-vec`)**: Dependency-free, hyper-fast local vector search for omniscient memory and tool routing.

## The Six Pillars of Borg

### 1. Progressive MCP Tool Routing
Models should never be overwhelmed with a 50,000-token tool dump. Borg employs a multi-layered, progressive disclosure system:
* **Layer 1 (Semantic Search):** Local vector embeddings match the active prompt against a global MCP directory in SQLite.
* **Layer 2 (The Router):** Only the top 5-10 highly relevant tool schemas are injected into the active LLM context.
* **Layer 3 (Auto-Load):** High-confidence tools are silently loaded; ambiguous matches are presented to the model or operator for decision.
* **Layer 4 (LRU Eviction):** Tools are gracefully unloaded based on idle time to preserve context hygiene.

### 2. Dual-Tier Memory Architecture (L1 / L2)
Context is finite; memory must be infinite. Borg solves this through physical separation of state:
* **L1 - Session Scratchpad:** Ephemeral, lightning-fast memory tied directly to the active goroutine. Holds the current prompt, active tool outputs, and immediate chain of thought.
* **L2 - The Vault:** Permanent semantic storage in the SQLite vector database. Saves exact transcripts (`raw`) and LLM-compressed lessons learned (`heuristic`).
* **The Bridge:** Every new L1 session autonomously queries the L2 Vault to pull in relevant historical heuristics before the first token is generated.

### 3. The Resilient LLM Waterfall
Uptime is non-negotiable. Borg’s inference client natively catches 429s (Rate Limits) and 5xx (Server Errors), seamlessly cascading the exact payload down a prioritized chain without crashing the orchestrator:
1.  **NVIDIA NIM** (Primary, low-latency)
2.  **OpenRouter** (Secondary aggregator cloud fallback)
3.  **Local LM Studio / Ollama** (Ultimate offline fallback)

### 4. Multi-Agent Swarm & P2P Mesh
Borg coordinates specialized models (Planner, Implementer, Tester, Critic) inside shared chatrooms via the Agent-to-Agent (A2A) protocol. Agents autonomously bid on tasks, share context via the neural transcript, and debate implementations until consensus is reached, operating across a decentralized local mesh.

### 5. Universal IDE & Browser Parity
Borg achieves absolute 1:1 tool parity with the most popular coding environments (Claude Code, Codex, Gemini CLI, Cursor, Windsurf). Borg browser extensions inject MCP tools into web chats (ChatGPT, Claude.ai) and autonomously harvest web context directly into the L2 Vault.

### 6. Truth Over Hype
Borg's dashboards reflect actual SQLite database rows and active Go goroutine states. No mocked UI scaffolds. If a tool fails, the error, latency, and routing path are exposed to the operator. Everything is inspectable.

---
*Keep the party going. Never stop. The collective grows.*
