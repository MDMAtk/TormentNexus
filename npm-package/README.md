# HyperNexus

**AI Control Plane with persistent memory and 26,000+ MCP tools catalog**

[![GitHub Stars](https://img.shields.io/github/stars/MDMAtk/HyperNexus?style=social)](https://gitlab.com/robertpelloni/HyperNexus)
[![npm version](https://img.shields.io/npm/v/hypernexus)](https://www.npmjs.com/package/hypernexus)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## What is HyperNexus?

HyperNexus is an open-source AI control plane that gives your AI agents:

- **Persistent Memory** — Tiered memory system (L1→L2→L3→L4) that remembers context across sessions
- **26,000+ MCP Tools** — Searchable catalog of Model Context Protocol servers
- **Local LLM Support** — Works with LM Studio, Ollama, DeepSeek (no data leaves your machine)
- **Unified Dashboard** — Real-time monitoring of memory, tools, agents, and security

## Quick Start

### Install via npm

```bash
npx hypernexus serve
```

### Install via Docker

```bash
docker run -p 7778:7778 hypernexus/hypernexus
```

### Install from source

```bash
git clone https://gitlab.com/robertpelloni/HyperNexus.git
cd HyperNexus
go build -buildvcs=false -o hypernexus ./cmd/hypernexus
./hypernexus serve
```

## Dashboard

Open <http://127.0.0.1:7778> in your browser to access the dashboard.

### Features

- **Memory Explorer** — Browse and search your AI's memory
- **Tool Catalog** — Search 26,000+ MCP servers
- **Agent Monitor** — Track agent activity and performance
- **Security Center** — Monitor threats and access control

## Architecture

```
┌─────────────────────────────────────────┐
│  Dashboard (Next.js) :7779              │
├─────────────────────────────────────────┤
│  TN Kernel (Go) :7778                   │
│  ├─ Memory System (L1/L2/L3/L4)        │
│  ├─ MCP Catalog (26K+ entries)          │
│  ├─ Tool Registry                       │
│  └─ Agent Orchestrator                  │
├─────────────────────────────────────────┤
│  Storage                                │
│  ├─ SQLite + FTS5                       │
│  ├─ Vector Embeddings                   │
│  └─ Graph Relations                     │
└─────────────────────────────────────────┘
```

## Memory System

| Tier | Name | Purpose | Retention |
|------|------|---------|-----------|
| L1 | Session | Current conversation | Ephemeral |
| L2 | Hot Store | Active knowledge | 30 days |
| L3 | Cold Archive | Historical data | 1 year |
| L4 | Limbo | Soft-deleted items | 90 days |

## MCP Catalog

HyperNexus includes a searchable database of 26,000+ MCP servers:

```bash
# Search the catalog
curl "http://127.0.0.1:7778/api/backlog/search?q=postgres&limit=5"

# Get stats
curl "http://127.0.0.1:7778/api/backlog/stats"

# Get categories
curl "http://127.0.0.1:7778/api/backlog/categories"
```

## Configuration

Configuration file: `~/.hypernexus/config.yaml`

```yaml
host: 127.0.0.1
port: 7778

memory:
  l2_enabled: true
  l3_enabled: true
  l4_enabled: false

providers:
  deepseek:
    enabled: true
    api_key: "your-api-key"
  lmstudio:
    enabled: true
    url: http://127.0.0.1:1234
```

## API Reference

### Memory

- `GET /api/memory/stats` — Memory statistics
- `GET /api/memory/search?q=...` — Search memories
- `POST /api/memory/store` — Store a memory

### Catalog

- `GET /api/backlog/search?q=...` — Search MCP servers
- `GET /api/backlog/stats` — Catalog statistics
- `GET /api/backlog/categories` — List categories

### System

- `GET /health` — Health check
- `GET /api/status` — System status

## Templates

Quick-start templates for common workflows:

```bash
# Cursor-like coding assistant
hypernexus init --template=cursor-killer

# Research assistant
hypernexus init --template=research-assistant

# Code reviewer
hypernexus init --template=code-reviewer
```

## Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Links

- **Website:** <https://hypernexus.site>
- **GitHub:** <https://gitlab.com/robertpelloni/HyperNexus>
- **Documentation:** <https://hypernexus.site/docs>
- **Blog:** <https://hypernexus.site/blog>

## Support

- **Issues:** <https://gitlab.com/robertpelloni/HyperNexus/issues>
- **Discord:** (coming soon)
- **Email:** <dev@hypernexus.org>
