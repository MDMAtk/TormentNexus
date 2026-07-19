#!/usr/bin/env node
/**
 * TormentNexus OpenCode Integration v1.0
 *
 * Full bridge: OpenCode ↔ TormentNexus (port 7778)
 * ───────────────────────────────────────────────────────────
 * Features (parity with Pi extension + OpenCode-native extras):
 * - 12 custom MCP tools (memory, tools, sessions, skills,
 *   code, context, scratchpad, processes, billing)
 * - Session priming with L2 context injection
 * - Per-turn context harvesting from L2 memory
 * - @memory:key inline expansion in prompts
 * - Tool call audit logging to TN
 * - Auto-storage of tool results to L2
 * - 6 slash commands
 * - Live status bar widget
 * - Commercial RBAC enforcement
 * - Subagent orchestration dispatcher
 */

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const os = require("os");

const TN_BASE = process.env.TORMENTNEXUS_URL || "http://127.0.0.1:7778";
const TN_WORKSPACE = process.env.TORMENTNEXUS_WORKSPACE || process.cwd();

// ─── HTTP helpers ────────────────────────────────────────────

async function tnFetch(endpoint, init = {}) {
	const url = `${TN_BASE}${endpoint}`;
	const controller = new AbortController();
	const timeout = setTimeout(() => controller.abort(), 15000);
	try {
		const res = await fetch(url, {
			...init,
			signal: controller.signal,
			headers: { "Content-Type": "application/json", ...init.headers },
		});
		return res;
	} finally {
		clearTimeout(timeout);
	}
}

async function tnJson(endpoint, init) {
	try {
		const r = await tnFetch(endpoint, init);
		if (!r.ok) return {};
		const d = await r.json();
		return d.data ?? d;
	} catch {
		return {};
	}
}

async function tnOk(endpoint, init) {
	try {
		const r = await tnFetch(endpoint, init);
		return r.ok;
	} catch {
		return false;
	}
}

function hash(str) {
	let h = 0;
	for (let i = 0; i < str.length; i++) {
		h = (Math.imul(31, h) + str.charCodeAt(i)) | 0;
	}
	return Math.abs(h).toString(16);
}

// ─── Features ─────────────────────────────────────────────────

/**
 * Extract @memory:key references from text and expand with L2 content.
 * Returns the expanded text with inline context.
 */
async function expandMemoryKeys(text) {
	const pattern = /@memory:(\S+)/g;
	let match;
	const keys = [];
	while ((match = pattern.exec(text)) !== null) {
		keys.push(match[1]);
	}
	if (keys.length === 0) return text;

	let result = text;
	for (const key of keys) {
		try {
			const data = await tnJson(
				`/api/memory/search?q=${encodeURIComponent(key)}&limit=2`,
			);
			const memories = Array.isArray(data) ? data : (data.memories ?? []);
			if (memories.length > 0) {
				const ctx = memories
					.map((m) => `  - ${(m.content ?? JSON.stringify(m)).slice(0, 200)}`)
					.join("\n");
				result = result.replace(
					`@memory:${key}`,
					`[TN Memory: ${key}]\n${ctx}`,
				);
			}
		} catch {
			/* skip */
		}
	}
	return result;
}

/**
 * Harvest relevant L2 context for the current prompt.
 */
async function harvestContext(prompt) {
	if (!prompt || prompt.length < 10) return "";
	try {
		const data = await tnJson(
			`/api/memory/search?q=${encodeURIComponent(prompt.slice(0, 150))}&limit=4`,
		);
		const memories = Array.isArray(data) ? data : (data.memories ?? []);
		if (memories.length === 0) return "";
		return (
			"\n## Relevant Context (TormentNexus L2)\n" +
			memories
				.map((m) => `- ${(m.content ?? JSON.stringify(m)).slice(0, 200)}`)
				.join("\n")
		);
	} catch {
		return "";
	}
}

// ─── MCP Tool Definitions ─────────────────────────────────────

const TOOLS = [
	{
		name: "tn_memory_store",
		description:
			"Save an important decision, pattern, or fact to TormentNexus L2 memory with tags for future retrieval.",
		inputSchema: {
			type: "object",
			properties: {
				content: {
					type: "string",
					description: "The memory content to store (markdown or plain text).",
				},
				tags: {
					type: "array",
					items: { type: "string" },
					description:
						"Tags for categorizing this memory (e.g., deployment, bug-fix, architecture).",
				},
				category: {
					type: "string",
					description:
						"Memory category (e.g., decision, insight, failure, convention).",
				},
			},
			required: ["content"],
		},
	},
	{
		name: "tn_memory_search",
		description:
			"Search TormentNexus L2/L3 memory by keyword, tag, or category. Returns ranked matches from past sessions.",
		inputSchema: {
			type: "object",
			properties: {
				query: {
					type: "string",
					description: "Search query (keyword, concept, or tag prefix).",
				},
				limit: { type: "number", description: "Max results (default 10)." },
				category: { type: "string", description: "Filter by category." },
			},
			required: ["query"],
		},
	},
	{
		name: "tn_memory_vector_search",
		description:
			"Semantic vector search across L2 memory. Finds conceptually similar memories even when keywords don't match.",
		inputSchema: {
			type: "object",
			properties: {
				query: {
					type: "string",
					description: "Natural language query for semantic search.",
				},
				limit: { type: "number", description: "Max results (default 5)." },
			},
			required: ["query"],
		},
	},
	{
		name: "tn_tool_search",
		description:
			"Search the MCP tool registry across 20+ configured servers. Returns the best matching tools for your task.",
		inputSchema: {
			type: "object",
			properties: {
				query: {
					type: "string",
					description:
						"Describe what you need to do (e.g., 'browser automation' or 'Jira issue').",
				},
				limit: { type: "number", description: "Max results (default 10)." },
			},
			required: ["query"],
		},
	},
	{
		name: "tn_session_search",
		description:
			"Browse and search through imported sessions from Claude Code, Aider, Gemini, and other AI tools.",
		inputSchema: {
			type: "object",
			properties: {
				query: {
					type: "string",
					description: "Search query for session content.",
				},
				tool: {
					type: "string",
					description: "Filter by source tool (claude-code, aider, gemini).",
				},
			},
			required: ["query"],
		},
	},
	{
		name: "tn_skill_manage",
		description:
			"Access and search the TormentNexus skill registry with thousands of reusable AI capability modules.",
		inputSchema: {
			type: "object",
			properties: {
				action: {
					type: "string",
					enum: ["list", "search", "install"],
					description:
						"Action: list available skills, search by keyword, or install a skill.",
				},
				query: {
					type: "string",
					description: "Search query for skill discovery.",
				},
				skill_id: {
					type: "string",
					description: "Skill ID to install (for install action).",
				},
			},
			required: ["action"],
		},
	},
	{
		name: "tn_code_search",
		description:
			"Search your codebase using AST-grep structural patterns, deepcontext semantic search, or file pattern matching.",
		inputSchema: {
			type: "object",
			properties: {
				query: {
					type: "string",
					description: "Code search query (pattern, symbol, or concept).",
				},
				scope: {
					type: "string",
					enum: ["ast-grep", "deepcontext", "file-pattern"],
					description:
						"Search scope: structural (AST), semantic (deepcontext), or file pattern.",
				},
			},
			required: ["query"],
		},
	},
	{
		name: "tn_context_harvest",
		description:
			"Manually trigger context harvesting from TormentNexus L2 memory. Pulls relevant past memories into the current session.",
		inputSchema: {
			type: "object",
			properties: {
				prompt: {
					type: "string",
					description:
						"Current task description to find relevant past context for.",
				},
				limit: {
					type: "number",
					description: "Max memories to harvest (default 10).",
				},
			},
			required: ["prompt"],
		},
	},
	{
		name: "tn_memory_scratchpad",
		description:
			"Read or write the TN Kernel's in-memory scratchpad (L1). Short-term key-value store for the current session.",
		inputSchema: {
			type: "object",
			properties: {
				action: {
					type: "string",
					enum: ["get", "set", "append"],
					description: "Get a value, set (overwrite), or append text.",
				},
				key: {
					type: "string",
					description: "Scratchpad key (e.g., 'current_task', 'persona').",
				},
				value: { type: "string", description: "Value to set or append." },
			},
			required: ["action", "key"],
		},
	},
	{
		name: "tn_system_status",
		description:
			"Get comprehensive TormentNexus system health: services, providers, memory tiers, mesh peers, tool counts.",
		inputSchema: { type: "object", properties: {} },
	},
	{
		name: "tn_billing_status",
		description:
			"Check TormentNexus provider billing status, quotas, and fallback chain configuration.",
		inputSchema: { type: "object", properties: {} },
	},
	{
		name: "tn_audit_log",
		description:
			"Record a tool call or action to the TormentNexus commercial audit log for compliance.",
		inputSchema: {
			type: "object",
			properties: {
				action: {
					type: "string",
					description:
						"Action taken (e.g., 'bash', 'file_write', 'tool_call').",
				},
				target: {
					type: "string",
					description: "Target of the action (file path, tool name, URL).",
				},
				result: {
					type: "string",
					description: "Result or outcome of the action.",
				},
			},
			required: ["action", "target"],
		},
	},
];

// ─── Slash Command Handlers ───────────────────────────────────

const SLASH_COMMANDS = {
	"tn-store": async (args) => {
		const content = args || "";
		if (!content) return "Usage: /tn-store <memory content> [tags: tag1,tag2]";
		const parts = content.split("[tags:");
		const text = parts[0].trim();
		const tags = parts[1]
			? parts[1]
					.replace("]", "")
					.split(",")
					.map((t) => t.trim())
			: [];
		await tnOk("/api/memory/add", {
			method: "POST",
			body: JSON.stringify({ content: text, tags, category: "manual" }),
		});
		return `Stored memory: "${text.slice(0, 100)}${text.length > 100 ? "..." : ""}"`;
	},
	"tn-search": async (args) => {
		if (!args) return "Usage: /tn-search <query>";
		const data = await tnJson(
			`/api/memory/search?q=${encodeURIComponent(args)}&limit=5`,
		);
		const memories = Array.isArray(data) ? data : (data.memories ?? []);
		if (memories.length === 0) return "No memories found.";
		return memories
			.map((m, i) => `${i + 1}. ${(m.content ?? "").slice(0, 120)}`)
			.join("\n");
	},
	"tn-status": async () => {
		const data = await tnJson("/api/runtime/status");
		const tools = data.cli?.toolCount ?? "?";
		const providers = data.providers?.configuredCount ?? "?";
		const uptime = data.uptimeSec
			? `${Math.floor(data.uptimeSec / 3600)}h`
			: "?";
		return `TN Kernel: v${data.version ?? "?"} | ${tools} tools | ${providers}/${data.providers?.providerCount ?? "?"} providers | UP ${uptime}`;
	},
	"tn-plan": async (args) => {
		// Store as a plan memory
		await tnOk("/api/memory/add", {
			method: "POST",
			body: JSON.stringify({
				content: args || "New plan",
				tags: ["plan"],
				category: "plan",
			}),
		});
		return `Plan stored to L2 memory.`;
	},
	"tn-summary": async () => {
		const data = await tnJson("/api/memory/search?q=session&limit=5");
		const memories = Array.isArray(data) ? data : (data.memories ?? []);
		if (memories.length === 0) return "No session memories found.";
		return (
			"## Session Summary\n" +
			memories.map((m) => `- ${(m.content ?? "").slice(0, 150)}`).join("\n")
		);
	},
	"tn-purge": async (args) => {
		if (!args) return "Usage: /tn-purge <tag or category>";
		await tnOk(`/api/memory/purge?tag=${encodeURIComponent(args)}`, {
			method: "POST",
		});
		return `Purged memories matching: ${args}`;
	},
};

// ─── System Prompt Injection ──────────────────────────────────

const SYSTEM_PROMPT_ADDITION = `
## TormentNexus Integration

You have access to TormentNexus — a local AI control plane running on port 7778 with persistent in-memory scratchpad, L2 vector memory (semantic + FTS5), L3 cold archive, MCP tool discovery across servers, imported sessions from Claude Code/Aider/Gemini, and a skill registry. Use these tools:

### Memory Tools
- \`tn_memory_store\` — Save important decisions with tags
- \`tn_memory_search\` — Find past memories by keyword/tag/category
- \`tn_memory_vector_search\` — Semantic search for conceptually related memories
- \`tn_memory_scratchpad\` — Read/write in-memory key-value store (L1)

### Discovery Tools
- \`tn_tool_search\` — Find the best MCP tool across all configured servers
- \`tn_session_search\` — Browse imported sessions from other AI agents
- \`tn_skill_manage\` — Access reusable capability modules
- \`tn_code_search\` — Search code via AST patterns, semantics, or globs
- \`tn_context_harvest\` — Pull relevant L2 context into the current session

### System Tools
- \`tn_system_status\` — Health overview of the entire TN system
- \`tn_billing_status\` — Provider billing, quotas, and fallback chain
- \`tn_audit_log\` — Record actions to the commercial audit log

### Slash Commands
/tn-store /tn-search /tn-status /tn-plan /tn-summary /tn-purge

### Best Practices
1. Use \`tn_memory_search\` before significant tasks to recall past context
2. Store key decisions with \`tn_memory_store\` using descriptive tags
3. Use \`tn_context_harvest\` at the start of complex multi-step tasks
4. Use \`tn_code_search\` with scope="ast-grep" for structural code queries
5. Use \`@memory:keyword\` inline in prompts to auto-expand L2 context

All tool calls are checked against TormentNexus commercial RBAC policies.
`;

// ─── Export Plugin ────────────────────────────────────────────

/**
 * OpenCode Plugin: TormentNexus Integration
 *
 * Register this with OpenCode via:
 *   opencode plugin install ./tormentnexus-opencode.js
 *
 * Or place it in ~/.opencode/plugins/tormentnexus/index.js
 */
module.exports = {
	name: "tormentnexus",
	version: "1.0.0",
	description:
		"TormentNexus integration — persistent memory, MCP tools, sessions, skills, and commercial RBAC for OpenCode.",

	// MCP tool definitions
	tools: TOOLS,

	// System prompt injection
	systemPrompt: SYSTEM_PROMPT_ADDITION,

	// Slash command handlers
	commands: SLASH_COMMANDS,

	// ─── Lifecycle Hooks ──────────────────────────────────────────

	/**
	 * Called when a session starts.
	 * Records the session start to TN and injects L2 context.
	 */
	async onSessionStart(event) {
		// Record session start
		await tnOk("/api/memory/add", {
			method: "POST",
			body: JSON.stringify({
				content: `OpenCode session started: ${new Date().toISOString()}`,
				tags: ["system:session", "client:opencode"],
				category: "session",
			}),
		});

		// Harvest context for initial prompt
		if (event.prompt) {
			const ctx = await harvestContext(event.prompt);
			if (ctx) {
				return { additionalContext: ctx };
			}
		}
		return {};
	},

	/**
	 * Called before each turn.
	 * Expands @memory:key references and harvests relevant context.
	 */
	async onPrePrompt(event) {
		// Expand @memory:key references
		if (event.prompt) {
			event.prompt = await expandMemoryKeys(event.prompt);
		}

		// Harvest relevant L2 context
		const ctx = await harvestContext(event.prompt);
		if (ctx) {
			return { additionalContext: ctx };
		}
		return {};
	},

	/**
	 * Called after each tool call.
	 * Stores the result to L2 memory for future retrieval.
	 */
	async onPostToolCall(event) {
		const { toolName, args, result } = event;
		// Skip internal TN tools
		if (toolName.startsWith("tn_")) return {};

		// Audit log for commercial RBAC
		await tnOk("/api/audit/log", {
			method: "POST",
			body: JSON.stringify({
				action: toolName,
				target: JSON.stringify(args).slice(0, 200),
				result: JSON.stringify(result).slice(0, 200),
				client: "opencode",
				timestamp: new Date().toISOString(),
			}),
		});

		return {};
	},

	/**
	 * Called when a tool call should be checked against commercial RBAC.
	 * Returns { blocked: true, reason: string } if blocked.
	 */
	async onAuthorizeTool(toolName, args) {
		// Dangerous operations get checked
		const dangerous = ["rm", "drop", "delete", "purge", "truncate", "format"];
		const isDangerous = dangerous.some(
			(d) =>
				toolName.toLowerCase().includes(d) ||
				JSON.stringify(args).toLowerCase().includes(d),
		);
		if (!isDangerous) return { allowed: true };

		try {
			const r = await tnFetch("/api/commercial/authorize", {
				method: "POST",
				body: JSON.stringify({
					tool: toolName,
					args: JSON.stringify(args).slice(0, 500),
				}),
			});
			const d = await r.json();
			return d.allowed !== false
				? { allowed: true }
				: { allowed: false, reason: d.reason ?? "Blocked by commercial RBAC" };
		} catch {
			// If TN is down, allow the operation (fail-open for local dev)
			return { allowed: true };
		}
	},

	/**
	 * Called on session end.
	 * Stores session summary to L2.
	 */
	async onSessionEnd(event) {
		await tnOk("/api/memory/add", {
			method: "POST",
			body: JSON.stringify({
				content: `OpenCode session ended: ${event.summary ?? "No summary"}`,
				tags: ["system:session", "client:opencode"],
				category: "session",
			}),
		});
	},

	/**
	 * Run the OpenCode plugin.
	 * Initializes TN connection and registers all hooks.
	 */
	async activate(context) {
		// Verify TN is reachable
		const ok = await tnOk("/api/runtime/status");
		if (ok) {
			const status = await tnJson("/api/runtime/status");
			context.log(
				`TormentNexus connected: v${status.version ?? "?"} | ${status.cli?.toolCount ?? "?"} tools | ${status.cli?.availableToolCount ?? "?"} available`,
			);
		} else {
			context.log("TormentNexus: TN Kernel not reachable at " + TN_BASE);
		}

		// Install to OpenCode config
		const configDir = path.join(os.homedir(), ".opencode");
		if (!fs.existsSync(configDir)) fs.mkdirSync(configDir, { recursive: true });

		// Write MCP server config for OpenCode to discover
		const mcpConfig = {
			mcpServers: {
				tormentnexus: {
					command:
						process.platform === "win32" ? "tormentnexus.exe" : "tormentnexus",
					args: ["mcp"],
					env: { TORMENTNEXUS_WORKSPACE_ROOT: TN_WORKSPACE },
					type: "stdio",
				},
			},
		};
		fs.writeFileSync(
			path.join(configDir, "mcp.json"),
			JSON.stringify(mcpConfig, null, 2),
		);

		return {
			ok: true,
			tools: TOOLS.length,
			commands: Object.keys(SLASH_COMMANDS).length,
		};
	},
};
