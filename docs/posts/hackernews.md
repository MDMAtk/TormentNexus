Title: Show HN: TormentNexus — Open-source AI control plane: persistent memory, 20K MCP tools, multi-agent swarm

Hey HN,

I've been building TormentNexus for months — a local-first, open-source (MIT) universal AI control plane. It's a Go + TypeScript monorepo that gives LLM agents persistent memory, MCP tool orchestration, and multi-agent coordination.

**The problem:** Every AI coding agent starts from scratch. No memory of past decisions. No shared context. Tools are scattered across dozens of MCP servers with no intelligent routing.

**What TormentNexus does:**

- Persistent multi-tier memory (L1 session scratchpad → L2 vector vault → L3 cold archive)
- Progressive semantic tool routing — searches 20,000+ MCP servers but only injects the relevant ones into context
- Auto-installs for 38+ AI clients (Claude, Gemini, Cursor, Windsurf, OpenHands, Aider, CodeWhale...)
- Multi-agent swarm with A2A protocol and consensus engine
- LLM waterfall: NVIDIA NIM → OpenRouter → Local Ollama

```bash
npx @tormentnexus/install  # One command, everything wired
```

**Tech:** Go sidecar (kernel), Next.js 16 dashboard, SQLite with vector search, Progressive disclosure MCP router, Ed25519 enterprise licensing.

The project runs in production on Hetzner with Stripe billing and per-tenant Docker isolation.

I've been iterating with a multi-AI-agent workflow (Gemini for speed, Claude for UI depth, GPT for architecture) coordinated through a shared AGENTS.md protocol. This project was built by AI agents using their own tooling — some inception-level recursion.

Would love technical feedback, especially on:

- The progressive MCP routing approach
- Memory tier architecture best practices
- Go sidecar design patterns

Repo: <https://github.com/MDMAtk/TormentNexus>
npm: @tormentnexus/install, @tormentnexus/cli, @tormentnexus/core, @tormentnexus/openhands
