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
  description: "Harvest relevant L2 context from HyperNexus memory for the current task and return it for injection.\nCall before complex multi-step tasks to pull in prior related work.\nArguments:\n- topic: the task topic / question to harvest context for.\n- depth: harvest depth 'quick' | 'standard' | 'deep' (default standard).",
  args: { topic: z.string().min(1), depth: z.enum(["quick", "standard", "deep"]).optional() },
  async execute(args) {
    const data = await hnFetch("/api/memory/context/harvest", {
      method: "POST",
      body: JSON.stringify({ topic: args.topic, depth: args.depth || "standard" }),
    });
    if (data.error) {
      const alt = await hnFetch("/api/memory/search?query=" + encodeURIComponent(args.topic) + "&limit=15");
      if (alt.error) return data.error;
      return fmt(alt);
    }
    return fmt(data);
  },
};
