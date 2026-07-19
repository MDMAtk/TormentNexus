#!/usr/bin/env node
/**
 * @tormentnexus/install — setup.js
 * Self-contained installer. No external dependencies, no relative paths to repo.
 * Runs automatically on `npm install @tormentnexus/install`.
 */
const fs = require("fs");
const path = require("path");
const os = require("os");

const HOME = os.homedir();
const TN_URL = process.env.TORMENTNEXUS_URL || "http://127.0.0.1:7778";

// MCP config for every AI client
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

const SKILL_MD = `# TormentNexus Skill — Universal AI Control Plane

## Overview
TormentNexus is your local AI control plane running on port 7778. It provides persistent
multi-tier memory (L1 scratchpad, L2 vector store, L3 cold archive), MCP tool routing
across 20+ servers, and session import from Claude Code/Aider/Gemini.

## Quick Start
1. Ensure TN Kernel is running: \`http://127.0.0.1:7778/api/runtime/status\`
2. Use \`tn_memory_search\` before any significant task
3. Store key decisions with \`tn_memory_store\`
4. Use \`tn_tool_search\` to find the right tool for any job

## Available Tools
- \`tn_memory_store\` — Save important decisions with tags
- \`tn_memory_search\` — Find past memories by keyword, tag, or category
- \`tn_memory_vector_search\` — Semantic vector search
- \`tn_tool_search\` — Discover tools across 20+ MCP servers
- \`tn_session_search\` — Browse imported sessions
- \`tn_skill_manage\` — Access 5,776 reusable skill modules
- \`tn_code_search\` — Search code via AST-grep or pattern matching
- \`tn_context_harvest\` — Pull relevant L2 context
`;

// All AI clients and their config directories
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
  console.log("\n⚡ TormentNexus Universal Installer\n");

  let count = 0;
  for (const dir of CLIENTS) {
    const base = path.join(HOME, dir, "tormentnexus");
    try {
      // MCP config
      fs.mkdirSync(path.join(base, "mcp"), { recursive: true });
      fs.writeFileSync(
        path.join(base, "mcp", "servers.json"),
        JSON.stringify(MCP_CONFIG, null, 2)
      );

      // Skill
      fs.mkdirSync(path.join(base, "skills"), { recursive: true });
      fs.writeFileSync(path.join(base, "skills", "SKILL.md"), SKILL_MD);

      count++;
    } catch (e) {
      // Skip clients that don't exist
    }
  }

  console.log(`✅ ${count} AI clients configured`);
  console.log("   MCP servers wired to TormentNexus");
  console.log("   Skills installed for all agents");
  console.log("\nNext: npm install -g @tormentnexus/cli && tn status\n");
}

install();
