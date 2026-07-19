"""
TormentNexus MCP SSE Bridge
Launches `tormentnexus.exe mcp` as a subprocess and exposes it via SSE transport
for OpenHands to connect to. This bridges the stdio-based TN MCP server to HTTP SSE.

Usage:
  python tn_mcp_sse_bridge.py [--port 9090] [--tn-binary path\to\tormentnexus.exe]
"""

import asyncio
import json
import os
import subprocess
import sys
import uuid
from fastapi import FastAPI, Request
from fastapi.responses import StreamingResponse
import uvicorn

app = FastAPI(title="TormentNexus MCP SSE Bridge")

TN_BINARY = os.environ.get(
    "TN_MCP_BINARY", r"C:\Users\hyper\workspace\tormentnexus\bin\tormentnexus.exe"
)
TN_WORKSPACE = os.environ.get(
    "TORMENTNEXUS_WORKSPACE_ROOT", r"C:\Users\hyper\workspace\tormentnexus"
)

# Store persistent TN MCP subprocess
tn_process = None
tn_lock = asyncio.Lock()


async def get_tn_process():
    """Get or create the TN MCP subprocess singleton."""
    global tn_process
    async with tn_lock:
        if tn_process is None or tn_process.poll() is not None:
            tn_process = subprocess.Popen(
                [TN_BINARY, "mcp"],
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                cwd=TN_WORKSPACE,
                env={**os.environ, "TORMENTNEXUS_WORKSPACE_ROOT": TN_WORKSPACE},
                text=True,
                bufsize=1,
            )
            print(f"[TN MCP Bridge] Started TN MCP process (PID: {tn_process.pid})")
        return tn_process


async def call_tn_mcp(request: dict) -> dict:
    """Send a JSON-RPC request to the TN MCP process and get the response."""
    proc = await get_tn_process()
    line = json.dumps(request)
    proc.stdin.write(line + "\n")
    proc.stdin.flush()
    response_line = proc.stdout.readline()
    if not response_line:
        proc.terminate()
        tn_process = None
        raise RuntimeError("TN MCP process died")
    return json.loads(response_line)


@app.get("/health")
async def health():
    return {"status": "ok", "service": "tormentnexus-mcp-sse-bridge"}


@app.get("/sse")
async def sse_endpoint(request: Request):
    """SSE endpoint for MCP transport."""
    session_id = str(uuid.uuid4())

    async def event_generator():
        # Send endpoint event first
        yield f"event: endpoint\ndata: {json.dumps({'endpoint': f'/messages?session_id={session_id}'})}\n\n"
        try:
            while True:
                if await request.is_disconnected():
                    break
                await asyncio.sleep(30)  # Keep-alive ping interval
                yield f"event: keepalive\ndata: {json.dumps({'session': session_id})}\n\n"
        except asyncio.CancelledError:
            pass

    return StreamingResponse(
        event_generator(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no",
        },
    )


@app.post("/messages")
async def messages(request: Request):
    """Handle JSON-RPC messages from the MCP client."""
    body = await request.json()
    method = body.get("method", "")
    params = body.get("params", {})

    if method == "initialize":
        # Return TN MCP server capabilities
        return {
            "jsonrpc": "2.0",
            "id": body.get("id"),
            "result": {
                "protocolVersion": "2024-11-05",
                "capabilities": {
                    "tools": {},
                    "resources": {},
                    "prompts": {},
                },
                "serverInfo": {
                    "name": "tormentnexus",
                    "version": "1.0.0",
                },
            },
        }

    elif method == "tools/list":
        # Get tools from TN MCP subprocess
        try:
            resp = await call_tn_mcp(
                {
                    "jsonrpc": "2.0",
                    "id": body.get("id", 1),
                    "method": "tools/list",
                    "params": {},
                }
            )
            # Add TN API tools (memory, skills, sessions)
            extra_tools = [
                {
                    "name": "tn_memory_store",
                    "description": "Store a memory in TormentNexus L2 vault",
                    "inputSchema": {
                        "type": "object",
                        "properties": {
                            "content": {
                                "type": "string",
                                "description": "Memory content",
                            },
                            "tags": {
                                "type": "array",
                                "items": {"type": "string"},
                                "description": "Optional tags",
                            },
                            "category": {
                                "type": "string",
                                "description": "Category: pattern, decision, convention, etc.",
                            },
                        },
                        "required": ["content"],
                    },
                },
                {
                    "name": "tn_memory_search",
                    "description": "Search L2 memory by keyword, tag, or category",
                    "inputSchema": {
                        "type": "object",
                        "properties": {
                            "query": {
                                "type": "string",
                                "description": "Search keyword",
                            },
                            "tag": {
                                "type": "string",
                                "description": "Filter by tag prefix",
                            },
                            "limit": {"type": "number", "description": "Max results"},
                        },
                    },
                },
                {
                    "name": "tn_skill_search",
                    "description": "Search TormentNexus skill registry (5,776+ skills)",
                    "inputSchema": {
                        "type": "object",
                        "properties": {
                            "query": {"type": "string", "description": "Search query"},
                            "limit": {
                                "type": "number",
                                "description": "Max results (default 10)",
                            },
                        },
                        "required": ["query"],
                    },
                },
                {
                    "name": "tn_session_search",
                    "description": "Search 542+ imported AI coding sessions",
                    "inputSchema": {
                        "type": "object",
                        "properties": {
                            "query": {"type": "string", "description": "Search query"},
                            "limit": {
                                "type": "number",
                                "description": "Max results (default 10)",
                            },
                        },
                        "required": ["query"],
                    },
                },
                {
                    "name": "tn_context_harvest",
                    "description": "Harvest relevant context from L2 memory + skills + sessions",
                    "inputSchema": {
                        "type": "object",
                        "properties": {
                            "query": {
                                "type": "string",
                                "description": "What you're working on",
                            },
                        },
                        "required": ["query"],
                    },
                },
            ]
            existing_tools = resp.get("result", {}).get("tools", [])
            resp["result"]["tools"] = existing_tools + extra_tools
            return resp
        except Exception as e:
            return {
                "jsonrpc": "2.0",
                "id": body.get("id"),
                "result": {"tools": []},
                "error": {"code": -32000, "message": str(e)},
            }

    elif method == "tools/call":
        tool_name = params.get("name", "")
        tool_args = params.get("arguments", {})
        tn_api = "http://127.0.0.1:7778"

        try:
            if tool_name.startswith("tn_"):
                # Handle TormentNexus API tools
                import httpx

                async with httpx.AsyncClient() as client:
                    if tool_name == "tn_memory_store":
                        content = tool_args.get("content", "")
                        tags = tool_args.get("tags", [])
                        category = tool_args.get("category", "general")
                        payload = {
                            "content": json.dumps(
                                {
                                    "content": content,
                                    "tags": tags,
                                    "category": category,
                                    "timestamp": __import__("datetime")
                                    .datetime.now()
                                    .isoformat(),
                                }
                            )
                        }
                        r = await client.post(f"{tn_api}/api/memory/add", json=payload)
                        return {
                            "jsonrpc": "2.0",
                            "id": body.get("id"),
                            "result": {
                                "content": [
                                    {
                                        "type": "text",
                                        "text": f"✅ Memory stored ({category})",
                                    }
                                ]
                            },
                        }
                    elif tool_name == "tn_memory_search":
                        query = tool_args.get("query", "")
                        tag = tool_args.get("tag", "")
                        limit = tool_args.get("limit", 20)
                        r = await client.get(f"{tn_api}/api/memory/list")
                        memories = r.json() if r.status_code == 200 else []
                        results = []
                        for m in memories if isinstance(memories, list) else []:
                            try:
                                p = json.loads(m)
                            except:
                                p = {"content": m, "tags": [], "category": "general"}
                            if (
                                query
                                and query.lower() not in p.get("content", "").lower()
                            ):
                                continue
                            if tag and not any(
                                t.startswith(tag) for t in p.get("tags", [])
                            ):
                                continue
                            results.append(p)
                        results = results[:limit]
                        text = (
                            f"📚 {len(results)} memories:\n\n"
                            + "\n\n".join(
                                f"{i + 1}. [{r.get('category', 'general')}] {r.get('tags', [])} {r.get('content', '')[:200]}"
                                for i, r in enumerate(results)
                            )
                            if results
                            else "No matching memories."
                        )
                        return {
                            "jsonrpc": "2.0",
                            "id": body.get("id"),
                            "result": {"content": [{"type": "text", "text": text}]},
                        }
                    elif tool_name == "tn_skill_search":
                        query = tool_args.get("query", "")
                        limit = tool_args.get("limit", 10)
                        r = await client.get(
                            f"{tn_api}/api/skills/search?q={query}&limit={limit}"
                        )
                        data = r.json() if r.status_code == 200 else {}
                        skills = data.get(
                            "skills", data.get("data", {}).get("skills", [])
                        )
                        text = (
                            f"🔍 {len(skills)} skills:\n"
                            + "\n".join(
                                f"  {i + 1}. {s.get('id', s)}"
                                for i, s in enumerate(skills[:limit])
                            )
                            if skills
                            else f"No skills for '{query}'."
                        )
                        return {
                            "jsonrpc": "2.0",
                            "id": body.get("id"),
                            "result": {"content": [{"type": "text", "text": text}]},
                        }
                    elif tool_name == "tn_session_search":
                        query = tool_args.get("query", "")
                        limit = tool_args.get("limit", 10)
                        r = await client.get(
                            f"{tn_api}/api/sessions/imported/list?limit={limit}"
                        )
                        data = r.json() if r.status_code == 200 else []
                        sessions = data if isinstance(data, list) else []
                        text = (
                            f"📋 {len(sessions)} sessions:\n"
                            + "\n".join(
                                f"  {i + 1}. {s.get('sourceTool', '?')} ({s.get('sessionFormat', '?')})"
                                for i, s in enumerate(sessions[:limit])
                            )
                            if sessions
                            else "No sessions found."
                        )
                        return {
                            "jsonrpc": "2.0",
                            "id": body.get("id"),
                            "result": {"content": [{"type": "text", "text": text}]},
                        }
                    elif tool_name == "tn_context_harvest":
                        query = tool_args.get("query", "")
                        # Harvest from memory + skills
                        mem_r = await client.get(
                            f"{tn_api}/api/memory/search?q={query}"
                        )
                        skill_r = await client.get(
                            f"{tn_api}/api/skills/search?q={query}"
                        )
                        mem_data = mem_r.json() if mem_r.status_code == 200 else {}
                        skill_data = (
                            skill_r.json() if skill_r.status_code == 200 else {}
                        )
                        memories = mem_data.get("data", [])
                        skills = skill_data.get(
                            "skills", skill_data.get("data", {}).get("skills", [])
                        )
                        parts = []
                        if memories:
                            parts.append(
                                "## L2 Memory\n"
                                + "\n".join(
                                    f"  • {(m.get('content') or json.dumps(m))[:200]}"
                                    for m in memories[:5]
                                )
                            )
                        if skills:
                            parts.append(
                                "## Skills\n"
                                + "\n".join(f"  • {s.get('id', s)}" for s in skills[:5])
                            )
                        text = (
                            "🌾 Context harvested:\n\n" + "\n\n".join(parts)
                            if parts
                            else f"No context for '{query}'."
                        )
                        return {
                            "jsonrpc": "2.0",
                            "id": body.get("id"),
                            "result": {"content": [{"type": "text", "text": text}]},
                        }
            else:
                # Forward to TN MCP subprocess
                resp = await call_tn_mcp(body)
                return resp
        except Exception as e:
            return {
                "jsonrpc": "2.0",
                "id": body.get("id"),
                "error": {"code": -32000, "message": str(e)},
            }

    else:
        return {
            "jsonrpc": "2.0",
            "id": body.get("id"),
            "result": {},
        }


if __name__ == "__main__":
    port = int(sys.argv[sys.argv.index("--port") + 1]) if "--port" in sys.argv else 9090
    if "--tn-binary" in sys.argv:
        idx = sys.argv.index("--tn-binary")
        TN_BINARY = sys.argv[idx + 1]
    print(f"[TN MCP Bridge] Starting SSE bridge on port {port}")
    print(f"[TN MCP Bridge] TN Binary: {TN_BINARY}")
    print(f"[TN MCP Bridge] TN Workspace: {TN_WORKSPACE}")
    print(f"[TN MCP Bridge] Connect OpenHands to: http://127.0.0.1:{port}/sse")
    uvicorn.run(app, host="0.0.0.0", port=port, log_level="info")
