#!/usr/bin/env python3
"""Generate the TormentNexus catalog page"""

import sqlite3

db = sqlite3.connect("/opt/tormentnexus/catalog.db")

# Get sample entries
samples = {}
for src, cnt in db.execute(
    "SELECT source, count(*) FROM links_backlog GROUP BY source ORDER BY count(*) DESC"
):
    samples[src] = cnt

total_mcp = db.execute("SELECT count(*) FROM links_backlog").fetchone()[0]
db.close()

db2 = sqlite3.connect("/root/.tormentnexus/catalog.db")
skills = db2.execute("SELECT count(*) FROM published_skills").fetchone()[0]
db2.close()

html = f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>TormentNexus — The Index</title>
<style>
*{{margin:0;padding:0;box-sizing:border-box}}
body{{background:#07070f;color:#e0e0e0;font-family:system-ui,-apple-system,sans-serif;line-height:1.6}}
.hero{{text-align:center;padding:80px 20px 40px;background:linear-gradient(180deg,#0a0a1a 0%,#07070f 100%)}}
.hero h1{{font-size:3.5rem;font-weight:800;background:linear-gradient(135deg,#ff006e,#8338ec,#3a86ff);-webkit-background-clip:text;-webkit-text-fill-color:transparent;margin-bottom:0.5rem}}
.hero p{{color:#888;font-size:1.2rem;max-width:600px;margin:0 auto}}
.stats{{display:flex;justify-content:center;gap:40px;padding:30px 20px;flex-wrap:wrap}}
.stat{{text-align:center}}
.stat .num{{font-size:2.5rem;font-weight:700;background:linear-gradient(135deg,#ff006e,#8338ec);-webkit-background-clip:text;-webkit-text-fill-color:transparent}}
.stat .label{{color:#666;font-size:0.85rem;text-transform:uppercase;letter-spacing:1px}}
.btn{{display:inline-block;padding:14px 32px;border-radius:8px;font-weight:600;text-decoration:none;transition:all .2s;margin:6px}}
.btn-gh{{background:#238636;color:#fff;font-size:1.2rem;border:2px solid #2ea043}}
.btn-gh:hover{{background:#2ea043;transform:translateY(-2px)}}
.btn-npm{{background:#cb3837;color:#fff}}
.btn-npm:hover{{background:#e03e36}}
.btn-docs{{background:transparent;color:#8338ec;border:2px solid #8338ec}}
.btn-docs:hover{{background:#8338ec;color:#fff}}
.container{{max-width:1200px;margin:0 auto;padding:0 20px}}
.section{{padding:60px 0}}
.section h2{{font-size:2rem;margin-bottom:8px;background:linear-gradient(135deg,#ff006e,#3a86ff);-webkit-background-clip:text;-webkit-text-fill-color:transparent}}
.section .sub{{color:#666;margin-bottom:30px}}
.grid{{display:grid;grid-template-columns:repeat(auto-fill,minmax(280px,1fr));gap:16px}}
.card{{background:#111122;border:1px solid #1a1a3a;border-radius:10px;padding:20px;transition:all .2s}}
.card:hover{{border-color:#8338ec;transform:translateY(-2px);box-shadow:0 4px 20px rgba(131,56,236,0.15)}}
.card h3{{font-size:1rem;margin-bottom:6px;color:#fff}}
.card p{{color:#888;font-size:0.85rem;margin-bottom:8px}}
.card .meta{{display:flex;gap:8px;flex-wrap:wrap}}
.badge{{display:inline-block;padding:2px 8px;border-radius:4px;font-size:0.7rem;font-weight:600}}
.badge-mcp{{background:rgba(131,56,236,0.2);color:#a78bfa}}
.badge-npm{{background:rgba(203,56,55,0.2);color:#f87171}}
.badge-gh{{background:rgba(35,134,54,0.2);color:#4ade80}}
.badge-pypi{{background:rgba(59,130,246,0.2);color:#60a5fa}}
.badge-prompt{{background:rgba(245,158,11,0.2);color:#fbbf24}}
.badge-agent{{background:rgba(6,182,212,0.2);color:#22d3ee}}
.badge-skill{{background:rgba(236,72,153,0.2);color:#f472b6}}
.clients{{display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:12px}}
.client{{display:flex;align-items:center;gap:10px;padding:12px 16px;background:#111122;border:1px solid #1a1a3a;border-radius:8px;font-size:0.9rem}}
.client .icon{{width:24px;height:24px;background:linear-gradient(135deg,#8338ec,#3a86ff);border-radius:4px;display:flex;align-items:center;justify-content:center;font-size:0.7rem;font-weight:700;color:#fff}}
.matrix{{width:100%;border-collapse:collapse;margin:20px 0}}
.matrix th{{background:#111122;padding:10px 12px;text-align:left;font-size:0.8rem;color:#888;text-transform:uppercase;border-bottom:1px solid #1a1a3a}}
.matrix td{{padding:10px 12px;border-bottom:1px solid #111;font-size:0.9rem}}
.matrix tr:hover{{background:#111122}}
.check{{color:#4ade80;font-weight:700}}
.cross{{color:#333}}
footer{{border-top:1px solid #1a1a3a;padding:40px 20px;text-align:center;color:#444;font-size:0.85rem}}
footer a{{color:#8338ec;text-decoration:none}}
footer a:hover{{text-decoration:underline}}
</style>
</head>
<body>

<div class="hero">
<h1>The Index</h1>
<p>The largest open-source catalog of MCP servers, skills, prompts, and agent tools. Curated from 45+ sources.</p>
<div style="margin-top:30px">
<a href="https://github.com/MDMAtk/TormentNexus" class="btn btn-gh">⭐ GitHub</a>
<a href="https://www.npmjs.com/search?q=%40tormentnexus" class="btn btn-npm">📦 npm</a>
<a href="https://hypernexus.site/docs" class="btn btn-docs">📖 Docs</a>
</div>
</div>

<div class="stats">
<div class="stat"><div class="num">{total_mcp:,}</div><div class="label">MCP Servers</div></div>
<div class="stat"><div class="num">{skills:,}</div><div class="label">Skills</div></div>
<div class="stat"><div class="num">5,668</div><div class="label">Go Handlers</div></div>
<div class="stat"><div class="num">295</div><div class="label">Prompts</div></div>
<div class="stat"><div class="num">38+</div><div class="label">AI Clients</div></div>
</div>

<div class="container">

<div class="section">
<h2>MCP Server Index</h2>
<p class="sub">{total_mcp:,} servers from 20+ registries. The largest MCP catalog on the internet.</p>
<div class="grid">
<div class="card"><h3>awesome-mcp-servers</h3><p>Community-curated list of MCP servers. The definitive source.</p><div class="meta"><span class="badge badge-gh">GitHub</span><span class="badge badge-mcp">{samples.get("awesome-mcp-servers", 0):,} servers</span></div></div>
<div class="card"><h3>awesome-mcp-fork</h3><p>10,000+ forks of the awesome-mcp-servers repository.</p><div class="meta"><span class="badge badge-gh">GitHub</span><span class="badge badge-mcp">{samples.get("awesome-mcp-fork", 0):,} servers</span></div></div>
<div class="card"><h3>npm MCP Packages</h3><p>MCP server packages published to the npm registry.</p><div class="meta"><span class="badge badge-npm">npm</span><span class="badge badge-mcp">{samples.get("npm", 0):,} packages</span></div></div>
<div class="card"><h3>GitHub Search</h3><p>MCP servers discovered via GitHub API search across 20+ queries.</p><div class="meta"><span class="badge badge-gh">GitHub</span><span class="badge badge-mcp">{samples.get("github-search", 0):,} repos</span></div></div>
<div class="card"><h3>PyPI</h3><p>Python MCP server packages from the Python Package Index.</p><div class="meta"><span class="badge badge-pypi">PyPI</span><span class="badge badge-mcp">{samples.get("pypi", 0):,} packages</span></div></div>
<div class="card"><h3>crates.io</h3><p>Rust MCP server packages from the Rust package registry.</p><div class="meta"><span class="badge badge-gh">crates.io</span><span class="badge badge-mcp">{samples.get("crates-io", 0):,} crates</span></div></div>
<div class="card"><h3>Docker Hub</h3><p>Containerized MCP server images ready to deploy.</p><div class="meta"><span class="badge badge-gh">Docker</span><span class="badge badge-mcp">{samples.get("docker-hub", 0):,} images</span></div></div>
<div class="card"><h3>Glama.ai</h3><p>MCP servers from the Glama.ai directory.</p><div class="meta"><span class="badge badge-gh">Glama</span><span class="badge badge-mcp">{samples.get("glama-html", 0):,} servers</span></div></div>
</div>
</div>

<div class="section">
<h2>Skills Index</h2>
<p class="sub">{skills:,} reusable skill modules for code review, Terraform, database migrations, and more.</p>
<div class="grid">
<div class="card"><h3>prompt_engineering</h3><p>Advanced prompt engineering patterns and techniques.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>code_review</h3><p>Automated code review with best practices and security checks.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>terraform_planner</h3><p>Infrastructure as code planning and validation.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>database_migration</h3><p>Safe database schema migrations with rollback support.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>security_audit</h3><p>Security vulnerability scanning and remediation.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>api_design</h3><p>RESTful and GraphQL API design patterns.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>testing_strategy</h3><p>Comprehensive testing strategy and coverage analysis.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>documentation</h3><p>Auto-generate and maintain project documentation.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>performance_optimization</h3><p>Application performance profiling and optimization.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>ci_cd_pipeline</h3><p>Continuous integration and deployment pipelines.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>docker_optimization</h3><p>Container optimization and best practices.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
<div class="card"><h3>monitoring_setup</h3><p>Observability and monitoring configuration.</p><div class="meta"><span class="badge badge-skill">Skill</span></div></div>
</div>
</div>

<div class="section">
<h2>Prompt & System Prompt Index</h2>
<p class="sub">295 curated prompts and system prompts from 12 sources.</p>
<div class="grid">
<div class="card"><h3>Fabric Patterns</h3><p>45 production-ready prompt patterns from danielmiessler/Fabric.</p><div class="meta"><span class="badge badge-prompt">Prompt</span><span class="badge badge-mcp">45 patterns</span></div></div>
<div class="card"><h3>System Prompts Leaks</h3><p>Leaked system prompts from major AI tools and platforms.</p><div class="meta"><span class="badge badge-prompt">Prompt</span></div></div>
<div class="card"><h3>Claude Code Prompts</h3><p>System prompts and instructions for Claude Code.</p><div class="meta"><span class="badge badge-prompt">Prompt</span></div></div>
<div class="card"><h3>ChatGPT System Prompts</h3><p>Collection of ChatGPT system prompts and configurations.</p><div class="meta"><span class="badge badge-prompt">Prompt</span></div></div>
<div class="card"><h3>AI System Prompts</h3><p>Curated system prompts for top AI tools and platforms.</p><div class="meta"><span class="badge badge-prompt">Prompt</span></div></div>
<div class="card"><h3>Big Prompt Library</h3><p>Comprehensive collection of prompts and LLM instructions.</p><div class="meta"><span class="badge badge-prompt">Prompt</span></div></div>
<div class="card"><h3>MCP Design Prompts</h3><p>Prompts for designing MCP servers and tools.</p><div class="meta"><span class="badge badge-prompt">Prompt</span></div></div>
<div class="card"><h3>Agent Prompts</h3><p>System prompts for AI agent frameworks.</p><div class="meta"><span class="badge badge-prompt">Prompt</span></div></div>
</div>
</div>

<div class="section">
<h2>Supported AI Clients</h2>
<p class="sub">TormentNexus installs and configures 38+ AI coding agents automatically.</p>
<div class="clients">
<div class="client"><div class="icon">C</div><span>Claude Code</span></div>
<div class="client"><div class="icon">C</div><span>Cursor</span></div>
<div class="client"><div class="icon">W</div><span>Windsurf</span></div>
<div class="client"><div class="icon">G</div><span>Gemini</span></div>
<div class="client"><div class="icon">C</div><span>Codex CLI</span></div>
<div class="client"><div class="icon">G</div><span>Grok Build</span></div>
<div class="client"><div class="icon">A</div><span>Antigravity</span></div>
<div class="client"><div class="icon">A</div><span>Aider</span></div>
<div class="client"><div class="icon">C</div><span>CodeWhale</span></div>
<div class="client"><div class="icon">G</div><span>Goose</span></div>
<div class="client"><div class="icon">P</div><span>Pi</span></div>
<div class="client"><div class="icon">C</div><span>Cline / Roo</span></div>
<div class="client"><div class="icon">C</div><span>Continue</span></div>
<div class="client"><div class="icon">Z</div><span>Zed</span></div>
<div class="client"><div class="icon">V</div><span>VS Code</span></div>
<div class="client"><div class="icon">J</div><span>JetBrains</span></div>
<div class="client"><div class="icon">O</div><span>OpenHands</span></div>
<div class="client"><div class="icon">A</div><span>Augment</span></div>
<div class="client"><div class="icon">C</div><span>Copilot</span></div>
<div class="client"><div class="icon">T</div><span>Tabnine</span></div>
<div class="client"><div class="icon">C</div><span>Codeium</span></div>
<div class="client"><div class="icon">S</div><span>Sourcegraph</span></div>
<div class="client"><div class="icon">R</div><span>Replit</span></div>
<div class="client"><div class="icon">P</div><span>Phind</span></div>
<div class="client"><div class="icon">B</div><span>Blackbox</span></div>
<div class="client"><div class="icon">M</div><span>MarsCode</span></div>
<div class="client"><div class="icon">T</div><span>Trae</span></div>
<div class="client"><div class="icon">F</div><span>Factory</span></div>
<div class="client"><div class="icon">K</div><span>Kiro</span></div>
<div class="client"><div class="icon">K</div><span>Kimi-Code</span></div>
<div class="client"><div class="icon">Q</div><span>Qwen-Code</span></div>
<div class="client"><div class="icon">O</div><span>OmniGent</span></div>
<div class="client"><div class="icon">C</div><span>Citadel</span></div>
<div class="client"><div class="icon">A</div><span>Agent-Fusion</span></div>
<div class="client"><div class="icon">H</div><span>Herdr</span></div>
<div class="client"><div class="icon">C</div><span>Claude-Squad</span></div>
<div class="client"><div class="icon">O</div><span>OpenCode</span></div>
<div class="client"><div class="icon">C</div><span>Cursor Nightly</span></div>
</div>
</div>

<div class="section">
<h2>Client Integration Matrix</h2>
<p class="sub">What each client receives when TormentNexus is installed.</p>
<table class="matrix">
<thead>
<tr><th>Client</th><th>SKILL</th><th>MCP</th><th>CMD</th><th>HOOK</th><th>EXT</th><th>AGENT</th></tr>
</thead>
<tbody>
<tr><td>Claude Code</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Cursor</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Windsurf</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Gemini</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td></tr>
<tr><td>Codex CLI</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Grok Build</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Antigravity</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td></tr>
<tr><td>Aider</td><td><span class="cross">—</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td></tr>
<tr><td>CodeWhale</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Goose</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Pi</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Cline / Roo</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Continue</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Zed</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>VS Code</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>JetBrains</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>OpenHands</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Augment</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
<tr><td>Copilot</td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td><td><span class="check">✓</span></td><td><span class="cross">—</span></td></tr>
</tbody>
</table>
</div>

<div class="section">
<h2>Quick Install</h2>
<p class="sub">One command. 38+ clients. Zero configuration.</p>
<div style="text-align:center">
<pre style="background:#111122;padding:20px;border-radius:8px;display:inline-block;text-align:left;border:1px solid #1a1a3a"><code style="color:#a78bfa;font-size:1.1rem">npx @tormentnexus/install</code></pre>
</div>
</div>

</div>

<footer>
<p>TormentNexus v1.0.0-b2 · MIT License · <a href="https://github.com/MDMAtk/TormentNexus">GitHub</a> · <a href="https://hypernexus.site">Corporate Edition</a></p>
</footer>

</body>
</html>"""

with open("/var/www/tormentnexus.site/catalog.html", "w") as f:
    f.write(html)

print(f"Catalog page written: {len(html)} bytes")
