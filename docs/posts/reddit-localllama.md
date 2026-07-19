Title: Show r/LocalLLaMA: TormentNexus — Universal AI Control Plane with persistent memory and 20,000+ MCP tools

I built an open-source AI control plane that gives LLM agents persistent memory, tool orchestration, and multi-agent coordination. MIT licensed, one-command install.

**What it does:**

- **Persistent memory:** Multi-tier L1/L2/L3 storage. Every session searches past decisions. No more "I don't remember that."
- **20,000+ MCP tools:** Largest local MCP registry with progressive semantic routing — only relevant tools enter the context window. NASA, arXiv, CoinGecko, PokéAPI, thousands more.
- **38 AI client integrations:** Auto-detects Claude, Gemini, Cursor, Windsurf, OpenHands, Aider, CodeWhale, and more. Wires MCP, skills, and commands automatically.
- **Multi-agent swarm:** A2A protocol with Planner→Checker→Implementer→Critic role rotation
- **LLM waterfall:** NVIDIA NIM → OpenRouter → Local Ollama cascading failover
- **Self-healing:** Diagnose→Fix→Verify→Retry loop. Errors auto-remediated.

**Quick start:**

```bash
npx @tormentnexus/install
npm install -g @tormentnexus/cli
tn status
```

GitHub: <https://github.com/MDMAtk/TormentNexus>
npm: <https://www.npmjs.com/search?q=%40tormentnexus>

Built with Go + TypeScript + Next.js. Production-deployed on Hetzner with Stripe billing and Docker tenant isolation.

Would love feedback from this community — you're the exact people this is built for.
