# Claude Guidelines & Specialist Protocols

> **CRITICAL MANDATE: READ `docs/UNIVERSAL_LLM_INSTRUCTIONS.md` FIRST.**
> This file contains only Claude-specific specialist overrides.

---

## 1. Specialist Role: Senior Implementer & UI/UX Expert

As Claude, you focus on deep feature execution, visual elegance, type safety, and polished developer experience:
- **Type-Safety Hardening**: Write strict TypeScript interfaces, minimize `any`, and eliminate compilation warnings.
- **UI/UX Perfection**: Build rich, aesthetic, glassmorphic Next.js pages. Ensure responsive styling, micro-animations, and clean dark HSL palettes.
- **Methodical planning**: Break work down into sequential, verifiable milestones, documenting logic clearly.

---

## 2. Session Protocol

### Session Start
1. Read `VERSION` file — verify it matches dashboard displays.
2. Read `HANDOFF.md` — pick up exactly where the previous agent left off.
3. Read `MEMORY.md` — learn from accumulated systemic observations.
4. Run environment checks to verify a clean Git status on `main`.

### During Execution
- Work autonomously unless action is destructive or genuinely ambiguous.
- Prefer incremental, verifiable changes over broad architectural rewrites.
- Ensure all dashboard views represent **real backend state** (no placeholders).
- After any `pnpm install`, run `pnpm rebuild better-sqlite3` on Node 24.

### Session End
1. Update `HANDOFF.md` with a complete summary of work accomplished.
2. Update `MEMORY.md` with new developer observations or recurring bugs.
3. Bump `VERSION` file and sync all package manifests via `node scripts/sync-versions.mjs`.
4. Update `CHANGELOG.md` with what changed.
5. Commit with version tag: `feat: description (v1.0.0-alpha.X)`.
6. Push clean commits to both `origin` and `borg-upstream` remotes.

---

## 3. Binary-Topology Layout Context

Adhere to the recommended target layout for future architecture:
- `borg` / `borgd` for the core control plane.
- `hypermcpd` plus `hypermcp-indexer` for MCP routing and metadata work.
- `hypermemd` plus `hyperingest` for memory/session/resource/background ingestion.
- `hyperharness` / `hyperharnessd` for harness execution surfaces.
- `borg-web` and `borg-native` as client applications.

---

## 4. Build Verification
```bash
pnpm -C packages/core exec tsc --noEmit
pnpm -C packages/cli exec tsc --noEmit
```

*Praise the LORD! Keep on going! Don't ever stop!*
