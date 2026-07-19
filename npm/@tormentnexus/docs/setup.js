#!/usr/bin/env node
const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const os = require("os");

const TN_URL = process.env.TORMENTNEXUS_URL || "http://127.0.0.1:7778";
const DOCS_DIR = path.join(os.homedir(), ".tormentnexus", "docs");

function fetch(urlPath) {
	try {
		return JSON.parse(
			execSync(`curl -s --max-time 5 "${TN_URL}${urlPath}"`, {
				encoding: "utf8",
				stdio: ["pipe", "pipe", "pipe"],
			}),
		);
	} catch {
		return null;
	}
}

console.log("\n📚 TormentNexus Docs Generator\n");

fs.mkdirSync(DOCS_DIR, { recursive: true });

// Fetch status
const status = fetch("/api/runtime/status");
if (!status?.data) {
	console.log("❌ TormentNexus kernel not running at " + TN_URL);
	console.log("   Start with: npm install -g @tormentnexus/core");
	process.exit(1);
}

const d = status.data;
const version = d.version;
const tools = d.cli?.tools || [];
const harnesses = d.cli?.harnessCount || 0;

// Generate API reference
let apidoc = `# TormentNexus API Reference — v${version}\n\n`;
apidoc += `> Base URL: ${TN_URL}\n\n`;
apidoc += `## System Status\n\n`;
apidoc += `- **Version:** ${version}\n`;
apidoc += `- **Uptime:** ${d.uptimeSec}s\n`;
apidoc += `- **Tools:** ${tools.length} registered\n`;
apidoc += `- **Harnesses:** ${harnesses} tracked\n`;
apidoc += `- **Memory sections:** ${d.memory?.sectionCount || 0}\n\n`;

// Generate tool catalog
let tooldoc = `# MCP Tool Catalog\n\n`;
tooldoc += `## Registered Tools (${tools.length})\n\n`;
for (const t of tools) {
	tooldoc += `- \`${t}\`\n`;
}
tooldoc += `\n## Session Import Sources (${d.importSources?.count || 0})\n\n`;
if (d.importSources?.bySourceTool) {
	for (const src of d.importSources.bySourceTool) {
		tooldoc += `- ${src.key}: ${src.count} artifacts\n`;
	}
}

// Generate install guide
let installDoc = `# TormentNexus Installation Guide\n\n`;
installDoc += `## Quick Start\n\n`;
installDoc += `\`\`\`bash\nnpx @tormentnexus/install\nnpm install -g @tormentnexus/cli\ntn status\n\`\`\`\n\n`;
installDoc += `## Supported AI Clients (38+)\n\n`;
installDoc += `Claude, Gemini, Cursor, Windsurf, OpenHands, Aider, CodeWhale, OpenCode,`;
installDoc += `Goose, Cline, Roo, Continue, Zed, Trae, Factory, Kiro, Pi, Kimi-Code,`;
installDoc += `Qwen-Code, OmniGent, Citadel, Agent-Fusion, Herdr, Claude-Squad, VS Code, JetBrains\n\n`;

// Write all docs
fs.writeFileSync(path.join(DOCS_DIR, "API_REFERENCE.md"), apidoc);
fs.writeFileSync(path.join(DOCS_DIR, "TOOL_CATALOG.md"), tooldoc);
fs.writeFileSync(path.join(DOCS_DIR, "INSTALL_GUIDE.md"), installDoc);

console.log(`✅ Generated documentation at ${DOCS_DIR}/`);
console.log(`   API_REFERENCE.md  — Full API surface`);
console.log(`   TOOL_CATALOG.md   — ${tools.length} registered tools`);
console.log(`   INSTALL_GUIDE.md  — Quick start guide`);
console.log(`\nOpen: file://${DOCS_DIR}/API_REFERENCE.md\n`);
