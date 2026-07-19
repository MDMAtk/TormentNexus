# Marketing Agent — Stripe Billing Integration

## Purpose

This directory bridges TormentNexus local billing with the cloud-hosted
hypernexus.site dashboard. The Go sidecar (`port 7778`) provides all
Stripe-related API endpoints; the marketing agent configures pricing,
checkout flows, and customer lifecycle.

## API Endpoints (Go sidecar → hypernexus.site dashboard)

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/billing/stripe/plans` | GET | List available subscription plans |
| `/api/billing/stripe/checkout` | POST | Create Stripe Checkout session |
| `/api/billing/stripe/portal` | POST | Generate Customer Portal link |
| `/api/billing/stripe/webhook` | POST | Receive Stripe lifecycle events |
| `/api/billing/stripe/subscription` | GET | Get current subscription status |
| `/api/billing/stripe/subscribe` | POST | Legacy subscribe (simulated) |

## Environment Variables

Set these on the hypernexus.site deployment:

```
STRIPE_SECRET_KEY=sk_live_xxx
STRIPE_WEBHOOK_SECRET=whsec_xxx
STRIPE_PRICE_ID_BASIC=price_basic_monthly
STRIPE_PRICE_ID_PRO=price_pro_monthly
STRIPE_PRICE_ID_ENTERPRISE=price_enterprise_monthly
TORMENTNEXUS_DASHBOARD_URL=https://hypernexus.site
TORMENTNEXUS_API_URL=https://api.hypernexus.site
```

## Plans

| Plan | Price | Features |
|------|-------|----------|
| Basic | $29/mo | 1 user, 100K tokens |
| Pro | $99/mo | 5 users, 1M tokens |
| Enterprise | $499/mo | Unlimited everything |

## Webhook Events Handled

- `checkout.session.completed` — New subscription → activate
- `customer.subscription.updated` — Status changes (active, past_due, canceled)
- `customer.subscription.deleted` — Cancelation
- `invoice.payment_succeeded` — Renewal → update next invoice date
- `invoice.payment_failed` — Failed payment → mark status

## Local Dev Mode

When `STRIPE_SECRET_KEY` is empty, all endpoints simulate responses.
The dashboard shows "Visa ending in 4242 (simulated)" and stores config
locally in `.tormentnexus/config.json`.
