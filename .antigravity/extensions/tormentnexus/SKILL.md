---
name: tormentnexus
description: TormentNexus AI control plane integration
version: 1.0.0
---

# TormentNexus Integration

TormentNexus is a local AI control plane on port 7778 with L2 memory, tool discovery, session import, skill registry, and code search.

## MCP Tools Available

All TN tools are available through the `mcp_tormentnexus_*` namespace:

### Memory & Context
- **mcp_tormentnexus_memory_scratchpad_*** — L1 scratchpad (get/set/append)
- **mcp_tormentnexus_memory_extract_relations** — knowledge graph extraction
- **mcp_tormentnexus_add_bookmark** — bookmark storage with tags

### Discovery & Routing
- **mcp_tormentnexus_mcp_list_servers/tools** — discover capabilities
- **mcp_tormentnexus_mcp_call_tool** — route to downstream MCP servers
- **mcp_tormentnexus_mcp_status** — check TN runtime health

### System
- **mcp_tormentnexus_bash** — shell execution
- **mcp_tormentnexus_read/write/edit/grep/find/ls** — file operations
- **mcp_tormentnexus_repomap** — repo map generation

### Integrations
- **mcp_tormentnexus_system_status** — system health
- **mcp_tormentnexus_code_interpreter** — code execution
- **mcp_tormentnexus_install_mcp_server** — install new MCP servers

## Workflow

1. At task start: check scratchpad for existing context
2. During work: use repomap for orientation, grep for search
3. When stuck: list tools for discovery, route through TN Kernel
4. After decisions: persist to scratchpad for cross-session recall
