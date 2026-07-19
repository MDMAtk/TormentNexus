---
description: Store a memory in TormentNexus L2. Use to persist key decisions, patterns, code conventions, and important context cross-session.
---

# Store in TormentNexus (/tn-store)

When I say "tn-store" or "store this in TN", help me save information to TormentNexus L2 memory.

1. **Identify Content** — Extract the key information to store (decision, pattern, convention, finding)
2. **Tag** — Add relevant tags (e.g., `project:foo`, `pattern`, `api`, `decision`)
3. **Store** — Use `mcp_tormentnexus_memory_scratchpad_set` with the content
4. **Confirm** — Summarize what was saved

Example: `/tn-store The project uses React 19 with Vite and pnpm workspaces.`

To check existing context: `mcp_tormentnexus_memory_scratchpad_get`
