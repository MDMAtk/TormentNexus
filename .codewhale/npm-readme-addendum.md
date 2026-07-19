
## TormentNexus Extension

CodeWhale 0.8.66+ includes a built-in **TormentNexus extension** (`crates/tn-extension`) compiled directly into the binary — no separate plugin or DLL needed. It provides:

### Lifecycle Hooks
- **SessionStart** — logs session info to TN L2 memory on port 7778
- **BeforeAgentStart** — injects TN system prompt guidance + searches L2 for relevant context
- **ToolCall** — logs tool calls; checks 6 dangerous patterns against TN commercial RBAC
- **ToolResult** — auto-stores substantial results from key tools to L2
- **TurnEnd** — logs tool usage per turn
- **Input** — expands `@memory:key` inline with L2 content
- **UserBash** — audit-logs shell commands to TN commercial audit
- **ModelSelect** — tracks model changes
- **SessionCompact** — preserves memory across compaction

### Custom Tools (via Extension API)
- 9 tool definitions: `tn_memory_store`, `tn_memory_search`, `tn_memory_vector_search`, `tn_tool_search`, `tn_session_search`, `tn_skill_manage`, `tn_code_search`, `tn_context_harvest`, `tn_scratchpad`
- 6 slash commands: `/tn-store`, `/tn-search`, `/tn-status`, `/tn-plan`, `/tn-summary`, `/tn-purge`
- 3 keyboard shortcuts: `Ctrl+Shift+M` (memory), `Ctrl+Shift+T` (tools), `Ctrl+Shift+P` (status)

### MCP Server Auto-Registration
The extension automatically registers the `tormentnexus` MCP server pointing at `tormentnexus.exe mcp`. Enable it by keeping the TormentNexus Kernel running on `http://127.0.0.1:7778`.

> Extension source: `crates/tn-extension` — implements `codewhale_extension::Extension` trait.
