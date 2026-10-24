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
  description: "List recent HyperNexus memory facts (newest first) with titles, tags, and timestamps.\nArguments:\n- limit: max entries (default 20).\n- namespace: optional namespace filter.",
  args: { limit: z.number().int().min(1).max(100).optional(), namespace: z.string().optional() },
  async execute(args) {
    const params = new URLSearchParams({ limit: String(args.limit || 20) });
    if (args.namespace) params.set("namespace", args.namespace);
    const data = await hnFetch("/api/memory/facts/list?" + params.toString());
    if (data.error) {
      const alt = await hnFetch("/api/memory/search?query=*&limit=" + (args.limit || 20));
      if (alt.error) return data.error;
      return fmt(alt);
    }
    const items = data.results || data.facts || data.items || (Array.isArray(data) ? data : []);
    if (!items.length) return "No memories stored yet.";
    return items.map((m, i) => (i + 1) + ". [" + (m.tags || []).join(",") + "] " + (m.title || "(untitled)") + " (" + (m.created_at || m.createdAt || "") + ")").join("\n");
  },
};
