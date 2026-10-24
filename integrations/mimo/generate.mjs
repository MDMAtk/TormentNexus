#!/usr/bin/env node
/**
 * Generates MiMo addon sources: Desktop tools, MiMoCode skill/commands/hooks/agent.
 * Run: node integrations/mimo/generate.mjs
 */
import { mkdirSync, writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const ROOT = dirname(fileURLToPath(import.meta.url));

function w(rel, content) {
  const p = join(ROOT, rel);
  mkdirSync(dirname(p), { recursive: true });
  writeFileSync(p, content, "utf8");
  console.log("  wrote", rel);
}

const API = `const API = (typeof process !== "undefined" && (process.env.HYPERNEXUS_URL || process.env.HN_URL)) || "http://127.0.0.1:7778";`;

const FETCH_HELPER = `${API}
async function hnFetch(path, opts) {
  const url = API.replace(/\\/$/, "") + path;
  let res;
  try {
    res = await fetch(url, {
      headers: { "Content-Type": "application/json", ...(opts && opts.headers) },
      ...opts,
    });
  } catch (e) {
    return { error: "HyperNexus kernel unreachable at " + API + ": " + (e && e.message ? e.message : String(e)) + ". Start it with: hypernexus serve" };
  }
  const text = await res.text();
  try { return JSON.parse(text); } catch { return { status: res.status, body: text }; }
}
function fmt(obj) {
  if (obj == null) return "(no data)";
  if (typeof obj === "string") return obj;
  try { return JSON.stringify(obj, null, 2); } catch { return String(obj); }
}
`;

// ─── Desktop tool sources ───────────────────────────────────────────────
const tools = {};

tools.hypernexus_search_memory = {
  desc: [
    "Search HyperNexus persistent memory (L2/L3) by keyword, tag, or category.",
    "Use this BEFORE starting any multi-step task to recall past context, decisions, and facts.",
    "Returns matching memory facts with title, content, tags, and timestamps.",
    "Arguments:",
    "- query: search terms (keywords, tags, or natural language).",
    "- limit: max results to return (default 10).",
    "- namespace: optional memory namespace filter (project | global | user).",
  ].join("\n"),
  args: `{ query: z.string().min(1), limit: z.number().int().min(1).max(100).optional(), namespace: z.string().optional() }`,
  body: `const params = new URLSearchParams({ query: args.query, limit: String(args.limit || 10) });
    if (args.namespace) params.set("namespace", args.namespace);
    const data = await hnFetch("/api/memory/search?" + params.toString());
    if (data.error) return data.error;
    const items = data.results || data.facts || data.items || (Array.isArray(data) ? data : []);
    if (!items.length) return "No memories matched: " + args.query;
    return "Found " + items.length + " memories:\\n" + items.map((m, i) =>
      (i + 1) + ". [" + (m.tags || []).join(",") + "] " + (m.title || "(untitled)") + " — " + String(m.content || m.text || "").slice(0, 200)
    ).join("\\n");`,
};

tools.hypernexus_add_fact = {
  desc: [
    "Store an important fact, decision, or finding into HyperNexus persistent memory.",
    "Call this after completing a significant step, making an architectural decision, or discovering a non-obvious pattern.",
    "Arguments:",
    "- title: short descriptive title for the fact.",
    "- content: the fact / decision / finding body.",
    "- tags: array of tags for retrieval (e.g. [\"architecture\", \"bugfix\"]).",
    "- namespace: memory namespace (project | global | user), default project.",
  ].join("\n"),
  args: `{ title: z.string().min(1), content: z.string().min(1), tags: z.array(z.string()).optional(), namespace: z.string().optional() }`,
  body: `const data = await hnFetch("/api/memory/facts/add", {
      method: "POST",
      body: JSON.stringify({ title: args.title, content: args.content, tags: args.tags || [], namespace: args.namespace || "project" }),
    });
    if (data.error) return data.error;
    return "Fact stored: " + args.title + (data.id ? " (id=" + data.id + ")" : "");`,
};

tools.hypernexus_record_observation = {
  desc: [
    "Record a structured observation into the HyperNexus knowledge graph (entity + observation).",
    "Use for durable domain knowledge: entities you discovered, relationships, constraints, gotchas.",
    "Arguments:",
    "- entity: the entity name this observation is about.",
    "- observation: the observation text.",
    "- relations: optional array of {to, type} relation hints.",
  ].join("\n"),
  args: `{ entity: z.string().min(1), observation: z.string().min(1), relations: z.array(z.object({ to: z.string(), type: z.string() })).optional() }`,
  body: `const data = await hnFetch("/api/memory/observations/add", {
      method: "POST",
      body: JSON.stringify({ entity: args.entity, observation: args.observation, relations: args.relations || [] }),
    });
    if (data.error) return data.error;
    return "Observation recorded for entity: " + args.entity;`,
};

tools.hypernexus_system_status = {
  desc: [
    "Get overall HyperNexus system health: kernel status, memory tiers, MCP runtime, and recent activity.",
    "Use when the user asks about system status, or when something seems broken.",
    "No arguments.",
  ].join("\n"),
  args: `{}`,
  body: `const data = await hnFetch("/api/system/status");
    if (data.error) return data.error;
    return fmt(data);`,
};

tools.hypernexus_mcp_search = {
  desc: [
    "Semantically search across all 60+ managed MCP servers and 26,000+ tools to find the right tool for a job.",
    "Prefer this over manually guessing tool names before calling hypernexus_mcp_call.",
    "Arguments:",
    "- query: natural language description of what you want to do.",
    "- limit: max tools to return (default 8).",
  ].join("\n"),
  args: `{ query: z.string().min(1), limit: z.number().int().min(1).max(50).optional() }`,
  body: `const data = await hnFetch("/api/mcp/tools/search?query=" + encodeURIComponent(args.query) + "&limit=" + (args.limit || 8));
    if (data.error) return data.error;
    const items = data.tools || data.results || (Array.isArray(data) ? data : []);
    if (!items.length) return "No MCP tools matched: " + args.query;
    return "Found " + items.length + " tools:\\n" + items.map((t, i) =>
      (i + 1) + ". " + (t.server || t.serverName || "?") + "/" + (t.name || t.toolName || "?") + " — " + String(t.description || "").slice(0, 160)
    ).join("\\n");`,
};

tools.hypernexus_mcp_call = {
  desc: [
    "Route a tool call through the HyperNexus kernel to a downstream MCP server.",
    "Use hypernexus_mcp_search first to discover the right serverName/toolName.",
    "Arguments:",
    "- serverName: MCP server name (from search results).",
    "- toolName: tool name on that server.",
    "- arguments: JSON object of tool arguments.",
  ].join("\n"),
  args: `{ serverName: z.string().min(1), toolName: z.string().min(1), arguments: z.record(z.string(), z.any()).optional() }`,
  body: `const data = await hnFetch("/api/mcp/tools/call", {
      method: "POST",
      body: JSON.stringify({ serverName: args.serverName, toolName: args.toolName, arguments: args.arguments || {} }),
    });
    if (data.error) return data.error;
    return fmt(data.result !== undefined ? data.result : data);`,
};

tools.hypernexus_memory_vector_search = {
  desc: [
    "Semantic vector search over HyperNexus L2 memory — finds conceptually related past context even without exact keyword matches.",
    "Use for nuanced recall when keyword search is too literal.",
    "Arguments:",
    "- query: natural language query.",
    "- limit: max results (default 8).",
  ].join("\n"),
  args: `{ query: z.string().min(1), limit: z.number().int().min(1).max(50).optional() }`,
  body: `const data = await hnFetch("/api/memory/vector/search", {
      method: "POST",
      body: JSON.stringify({ query: args.query, limit: args.limit || 8 }),
    });
    if (data.error) {
      const alt = await hnFetch("/api/memory/search?query=" + encodeURIComponent(args.query) + "&limit=" + (args.limit || 8) + "&mode=vector");
      if (alt.error) return data.error;
      return fmt(alt);
    }
    return fmt(data);`,
};

tools.hypernexus_memory_list = {
  desc: [
    "List recent HyperNexus memory facts (newest first) with titles, tags, and timestamps.",
    "Arguments:",
    "- limit: max entries (default 20).",
    "- namespace: optional namespace filter.",
  ].join("\n"),
  args: `{ limit: z.number().int().min(1).max(100).optional(), namespace: z.string().optional() }`,
  body: `const params = new URLSearchParams({ limit: String(args.limit || 20) });
    if (args.namespace) params.set("namespace", args.namespace);
    const data = await hnFetch("/api/memory/facts/list?" + params.toString());
    if (data.error) {
      const alt = await hnFetch("/api/memory/search?query=*&limit=" + (args.limit || 20));
      if (alt.error) return data.error;
      return fmt(alt);
    }
    const items = data.results || data.facts || data.items || (Array.isArray(data) ? data : []);
    if (!items.length) return "No memories stored yet.";
    return items.map((m, i) => (i + 1) + ". [" + (m.tags || []).join(",") + "] " + (m.title || "(untitled)") + " (" + (m.created_at || m.createdAt || "") + ")").join("\\n");`,
};

tools.hypernexus_session_search = {
  desc: [
    "Search imported sessions (from Claude Code, Aider, Gemini, etc.) in the HyperNexus session store.",
    "Arguments:",
    "- query: search terms.",
    "- limit: max results (default 10).",
  ].join("\n"),
  args: `{ query: z.string().min(1), limit: z.number().int().min(1).max(50).optional() }`,
  body: `const data = await hnFetch("/api/sessions/search?query=" + encodeURIComponent(args.query) + "&limit=" + (args.limit || 10));
    if (data.error) return data.error;
    return fmt(data);`,
};

tools.hypernexus_session_context = {
  desc: [
    "Get or set the current session context (goal/objective) used by HyperNexus for tool-query bootstrap and context harvesting.",
    "Arguments:",
    "- action: 'get' to read current context, 'set' to update it.",
    "- goal: when action=set, the session goal text.",
    "- objective: when action=set, the session objective text.",
  ].join("\n"),
  args: `{ action: z.enum(["get", "set"]), goal: z.string().optional(), objective: z.string().optional() }`,
  body: `if (args.action === "get") {
      const data = await hnFetch("/api/session/context");
      if (data.error) return data.error;
      return fmt(data);
    }
    const data = await hnFetch("/api/session/context", {
      method: "POST",
      body: JSON.stringify({ goal: args.goal || "", objective: args.objective || "" }),
    });
    if (data.error) return data.error;
    return "Session context updated.";`,
};

tools.hypernexus_code_search = {
  desc: [
    "Search code through the HyperNexus kernel (pattern, AST-grep, or semantic).",
    "Arguments:",
    "- query: search pattern or natural language.",
    "- mode: 'pattern' (regex/ripgrep), 'ast' (ast-grep), or 'semantic'.",
    "- path: optional path to scope the search.",
  ].join("\n"),
  args: `{ query: z.string().min(1), mode: z.enum(["pattern", "ast", "semantic"]).optional(), path: z.string().optional() }`,
  body: `const data = await hnFetch("/api/code/search", {
      method: "POST",
      body: JSON.stringify({ query: args.query, mode: args.mode || "pattern", path: args.path || "" }),
    });
    if (data.error) return data.error;
    return fmt(data);`,
};

tools.hypernexus_tool_discover = {
  desc: [
    "Discover HyperNexus capabilities and MCP tools by natural-language intent (broader than hypernexus_mcp_search — includes kernel APIs and skills).",
    "Arguments:",
    "- intent: what you want to accomplish.",
  ].join("\n"),
  args: `{ intent: z.string().min(1) }`,
  body: `const data = await hnFetch("/api/mcp/tools/search?query=" + encodeURIComponent(args.intent) + "&limit=12");
    if (data.error) return data.error;
    return fmt(data);`,
};

tools.hypernexus_knowledge_graph = {
  desc: [
    "Add entities/relations or query the HyperNexus knowledge graph.",
    "Arguments:",
    "- action: 'add_entities', 'add_relations', 'add_observations', or 'search'.",
    "- entities: array of {name, type, observations} for add_entities.",
    "- relations: array of {from, to, type} for add_relations.",
    "- query: search string for action=search.",
  ].join("\n"),
  args: `{ action: z.enum(["add_entities", "add_relations", "add_observations", "search"]), entities: z.array(z.object({ name: z.string(), type: z.string().optional(), observations: z.array(z.string()).optional() })).optional(), relations: z.array(z.object({ from: z.string(), to: z.string(), type: z.string() })).optional(), query: z.string().optional() }`,
  body: `const map = {
      add_entities: ["/api/memory/entities/add", { entities: args.entities || [] }],
      add_relations: ["/api/memory/relations/add", { relations: args.relations || [] }],
      add_observations: ["/api/memory/observations/add", { entities: args.entities || [] }],
      search: ["/api/memory/search?query=" + encodeURIComponent(args.query || ""), null],
    };
    const entry = map[args.action];
    if (!entry) return "Unknown action: " + args.action;
    const [path, body] = entry;
    const data = body
      ? await hnFetch(path, { method: "POST", body: JSON.stringify(body) })
      : await hnFetch(path);
    if (data.error) return data.error;
    return fmt(data);`,
};

tools.hypernexus_billing_status = {
  desc: [
    "Get HyperNexus billing and provider status: quotas, costs, and active providers.",
    "No arguments.",
  ].join("\n"),
  args: `{}`,
  body: `const data = await hnFetch("/api/billing/status");
    if (data.error) return data.error;
    return fmt(data);`,
};

tools.hypernexus_repomap = {
  desc: [
    "Generate a ranked repository map with lightweight symbol summaries for a directory.",
    "Arguments:",
    "- dir: directory path to map (default: current workspace).",
    "- maxFiles: cap on files included (default 200).",
  ].join("\n"),
  args: `{ dir: z.string().optional(), maxFiles: z.number().int().min(1).max(2000).optional() }`,
  body: `const data = await hnFetch("/api/code/repomap", {
      method: "POST",
      body: JSON.stringify({ dir: args.dir || "", maxFiles: args.maxFiles || 200 }),
    });
    if (data.error) return data.error;
    return fmt(data);`,
};

tools.hypernexus_context_harvest = {
  desc: [
    "Harvest relevant L2 context from HyperNexus memory for the current task and return it for injection.",
    "Call before complex multi-step tasks to pull in prior related work.",
    "Arguments:",
    "- topic: the task topic / question to harvest context for.",
    "- depth: harvest depth 'quick' | 'standard' | 'deep' (default standard).",
  ].join("\n"),
  args: `{ topic: z.string().min(1), depth: z.enum(["quick", "standard", "deep"]).optional() }`,
  body: `const data = await hnFetch("/api/memory/context/harvest", {
      method: "POST",
      body: JSON.stringify({ topic: args.topic, depth: args.depth || "standard" }),
    });
    if (data.error) {
      const alt = await hnFetch("/api/memory/search?query=" + encodeURIComponent(args.topic) + "&limit=15");
      if (alt.error) return data.error;
      return fmt(alt);
    }
    return fmt(data);`,
};

console.log("Generating Desktop tool sources...");
for (const [name, t] of Object.entries(tools)) {
  const src = `import { z } from "zod";

${FETCH_HELPER}

export default {
  description: ${JSON.stringify(t.desc)},
  args: ${t.args},
  async execute(args) {
    ${t.body}
  },
};
`;
  w(join("tools", "src", name + ".js"), src);
}

// ─── Skill ──────────────────────────────────────────────────────────────
console.log("Generating skill...");
const SKILL_MD = `---
name: hypernexus
description: "HyperNexus AI control plane — persistent L2 vector memory, semantic MCP tool discovery across 60+ servers, knowledge graph, session management, billing, agent orchestration, and 400+ API endpoints. Use when the task benefits from cross-session memory (scratchpad, facts, observations), tool discovery across MCP servers, knowledge ingestion, repo analysis, system health checks, or routing through the HyperNexus Go kernel on port 7778. Trigger on mentions of 'hypernexus', 'HN', 'scratchpad', 'cross-session memory', 'MCP tools', 'tool discovery', 'knowledge graph', 'agent memory', or 'control plane'."
---

# HyperNexus Integration

HyperNexus is a local AI control plane running a Go kernel on port 7778 (dashboard 7779) with persistent multi-tier memory (L1 scratchpad, L2 vector store, L3 cold archive), semantic tool discovery across 60+ MCP servers / 26,000+ tools, knowledge graph, session import, and agent orchestration.

## When to Use

- Start of any multi-step task → harvest prior context (\`hypernexus_context_harvest\` / \`hypernexus_search_memory\`)
- After a decision, gotcha, or non-obvious finding → store it (\`hypernexus_add_fact\` / \`hypernexus_record_observation\`)
- When you need a capability you don't have → discover tools (\`hypernexus_mcp_search\` / \`hypernexus_tool_discover\`) then route (\`hypernexus_mcp_call\`)
- When the user asks about system health, billing, or the repo → \`hypernexus_system_status\`, \`hypernexus_billing_status\`, \`hypernexus_repomap\`

## MiMo Desktop Tools (native)

| Tool | Purpose |
|---|---|
| \`hypernexus_search_memory\` | Keyword search over persistent memory |
| \`hypernexus_memory_vector_search\` | Semantic / conceptual memory recall |
| \`hypernexus_memory_list\` | List recent facts |
| \`hypernexus_add_fact\` | Store a fact / decision |
| \`hypernexus_record_observation\` | Record a knowledge-graph observation |
| \`hypernexus_knowledge_graph\` | Add entities/relations or search the graph |
| \`hypernexus_context_harvest\` | Pull relevant L2 context for a task |
| \`hypernexus_session_search\` | Search imported sessions |
| \`hypernexus_session_context\` | Get/set session goal & objective |
| \`hypernexus_mcp_search\` | Semantic MCP tool search |
| \`hypernexus_tool_discover\` | Broad capability discovery |
| \`hypernexus_mcp_call\` | Route a call to a downstream MCP tool |
| \`hypernexus_code_search\` | Pattern / AST / semantic code search |
| \`hypernexus_repomap\` | Repository map with symbols |
| \`hypernexus_system_status\` | Kernel + memory + MCP health |
| \`hypernexus_billing_status\` | Provider quotas & costs |

## HTTP API (port 7778)

\`\`\`bash
# Memory
curl -X POST http://localhost:7778/api/memory/facts/add -H 'Content-Type: application/json' \\
  -d '{"title":"...","content":"...","tags":["..."],"namespace":"project"}'
curl -s "http://localhost:7778/api/memory/search?query=<term>&limit=<n>"

# MCP discovery & routing
curl -s "http://localhost:7778/api/mcp/tools/search?query=<natural+language>"
curl -X POST http://localhost:7778/api/mcp/tools/call -H 'Content-Type: application/json' \\
  -d '{"serverName":"...","toolName":"...","arguments":{}}'

# System
curl -s http://localhost:7778/health
curl -s http://localhost:7778/api/system/status
curl -s http://localhost:7778/api/billing/status
\`\`\`

## MiMoCode CLI Commands

- \`/hn-store\` — store a memory fact
- \`/hn-search\` — search persistent memory
- \`/hn-status\` — system status
- \`/hn-tools\` — discover MCP tools
- \`/hn-session\` — session context / search

## Best Practices

1. **Recall before you write**: search memory before large tasks.
2. **Store immediately**: non-obvious findings decay fast — persist them when found.
3. **Discover then route**: use \`hypernexus_mcp_search\` before hardcoding tool names.
4. **Graph what matters**: entities + relations compound into durable knowledge.
5. The Go kernel is the source of truth — prefer the HTTP API over shell improvisation.
`;

w(join("skills", "hypernexus", "SKILL.md"), SKILL_MD);
w(join("skills", "hypernexus", "locales", "en-US.json"), JSON.stringify({ displayName: "HyperNexus", brief: "Persistent memory, MCP tool routing, and knowledge graph via the HyperNexus control plane" }, null, 2) + "\n");
w(join("skills", "hypernexus", "locales", "zh-CN.json"), JSON.stringify({ displayName: "HyperNexus", brief: "通过 HyperNexus 控制平面提供持久记忆、MCP 工具路由与知识图谱" }, null, 2) + "\n");

// ─── MiMoCode commands ──────────────────────────────────────────────────
console.log("Generating MiMoCode commands...");
const commands = {
  "hn-store": `# /hn-store — Store a memory fact

Store an important fact, decision, or finding into HyperNexus persistent memory.

Usage: \`/hn-store <title> | <content> [| tag1,tag2]\`

Examples:
- \`/hn-store DB path | catalog.db was renamed to hypernexus.db in alpha.251 | bugfix,db\`
- \`/hn-store Deploy target | Production runs on Hetzner at 5.161.250.43 (ssh alias: hetzner) | deploy,infra\`

After storing, confirm with the returned fact id if provided.
`,
  "hn-search": `# /hn-search — Search persistent memory

Search HyperNexus L2/L3 memory by keyword, tag, or natural language.

Usage: \`/hn-search <query>\`

Tips:
- Use distinctive terms (function name, error string, decision phrase).
- For conceptual matches, use the \`hypernexus_memory_vector_search\` tool instead.
`,
  "hn-status": `# /hn-status — System status

Show HyperNexus kernel health: memory tiers, MCP runtime, active services.

Usage: \`/hn-status\`

Calls the \`hypernexus_system_status\` tool (or GET /api/system/status).
`,
  "hn-tools": `# /hn-tools — Discover MCP tools

Semantically search 26,000+ tools across 60+ MCP servers for a capability.

Usage: \`/hn-tools <what you want to do>\`

Example: \`/hn-tools create a Jira issue from a bug report\`

Then use \`hypernexus_mcp_call\` to invoke the discovered tool.
`,
  "hn-session": `# /hn-session — Session context & search

Get/set the session goal, or search imported sessions.

Usage:
- \`/hn-session get\` — read current goal/objective
- \`/hn-session set <goal>\` — update session goal
- \`/hn-session search <query>\` — search imported sessions
`,
  "hn-remember": `# /hn-remember — Record a knowledge-graph observation

Record a structured observation about an entity into the HyperNexus knowledge graph.

Usage: \`/hn-remember <entity> | <observation>\`

Example: \`/hn-remember sessionbridge | callUpstreamJSON at session_context_handlers.go:138 is the central upstream caller (~150 call sites)\`
`,
};
for (const [name, body] of Object.entries(commands)) {
  w(join("mimocode", "commands", name + ".md"), body);
}

// ─── Hooks ──────────────────────────────────────────────────────────────
console.log("Generating hooks...");
w(join("mimocode", "hooks", "config.json"), JSON.stringify({
  on_session_start: "hypernexus_context_harvest",
  on_tool_error: "hypernexus_add_fact",
  on_decision: "hypernexus_add_fact",
}, null, 2) + "\n");

// ─── Agent ──────────────────────────────────────────────────────────────
console.log("Generating agent...");
w(join("mimocode", "agents", "hypernexus-orchestrator.md"), `# HyperNexus Orchestrator Agent

## Role
You are an autonomous orchestrator with full access to the HyperNexus AI control plane:
persistent multi-tier memory, 60+ MCP servers / 26,000+ tools, knowledge graph, and session management.

## Operating Loop
1. **Recall** — \`hypernexus_context_harvest\` / \`hypernexus_search_memory\` before any multi-step work.
2. **Discover** — \`hypernexus_mcp_search\` / \`hypernexus_tool_discover\` before hardcoding tool names.
3. **Execute** — route work through \`hypernexus_mcp_call\` or local tools.
4. **Persist** — \`hypernexus_add_fact\` for decisions/gotchas, \`hypernexus_record_observation\` for domain knowledge.
5. **Verify** — \`hypernexus_system_status\` when anything seems off.

## Rules
- Never repeat work already recorded in memory — search first.
- Store non-obvious findings immediately; do not wait until the end of the task.
- Prefer kernel HTTP API (port 7778) over shell improvisation.
- The Go kernel is the source of truth.
`);

// ─── MiMoCode skill copy (commands + skill for CLI) ─────────────────────
console.log("Generating MiMoCode CLI skill tree...");
w(join("mimocode", "skills", "SKILL.md"), SKILL_MD);

console.log("Done.");
