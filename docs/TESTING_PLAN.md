# TormentNexus — Full Platform Testing Plan

> **Created:** 2026-07-14
> **Purpose:** Verify billing, deployment, and full platform functionality
> **Status:** In Progress

---

## 🎯 Testing Philosophy

**Test everything. Trust nothing.** Every component must be verified end-to-end before we can claim the platform is production-ready.

---

## 📋 Test Categories

### 1. Infrastructure & Deployment

### 2. API Endpoints

### 3. Billing & Payments (Stripe)

### 4. Memory System (L1-L4)

### 5. Catalog & Search

### 6. Dashboard & UI

### 7. Security & Authentication

### 8. Multi-Tenant Support

### 9. Performance & Load

### 10. End-to-End Workflows

---

## 1. Infrastructure & Deployment

### 1.1 Health Checks

- [ ] TN Kernel health (`/health`)
- [ ] Dashboard health (port 7779)
- [ ] Demo container health (port 7790)
- [ ] Nginx reverse proxy
- [ ] SSL certificate validity

### 1.2 Service Status

- [ ] PM2 process status
- [ ] Docker container status
- [ ] Port accessibility (7778, 7779, 7790)
- [ ] Log file rotation

### 1.3 Deployment Pipeline

- [ ] Git push triggers CI
- [ ] CI builds successfully
- [ ] Docker image builds
- [ ] Docker image deploys
- [ ] Health check passes after deploy

---

## 2. API Endpoints

### 2.1 Core APIs

- [ ] `GET /health` — Health check
- [ ] `GET /api/status` — System status
- [ ] `GET /api/index` — API index

### 2.2 Memory APIs

- [ ] `GET /api/memory/stats` — Memory statistics
- [ ] `GET /api/memory/search?q=...` — Search memories
- [ ] `POST /api/memory/store` — Store memory
- [ ] `GET /api/memory/project/sync` — Project sync

### 2.3 Catalog APIs

- [ ] `GET /api/backlog/search?q=...` — Search catalog
- [ ] `GET /api/backlog/stats` — Catalog statistics
- [ ] `GET /api/backlog/categories` — List categories

### 2.4 Account APIs

- [ ] `POST /api/account/register` — Register user
- [ ] `POST /api/account/login` — Login
- [ ] `GET /api/account/status` — Account status

### 2.5 Provider APIs

- [ ] `GET /api/providers/catalog` — Provider catalog
- [ ] `GET /api/providers/routing-summary` — Routing summary

### 2.6 Session APIs

- [ ] `GET /api/sessions/` — List sessions
- [ ] `POST /api/sessions/` — Create session

---

## 3. Billing & Payments (Stripe)

### 3.1 Checkout Flow

- [ ] Pricing page loads
- [ ] "Get Started" button works
- [ ] Stripe checkout session creates
- [ ] Payment form loads
- [ ] Test card (4242...) works
- [ ] Payment succeeds
- [ ] Redirect to success page

### 3.2 Webhook Handling

- [ ] Webhook endpoint exists
- [ ] Webhook signature verification
- [ ] `checkout.session.completed` event
- [ ] `customer.subscription.created` event
- [ ] `customer.subscription.updated` event
- [ ] `customer.subscription.deleted` event
- [ ] `invoice.payment_succeeded` event
- [ ] `invoice.payment_failed` event

### 3.3 Subscription Management

- [ ] Subscription created in database
- [ ] Subscription status updates
- [ ] Plan upgrade works
- [ ] Plan downgrade works
- [ ] Cancellation works
- [ ] Re-subscription works

### 3.4 Billing Portal

- [ ] Portal link generation
- [ ] Portal loads
- [ ] Update payment method
- [ ] View invoices
- [ ] Download invoices

---

## 4. Memory System (L1-L4)

### 4.1 L1 Session Memory

- [ ] Store session memory
- [ ] Retrieve session memory
- [ ] Session isolation
- [ ] Session cleanup

### 4.2 L2 Hot Store

- [ ] Store to hot store
- [ ] Search hot store
- [ ] Vector similarity search
- [ ] TTL expiration (30 days)

### 4.3 L3 Cold Archive

- [ ] Archive to cold storage
- [ ] Search cold archive
- [ ] Retrieval from cold archive
- [ ] TTL expiration (1 year)

### 4.4 L4 Limbo

- [ ] Soft delete to limbo
- [ ] Restore from limbo
- [ ] Permanent delete
- [ ] TTL expiration (90 days)

### 4.5 Graph Relations

- [ ] Create relation
- [ ] Query relations
- [ ] Traverse graph
- [ ] Delete relation

---

## 5. Catalog & Search

### 5.1 Search Functionality

- [ ] Full-text search works
- [ ] Category filter works
- [ ] Source filter works
- [ ] Limit parameter works
- [ ] Empty results handled

### 5.2 Data Integrity

- [ ] Total count accurate
- [ ] No duplicate entries
- [ ] URLs are valid
- [ ] Descriptions present (enrichment)

### 5.3 Performance

- [ ] Search response < 100ms
- [ ] Stats response < 50ms
- [ ] Categories response < 50ms

---

## 6. Dashboard & UI

### 6.1 Page Loads

- [ ] Home page loads
- [ ] Memory tab loads
- [ ] Tools tab loads
- [ ] Catalog tab loads
- [ ] Agents tab loads
- [ ] Security tab loads
- [ ] Infrastructure tab loads
- [ ] Commercial tab loads

### 6.2 Functionality

- [ ] Memory search works
- [ ] Catalog search works
- [ ] Tool execution works
- [ ] Agent monitoring works
- [ ] Security alerts work

### 6.3 Responsiveness

- [ ] Desktop layout
- [ ] Tablet layout
- [ ] Mobile layout

---

## 7. Security & Authentication

### 7.1 SSL/TLS

- [ ] HTTPS enforced
- [ ] Certificate valid
- [ ] Wildcard cert works
- [ ] HTTP redirects to HTTPS

### 7.2 Headers

- [ ] X-Frame-Options
- [ ] X-Content-Type-Options
- [ ] Content-Security-Policy
- [ ] Strict-Transport-Security

### 7.3 Authentication

- [ ] Registration works
- [ ] Login works
- [ ] Logout works
- [ ] Password hashing
- [ ] Session management
- [ ] Token validation

### 7.4 Authorization

- [ ] Protected routes require auth
- [ ] Role-based access works
- [ ] API key authentication works

---

## 8. Multi-Tenant Support

### 8.1 Tenant Provisioning

- [ ] Tenant creation
- [ ] Subdomain assignment
- [ ] Database isolation
- [ ] Config isolation

### 8.2 Tenant Management

- [ ] Tenant listing
- [ ] Tenant status
- [ ] Tenant deletion
- [ ] Tenant suspension

### 8.3 Tenant Isolation

- [ ] Data isolation
- [ ] Process isolation
- [ ] Network isolation

---

## 9. Performance & Load

### 9.1 Response Times

- [ ] Health check < 50ms
- [ ] API calls < 200ms
- [ ] Search < 500ms
- [ ] Page load < 2s

### 9.2 Resource Usage

- [ ] CPU usage < 80%
- [ ] Memory usage < 80%
- [ ] Disk usage < 80%
- [ ] Network bandwidth

### 9.3 Concurrent Users

- [ ] 10 concurrent users
- [ ] 100 concurrent users
- [ ] 1000 concurrent users

---

## 10. End-to-End Workflows

### 10.1 New User Flow

1. Visit landing page
2. Click "Try Demo"
3. Explore demo
4. Click "Sign Up"
5. Register account
6. Choose plan
7. Complete payment
8. Access dashboard
9. Create first project
10. Store first memory

### 10.2 Developer Flow

1. Visit GitHub
2. Clone repository
3. Run `npx tormentnexus serve`
4. Open dashboard
5. Search catalog
6. Install MCP tool
7. Configure tool
8. Use tool

### 10.3 Enterprise Flow

1. Contact sales
2. Get custom quote
3. Deploy to own infrastructure
4. Configure SSO
5. Set up RBAC
6. Enable audit logs
7. Go live

---

## 🧪 Test Scripts

### Automated Tests

```bash
# Run all tests
./scripts/test-all.sh

# Run specific category
./scripts/test-api.sh
./scripts/test-billing.sh
./scripts/test-memory.sh
./scripts/test-catalog.sh
./scripts/test-security.sh
```

### Manual Tests

- Follow the checklist above
- Document results
- Report issues

---

## 📊 Test Results Template

| Test | Status | Notes |
|------|--------|-------|
| Health check | ✅ PASS | 200 OK |
| API search | ✅ PASS | Returns results |
| Billing checkout | ❌ FAIL | Stripe not configured |
| Memory store | ✅ PASS | Stores correctly |
| ... | ... | ... |

---

## 🐛 Bug Report Template

```markdown
## Bug Description
[Clear description]

## Steps to Reproduce
1. Go to '...'
2. Click on '...'
3. See error

## Expected Behavior
[What should happen]

## Actual Behavior
[What actually happens]

## Environment
- OS: [e.g. Windows 11]
- Browser: [e.g. Chrome 120]
- Version: [e.g. 1.0.0-alpha.257]

## Logs
```

[Paste relevant logs]

```
```

---

## ✅ Sign-off

- [ ] All tests pass
- [ ] No critical bugs
- [ ] Performance acceptable
- [ ] Security verified
- [ ] Documentation updated
- [ ] Ready for production
