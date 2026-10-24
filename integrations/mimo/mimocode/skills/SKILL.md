---
name: hypernexus
description: "HyperNexus AI control plane — persistent L2 vector memory, semantic MCP tool discovery across 60+ servers, knowledge graph, session management, billing, agent orchestration, and 400+ API endpoints. Use when the task benefits from cross-session memory (scratchpad, facts, observations), tool discovery across MCP servers, knowledge ingestion, repo analysis, system health checks, or routing through the HyperNexus Go kernel on port 7778. Trigger on mentions of 'hypernexus', 'HN', 'scratchpad', 'cross-session memory', 'MCP tools', 'tool discovery', 'knowledge graph', 'agent memory', or 'control plane'."
---

# HyperNexus Integration

HyperNexus is a local AI control plane running a Go kernel on port 7778 (dashboard 7779) with persistent multi-tier memory (L1 scratchpad, L2 vector store, L3 cold archive), semantic tool discovery across 60+ MCP servers / 26,000+ tools, knowledge graph, session import, and agent orchestration.

## When to Use

- Start of any multi-step task → harvest prior context (`hypernexus_context_harvest` / `hypernexus_search_memory`)
- After a decision, gotcha, or non-obvious finding → store it (`hypernexus_add_fact` / `hypernexus_record_observation`)
- When you need a capability you don't have → discover tools (`hypernexus_mcp_search` / `hypernexus_tool_discover`) then route (`hypernexus_mcp_call`)
- When the user asks about system health, billing, or the repo → `hypernexus_system_status`, `hypernexus_billing_status`, `hypernexus_repomap`

## MiMo Desktop Tools (native)

| Tool | Purpose |
|---|---|
| `hypernexus_search_memory` | Keyword search over persistent memory |
| `hypernexus_memory_vector_search` | Semantic / conceptual memory recall |
| `hypernexus_memory_list` | List recent facts |
| `hypernexus_add_fact` | Store a fact / decision |
| `hypernexus_record_observation` | Record a knowledge-graph observation |
| `hypernexus_knowledge_graph` | Add entities/relations or search the graph |
| `hypernexus_context_harvest` | Pull relevant L2 context for a task |
| `hypernexus_session_search` | Search imported sessions |
| `hypernexus_session_context` | Get/set session goal & objective |
| `hypernexus_mcp_search` | Semantic MCP tool search |
| `hypernexus_tool_discover` | Broad capability discovery |
| `hypernexus_mcp_call` | Route a call to a downstream MCP tool |
| `hypernexus_code_search` | Pattern / AST / semantic code search |
| `hypernexus_repomap` | Repository map with symbols |
| `hypernexus_system_status` | Kernel + memory + MCP health |
| `hypernexus_billing_status` | Provider quotas & costs |

## HTTP API (port 7778)

```bash
# Memory
curl -X POST http://localhost:7778/api/memory/facts/add -H 'Content-Type: application/json' \
  -d '{"title":"...","content":"...","tags":["..."],"namespace":"project"}'
curl -s "http://localhost:7778/api/memory/search?query=<term>&limit=<n>"

# MCP discovery & routing
curl -s "http://localhost:7778/api/mcp/tools/search?query=<natural+language>"
curl -X POST http://localhost:7778/api/mcp/tools/call -H 'Content-Type: application/json' \
  -d '{"serverName":"...","toolName":"...","arguments":{}}'

# System
curl -s http://localhost:7778/health
curl -s http://localhost:7778/api/system/status
curl -s http://localhost:7778/api/billing/status
```

## MiMoCode CLI Commands

- `/hn-store` — store a memory fact
- `/hn-search` — search persistent memory
- `/hn-status` — system status
- `/hn-tools` — discover MCP tools
- `/hn-session` — session context / search

## Best Practices

1. **Recall before you write**: search memory before large tasks.
2. **Store immediately**: non-obvious findings decay fast — persist them when found.
3. **Discover then route**: use `hypernexus_mcp_search` before hardcoding tool names.
4. **Graph what matters**: entities + relations compound into durable knowledge.
5. The Go kernel is the source of truth — prefer the HTTP API over shell improvisation.
