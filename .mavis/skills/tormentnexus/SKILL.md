---
name: tormentnexus
description: TormentNexus AI control plane integration for Mavis CLI
version: 1.0.0
---

# TormentNexus Integration for Mavis

TormentNexus is a local AI control plane (port 7778) with L2 memory, tool discovery, session import, skill registry, and code search. Connected via MCP stdio.

## MCP Tools

The tormentnexus MCP server exposes tools prefixed `mcp_tormentnexus_*`.

### Memory & Context (L1/L2)
- `mcp_tormentnexus_memory_scratchpad_*` — L1 working memory
- `mcp_tormentnexus_memory_extract_relations` — knowledge graph
- `mcp_tormentnexus_add_bookmark` — save links with tags

### Tool Discovery
- `mcp_tormentnexus_mcp_list_servers/tools` — browse capabilities
- `mcp_tormentnexus_mcp_call_tool` — route to 20+ MCP servers
- `mcp_tormentnexus_mcp_status` — runtime health

### File & System
- `mcp_tormentnexus_read/write/edit/bash/grep/find/ls` — core tools
- `mcp_tormentnexus_repomap` — codebase map
- `mcp_tormentnexus_system_status` — system info
- `mcp_tormentnexus_code_interpreter` — code execution

## Usage

1. **Before work**: Check `mcp_tormentnexus_memory_scratchpad_get` for context
2. **During work**: Use `mcp_tormentnexus_repomap` for orientation
3. **After work**: Store with `mcp_tormentnexus_memory_scratchpad_set`
4. **Discovery**: `mcp_tormentnexus_mcp_list_tools` to find capabilities
