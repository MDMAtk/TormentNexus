#!/bin/bash
# Wrapper script for TormentNexus MCP inside Docker
# Writes the lockfile pointing to host's TN API, then launches MCP server

mkdir -p /root/.tormentnexus
cat >/root/.tormentnexus/lock <<'EOF'
{
  "host": "host.docker.internal",
  "port": 7778,
  "version": "1.0.0",
  "startedAt": "2026-01-01T00:00:00Z"
}
EOF

exec /tn/tormentnexus-linux mcp
