# Memory Maintenance System

## Overview

HyperNexus uses a multi-tier memory architecture with automatic lifecycle management. Memories flow through tiers based on heat scores and access patterns.

```
L2 Vault (hot/warm) ──decay──▶ L3 Cold Archive (cool/cold) ──burial──▶ L4 Limbo (dead)
     │                                                                       │
     ◀─────────────── promotion ◀────────────────────────────────────────────┘
```

## Memory Tiers

| Tier | Name | Purpose | Retention |
|------|------|---------|-----------|
| L1 | Working Memory | Ephemeral scratchpad (current session) | Session-scoped |
| L2 | Long-Term Vault | Active memories with heat scores | Indefinite (while heat > 10) |
| L3 | Cold Archive | Compressed long-term storage | Indefinite |
| L4 | Limbo | Discarded/lost/orphaned memories | 90 days |
| Project | .memdb | Portable git-tracked per-project file | Indefinite |

## Heat Score System

Every memory in L2 has a **heat score** (0-100+) that determines its lifecycle:

| Heat Range | Status | Action |
|-----------|--------|--------|
| 80-100 | Hot | Frequently accessed, stays in L2 |
| 50-79 | Warm | Normal active memory |
| 20-49 | Cool | Decaying, at risk of archiving |
| 10-19 | Cold | Candidate for L3 cold archive |
| < 10 | Archive | Auto-demoted to L3 |

Heat scores decay via the **Ebbinghaus forgetting curve**:

```
heat = heat × e^(-0.0288 × hours_since_last_access)
```

This means:

- After 24 hours: heat × **0.50** (halved)
- After 48 hours: heat × **0.25** (quartered)
- After 72 hours: heat × **0.12** (reduced 88%)

## Maintenance Cycle

The full maintenance cycle runs **automatically every 4 hours** (or on demand via API).

### 1. ForgettingCurveDecay

Applies the Ebbinghaus decay formula to all non-archive memories. Memories with heat < 10 are moved to L3 Cold Archive.

### 2. ConsolidateMemories

Finds pairs of memories with similar content using Jaccard similarity (>0.8 threshold) and merges them. Reduces redundancy.

### 3. ApplyDecay

Secondary decay pass. Also promotes working memories (heat > 80) to long-term, and demotes long-term (heat < 20) to archive status.

### 4. BuryOrphanedMemories

Finds memories whose session no longer exists and whose heat < 15, moves them to L4 Limbo with reason "lost".

### 5. DreamCycle

Auto-reviews due spaced-repetition memories. Uses SM-2 algorithm to calculate next review interval. Successfully recalled memories get a +2 heat boost.

## Triggering Maintenance

**Automatic**: Runs 2 minutes after server start, then every 4 hours.

**Manual**:

```bash
curl -X POST http://127.0.0.1:7778/api/memory/maintenance
```

Response:

```json
{
  "success": true,
  "data": [
    "forgetting-curve-decay: complete",
    "consolidation: complete",
    "apply-decay: complete",
    "orphan-burial: complete",
    "dream-cycle: complete"
  ],
  "vault": 17444,
  "cold": 0
}
```

## Monitoring

- **Dashboard**: `/dashboard/memory-analytics` — tier stats, heat distribution, lifecycle pipeline
- **API**: `GET /api/memory/cold-archive/count` — L3 count
- **API**: `POST /api/memory/maintenance` — trigger cycle, returns vault + cold counts
- **Pi extension**: `/tn-status` slash command or `Ctrl+Shift+P` shortcut

## Spaced Repetition (Dream Cycle)

The dream cycle uses the SM-2 algorithm (same as SuperMemo/Anki):

| Quality | Meaning | Interval Change |
|---------|---------|----------------|
| 5 | Perfect recall | Interval × EF |
| 4 | Correct with hesitation | Interval × EF |
| 3 | Correct with difficulty | Reset to 1 day, EF -0.2 |
| 2 | Incorrect, close | Reset to 1 day, EF -0.2 |
| 1 | Complete blackout | Reset to 1 day, EF -0.2 |

- **EF (Easiness Factor)**: Starts at 2.5, minimum 1.3
- **First review**: 1 day
- **Successful review**: Increases by EF multiplier

## P2P Mesh (Cross-Machine Memory Sync)

### How It Works

The gossip protocol uses a push-pull anti-entropy model over UDP:

```
Node A ──sync request──▶ Node B
Node A ◀──sync response── Node B
Node A ◀──state update──── Node B (if B has newer data)
```

### Getting Peers

1. Start another TN instance on the same network:

   ```bash
   bin/hypernexus.exe serve --port 7778
   ```

2. The UDP discovery service automatically broadcasts presence every 30 seconds.

3. Peers appear within ~30-60 seconds under `/dashboard/mesh`.

### Mesh API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /api/mesh/status` | Node ID + peer count |
| `GET /api/mesh/peers` | List of peer node IDs |
| `GET /api/mesh/capabilities` | All peer capabilities |
| `GET /api/mesh/query-capabilities?nodeId=` | Single peer capabilities |
| `GET /api/mesh/find-peer?capability=` | Find peer with capability |
| `POST /api/mesh/broadcast` | Broadcast message to all peers |

### Current Status

- Node: `hypernexus-go-4c566d76db82`
- Peers: 0 (no other TN instances running)
- Gossip port: 7778 (UDP)
- Protocol: Ed25519-signed messages with vector clocks for conflict resolution

## Per-Project .memdb System

### Overview

Every HyperNexus project can have a `.memdb` file in its root directory — a portable, git-tracked SQLite database containing project-specific memories. These files survive clones, branches, and merges.

```
project-root/
  .memdb              ← git-tracked, portable memories for this project
  src/...

~/.hypernexus/
  memory.db           ← global unified index (aggregates all .memdb files)
```

### How Memories Flow

```
tn_memory_store(content="fix: build bug", project="hypernexus", tags=["pattern:build"])
  ├── Global L2 vault (memory.db) — for unified search
  └── Project tag added ──→ triggers workspace .memdb scan
                              └── hypernexus/.memdb created/updated
                                  └── imported into global index on next scan
```

### Server Scanning

- **Startup**: 10s after boot, scans `workspace/` for all `.memdb` files and imports into global L2
- **Hourly**: Rescans to pick up new projects from `git pull`, `git clone`, or manual creation
- **Manual**: `POST /api/memory/project/sync` triggers on-demand
- **Retroactive split**: `POST /api/memory/project/split` reads existing global memories with `project:` tags and writes them into per-project `.memdb` files

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/memory/project/sync` | POST | Scan workspace, import all `.memdb` files into global index
| `/api/memory/project/split` | POST | Split global memories by project tag into per-project `.memdb` files

### Pi Extension Usage

The `tn_memory_store` tool accepts an optional `project` parameter:

```
tn_memory_store(
  content="Architecture decision: use FTS5 for full-text search",
  project="hypernexus",
  tags=["pattern:search", "decision:architecture"],
  category="decision"
)
```

This stores the memory in both the global L2 vault AND triggers a workspace scan that writes it to `hypernexus/.memdb`. The `.memdb` file can be committed to git and will be auto-imported when the project is cloned on another machine.

## Configuration

- Config directory: `~/.hypernexus` (auto-migrated from `~/.hypernexus-go`)
- Override: `HYPERNEXUS_CONFIG_DIR` env var
- Maintenance interval: 4 hours (hardcoded)
- Initial delay: 2 minutes
