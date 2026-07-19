---
description: Check TormentNexus system status — verify the MCP server is running and healthy.
---

# TormentNexus Status (/tn-status)

When I say "tn-status" or "check TN", check the TormentNexus system health.

1. **Check MCP** — Verify the tormentnexus MCP server is connected: `mcp_tormentnexus_system_status`
2. **Check Server** — Test connectivity: `mcp_tormentnexus_mcp_status`
3. **Report** — Show whether TN is running, what tools are available, and any issues

If TN isn't running: `cd ~/workspace/tormentnexus && start tormentnexus.exe serve`
