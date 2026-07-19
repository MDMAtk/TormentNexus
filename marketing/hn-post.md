# Hacker News Post Draft

## Title

Show HN: TormentNexus – Open-source AI control plane with 26K+ MCP tools catalog

## URL
<https://github.com/MDMAtk/TormentNexus>

## Body (first comment)

Hi HN! I built TormentNexus — an open-source AI control plane that gives AI agents persistent memory and access to 26,000+ MCP (Model Context Protocol) servers.

### The Problem

Current AI assistants lose context between sessions. They can't access external tools easily. And if you want to use local LLMs, you're stuck with basic chat interfaces.

### The Solution

TormentNexus provides:

1. **Persistent Memory**: A tiered memory system (L1→L2→L3→L4) that remembers context across sessions. Uses vector embeddings for semantic search.

2. **MCP Tool Catalog**: 26,000+ searchable MCP servers — databases, filesystems, browsers, APIs, DevOps tools. One-click install.

3. **Local LLM Support**: Works with LM Studio, Ollama, DeepSeek. No data leaves your machine.

4. **Unified Dashboard**: Real-time monitoring of memory, tools, agents, security, and infrastructure.

### Tech Stack

- Go backend (single binary, fast startup)
- Next.js dashboard
- SQLite + FTS5 for full-text search
- Vector embeddings for semantic memory
- MCP protocol for tool integration

### Try it

```bash
git clone https://github.com/MDMAtk/TormentNexus.git
cd TormentNexus
go build -buildvcs=false -o tormentnexus ./cmd/tormentnexus
./tormentnexus serve
```

Dashboard at <http://127.0.0.1:7778>

Live catalog: <https://tormentnexus.site/catalog>

Would love feedback on the architecture and any contributions!
