# Pricing Strategy — HyperNexus

> **Last Updated:** 2026-07-18
> **Status:** Draft — Ready for implementation

---

## 💰 Pricing Models

### Option A: Local License (One-Time Purchase)

| | Details |
|---|---|
| **Price** | $50 |
| **Duration** | Lifetime for current major version |
| **Upgrades** | $25/year for new major versions |
| **Storage** | Local SQLite on user's machine |
| **Backups** | User's responsibility |
| **Support** | Community (GitHub/GitLab issues) |
| **Updates** | Security patches free forever |

**What You Get:**

- Full HyperNexus agent + dashboard
- L1/L2 memory system (local)
- 20,000+ MCP tool routing
- Multi-agent coordination
- Native tools (filesystem, web, etc.)
- All future 1.x updates
- Community support

**Best For:**

- Privacy-focused developers
- Teams with own infrastructure
- Offline/air-gapped environments
- Budget-conscious users
- Hobbyists and learners

---

### Option B: Cloud Subscription (Monthly)

| | Details |
|---|---|
| **Price** | $5/month |
| **Billing** | Monthly, auto-renew |
| **Cancellation** | Immediate loss of access |
| **Storage** | Cloud-managed, encrypted |
| **Backups** | Automatic daily backups |
| **Support** | Email support |
| **Updates** | Always latest version |

**What You Get:**

- Everything in Local License, PLUS:
- Cloud-synced memory across devices
- Sandboxed MCP container execution
- Automatic backups (30-day retention)
- Usage dashboard and analytics
- Priority email support
- Team sharing (coming soon)

**Best For:**

- Cloud-first developers
- Multi-device users
- Teams wanting managed infrastructure
- Users who don't want to manage backups
- Convenience-focused workflows

---

### Option C: Enterprise (Custom Pricing)

| Tier | Users | Price | Includes |
|------|-------|-------|----------|
| **Team** | 5-20 | $500/month | SSO, shared memory, audit logs |
| **Business** | 20-100 | $2,000/month | + RBAC, custom domains, SLA |
| **Enterprise** | 100+ | Custom | + On-prem, air-gapped, dedicated support |

**What You Get:**

- Everything in Cloud, PLUS:
- Single Sign-On (SAML/OIDC)
- Role-Based Access Control
- Audit logs for compliance
- Shared team memory pools
- Custom domain support
- Dedicated support channel
- SLA guarantees (99.9% uptime)
- On-premises deployment option

**Best For:**

- Regulated industries (finance, healthcare)
- Companies with compliance requirements
- Teams needing shared AI knowledge
- Organizations wanting vendor support

---

## 🧠 Psychology of Pricing

### Why $50 is the Right Price

1. **Impulse buy territory** — No procurement process needed
2. **Less than a dinner out** — Easy to justify
3. **Lifetime feels permanent** — Builds trust and loyalty
4. **Lower than competitors** — Cursor costs $20/month, GitHub Copilot $10/month
5. **One-time pain** — No recurring mental overhead

### Why $5/month for Cloud

1. **Cheaper than coffee** — Starbucks costs more
2. **Sticky once adopted** — Memory makes switching hard
3. **Recurring revenue** — Funds ongoing development
4. **Lower barrier than $50** — "I'll try it for a month"
5. **Automatic value** — Backups, sync, updates included

### The Expiration Psychology

**Local License (Never Expires):**

- "I own this forever" — feels safe
- No anxiety about cancellation
- Trust in the product
- Willing to pay more upfront

**Cloud Subscription (Expires if Unpaid):**

- "I'm renting convenience" — expected
- Data export tools before cancellation
- Grace period (30 days) before deletion
- Clear value for ongoing payment

---

## 📊 Competitive Analysis

| Product | Price | Memory | Tools | Cloud |
|---------|-------|--------|-------|-------|
| **HyperNexus Local** | $50 lifetime | ✅ L1/L2 | 20,000+ | ❌ |
| **HyperNexus Cloud** | $5/month | ✅ L1/L2/L3 | 20,000+ | ✅ |
| Cursor | $20/month | ❌ | Limited | ✅ |
| GitHub Copilot | $10/month | ❌ | Limited | ✅ |
| Claude Code | $20/month | ❌ | Limited | ✅ |
| Aider | Free | ❌ | Limited | ❌ |
| Continue | Free | ❌ | Limited | ❌ |

**HyperNexus Differentiation:**

- Only product with persistent multi-tier memory
- Largest MCP tool catalog (20,000+)
- Works with 38+ AI agents (not locked to one)
- Local-first with optional cloud
- Open-source core

---

## 🚀 Implementation Plan

### Phase 1: Stripe Integration (Week 1)

- [ ] Create Stripe products (local license, cloud subscription)
- [ ] Implement checkout flow on hypernexus.site
- [ ] License key generation (Ed25519-signed tokens)
- [ ] Email delivery of license keys
- [ ] Account creation and management

### Phase 2: License Validation (Week 2)

- [ ] License key validation in Go sidecar
- [ ] Grace period handling (7 days for cloud)
- [ ] License expiration notifications
- [ ] Offline validation for local licenses
- [ ] Upgrade path from v1 to v2

### Phase 3: Cloud Billing (Week 3)

- [ ] Usage metering (API calls, storage)
- [ ] Invoice generation
- [ ] Payment failure handling
- [ ] Account suspension flow
- [ ] Data export before deletion

### Phase 4: Enterprise Sales (Week 4)

- [ ] Sales one-pager
- [ ] Corporate pitch deck
- [ ] Demo environment
- [ ] Pilot program structure
- [ ] Custom pricing calculator

---

## 💡 Pricing Experiments

### A/B Tests to Run

1. **Price sensitivity** — Test $40 vs $50 vs $60
2. **Annual discount** — Test $50/year vs $5/month
3. **Trial length** — Test 7-day vs 14-day vs 30-day
4. **Feature gating** — Test what drives upgrades

### Metrics to Track

| Metric | Target | Measurement |
|--------|--------|-------------|
| Conversion rate | 5% | Visitors → signups |
| Trial-to-paid | 20% | Trial → paying |
| Monthly churn | <5% | Cancellations / active |
| LTV:CAC | >3:1 | Lifetime value / acquisition cost |
| Expansion revenue | 20% | Upgrades / total revenue |

---

## 🎯 Go-to-Market Strategy

### Launch Sequence

1. **Product Hunt** — Free tier generates awareness
2. **Reddit (r/LocalLLaMA)** — Community-driven adoption
3. **Hacker News** — Technical credibility
4. **Twitter/X** — Viral content from marketing agent
5. **YouTube** — Demo videos and tutorials
6. **Blog** — SEO content for long-tail traffic

### Messaging Framework

**For Developers:**
> "Your AI finally remembers everything. Install in 30 seconds, works with every tool you already use."

**For Teams:**
> "Shared AI memory for your entire team. Stop re-explaining your codebase to every new hire."

**For Enterprise:**
> "Compliance-ready AI infrastructure. Full audit logs, SSO, and sandboxed execution."

---

## 📋 FAQ

**Q: What happens to my data if I stop paying for cloud?**
A: You have 30 days to export everything. After that, data is permanently deleted. Local licenses never have this issue.

**Q: Can I switch from cloud to local?**
A: Yes! Export your memory from cloud and import locally. Your license converts automatically.

**Q: Is there a free tier?**
A: The open-source core is free forever. Cloud features require a subscription.

**Q: Do you offer refunds?**
A: Yes, 30-day money-back guarantee on all plans.

**Q: Can I self-host the cloud version?**
A: Enterprise tier includes on-premises deployment. Contact <sales@hypernexus.site>.
