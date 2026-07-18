---
description: Check HyperNexus system status — verify the MCP server is running and healthy.
---

# HyperNexus Status (/tn-status)

When I say "tn-status" or "check TN", check the HyperNexus system health.

1. **Check MCP** — Verify the hypernexus MCP server is connected: `mcp_hypernexus_system_status`
2. **Check Server** — Test connectivity: `mcp_hypernexus_mcp_status`
3. **Report** — Show whether TN is running, what tools are available, and any issues

If TN isn't running: `cd ~/workspace/hypernexus && start hypernexus.exe serve`
