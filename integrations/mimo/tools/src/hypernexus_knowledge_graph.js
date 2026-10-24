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
  description: "Add entities/relations or query the HyperNexus knowledge graph.\nArguments:\n- action: 'add_entities', 'add_relations', 'add_observations', or 'search'.\n- entities: array of {name, type, observations} for add_entities.\n- relations: array of {from, to, type} for add_relations.\n- query: search string for action=search.",
  args: { action: z.enum(["add_entities", "add_relations", "add_observations", "search"]), entities: z.array(z.object({ name: z.string(), type: z.string().optional(), observations: z.array(z.string()).optional() })).optional(), relations: z.array(z.object({ from: z.string(), to: z.string(), type: z.string() })).optional(), query: z.string().optional() },
  async execute(args) {
    const map = {
      add_entities: ["/api/memory/entities/add", { entities: args.entities || [] }],
      add_relations: ["/api/memory/relations/add", { relations: args.relations || [] }],
      add_observations: ["/api/memory/observations/add", { entities: args.entities || [] }],
      search: ["/api/memory/search?query=" + encodeURIComponent(args.query || ""), null],
    };
    const entry = map[args.action];
    if (!entry) return "Unknown action: " + args.action;
    const [path, body] = entry;
    const data = body
      ? await hnFetch(path, { method: "POST", body: JSON.stringify(body) })
      : await hnFetch(path);
    if (data.error) return data.error;
    return fmt(data);
  },
};
