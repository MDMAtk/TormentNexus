---
name: tormentnexus
description: TormentNexus AI control plane — L2 vector memory, semantic tool discovery, skill registry, session import, code search, subagent orchestration, and commercial RBAC
metadata:
  short-description: TormentNexus AI control plane
---

# TormentNexus Integration

TormentNexus is a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, a skill registry, commercial RBAC, and subagent orchestration.

CodeWhale connects to it via MCP stdio (`tormentnexus.exe mcp`). The MCP server is auto-configured in `~/.codewhale/mcp.json` by the TN installer.

## Architecture Overview

```
CodeWhale ⇄ MCP (tormentnexus.exe mcp) ⇄ TN Kernel (port 7778)
                                     ⇄ 20+ downstream MCP servers
                                     ⇄ SQLite L2 vector memory
                                     ⇄ Session import store (542+ sessions)
                                     ⇄ Skill registry (5,776+ modules)
                                     ⇄ Commercial RBAC
```

The Kernel runs as the "TormentNexusKernel" Windows service on port 7778.

---

## MCP Tools

The `tormentnexus` MCP server provides 49+ tools. Tools are prefixed `mcp_tormentnexus_*`:

### Memory & Knowledge
| Tool | Description |
|------|-------------|
| `mcp_tormentnexus_memory_scratchpad_get` | Read L1 ephemeral working memory |
| `mcp_tormentnexus_memory_scratchpad_set` | Write to L1 scratchpad |
| `mcp_tormentnexus_memory_scratchpad_append` | Append to L1 scratchpad |
| `mcp_tormentnexus_memory_extract_relations` | Extract knowledge graph relations |
| `mcp_tormentnexus_add_bookmark` | Store a bookmark (lighter than full memory) |

### File & System
| Tool | Description |
|------|-------------|
| `mcp_tormentnexus_read` | Read file contents |
| `mcp_tormentnexus_write` | Write file contents |
| `mcp_tormentnexus_edit` | Edit file |
| `mcp_tormentnexus_grep` | Search file contents |
| `mcp_tormentnexus_find` | Find files |
| `mcp_tormentnexus_ls` | List directory |
| `mcp_tormentnexus_bash` | Execute shell command (sandboxed) |
| `mcp_tormentnexus_repomap` | Generate repository map |

### MCP Routing (via TN Kernel)
| Tool | Description |
|------|-------------|
| `mcp_tormentnexus_mcp_list_servers` | List all registered MCP servers |
| `mcp_tormentnexus_mcp_list_tools` | List tools from a specific server |
| `mcp_tormentnexus_mcp_call_tool` | Call a tool on a downstream server |
| `mcp_tormentnexus_mcp_server_test` | Check server health |
| `mcp_tormentnexus_mcp_status` | Overall MCP mesh status |

### Windows & UI Automation
| Tool | Description |
|------|-------------|
| `mcp_tormentnexus_inspect_window_ui` | Inspect accessibility tree |
| `mcp_tormentnexus_simulate_input` | Simulate keyboard/mouse input |
| `mcp_tormentnexus_list_processes` | List running processes |
| `mcp_tormentnexus_kill_process` | Kill a process |
| `mcp_tormentnexus_detect_chat_surface` | Detect AI chat windows |
| `mcp_tormentnexus_detect_chat_state` | Check chat state |
| `mcp_tormentnexus_click_chat_button` | Click chat UI buttons |
| `mcp_tormentnexus_set_chat_input` | Set chat input field |
| `mcp_tormentnexus_submit_chat_input` | Submit chat input |
| `mcp_tormentnexus_advance_chat` | Advance chat session |

### Commercial & Integrations
| Tool | Description |
|------|-------------|
| `mcp_tormentnexus_jira_create_issue` | Create a Jira issue |
| `mcp_tormentnexus_confluence_search` | Search Confluence |
| `mcp_tormentnexus_cloud_troubleshoot` | Cloud troubleshooting |
| `mcp_tormentnexus_generate_devops_pipeline` | Generate CI/CD pipeline |
| `mcp_tormentnexus_install_mcp_server` | Install new MCP servers |

### System
| Tool | Description |
|------|-------------|
| `mcp_tormentnexus_code_interpreter` | Run code in sandbox |
| `mcp_tormentnexus_download_llamafile` | Download a llamafile |
| `mcp_tormentnexus_system_status` | Show system status |
| `mcp_tormentnexus_billing_status` | Show billing info |
| `mcp_tormentnexus_get_system_stats` | Get system statistics |

---

## 🧠 Memory System

TN has three memory tiers:

| Tier | Name | Scope | Persistence |
|------|------|-------|-------------|
| L1 | Scratchpad | Current session | Ephemeral |
| L2 | Vector Vault | Cross-session | Persistent (SQLite + embeddings) |
| L3 | Cold Archive | Long-term | Archived |

### L2 Memory Operations (via TN Kernel API)

For operations not exposed through MCP tools, use the TN kernel REST API at `http://127.0.0.1:7778`:

**Store a memory:**
```bash
curl -X POST http://127.0.0.1:7778/api/memory/add \
  -H "Content-Type: application/json" \
  -d '{"content": "{\"content\": \"...\", \"tags\": [\"project:foo\", \"pattern:build\"], \"category\": \"decision\"}"}'
```

**Search memories:**
```
GET http://127.0.0.1:7778/api/memory/search?q=<query>&limit=<n>
```

**List all memories:**
```
GET http://127.0.0.1:7778/api/memory/list
```

**FTS search (full count):**
```
GET http://127.0.0.1:7778/api/memory/fts-search?q=<query>&limit=<n>
```

**Cold archive:**
```
GET http://127.0.0.1:7778/api/memory/cold-archive/count
```

**Mesh status:**
```
GET http://127.0.0.1:7778/api/mesh/status
```

**Tool search (semantic MCP tool discovery):**
```
GET http://127.0.0.1:7778/api/mcp/native/search?query=<natural language>
```

**Skill search:**
```
GET http://127.0.0.1:7778/api/skills/search?q=<query>
```

**Session search:**
```
GET http://127.0.0.1:7778/api/memory/search?q=<query>&type=session
```

**Project memdb sync:**
```
POST http://127.0.0.1:7778/api/memory/project/sync
```

**Code search:**
```
GET http://127.0.0.1:7778/api/code/search?q=<query>&scope=<ast-grep|deepcontext|file>
```

---

## Best Practices

### 1. Pre-Task Context Harvesting
Before starting significant work, harvest context from L2 memory:
1. Use `mcp_tormentnexus_memory_scratchpad_get` to check current L1 scratchpad
2. Search L2 via Kernel API: `GET /api/memory/search?q=<task description>&limit=5`
3. Optionally search skills: `GET /api/skills/search?q=<task description>&limit=5`
4. Store the findings in the L1 scratchpad (`mcp_tormentnexus_memory_scratchpad_append`)

### 2. Storing Knowledge
After key decisions, fixes, or discoveries:
1. **For lightweight**: Use `mcp_tormentnexus_add_bookmark` 
2. **For full memory**: POST to `/api/memory/add` with structured JSON content including tags
3. **For project-specific memories**: Include a `project:<dirname>` tag so they sync to `.memdb`

Good candidates: architectural decisions, bug fixes, build procedures, tool quirks, conventions.

Tag conventions:
- `project:<name>` — scope to a project
- `pattern:<name>` — reusable pattern
- `failure:<description>` — what went wrong
- `convention:<topic>` — team/project conventions
- `system:*` — system-level metadata tags

### 3. Tool Discovery
When unsure what tool to use:
1. List available MCP servers: `mcp_tormentnexus_mcp_list_servers`
2. List tools from a server: `mcp_tormentnexus_mcp_list_tools` with server name
3. Semantic search: `GET /api/mcp/native/search?query=<natural language description>`
4. Route through Kernel: `mcp_tormentnexus_mcp_call_tool` for downstream servers

### 4. Session & Skill Discovery
- Search imported sessions: `GET /api/memory/search?q=<query>&type=session`
- Search skill registry: `GET /api/skills/search?q=<query>`
- Browse all sessions: varies by Kernel endpoint

### 5. Code Search (Multi-Engine)
When searching code across the workspace:
- **AST patterns**: Use the TN Kernel's code search with `scope=ast-grep`
- **Semantic search**: Use `scope=deepcontext`
- **File/pattern**: Use `scope=file`
- Or directly use `mcp_tormentnexus_grep` for simple regex

### 6. Commercial RBAC
The TN Kernel enforces commercial authorization for dangerous operations. If a destructive tool call is blocked, use `mcp_tormentnexus_memory_scratchpad_set` or the TN Kernel memory API to record what you were doing and why before proceeding with explicitly authorized alternatives.

### 7. Per-Turn Memory Injection
After each significant user request, consider:
1. What relevant context might L2 have?
2. Search Kernel memory API for related memories
3. If found, use `mcp_tormentnexus_memory_scratchpad_append` to bring it into the current L1 context

---

## Slash Command Handling

When the user types a `/tn-*` slash command, handle it as follows:

### `/tn-store` — Interactive Memory Store
Ask the user for:
1. Content to store
2. Category (general, pattern, decision, convention, insight, failure, correction, preference)
3. Optional project name (for per-project .memdb sync)
4. Optional tags

Then POST to `/api/memory/add` and confirm storage.

### `/tn-search` — Interactive Memory Search
Ask the user for:
1. Search query (keyword or natural language)
2. Optional tag filter
3. Optional category

Then GET from `/api/memory/search` and display results.

### `/tn-status` — System Status
Fetch and display:
- L2 vault count: `GET /api/memory/fts-search?q=the&limit=1`
- L3 cold archive count: `GET /api/memory/cold-archive/count`
- Mesh peers: `GET /api/mesh/status`
- MCP server status: `mcp_tormentnexus_mcp_status`

### `/tn-plan` — Plan Management
Interactive multi-step:
1. Ask: create, list, view, or complete
2. **Create**: Title + markdown steps → POST to `/api/memory/add` with `category: "plan"` and `tag: "plan:<slug>"`
3. **List**: `GET /api/memory/list` → filter `category === "plan"` → display
4. **View**: List plans → pick one → display content
5. **Complete**: List active plans → mark complete

### `/tn-summary` — Session Summary
Summarize the current session using the conversation history and L1 scratchpad. Optionally store the summary to L2 memory.

### `/tn-purge` — Purge Stale Memories
Ask what to purge (by tag, category, or query). Use Kernel API to remove or archive matching memories.

---

## Session Lifecycle

### Session Start
- Read the L1 scratchpad for any context left from previous sessions
- Optionally fetch recent L2 memories related to the current workspace

### Before Agent / Sub-Agent Tasks
- For sub-agents: include relevant L2 context in the sub-agent prompt so the child has memory context
- For new tasks in the parent: check L2 for relevant past decisions

### Context Harvesting Workflow
A structured pre-task routine:
1. Identify the task domain (project, tool, pattern)
2. Search L2 memory for relevant entries: `GET /api/memory/search?q=<domain>&limit=5`
3. Search skills: `GET /api/skills/search?q=<domain>&limit=3`
4. If you find relevant context, append it to the L1 scratchpad
5. Reference the scratchpad content during the task

### At Session End
- Store any important decisions or findings to L2 via the TN Kernel API
- Optionally store a session summary as a memory entry

---

## Configuration

The MCP server is auto-configured. To verify:

```
codewhale mcp list        # Should show "tormentnexus"
codewhale mcp connect tormentnexus  # Test the connection
codewhale mcp tools       # Show all tools (includes mcp_tormentnexus_*)
```

Kernel health check:
```bash
curl http://127.0.0.1:7778/api/health
```

Dashboard (if running): `http://127.0.0.1:7779/dashboard`

---

## Security Notes

- The Kernel runs on `127.0.0.1:7778` — localhost only
- Commercial RBAC is enforced at the TN Kernel level
- Dangerous operations (`rm -rf`, `sudo`, `DROP TABLE`, `git push --force`) are blocked by default unless explicitly authorized
- All tool calls are audited in the TN Kernel's commercial audit log
- User-initiated shell commands are also audited
