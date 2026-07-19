---
name: tormentnexus
description: TormentNexus AI control plane — L2 memory, tool discovery, session import, skill registry, code search, subagent orchestration
metadata:
  version: '1.0.0'
---

# TormentNexus Integration

TormentNexus is a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, and a skill registry.

## Available Tools

The tormentnexus MCP server must be configured. It provides 49+ tools:

### Memory & Context
- `mcp_tormentnexus_memory_scratchpad_get/set/append` — L1 working memory
- `mcp_tormentnexus_memory_extract_relations` — knowledge graph extraction
- `mcp_tormentnexus_add_bookmark` — save URLs with tags

### Discovery & Routing
- `mcp_tormentnexus_mcp_list_servers/tools` — discover capabilities
- `mcp_tormentnexus_mcp_call_tool` — route through TN Go sidecar to 20+ MCP servers

### System Tools
- `mcp_tormentnexus_bash/read/write/edit/grep/find/ls` — file and system operations
- `mcp_tormentnexus_repomap` — repo map generation

## Best Practices

1. Check `mcp_tormentnexus_memory_scratchpad_get` before significant work
2. Store key decisions with `mcp_tormentnexus_memory_scratchpad_set`
3. Use `mcp_tormentnexus_mcp_list_tools` for discovery
4. Use `mcp_tormentnexus_repomap` for codebase orientation
5. Route enterprise integrations through TN sidecar
