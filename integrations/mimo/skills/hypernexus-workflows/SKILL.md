---
name: hypernexus-workflows
description: "Multi-step orchestration workflows composing the 16 HyperNexus Desktop tools into repeatable patterns — memory-first task execution, tool discovery pipelines, knowledge graph building, session handoff, production runbook, and deep context harvesting. Use when the user wants to run a structured multi-step workflow, orchestrate HyperNexus tools together, set up a memory-driven task loop, build the knowledge graph incrementally, hand off session context cleanly, or execute a production runbook. Trigger on 'workflow', 'orchestrate', 'runbook', 'pipeline', 'multi-step', 'memory-first', 'knowledge graph build', 'session handoff', 'context harvest', or when a task spans memory recall + tool routing + persistence."
---

# HyperNexus Workflows — Orchestration Templates

Composable multi-step workflows built from the 16 HyperNexus Desktop tools. Each workflow is a named pattern with clear entry criteria, tool sequence, and exit criteria.

## Workflow 1: Memory-First Task Execution (MEM-FIRST)

**When**: Starting any multi-step task or feature build.

| Step | Tool | Purpose |
|---|---|---|
| 1 | `hypernexus_context_harvest` | Pull relevant L2 context for the task topic |
| 2 | `hypernexus_search_memory` | Keyword search for prior decisions/gotchas |
| 3 | `hypernexus_memory_vector_search` | Conceptual recall for nuanced context |
| 4 | Execute the task | Do the work |
| 5 | `hypernexus_add_fact` | Store decisions, gotchas, non-obvious findings |
| 6 | `hypernexus_record_observation` | Record domain knowledge in the graph |

**Exit criteria**: Facts stored for all non-obvious findings.

## Workflow 2: Tool Discovery Pipeline (DISCOVER)

**When**: You need a capability you don't have locally.

| Step | Tool | Purpose |
|---|---|---|
| 1 | `hypernexus_tool_discover` | Broad capability search by intent |
| 2 | `hypernexus_mcp_search` | Narrow search across 26k+ MCP tools |
| 3 | `hypernexus_mcp_call` | Route the call through the kernel |
| 4 | `hypernexus_add_fact` | Store what worked for future runs |

**Exit criteria**: Tool found and called successfully.

## Workflow 3: Knowledge Graph Building (GRAPH)

**When**: Learning a new codebase, system, or domain.

| Step | Tool | Purpose |
|---|---|---|
| 1 | `hypernexus_repomap` | Get repository map with symbol summaries |
| 2 | `hypernexus_code_search` | Find key patterns and implementations |
| 3 | `hypernexus_knowledge_graph` (add_entities) | Register discovered entities |
| 4 | `hypernexus_knowledge_graph` (add_relations) | Wire entity relationships |
| 5 | `hypernexus_record_observation` | Attach observations to entities |

**Exit criteria**: Entities + relations + observations persisted.

## Workflow 4: Session Handoff (HANDOFF)

**When**: End of a work session or context switch.

| Step | Tool | Purpose |
|---|---|---|
| 1 | `hypernexus_session_context` (get) | Read current goal/objective |
| 2 | `hypernexus_add_fact` | Store session summary |
| 3 | `hypernexus_session_context` (set) | Update goal for next session |
| 4 | `hypernexus_record_observation` | Record session learnings |

**Exit criteria**: Session state persisted and goal updated.

## Workflow 5: Production Runbook (OPS)

**When**: Deployments, incidents, health checks.

| Step | Tool | Purpose |
|---|---|---|
| 1 | `hypernexus_system_status` | Check kernel + memory + MCP health |
| 2 | `hypernexus_billing_status` | Verify provider quotas |
| 3 | Execute ops action | Restart, deploy, fix |
| 4 | `hypernexus_system_status` | Verify health after action |
| 5 | `hypernexus_add_fact` | Record what was done and why |

**Exit criteria**: Health verified post-action and change recorded.

## Workflow 6: Deep Context Harvest (HARVEST)

**When**: Starting a complex multi-day task or resuming after a break.

| Step | Tool | Purpose |
|---|---|---|
| 1 | `hypernexus_context_harvest` (depth=deep) | Pull deep L2 context |
| 2 | `hypernexus_memory_list` | Review recent facts |
| 3 | `hypernexus_session_search` | Find related past sessions |
| 4 | `hypernexus_memory_vector_search` | Conceptual recall |
| 5 | `hypernexus_search_memory` | Keyword recall |

**Exit criteria**: Full context assembled before starting work.

## Composition Rules

1. **Always start with recall** (Workflow 1 or 6) before any multi-step work.
2. **Always end with persist** (store fact + observation) for anything non-obvious.
3. **Chain workflows**: HARVEST → MEM-FIRST → HANDOFF for long sessions.
4. **Use DISCOVER mid-flow** when you hit a capability gap.
5. **Use GRAPH when** the task involves learning a system or domain.
6. **Use OPS before and after** any production change.

## The Operating Loop (all workflows share this)

```
RECALL → DISCOVER (if needed) → EXECUTE → PERSIST → VERIFY
```
