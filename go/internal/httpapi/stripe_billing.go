package httpapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ── Configuration ────────────────────────────────────────────────
// These are read from environment variables for cloud deployments.
// When running locally, they fall back to config DB values.
//   STRIPE_SECRET_KEY         — sk_live_xxx or sk_test_xxx
//   STRIPE_WEBHOOK_SECRET     — whsec_xxx
//   STRIPE_PRICE_ID_BASIC     — price_xxx for $29/mo plan
//   STRIPE_PRICE_ID_PRO       — price_xxx for $99/mo plan
//   STRIPE_PRICE_ID_COMMERCIAL — price_xxx for $499/mo plan
//   TORMENTNEXUS_DASHBOARD_URL — https://hypernexus.site
//   TORMENTNEXUS_API_URL      — https://api.hypernexus.site

func stripeSecretKey() string {
	if k := os.Getenv("STRIPE_SECRET_KEY"); k != "" {
		return k
	}
	return "" // Will fail gracefully — local mode
}

func stripeWebhookSecret() string {
	return os.Getenv("STRIPE_WEBHOOK_SECRET")
}

func stripePriceID(plan string) string {
	switch plan {
	case "basic":
		return os.Getenv("STRIPE_PRICE_ID_BASIC")
	case "pro":
		return os.Getenv("STRIPE_PRICE_ID_PRO")
	case "commercial":
		return os.Getenv("STRIPE_PRICE_ID_COMMERCIAL")
	}
	return os.Getenv("STRIPE_PRICE_ID_COMMERCIAL")
}

func dashboardURL() string {
	if u := os.Getenv("TORMENTNEXUS_DASHBOARD_URL"); u != "" {
		return u
	}
	return "https://hypernexus.site"
}

func apiBaseURL() string {
	if u := os.Getenv("TORMENTNEXUS_API_URL"); u != "" {
		return u
	}
	return "http://127.0.0.1:7778"
}

// ── Plan Definitions ─────────────────────────────────────────────

type StripePlan struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Interval    string   `json:"interval"`
	PriceID     string   `json:"priceId"`
	Features    []string `json:"features"`
}

func availablePlans() []StripePlan {
	return []StripePlan{
		{
			ID:          "basic",
			Name:        "Basic",
			Description: "For individuals and small teams",
			Price:       29.00,
			Interval:    "month",
			PriceID:     stripePriceID("basic"),
			Features:    []string{"1 user", "100K tokens/month", "Basic support", "Community access"},
		},
		{
			ID:          "pro",
			Name:        "Pro",
			Description: "For professional developers",
			Price:       99.00,
			Interval:    "month",
			PriceID:     stripePriceID("pro"),
			Features:    []string{"5 users", "1M tokens/month", "Priority support", "API access", "Custom models"},
		},
		{
			ID:          "commercial",
			Name:        "Commercial Cloud SaaS",
			Description: "For organizations with advanced needs",
			Price:       499.00,
			Interval:    "month",
			PriceID:     stripePriceID("commercial"),
			Features:    []string{"Unlimited users", "Unlimited tokens", "24/7 support", "SSO/SAML", "Dedicated infrastructure", "SLA guarantees"},
		},
	}
}

// ── HTTP Handlers ────────────────────────────────────────────────

// handleStripePlans returns available subscription plans.
func (s *Server) handleStripePlans(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    availablePlans(),
	})
}

// handleStripeCreateCheckout creates a Stripe Checkout Session for subscription.
// POST /api/billing/stripe/checkout
// Body: { "plan": "pro", "successUrl": "...", "cancelUrl": "..." }
func (s *Server) handleStripeCreateCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"success": false, "error": "POST required"})
		return
	}

	var payload struct {
		Plan       string `json:"plan"`
		SuccessURL string `json:"successUrl"`
		CancelURL  string `json:"cancelUrl"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON"})
		return
	}

	// Validate plan
	priceID := stripePriceID(payload.Plan)
	if priceID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": fmt.Sprintf("unknown plan: %q or missing STRIPE_PRICE_ID_%s env var", payload.Plan, strings.ToUpper(payload.Plan))})
		return
	}

	// If no Stripe key is configured, return a simulated response for local dev
	if stripeSecretKey() == "" {
		s.handleLocalCheckoutSimulation(w, r, payload.Plan)
		return
	}

	// Real Stripe API call
	successURL := payload.SuccessURL
	if successURL == "" {
		successURL = dashboardURL() + "/billing?checkout=success"
	}
	cancelURL := payload.CancelURL
	if cancelURL == "" {
		cancelURL = dashboardURL() + "/billing?checkout=canceled"
	}

	// Get or create customer ID from local config
	customerID, _ := s.localConfigValue("stripe.customerID")
	custID := ""
	if customerID != nil {
		custID = customerID.(string)
	}

	body := fmt.Sprintf(
		"mode=subscription&success_url=%s&cancel_url=%s&line_items[0][price]=%s&line_items[0][quantity]=1&line_items[0][adjustable_quantity][enabled]=true&line_items[0][adjustable_quantity][minimum]=1&line_items[0][adjustable_quantity][maximum]=999",
		successURL, cancelURL, priceID,
	)
	if custID != "" {
		body += "&customer=" + custID
	}

	req, _ := http.NewRequestWithContext(r.Context(), "POST", "https://api.stripe.com/v1/checkout/sessions", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+stripeSecretKey())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"success": false, "error": "Stripe API call failed: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var stripeResp map[string]any
	if err := json.Unmarshal(respBody, &stripeResp); err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"success": false, "error": "failed to parse Stripe response", "raw": string(respBody)})
		return
	}

	if resp.StatusCode >= 400 {
		writeJSON(w, http.StatusBadGateway, map[string]any{
			"success": false,
			"error":   "Stripe API error",
			"stripe":  stripeResp,
		})
		return
	}

	// Store the session ID and customer ID from Stripe
	if sessionID, ok := stripeResp["id"].(string); ok {
		s.setLocalConfigValue("stripe.checkoutSessionID", sessionID)
	}
	if custID, ok := stripeResp["customer"].(string); ok {
		s.setLocalConfigValue("stripe.customerID", custID)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"sessionId":  stripeResp["id"],
			"sessionUrl": stripeResp["url"],
			"customerID": stripeResp["customer"],
			"plan":       payload.Plan,
			"priceID":    priceID,
		},
	})
}

// handleLocalCheckoutSimulation provides a dev-mode checkout without Stripe.
func (s *Server) handleLocalCheckoutSimulation(w http.ResponseWriter, r *http.Request, plan string) {
	planName := plan
	for _, p := range availablePlans() {
		if p.ID == plan {
			planName = p.Name
			break
		}
	}

	customerID := "cus_" + strconv.FormatInt(time.Now().Unix(), 36)
	sessionID := "cs_" + strconv.FormatInt(time.Now().Unix(), 36)

	s.setLocalConfigValue("stripe.plan", planName)
	s.setLocalConfigValue("stripe.status", "ACTIVE (PAID)")
	s.setLocalConfigValue("stripe.customerID", customerID)
	s.setLocalConfigValue("stripe.checkoutSessionID", sessionID)

	for _, p := range availablePlans() {
		if p.ID == plan {
			s.setLocalConfigValue("stripe.price", fmt.Sprintf("$%.2f / month", p.Price))
			break
		}
	}
	s.setLocalConfigValue("stripe.nextInvoice", time.Now().AddDate(0, 1, 0).Format("January 02, 2006"))
	s.setLocalConfigValue("stripe.paymentSource", "Visa ending in 4242 (simulated)")

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"sessionId":  sessionID,
			"sessionUrl": dashboardURL() + "/billing?checkout=simulated&plan=" + plan,
			"customerID": customerID,
			"plan":       plan,
			"simulated":  true,
			"message":    "Local dev mode — no Stripe API key configured. Subscription simulated.",
		},
	})
}

// handleStripeCustomerPortal creates a Stripe Customer Portal session.
// POST /api/billing/stripe/portal
func (s *Server) handleStripeCustomerPortal(w http.ResponseWriter, r *http.Request) {
	if stripeSecretKey() == "" {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data": map[string]any{
				"url":       dashboardURL() + "/billing",
				"simulated": true,
				"message":   "Local dev mode — no Stripe API key configured.",
			},
		})
		return
	}

	customerID, _ := s.localConfigValue("stripe.customerID")
	custID := ""
	if customerID != nil {
		custID = customerID.(string)
	}
	if custID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "no customer ID found. subscribe first."})
		return
	}

	body := fmt.Sprintf("customer=%s&return_url=%s/billing", custID, dashboardURL())
	req, _ := http.NewRequestWithContext(r.Context(), "POST", "https://api.stripe.com/v1/billing_portal/sessions", bytes.NewBufferString(body))
	req.Header.Set("Authorization", "Bearer "+stripeSecretKey())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"success": false, "error": "Stripe portal API failed: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var portalResp map[string]any
	json.Unmarshal(respBody, &portalResp)

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"url":        portalResp["url"],
			"customerID": custID,
		},
	})
}

// handleStripeWebhook receives Stripe webhook events.
// POST /api/billing/stripe/webhook
func (s *Server) handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"success": false, "error": "POST required"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "failed to read body"})
		return
	}

	// Verify webhook signature if configured
	whSecret := stripeWebhookSecret()
	if whSecret != "" {
		sigHeader := r.Header.Get("Stripe-Signature")
		if sigHeader == "" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "missing Stripe-Signature header"})
			return
		}
		if !verifyStripeSignature(whSecret, string(body), sigHeader) {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"success": false, "error": "invalid webhook signature"})
			return
		}
	}

	var event struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID            string `json:"id"`
				Customer      string `json:"customer"`
				Status        string `json:"status"`
				Subscription  string `json:"subscription"`
				PaymentIntent string `json:"payment_intent"`
				BillingReason string `json:"billing_reason"`
				Items         struct {
					Data []struct {
						Price struct {
							ID         string  `json:"id"`
							Nickname   string  `json:"nickname"`
							UnitAmount float64 `json:"unit_amount"`
						} `json:"price"`
					} `json:"data"`
				} `json:"lines"`
			} `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &event); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid event JSON"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		obj := event.Data.Object
		planName := "Pro"
		if len(obj.Items.Data) > 0 {
			for _, p := range availablePlans() {
				if p.PriceID == obj.Items.Data[0].Price.ID {
					planName = p.Name
					break
				}
			}
		}
		s.setLocalConfigValue("stripe.customerID", obj.Customer)
		s.setLocalConfigValue("stripe.subscriptionID", obj.Subscription)
		s.setLocalConfigValue("stripe.plan", planName)
		s.setLocalConfigValue("stripe.status", "ACTIVE (PAID)")

	case "customer.subscription.updated":
		obj := event.Data.Object
		status := obj.Status
		switch status {
		case "active":
			s.setLocalConfigValue("stripe.status", "ACTIVE (PAID)")
		case "past_due":
			s.setLocalConfigValue("stripe.status", "PAST DUE")
		case "canceled":
			s.setLocalConfigValue("stripe.status", "CANCELED")
		case "trialing":
			s.setLocalConfigValue("stripe.status", "TRIALING")
		case "incomplete", "incomplete_expired":
			s.setLocalConfigValue("stripe.status", strings.ToUpper(status))
		default:
			s.setLocalConfigValue("stripe.status", strings.ToUpper(status))
		}
		s.setLocalConfigValue("stripe.subscriptionID", obj.ID)
		s.setLocalConfigValue("stripe.customerID", obj.Customer)

	case "customer.subscription.deleted":
		s.setLocalConfigValue("stripe.status", "CANCELED")
		s.setLocalConfigValue("stripe.subscriptionID", event.Data.Object.ID)

	case "invoice.payment_succeeded":
		obj := event.Data.Object
		s.setLocalConfigValue("stripe.status", "ACTIVE (PAID)")
		s.setLocalConfigValue("stripe.lastInvoiceID", obj.ID)
		if obj.PaymentIntent != "" {
			s.setLocalConfigValue("stripe.lastPaymentIntent", obj.PaymentIntent)
		}
		nextInvoice := time.Now().AddDate(0, 1, 0).Format("January 02, 2006")
		s.setLocalConfigValue("stripe.nextInvoice", nextInvoice)
		if obj.BillingReason == "subscription_create" {
			s.setLocalConfigValue("stripe.nextInvoice", time.Now().AddDate(0, 1, 0).Format("January 02, 2006"))
		}

	case "invoice.payment_failed":
		s.setLocalConfigValue("stripe.status", "PAYMENT FAILED")
		s.setLocalConfigValue("stripe.lastFailedInvoiceID", event.Data.Object.ID)
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true, "received": true, "event": event.Type})
}

// handleStripeGetSubscription returns the current subscription status.
func (s *Server) handleStripeGetSubscription(w http.ResponseWriter, r *http.Request) {
	plan, _ := s.localConfigValue("stripe.plan")
	status, _ := s.localConfigValue("stripe.status")
	price, _ := s.localConfigValue("stripe.price")
	invoice, _ := s.localConfigValue("stripe.nextInvoice")
	source, _ := s.localConfigValue("stripe.paymentSource")
	customerID, _ := s.localConfigValue("stripe.customerID")
	subscriptionID, _ := s.localConfigValue("stripe.subscriptionID")

	d := map[string]any{
		"plan":                  plan,
		"status":                status,
		"price":                 price,
		"nextInvoice":           invoice,
		"paymentSource":         source,
		"customerID":            customerID,
		"subscriptionID":        subscriptionID,
		"hasActiveSubscription": status != nil && (status.(string) == "ACTIVE (PAID)" || status.(string) == "TRIALING"),
	}

	// Try live Stripe lookup if configured
	if stripeSecretKey() != "" && customerID != nil && customerID.(string) != "" {
		req, _ := http.NewRequestWithContext(r.Context(), "GET", "https://api.stripe.com/v1/customers/"+customerID.(string), nil)
		req.Header.Set("Authorization", "Bearer "+stripeSecretKey())
		if resp, err := http.DefaultClient.Do(req); err == nil {
			defer resp.Body.Close()
			if respBody, err := io.ReadAll(resp.Body); err == nil {
				var cust map[string]any
				if json.Unmarshal(respBody, &cust) == nil {
					if email, ok := cust["email"].(string); ok {
						d["email"] = email
					}
					if name, ok := cust["name"].(string); ok {
						d["customerName"] = name
					}
				}
			}
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    d,
	})
}

// ── Webhook Signature Verification ───────────────────────────────

func verifyStripeSignature(secret, payload, signatureHeader string) bool {
	// Stripe sends: t=timestamp,v1=signature
	parts := strings.Split(signatureHeader, ",")
	var timestamp string
	var expectedSig string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "t=") {
			timestamp = part[2:]
		} else if strings.HasPrefix(part, "v1=") {
			expectedSig = part[3:]
		}
	}
	if timestamp == "" || expectedSig == "" {
		return false
	}

	signedPayload := timestamp + "." + payload
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	computedSig := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(computedSig), []byte(expectedSig))
}
