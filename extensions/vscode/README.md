# HyperNexus VS Code Extension

AI Control Plane with Persistent Memory & 26,000+ MCP Tools

## Features

- 🧠 **Memory Explorer** — View and search your persistent memory (L1-L4)
- 🔧 **MCP Tool Search** — Search 26,000+ tools from the sidebar
- 📊 **Status Dashboard** — Monitor connection and server health
- ⚡ **Quick Commands** — Add memory, search tools, open dashboard

## Installation

1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "HyperNexus"
4. Click Install

## Usage

### Sidebar

- Click the HyperNexus icon in the activity bar
- View memory tiers, tool stats, and connection status

### Commands

- `HyperNexus: Connect to Server` — Connect to local server
- `HyperNexus: Search MCP Tools` — Search the tool catalog
- `HyperNexus: Add Memory` — Save a memory entry
- `HyperNexus: Search Memory` — Search your memories
- `HyperNexus: Open Dashboard` — Open web dashboard

### Configuration

```json
{
  "hypernexus.serverUrl": "http://localhost:7778",
  "hypernexus.autoConnect": true
}
```

## Requirements

- HyperNexus server running locally (`hypernexus serve`)
- Node.js 18+

## Links

- [GitHub](https://gitlab.com/robertpelloni/HyperNexus)
- [Website](https://hypernexus.site)
- [Discord](https://discord.gg/Hj9P3GbVxR)

## License

MIT
