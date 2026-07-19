---
name: "tormentnexus"
description: "TormentNexus AI control plane — persistent L2 vector memory, semantic tool discovery, imported sessions, skill registry, code search, subagent orchestration. Use when the task benefits from cross-session memory, tool discovery across 20+ MCP servers, or access to imported sessions and skill modules."
---

# TormentNexus Integration

TormentNexus is a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, and a skill registry. Connected via MCP stdio.

## MCP Server

The tormentnexus MCP server is configured and exposes tools via the `mcp_tormentnexus_*` namespace.

### Memory & Context
- `mcp_tormentnexus_memory_scratchpad_get/set/append` — L1 working memory
- `mcp_tormentnexus_memory_extract_relations` — knowledge graph extraction
- `mcp_tormentnexus_add_bookmark` — save URLs with tags

### Discovery & Routing
- `mcp_tormentnexus_mcp_list_servers/tools` — discover available capabilities
- `mcp_tormentnexus_mcp_call_tool` — route through TN's TN kernel to 20+ MCP servers
- `mcp_tormentnexus_mcp_status` — check runtime health

### System Tools
- `mcp_tormentnexus_bash` — shell execution
- `mcp_tormentnexus_read/write/edit` — file I/O
- `mcp_tormentnexus_grep/find/ls` — search
- `mcp_tormentnexus_repomap` — repo map generation

### Integration Tools
- `mcp_tormentnexus_system_status` — system health
- `mcp_tormentnexus_code_interpreter` — code execution
- `mcp_tormentnexus_jira_create_issue` — create Jira issues
- `mcp_tormentnexus_confluence_search` — search Confluence

### Commercial Tools
- `mcp_tormentnexus_cloud_troubleshoot` — cloud infrastructure diagnostics
- `mcp_tormentnexus_generate_devops_pipeline` — CI/CD generation
- `mcp_tormentnexus_install_mcp_server` — install new MCP servers

## Best Practices

1. **Before significant work**: Check `mcp_tormentnexus_memory_scratchpad_get` for relevant context
2. **During development**: Use `mcp_tormentnexus_repomap` for codebase orientation
3. **After decisions**: Store key patterns with `mcp_tormentnexus_memory_scratchpad_set`
4. **Tool discovery**: Use `mcp_tormentnexus_mcp_list_tools` when unsure what's available
5. **Complex tasks**: Route through TN Kernel via `mcp_tormentnexus_mcp_call_tool`
6. **Cross-session**: Use scratchpad tools to persist context between sessions

## Rules

- Always check scratchpad memory before starting complex multi-step tasks
- Store important decisions and patterns after completing significant work
- Use repomap for codebase orientation in unfamiliar projects
- Route commercial integrations (Jira, Confluence) through TN Kernel
