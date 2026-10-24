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
  description: "Search HyperNexus persistent memory (L2/L3) by keyword, tag, or category.\nUse this BEFORE starting any multi-step task to recall past context, decisions, and facts.\nReturns matching memory facts with title, content, tags, and timestamps.\nArguments:\n- query: search terms (keywords, tags, or natural language).\n- limit: max results to return (default 10).\n- namespace: optional memory namespace filter (project | global | user).",
  args: { query: z.string().min(1), limit: z.number().int().min(1).max(100).optional(), namespace: z.string().optional() },
  async execute(args) {
    const params = new URLSearchParams({ query: args.query, limit: String(args.limit || 10) });
    if (args.namespace) params.set("namespace", args.namespace);
    const data = await hnFetch("/api/memory/search?" + params.toString());
    if (data.error) return data.error;
    const items = data.results || data.facts || data.items || (Array.isArray(data) ? data : []);
    if (!items.length) return "No memories matched: " + args.query;
    return "Found " + items.length + " memories:\n" + items.map((m, i) =>
      (i + 1) + ". [" + (m.tags || []).join(",") + "] " + (m.title || "(untitled)") + " — " + String(m.content || m.text || "").slice(0, 200)
    ).join("\n");
  },
};
