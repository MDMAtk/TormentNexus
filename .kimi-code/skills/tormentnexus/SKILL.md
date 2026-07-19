---
name: tormentnexus
description: TormentNexus AI control plane — L2 memory, tool discovery, session import, skill registry, code search
---

# TormentNexus Integration

TormentNexus is a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, and a skill registry.

## MCP Server

Configure the tormentnexus MCP server:
- **Command**: `C:\Users\hyper\workspace\tormentnexus\tormentnexus.exe`
- **Args**: `["mcp"]`

## Tools
- `mcp_tormentnexus_memory_scratchpad_*` — L1 working memory
- `mcp_tormentnexus_mcp_list_tools/servers` — discovery
- `mcp_tormentnexus_bash/read/write/edit/grep/find/ls` — system tools

## Best Practices
1. Check scratchpad before significant work
2. Store decisions after key moments
3. Use repomap for codebase orientation
