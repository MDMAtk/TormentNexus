# YouTube Video Script: "I Built an AI with 26,000 Tools"

## Video Details

- **Title**: I Built an AI with 26,000 Tools (Open Source)
- **Duration**: 8-10 minutes
- **Style**: Faceless screen recording with voiceover
- **Thumbnail**: Dark background, green text "26K+ AI TOOLS", screenshot of dashboard

---

## Script

### [0:00-0:30] HOOK

**[Screen: Quick montage of dashboard, tool catalog, memory system]**

"What if your AI assistant could remember everything from previous sessions? What if it had access to 26,000 tools? And what if it all ran locally on your machine, completely free?

Today I'm going to show you TormentNexus — an open-source AI control plane that does exactly that."

---

### [0:30-2:00] THE PROBLEM

**[Screen: Show frustration with ChatGPT forgetting context]**

"Here's the problem with current AI assistants:

1. They forget everything between sessions
2. They can't access external tools easily
3. If you want to use local LLMs, you're stuck with basic chat interfaces
4. There's no way to give them persistent memory

I got tired of explaining my project to ChatGPT every time I opened my laptop. So I built TormentNexus."

---

### [2:00-4:00] THE SOLUTION

**[Screen: Show TormentNexus architecture diagram]**

"TormentNexus is an open-source AI control plane that gives your AI:

**1. Persistent Memory** — A tiered memory system that remembers across sessions

- L1: Session memory (ephemeral)
- L2: Hot store (30 days)
- L3: Cold archive (1 year)
- L4: Limbo (soft delete)

**2. 26,000+ MCP Tools** — A searchable catalog of Model Context Protocol servers

- Databases, filesystems, browsers, APIs
- One-click install
- Auto-categorized

**3. Local LLM Support** — Works with LM Studio, Ollama, DeepSeek

- No data leaves your machine
- OpenAI-compatible API
- Multi-model support"

---

### [4:00-6:00] LIVE DEMO

**[Screen: Live demo of the dashboard]**

"Let me show you how it works:

**1. Search the catalog** — Go to demo.hypernexus.site

- Search for 'postgres'
- See the results with descriptions
- Click through to GitHub

**2. Try the API** — curl commands

- /api/backlog/search
- /api/backlog/stats
- /api/backlog/categories

**3. Install locally** — One command

- npx @tormentnexus/core serve
- Open dashboard at localhost:7778
- Search, install, configure tools"

---

### [6:00-8:00] ARCHITECTURE

**[Screen: Show architecture diagram]**

"Let me show you the architecture:

- **Go Backend** — Single binary, fast startup, low resources
- **Next.js Dashboard** — Real-time monitoring
- **SQLite + FTS5** — Full-text search
- **Vector Embeddings** — Semantic memory
- **MCP Protocol** — Standard tool integration

Everything runs locally. No cloud dependency."

---

### [8:00-9:00] HOW TO GET STARTED

**[Screen: Show installation steps]**

"Getting started is easy:

**Option 1: npm**

```bash
npx @tormentnexus/core serve
```

**Option 2: Docker**

```bash
docker run -p 7778:7778 ghcr.io/mdmatk/tormentnexus:latest
```

**Option 3: Download**

- Windows: tormentnexus-setup.exe
- macOS/Linux: Install script

Links in the description below."

---

### [9:00-10:00] CALL TO ACTION

**[Screen: Show GitHub page]**

"If you found this useful:

- ⭐ Star the repo on GitHub
- 💬 Join our Discord
- 🔔 Subscribe for more AI tooling content

Thanks for watching!"

---

## Thumbnail Design

```
Background: Dark (#09090b)
Main text: "26K+ AI TOOLS" (green gradient)
Subtext: "Open Source • Local • Free"
Image: Screenshot of dashboard
Logo: TormentNexus logo
```

## Tags

```
ai, mcp, tools, local llm, ollama, lm studio, persistent memory, 
ai agent, open source, developer tools, coding, programming
```

## Description

```
I built an open-source AI control plane with 26,000+ MCP tools and persistent memory.

🔗 Links:
- GitHub: https://github.com/MDMAtk/TormentNexus
- Website: https://tormentnexus.site
- Demo: https://demo.hypernexus.site
- Discord: https://discord.gg/Hj9P3GbVxR

⚡ Quick Start:
npx @tormentnexus/core serve

📦 Features:
- 26,000+ MCP tools
- Persistent memory (4-tier)
- Local-first architecture
- Works with Ollama, LM Studio, DeepSeek
- Single binary, zero dependencies

#ai #mcp #opensource #developer #tools
```
