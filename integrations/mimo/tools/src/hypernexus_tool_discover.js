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
  description: "Discover HyperNexus capabilities and MCP tools by natural-language intent (broader than hypernexus_mcp_search — includes kernel APIs and skills).\nArguments:\n- intent: what you want to accomplish.",
  args: { intent: z.string().min(1) },
  async execute(args) {
    const data = await hnFetch("/api/mcp/tools/search?query=" + encodeURIComponent(args.intent) + "&limit=12");
    if (data.error) return data.error;
    return fmt(data);
  },
};
