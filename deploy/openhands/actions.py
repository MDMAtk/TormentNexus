"""
TormentNexus OpenHands Actions
Exposes all TN MCP tools as OpenHands actions for LLM agents.
Place in: ~/.openhands/plugins/tormentnexus/actions.py
"""

import json
import urllib.request
from typing import Any, Dict, List, Optional

TN_URL = "http://127.0.0.1:7778"


def _tn(endpoint: str, data: Optional[Dict] = None) -> Any:
    url = f"{TN_URL}/api/{endpoint}"
    body = json.dumps(data).encode() if data else None
    req = urllib.request.Request(
        url,
        body,
        {"Content-Type": "application/json"},
        method="POST" if data else "GET",
    )
    with urllib.request.urlopen(req, timeout=10) as r:
        return json.loads(r.read())


# ─── Memory Actions ───


def tn_memory_store(content: str, tags: List[str], category: str = "general") -> str:
    """Store a fact, decision, or pattern to TN L2 memory.
    Args: content (str): What to store. tags (list[str]): Tags for retrieval.
          category (str): general, decision, pattern, bug, or architecture."""
    return json.dumps(
        _tn("memory/store", {"content": content, "tags": tags, "category": category})
    )


def tn_memory_search(
    query: str, tags: Optional[List[str]] = None, limit: int = 10
) -> str:
    """Search TN L2 memory by keyword, tag, or category.
    Args: query (str): Search terms. tags (list[str], optional): Filter by tags.
          limit (int): Max results (default 10)."""
    return json.dumps(
        _tn("memory/search", {"query": query, "tags": tags or [], "limit": limit})
    )


def tn_memory_vector_search(query: str, limit: int = 5) -> str:
    """Semantic vector search across L2 memory (finds conceptually similar entries).
    Args: query (str): Natural language query. limit (int): Max results (default 5)."""
    return json.dumps(_tn("memory/vector-search", {"query": query, "limit": limit}))


# ─── Tool Discovery ───


def tn_tool_search(query: str) -> str:
    """Find MCP tools across all configured servers by describing what you need.
    Args: query (str): Description of the tool you need (e.g., 'search code', 'query database')."""
    return json.dumps(_tn("tools/search", {"query": query}))


def tn_tool_execute(tool: str, arguments: Dict[str, Any]) -> str:
    """Execute a named TN MCP tool with arguments.
    Args: tool (str): Tool name (from tn_tool_search). arguments (dict): Tool parameters."""
    return json.dumps(_tn("mcp/execute", {"tool": tool, "arguments": arguments}))


# ─── Sessions ───


def tn_session_search(query: str) -> str:
    """Search 542+ imported sessions from Claude Code, Aider, Gemini, etc.
    Args: query (str): What topic, tool, or pattern to search for."""
    return json.dumps(_tn("sessions/search", {"query": query}))


def tn_session_import(source: str, path: str) -> str:
    """Import a session from an external AI tool into TN.
    Args: source (str): claude, aider, gemini, copilot, or cursor.
          path (str): Path to the session file."""
    return json.dumps(_tn("sessions/import", {"source": source, "path": path}))


# ─── Code Search ───


def tn_code_search(query: str, scope: str = "all", method: str = "ast-grep") -> str:
    """Search code via AST-grep rules, semantic search, or file patterns.
    Args: query (str): Pattern or query. scope (str): all, go, ts, py, or rust.
          method (str): ast-grep, semantic, or pattern."""
    return json.dumps(
        _tn("code/search", {"query": query, "scope": scope, "method": method})
    )


# ─── System ───


def tn_system_status() -> str:
    """Get TN health overview: version, uptime, memory tiers, connected servers, provider quota."""
    return json.dumps(_tn("runtime/status"))


def tn_context_harvest(prompt: str) -> str:
    """Pull relevant L2 memory context for the given task/prompt.
    Args: prompt (str): What you're working on — TN finds related memories."""
    return json.dumps(_tn("context/harvest", {"prompt": prompt}))


def tn_skill_manage(action: str, query: str = "") -> str:
    """Access the TN skill registry (5,776+ modules).
    Args: action (str): list, search, install, remove, or describe.
          query (str): Skill name or search term."""
    return json.dumps(_tn("skills/manage", {"action": action, "query": query}))


# OpenHands action definitions (used by the action parser)
ACTIONS = {
    "tn_memory_store": {
        "function": tn_memory_store,
        "description": "Store a fact, decision, or pattern to TormentNexus persistent L2 memory",
        "args_schema": {
            "content": {"type": "string", "description": "Content to store"},
            "tags": {
                "type": "array",
                "items": {"type": "string"},
                "description": "Tags for retrieval",
            },
            "category": {
                "type": "string",
                "default": "general",
                "description": "general, decision, pattern, bug, architecture",
            },
        },
    },
    "tn_memory_search": {
        "function": tn_memory_search,
        "description": "Search TormentNexus L2 memory by keyword, tag, or category",
        "args_schema": {
            "query": {"type": "string", "description": "Search terms"},
            "tags": {
                "type": "array",
                "items": {"type": "string"},
                "description": "Filter tags",
            },
            "limit": {"type": "integer", "default": 10, "description": "Max results"},
        },
    },
    "tn_memory_vector_search": {
        "function": tn_memory_vector_search,
        "description": "Semantic vector search across TormentNexus L2 memory",
        "args_schema": {
            "query": {"type": "string", "description": "Natural language query"},
            "limit": {"type": "integer", "default": 5, "description": "Max results"},
        },
    },
    "tn_tool_search": {
        "function": tn_tool_search,
        "description": "Find MCP tools across all TormentNexus servers by describing what you need",
        "args_schema": {
            "query": {
                "type": "string",
                "description": "What kind of tool or capability you need",
            },
        },
    },
    "tn_tool_execute": {
        "function": tn_tool_execute,
        "description": "Execute a named TormentNexus MCP tool",
        "args_schema": {
            "tool": {"type": "string", "description": "Tool name"},
            "arguments": {"type": "object", "description": "Tool parameters"},
        },
    },
    "tn_session_search": {
        "function": tn_session_search,
        "description": "Search 542+ imported AI sessions for patterns and solutions",
        "args_schema": {
            "query": {
                "type": "string",
                "description": "Topic, tool, or pattern to find in past sessions",
            },
        },
    },
    "tn_code_search": {
        "function": tn_code_search,
        "description": "Search codebase via AST-grep, semantic, or pattern matching",
        "args_schema": {
            "query": {"type": "string", "description": "Pattern or query"},
            "scope": {
                "type": "string",
                "default": "all",
                "description": "all, go, ts, py, rust",
            },
            "method": {
                "type": "string",
                "default": "ast-grep",
                "description": "ast-grep, semantic, pattern",
            },
        },
    },
    "tn_system_status": {
        "function": tn_system_status,
        "description": "Get TormentNexus health overview: version, uptime, memory, quota",
        "args_schema": {},
    },
    "tn_context_harvest": {
        "function": tn_context_harvest,
        "description": "Pull relevant L2 memory context for your current task",
        "args_schema": {
            "prompt": {"type": "string", "description": "What you're working on"},
        },
    },
    "tn_skill_manage": {
        "function": tn_skill_manage,
        "description": "Access the TormentNexus skill registry (5,776+ modules)",
        "args_schema": {
            "action": {
                "type": "string",
                "description": "list, search, install, remove, describe",
            },
            "query": {
                "type": "string",
                "default": "",
                "description": "Skill name or search term",
            },
        },
    },
}
