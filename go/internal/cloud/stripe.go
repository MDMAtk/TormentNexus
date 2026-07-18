package cloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// StripeConfig holds Stripe configuration
type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	PriceID        string
	SuccessURL     string
	CancelURL      string
}

// StripeCheckoutSession represents a Stripe checkout session request
type StripeCheckoutSession struct {
	Email    string `json:"email"`
	Quantity int    `json:"quantity"`
}

// StripeCheckoutResponse represents the response from creating a checkout session
type StripeCheckoutResponse struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

// handleStripeCheckout creates a Stripe checkout session
func (s *CloudServer) handleStripeCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req StripeCheckoutSession
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	if req.Quantity < 1 {
		req.Quantity = 1
	}

	// Get Stripe config from environment
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	priceID := os.Getenv("STRIPE_PRICE_ID")

	if stripeKey == "" || priceID == "" {
		// Return placeholder response for demo
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(StripeCheckoutResponse{
			SessionID: "demo_session",
			URL:       "https://buy.stripe.com/demo",
		})
		return
	}

	// In production, you would call the Stripe API here
	// For now, return a placeholder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StripeCheckoutResponse{
		SessionID: "placeholder",
		URL:       fmt.Sprintf("https://checkout.stripe.com/pay/%s", priceID),
	})
}

// handleStripeWebhook handles Stripe webhook events
func (s *CloudServer) handleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// In production, you would verify the webhook signature
	// and handle events like checkout.session.completed

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// RegisterStripeRoutes registers Stripe-related routes
func (s *CloudServer) RegisterStripeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/cloud/checkout", s.handleStripeCheckout)
	mux.HandleFunc("/api/cloud/webhook", s.handleStripeWebhook)
}
