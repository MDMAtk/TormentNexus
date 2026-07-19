---
sidebar_position: 1
---

# Getting Started

Welcome to TormentNexus! This guide will help you get up and running in minutes.

## Quick Start

### Option 1: npm (Recommended)

```bash
# Install globally
npm install -g @tormentnexus/core

# Start the server
tormentnexus serve

# Open dashboard
open http://localhost:7778
```

### Option 2: npx

```bash
npx @tormentnexus/core serve
```

### Option 3: Docker

```bash
docker run -p 7778:7778 ghcr.io/mdmatk/tormentnexus:latest
```

### Option 4: Download

- **Windows**: Download `tormentnexus-setup.exe` from [GitHub Releases](https://github.com/MDMAtk/TormentNexus/releases)
- **macOS**: Download `tormentnexus-darwin-*.tar.gz`
- **Linux**: Download `tormentnexus-linux-*.tar.gz`

## First Steps

### 1. Start the Server

```bash
tormentnexus serve
```

The server starts on `http://localhost:7778` by default.

### 2. Open the Dashboard

Open your browser and go to `http://localhost:7778`

You'll see:

- **Mission Control** — System status and controls
- **Memory Explorer** — View your persistent memory
- **MCP Tool Catalog** — Browse 26,000+ tools
- **Settings** — Configure providers and preferences

### 3. Search for Tools

```bash
# Search via API
curl http://localhost:7778/api/backlog/search?q=postgres

# Get stats
curl http://localhost:7778/api/backlog/stats
```

### 4. Add a Memory

```bash
curl -X POST http://localhost:7778/api/memory/add \
  -H "Content-Type: application/json" \
  -d '{"content": "My first memory!", "tags": ["test"]}'
```

## Configuration

Create a config file at `~/.tormentnexus/config.json`:

```json
{
  "server": {
    "host": "127.0.0.1",
    "port": 7778
  },
  "memory": {
    "enabled": true,
    "tiers": ["L1", "L2", "L3", "L4"]
  },
  "mcp": {
    "catalog": true,
    "autoInstall": false
  }
}
```

## Next Steps

- [Architecture](/docs/architecture) — Learn how TormentNexus works
- [API Reference](/docs/api) — Explore the REST API
- [MCP Tools](/docs/tools) — Browse the tool catalog
- [Memory System](/docs/memory) — Understand the 4-tier memory

## Need Help?

- [GitHub Issues](https://github.com/MDMAtk/TormentNexus/issues)
- [Discord](https://discord.gg/Hj9P3GbVxR)
