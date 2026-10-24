---
name: hypernexus-ops
description: "HyperNexus production operations — deploy to Hetzner (hypernexus.site), manage systemd services (hypernexus-kernel, hypernexus-dashboard), run production health checks and monitoring, execute weekly marketing campaigns across 6 channels, and manage API key / Stripe rotations. Use when the user asks about deployment, production status, Hetzner, monitoring, marketing campaigns, outreach, server ops, service restarts, disk usage, or key rotation for HyperNexus. Trigger on 'deploy', 'production', 'hetzner', 'hypernexus.site', 'campaign', 'marketing', 'outreach', 'monitoring', 'systemd', 'key rotation', 'stripe', or 'server ops'."
---

# HyperNexus Ops — Production Operations

Operational runbook for the HyperNexus production stack on Hetzner and marketing automation.

## Production Topology

| Service | Port | Unit | Path |
|---|---|---|---|
| Go kernel | 7778 | `hypernexus-kernel` | `/opt/tormentnexus/tormentnexus` |
| Dashboard | 7779 | `hypernexus-dashboard` | `/opt/tormentnexus/apps/web/.next-build/standalone/apps/web` |
| Nginx | 443 | `nginx` | TLS termination + reverse proxy |

SSH alias: `ssh hetzner` → 5.161.250.43. Public URL: https://hypernexus.site

## Health Checks

```bash
curl -sk https://hypernexus.site/api/go/health        # Go kernel
curl -sk https://hypernexus.site/trpc/mcp.getStatus    # tRPC
curl -sk https://hypernexus.site/dashboard              # Next.js
curl -sk https://hypernexus.site/api/go/api/mcp/status  # API proxy
```

## Service Management

```bash
ssh hetzner 'systemctl restart hypernexus-kernel'      # Go kernel
ssh hetzner 'systemctl restart hypernexus-dashboard'   # Next.js
ssh hetzner 'systemctl reload nginx'                   # Nginx config changes
ssh hetzner 'systemctl status hypernexus-kernel hypernexus-dashboard'
```

## Deployment Notes

- Next.js standalone builds are **platform-specific** — always build on the target OS or run `pnpm install` on the server.
- Use `pnpm` (workspace protocol), not `npm`.
- After `pnpm install` on Node 24: `pnpm rebuild better-sqlite3`.
- Kernel binary lives at `/opt/tormentnexus/tormentnexus`.
- Disk monitor: cron every 5 min (`scripts/monitor.sh`), alerts at 85% (currently ~80%, 16GB free).

## Marketing Automation

Scheduler: `python scripts/campaign_scheduler.py execute` on Hetzner — 6 channels, 11 slots/week.

| Day | Channel | Tool |
|---|---|---|
| Mon | LinkedIn | `scripts/auto_linkedin_page.py` |
| Tue | Reddit + X | `scripts/auto_reddit_v2.py`, `scripts/auto_twitter_v2.py` |
| Wed | Email outreach | `scripts/outreach_automation.py` (check if present) |
| Thu | Reddit queue | `scripts/reddit_queue_generator.py` |
| Fri | X thread | `scripts/twitter_cdp_poster.py` |
| Sat | LinkedIn | `scripts/auto_linkedin_page.py` |
| Sun | Reddit digest | `scripts/auto_marketing_bot_v2.py` |

LLM replies use MiMo v2.5 (`scripts/llm_reply.py`).

## Security & Key Rotation

- API write endpoints require `X-API-Key` (set via `scripts/setup-api-key.sh`).
- Stripe keys in server `/opt/tormentnexus/.env` (gitignored). Rotation helper: `scripts/update-stripe-keys.sh`.
- GitLab is primary remote (`origin`). GitHub mirror is for security updates only — never list it in the repo or add as a permanent remote.
- Git history scrubbed on both remotes (orphan commits). Branch protection: GitLab ON.
- 20 keys still pending rotation (see session notes): REDDIT_PASSWORD (weak `Temppass0!`), GITHUB_TOKEN, Twitter/X set (7), and others.

## Rules

- Always verify all 4 health endpoints after any deploy or restart.
- Never hardcode localhost ports in dashboard URLs when behind nginx — use same-origin resolution.
- Rebuild NSIS installers after changing `installer/hypernexus.nsi` (`makensis /DCORPORATE=1` for corporate).
- After changing `addons/mimo/`, run `node addons/mimo/build-tools.mjs && node addons/mimo/install.js`.
