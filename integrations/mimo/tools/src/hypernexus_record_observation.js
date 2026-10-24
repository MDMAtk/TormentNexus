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
  description: "Record a structured observation into the HyperNexus knowledge graph (entity + observation).\nUse for durable domain knowledge: entities you discovered, relationships, constraints, gotchas.\nArguments:\n- entity: the entity name this observation is about.\n- observation: the observation text.\n- relations: optional array of {to, type} relation hints.",
  args: { entity: z.string().min(1), observation: z.string().min(1), relations: z.array(z.object({ to: z.string(), type: z.string() })).optional() },
  async execute(args) {
    const data = await hnFetch("/api/memory/observations/add", {
      method: "POST",
      body: JSON.stringify({ entity: args.entity, observation: args.observation, relations: args.relations || [] }),
    });
    if (data.error) return data.error;
    return "Observation recorded for entity: " + args.entity;
  },
};
