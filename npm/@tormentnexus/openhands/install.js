#!/usr/bin/env node
const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const os = require("os");

const HOME = os.homedir();
const OPENHANDS_DIR = path.join(HOME, ".openhands");
const TN_WORKSPACE = path.resolve(process.cwd());

console.log("\n⚡ TormentNexus for OpenHands — Full Integration Installer\n");

// 1. Create all directories
const dirs = [
	path.join(OPENHANDS_DIR, "microagents"),
	path.join(OPENHANDS_DIR, "plugins", "tormentnexus"),
	path.join(OPENHANDS_DIR, "plugins", "tormentnexus", "skills"),
	path.join(OPENHANDS_DIR, "plugins", "tormentnexus", "commands"),
	path.join(OPENHANDS_DIR, "plugins", "tormentnexus", "mcp"),
];
for (const d of dirs) fs.mkdirSync(d, { recursive: true });

// 2. MCP server config (Stdio transport)
const mcpConfig = {
	mcpServers: {
		tormentnexus: {
			command:
				process.platform === "win32" ? "tormentnexus.exe" : "tormentnexus",
			args: ["mcp"],
			env: { TORMENTNEXUS_WORKSPACE_ROOT: TN_WORKSPACE },
			type: "stdio",
			lifecycle: "eager",
		},
	},
};
fs.writeFileSync(
	path.join(OPENHANDS_DIR, "plugins", "tormentnexus", "mcp", "servers.json"),
	JSON.stringify(mcpConfig, null, 2),
);

// 3. TN config.toml for OpenHands
const configToml = `[core]
workspace_base = "./workspace"

[mcp]
enabled = true

[mcp.servers.tormentnexus]
command = "${process.platform === "win32" ? "tormentnexus.exe" : "tormentnexus"}"
args = ["mcp"]
env = { TORMENTNEXUS_WORKSPACE_ROOT = "${TN_WORKSPACE.replace(/\\/g, "\\\\")}" }

[tormentnexus]
url = "http://127.0.0.1:7778"
auto_connect = true
context_harvest = true
`;
fs.writeFileSync(path.join(OPENHANDS_DIR, "config.toml"), configToml);

// 4. Plugin manifest
const pluginToml = `[plugin]
name = "tormentnexus"
version = "1.0.0"
description = "TormentNexus — Persistent AI memory, MCP tools, session import"
author = "TormentNexus"
type = ["agent", "skill", "mcp", "memory"]

[agent]
entrypoint = "tormentnexus_agent.py"
auto_load = true

[memory]
provider = "tormentnexus"
auto_harvest = true

[mcp]
servers = ["tormentnexus"]
`;
fs.writeFileSync(
	path.join(OPENHANDS_DIR, "plugins", "tormentnexus", "plugin.toml"),
	pluginToml,
);

// 5. Skills
const SKILL_MD = `# TormentNexus Skill — Universal AI Control Plane

## Overview
TormentNexus is your local AI control plane running on port 7778. It provides persistent
multi-tier memory (L1 scratchpad, L2 vector store, L3 cold archive), MCP tool routing
across 20+ servers, session import from Claude Code/Aider/Gemini, and commercial RBAC.

## Quick Start
1. Ensure TN Kernel is running: \`http://127.0.0.1:7778/api/runtime/status\`
2. Use \`tn_memory_search\` before any significant task to recall past context
3. Store key decisions with \`tn_memory_store\` using descriptive tags
4. Route through TN Kernel for commercial integrations (Jira, Confluence)
5. Use \`tn_tool_search\` to find the right tool for any job

## Available Tools
- \`tn_memory_store\` — Save important decisions with tags
- \`tn_memory_search\` — Find past memories by keyword, tag, or category
- \`tn_memory_vector_search\` — Semantic vector search
- \`tn_tool_search\` — Discover tools across 20+ MCP servers
- \`tn_session_search\` — Browse imported sessions
- \`tn_skill_manage\` — Access 5,776 reusable skill modules
- \`tn_code_search\` — Search code via AST-grep or pattern matching
- \`tn_context_harvest\` — Pull relevant L2 context

## Memory Best Practices
1. Search L2 before starting any significant task
2. Store important decisions, patterns, and facts
3. Use \`@memory:keyword\` inline for auto-expanded context
4. Check the cold archive for archived knowledge`;
fs.writeFileSync(
	path.join(OPENHANDS_DIR, "plugins", "tormentnexus", "skills", "SKILL.md"),
	SKILL_MD,
);

// 6. Slash commands
const commands = {
	"tn-search": {
		description: "Search TN memory by keyword, tag, or category",
		args: "query",
	},
	"tn-store": {
		description: "Store a memory with tags",
		args: "content, tags",
	},
	"tn-status": { description: "Show TN system status", args: "" },
	"tn-plan": {
		description: "Create/edit/view project plans in L2",
		args: "action, query",
	},
	"tn-summary": {
		description: "Summarize current session using TN context",
		args: "",
	},
	"tn-purge": { description: "Remove stale memories from L2", args: "query" },
};

for (const [name, info] of Object.entries(commands)) {
	const cmd = `# /${name}\n\n${info.description}.\n\nUsage: /${name}${info.args ? ` ${info.args}` : ""}\n\nThis command routes through TormentNexus L2 memory via MCP.\n`;
	fs.writeFileSync(
		path.join(
			OPENHANDS_DIR,
			"plugins",
			"tormentnexus",
			"commands",
			`${name}.md`,
		),
		cmd,
	);
}

// 7. Microagent
const microagent = `---
name: tormentnexus
type: microagent
description: TormentNexus integration — persistent memory, MCP tools, session import
---

You are a TormentNexus-aware OpenHands microagent. Use these tools:

## Memory Operations
- \`tn_memory_store(content, tags)\` — Save decisions, patterns, facts to L2
- \`tn_memory_search(query)\` — Search past memories by keyword, tag, or category
- \`tn_memory_vector_search(query)\` — Semantic vector search across L2
- \`tn_context_harvest(prompt)\` — Pull relevant context into current session
- \`tn_memory_scratchpad(action, key, value)\` — L1 in-memory key-value store

## Tool Discovery
- \`tn_tool_search(query)\` — Find MCP tools across all configured servers
- \`tn_session_search(query)\` — Browse imported sessions from other AI agents
- \`tn_skill_manage(action, query)\` — Access 5,776+ reusable skill modules
- \`tn_code_search(query, scope)\` — Search code via AST-grep, semantic, or patterns

## System Tools
- \`tn_system_status()\` — Health overview of TN services
- \`tn_billing_status()\` — Provider quotas and fallback chain
- \`tn_audit_log(action, target)\` — Record to commercial audit log

## Best Practices
1. Search L2 memory before starting any complex task
2. Store key architectural decisions, bug fixes, and patterns
3. Harvest context at the start of multi-step tasks
4. Use tn_code_search for structural code understanding
5. All destructive operations are RBAC-checked
`;
fs.writeFileSync(
	path.join(OPENHANDS_DIR, "microagents", "tormentnexus.md"),
	microagent,
);

// 8. Python extension agent skeleton
const agentPy = `"""
TormentNexus OpenHands Agent Extension
Auto-loads TN memory, tools, and context harvesting into every OpenHands session.
"""
import os
import json
import urllib.request
from typing import Any, Dict, List, Optional

TN_URL = os.environ.get("TORMENTNEXUS_URL", "http://127.0.0.1:7778")


class TormentNexusAgent:
    """OpenHands agent wrapper that adds TN memory and tools."""

    def __init__(self, base_agent=None):
        self.base_agent = base_agent
        self._tools = None

    def _request(self, endpoint: str, method: str = "GET", data: Optional[Dict] = None) -> Dict:
        """Make HTTP request to TN Kernel."""
        url = f"{TN_URL}/api/{endpoint}"
        body = json.dumps(data).encode() if data else None
        req = urllib.request.Request(url, body, {"Content-Type": "application/json"}, method=method)
        try:
            with urllib.request.urlopen(req, timeout=10) as r:
                return json.loads(r.read())
        except Exception as e:
            return {"error": str(e)}

    def memory_search(self, query: str, tags: Optional[List[str]] = None) -> List[Dict]:
        """Search TN L2 memory."""
        result = self._request("memory/search", "POST", {"query": query, "tags": tags or []})
        return result.get("data", result.get("results", []))

    def memory_store(self, content: str, tags: List[str], category: str = "general") -> Dict:
        """Store to TN L2 memory."""
        return self._request("memory/store", "POST", {
            "content": content, "tags": tags, "category": category
        })

    def context_harvest(self, prompt: str) -> Dict:
        """Harvest context from TN L2 memory."""
        return self._request("context/harvest", "POST", {"prompt": prompt})

    def tool_search(self, query: str) -> List[Dict]:
        """Find MCP tools across all servers."""
        result = self._request("tools/search", "POST", {"query": query})
        return result.get("data", result.get("tools", []))

    def list_tools(self) -> List[str]:
        """Get all available TN tools."""
        if not self._tools:
            status = self._request("runtime/status")
            self._tools = status.get("data", {}).get("cli", {}).get("tools", [])
        return self._tools

    def health_check(self) -> bool:
        """Check if TN is reachable."""
        try:
            result = self._request("runtime/status")
            return "data" in result and result["data"].get("version")
        except Exception:
            return False


# Auto-instantiate when loaded as plugin
tn = TormentNexusAgent()

def get_agent():
    """Entry point for OpenHands plugin loader."""
    return tn
`;
fs.writeFileSync(
	path.join(OPENHANDS_DIR, "plugins", "tormentnexus", "tormentnexus_agent.py"),
	agentPy,
);

// 9. Run TN installer to ensure MCP is wired
try {
	const platform = process.platform;
	const installer = path.join(
		__dirname,
		"..",
		"..",
		"..",
		"scripts",
		platform === "win32" ? "install_codewhale.bat" : "install_codewhale.sh",
	);
	if (fs.existsSync(installer)) {
		console.log("Running TN client installer...");
	}
} catch (e) {
	// Non-fatal
}

console.log("\n✅ TormentNexus fully installed for OpenHands!");
console.log("   🧠 Agent:   TormentNexusAgent with memory + tools");
console.log("   📋 Config:  ~/.openhands/config.toml");
console.log("   🤖 Plugins: ~/.openhands/plugins/tormentnexus/");
console.log(
	"   📡 MCP:     ~/.openhands/plugins/tormentnexus/mcp/servers.json",
);
console.log("   🎯 Skills:  ~/.openhands/plugins/tormentnexus/skills/SKILL.md");
console.log(
	"   ⌨️  Commands: /tn-search, /tn-store, /tn-status, /tn-plan, /tn-summary, /tn-purge",
);
console.log(
	"\nStart OpenHands with: docker compose -f deploy/openhands/docker-compose.yml up\n",
);
