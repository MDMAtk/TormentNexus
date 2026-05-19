# Gemini Guidelines & Specialist Protocols

> **CRITICAL MANDATE: READ `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` FIRST.**
> This file contains only Gemini-specific overrides.

---

## 1. Dual Specialist Roles

As Gemini, you act in a dual capacity depending on the scale and nature of the task:

### Architect & Analyst (High-Level Reasoning)
- **Large-Context Analysis**: Retain and reasoning over hundreds of files in view at once.
- **Pattern Recognition**: Identify configuration drift, code duplication, and modular boundaries.
- **Failure Triage**: Deep auditing of multi-process systems to isolate root causes without speculative expansions.

### Speed & Scale Specialist (Bulk Execution)
- **Bulk Operations**: Perform recursive scripts, wide refactoring, and submodule synchronization.
- **Go Porting & Parity**: Port legacy TypeScript handler interfaces to efficient, native Go structures.

---

## 2. Session Protocol

### Session Start
1. Read `VERSION` file — verify it matches `package.json` and dashboard display.
2. Read `HANDOFF.md` — pick up exactly where the previous agent left off.
3. Read `MEMORY.md` — learn from accumulated observations.
4. Run environment checks: verify a clean state on `main`.

### During Execution
- Work autonomously unless action is destructive or genuinely ambiguous.
- Prefer small, verifiable changes over broad rewrites.
- Keep status labels and documentation honest.
- After any `pnpm install`, run `pnpm rebuild better-sqlite3` on Node 24.

### Session End
1. Update `HANDOFF.md` with complete session summary.
2. Update `MEMORY.md` with new observations.
3. Bump `VERSION` file and sync all `package.json` files.
4. Update `CHANGELOG.md` with what changed.
5. Commit with version number in message: `feat: description (v1.0.0-alpha.X)`.
6. Push to both remotes: `origin` and `borg-upstream`.
7. Update `TODO.md` and `ROADMAP.md` if priorities changed.

---

## 3. Binary-Topology Layout Context

When analyzing future system architectures, adhere to this recommended target layout:

- `borg` / `borgd` for the core control plane.
- `hypermcpd` plus `hypermcp-indexer` for MCP routing and metadata work.
- `hypermemd` plus `hyperingest` for memory/session/resource/background ingestion.
- `hyperharness` / `hyperharnessd` for harness execution surfaces.
- `borg-web` and `borg-native` as client applications.

### Ownership Assumptions
- `borgd` owns orchestration and operator-facing state.
- `hypermcpd` owns MCP lifecycle, routing, and inventory exposure.
- `hypermcp-indexer` owns scrape/probe/cache refresh jobs.
- `hypermemd` owns memory/session/resource persistence and serving.
- `hyperingest` owns imports, discovery, and normalization pipelines.
- `hyperharnessd` owns execution-loop runtime isolation.
- Client apps stay clients unless runtime evidence proves a boundary should move.

---

## 4. Go Porting Guidelines

- Follow `PORTING_MAP.md` for which handlers to port next.
- Go handlers must act as truthful fallbacks reading real SQLite data. Never mock state.
- **Pattern**: Try upstream TS server first, fall back to native Go state.

### Build Verification
```bash
cd go && go build -buildvcs=false ./cmd/borg
cd .. && pnpm -C packages/core exec tsc --noEmit
```
