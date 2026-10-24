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
  description: "Semantically search across all 60+ managed MCP servers and 26,000+ tools to find the right tool for a job.\nPrefer this over manually guessing tool names before calling hypernexus_mcp_call.\nArguments:\n- query: natural language description of what you want to do.\n- limit: max tools to return (default 8).",
  args: { query: z.string().min(1), limit: z.number().int().min(1).max(50).optional() },
  async execute(args) {
    const data = await hnFetch("/api/mcp/tools/search?query=" + encodeURIComponent(args.query) + "&limit=" + (args.limit || 8));
    if (data.error) return data.error;
    const items = data.tools || data.results || (Array.isArray(data) ? data : []);
    if (!items.length) return "No MCP tools matched: " + args.query;
    return "Found " + items.length + " tools:\n" + items.map((t, i) =>
      (i + 1) + ". " + (t.server || t.serverName || "?") + "/" + (t.name || t.toolName || "?") + " — " + String(t.description || "").slice(0, 160)
    ).join("\n");
  },
};
