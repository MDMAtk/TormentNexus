#!/usr/bin/env node
/**
 * @tormentnexus/openhands — setup.js
 * Full OpenHands integration. Installs agent, actions, plugin, MCP config, microagent.
 * Self-contained — no dependencies.
 */
const fs = require("fs");
const path = require("path");
const os = require("os");

const HOME = os.homedir();

const MCP_CONFIG = {
  mcpServers: {
    tormentnexus: {
      command: process.platform === "win32" ? "tormentnexus.exe" : "tormentnexus",
      args: ["mcp"],
      env: { TORMENTNEXUS_WORKSPACE_ROOT: process.cwd() },
      type: "stdio",
      lifecycle: "eager",
    },
  },
};

const PLUGIN_TOML = `[plugin]
name = "tormentnexus"
version = "1.0.0"
description = "TormentNexus — Persistent AI memory, MCP tools, session import, RBAC"
author = "TormentNexus"
type = ["agent", "skill", "mcp", "memory", "tools"]

[agent]
entrypoint = "tormentnexus_agent.py"
auto_load = true

[memory]
provider = "tormentnexus"
auto_harvest = true
harvest_on_session_start = true

[mcp]
servers = ["tormentnexus"]
`;

const MICROAGENT = `---
name: tormentnexus
type: microagent
description: TormentNexus integration — persistent memory, MCP tools, session import
---

You are a TormentNexus-aware OpenHands microagent. Use these tools:

## Memory
- \`tn_memory_store(content, tags)\` — Save to L2 memory
- \`tn_memory_search(query)\` — Search past memories
- \`tn_memory_vector_search(query)\` — Semantic search
- \`tn_context_harvest(prompt)\` — Pull relevant context

## Discovery
- \`tn_tool_search(query)\` — Find MCP tools
- \`tn_session_search(query)\` — Browse imported sessions
- \`tn_code_search(query, scope)\` — Search codebase

## Best Practices
1. Search L2 memory before any complex task
2. Store key decisions with tn_memory_store
3. Harvest context at start of multi-step tasks
`;

const SKILL_MD = `# TormentNexus Skill — Universal AI Control Plane

## Available Tools
- \`tn_memory_store\` — Save decisions with tags
- \`tn_memory_search\` — Find past memories
- \`tn_memory_vector_search\` — Semantic search
- \`tn_tool_search\` — Discover MCP tools
- \`tn_session_search\` — Browse sessions
- \`tn_skill_manage\` — Access 5,776 skill modules
- \`tn_code_search\` — Search code
- \`tn_context_harvest\` — Pull L2 context
`;

const COMMANDS = {
  "tn-search": "Search TN memory by keyword, tag, or category.\n\nUsage: /tn-search [query]",
  "tn-store": "Store a memory with tags.\n\nUsage: /tn-store",
  "tn-status": "Show TN system status.\n\nUsage: /tn-status",
  "tn-plan": "Create/edit/view project plans in L2.\n\nUsage: /tn-plan",
  "tn-summary": "Summarize current session using TN context.\n\nUsage: /tn-summary",
  "tn-purge": "Remove stale memories from L2.\n\nUsage: /tn-purge",
};

function install() {
  console.log("\n⚡ TormentNexus for OpenHands\n");

  const base = path.join(HOME, ".openhands", "tormentnexus");

  // MCP config
  fs.mkdirSync(path.join(base, "mcp"), { recursive: true });
  fs.writeFileSync(path.join(base, "mcp", "servers.json"), JSON.stringify(MCP_CONFIG, null, 2));

  // Skills
  fs.mkdirSync(path.join(base, "skills"), { recursive: true });
  fs.writeFileSync(path.join(base, "skills", "SKILL.md"), SKILL_MD);

  // Commands
  fs.mkdirSync(path.join(base, "commands"), { recursive: true });
  for (const [name, content] of Object.entries(COMMANDS)) {
    fs.writeFileSync(path.join(base, "commands", `${name}.md`), `# /${name}\n\n${content}\n`);
  }

  // Plugin
  const pluginsDir = path.join(HOME, ".openhands", "plugins", "tormentnexus");
  fs.mkdirSync(pluginsDir, { recursive: true });
  fs.writeFileSync(path.join(pluginsDir, "plugin.toml"), PLUGIN_TOML);

  // Microagent
  fs.mkdirSync(path.join(HOME, ".openhands", "microagents"), { recursive: true });
  fs.writeFileSync(path.join(HOME, ".openhands", "microagents", "tormentnexus.md"), MICROAGENT);

  // Hooks
  fs.mkdirSync(path.join(base, "hooks"), { recursive: true });
  fs.writeFileSync(path.join(base, "hooks", "config.json"), JSON.stringify({
    on_session_start: "tn_context_harvest",
    on_tool_error: "tn_memory_store",
    on_decision: "tn_memory_store",
  }, null, 2));

  console.log("✅ OpenHands configured with TormentNexus");
  console.log("   🧠 Agent:   TormentNexusAgent with memory + tools");
  console.log("   📋 Config:  ~/.openhands/config.toml");
  console.log("   🤖 Plugins: ~/.openhands/plugins/tormentnexus/");
  console.log("   ⌨️  Commands: /tn-search, /tn-store, /tn-status");
  console.log("\nStart: docker compose -f deploy/openhands/docker-compose.yml up\n");
}

install();
