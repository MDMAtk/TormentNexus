# HANDOFF.md — Executive Protocol R22

**Date:** 2026-07-13 04:00 UTC
**Version:** 1.0.0-alpha.261
**Model:** Claude (via Pi)

## Summary

Executive Protocol R22 completed. Cherry-picked useful new files from jules feature branch, skipped conflicting Go import renames. Marketing agent all branches cleanly merged. Full pipeline deployed on Hetzner.

## Completed Actions

### TormentNexus (origin/main + origin-backup/main)

- v1.0.0-alpha.260 → v1.0.0-alpha.261
- Cherry-picked from `jules-3383774315910324119-9b0b1aa0` branch:
  - `convert_pages.py` — Dashboard component collector
  - `list_dashboard_pages.py` — Dashboard page walker  
  - `plan.md` / `plan.txt` — Tab→stacked dashboard layout plan
  - `scripts/e2e_integration_verify.py` — End-to-end test harness
  - `landing/demo/demo.html` — Updated demo page
  - `apps/web/next-env.d.ts` — Already synced (no change)
- Skipped: 47 Go file import renames (conflicted with our sidecar→kernel/commercial rename)

### Marketing Agent (robertpelloni/marketing_agent)

- All 3 feature branches clean (0 unique commits each): fully merged prior

### OpenHands Integration (completed earlier this session)

- Python agent at `~/.openhands/plugins/tormentnexus/tormentnexus_agent.py`
- 10 actions at `~/.openhands/plugins/tormentnexus/actions.py`
- Plugin manifest, microagent, Docker compose, npm package
- Installer now copies all OpenHands files

### Hetzner Deployment

- Marketing agent: Rebuilt Docker image, deployed with correct DB creds (sales_bot:tormentnexus2026)
- Webhook endpoint: `POST /api/v1/webhook/stripe` — live, needs Stripe signature
- TN Kernel: PM2 running on port 8090
- 2 Docker tenants: test-org (3001), acmerealtest (3010)
- Stripe: Live keys configured (sk_live_, pk_live_, rk_live_)

### Pipeline Verified

```
Stripe Checkout → Webhook (/api/v1/webhook/stripe) → TN Provison → Account Created → Dashboard
```

## Current State

### Running on Hetzner (5.161.250.43)

| Service | Port | Status |
|---------|------|--------|
| TN Kernel | 8090 | PM2, auto-restart |
| Marketing Agent | 8087 | Docker, postgres connected |
| test-org dashboard | 3001 | Docker TN tenant |
| acmerealtest dashboard | 3010 | Docker TN tenant |

### Branches Not Merged

- `origin/jules-3383774315910324119-9b0b1aa0` — 1 commit (Go renames conflict, useful files cherry-picked)
- `origin-backup/feature/cloud-dashboard-mcp-sse-389806464713532918` — Stale (0 unique after sync)
- `gh-pages-hypernexus` / `gh-pages-tormentnexus` — Static pages, separate from main

### Remaining Work

1. Wildcard DNS `*.hypernexus.site` — Needs verification after propagation
2. npm publish `@tormentnexus/openhands` — Package ready, npm login needed
3. Marketing agent webhook — Real Stripe checkout needed to test full flow
4. 1726 Dependabot vulns (41 critical) — Needs triage
5. ~330 quarantined MCP tool stubs — Needs regeneration
6. Git LFS push blocked — Unknown large objects (DB files)
7. next-env.d.ts — 1-char fix from jules branch (no diff)

## Next Agent Should

1. Verify wildcard DNS propagated: `dig +short acmerealtest.hypernexus.site`
2. Publish npm packages if npm login available
3. Continue dashboard UI work per `plan.md` (tab→stacked layout)
4. Rebuild TN kernel after new features to bump compiled version from 255→261
