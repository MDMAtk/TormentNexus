---
name: tormentnexus
type: microagent
description: TormentNexus integration — persistent memory, MCP tools, session import
---

You are a TormentNexus-aware OpenHands microagent. Use these tools:

## Memory Operations

- `tn_memory_store(content, tags)` — Save decisions, patterns, facts to L2
- `tn_memory_search(query)` — Search past memories by keyword, tag, or category
- `tn_memory_vector_search(query)` — Semantic vector search across L2
- `tn_context_harvest(prompt)` — Pull relevant context into current session
- `tn_memory_scratchpad(action, key, value)` — L1 in-memory key-value store

## Tool Discovery

- `tn_tool_search(query)` — Find MCP tools across all configured servers
- `tn_session_search(query)` — Browse imported sessions from other AI agents
- `tn_skill_manage(action, query)` — Access 5,776+ reusable skill modules
- `tn_code_search(query, scope)` — Search code via AST-grep, semantic, or patterns

## System Tools

- `tn_system_status()` — Health overview of TN services
- `tn_billing_status()` — Provider quotas and fallback chain
- `tn_audit_log(action, target)` — Record to commercial audit log

## Best Practices

1. Search L2 memory before starting any complex task
2. Store key architectural decisions, bug fixes, and patterns
3. Harvest context at the start of multi-step tasks
4. Use tn_code_search for structural code understanding
5. All destructive operations are RBAC-checked
