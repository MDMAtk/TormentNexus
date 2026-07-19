# TormentNexus VS Code Extension

AI Control Plane with Persistent Memory & 26,000+ MCP Tools

## Features

- 🧠 **Memory Explorer** — View and search your persistent memory (L1-L4)
- 🔧 **MCP Tool Search** — Search 26,000+ tools from the sidebar
- 📊 **Status Dashboard** — Monitor connection and server health
- ⚡ **Quick Commands** — Add memory, search tools, open dashboard

## Installation

1. Open VS Code
2. Go to Extensions (Ctrl+Shift+X)
3. Search for "TormentNexus"
4. Click Install

## Usage

### Sidebar

- Click the TormentNexus icon in the activity bar
- View memory tiers, tool stats, and connection status

### Commands

- `TormentNexus: Connect to Server` — Connect to local server
- `TormentNexus: Search MCP Tools` — Search the tool catalog
- `TormentNexus: Add Memory` — Save a memory entry
- `TormentNexus: Search Memory` — Search your memories
- `TormentNexus: Open Dashboard` — Open web dashboard

### Configuration

```json
{
  "tormentnexus.serverUrl": "http://localhost:7778",
  "tormentnexus.autoConnect": true
}
```

## Requirements

- TormentNexus server running locally (`tormentnexus serve`)
- Node.js 18+

## Links

- [GitHub](https://github.com/MDMAtk/TormentNexus)
- [Website](https://tormentnexus.site)
- [Discord](https://discord.gg/Hj9P3GbVxR)

## License

MIT
