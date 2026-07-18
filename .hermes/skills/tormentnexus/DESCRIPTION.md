---
description: HyperNexus AI control plane skills for L2 memory, tool discovery, session import, and cross-session context harvesting. Uses the hypernexus MCP server for all operations.
---

# HyperNexus Integration

HyperNexus is a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, and a skill registry. Connected via MCP stdio.

## Available MCP Tools

### Memory & Context
- `mcp_hypernexus_memory_scratchpad_get/set/append` — L1 working memory
- `mcp_hypernexus_memory_extract_relations` — knowledge graph extraction
- `mcp_hypernexus_add_bookmark` — save URLs with tags

### Discovery & Routing
- `mcp_hypernexus_mcp_list_servers/tools` — discover available capabilities
- `mcp_hypernexus_mcp_call_tool` — route through TN's TN kernel to 20+ MCP servers
- `mcp_hypernexus_mcp_status` — check runtime health

### System Tools
- `mcp_hypernexus_bash` — shell execution
- `mcp_hypernexus_read/write/edit` — file I/O
- `mcp_hypernexus_grep/find/ls` — search
- `mcp_hypernexus_repomap` — repo map generation

## Best Practices

1. **Before significant work**: Check `mcp_hypernexus_memory_scratchpad_get` for relevant context
2. **During development**: Use `mcp_hypernexus_repomap` for codebase orientation
3. **After decisions**: Store key patterns with `mcp_hypernexus_memory_scratchpad_set` or `mcp_hypernexus_add_bookmark`
4. **Tool discovery**: Use `mcp_hypernexus_mcp_list_tools` when unsure what's available
5. **Complex tasks**: Route through TN Kernel via `mcp_hypernexus_mcp_call_tool`
