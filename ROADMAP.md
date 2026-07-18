# ROADMAP — HyperNexus

> **Vision:** Make HyperNexus the "WordPress of AI" — easy to install, infinitely extensible, and the default choice for anyone who wants to give their AI persistent memory and tools.

---

## 💰 Pricing Strategy

### Two Models, One Product

| | **Local License** | **Cloud Subscription** |
|---|---|---|
| **Price** | $50 one-time (lifetime per major version) | $5/month |
| **Renewal** | $25/year for upgrades to new major versions | Auto-renews |
| **Storage** | Local SQLite + your own backups | Managed cloud storage |
| **Memory** | Your machine, your control | Synced across devices |
| **MCP Sandboxing** | Runs on your hardware | Isolated cloud containers |
| **Native Tools** | ✅ Full access | ✅ Full access |
| **Expiration** | Never (but upgrades cost $25/year) | Stops if you stop paying |
| **Best For** | Privacy-focused devs, teams with own infra | Convenience, teams, cloud-first |

### Why This Pricing Works

1. **$50 is impulse-buy territory** — No procurement process needed
2. **$5/month is cheaper than coffee** — Easy to justify for persistent AI memory
3. **Lifetime license builds trust** — You own it forever
4. **Cloud subscription funds development** — Recurring revenue sustains the project
5. **Both options exist** — Users choose what fits their workflow

### Revenue Projections

| Milestone | Users | MRR | ARR |
|-----------|-------|-----|-----|
| Launch (Month 1) | 50 | $250 | $3,000 |
| Growth (Month 3) | 500 | $2,500 | $30,000 |
| Scale (Month 6) | 2,000 | $10,000 | $120,000 |
| Maturity (Year 1) | 10,000 | $50,000 | $600,000 |

---

## 🎯 Strategic Pillars

### 1. **Distribution** — Be Everywhere

- GitLab (primary), GitHub (mirror)
- Docker Hub, npm, Homebrew
- VS Code Marketplace
- One-click install on any platform

### 2. **Community** — Build the Ecosystem

- MCP tool marketplace
- Discord community
- Open contribution model
- Plugin/extension system

### 3. **Product** — Solve Real Problems

- Persistent AI memory (L1/L2)
- 20,000+ MCP tool routing
- Multi-agent coordination
- Cloud sync and collaboration

### 4. **Monetization** — Sustainable Growth

- **Local License:** $50 lifetime
- **Cloud Subscription:** $5/month
- **Enterprise:** Custom pricing (SSO, RBAC, SLA)
- **Training:** Certification program

---

## 📅 Timeline

### **Phase 1: Launch** — July 2026

**Goal:** Public launch, first paying customers

| Week | Task | Status |
|------|------|--------|
| 1 | Cloud infrastructure deployed | ✅ |
| 1 | GitLab CI/CD pipeline | ✅ |
| 1 | SSL and monitoring | ✅ |
| 2 | Cloud landing page | 🔄 |
| 2 | Pricing page + Stripe | ⏳ |
| 2 | Product Hunt launch | ⏳ |
| 2 | Reddit/HN posts | ⏳ |

**Deliverables:**

- [x] Cloud API running at cloud.hypernexus.site
- [x] 14 services on Hetzner
- [x] CI/CD with auto-deploy
- [ ] Landing page with pricing
- [ ] Stripe checkout flow
- [ ] Product Hunt page

---

### **Phase 2: Growth** — August 2026

**Goal:** Build community, improve product, scale

| Task | Priority | Impact |
|------|----------|--------|
| MCP client examples | P0 | High |
| API documentation | P0 | High |
| VS Code extension | P0 | High |
| Community Discord | P1 | Medium |
| Video tutorials | P1 | Medium |
| Blog series | P1 | Medium |
| Enterprise outreach | P2 | High |

**Metrics:**

- 500+ GitLab stars
- 50 paying customers
- $500 MRR
- 50+ Discord members

---

### **Phase 3: Scale** — Q4 2026

**Goal:** Enterprise adoption, team features

| Task | Priority | Impact |
|------|----------|--------|
| Team accounts | P0 | High |
| SSO/SAML integration | P0 | High |
| Audit logs | P0 | Medium |
| Custom domains | P1 | Medium |
| White-label support | P1 | High |
| Mobile companion app | P2 | Medium |

**Metrics:**

- 5,000+ GitLab stars
- 500 paying customers
- $5,000 MRR
- 1+ enterprise pilot

---

### **Phase 4: Maturity** — 2027

**Goal:** Market leader, sustainable business

| Task | Priority | Impact |
|------|----------|--------|
| Plugin marketplace | P0 | High |
| Training certification | P0 | Medium |
| International expansion | P1 | Medium |
| Acquisition readiness | P2 | High |

**Metrics:**

- 20,000+ GitLab stars
- 2,000 paying customers
- $20,000 MRR
- Profitable operations

---

## 🏢 Enterprise Value Proposition

### Why Corporations Should Pay

**Problem:** Every developer using AI tools (Cursor, Claude, Copilot) has:

- No persistent memory between sessions
- No shared knowledge across team
- No audit trail of AI actions
- No sandboxing of tool execution
- No compliance controls

**Solution:** HyperNexus provides:

- ✅ Persistent memory that survives sessions
- ✅ Team-wide knowledge sharing
- ✅ Full audit logs of all AI actions
- ✅ Sandboxed MCP tool execution
- ✅ SSO/RBAC for compliance

### ROI for Enterprises

| Metric | Without HyperNexus | With HyperNexus | Improvement |
|--------|-------------------|-----------------|-------------|
| Context setup time | 15 min/session | 0 min (automatic) | 100% |
| Knowledge reuse | 0% (lost each session) | 80%+ | ∞ |
| Debugging time | 2 hours/issue | 30 min/issue | 75% |
| Onboarding new devs | 2 weeks | 3 days | 70% |
| AI tool compliance | None | Full audit trail | 100% |

### Corporate Pricing (Custom)

| Tier | Users | Price | Includes |
|------|-------|-------|----------|
| **Team** | 5-20 | $500/month | SSO, shared memory, audit logs |
| **Business** | 20-100 | $2,000/month | + RBAC, custom domains, SLA |
| **Enterprise** | 100+ | Custom | + On-prem, air-gapped, dedicated support |

---

## 🔧 Technical Roadmap

### Memory System (Stable)

- [x] L1 Session Memory (ephemeral)
- [x] L2 Semantic Vault (SQLite + vectors)
- [ ] L3 Cloud Sync (encrypted)
- [ ] L4 Team Knowledge Graph

### MCP Ecosystem (Stable)

- [x] 20,000+ servers indexed
- [x] Semantic tool routing
- [ ] Tool marketplace
- [ ] Custom tool builder

### Cloud Platform (Beta)

- [x] User accounts
- [x] Docker provisioning
- [x] MCP Streamable HTTP
- [ ] Memory sync
- [ ] Team collaboration
- [ ] Usage metering

### Agent Coordination (Beta)

- [x] A2A protocol
- [x] Role rotation
- [ ] Persistent agent identities
- [ ] Cross-session learning
- [ ] Agent marketplace

---

## 📊 Key Metrics

### Product Metrics

| Metric | Current | Target (6mo) |
|--------|---------|--------------|
| MCP Servers Indexed | 20,000+ | 50,000+ |
| Native Tools | 50+ | 200+ |
| Supported Agents | 38+ | 50+ |
| Test Coverage | 20% | 80% |
| API Latency (p95) | 50ms | 20ms |

### Business Metrics

| Metric | Current | Target (6mo) |
|--------|---------|--------------|
| Total Users | ~10 | 5,000 |
| Paying Customers | 0 | 500 |
| MRR | $0 | $5,000 |
| Churn Rate | N/A | <5% |
| NPS | N/A | 50+ |

---

## 🎯 Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-07-18 | Migrate to GitLab | Better CI/CD, avoid GitHub vendor lock-in |
| 2026-07-18 | $50 lifetime + $5/month cloud | Impulse-buy pricing, both options available |
| 2026-07-18 | Hetzner for hosting | Cost-effective, good performance |
| 2026-07-18 | MCP Streamable HTTP | Standard protocol, broad compatibility |
| 2026-07-18 | Sandboxed cloud containers | Security, isolation, compliance |
