# TormentNexus — Billing & Deployment Test Plan

> **Created:** 2026-07-14
> **Purpose:** Verify billing, deployment, and FULL platform functionality
> **Status:** Executing Now

---

## 🎯 Test Categories

### A. Billing & Payments

### B. Deployment & Infrastructure

### C. API Endpoints (ALL)

### D. Memory System

### E. Catalog & Search

### F. Dashboard & UI

### G. Security & Auth

### H. Multi-Tenant

### I. Performance

### J. End-to-End Workflows

---

## A. Billing & Payments

### A1. Stripe Configuration

- [ ] Stripe API key configured
- [ ] Webhook secret configured
- [ ] Price IDs configured
- [ ] Test mode enabled

### A2. Checkout Flow

- [ ] Pricing page loads
- [ ] "Get Started" button works
- [ ] Stripe checkout session creates
- [ ] Payment form loads
- [ ] Test card (4242 4242 4242 4242) works
- [ ] Payment succeeds
- [ ] Redirect to success page
- [ ] Receipt email sent

### A3. Subscription Management

- [ ] Subscription created in database
- [ ] Subscription status updates
- [ ] Plan upgrade works
- [ ] Plan downgrade works
- [ ] Cancellation works
- [ ] Re-subscription works

### A4. Webhook Handling

- [ ] Webhook endpoint exists
- [ ] Webhook signature verification
- [ ] checkout.session.completed
- [ ] customer.subscription.created
- [ ] customer.subscription.updated
- [ ] customer.subscription.deleted
- [ ] invoice.payment_succeeded
- [ ] invoice.payment_failed

### A5. Billing Portal

- [ ] Portal link generation
- [ ] Portal loads
- [ ] Update payment method
- [ ] View invoices
- [ ] Download invoices

---

## B. Deployment & Infrastructure

### B1. Docker

- [ ] Docker image builds
- [ ] Docker image pushes to registry
- [ ] Container starts
- [ ] Health check passes
- [ ] Port mapping works
- [ ] Volume mounts work
- [ ] Environment variables work

### B2. PM2 Services

- [ ] tn-kernel online
- [ ] fwber-backend online
- [ ] fwber-frontend online
- [ ] Auto-restart works

### B3. Nginx

- [ ] Reverse proxy works
- [ ] SSL termination works
- [ ] HTTP->HTTPS redirect
- [ ] Security headers present
- [ ] Rate limiting configured

### B4. DNS

- [ ] tormentnexus.site resolves
- [ ] cloud.hypernexus.site resolves
- [ ] demo.hypernexus.site resolves
- [ ] *.hypernexus.site resolves

### B5. CI/CD

- [ ] GitHub Actions triggers on push
- [ ] Go build succeeds
- [ ] Docker build succeeds
- [ ] Tests pass
- [ ] Deployment succeeds

---

## C. API Endpoints (ALL)

### C1. Core APIs

- [ ] GET /health
- [ ] GET /api/index
- [ ] GET /api/status

### C2. Account APIs

- [ ] POST /api/account/register
- [ ] POST /api/account/login
- [ ] GET /api/account/status
- [ ] POST /api/account/provision

### C3. Catalog APIs

- [ ] GET /api/backlog/search
- [ ] GET /api/backlog/stats
- [ ] GET /api/backlog/categories

### C4. Memory APIs

- [ ] GET /api/memory/stats
- [ ] GET /api/memory/search
- [ ] POST /api/memory/store

### C5. Provider APIs

- [ ] GET /api/providers/catalog
- [ ] GET /api/providers/routing-summary

### C6. Session APIs

- [ ] GET /api/sessions
- [ ] POST /api/sessions

### C7. MCP APIs

- [ ] GET /api/mcp/servers
- [ ] POST /api/mcp/servers
- [ ] GET /api/mcp/tools

---

## D. Memory System

### D1. L1 Session Memory

- [ ] Store session memory
- [ ] Retrieve session memory
- [ ] Session isolation

### D2. L2 Hot Store

- [ ] Store to hot store
- [ ] Search hot store
- [ ] Vector similarity search

### D3. L3 Cold Archive

- [ ] Archive to cold storage
- [ ] Search cold archive

### D4. L4 Limbo

- [ ] Soft delete to limbo
- [ ] Restore from limbo

---

## E. Catalog & Search

### E1. Search

- [ ] Full-text search works
- [ ] Category filter works
- [ ] Source filter works
- [ ] Limit parameter works
- [ ] Empty results handled

### E2. Data

- [ ] Total count accurate
- [ ] No duplicates
- [ ] URLs valid
- [ ] Descriptions present

### E3. Performance

- [ ] Search < 100ms
- [ ] Stats < 50ms
- [ ] Categories < 50ms

---

## F. Dashboard & UI

### F1. Pages

- [ ] Home loads
- [ ] Memory tab loads
- [ ] Tools tab loads
- [ ] Catalog tab loads
- [ ] Agents tab loads
- [ ] Security tab loads

### F2. Functionality

- [ ] Memory search works
- [ ] Catalog search works
- [ ] Tool execution works

---

## G. Security & Auth

### G1. SSL

- [ ] HTTPS enforced
- [ ] Certificate valid
- [ ] Wildcard cert works

### G2. Headers

- [ ] X-Frame-Options
- [ ] X-Content-Type-Options
- [ ] Content-Security-Policy

### G3. Auth

- [ ] Registration works
- [ ] Login works
- [ ] Logout works
- [ ] Token validation

---

## H. Multi-Tenant

### H1. Provisioning

- [ ] Tenant creation
- [ ] Subdomain assignment
- [ ] Database isolation

### H2. Management

- [ ] Tenant listing
- [ ] Tenant status
- [ ] Tenant deletion

---

## I. Performance

### I1. Response Times

- [ ] Health < 50ms
- [ ] API < 200ms
- [ ] Search < 500ms
- [ ] Page load < 2s

### I2. Resources

- [ ] CPU < 80%
- [ ] Memory < 80%
- [ ] Disk < 80%

---

## J. End-to-End Workflows

### J1. New User Flow

1. Visit landing page
2. Click "Try Demo"
3. Search catalog
4. Click "Sign Up"
5. Register account
6. Choose plan
7. Complete payment
8. Access dashboard

### J2. Developer Flow

1. Visit GitHub
2. Clone repo
3. Run `npx tormentnexus serve`
4. Open dashboard
5. Search catalog
6. Install MCP tool

### J3. Enterprise Flow

1. Contact sales
2. Get quote
3. Deploy
4. Configure SSO
5. Go live

---

## 🧪 Automated Test Commands

```bash
# Full test suite
bash scripts/test-final.sh

# Billing tests
bash scripts/test-billing.sh

# Deployment tests
bash scripts/test-deployment.sh

# API tests
bash scripts/test-api.sh

# Performance tests
bash scripts/test-performance.sh
```

---

## ✅ Sign-off

- [ ] All tests pass
- [ ] No critical bugs
- [ ] Performance acceptable
- [ ] Security verified
- [ ] Documentation updated
- [ ] Ready for production
