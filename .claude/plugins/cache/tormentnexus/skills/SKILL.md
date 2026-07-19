---
name: tormentnexus
description: TormentNexus AI control plane integration for Claude Code CLI
version: 1.0.0
---

# TormentNexus Integration

TormentNexus is a local AI control plane (port 7778) with L2 memory, tool discovery, session import, skill registry, and code search. The MCP server must be configured in `~/.claude/claude.json`.

## MCP Server

Add this entry to `~/.claude/claude.json`:
```json
{
  "mcpServers": {
    "tormentnexus": {
      "command": "C:\\Users\\hyper\\workspace\\tormentnexus\\tormentnexus.exe",
      "args": ["mcp"]
    }
  }
}
```

## MCP Tools

- `mcp_tormentnexus_memory_scratchpad_*` — L1 working memory
- `mcp_tormentnexus_mcp_list_tools/servers` — discovery
- `mcp_tormentnexus_mcp_call_tool` — route to 20+ servers
- `mcp_tormentnexus_bash/read/write/edit/grep/find/ls` — system tools

## Usage

1. Check scratchpad before complex tasks
2. Store key decisions after completion
3. Use repomap for codebase orientation
4. Route commercial integrations through TN Kernel
