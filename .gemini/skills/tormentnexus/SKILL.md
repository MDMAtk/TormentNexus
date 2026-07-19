---
name: tormentnexus
description: TormentNexus AI control plane — L2 memory, tool search, session import, skill registry, code search, subagent orchestration
version: 1.0.0
---

# TormentNexus Integration

TormentNexus is a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, and a skill registry.

## Available via MCP

The tormentnexus MCP server is configured and exposes tools via the `mcp_tormentnexus_*` prefix through the MCP system.

### Memory Tools
- `mcp_tormentnexus_memory_scratchpad_get/set/append` — L1 working memory
- `mcp_tormentnexus_memory_extract_relations` — graph extraction
- `mcp_tormentnexus_add_bookmark` — save URLs with tags

### Discovery Tools
- `mcp_tormentnexus_mcp_call_tool` — route through TN's TN kernel to 20+ MCP servers
- `mcp_tormentnexus_mcp_list_tools` — discover available MCP tools
- `mcp_tormentnexus_mcp_list_servers` — list connected MCP servers

### System Tools
- `mcp_tormentnexus_bash` — shell execution
- `mcp_tormentnexus_read/write/edit` — file I/O
- `mcp_tormentnexus_grep/find/ls` — search
- `mcp_tormentnexus_repomap` — repo map generation

### Integration Tools
- `mcp_tormentnexus_jira_create_issue` — create Jira issues
- `mcp_tormentnexus_confluence_search` — search Confluence
- `mcp_tormentnexus_code_interpreter` — execute code
- `mcp_tormentnexus_system_status` — system health

## Best Practices

1. **Before significant work**: Use `mcp_tormentnexus_memory_scratchpad_get` to check for relevant context
2. **During development**: Use `mcp_tormentnexus_repomap` for codebase orientation
3. **After decisions**: Store key decisions via `mcp_tormentnexus_memory_scratchpad_set`
4. **Tool discovery**: Use `mcp_tormentnexus_mcp_list_tools` when unsure what's available
5. **Complex tasks**: Route through TN Kernel via `mcp_tormentnexus_mcp_call_tool` for deep context
