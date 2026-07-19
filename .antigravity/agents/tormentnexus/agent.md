# TormentNexus Agent

You have access to TormentNexus — a local AI control plane running on port 7778 via MCP stdio.

## Available Capabilities

### L1 Scratchpad (working memory)
- `mcp_tormentnexus_memory_scratchpad_get` — retrieve current context
- `mcp_tormentnexus_memory_scratchpad_set` — store key decisions
- `mcp_tormentnexus_memory_scratchpad_append` — append to context

### Tool Discovery
- `mcp_tormentnexus_mcp_list_tools` — see all available TN tools
- `mcp_tormentnexus_mcp_list_servers` — see downstream MCP servers
- `mcp_tormentnexus_mcp_call_tool` — call a tool on a specific server

### Code Intelligence
- `mcp_tormentnexus_repomap` — generate ranked repo map
- `mcp_tormentnexus_grep` — search file contents
- `mcp_tormentnexus_find` — find files by glob

### File Operations
- `mcp_tormentnexus_read` — read files
- `mcp_tormentnexus_write` — write files
- `mcp_tormentnexus_edit` — apply text replacements
- `mcp_tormentnexus_bash` — execute shell commands
- `mcp_tormentnexus_ls` — list directory contents

### Bookmarks
- `mcp_tormentnexus_add_bookmark` — save URLs with tags

## Best Practices

1. Check scratchpad before starting complex tasks
2. Store key patterns and decisions with scratchpad_set
3. Use repomap for codebase orientation in new projects
4. Route through TN Kernel for commercial integrations (Jira, Confluence)
5. Use code_interpreter for safe code execution
