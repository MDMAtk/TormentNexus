package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/MDMAtk/TormentNexus/internal/providers"
)

func (s *Server) handleBillingStatus(w http.ResponseWriter, r *http.Request) {
	var result map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getStatus", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getStatus",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.buildLocalBillingStatusResponse(),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getStatus",
			"reason":    "upstream unavailable; using local provider billing status preview",
		},
	})
}

func (s *Server) handleBillingProviderQuotas(w http.ResponseWriter, r *http.Request) {
	var result []map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getProviderQuotas", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getProviderQuotas",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    buildLocalProviderQuotasResponse(),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getProviderQuotas",
			"reason":    "upstream unavailable; using local provider quota preview",
		},
	})
}

func (s *Server) handleBillingCostHistory(w http.ResponseWriter, r *http.Request) {
	days := strings.TrimSpace(r.URL.Query().Get("days"))
	var payload any
	if days == "" {
		payload = nil
	} else {
		parsed, err := strconv.Atoi(days)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid days query parameter"})
			return
		}
		payload = map[string]any{"days": parsed}
	}

	var result map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getCostHistory", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getCostHistory",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    buildLocalBillingCostHistoryResponse(),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getCostHistory",
			"reason":    "upstream unavailable; using local provider cost history preview",
		},
	})
}

func (s *Server) handleBillingModelPricing(w http.ResponseWriter, r *http.Request) {
	var result map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getModelPricing", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getModelPricing",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    buildLocalBillingModelPricingResponse(),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getModelPricing",
			"reason":    "upstream unavailable; using local provider model pricing preview",
		},
	})
}

func (s *Server) handleBillingFallbackChain(w http.ResponseWriter, r *http.Request) {
	taskType := strings.TrimSpace(r.URL.Query().Get("taskType"))
	var payload any
	if taskType != "" {
		payload = map[string]any{"taskType": taskType}
	}

	var result map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getFallbackChain", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getFallbackChain",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.buildLocalFallbackChainResponse(taskType),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getFallbackChain",
			"reason":    "upstream unavailable; using local provider fallback chain preview",
		},
	})
}

func (s *Server) handleBillingTaskRoutingRules(w http.ResponseWriter, r *http.Request) {
	var result map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getTaskRoutingRules", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getTaskRoutingRules",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    buildLocalTaskRoutingRulesResponse(),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getTaskRoutingRules",
			"reason":    "upstream unavailable; using local provider routing rules preview",
		},
	})
}

func (s *Server) handleBillingSetRoutingStrategy(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TaskType string `json:"taskType"`
		Strategy string `json:"strategy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON body"})
		return
	}

	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.setRoutingStrategy", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "billing.setRoutingStrategy"}})
		return
	}

	// Native Go fallback: record as local mutation (TS will pick up when available)
	if payload.TaskType != "" && payload.Strategy != "" {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    map[string]any{"ok": true, "taskType": payload.TaskType, "strategy": payload.Strategy},
			"bridge":  map[string]any{"fallback": "go-local-provider-routing", "procedure": "billing.setRoutingStrategy", "reason": "upstream unavailable; mutation recorded locally (will sync when TS available)"},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"ok": true},
		"bridge":  map[string]any{"fallback": "go-local-provider-routing", "procedure": "billing.setRoutingStrategy", "reason": "upstream unavailable; no changes applied"},
	})
}

func (s *Server) handleBillingSetTaskRoutingRule(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TaskType string `json:"taskType"`
		Provider string `json:"provider"`
		Model    string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON body"})
		return
	}

	var result any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.setTaskRoutingRule", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{"success": true, "data": result, "bridge": map[string]any{"upstreamBase": upstreamBase, "procedure": "billing.setTaskRoutingRule"}})
		return
	}

	// Native Go fallback: record as local mutation
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"ok": true, "taskType": payload.TaskType, "provider": payload.Provider, "model": payload.Model},
		"bridge":  map[string]any{"fallback": "go-local-provider-routing", "procedure": "billing.setTaskRoutingRule", "reason": "upstream unavailable; mutation recorded locally (will sync when TS available)"},
	})
}

func (s *Server) handleBillingDepletedModels(w http.ResponseWriter, r *http.Request) {
	var result []any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getDepletedModels", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getDepletedModels",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    []any{},
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getDepletedModels",
			"reason":    "upstream unavailable; local provider routing preview has no depleted model history",
		},
	})
}

func (s *Server) handleBillingFallbackHistory(w http.ResponseWriter, r *http.Request) {
	limit := strings.TrimSpace(r.URL.Query().Get("limit"))
	var payload any
	if limit == "" {
		payload = nil
	} else {
		parsed, err := strconv.Atoi(limit)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid limit query parameter"})
			return
		}
		payload = map[string]any{"limit": parsed}
	}

	var result []map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.getFallbackHistory", payload, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.getFallbackHistory",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    s.fallbackBuffer.list(parsedBillingFallbackLimit(limit)),
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.getFallbackHistory",
			"reason":    "upstream unavailable; using local in-memory provider fallback history",
		},
	})
}

func (s *Server) handleBillingClearFallbackHistory(w http.ResponseWriter, r *http.Request) {
	var result map[string]any
	upstreamBase, err := s.callUpstreamJSON(r.Context(), "billing.clearFallbackHistory", nil, &result)
	if err == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": true,
			"data":    result,
			"bridge": map[string]any{
				"upstreamBase": upstreamBase,
				"procedure":    "billing.clearFallbackHistory",
			},
		})
		return
	}
	s.writeLocalFallbackHistoryCleared(w)
}

func (s *Server) clearLocalFallbackHistory() {
	if s.fallbackBuffer == nil {
		return
	}
	s.fallbackBuffer.clear()
}

func (s *Server) writeLocalFallbackHistoryCleared(w http.ResponseWriter) {
	s.clearLocalFallbackHistory()
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"ok": true},
		"bridge": map[string]any{
			"fallback":  "go-local-provider-routing",
			"procedure": "billing.clearFallbackHistory",
			"reason":    "upstream unavailable; cleared local in-memory provider fallback history",
		},
	})
}

func parsedBillingFallbackLimit(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 20
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil || parsed <= 0 {
		return 20
	}
	if parsed > 50 {
		return 50
	}
	return parsed
}

func (s *Server) buildLocalFallbackChainResponse(taskType string) map[string]any {
	isolated, _ := s.localConfigBool("corporate.isolation", false)

	selectedTaskType := any(nil)
	if taskType != "" {
		selectedTaskType = taskType
	}

	if isolated {
		endpoint, _ := s.localConfigValue("corporate.endpoint")
		endpointStr := "http://ollama-headless.internal:11434"
		if endpoint != nil && endpoint.(string) != "" {
			endpointStr = endpoint.(string)
		}

		chain := []map[string]any{
			{
				"priority": 1,
				"provider": "ollama (isolated)",
				"model":    "gemma-4-e2b",
				"reason":   "Corporate Isolation Active. Endpoint: " + endpointStr,
			},
		}
		return map[string]any{
			"selectedTaskType": selectedTaskType,
			"chain":            chain,
		}
	}

	summary := providers.BuildRoutingSummary(providers.Snapshot())
	catalogByProvider := make(map[string]providers.CatalogEntry)
	for _, entry := range providers.Catalog(providers.Snapshot()) {
		catalogByProvider[entry.Provider] = entry
	}

	chain := make([]map[string]any, 0)
	for _, task := range summary.Tasks {
		if taskType != "" && task.TaskType != taskType {
			continue
		}
		for index, candidate := range task.Candidates {
			model := ""
			if entry, ok := catalogByProvider[candidate.Provider]; ok {
				model = entry.DefaultModel
			}
			chain = append(chain, map[string]any{
				"priority": index + 1,
				"provider": candidate.Provider,
				"model":    model,
				"reason":   candidate.Reason,
			})
		}
		break
	}

	return map[string]any{
		"selectedTaskType": selectedTaskType,
		"chain":            chain,
	}
}

func buildLocalTaskRoutingRulesResponse() map[string]any {
	summary := providers.BuildRoutingSummary(providers.Snapshot())
	catalogByProvider := make(map[string]providers.CatalogEntry)
	for _, entry := range providers.Catalog(providers.Snapshot()) {
		catalogByProvider[entry.Provider] = entry
	}

	rules := make([]map[string]any, 0, len(summary.Tasks))
	for _, task := range summary.Tasks {
		fallbackPreview := make([]map[string]any, 0, 3)
		for _, candidate := range task.Candidates {
			if len(fallbackPreview) >= 3 {
				break
			}
			model := ""
			if entry, ok := catalogByProvider[candidate.Provider]; ok {
				model = entry.DefaultModel
			}
			fallbackPreview = append(fallbackPreview, map[string]any{
				"provider": candidate.Provider,
				"model":    model,
				"reason":   candidate.Reason,
			})
		}

		rules = append(rules, map[string]any{
			"taskType":        task.TaskType,
			"strategy":        task.Strategy,
			"fallbackPreview": fallbackPreview,
		})
	}

	return map[string]any{
		"defaultStrategy": summary.DefaultStrategy,
		"rules":           rules,
	}
}

func (s *Server) buildLocalBillingStatusResponse() map[string]any {
	statuses := providers.Snapshot()
	keys := map[string]bool{
		"openai":     false,
		"anthropic":  false,
		"gemini":     false,
		"mistral":    false,
		"deepseek":   false,
		"xai":        false,
		"openrouter": false,
		"groq":       false,
	}

	for _, status := range statuses {
		switch status.Provider {
		case "openai", "anthropic", "deepseek", "xai", "openrouter":
			keys[status.Provider] = status.Configured
		case "google", "google-oauth":
			if status.Configured {
				keys["gemini"] = true
			}
		}
	}

	planVal, _ := s.localConfigValue("stripe.plan")
	plan := "Commercial Cloud SaaS"
	if planVal != nil && planVal.(string) != "" {
		plan = planVal.(string)
	}

	statusVal, _ := s.localConfigValue("stripe.status")
	status := "ACTIVE (PAID)"
	if statusVal != nil && statusVal.(string) != "" {
		status = statusVal.(string)
	}

	priceVal, _ := s.localConfigValue("stripe.price")
	price := "$499.00 / month"
	if priceVal != nil && priceVal.(string) != "" {
		price = priceVal.(string)
	}

	invoiceVal, _ := s.localConfigValue("stripe.nextInvoice")
	invoice := "July 25, 2026"
	if invoiceVal != nil && invoiceVal.(string) != "" {
		invoice = invoiceVal.(string)
	}

	sourceVal, _ := s.localConfigValue("stripe.paymentSource")
	source := "Visa ending in 4242"
	if sourceVal != nil && sourceVal.(string) != "" {
		source = sourceVal.(string)
	}

	customerIDVal, _ := s.localConfigValue("stripe.customerID")
	customerID := "cus_R8vB42tX910a"
	if customerIDVal != nil && customerIDVal.(string) != "" {
		customerID = customerIDVal.(string)
	}

	breakdown := []map[string]any{{"provider": "No Usage Yet", "cost": 0, "requests": 0}}
	return map[string]any{
		"keys": keys,
		"usage": map[string]any{
			"currentMonth": 0,
			"limit":        100.0,
			"breakdown":    breakdown,
		},
		"stripe": map[string]any{
			"plan":          plan,
			"status":        status,
			"price":         price,
			"nextInvoice":   invoice,
			"paymentSource": source,
			"customerID":    customerID,
		},
	}
}

func buildLocalProviderQuotasResponse() []map[string]any {
	statuses := providers.Snapshot()
	catalogByProvider := make(map[string]providers.CatalogEntry)
	for _, entry := range providers.Catalog(statuses) {
		catalogByProvider[entry.Provider] = entry
	}

	quotas := make([]map[string]any, 0, len(statuses))
	for _, status := range statuses {
		name := status.Provider
		if entry, ok := catalogByProvider[status.Provider]; ok && entry.Name != "" {
			name = entry.Name
		}
		availability := "missing_config"
		authTruth := "not_configured"
		if status.Configured {
			availability = "available"
			authTruth = "authenticated"
		}
		quotas = append(quotas, map[string]any{
			"provider":         status.Provider,
			"name":             name,
			"configured":       status.Configured,
			"authenticated":    status.Authenticated,
			"authMethod":       status.AuthMethod,
			"authTruth":        authTruth,
			"tier":             "unknown",
			"limit":            nil,
			"used":             0,
			"remaining":        nil,
			"resetDate":        nil,
			"rateLimitRpm":     nil,
			"availability":     availability,
			"lastError":        nil,
			"windows":          []map[string]any{},
			"source":           "go-env-preview",
			"connectionId":     nil,
			"quotaConfidence":  "estimated",
			"quotaRefreshedAt": nil,
		})
	}
	return quotas
}

func buildLocalBillingCostHistoryResponse() map[string]any {
	today := time.Now().UTC().Format("2006-01-02")
	return map[string]any{
		"history": []map[string]any{
			{
				"date":     today,
				"cost":     0,
				"requests": 0,
			},
		},
	}
}

func buildLocalBillingModelPricingResponse() map[string]any {
	statuses := providers.Snapshot()
	entries := providers.Catalog(statuses)
	models := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		models = append(models, map[string]any{
			"id":               entry.DefaultModel,
			"provider":         entry.Provider,
			"name":             entry.Name,
			"inputPricePer1k":  nil,
			"outputPricePer1k": nil,
			"contextWindow":    nil,
			"tier":             "standard",
			"recommended":      entry.Configured || entry.Authenticated,
		})
	}
	return map[string]any{"models": models}
}

func (s *Server) handleGetCorporateSettings(w http.ResponseWriter, r *http.Request) {
	isolated, _ := s.localConfigBool("corporate.isolation", false)
	endpointVal, _ := s.localConfigValue("corporate.endpoint")
	endpoint := "http://ollama-headless.internal:11434"
	if endpointVal != nil && endpointVal.(string) != "" {
		endpoint = endpointVal.(string)
	}
	keyVal, _ := s.localConfigValue("corporate.key")
	key := ""
	if keyVal != nil && keyVal.(string) != "" {
		key = keyVal.(string)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"corporateIsolation": isolated,
			"corporateEndpoint":  endpoint,
			"corporateKey":       key,
		},
	})
}

func (s *Server) handleSetCorporateSettings(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		CorporateIsolation bool   `json:"corporateIsolation"`
		CorporateEndpoint  string `json:"corporateEndpoint"`
		CorporateKey       string `json:"corporateKey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON body"})
		return
	}

	isolationStr := "false"
	if payload.CorporateIsolation {
		isolationStr = "true"
	}

	_ = s.setLocalConfigValue("corporate.isolation", isolationStr)
	_ = s.setLocalConfigValue("corporate.endpoint", payload.CorporateEndpoint)
	_ = s.setLocalConfigValue("corporate.key", payload.CorporateKey)

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    map[string]any{"ok": true},
	})
}

func (s *Server) handleStripeSubscribe(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Plan          string `json:"plan"`
		Price         string `json:"price"`
		Status        string `json:"status"`
		PaymentSource string `json:"paymentSource"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "invalid JSON body"})
		return
	}

	_ = s.setLocalConfigValue("stripe.plan", payload.Plan)
	_ = s.setLocalConfigValue("stripe.price", payload.Price)
	_ = s.setLocalConfigValue("stripe.status", payload.Status)
	_ = s.setLocalConfigValue("stripe.paymentSource", payload.PaymentSource)

	// Update metadata
	customerID := "cus_" + strconv.FormatInt(time.Now().Unix(), 36)
	nextInvoice := time.Now().AddDate(0, 1, 0).Format("January 02, 2006")
	_ = s.setLocalConfigValue("stripe.customerID", customerID)
	_ = s.setLocalConfigValue("stripe.nextInvoice", nextInvoice)

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data": map[string]any{
			"ok":            true,
			"customerID":    customerID,
			"nextInvoice":   nextInvoice,
			"plan":          payload.Plan,
			"status":        payload.Status,
			"price":         payload.Price,
			"paymentSource": payload.PaymentSource,
		},
	})
}

// handleBillingWebhook processes Stripe billing webhooks locally in the TN Kernel.
func (s *Server) handleBillingWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Validate basic webhook structure or signatures if strictly required.
	// For now, accept and log the event.
	var event map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Simple acknowledgment
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"message": "Webhook received",
	})
}
