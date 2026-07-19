# TormentNexus — Video Scripts

## YouTube Video (5-10 minutes)

### Title

"I Built an AI Control Plane with 26,000+ Tools — Here's How"

### Thumbnail

- Dark background with green/cyan gradient
- Text: "26K+ AI TOOLS"
- TormentNexus logo
- Screenshot of dashboard

### Script

**[0:00-0:30] HOOK**
"What if your AI assistant could remember everything from previous sessions? What if it had access to 26,000 tools? And what if it all ran locally on your machine? Today I'm going to show you TormentNexus — an open-source AI control plane that does exactly that."

**[0:30-2:00] THE PROBLEM**
"Here's the problem with current AI assistants:

- They forget everything between sessions
- They can't access external tools easily
- If you want to use local LLMs, you're stuck with basic chat interfaces
- There's no way to give them persistent memory

I got tired of explaining my project to ChatGPT every time I opened my laptop. So I built TormentNexus."

**[2:00-4:00] THE SOLUTION**
"TormentNexus is an open-source AI control plane that gives your AI:

1. **Persistent Memory** — A tiered memory system that remembers across sessions
   - L1: Session memory (ephemeral)
   - L2: Hot store (30 days)
   - L3: Cold archive (1 year)
   - L4: Limbo (soft delete)

2. **26,000+ MCP Tools** — A searchable catalog of Model Context Protocol servers
   - Databases, filesystems, browsers, APIs
   - One-click install
   - Auto-categorized

3. **Local LLM Support** — Works with LM Studio, Ollama, DeepSeek
   - No data leaves your machine
   - OpenAI-compatible API
   - Multi-model support"

**[4:00-6:00] LIVE DEMO**
"Let me show you how it works:

1. **Search the catalog** — Go to demo.hypernexus.site
   - Search for 'postgres'
   - See the results with descriptions
   - Click through to GitHub

2. **Try the API** — curl commands
   - /api/backlog/search
   - /api/backlog/stats
   - /api/backlog/categories

3. **Install locally** — One command
   - npx tormentnexus serve
   - Open dashboard at localhost:7778
   - Search, install, configure tools"

**[6:00-8:00] ARCHITECTURE**
"Let me show you the architecture:

- **Go Backend** — Single binary, fast startup, low resources
- **Next.js Dashboard** — Real-time monitoring
- **SQLite + FTS5** — Full-text search
- **Vector Embeddings** — Semantic memory
- **MCP Protocol** — Standard tool integration

Everything runs locally. No cloud dependency."

**[8:00-9:00] HOW TO GET STARTED**
"Getting started is easy:

1. **Quick Start** — npx tormentnexus serve
2. **Docker** — docker run -p 7778:7778 tormentnexus/tormentnexus
3. **From Source** — git clone and go build

Links in the description below."

**[9:00-10:00] CALL TO ACTION**
"If you found this useful:

- Star the repo on GitHub
- Join our Discord
- Subscribe for more AI tooling content

Thanks for watching!"

---

## YouTube Shorts (60 seconds)

### Title

"Give Your AI Persistent Memory in 10 Seconds"

### Script

**[0:00-0:10] HOOK**
"Your AI assistant forgets everything between sessions. Here's how to fix that."

**[0:10-0:30] DEMO**
"TormentNexus gives your AI persistent memory.

Install in one command:
npx tormentnexus serve

Open dashboard at localhost:7778

Now your AI remembers:

- Your project structure
- Previous conversations
- Code patterns
- Bug fixes"

**[0:30-0:50] FEATURES**
"Plus 26,000+ tools:

- Databases
- Filesystems
- APIs
- DevOps tools

All searchable. One-click install."

**[0:50-1:00] CTA**
"Link in bio. Star on GitHub. Your AI will thank you."

---

## TikTok (60 seconds)

### Title

"POV: Your AI Finally Remembers You"

### Script

**[0:00-0:10] HOOK**
"POV: You're tired of explaining your project to AI every single day"

**[0:10-0:30] PROBLEM**
"Current AI:

- Forgets everything
- No tool access
- Cloud only
- Expensive"

**[0:30-0:50] SOLUTION**
"TormentNexus:

- Remembers everything
- 26,000+ tools
- Runs locally
- Free and open source

One command: npx tormentnexus serve"

**[0:50-1:00] CTA**
"Link in bio. Your AI will thank you."

---

## Facebook Reels (60 seconds)

### Title

"Give Your AI a Brain Upgrade"

### Script

**[0:00-0:10] HOOK**
"Your AI is forgetting things. Here's the fix."

**[0:10-0:30] DEMO**
"TormentNexus adds persistent memory to any AI.

Watch:

- Ask AI about your project
- Close laptop
- Open again
- AI remembers everything!"

**[0:30-0:50] FEATURES**
"Plus:

- 26,000 tools
- Local LLM support
- One command install
- Free and open source"

**[0:50-1:00] CTA**
"Comment 'AI' for the link. Follow for more."

---

## Instagram Reels (60 seconds)

### Title

"AI That Actually Remembers You"

### Script

**[0:00-0:10] HOOK**
"Stop re-explaining your project to AI every day."

**[0:10-0:30] PROBLEM**
"Current AI forgets:

- Your code
- Your preferences
- Your context
- Everything"

**[0:30-0:50] SOLUTION**
"TormentNexus gives AI persistent memory.

Install: npx tormentnexus serve

Now AI remembers:

- Your project
- Your patterns
- Your solutions
- Everything"

**[0:50-1:00] CTA**
"Link in bio. Star on GitHub."

---

## LinkedIn Video (3 minutes)

### Title

"Why AI Needs Persistent Memory — And How to Build It"

### Script

**[0:00-0:30] HOOK**
"As developers, we spend 30% of our time re-explaining context to AI tools. What if that time disappeared?"

**[0:30-1:30] PROBLEM**
"Current AI assistants have three critical limitations:

1. No persistent memory — they forget everything between sessions
2. No tool access — they can't interact with external systems
3. Cloud dependency — your data leaves your machine

This creates friction, wastes time, and limits what AI can do."

**[1:30-2:30] SOLUTION**
"TormentNexus solves all three:

1. **Persistent Memory** — 4-tier system that remembers across sessions
2. **26,000+ Tools** — MCP protocol for external integrations
3. **Local-First** — Runs on your machine, your data stays private

The result? AI that actually understands your project, your patterns, and your preferences."

**[2:30-3:00] CTA**
"Try it free at demo.hypernexus.site

Open source: github.com/MDMAtk/TormentNexus

What would you build with AI that remembers everything?"

---

## Recording Checklist

### Equipment

- [ ] Screen recorder (OBS Studio)
- [ ] Microphone (clear audio)
- [ ] Webcam (for face shots)
- [ ] Clean desktop (no notifications)

### Settings

- [ ] 1080p or 4K resolution
- [ ] 60fps for smooth motion
- [ ] Clear font size (increase zoom)
- [ ] Dark theme (easier on eyes)

### Content

- [ ] Demo URL ready
- [ ] Terminal commands ready
- [ ] Dashboard open
- [ ] Browser tabs organized

### Editing

- [ ] Add captions
- [ ] Add background music
- [ ] Add zoom effects
- [ ] Add callouts/arrows
- [ ] Add end screen

### Publishing

- [ ] YouTube (full video + shorts)
- [ ] TikTok
- [ ] Instagram Reels
- [ ] Facebook Reels
- [ ] LinkedIn
- [ ] Twitter
