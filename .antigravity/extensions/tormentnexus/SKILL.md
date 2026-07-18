---
name: hypernexus
description: HyperNexus AI control plane integration
version: 1.0.0
---

# HyperNexus Integration

HyperNexus is a local AI control plane on port 7778 with L2 memory, tool discovery, session import, skill registry, and code search.

## MCP Tools Available

All TN tools are available through the `mcp_hypernexus_*` namespace:

### Memory & Context
- **mcp_hypernexus_memory_scratchpad_*** — L1 scratchpad (get/set/append)
- **mcp_hypernexus_memory_extract_relations** — knowledge graph extraction
- **mcp_hypernexus_add_bookmark** — bookmark storage with tags

### Discovery & Routing
- **mcp_hypernexus_mcp_list_servers/tools** — discover capabilities
- **mcp_hypernexus_mcp_call_tool** — route to downstream MCP servers
- **mcp_hypernexus_mcp_status** — check TN runtime health

### System
- **mcp_hypernexus_bash** — shell execution
- **mcp_hypernexus_read/write/edit/grep/find/ls** — file operations
- **mcp_hypernexus_repomap** — repo map generation

### Integrations
- **mcp_hypernexus_system_status** — system health
- **mcp_hypernexus_code_interpreter** — code execution
- **mcp_hypernexus_install_mcp_server** — install new MCP servers

## Workflow

1. At task start: check scratchpad for existing context
2. During work: use repomap for orientation, grep for search
3. When stuck: list tools for discovery, route through TN Kernel
4. After decisions: persist to scratchpad for cross-session recall
