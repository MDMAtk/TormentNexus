"""
TormentNexus OpenHands Agent Extension
Auto-loads TN memory, tools, and context harvesting into every OpenHands session.

Place in: ~/.openhands/plugins/tormentnexus/tormentnexus_agent.py
"""

import json
import os
import urllib.request
from typing import Any, Dict, List, Optional

TN_URL = os.environ.get("TORMENTNEXUS_URL", "http://127.0.0.1:7778")


class TormentNexusAgent:
    """OpenHands agent wrapper that adds TN memory and tools to every session.

    Usage in OpenHands config:
        [agent]
        name = "tormentnexus"
        base_agent = "CodeActAgent"

    Or auto-loaded via plugin.toml:
        [agent]
        entrypoint = "tormentnexus_agent.py"
        auto_load = true
    """

    def __init__(self, base_agent: Any = None) -> None:
        self.base_agent = base_agent
        self._tools_cache: Optional[List[str]] = None
        self._healthy: Optional[bool] = None

    # ─── HTTP helpers ───

    def _tn_request(
        self, endpoint: str, method: str = "GET", data: Optional[Dict] = None
    ) -> Dict:
        url = f"{TN_URL}/api/{endpoint}"
        body = json.dumps(data).encode() if data else None
        req = urllib.request.Request(
            url, body, {"Content-Type": "application/json"}, method=method
        )
        try:
            with urllib.request.urlopen(req, timeout=10) as r:
                return json.loads(r.read())
        except Exception as e:
            return {"error": str(e)}

    # ─── Health ───

    def is_healthy(self) -> bool:
        """Check if TN Kernel is reachable and healthy."""
        if self._healthy is not None:
            return self._healthy
        try:
            result = self._tn_request("runtime/status")
            self._healthy = "data" in result and "version" in result.get("data", {})
            return self._healthy
        except Exception:
            self._healthy = False
            return False

    def get_status(self) -> Dict:
        """Full TN system status (version, uptime, tools, memory tiers)."""
        return self._tn_request("runtime/status").get("data", {})

    # ─── L2 Memory ───

    def memory_store(
        self, content: str, tags: List[str], category: str = "general"
    ) -> Dict:
        """Store a fact, decision, or pattern to TN L2 vector memory."""
        return self._tn_request(
            "memory/store",
            "POST",
            {"content": content, "tags": tags, "category": category},
        )

    def memory_search(
        self, query: str, tags: Optional[List[str]] = None, limit: int = 10
    ) -> List[Dict]:
        """Keyword + tag search across L2 memory."""
        result = self._tn_request(
            "memory/search",
            "POST",
            {"query": query, "tags": tags or [], "limit": limit},
        )
        return result.get("data", result.get("results", []))

    def memory_vector_search(self, query: str, limit: int = 5) -> List[Dict]:
        """Semantic vector search across L2 memory."""
        result = self._tn_request(
            "memory/vector-search",
            "POST",
            {"query": query, "limit": limit},
        )
        return result.get("data", result.get("results", []))

    def context_harvest(self, prompt: str) -> Dict:
        """Pull relevant context from L2 memory for the given prompt."""
        return self._tn_request("context/harvest", "POST", {"prompt": prompt})

    # ─── Tools ───

    def tool_search(self, query: str) -> List[Dict]:
        """Discover MCP tools across all configured servers by description."""
        result = self._tn_request("tools/search", "POST", {"query": query})
        return result.get("data", result.get("tools", []))

    def list_tools(self) -> List[str]:
        """Get all available TN MCP tool names."""
        if not self._tools_cache:
            status = self._tn_request("runtime/status")
            self._tools_cache = status.get("data", {}).get("cli", {}).get("tools", [])
        return self._tools_cache or []

    def execute_tool(self, tool_name: str, args: Dict) -> Dict:
        """Execute a named TN MCP tool with arguments."""
        return self._tn_request(
            "mcp/execute", "POST", {"tool": tool_name, "arguments": args}
        )

    # ─── Sessions ───

    def session_search(self, query: str) -> List[Dict]:
        """Search imported sessions from other AI agents (Claude Code, Aider, etc)."""
        result = self._tn_request("sessions/search", "POST", {"query": query})
        return result.get("data", result.get("sessions", []))

    # ─── Skills ───

    def skill_manage(self, action: str, query: str = "") -> Dict:
        """Access the TN skill registry (list, search, install, remove)."""
        return self._tn_request(
            "skills/manage", "POST", {"action": action, "query": query}
        )

    # ─── Code Search ───

    def code_search(
        self, query: str, scope: str = "all", method: str = "ast-grep"
    ) -> List[Dict]:
        """Search code via AST-grep, semantic, or pattern matching."""
        return self._tn_request(
            "code/search",
            "POST",
            {"query": query, "scope": scope, "method": method},
        ).get("data", [])

    # ─── OpenHands integration hooks ───

    def on_session_start(self, task: str) -> Dict:
        """Called when an OpenHands session starts. Auto-harvests TN context."""
        if self.is_healthy():
            harvested = self.context_harvest(task)
            return {"harvested": harvested, "tools_loaded": len(self.list_tools())}
        return {"healthy": False, "warning": "TormentNexus not reachable"}

    def on_session_end(self, summary: str, decisions: List[str]) -> Dict:
        """Called when an OpenHands session ends. Stores decisions to L2."""
        results = []
        if self.is_healthy():
            for d in decisions:
                results.append(self.memory_store(d, ["openhands", "decision", "auto"]))
            results.append(
                self.memory_store(summary, ["openhands", "session-summary", "auto"])
            )
        return {"stored": len(results)}
