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
  description: "Semantic vector search over HyperNexus L2 memory — finds conceptually related past context even without exact keyword matches.\nUse for nuanced recall when keyword search is too literal.\nArguments:\n- query: natural language query.\n- limit: max results (default 8).",
  args: { query: z.string().min(1), limit: z.number().int().min(1).max(50).optional() },
  async execute(args) {
    const data = await hnFetch("/api/memory/vector/search", {
      method: "POST",
      body: JSON.stringify({ query: args.query, limit: args.limit || 8 }),
    });
    if (data.error) {
      const alt = await hnFetch("/api/memory/search?query=" + encodeURIComponent(args.query) + "&limit=" + (args.limit || 8) + "&mode=vector");
      if (alt.error) return data.error;
      return fmt(alt);
    }
    return fmt(data);
  },
};
