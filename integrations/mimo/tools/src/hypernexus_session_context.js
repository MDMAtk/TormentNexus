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
  description: "Get or set the current session context (goal/objective) used by HyperNexus for tool-query bootstrap and context harvesting.\nArguments:\n- action: 'get' to read current context, 'set' to update it.\n- goal: when action=set, the session goal text.\n- objective: when action=set, the session objective text.",
  args: { action: z.enum(["get", "set"]), goal: z.string().optional(), objective: z.string().optional() },
  async execute(args) {
    if (args.action === "get") {
      const data = await hnFetch("/api/session/context");
      if (data.error) return data.error;
      return fmt(data);
    }
    const data = await hnFetch("/api/session/context", {
      method: "POST",
      body: JSON.stringify({ goal: args.goal || "", objective: args.objective || "" }),
    });
    if (data.error) return data.error;
    return "Session context updated.";
  },
};
