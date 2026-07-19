#!/usr/bin/env node
/**
 * @tormentnexus/core — setup.js
 * Core installer. Wires MCP + skills for all AI clients.
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

const SKILL_MD = `# TormentNexus Core Skill

## Memory Operations
- \`tn_memory_store(content, tags)\` — Save to L2 memory
- \`tn_memory_search(query)\` — Search by keyword/tag
- \`tn_memory_vector_search(query)\` — Semantic search
- \`tn_context_harvest(prompt)\` — Pull relevant context

## Tool Discovery
- \`tn_tool_search(query)\` — Find MCP tools
- \`tn_session_search(query)\` — Browse past sessions

## Best Practices
1. Search memory before any complex task
2. Store key decisions as you work
3. Harvest context for multi-step tasks
`;

const CLIENTS = [
  ".claude", ".gemini", ".codex", ".grok", ".antigravity",
  ".aider", ".opencode", ".openclaw", ".goose", ".iflow",
  ".roo", ".cline", ".cursor", ".windsurf", ".zed", ".trae",
  ".continue", ".factory", ".openhands", ".kiro", ".codewhale",
  ".omnigent", ".citadel", ".agent-fusion", ".herdr", ".claude-squad",
  ".qwen-code", ".qwen", ".pi", ".kimi-code", ".moonshot",
  ".cliproxyapi", ".vscode", ".jetbrains", ".hermes",
];

function install() {
  console.log("\n🧠 TormentNexus Core Installer\n");

  let count = 0;
  for (const dir of CLIENTS) {
    const base = path.join(HOME, dir, "tormentnexus");
    try {
      fs.mkdirSync(path.join(base, "mcp"), { recursive: true });
      fs.writeFileSync(path.join(base, "mcp", "servers.json"), JSON.stringify(MCP_CONFIG, null, 2));
      fs.mkdirSync(path.join(base, "skills"), { recursive: true });
      fs.writeFileSync(path.join(base, "skills", "SKILL.md"), SKILL_MD);
      count++;
    } catch { /* skip */ }
  }

  console.log(`✅ ${count} AI clients wired to TormentNexus`);
  console.log("   Use: tn search <query> to access memory");
  console.log("   Install CLI: npm install -g @tormentnexus/cli\n");
}

install();
