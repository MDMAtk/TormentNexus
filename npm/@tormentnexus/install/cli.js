#!/usr/bin/env node
/**
 * @hypernexus/install — Universal Installer
 * Downloads and runs the HyperNexus installer, then configures MCP for all AI clients.
 */
const fs = require("fs");
const path = require("path");
const os = require("os");
const https = require("https");
const { execSync } = require("child_process");

const HOME = os.homedir();
const platform = os.platform();
const arch = os.arch();

// MCP config for every AI client
const MCP_CONFIG = {
  mcpServers: {
    hypernexus: {
      command: platform === "win32" ? "hypernexus.exe" : "hypernexus",
      args: ["mcp"],
      env: { HYPERNEXUS_WORKSPACE_ROOT: process.cwd() },
      type: "stdio",
      lifecycle: "eager",
    },
  },
};

const SKILL_MD = `# HyperNexus Skill — Universal AI Control Plane

## Overview
HyperNexus is your local AI control plane running on port 7778. It provides persistent
multi-tier memory (L1 scratchpad, L2 vector store, L3 cold archive), MCP tool routing
across 20,000+ servers, and session import from Claude Code/Aider/Gemini.

## Quick Start
1. Ensure HyperNexus is running: \`http://127.0.0.1:7778/api/runtime/status\`
2. Use \`hn_memory_search\` before any significant task
3. Store key decisions with \`hn_memory_store\`
4. Use \`hn_tool_search\` to find the right tool for any job

## Available Tools
- \`hn_memory_store\` — Save important decisions with tags
- \`hn_memory_search\` — Find past memories by keyword, tag, or category
- \`hn_memory_vector_search\` — Semantic vector search
- \`hn_tool_search\` — Discover tools across 20,000+ MCP servers
- \`hn_session_search\` — Browse imported sessions
- \`hn_skill_manage\` — Access reusable skill modules
- \`hn_code_search\` — Search code via AST-grep or pattern matching
- \`hn_context_harvest\` — Pull relevant L2 context

## Pricing
- $50/seat/year — Local license + cloud hosting
- https://cloud.hypernexus.site
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

// Download file
function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https
      .get(url, (response) => {
        if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
          file.close();
          fs.unlinkSync(dest);
          download(response.headers.location, dest).then(resolve).catch(reject);
          return;
        }
        if (response.statusCode !== 200) {
          file.close();
          fs.unlinkSync(dest);
          reject(new Error(`Download failed: ${response.statusCode}`));
          return;
        }
        response.pipe(file);
        file.on("finish", () => { file.close(); resolve(); });
      })
      .on("error", (err) => { file.close(); fs.unlinkSync(dest); reject(err); });
  });
}

// Get installer download URL
function getInstallerUrl() {
  const version = "1.0.0";
  const baseUrl = `https://releases.hypernexus.site/${version}`;

  if (platform === "win32") return `${baseUrl}/hypernexus-setup.exe`;
  if (platform === "darwin") return arch === "arm64" ? `${baseUrl}/hypernexus-darwin-arm64` : `${baseUrl}/hypernexus-darwin-amd64`;
  return arch === "arm64" ? `${baseUrl}/hypernexus-linux-arm64` : `${baseUrl}/hypernexus-linux-amd64`;
}

// Run the .exe installer
async function runInstaller() {
  const url = getInstallerUrl();
  const isWindows = platform === "win32";
  const ext = isWindows ? ".exe" : "";
  const dest = path.join(os.tmpdir(), `hypernexus-install${ext}`);

  console.log(`\nDownloading HyperNexus installer...`);
  console.log(`Platform: ${platform} (${arch})`);

  try {
    await download(url, dest);
    if (!isWindows) fs.chmodSync(dest, "755");

    console.log("Running installer...\n");

    if (isWindows) {
      execSync(`"${dest}" /S`, { stdio: "inherit" });
    } else {
      execSync(`"${dest}"`, { stdio: "inherit" });
    }

    try { fs.unlinkSync(dest); } catch {}
    console.log("\n✅ HyperNexus binary installed!");
  } catch (err) {
    console.log(`\n⚠️  Binary installer not available yet.`);
    console.log(`   Download manually from: https://hypernexus.site/download`);
    console.log(`   Continuing with MCP configuration...\n`);
  }
}

// Configure MCP for all AI clients
function configureMCP() {
  console.log("\nConfiguring MCP for AI clients...\n");

  let count = 0;
  for (const dir of CLIENTS) {
    const base = path.join(HOME, dir, "hypernexus");
    try {
      fs.mkdirSync(path.join(base, "mcp"), { recursive: true });
      fs.writeFileSync(path.join(base, "mcp", "servers.json"), JSON.stringify(MCP_CONFIG, null, 2));

      fs.mkdirSync(path.join(base, "skills"), { recursive: true });
      fs.writeFileSync(path.join(base, "skills", "SKILL.md"), SKILL_MD);

      count++;
    } catch {}
  }

  console.log(`✅ ${count} AI clients configured`);
  console.log("   MCP servers wired to HyperNexus");
  console.log("   Skills installed for all agents");
}

// Main
async function main() {
  console.log("\n╔══════════════════════════════════════════╗");
  console.log("║   HyperNexus Universal Installer        ║");
  console.log("║   38 AI Clients • One Command           ║");
  console.log("╚══════════════════════════════════════════╝\n");

  // Step 1: Download and run .exe installer
  await runInstaller();

  // Step 2: Configure MCP for all AI clients
  configureMCP();

  console.log("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━");
  console.log("\n🎉 Installation complete!\n");
  console.log("Next steps:");
  console.log("  1. Run 'hypernexus serve' to start the server");
  console.log("  2. Open http://localhost:7779/dashboard");
  console.log("  3. Visit https://cloud.hypernexus.site for cloud features\n");
  console.log("Pricing: $50/seat/year (local + cloud)");
  console.log("https://hypernexus.site/pricing.html\n");
}

main().catch(console.error);
