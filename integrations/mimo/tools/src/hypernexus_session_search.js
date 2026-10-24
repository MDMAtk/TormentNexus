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
  description: "Search imported sessions (from Claude Code, Aider, Gemini, etc.) in the HyperNexus session store.\nArguments:\n- query: search terms.\n- limit: max results (default 10).",
  args: { query: z.string().min(1), limit: z.number().int().min(1).max(50).optional() },
  async execute(args) {
    const data = await hnFetch("/api/sessions/search?query=" + encodeURIComponent(args.query) + "&limit=" + (args.limit || 10));
    if (data.error) return data.error;
    return fmt(data);
  },
};
