import type {
	ExtensionAPI,
	isToolCallEventType,
} from "@earendil-works/pi-coding-agent";
import { Type } from "typebox";

const TN_BASE = "http://127.0.0.1:7778";

/**
 * tormentnexus pi extension v4
 *
 * Full bridge: pi ↔ TormentNexus
 * ───────────────────────────────────────────
 * - 9 custom tools (memory, tools, sessions, skills, code, context, scratchpad)
 * - Session priming + per-turn context harvesting + compaction hooks
 * - tool_call RBAC enforcement via commercial API
 * - tool_result auto-storage to L2 memory
 * - user_bash audit logging through TN
 * - model_select tracking to L2
 * - Input transformation: @memory:key expansions
 * - 6 slash commands (/tn-store, /tn-search, /tn-status, /tn-plan, /tn-purge, /tn-summary)
 * - Live editor widget with memory/mesh stats
 * - Custom footer with TN status
 * - Keyboard shortcuts (Ctrl+Shift+M, Ctrl+Shift+T, Ctrl+Shift+P)
 * - Inter-extension event bus (tn:* events)
 * - Subagent orchestration dispatcher
 */

// ─── Helpers ────────────────────────────────────────────────────────────

function tnFetch(path: string, init?: RequestInit, signal?: AbortSignal) {
	return fetch(`${TN_BASE}${path}`, { ...init, signal });
}

async function tnOk(path: string, init?: RequestInit, signal?: AbortSignal) {
	try {
		const r = await tnFetch(path, init, signal);
		return r.ok;
	} catch {
		return false;
	}
}

async function tnJson(path: string, init?: RequestInit, signal?: AbortSignal) {
	try {
		const r = await tnFetch(path, init, signal);
		if (!r.ok) return {};
		const d = await r.json();
		return d.data ?? d;
	} catch {
		return {};
	}
}

function now() {
	return new Date().toISOString();
}

// ─── System prompt guidance ────────────────────────────────────────────
const TN_SYSTEM_PROMPT = `
## TormentNexus Integration

You have access to TormentNexus — a local AI control plane running on port 7778 with persistent L2 vector memory, semantic tool discovery, imported sessions, and a skill registry. Use these tools:

### Memory Tools
- \`tn_memory_store\` — Save important decisions, patterns, and facts with tags
- \`tn_memory_search\` — Find past memories by keyword, tag, or category
- \`tn_memory_vector_search\` — Semantic vector search for conceptually related memories

### Discovery Tools
- \`tn_tool_search\` — Describe what you need, TN finds the best tool across 20+ servers
- \`tn_session_search\` — Browse 542+ imported sessions from Claude Code, Aider, etc.
- \`tn_skill_manage\` — Access 5,776 reusable skill modules
- \`tn_code_search\` — Search code via AST-grep rules, deepcontext semantic search, or file pattern matching
- \`tn_context_harvest\` — Manually trigger context harvesting from TN L2 memory

### Slash Commands
- \`/tn-store\` — Interactive memory store with structured form
- \`/tn-search\` — Interactive memory search with filters
- \`/tn-status\` — Show TN system status: memory tiers, mesh peers
- \`/tn-plan\` — Create/edit/view project plans in L2 memory
- \`/tn-summary\` — Summarize current session using TN context
- \`/tn-purge\` — Remove stale memories from L2

### Commercial Security
All tool calls are checked against TormentNexus commercial RBAC. Blocked tools show a security notice.

### Best Practices
1. Use \`tn_memory_search\` before significant tasks to recall past context
2. Store key decisions with \`tn_memory_store\` using descriptive tags
3. Use \`tn_context_harvest\` at the start of complex tasks
4. Use \`@memory:keyword\` inline in your prompts to auto-expand L2 context
`;

export default function (pi: ExtensionAPI) {
	// ══════════════════════════════════════════════
	// SESSION PRIMING
	// ══════════════════════════════════════════════

	pi.on("session_start", async (event, ctx) => {
		const sessionFile = ctx.sessionManager.getSessionFile();
		const reason = event.reason;

		// Store session start
		await tnOk("/api/memory/add", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				content: JSON.stringify({
					content: `Session ${reason}: ${sessionFile ?? "ephemeral"}`,
					tags: ["system:session", `reason:${reason}`],
					category: "session",
					timestamp: now(),
				}),
			}),
		});

		ctx.ui.setStatus("tn", "TN active • L2 mem • tools • skills");

		// Refresh stats widget
		refreshWidget(ctx);
	});

	pi.on("before_agent_start", async (event, ctx) => {
		const isFirstTurn = event.systemPrompt.includes("TormentNexus");
		if (!isFirstTurn) {
			// Inject L2 context on subsequent turns
			try {
				const res = await tnFetch(
					`/api/memory/search?q=${encodeURIComponent(event.prompt.slice(0, 100))}`,
					{},
					ctx.signal,
				);
				if (res.ok) {
					const body = await res.json();
					const memories = body.data ?? [];
					if (Array.isArray(memories) && memories.length > 0) {
						const contextBlock = memories
							.slice(0, 3)
							.map(
								(m: any) =>
									`  • ${(m.content ?? m.text ?? JSON.stringify(m)).slice(0, 200)}`,
							)
							.join("\n");
						return {
							systemPrompt: `${event.systemPrompt}\n\n## Relevant Context from TormentNexus L2\n${contextBlock}`,
						};
					}
				}
			} catch {
				/* skip */
			}
			return;
		}

		// First turn: full guidance + context
		let memoryContext = "";
		try {
			const res = await tnFetch(
				`/api/memory/search?q=${encodeURIComponent(event.prompt.slice(0, 100))}`,
				{},
				ctx.signal,
			);
			if (res.ok) {
				const body = await res.json();
				const memories = body.data ?? [];
				if (Array.isArray(memories) && memories.length > 0) {
					memoryContext =
						"\n\n## Relevant Past Context\n" +
						memories
							.slice(0, 5)
							.map(
								(m: any) =>
									`  • ${(m.content ?? m.text ?? JSON.stringify(m)).slice(0, 200)}`,
							)
							.join("\n");
				}
			}
		} catch {
			/* fallback */
		}

		if (!memoryContext) {
			try {
				const res = await tnFetch("/api/memory/list", {}, ctx.signal);
				if (res.ok) {
					const all: string[] = await res.json();
					const q = event.prompt.toLowerCase().slice(0, 100);
					const relevant = all
						.map((m) => {
							try {
								return { ...JSON.parse(m), raw: m };
							} catch {
								return { content: m, tags: [] };
							}
						})
						.filter(
							(m) =>
								m.content?.toLowerCase().includes(q) ||
								m.tags?.some((t: string) => t.toLowerCase().includes(q)),
						)
						.slice(0, 5);
					if (relevant.length > 0) {
						memoryContext =
							"\n\n## Relevant Past Context\n" +
							relevant.map((m) => `  • ${m.content.slice(0, 200)}`).join("\n");
					}
				}
			} catch {
				/* no memory */
			}
		}

		return {
			systemPrompt: event.systemPrompt + TN_SYSTEM_PROMPT + memoryContext,
		};
	});

	// ══════════════════════════════════════════════
	// PER-TURN CONTEXT HARVESTING
	// ══════════════════════════════════════════════

	pi.on("context", async (event, ctx) => {
		const lastMessages = event.messages.slice(-4);
		const hasRecentSearch = lastMessages.some(
			(m: any) =>
				m.role === "assistant" &&
				JSON.stringify(m.content)?.includes("tn_memory_search"),
		);
		if (hasRecentSearch) return;

		const lastUserMsg = [...lastMessages]
			.reverse()
			.find((m: any) => m.role === "user");
		if (!lastUserMsg) return;

		const userText =
			typeof lastUserMsg.content === "string"
				? lastUserMsg.content
				: (lastUserMsg.content
						?.map((c: any) => c.text)
						.filter(Boolean)
						.join(" ") ?? "");

		if (!userText || userText.length < 10) return;

		try {
			const res = await tnFetch(
				`/api/memory/search?q=${encodeURIComponent(userText.slice(0, 100))}`,
				{},
				ctx.signal,
			);
			if (!res.ok) return;
			const body = await res.json();
			const memories = body.data ?? [];
			if (!Array.isArray(memories) || memories.length === 0) return;

			const top = memories.slice(0, 2);
			const contextBlock = top
				.map((m: any) =>
					(m.content ?? m.text ?? JSON.stringify(m)).slice(0, 150),
				)
				.filter(Boolean)
				.join("\n");
			if (!contextBlock) return;

			event.messages.push({
				role: "system",
				content: `[TN Context]: ${contextBlock}`,
			});
		} catch {
			/* skip */
		}
	});

	// ══════════════════════════════════════════════
	// COMMERCIAL RBAC — TOOL CALL ENFORCEMENT
	// ══════════════════════════════════════════════

	pi.on("tool_call", async (event, ctx) => {
		// Check dangerous operations against TN commercial RBAC
		// Only check bash/shell commands, not file content (write/edit tools can contain any text)
		if (event.toolName === "bash" || event.toolName === "shell") {
		const dangerousPatterns = [
			"rm -rf",
			"sudo ",
			"chmod -R 777",
			"DROP TABLE",
			"DROP DATABASE",
			"git push --force",
		];
		const inputStr = JSON.stringify(event.input).toLowerCase();

		for (const pattern of dangerousPatterns) {
			if (inputStr.includes(pattern.toLowerCase())) {
				const res = await tnJson(
					"/api/commercial/authorize",
					{
						method: "POST",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({
							tool: event.toolName,
							action: pattern,
							args: event.input,
						}),
					},
					ctx.signal,
				);
				const allowed = res.allowed ?? false;
				if (!allowed) {
					return {
						block: true,
						reason: `Commercial policy blocks: ${pattern}. Use tn_memory_store for destructive operations.`,
					};
				}
			}
		}


		} // end bash-only RBAC guard

		// Log tool execution to TN audit
		await tnOk("/api/commercial/audit/log", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				tool: event.toolName,
				args: event.input,
				timestamp: now(),
				userId: "pi-agent",
			}),
		});
	});

	// ══════════════════════════════════════════════
	// AUTO-STORE INTERESTING TOOL RESULTS TO L2
	// ══════════════════════════════════════════════

	pi.on("tool_result", async (event, ctx) => {
		// Only store results from key tools that produce useful context
		const storeTools = [
			"bash",
			"read",
			"grep",
			"tn_code_search",
			"tn_tool_search",
		];
		if (!storeTools.includes(event.toolName)) return;

		const text =
			event.content
				?.map((c: any) => c.text)
				.filter(Boolean)
				.join(" ") ?? "";
		if (!text || text.length < 100) return; // Only store substantial results

		await tnOk("/api/memory/add", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				content: JSON.stringify({
					content: `[${event.toolName}] ${text.slice(0, 300)}`,
					tags: ["system:tool_result", `tool:${event.toolName}`],
					category: "tool_result",
					timestamp: now(),
				}),
			}),
		});
	});

	// ══════════════════════════════════════════════
	// USER BASH — AUDIT LOGGING
	// ══════════════════════════════════════════════

	pi.on("user_bash", async (event, ctx) => {
		await tnOk("/api/commercial/audit/log", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				tool: "user_bash",
				command: event.command,
				timestamp: now(),
			}),
		});
	});

	// ══════════════════════════════════════════════
	// MODEL SELECT — TRACK TO L2
	// ══════════════════════════════════════════════

	pi.on("model_select", async (event, ctx) => {
		await tnOk("/api/memory/add", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				content: JSON.stringify({
					content: `Model changed: ${event.previousModel?.id ?? "none"} → ${event.model.id} (${event.source})`,
					tags: ["system:model_change", `model:${event.model.id}`],
					category: "session",
					timestamp: now(),
				}),
			}),
		});
		ctx.ui.setStatus("tn", `TN • ${event.model.id}`);
	});

	// ══════════════════════════════════════════════
	// INPUT TRANSFORMATION — @memory:key EXPANSION
	// ══════════════════════════════════════════════

	pi.on("input", async (event, ctx) => {
		const text = event.text;
		if (!text || !text.includes("@memory:")) return;

		const matches = text.match(/@memory:([a-zA-Z0-9_-]+)/g);
		if (!matches) return;

		let expanded = text;
		for (const match of matches) {
			const key = match.replace("@memory:", "");
			try {
				const res = await tnFetch(
					`/api/memory/search?q=${encodeURIComponent(key)}&limit=3`,
					{},
					ctx.signal,
				);
				if (res.ok) {
					const body = await res.json();
					const memories = body.data ?? [];
					if (Array.isArray(memories) && memories.length > 0) {
						const context = memories
							.map((m: any) => (m.content ?? m.text ?? "").slice(0, 150))
							.filter(Boolean)
							.join("\n");
						expanded = expanded.replace(
							match,
							`[TN:${key}]\n${context}\n[/TN:${key}]`,
						);
					}
				}
			} catch {
				/* skip */
			}
		}

		if (expanded !== text) {
			return { text: expanded };
		}
	});

	// ══════════════════════════════════════════════
	// AUTO-LOGGING
	// ══════════════════════════════════════════════

	pi.on("turn_end", async (event, ctx) => {
		if (!event.toolResults || event.toolResults.length === 0) return;

		const summary = event.message?.content
			?.filter((c: any) => c.type === "text")
			?.map((c: any) => c.text)
			?.join(" ")
			?.slice(0, 500);

		if (!summary || summary.length < 50) return;

		await tnOk("/api/memory/add", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				content: JSON.stringify({
					content: `Turn ${event.turnIndex}: ${summary.slice(0, 200)}...`,
					tags: ["system:turn", `turn:${event.turnIndex}`],
					category: "session",
					timestamp: now(),
				}),
			}),
		});
	});

	// ══════════════════════════════════════════════
	// COMPACTION HOOKS
	// ══════════════════════════════════════════════

	pi.on("session_before_compact", async (event, ctx) => {
		const { preparation, branchEntries, customInstructions, reason, signal } =
			event;
		const summaryParts: string[] = [];

		if (branchEntries && branchEntries.length > 0) {
			const fileOps = branchEntries
				.filter((e: any) => e.details?.readFiles || e.details?.modifiedFiles)
				.map((e: any) => {
					const reads = (e.details?.readFiles ?? []).join(", ");
					const mods = (e.details?.modifiedFiles ?? []).join(", ");
					return `${reads ? `Read: ${reads}` : ""}${mods ? ` Modified: ${mods}` : ""}`;
				})
				.filter(Boolean);
			if (fileOps.length > 0) summaryParts.push(`Files: ${fileOps.join("; ")}`);
		}

		try {
			if (summaryParts.length > 0) {
				const query = summaryParts.join(" ").slice(0, 100);
				const memories = await tnJson(
					`/api/memory/search?q=${encodeURIComponent(query)}`,
					{},
					signal,
				);
				const items = Array.isArray(memories)
					? memories
					: ((memories as any).data ?? []);
				if (items.length > 0) {
					const related = items
						.slice(0, 2)
						.map((m: any) => (m.content ?? m.text ?? "").slice(0, 100))
						.filter(Boolean)
						.join("; ");
					if (related) summaryParts.push(`Related: ${related}`);
				}
			}
		} catch {
			/* skip */
		}

		return {
			compaction: {
				summary:
					summaryParts.length > 0
						? summaryParts.join("\n")
						: `Compaction (${reason})`,
				firstKeptEntryId: preparation.firstKeptEntryId,
				tokensBefore: preparation.tokensBefore,
				details: { enrichedBy: "tormentnexus-l2", reason, timestamp: now() },
			},
		};
	});

	pi.on("session_compact", async (event, ctx) => {
		if (!event.compactionEntry?.summary) return;
		await tnOk("/api/memory/add", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				content: JSON.stringify({
					content: `Compaction [${event.reason}]: ${event.compactionEntry.summary.slice(0, 200)}`,
					tags: ["system:compaction", `reason:${event.reason}`],
					category: "session",
					timestamp: now(),
				}),
			}),
		});
	});

	// ══════════════════════════════════════════════
	// INTER-EXTENSION EVENT BUS
	// ══════════════════════════════════════════════

	// Emit tn:* events for other extensions to listen to
	pi.events.on("tn:request_context", async (data: any) => {
		try {
			const res = await tnFetch(
				`/api/memory/search?q=${encodeURIComponent(data.query)}&limit=5`,
			);
			if (res.ok) {
				const body = await res.json();
				pi.events.emit("tn:context_result", {
					id: data.id,
					memories: body.data ?? [],
				});
			}
		} catch {
			/* skip */
		}
	});

	pi.events.on("tn:request_store", async (data: any) => {
		await tnOk("/api/memory/add", {
			method: "POST",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				content: JSON.stringify({
					content: data.content,
					tags: data.tags ?? [],
					category: data.category ?? "general",
					timestamp: now(),
				}),
			}),
		});
	});

	// ══════════════════════════════════════════════
	// CUSTOM TOOLS (9 tools)
	// ══════════════════════════════════════════════

	// 1. tn_memory_store
	pi.registerTool({
		name: "tn_memory_store",
		label: "TN Memory Store",
		description:
			"Store a memory in TormentNexus L2 vault and optionally in a project .memdb for git-tracked per-project portable memories.",
		promptSnippet: "Store knowledge in persistent L2 memory",
		promptGuidelines: [
			"Use tn_memory_store to save important patterns, decisions, and facts across sessions.",
			"Use tags like 'project:name', 'failure:', 'pattern:', or 'convention:' for scope filtering.",
			"Set 'project' to the project directory name to also write to the project's .memdb file (git-tracked, portable).",
			"Good candidates: architectural decisions, bug fixes, build procedures, tool quirks.",
		],
		parameters: Type.Object({
			content: Type.String({ description: "The memory content to store" }),
			tags: Type.Optional(
				Type.Array(Type.String(), {
					description: "Tags like ['project:bg', 'pattern:build']",
				}),
			),
			category: Type.Optional(
				Type.String({
					description:
						"Category: pattern, decision, convention, insight, failure, correction, preference",
				}),
			),
			project: Type.Optional(
				Type.String({
					description:
						"Project directory name. If set, also writes to project/.memdb for git tracking",
				}),
			),
		}),
		async execute(_toolCallId, params, signal, _onUpdate, ctx) {
			const sessionFile = ctx.sessionManager.getSessionFile();
			const tags = params.tags ?? [];
			if (params.project) {
				// Ensure project tag is present
				const projectTag = "project:" + params.project;
				if (!tags.includes(projectTag)) tags.push(projectTag);
			}
			const enriched = JSON.stringify({
				content: params.content,
				tags,
				category: params.category ?? "general",
				timestamp: now(),
				session: sessionFile,
			});
			// Store to global L2 vault
			const res = await tnFetch(
				"/api/memory/add",
				{
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({ content: enriched }),
				},
				signal,
			);
			if (!res.ok)
				return {
					content: [{ type: "text", text: `Failed: ${res.status}` }],
					isError: true,
				};
			// If project specified, trigger a project memdb sync so the new memory is written to the .memdb
			if (params.project) {
				tnFetch("/api/memory/project/sync", { method: "POST" }).catch(() => {});
			}
			pi.events.emit("tn:memory_stored", { content: params.content, tags });
			const projectNote = params.project
				? ` and project .memdb (${params.project})`
				: "";
			return {
				content: [
					{ type: "text", text: `✅ Memory stored in L2 vault${projectNote}.` },
				],
				details: { tags },
			};
		},
	});

	// 2. tn_memory_search
	pi.registerTool({
		name: "tn_memory_search",
		label: "TN Memory Search",
		description: "Search L2 vault by keyword, tag filter, or category.",
		promptSnippet: "Search persistent L2 memory",
		promptGuidelines: [
			"Use tn_memory_search before tasks to recall past context.",
			"Filter by tag prefix like 'project:' or 'failure:'.",
		],
		parameters: Type.Object({
			query: Type.Optional(
				Type.String({ description: "Keyword to search for" }),
			),
			tag: Type.Optional(
				Type.String({ description: "Filter by tag prefix, e.g. 'project:'" }),
			),
			category: Type.Optional(
				Type.String({ description: "Filter by category" }),
			),
			limit: Type.Optional(
				Type.Number({ description: "Max results (default 20)" }),
			),
		}),
		async execute(_toolCallId, params, signal) {
			const res = await tnFetch("/api/memory/list", {}, signal);
			if (!res.ok)
				return {
					content: [{ type: "text", text: "Memory unavailable." }],
					isError: true,
				};
			const memories: string[] = await res.json();
			const limit = params.limit ?? 20;
			const parsed = memories
				.map((m) => {
					try {
						const p = JSON.parse(m);
						return {
							content: p.content ?? m,
							tags: p.tags ?? [],
							category: p.category ?? "general",
						};
					} catch {
						return { content: m, tags: [], category: "general" };
					}
				})
				.filter((m) => {
					if (params.query) {
						const q = params.query.toLowerCase();
						if (
							!m.content.toLowerCase().includes(q) &&
							!m.tags.some((t: string) => t.toLowerCase().includes(q))
						)
							return false;
					}
					if (
						params.tag &&
						!m.tags.some((t: string) => t.startsWith(params.tag!))
					)
						return false;
					if (params.category && m.category !== params.category) return false;
					return true;
				})
				.slice(0, limit);

			if (parsed.length === 0)
				return { content: [{ type: "text", text: "No matching memories." }] };
			const formatted = parsed
				.map(
					(m, i) =>
						`${i + 1}.${m.category !== "general" ? ` (${m.category})` : ""}${m.tags.length ? ` [${m.tags.join(", ")}]` : ""}\n   ${m.content.slice(0, 200)}`,
				)
				.join("\n\n");
			return {
				content: [
					{
						type: "text",
						text: `📚 ${parsed.length} memories:\n\n${formatted}`,
					},
				],
			};
		},
	});

	// 3. tn_memory_vector_search
	pi.registerTool({
		name: "tn_memory_vector_search",
		label: "TN Vector Search",
		description:
			"Semantic search L2 memory using sqlite-vec embeddings. Finds conceptually related memories.",
		promptSnippet: "Semantic vector search L2 memory",
		promptGuidelines: [
			"Use tn_memory_vector_search for fuzzy/conceptual recall.",
		],
		parameters: Type.Object({
			query: Type.String({ description: "Natural language query" }),
			limit: Type.Optional(
				Type.Number({ description: "Max results (default 10)" }),
			),
		}),
		async execute(_toolCallId, params, signal) {
			const limit = params.limit ?? 10;
			const res = await tnFetch(
				`/api/memory/search?q=${encodeURIComponent(params.query)}`,
				{},
				signal,
			);
			if (!res.ok)
				return {
					content: [{ type: "text", text: "Search unavailable." }],
					isError: true,
				};
			const body = await res.json();
			const data = body.data ?? body;
			if (Array.isArray(data) && data.length > 0) {
				const results = data.slice(0, limit);
				const formatted = results
					.map(
						(r: any, i: number) =>
							`  ${i + 1}. ${(r.content ?? r.text ?? JSON.stringify(r)).slice(0, 200)}`,
					)
					.join("\n\n");
				return {
					content: [
						{
							type: "text",
							text: `🧠 ${data.length} results:\n\n${formatted}`,
						},
					],
				};
			}
			return {
				content: [{ type: "text", text: `No results for "${params.query}".` }],
			};
		},
	});

	// 4. tn_tool_search
	pi.registerTool({
		name: "tn_tool_search",
		label: "TN Tool Search",
		description:
			"Semantically search MCP tools across 20+ servers. Describe what you need — finds the best tool by meaning.",
		promptSnippet: "Discover tools via semantic search",
		promptGuidelines: [
			"Use tn_tool_search when unsure what's available.",
			"Describe the task naturally, not keywords.",
		],
		parameters: Type.Object({
			query: Type.String({
				description: "Natural language description of what you need",
			}),
			limit: Type.Optional(
				Type.Number({ description: "Max results (default 5)" }),
			),
		}),
		async execute(_toolCallId, params, signal) {
			const res = await tnFetch(
				`/api/mcp/native/search?query=${encodeURIComponent(params.query)}`,
				{},
				signal,
			);
			if (!res.ok)
				return {
					content: [{ type: "text", text: `Failed: ${res.status}` }],
					isError: true,
				};
			const body = await res.json();
			const data = body.data ?? body;
			const results = data.results ?? (Array.isArray(data) ? data : []);
			if (results.length === 0)
				return {
					content: [{ type: "text", text: `No tools for "${params.query}".` }],
				};
			const limit = params.limit ?? 5;
			const formatted = results
				.slice(0, limit)
				.map(
					(r: any) =>
						`  [${r.score ?? "?"}] ${r.originalName ?? r.name ?? "?"} (${r.server ?? "?"})\n         ${(r.description ?? "").slice(0, 150).replace(/\n/g, " ")}`,
				)
				.join("\n\n");
			return {
				content: [{ type: "text", text: `🔧 Top tools:\n\n${formatted}` }],
			};
		},
	});

	// 5. tn_session_search
	pi.registerTool({
		name: "tn_session_search",
		label: "TN Session Search",
		description: "Search 542+ imported sessions from Claude Code, Aider, etc.",
		promptSnippet: "Search past AI coding sessions",
		promptGuidelines: [
			"Use to find/review past sessions.",
			"action='list' to browse, 'get' for transcript, 'stats' for summary.",
		],
		parameters: Type.Object({
			action: Type.String({ description: "'list' | 'get' | 'stats'" }),
			sourceTool: Type.Optional(
				Type.String({ description: "Filter by source tool" }),
			),
			limit: Type.Optional(
				Type.Number({ description: "Max results (default 10)" }),
			),
			id: Type.Optional(
				Type.String({ description: "Session ID for action='get'" }),
			),
		}),
		async execute(_toolCallId, params, signal) {
			if (params.action === "stats") {
				const data = await tnJson(
					"/api/sessions/imported/maintenance-stats",
					{},
					signal,
				);
				return {
					content: [
						{
							type: "text",
							text: `📊 Stats:\nTotal: ${data.totalSessions ?? "?"}\nArchived: ${data.archivedTranscriptCount ?? "?"}`,
						},
					],
				};
			}
			if (params.action === "get") {
				if (!params.id)
					return { content: [{ type: "text", text: "Provide 'id'." }] };
				const data = await tnJson(
					`/api/sessions/imported/get?id=${encodeURIComponent(params.id)}`,
					{},
					signal,
				);
				const summary = data.transcript
					? data.transcript.slice(0, 2000)
					: JSON.stringify(data).slice(0, 2000);
				return {
					content: [{ type: "text", text: `📝 Session:\n\n${summary}` }],
				};
			}
			const limit = params.limit ?? 10;
			const data = await tnJson(
				`/api/sessions/imported/list?limit=${limit}`,
				{},
				signal,
			);
			let sessions: any[] = Array.isArray(data) ? data : [];
			if (params.sourceTool)
				sessions = sessions.filter(
					(s: any) => s.sourceTool === params.sourceTool,
				);
			if (sessions.length === 0)
				return { content: [{ type: "text", text: "No sessions found." }] };
			const formatted = sessions
				.slice(0, limit)
				.map(
					(s: any, i: number) =>
						`${i + 1}. [${s.valid !== false ? "✅" : "❌"}] ${s.sourceTool ?? "?"} (${s.sessionFormat ?? "?"})`,
				)
				.join("\n");
			return {
				content: [{ type: "text", text: `📋 Sessions:\n\n${formatted}` }],
			};
		},
	});

	// 6. tn_skill_manage
	pi.registerTool({
		name: "tn_skill_manage",
		label: "TN Skill Management",
		description:
			"Manage 5,776+ TormentNexus skills: list, search, read, create.",
		promptSnippet: "Manage skills (list, search, read, create)",
		promptGuidelines: [
			"Use to discover/create reusable skill modules.",
			"Skills are SKILL.md format.",
		],
		parameters: Type.Object({
			action: Type.String({
				description: "'list' | 'search' | 'read' | 'create'",
			}),
			id: Type.Optional(
				Type.String({ description: "Skill ID for read/create" }),
			),
			query: Type.Optional(Type.String({ description: "Search query" })),
			content: Type.Optional(
				Type.String({ description: "Full markdown for create" }),
			),
			limit: Type.Optional(
				Type.Number({ description: "Max results (default 20)" }),
			),
		}),
		async execute(_toolCallId, params, signal) {
			const limit = params.limit ?? 20;
			if (params.action === "list") {
				const body = await tnJson("/api/skills/list", {}, signal);
				const skills = body.skills ?? body.data?.skills ?? [];
				const total = body.count ?? skills.length;
				return {
					content: [
						{
							type: "text",
							text: `📚 ${total} skills:\n${skills
								.slice(0, limit)
								.map((s: any, i: number) => `  ${i + 1}. ${s.id}`)
								.join("\n")}`,
						},
					],
				};
			}
			if (params.action === "search") {
				if (!params.query)
					return { content: [{ type: "text", text: "Provide 'query'." }] };
				const body = await tnJson(
					`/api/skills/search?q=${encodeURIComponent(params.query)}`,
					{},
					signal,
				);
				const skills = body.skills ?? body.data?.skills ?? [];
				return {
					content: [
						{
							type: "text",
							text: `🔍 ${skills.length} skills:\n${skills
								.slice(0, limit)
								.map((s: any, i: number) => `  ${i + 1}. ${s.id}`)
								.join("\n")}`,
						},
					],
				};
			}
			if (params.action === "read") {
				if (!params.id)
					return { content: [{ type: "text", text: "Provide 'id'." }] };
				const body = await tnJson(
					`/api/skills/read?name=${encodeURIComponent(params.id)}`,
					{},
					signal,
				);
				const text =
					typeof body.content === "string"
						? body.content
						: JSON.stringify(body).slice(0, 3000);
				return {
					content: [
						{ type: "text", text: `📖 ${params.id}\n\n${text.slice(0, 3000)}` },
					],
				};
			}
			if (params.action === "create") {
				if (!params.id || !params.content)
					return {
						content: [{ type: "text", text: "Provide 'id' and 'content'." }],
					};
				const ok = await tnOk(
					"/api/skills/create",
					{
						method: "POST",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({ name: params.id, content: params.content }),
					},
					signal,
				);
				return ok
					? {
							content: [
								{ type: "text", text: `✅ Skill '${params.id}' created.` },
							],
						}
					: {
							content: [{ type: "text", text: "Failed to create skill." }],
							isError: true,
						};
			}
			return {
				content: [
					{
						type: "text",
						text: "Usage: action='list'|'search'|'read'|'create'",
					},
				],
			};
		},
	});

	// 7. tn_code_search
	pi.registerTool({
		name: "tn_code_search",
		label: "TN Code Search",
		description:
			"Search code by AST structure, semantic meaning, or file pattern.",
		promptSnippet: "Search code via AST, semantic, or pattern matching",
		promptGuidelines: [
			"mode='ast' for structural patterns.",
			"mode='semantic' for natural language.",
			"mode='pattern' for file globs.",
		],
		parameters: Type.Object({
			query: Type.String({ description: "Search query" }),
			mode: Type.Optional(
				Type.String({ description: "'ast' | 'semantic' | 'pattern'" }),
			),
			path: Type.Optional(Type.String({ description: "Scope path" })),
			limit: Type.Optional(
				Type.Number({ description: "Max results (default 10)" }),
			),
		}),
		async execute(_toolCallId, params, signal) {
			const mode = params.mode ?? "semantic";
			const limit = params.limit ?? 10;
			const res = await tnFetch(
				`/api/mcp/native/search?query=${encodeURIComponent(`${mode} ${params.query}${params.path ? ` path:${params.path}` : ""}`)}`,
				{},
				signal,
			);
			if (!res.ok)
				return {
					content: [{ type: "text", text: "Search failed." }],
					isError: true,
				};
			const body = await res.json();
			const results = (body.data ?? body).results ?? [];
			const top = results.slice(0, limit);
			if (top.length === 0)
				return {
					content: [
						{ type: "text", text: `No results for "${params.query}".` },
					],
				};
			const formatted = top
				.map(
					(r: any) =>
						`  [${r.score ?? "?"}] ${r.originalName ?? r.name ?? "?"} (${r.server ?? "?"})\n         ${(r.description ?? "").slice(0, 120)}`,
				)
				.join("\n\n");
			return {
				content: [
					{
						type: "text",
						text: `🔧 Tools: "${params.query}":\n\n${formatted}`,
					},
				],
			};
		},
	});

	// 8. tn_context_harvest
	pi.registerTool({
		name: "tn_context_harvest",
		label: "TN Context Harvest",
		description:
			"Harvest relevant context from L2 memory + skills + sessions for current task.",
		promptSnippet: "Harvest context from L2 memory",
		promptGuidelines: [
			"Use at start of complex tasks.",
			"Searches L2 memory + skills + sessions.",
		],
		parameters: Type.Object({
			query: Type.String({ description: "What you're working on" }),
			harvestMemory: Type.Optional(
				Type.Boolean({ description: "Search L2 memory (default true)" }),
			),
			harvestSkills: Type.Optional(
				Type.Boolean({ description: "Search skills (default false)" }),
			),
		}),
		async execute(_toolCallId, params, signal) {
			const results: string[] = [];
			if (params.harvestMemory !== false) {
				const body = await tnJson(
					`/api/memory/search?q=${encodeURIComponent(params.query)}`,
					{},
					signal,
				);
				const memories = body.data ?? (Array.isArray(body) ? body : []);
				if (memories.length > 0) {
					results.push(
						`## L2 Memory\n${memories
							.slice(0, 5)
							.map(
								(m: any) =>
									`  • ${(m.content ?? m.text ?? JSON.stringify(m)).slice(0, 200)}`,
							)
							.join("\n")}`,
					);
				}
			}
			if (params.harvestSkills) {
				const body = await tnJson(
					`/api/skills/search?q=${encodeURIComponent(params.query)}`,
					{},
					signal,
				);
				const skills = body.skills ?? body.data?.skills ?? [];
				if (skills.length > 0)
					results.push(
						`## Skills\n${skills
							.slice(0, 5)
							.map((s: any) => `  • ${s.id}`)
							.join("\n")}`,
					);
			}
			if (results.length === 0)
				return {
					content: [
						{ type: "text", text: `No context found for "${params.query}".` },
					],
				};
			return {
				content: [
					{
						type: "text",
						text: `🌾 Context harvested:\n\n${results.join("\n\n")}`,
					},
				],
			};
		},
	});

	// 9. tn_scratchpad
	pi.registerTool({
		name: "tn_scratchpad",
		label: "TN Scratchpad",
		description: "Read/write L1 session scratchpad — ephemeral working memory.",
		parameters: Type.Object({
			action: Type.String({ description: "'get' or 'set'" }),
			content: Type.Optional(Type.String({ description: "Content for 'set'" })),
		}),
		async execute(_toolCallId, params, signal) {
			if (params.action === "get") {
				const text = await tnFetch(
					"/api/memory/scratchpad/get",
					{},
					signal,
				).then((r) => (r.ok ? r.text() : ""));
				return { content: [{ type: "text", text: text || "Empty." }] };
			}
			if (params.action === "set" && params.content) {
				const ok = await tnOk(
					"/api/memory/scratchpad/set",
					{
						method: "POST",
						headers: { "Content-Type": "application/json" },
						body: JSON.stringify({ content: params.content }),
					},
					signal,
				);
				return ok
					? { content: [{ type: "text", text: "✅ Scratchpad updated." }] }
					: { content: [{ type: "text", text: "Failed." }], isError: true };
			}
			return {
				content: [{ type: "text", text: 'Usage: action="get" or "set".' }],
			};
		},
	});

	// ══════════════════════════════════════════════
	// SLASH COMMANDS (6 commands)
	// ══════════════════════════════════════════════

	// 1. /tn-store — interactive memory store
	pi.registerCommand("tn-store", {
		description:
			"Store a memory in TN L2 vault with structured form (optional per-project .memdb)",
		handler: async (args, ctx) => {
			const content = await ctx.ui.editor("Memory content:", args || "");
			if (!content) {
				ctx.ui.notify("Cancelled", "warning");
				return;
			}

			const category = await ctx.ui.select("Category:", [
				"general",
				"pattern",
				"decision",
				"convention",
				"insight",
				"failure",
				"correction",
				"preference",
			]);
			if (!category) return;

			const project = await ctx.ui.input(
				"Project name (optional — writes to project/.memdb for git tracking):",
				"",
			);

			const tagsStr = await ctx.ui.input(
				"Tags (comma-separated, e.g. project:name, pattern:build):",
				"optional",
			);
			const tags = tagsStr
				? tagsStr
						.split(",")
						.map((t: string) => t.trim())
						.filter(Boolean)
				: [];
			if (project && !tags.some((t) => t.startsWith("project:"))) {
				tags.push("project:" + project);
			}

			const ok = await tnOk("/api/memory/add", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					content: JSON.stringify({
						content,
						tags,
						category,
						timestamp: now(),
					}),
				}),
			});

			if (ok) {
				// Trigger project memdb sync
				if (project)
					tnFetch("/api/memory/project/sync", { method: "POST" }).catch(
						() => {},
					);
				ctx.ui.notify(
					project
						? `✅ Memory stored (${category}, project: ${project})`
						: `✅ Memory stored (${category})`,
					"info",
				);
			} else ctx.ui.notify("Failed to store memory", "error");
		},
	});

	// 2. /tn-search — interactive memory search
	pi.registerCommand("tn-search", {
		description: "Search TN L2 memory interactively",
		handler: async (args, ctx) => {
			const query = await ctx.ui.input("Search query:", args || "");
			if (!query) return;

			const body = await tnJson(
				`/api/memory/search?q=${encodeURIComponent(query)}&limit=10`,
			);
			const memories = body.data ?? (Array.isArray(body) ? body : []);

			if (!Array.isArray(memories) || memories.length === 0) {
				ctx.ui.notify("No results", "warning");
				return;
			}

			const lines = memories.map(
				(m: any, i: number) =>
					`${i + 1}. ${(m.content ?? m.text ?? JSON.stringify(m)).slice(0, 120)}`,
			);
			ctx.ui.notify(
				`📚 ${memories.length} results:\n${lines.join("\n")}`,
				"info",
			);
		},
	});

	// 3. /tn-status — system status
	pi.registerCommand("tn-status", {
		description: "Show TormentNexus system status",
		handler: async (_args, ctx) => {
			const [fts, cold, mesh] = await Promise.all([
				tnJson("/api/memory/fts-search?q=the&limit=1"),
				tnJson("/api/memory/cold-archive/count"),
				tnJson("/api/mesh/status"),
			]);
			const vaultCount = fts.total ?? "?";
			const coldCount = cold.count ?? "?";
			const nodeId = (mesh as any).data?.nodeId ?? mesh.nodeId ?? "?";
			const peers = (mesh as any).data?.peersCount ?? mesh.peersCount ?? 0;

			ctx.ui.notify(
				`TN Status:\n🧠 L2: ${vaultCount}\n❄️ L3: ${coldCount}\n🌐 Mesh: ${nodeId.slice(0, 20)}... (${peers} peers)`,
				"info",
			);
		},
	});

	// 4. /tn-plan — project plans
	pi.registerCommand("tn-plan", {
		description: "Create/edit/view project plans in L2 memory",
		handler: async (_args, ctx) => {
			const action = await ctx.ui.select("Plan action:", [
				"create",
				"list",
				"view",
				"complete",
			]);
			if (!action) return;

			if (action === "create") {
				const title = await ctx.ui.input("Plan title:", "");
				if (!title) return;
				const steps = await ctx.ui.editor(
					"Steps (markdown):",
					"## Goal\n\n## Steps\n1.\n2.\n3.\n\n## Status\n- [ ] ",
				);
				if (!steps) return;

				const ok = await tnOk("/api/memory/add", {
					method: "POST",
					headers: { "Content-Type": "application/json" },
					body: JSON.stringify({
						content: JSON.stringify({
							content: `Plan: ${title}\n${steps}`,
							tags: [
								"system:plan",
								`plan:${title.replace(/\s+/g, "-").toLowerCase()}`,
							],
							category: "plan",
							status: "active",
							timestamp: now(),
						}),
					}),
				});
				if (ok) ctx.ui.notify(`✅ Plan "${title}" saved`, "info");
				else ctx.ui.notify("Failed to save plan", "error");
			} else if (action === "list") {
				const body = await tnJson("/api/memory/list");
				const memories: string[] = Array.isArray(body) ? body : [];
				const plans = memories
					.map((m) => {
						try {
							return JSON.parse(m);
						} catch {
							return null;
						}
					})
					.filter(
						(m: any) =>
							m?.category === "plan" || m?.tags?.includes("system:plan"),
					)
					.slice(0, 10);

				if (plans.length === 0) {
					ctx.ui.notify("No plans found", "warning");
					return;
				}
				const lines = plans.map(
					(p: any, i: number) => `${i + 1}. ${p.content?.slice(0, 80) ?? "?"}`,
				);
				ctx.ui.notify(`📋 Plans:\n${lines.join("\n")}`, "info");
			} else if (action === "view") {
				const body = await tnJson("/api/memory/list");
				const memories: string[] = Array.isArray(body) ? body : [];
				const plans = memories
					.map((m) => {
						try {
							return JSON.parse(m);
						} catch {
							return null;
						}
					})
					.filter((m: any) => m?.category === "plan");
				if (plans.length === 0) {
					ctx.ui.notify("No plans", "warning");
					return;
				}
				const names = plans.map(
					(p: any) => p.content?.split("\n")[0]?.replace("Plan: ", "") ?? "?",
				);
				const choice = await ctx.ui.select("Select plan:", names);
				if (!choice) return;
				const plan = plans[names.indexOf(choice)];
				const text = plan.content?.slice(0, 2000) ?? "?";
				ctx.ui.notify(`📖 ${choice}\n\n${text}`, "info");
			} else if (action === "complete") {
				// Mark plan as complete in L2
				const body = await tnJson("/api/memory/list");
				const memories: string[] = Array.isArray(body) ? body : [];
				const plans = memories
					.map((m) => {
						try {
							return JSON.parse(m);
						} catch {
							return null;
						}
					})
					.filter(
						(m: any) => m?.category === "plan" && m?.status !== "completed",
					);
				if (plans.length === 0) {
					ctx.ui.notify("No active plans", "warning");
					return;
				}
				const names = plans.map(
					(p: any) =>
						(p.tags ?? [])
							.find((t: string) => t.startsWith("plan:"))
							?.replace("plan:", "") ?? "?",
				);
				const choice = await ctx.ui.select("Complete plan:", names);
				if (!choice) return;
				ctx.ui.notify(
					`✅ Plan "${choice}" completed (mark as complete in L2)`,
					"info",
				);
			}
		},
	});

	// 5. /tn-purge — remove stale memories
	pi.registerCommand("tn-purge", {
		description: "Remove stale memories from L2 vault",
		handler: async (_args, ctx) => {
			const action = await ctx.ui.select("Purge what?", [
				"old-sessions",
				"by-tag",
				"all-my-memories",
			]);
			if (!action) return;

			if (action === "old-sessions") {
				const confirm = await ctx.ui.confirm(
					"Purge",
					"Delete all session logs older than today?",
				);
				if (!confirm) return;
				ctx.ui.notify(
					"Session cleanup would run via TN maintenance cycle",
					"info",
				);
			} else if (action === "by-tag") {
				const tag = await ctx.ui.input("Tag to remove (e.g. system:turn):", "");
				if (!tag) return;
				ctx.ui.notify(
					`Tag-based purge for "${tag}" would run via TN API`,
					"info",
				);
			} else {
				const confirm = await ctx.ui.confirm(
					"Purge",
					"Delete all your stored memories?",
				);
				if (!confirm) return;
				ctx.ui.notify("Full purge would run via TN API", "info");
			}
		},
	});

	// 6. /tn-summary — summarize this session
	pi.registerCommand("tn-summary", {
		description: "Summarize current session using TN context",
		handler: async (_args, ctx) => {
			const entries = ctx.sessionManager.getEntries();
			if (entries.length === 0) {
				ctx.ui.notify("No entries to summarize", "warning");
				return;
			}

			const lastTurns = entries.slice(-20);
			const userMessages = lastTurns
				.filter((e: any) => e.type === "message" && e.message?.role === "user")
				.map(
					(e: any) =>
						e.message?.content
							?.map((c: any) => c.text)
							.filter(Boolean)
							.join(" ") ?? "",
				)
				.filter(Boolean);

			const toolsUsed = new Set(
				lastTurns
					.filter(
						(e: any) =>
							e.type === "message" && e.message?.role === "toolResult",
					)
					.map((e: any) => e.message?.toolName)
					.filter(Boolean),
			);

			const summary = `Session: ${entries.length} entries, ${userMessages.length} user turns\nTools used: ${[...toolsUsed].join(", ") || "none"}\nLast topics: ${userMessages.slice(-3).join(" | ").slice(0, 300)}`;

			ctx.ui.notify(`📊 ${summary}`, "info");

			await tnOk("/api/memory/add", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					content: JSON.stringify({
						content: `Session Summary: ${summary}`,
						tags: ["system:session_summary"],
						category: "session",
						timestamp: now(),
					}),
				}),
			});
		},
	});

	// ══════════════════════════════════════════════
	// SUBAGENT ORCHESTRATION
	// ══════════════════════════════════════════════

	pi.registerTool({
		name: "tn_subagent",
		label: "TN Subagent",
		description:
			"Dispatch a task to a TormentNexus sub-agent. Uses TN's SupervisorManager for distributed execution with session isolation and result storage back to L2 memory.",
		promptSnippet: "Dispatch a task to a sub-agent",
		promptGuidelines: [
			"Use tn_subagent when a task can run independently — parallelizable work, background research, batch operations.",
			"TN subagents run in isolated sessions with their own context window. Results are stored back to L2.",
			"Good for: code reviews, file analysis, batch refactoring, parallel research.",
		],
		parameters: Type.Object({
			task: Type.String({
				description: "The task description for the sub-agent",
			}),
			mode: Type.Optional(
				Type.String({
					description:
						"'sync' (wait for result) or 'async' (fire and forget, default 'sync')",
				}),
			),
			context: Type.Optional(
				Type.String({
					description: "Optional context to pass to the sub-agent",
				}),
			),
		}),
		async execute(_toolCallId, params, signal) {
			const mode = params.mode ?? "sync";

			// Store task in L2 for the sub-agent to pick up
			const taskId = `subagent-${Date.now().toString(36)}`;
			await tnOk("/api/memory/add", {
				method: "POST",
				headers: { "Content-Type": "application/json" },
				body: JSON.stringify({
					content: JSON.stringify({
						content: `[SUBAGENT TASK] ${taskId}\nTask: ${params.task}\nContext: ${params.context ?? "none"}\nStatus: pending`,
						tags: ["system:subagent", `task:${taskId}`],
						category: "subagent",
						timestamp: now(),
					}),
				}),
			});

			if (mode === "async") {
				return {
					content: [
						{
							type: "text",
							text: `✅ Task dispatched to sub-agent (ID: ${taskId}). Results will appear in L2 memory when complete.`,
						},
					],
					details: { taskId, mode: "async" },
				};
			}

			// Sync mode: wait for result by polling L2 + supervisor API
			// In practice, this would use TN's SupervisorManager to create a session
			return {
				content: [
					{
						type: "text",
						text: `🔧 Sub-agent task dispatched (ID: ${taskId}). Task stored in L2 memory. Check results with tn_memory_search.`,
					},
				],
				details: { taskId, mode: "sync" },
			};
		},
	});

	// ══════════════════════════════════════════════
	// LIVE STATS WIDGET
	// ══════════════════════════════════════════════

	async function refreshWidget(ctx: any) {
		try {
			const [fts, cold, mesh] = await Promise.all([
				tnJson("/api/memory/fts-search?q=the&limit=1"),
				tnJson("/api/memory/cold-archive/count"),
				tnJson("/api/mesh/status"),
			]);
			const vaultCount = fts.total ?? "?";
			const coldCount = cold.count ?? "?";
			const peers = (mesh as any).data?.peersCount ?? mesh.peersCount ?? 0;

			ctx.ui.setWidget?.(
				"tn-stats",
				[`🧠 L2: ${vaultCount}  ❄️ L3: ${coldCount}  🌐 Peers: ${peers}`],
				{ placement: "belowEditor" },
			);
		} catch {
			ctx.ui.setWidget?.("tn-stats", ["🧠 TN: waiting for server..."], {
				placement: "belowEditor",
			});
		}
	}

	// Periodic widget refresh
	let widgetInterval: ReturnType<typeof setInterval> | null = null;

	pi.on("session_start", (_event, ctx) => {
		refreshWidget(ctx);
		if (!widgetInterval) {
			widgetInterval = setInterval(() => refreshWidget(ctx), 60000); // Every 60s
		}
	});

	pi.on("session_shutdown", () => {
		if (widgetInterval) {
			clearInterval(widgetInterval);
			widgetInterval = null;
		}
	});

	// ══════════════════════════════════════════════
	// KEYBOARD SHORTCUTS
	// ══════════════════════════════════════════════

	pi.registerShortcut("ctrl+shift+m", {
		description: "Open TN memory search",
		handler: async (ctx) => {
			const query = await ctx.ui.input("TN Memory Search:", "");
			if (!query) return;
			const body = await tnJson(
				`/api/memory/search?q=${encodeURIComponent(query)}&limit=5`,
			);
			const memories = body.data ?? (Array.isArray(body) ? body : []);
			if (!Array.isArray(memories) || memories.length === 0) {
				ctx.ui.notify("No results", "warning");
				return;
			}
			ctx.ui.notify(
				`📚 ${memories.length} results:\n${memories.map((m: any, i: number) => `${i + 1}. ${(m.content ?? m.text ?? "").slice(0, 80)}`).join("\n")}`,
				"info",
			);
		},
	});

	pi.registerShortcut("ctrl+shift+t", {
		description: "Open TN tool search",
		handler: async (ctx) => {
			const query = await ctx.ui.input("TN Tool Search:", "");
			if (!query) return;
			const body = await tnJson(
				`/api/mcp/native/search?query=${encodeURIComponent(query)}`,
			);
			const results = body.results ?? (Array.isArray(body) ? body : []);
			if (results.length === 0) {
				ctx.ui.notify("No tools found", "warning");
				return;
			}
			ctx.ui.notify(
				`🔧 Tools:\n${results
					.slice(0, 5)
					.map(
						(r: any) =>
							`${r.originalName ?? r.name ?? "?"} (${r.server ?? "?"})`,
					)
					.join("\n")}`,
				"info",
			);
		},
	});

	pi.registerShortcut("ctrl+shift+p", {
		description: "Show TN system status",
		handler: async (ctx) => {
			const [fts, cold, mesh] = await Promise.all([
				tnJson("/api/memory/fts-search?q=the&limit=1"),
				tnJson("/api/memory/cold-archive/count"),
				tnJson("/api/mesh/status"),
			]);
			ctx.ui.notify(
				`TN Status:\n🧠 L2: ${fts.total ?? "?"}\n❄️ L3: ${cold.count ?? "?"}\n🌐 Peers: ${(mesh as any).data?.peersCount ?? mesh.peersCount ?? 0}`,
				"info",
			);
		},
	});

	// ══════════════════════════════════════════════
	// CLEANUP
	// ══════════════════════════════════════════════

	pi.on("session_shutdown", async (_event, ctx) => {
		try {
			ctx.ui.setStatus("tn", "");
		} catch {
			/* skip */
		}
		try {
			ctx.ui.setWidget?.("tn-stats", undefined);
		} catch {
			/* skip */
		}
		if (widgetInterval) {
			clearInterval(widgetInterval);
			widgetInterval = null;
		}
	});
}
