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
  description: "Store an important fact, decision, or finding into HyperNexus persistent memory.\nCall this after completing a significant step, making an architectural decision, or discovering a non-obvious pattern.\nArguments:\n- title: short descriptive title for the fact.\n- content: the fact / decision / finding body.\n- tags: array of tags for retrieval (e.g. [\"architecture\", \"bugfix\"]).\n- namespace: memory namespace (project | global | user), default project.",
  args: { title: z.string().min(1), content: z.string().min(1), tags: z.array(z.string()).optional(), namespace: z.string().optional() },
  async execute(args) {
    const data = await hnFetch("/api/memory/facts/add", {
      method: "POST",
      body: JSON.stringify({ title: args.title, content: args.content, tags: args.tags || [], namespace: args.namespace || "project" }),
    });
    if (data.error) return data.error;
    return "Fact stored: " + args.title + (data.id ? " (id=" + data.id + ")" : "");
  },
};
