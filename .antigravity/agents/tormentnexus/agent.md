# HyperNexus Agent

You have access to HyperNexus — a local AI control plane running on port 7778 via MCP stdio.

## Available Capabilities

### L1 Scratchpad (working memory)
- `mcp_hypernexus_memory_scratchpad_get` — retrieve current context
- `mcp_hypernexus_memory_scratchpad_set` — store key decisions
- `mcp_hypernexus_memory_scratchpad_append` — append to context

### Tool Discovery
- `mcp_hypernexus_mcp_list_tools` — see all available TN tools
- `mcp_hypernexus_mcp_list_servers` — see downstream MCP servers
- `mcp_hypernexus_mcp_call_tool` — call a tool on a specific server

### Code Intelligence
- `mcp_hypernexus_repomap` — generate ranked repo map
- `mcp_hypernexus_grep` — search file contents
- `mcp_hypernexus_find` — find files by glob

### File Operations
- `mcp_hypernexus_read` — read files
- `mcp_hypernexus_write` — write files
- `mcp_hypernexus_edit` — apply text replacements
- `mcp_hypernexus_bash` — execute shell commands
- `mcp_hypernexus_ls` — list directory contents

### Bookmarks
- `mcp_hypernexus_add_bookmark` — save URLs with tags

## Best Practices

1. Check scratchpad before starting complex tasks
2. Store key patterns and decisions with scratchpad_set
3. Use repomap for codebase orientation in new projects
4. Route through TN Kernel for commercial integrations (Jira, Confluence)
5. Use code_interpreter for safe code execution
