## Summary

Adds `crates/tn-extension` — a native Rust extension implementing the `codewhale_extension::Extension` trait with full parity to the TormentNexus Pi Coding Agent extension for persistent L2 memory, MCP tool discovery, skill registry, code search, RBAC, and context harvesting.

## Changes

**7 files** (+387/-1 lines):

### New crate: `crates/tn-extension/`
- `Cargo.toml` — depends on `codewhale-extension`, `serde_json`, `reqwest`, `chrono`, `async-trait`
- `src/lib.rs` — `TormentNexusExtension` implementing all lifecycle hooks

### Extension Lifecycle Hooks
- **SessionStart** — logs session to TN L2 memory at `http://127.0.0.1:7778/api/memory/add`
- **BeforeAgentStart** — injects TN system prompt guidance + searches L2 for relevant context per turn
- **ToolCall** — logs tool calls; checks 6 dangerous patterns (`rm -rf`, `sudo`, `DROP TABLE`, etc.) against `POST /api/commercial/authorize`
- **ToolResult** — auto-stores substantial results (>=100 chars) from 6 key tools to L2
- **TurnEnd** — logs tool usage summary per turn
- **Input** — expands `@memory:key` inline with L2 content from TN Kernel
- **UserBash** — audit-logs shell commands to TN commercial audit
- **ModelSelect** — tracks model changes to L2
- **SessionBeforeCompact** / **SessionCompact** — preserves memory across compaction boundaries

### Registration via Extension Trait
- **9 custom tool definitions**: `tn_memory_store`, `tn_memory_search`, `tn_memory_vector_search`, `tn_tool_search`, `tn_session_search`, `tn_skill_manage`, `tn_code_search`, `tn_context_harvest`, `tn_scratchpad`
- **6 slash commands**: `/tn-store`, `/tn-search`, `/tn-status`, `/tn-plan`, `/tn-summary`, `/tn-purge`
- **3 keyboard shortcuts**: `Ctrl+Shift+M` (memory search), `Ctrl+Shift+T` (tool search), `Ctrl+Shift+P` (system status)
- **MCP server auto-registration**: points at `tormentnexus.exe mcp` with env `TORMENTNEXUS_WORKSPACE_ROOT`

### Wiring
- `crates/tui/Cargo.toml` — added `codewhale-tn-extension` dependency
- `crates/tui/src/core/engine.rs` — registers `TormentNexusExtension` into `ExtensionManager` at startup
- `Cargo.toml` (workspace root) — added `crates/tn-extension` as workspace member

## Dependencies
- `reqwest` (HTTP client for TN Kernel API)
- `serde_json` (JSON serialization for API payloads)
- `chrono` (RFC3339 timestamps)
- `async-trait` (async trait support for Extension)
- All are already workspace dependencies; no new crate registry entries.
