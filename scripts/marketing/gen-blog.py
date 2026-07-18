blogdir = "/var/www/tormentnexus.site/blog/tormentnexus"

posts = [
    {
        "slug": "building-ai-agents-that-survive-restarts-persistent-memory-done-right",
        "title": "Building AI Agents That Survive Restarts: Persistent Memory Done Right",
        "date": "July 13, 2026",
        "body": '<p>Most AI agents lose everything when they crash. Session context, learned patterns, user preferences — all gone. TormentNexus solves this with a multi-tier memory architecture that persists across restarts.</p>\n<h2>The Problem: Ephemeral Intelligence</h2>\n<p>Every time an AI agent restarts, it starts from zero. No memory of past conversations, no knowledge of previous decisions. This is like hiring a new employee every morning who has never seen your codebase.</p>\n<h2>The Solution: Multi-Tier Persistent Memory</h2>\n<p>TormentNexus implements four memory tiers:</p>\n<ul style="margin:1rem 0 1.5rem 1.5rem;color:#b8b8c8">\n<li><strong>L1 Scratchpad</strong> — Working memory for the current session</li>\n<li><strong>L2 Vault</strong> — Persistent vector-indexed memories with semantic search</li>\n<li><strong>L3 Cold Archive</strong> — Compressed long-term storage for older memories</li>\n<li><strong>L4 Limbo</strong> — Quarantined memories pending review or deletion</li>\n</ul>\n<h2>How It Works</h2>\n<p>When an agent encounters something worth remembering, it stores a memory entry with content, tags, category, and importance score. The L2 vault uses SQLite with FTS5 full-text search and vector embeddings for semantic retrieval.</p>\n<h2>Spaced Repetition for Consolidation</h2>\n<p>Inspired by how human brains consolidate memories during sleep, TormentNexus runs periodic memory maintenance. Memories that lose relevance gradually cool and move to cold archive. Frequently accessed memories stay hot and rank higher in search results.</p>\n<h2>The Result</h2>\n<p>An agent that remembers. After a restart, the Go Coder agent picks up where it left off — same workspace, same memories, same learned patterns. This is institutional memory that survives across sessions, crashes, and migrations.</p>',
    },
    {
        "slug": "the-llm-waterfall-pattern-escaping-api-rate-limits-with-provider-failover",
        "title": "The LLM Waterfall Pattern: Escaping API Rate Limits with Provider Failover",
        "date": "July 13, 2026",
        "body": '<p>Your AI agent hits a rate limit at 2am. Instead of waiting, it automatically fails over to a different provider — no human intervention, no downtime. This is the LLM Waterfall pattern.</p>\n<h2>Why Single-Provider Architectures Fail</h2>\n<p>Every LLM provider has rate limits, outages, and maintenance windows. If your agent depends on a single provider, any outage is a full stop.</p>\n<h2>The Waterfall Architecture</h2>\n<p>TormentNexus cascades through multiple providers in priority order:</p>\n<ul style="margin:1rem 0 1.5rem 1.5rem;color:#b8b8c8">\n<li><strong>Primary:</strong> NVIDIA NIM (fastest, free tier)</li>\n<li><strong>Secondary:</strong> OpenRouter (broadest model selection)</li>\n<li><strong>Tertiary:</strong> Local Ollama (zero latency, zero cost)</li>\n</ul>\n<p>The waterfall is transparent to the calling code. The agent does not know or care which provider answered.</p>\n<h2>Provider-Agnostic Tool Routing</h2>\n<p>The waterfall extends beyond LLM calls. MCP tool routing follows the same pattern: try the local registry first, fall back to remote, then generate inline.</p>\n<h2>Zero-Downtime Deployments</h2>\n<p>Combined with the Go sidecar architecture, the waterfall pattern means TormentNexus can update, restart, and redeploy without dropping a single agent request.</p>',
    },
    {
        "slug": "dual-tier-memory-architecture-for-ai-agents-why-your-agent-needs-an-l1-scratchpad-and-an-l2-vault",
        "title": "Dual-Tier Memory Architecture for AI Agents: L1 Scratchpad + L2 Vault",
        "date": "July 13, 2026",
        "body": "<p>Human brains have working memory and long-term memory. AI agents need the same split. TormentNexus implements this as L1 (fast, volatile scratchpad) and L2 (persistent, searchable vault).</p>\n<h2>L1: The Scratchpad</h2>\n<p>L1 is the agent working memory — fast, in-process, and ephemeral. It holds the current conversation context, intermediate reasoning steps, and scratch calculations. When the session ends, L1 clears.</p>\n<h2>L2: The Vault</h2>\n<p>L2 is where permanence lives. Every memory entry has content, tags, category, importance score, and a heat score that decays over time. The vault uses SQLite with FTS5 for full-text search and vector embeddings for semantic similarity.</p>\n<h2>The Bridge: Harvesting L1 into L2</h2>\n<p>The Memory Harvester agent periodically scans L1 contents and extracts high-value memories for L2 storage. It identifies patterns worth keeping: architecture decisions, bug fixes, user preferences.</p>\n<h2>Why Two Tiers Matter</h2>\n<p>Without L1, agents are slow — every piece of context requires a database query. Without L2, agents are forgetful — everything vanishes at session end. Together, they give agents the speed of working memory with the durability of long-term storage.</p>",
    },
    {
        "slug": "progressive-mcp-tool-routing-stop-drowning-your-agents-in-50k-tokens",
        "title": "Progressive MCP Tool Routing: Stop Drowning Your Agents in 50K Tokens",
        "date": "July 13, 2026",
        "body": '<p>The MCP ecosystem has over 20,000 registered tools. Dumping all of them into the LLM context window would consume 50,000+ tokens before the agent even starts thinking. Progressive routing solves this.</p>\n<h2>The Token Budget Problem</h2>\n<p>Every tool definition costs tokens. A typical tool schema is 200-500 tokens. With 20,000 tools, that is 4-10 million tokens of pure tool definitions — far beyond any model context window.</p>\n<h2>Progressive Routing: Three Stages</h2>\n<ul style="margin:1rem 0 1.5rem 1.5rem;color:#b8b8c8">\n<li><strong>Stage 1 — Semantic Match:</strong> Vector embeddings find the top 20 tools most relevant to the current task</li>\n<li><strong>Stage 2 — Context Filter:</strong> Narrow based on project context, recent tool usage, and agent capabilities</li>\n<li><strong>Stage 3 — Lazy Load:</strong> Only inject the final 5-10 tool definitions into the context window</li>\n</ul>\n<h2>The Result</h2>\n<p>Context window usage drops from 50K+ tokens to under 2K. Agents respond faster, make fewer mistakes, and stay within budget. This is the difference between giving someone a library card and dumping 20,000 books on their desk.</p>',
    },
    {
        "slug": "zero-trust-ai-architecture-authenticating-every-tool-call-memory-access-and-model-request",
        "title": "Zero-Trust AI Architecture: Authenticating Every Tool Call, Memory Access, and Model Request",
        "date": "July 13, 2026",
        "body": '<p>When you give an AI agent access to your filesystem, database, and API keys, you are trusting it with everything. Zero-trust architecture means every single operation is authenticated, authorized, and audited.</p>\n<h2>The Threat Model</h2>\n<p>AI agents can execute arbitrary code, access sensitive data, and make network requests. A compromised or misbehaving agent could exfiltrate data, delete files, or escalate privileges.</p>\n<h2>Three Pillars of Zero-Trust AI</h2>\n<ul style="margin:1rem 0 1.5rem 1.5rem;color:#b8b8c8">\n<li><strong>Authentication:</strong> Every agent, tool call, and API request carries cryptographic identity (Ed25519 signatures)</li>\n<li><strong>Authorization:</strong> RBAC policies determine what each role can do. Tool calls are checked against policies before execution</li>\n<li><strong>Audit:</strong> Every action is logged to an immutable audit trail with timestamps, agent IDs, and full argument payloads</li>\n</ul>\n<h2>RBAC Enforcement in Practice</h2>\n<p>TormentNexus intercepts every tool call before execution. Destructive shell operations are blocked by default unless explicitly authorized. The check runs on shell tools, not on file content — preventing false positives while maintaining security.</p>\n<h2>The Result</h2>\n<p>Security that scales with autonomy. As agents gain more capabilities, the zero-trust framework ensures they cannot exceed their authorized scope. Every action is visible, every permission is explicit, and every violation is caught in real-time.</p>',
    },
]

for p in posts:
    html = f"""<!DOCTYPE html>
<html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{p["title"]} — TormentNexus Blog</title>
<style>*{{margin:0;padding:0;box-sizing:border-box}}body{{background:#0a0a0f;color:#e0e0e0;font-family:system-ui,sans-serif;line-height:1.8}}.c{{max-width:740px;margin:0 auto;padding:60px 20px}}h1{{font-size:2.2rem;background:linear-gradient(135deg,#7c3aed,#3b82f6);-webkit-background-clip:text;-webkit-text-fill-color:transparent}}.m{{color:#6868a0;font-size:.85rem;margin-bottom:2rem}}h2{{font-size:1.4rem;margin:2rem 0 .8rem;color:#c4b5fd}}p{{margin-bottom:1.2rem;color:#b8b8c8}}a{{color:#a78bfa}}.b{{display:inline-block;margin-bottom:2rem;color:#6868a0;text-decoration:none}}</style>
</head><body><div class="c">
<a href="/" class="b">← TormentNexus</a>
<h1>{p["title"]}</h1>
<p class="m">{p["date"]} · TormentNexus Team</p>
{p["body"]}
<p style="margin-top:3rem;border-top:1px solid #1a1a2e;padding-top:2rem;color:#6868a0;font-size:.85rem">
<a href="https://github.com/MDMAtk/TormentNexus">GitHub</a> · <a href="https://hypernexus.site/docs">Docs</a></p>
</div></body></html>"""
    with open(f"{blogdir}/{p['slug']}.html", "w", encoding="utf-8") as f:
        f.write(html)

print(f"Wrote {len(posts)} blog posts")
