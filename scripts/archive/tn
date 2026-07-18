#!/usr/bin/env python3
"""
tn — TormentNexus CLI
─────────────────────
Universal command-line tool for TN Kernel interaction.
Works standalone — no AI agent required.

Usage:
  tn search "deployment strategy"
  tn store "Use blue-green deploys" --tags deployment,devops
  tn status
  tn sessions
  tn tools
  tn skills
  tn harvest "my current task"
  tn purge --tag stale
  tn serve
"""
import sys, os, json, urllib.request, urllib.error, subprocess, argparse

TN = os.environ.get("TN_URL", "http://127.0.0.1:7778")

def fetch(path, method="GET", body=None):
    url = f"{TN}{path}"
    data = json.dumps(body).encode() if body else None
    req = urllib.request.Request(url, data=data, method=method)
    req.add_header("Content-Type", "application/json")
    try:
        with urllib.request.urlopen(req, timeout=10) as r:
            return json.loads(r.read())
    except urllib.error.URLError as e:
        return {"error": f"TN Kernel unreachable at {TN}", "detail": str(e)}
    except json.JSONDecodeError:
        return {"error": "Invalid response from TN Kernel"}

def cmd_search(args):
    q = " ".join(args.query) if isinstance(args.query, list) else args.query
    result = fetch(f"/api/memory/search?q={urllib.parse.quote(q)}&limit={args.limit}")
    memories = result.get("memories", result.get("data", []))
    if isinstance(memories, dict):
        memories = memories.get("memories", [])
    if not memories:
        print("No memories found.")
        return
    for i, m in enumerate(memories[:args.limit]):
        content = (m.get("content") or str(m))[:120]
        tags = m.get("tags", [])
        tag_str = f" [{', '.join(tags)}]" if tags else ""
        print(f"{i+1}. {content}{tag_str}")

def cmd_store(args):
    content = " ".join(args.content) if isinstance(args.content, list) else args.content
    tags = [t.strip() for t in args.tags.split(",")] if args.tags else []
    result = fetch("/api/memory/add", method="POST", body={
        "content": content, "tags": tags, "category": "cli"
    })
    if result.get("success") or result.get("id"):
        print(f"Stored: {content[:80]}...")
    else:
        print(f"Error: {result.get('error', 'unknown')}")

def cmd_status(args):
    result = fetch("/api/runtime/status")
    if result.get("error"):
        print(f"Error: {result['error']}")
        return
    data = result.get("data", result)
    print(f"TN Kernel : v{data.get('version','?')}")
    print(f"Uptime   : {data.get('uptimeSec',0)//3600}h {data.get('uptimeSec',0)%3600//60}m")
    cli = data.get("cli", {})
    print(f"Tools    : {cli.get('toolCount','?')} total, {cli.get('availableToolCount','?')} available")
    prov = data.get("providers", {})
    print(f"Providers: {prov.get('configuredCount','?')}/{prov.get('providerCount','?')} configured")
    locks = data.get("locks", [])
    if locks:
        print(f"Locks    : {sum(1 for l in locks if l.get('running'))}/{len(locks)} running")

def cmd_sessions(args):
    result = fetch("/api/import/summary")
    data = result.get("data", result)
    print(f"Imported sessions: {data.get('count','?')} total, {data.get('validCount','?')} valid")
    for item in data.get("bySourceTool", []):
        print(f"  {item.get('key','?')}: {item.get('count','?')}")

def cmd_tools(args):
    result = fetch("/api/cli/summary")
    data = result.get("data", result)
    print(f"Tools: {data.get('toolCount','?')} total, {data.get('availableToolCount','?')} available")
    print(f"Harnesses: {data.get('harnessCount','?')} total, {data.get('installedHarnessCount','?')} installed")

def cmd_skills(args):
    print("Skills: 5,776 available via tn_skill_manage MCP tool")
    print("Use an AI agent with TN integration for interactive skill search.")

def cmd_harvest(args):
    q = " ".join(args.prompt) if isinstance(args.prompt, list) else args.prompt
    result = fetch(f"/api/memory/search?q={urllib.parse.quote(q)}&limit=5")
    memories = result.get("memories", result.get("data", []))
    if isinstance(memories, dict):
        memories = memories.get("memories", [])
    if not memories:
        print("No relevant context found.")
        return
    print("Relevant context from L2 memory:\n")
    for m in memories[:5]:
        content = (m.get("content") or str(m))[:200]
        print(f"  - {content}")

def cmd_purge(args):
    tag = args.tag or "stale"
    result = fetch(f"/api/memory/purge?tag={urllib.parse.quote(tag)}", method="POST")
    print(result.get("message", result.get("status", f"Purged memories tagged: {tag}")))

def cmd_serve(args):
    binary = "tormentnexus.exe" if sys.platform == "win32" else "tormentnexus"
    port = args.port or 7778
    print(f"Starting TN Kernel on port {port}...")
    subprocess.run([binary, "serve", "-port", str(port), "-host", "127.0.0.1"])

def main():
    parser = argparse.ArgumentParser(description="tn — TormentNexus CLI")
    sub = parser.add_subparsers(dest="command")

    p = sub.add_parser("search", help="Search L2 memory")
    p.add_argument("query", nargs="+", help="Search query")
    p.add_argument("--limit", type=int, default=10)

    p = sub.add_parser("store", help="Store a memory")
    p.add_argument("content", nargs="+", help="Memory content")
    p.add_argument("--tags", default="")

    sub.add_parser("status", help="System health")

    sub.add_parser("sessions", help="Imported session summary")

    sub.add_parser("tools", help="Tool inventory")

    sub.add_parser("skills", help="Skill registry info")

    p = sub.add_parser("harvest", help="Harvest L2 context")
    p.add_argument("prompt", nargs="+", help="Task description")

    p = sub.add_parser("purge", help="Remove stale memories")
    p.add_argument("--tag", default="stale")

    p = sub.add_parser("serve", help="Start TN Kernel")
    p.add_argument("--port", type=int, default=7778)

    args = parser.parse_args()
    if not args.command:
        parser.print_help()
        return

    {"search": cmd_search, "store": cmd_store, "status": cmd_status,
     "sessions": cmd_sessions, "tools": cmd_tools, "skills": cmd_skills,
     "harvest": cmd_harvest, "purge": cmd_purge, "serve": cmd_serve,
    }[args.command](args)

if __name__ == "__main__":
    main()
