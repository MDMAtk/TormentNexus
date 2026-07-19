# TormentNexus Registry & Catalog Documentation

> Last updated: 2026-07-13 | Version: 1.0.0-b2 | Total entries: 37,289

## Overview

TormentNexus aggregates MCP tools, skills, prompts, agent frameworks, and AI/dev tools from 45+ sources into a unified catalog. This document maps every registry, database, scraper script, and content source in the project.

---

## 1. Content Sources by Category

### 1.1 MCP Servers (19,464 entries)

| Source | Count | Type | URL/Location |
|---|---|---|---|
| awesome-mcp-fork | 10,289 | GitHub forks | `github.com/punkpeye/awesome-mcp-servers` forks |
| awesome-mcp-servers | 3,042 | GitHub | `github.com/punkpeye/awesome-mcp-servers` |
| npm | 1,391 | npm registry | `registry.npmjs.org` |
| multi-source-scrape | 1,260 | Mixed | Combined scraper run |
| github-search | 1,156 | GitHub API | `api.github.com/search/repositories` |
| npm-tool | 772 | npm registry | MCP-adjacent packages |
| github-mcp-name | 505 | GitHub API | Repos with "mcp" in name |
| github-specific | 270 | GitHub API | Category-specific searches |
| npm-expanded | 263 | npm registry | Expanded keyword searches |
| github-topic | 145 | GitHub API | `topic:mcp` repos |
| pypi | 101 | PyPI | `pypi.org` MCP packages |
| crates-io | 100 | crates.io | Rust MCP packages |
| github-topic-mcp | 54 | GitHub API | `topic:mcp` repos |
| github-mcp-desc | 44 | GitHub API | "mcp-server" in description |
| github-topic-mcp-proto | 32 | GitHub API | `topic:model-context-protocol` |
| glama-html | 24 | HTML scrape | `glama.ai/mcp/servers` |
| awesome-fork-v2 | 9 | GitHub forks | Additional awesome list forks |
| glama-mock | 3 | Hardcoded | Fallback presets |
| npm-mcp-sdk | 2 | npm registry | `@modelcontextprotocol` scope |
| mcp-get | 1 | GitHub | `michaellatman/mcp-get` |
| docker-hub | 1 | Docker Hub | MCP Docker images |

### 1.2 MCP by Language (249 entries)

| Language | Count | Source |
|---|---|---|
| C++ | 81 | `github-lang-c++` |
| Ruby | 76 | `github-lang-ruby` |
| Swift | 60 | `github-lang-swift` |
| Java | 21 | `github-lang-java` |
| Rust | 10 | `github-lang-rust` |
| TypeScript | 1 | `github-lang-typescript` |

### 1.3 MCP by Category (367 entries)

| Category | Count | Category | Count |
|---|---|---|---|
| Payment | 40 | CMS | 37 |
| Stripe | 36 | Social | 36 |
| E-commerce | 36 | CRM | 35 |
| Calendar | 19 | Logging | 17 |
| Email | 17 | Filesystem | 16 |
| Auth | 14 | Slack | 13 |
| Notion | 12 | Monitoring | 10 |
| AWS | 8 | Analytics | 8 |
| Discord | 6 | GitHub | 3 |
| Database | 2 | Browser | 2 |

### 1.4 AI/Dev Tools (5,582 entries)

| Source | Count | Type |
|---|---|---|
| awesome-general | 3,363 | General awesome lists |
| npm-ai-sdk | 784 | AI SDK packages |
| npm-tool-v2 | 693 | Tool packages |
| github-ai-tool | 453 | GitHub AI tool repos |
| npm-adjacent | 272 | MCP-adjacent npm packages |
| awesome-domain | 10 | Domain-specific awesome lists |
| awesome-ai-v2 | 7 | AI awesome lists |

### 1.5 Prompts/Templates (295 entries)

| Source | Count | Type |
|---|---|---|
| prompt-repo | 92 | Prompt repository READMEs |
| github-prompt-search | 87 | GitHub prompt searches |
| fabric-patterns | 45 | `danielmiessler/Fabric` patterns |
| awesome-ai-system-prompts | 21 | `dontriskit/awesome-ai-system-prompts` |
| chatgpt-system-prompts | 12 | `LouisShark/chatgpt_system_prompt` |
| system-prompts-leaks | 9 | `asgeirtj/system_prompts_leaks` |
| big-prompt-library | 9 | `0xeb/TheBigPromptLibrary` |
| meigen-ai-design-mcp | 7 | `jau123/MeiGen-AI-Design-MCP` |
| system-prompts-ai-tools | 5 | `x1xhlol/system-prompts-and-models-of-ai-tools` |
| claude-code-prompts | 5 | `Piebald-AI/claude-code-system-prompts` |
| claude-code-prompts-v2 | 2 | `repowise-dev/claude-code-prompts` |
| leaked-system-prompts | 1 | `jujumilk3/leaked-system-prompts` |

### 1.6 Agent Frameworks (223 entries)

| Source | Count | Type |
|---|---|---|
| awesome-ai-agent | 183 | Awesome AI agent lists |
| agent-framework | 40 | Agent framework repos |

### 1.7 Skills (5,441 entries)

| Location | Count |
|---|---|
| `/opt/marketing_agent/borg/` | 5,438 |
| `/opt/tormentnexus/.antigravity/` | 1 |
| `/opt/tormentnexus/.mavis/` | 1 |
| `/opt/tormentnexus/.codewhale/` | 1 |

### 1.8 Go Native Handlers (5,668 entries)

| File | Count | Description |
|---|---|---|
| `go/internal/mcpimpl/registry.go` | 5,668 | MCP tool handler definitions |

---

## 2. Database Inventory

### 2.1 Kernel Databases (`/root/.tormentnexus/`)

| Database | Tables | Purpose |
|---|---|---|
| `memory.db` | 56 | L2 vault, L3 cold archive, L4 limbo, GraphRAG, FTS5 |
| `catalog.db` | 4 | Published skills, registry tools, sync state, tool index |
| `accounts.db` | 3 | Tenant accounts and billing |
| `l3_cold_archive.db` | 14 | Cold storage for old memories |

### 2.2 Catalog Databases (dual location)

| Database | Tables | Count |
|---|---|---|
| `/opt/tormentnexus/catalog.db` | links_backlog | 26,180 |
| `/opt/tormentnexus/catalog.db` | published_mcp_servers | 8 |
| `/root/.tormentnexus/catalog.db` | published_skills | 5,441 |

### 2.3 Docker Tenant Databases

- Each tenant has its own `tormentnexus.db` in the container volume
- Located at `/var/lib/hypernexus/tenants/{tenant}/`

---

## 3. Scraper Scripts

### 3.1 `scripts/scrape-mega.py`

- **Purpose:** All-in-one mega scraper
- **Sources:** awesome-mcp-servers, GitHub (20 queries), npm (10 queries), PyPI, curated lists, Glama HTML, Smithery, Docker Hub, crates.io, MCP Hub
- **Run:** `python3 scripts/scrape-mega.py`
- **Last run:** 2026-07-13

### 3.2 `scripts/scrape-more.py`

- **Purpose:** Additional sources (language, category, SDK, forks)
- **Sources:** GitHub by language (7), GitHub by category (24), npm AI SDKs, awesome forks, AI awesome lists
- **Run:** `python3 scripts/scrape-more.py`

### 3.3 `scripts/scrape-all-registries.py`

- **Purpose:** Comprehensive multi-source scraper
- **Sources:** awesome-mcp-servers, GitHub topics, npm, PyPI, curated lists
- **Run:** `python3 scripts/scrape-all-registries.py`

### 3.4 `scripts/scrape-mcp-servers.py`

- **Purpose:** Basic scraper
- **Sources:** awesome-mcp-servers, GitHub topics
- **Run:** `python3 scripts/scrape-mcp-servers.py`

---

## 4. External Registry Sources

### 4.1 Active Sources (Working)

| Source | URL | Status | Last Scrape |
|---|---|---|---|
| awesome-mcp-servers | `github.com/punkpeye/awesome-mcp-servers` | ✅ Working | 2026-07-13 |
| GitHub Topics | `api.github.com/search/repositories` | ✅ Working | 2026-07-13 |
| GitHub Language | `api.github.com/search/repositories` | ✅ Working | 2026-07-13 |
| GitHub Category | `api.github.com/search/repositories` | ✅ Working | 2026-07-13 |
| npm | `registry.npmjs.org` | ✅ Working | 2026-07-13 |
| PyPI | `pypi.org` | ✅ Working | 2026-07-13 |
| crates.io | `crates.io/api/v1/crates` | ✅ Working | 2026-07-13 |
| Docker Hub | `hub.docker.com/v2/search/repositories` | ✅ Working | 2026-07-13 |
| Fabric | `github.com/danielmiessler/Fabric` | ✅ Working | 2026-07-13 |
| Prompt repos | Various GitHub repos | ✅ Working | 2026-07-13 |

### 4.2 Broken Sources (Need Fix)

| Source | URL | Issue |
|---|---|---|
| Glama.ai API | `glama.ai/api/v1/mcp/servers` | Returns HTML instead of JSON |
| Smithery | `smithery.ai/api/v1/servers` | Empty response |
| mcp.run | `mcp.run/api/catalog` | 301 redirect |

### 4.3 Sources Not Yet Integrated

| Source | URL | Notes |
|---|---|---|
| Toolhouse | `toolhouse.ai` | Agent tool marketplace |
| Composio | `composio.dev` | Integration platform |
| Mintlify | `mintlify.com` | API docs platform |
| LangChain tools | `langchain.com` | Tool registry |
| LlamaIndex tools | `llamaindex.ai` | Tool registry |

---

## 5. How to Add a New Registry Source

### Option A: Python Scraper (Quick)

1. Add a new function to `scripts/scrape-mega.py`
2. Fetch from the API/URL
3. Parse and normalize entries
4. Insert into `catalog.db` → `links_backlog`
5. Run the script

### Option B: Go Adapter (Production)

1. Create a new adapter in `go/internal/mcp/catalog_ingest.go`
2. Implement `CatalogSourceAdapter` interface (`Name()` + `Ingest()`)
3. Add to `adapters` slice in `IngestPublishedCatalog()`
4. Rebuild kernel

---

## 6. Quick Reference Commands

```bash
# Check catalog counts
ssh hetzner 'sqlite3 /opt/tormentnexus/catalog.db "SELECT source, count(*) FROM links_backlog GROUP BY source ORDER BY count(*) DESC;"'

# Run mega scraper
python3 scripts/scrape-mega.py

# Check active tools
curl -s http://127.0.0.1:8090/api/runtime/status | jq '.data.cli.toolCount'

# List skills
find .tormentnexus/skills/ -name "SKILL.md" | wc -l

# Check Stripe config
ssh hetzner 'grep STRIPE /opt/tormentnexus/.env'

# Full catalog status
ssh hetzner 'sqlite3 /opt/tormentnexus/catalog.db "SELECT count(*) FROM links_backlog;" && sqlite3 /root/.tormentnexus/catalog.db "SELECT count(*) FROM published_skills;"'
```
