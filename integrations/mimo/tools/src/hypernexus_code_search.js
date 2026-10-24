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
  description: "Search code through the HyperNexus kernel (pattern, AST-grep, or semantic).\nArguments:\n- query: search pattern or natural language.\n- mode: 'pattern' (regex/ripgrep), 'ast' (ast-grep), or 'semantic'.\n- path: optional path to scope the search.",
  args: { query: z.string().min(1), mode: z.enum(["pattern", "ast", "semantic"]).optional(), path: z.string().optional() },
  async execute(args) {
    const data = await hnFetch("/api/code/search", {
      method: "POST",
      body: JSON.stringify({ query: args.query, mode: args.mode || "pattern", path: args.path || "" }),
    });
    if (data.error) return data.error;
    return fmt(data);
  },
};
