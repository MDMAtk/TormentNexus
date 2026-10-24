import { z } from "zod";

const API = (typeof process !== "undefined" && (process.env.HYPERNEXUS_URL || process.env.HN_URL)) || "http://127.0.0.1:7778";
async function hnFetch(path, opts) {
  const url = API.replace(/\/$/, "") + path;
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


export default {
  description: "Route a tool call through the HyperNexus kernel to a downstream MCP server.\nUse hypernexus_mcp_search first to discover the right serverName/toolName.\nArguments:\n- serverName: MCP server name (from search results).\n- toolName: tool name on that server.\n- arguments: JSON object of tool arguments.",
  args: { serverName: z.string().min(1), toolName: z.string().min(1), arguments: z.record(z.string(), z.any()).optional() },
  async execute(args) {
    const data = await hnFetch("/api/mcp/tools/call", {
      method: "POST",
      body: JSON.stringify({ serverName: args.serverName, toolName: args.toolName, arguments: args.arguments || {} }),
    });
    if (data.error) return data.error;
    return fmt(data.result !== undefined ? data.result : data);
  },
};
