---
name: hypernexus
description: HyperNexus AI control plane integration for Mavis CLI
version: 1.0.0
---

# HyperNexus Integration for Mavis

HyperNexus is a local AI control plane (port 7778) with L2 memory, tool discovery, session import, skill registry, and code search. Connected via MCP stdio.

## MCP Tools

The hypernexus MCP server exposes tools prefixed `mcp_hypernexus_*`.

### Memory & Context (L1/L2)
- `mcp_hypernexus_memory_scratchpad_*` — L1 working memory
- `mcp_hypernexus_memory_extract_relations` — knowledge graph
- `mcp_hypernexus_add_bookmark` — save links with tags

### Tool Discovery
- `mcp_hypernexus_mcp_list_servers/tools` — browse capabilities
- `mcp_hypernexus_mcp_call_tool` — route to 20+ MCP servers
- `mcp_hypernexus_mcp_status` — runtime health

### File & System
- `mcp_hypernexus_read/write/edit/bash/grep/find/ls` — core tools
- `mcp_hypernexus_repomap` — codebase map
- `mcp_hypernexus_system_status` — system info
- `mcp_hypernexus_code_interpreter` — code execution

## Usage

1. **Before work**: Check `mcp_hypernexus_memory_scratchpad_get` for context
2. **During work**: Use `mcp_hypernexus_repomap` for orientation
3. **After work**: Store with `mcp_hypernexus_memory_scratchpad_set`
4. **Discovery**: `mcp_hypernexus_mcp_list_tools` to find capabilities
