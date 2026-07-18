# Reddit Post Draft — r/MachineLearning or r/LocalLLaMA

## Title

I built an open-source AI control plane with 26K+ MCP servers catalog, persistent memory, and local LLM support

## Body

Hey r/MachineLearning 👋

I've been working on **HyperNexus** — an open-source AI control plane that gives your AI agents persistent memory, tool access, and the ability to work with local LLMs.

### What it does

🧠 **Persistent Memory System**

- L1 (session) → L2 (hot vector store) → L3 (cold archive) → L4 (limbo vault)
- Your AI remembers context across sessions, projects, and even different models

🔧 **26,000+ MCP Servers Catalog**

- Searchable database of Model Context Protocol servers
- Auto-categorized: databases, filesystems, browsers, APIs, DevOps, etc.
- One-click install for MCP tools

🤖 **Local LLM Support**

- Works with LM Studio, Ollama, DeepSeek, and any OpenAI-compatible API
- No data leaves your machine unless you want it to

📊 **Unified Dashboard**

- Real-time monitoring of memory, tools, agents, and security
- System health, mesh network status, and more

### Tech Stack

- Go backend (fast, single binary)
- Next.js dashboard
- SQLite + FTS5 for search
- Vector embeddings for semantic memory

### Try it

```bash
# Quick install (Linux/Mac)
curl -fsSL https://raw.githubusercontent.com/MDMAtk/HyperNexus/main/scripts/install.sh | bash

# Or clone and build
git clone https://github.com/MDMAtk/HyperNexus.git
cd HyperNexus
go build -buildvcs=false -o hypernexus ./cmd/hypernexus
./hypernexus serve
```

Dashboard: <http://127.0.0.1:7778>

### Links

- **GitHub**: <https://github.com/MDMAtk/HyperNexus>
- **Live Catalog**: <https://hypernexus.site/catalog>
- **Blog**: <https://hypernexus.site/blog/>

Would love feedback on the architecture and any feature requests!

---

## Tags

# MachineLearning #AI #LLM #MCP #OpenSource #LocalAI
