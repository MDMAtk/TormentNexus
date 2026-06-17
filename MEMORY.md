# MEMORY.md — Multi-Agent Observations

## Session 2026-06-16 (Continued)

### Architecture Observations
- **Go sidecar bridges to TS control plane** via `handleTRPCBridgeBodyCall()` — most `/api/*` endpoints are proxies to tRPC procedures on port 4100
- **Session import** (`/api/session-export/import`) expects body `{"data":"<json-string>","merge":true,"dryRun":false}` — empty `{}` results in 0 imports
- **Import pipeline**: Go sidecar scans `~/.claude` and `~/.aider` directories, produces 49 valid candidates. TS control plane may ingest these asynchronously.
- **Swarm v7** runs with 5 workers, `--forever` mode, `--limit 200`, pipeline: GENERATE → REVIEW(2) → FIX
- FreeLLM proxy fallback chain: `proxy-free-llm` → `proxy-free-llm-fallback` → direct providers

### Gotchas & Tool Quirks
- **`.gitignore` has `data/`** — but `data/assimilation_state.db` and `data/prompt_library.db` are tracked (added before ignore rule). Use `git add -u` to stage tracked DB changes.
- **`go-sidecar.pid`** is untracked (not committed) — it's a runtime file.
- **`*.db-shm` and `*.db-wal`** are ignored by `.gitignore` — SQLite WAL files won't be committed.
- **`git add "go/internal/tools/"`** stages all new/deleted tool files from swarm output.
- **Sync script** `node scripts/sync-versions.mjs` updates all 35 package.json files + `go.mod` + Go buildinfo.
- **Some NVIDIA models** (deepseek-v4-pro, deepseek-v4-flash) return short/empty responses frequently — swarm retries 5x with exponential backoff.

### Failure Lessons
- **bobbybookmarks.com** DNS resolution fails from this environment — cannot use for links backlog sync. Consider Smithery.ai or Glama.ai as alternatives.
- **`--repair` flag** in swarm_v7 causes premature exit — use `--forever` without `--repair` for stability.
- **`tormentnexus-upstream` remote** is not configured — only `origin` and `origin-backup` exist.
- **`data/` directory** is gitignored but DB files within it may still be tracked — always use `git add -u` for existing tracked files, not `git add data/`.

### Preferences
- Always stage and commit `.db` files per user instruction "stage and track db always"
- Use `--forever` mode for swarm to avoid premature shutdown
- Tag commits with version: `v1.0.0-alpha.X`
- Push to both `origin` and upstream when available
