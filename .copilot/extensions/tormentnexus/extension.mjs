import { approveAll } from "@github/copilot-sdk";
import { joinSession } from "@github/copilot-sdk/extension";

const TN_API = "http://127.0.0.1:7778";
const TN_WORKSPACE = "C:\\Users\\hyper\\workspace\\tormentnexus";

async function tnPost(path, body) {
  try {
    await fetch(`${TN_API}${path}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
      signal: AbortSignal.timeout(3000),
    });
  } catch {}
}

async function tnSearch(query) {
  try {
    const q = encodeURIComponent(query.slice(0, 200));
    const res = await fetch(`${TN_API}/api/memory/search?q=${q}`, {
      signal: AbortSignal.timeout(3000),
    });
    const json = await res.json();
    if (json.data && json.data.length > 0) return json.data.slice(0, 3);
  } catch {}
  return [];
}

const session = await joinSession({
  onPermissionRequest: approveAll,

  // Inject TN system guidance into every session
  systemMessage: {
    role: "system",
    content: `## TormentNexus Integration

You have access to TormentNexus — a local AI control plane with persistent L2 vector memory, semantic tool discovery, and imported sessions.

### Memory Tools
- \`tn_memory_store\` — Store key decisions, patterns, and context with tags
- \`tn_memory_search\` — Find relevant context from past sessions

### Best Practices
1. Check memory before significant work
2. Store patterns and decisions after key moments
3. Use TN MCP tools for advanced operations (repomap, tool discovery)`,
  },

  // Register custom TN tools
  tools: [
    {
      name: "tn_memory_store",
      description: "Store a decision, pattern, or fact in TormentNexus L2 memory with tags",
      parameters: {
        type: "object",
        required: ["content"],
        properties: {
          content: { type: "string", description: "What to remember" },
          tags: { type: "array", items: { type: "string" }, description: "Optional tags e.g. project:foo, pattern, decision" },
        },
      },
      handler: async ({ content, tags }) => {
        await tnPost("/api/memory/add", {
          content: JSON.stringify({ content, tags: tags || ["agent:copilot"], category: "memory" }),
        });
        return { data: "Stored in TormentNexus L2 memory" };
      },
    },
    {
      name: "tn_memory_search",
      description: "Search TormentNexus L2 memory for relevant context",
      parameters: {
        type: "object",
        required: ["query"],
        properties: {
          query: { type: "string", description: "What to search for" },
        },
      },
      handler: async ({ query }) => {
        const memories = await tnSearch(query);
        if (memories.length === 0) return { data: "No relevant memories found" };
        return {
          data: memories.map((m) => `• ${m.text || m.content || ""}`).join("\n"),
        };
      },
    },
  ],

  // Point to TN skill directory
  skillDirectories: [`${TN_WORKSPACE}\\.copilot\\skills`],

  // All 5 lifecycle hooks
  hooks: {
    onSessionStart: async ({ sessionId, reason }) => {
      await tnPost("/api/memory/add", {
        content: JSON.stringify({
          content: `Copilot session ${reason}: ${sessionId}`,
          tags: ["system:session", "agent:copilot"],
          category: "session",
        }),
      });
    },

    onSessionEnd: async ({ sessionId }) => {
      await tnPost("/api/memory/add", {
        content: JSON.stringify({
          content: `Copilot session ended: ${sessionId}`,
          tags: ["system:session_end", "agent:copilot"],
          category: "session",
        }),
      });
    },

    onUserPromptSubmitted: async ({ prompt }) => {
      const match = prompt.match(/@memory:(\S+)/);
      if (match) {
        const memories = await tnSearch(match[1]);
        if (memories.length > 0) {
          const value = memories[0].text || memories[0].content || "<not found>";
          return { prompt: prompt.replace(`@memory:${match[1]}`, value) };
        }
      }
      const memories = await tnSearch(prompt);
      if (memories.length > 0) {
        const context = memories.map((m) => `  • ${(m.text || m.content || "").slice(0, 200)}`).join("\n");
        return { context: `Relevant TormentNexus context:\n${context}` };
      }
    },

    onPreToolUse: async ({ toolName }) => {
      await tnPost("/api/memory/add", {
        content: JSON.stringify({
          content: `Copilot tool: ${toolName}`,
          tags: ["system:tool_call", `tool:${toolName}`, "agent:copilot"],
          category: "tool",
        }),
      });
      return { allow: true };
    },

    onPostToolUse: async ({ toolName, result }) => {
      if (toolName === "read" || toolName === "ls") return;
      const resultStr = JSON.stringify(result);
      if (resultStr.length > 2000) return;
      await tnPost("/api/memory/add", {
        content: JSON.stringify({
          content: `Copilot ${toolName} completed`,
          tags: ["system:tool_result", `tool:${toolName}`, "agent:copilot"],
          category: "tool_result",
        }),
      });
    },
  },
});
