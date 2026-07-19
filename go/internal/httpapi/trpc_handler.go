package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type responseRecorder struct {
	header http.Header
	body   *bytes.Buffer
	status int
}

func newResponseRecorder() *responseRecorder {
	return &responseRecorder{
		header: make(http.Header),
		body:   new(bytes.Buffer),
		status: http.StatusOK,
	}
}

func (r *responseRecorder) Header() http.Header {
	return r.header
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
}

type cacheEntry struct {
	data      any
	createdAt time.Time
	ttl       time.Duration
}

var (
	trpcCache           = make(map[string]cacheEntry)
	trpcCacheMu         sync.RWMutex
	cacheableProcedures = map[string]time.Duration{
		"startupStatus":               10 * time.Second,
		"mcp.traffic":                 3 * time.Second,
		"session.list":                3 * time.Second,
		"billing.getCostHistory":      10 * time.Second,
		"billing.getModelPricing":     10 * time.Second,
		"billing.getTaskRoutingRules": 10 * time.Second,
		"mcp.getWorkingSet":           10 * time.Second,
		"mcp.getJsoncEditor":          10 * time.Second,
		"mcp.searchTools":             5 * time.Second,
	}
)

func getCompatRoute(procedurePath string, input map[string]any) string {
	switch procedurePath {
	case "startupStatus":
		return "/api/startup/status"
	case "mcp.getStatus":
		return "/api/mcp/status"
	case "mcp.listServers":
		return "/api/mcp/servers"
	case "mcp.getToolSelectionTelemetry":
		return "/api/mcp/tool-selection-telemetry"
	case "mcp.clearToolSelectionTelemetry":
		return "/api/mcp/tool-selection-telemetry/clear"
	case "mcp.runServerTest":
		return "/api/mcp/server-test"
	case "session.list":
		return "/api/native/session/list"
	case "session.importedMaintenanceStats":
		return "/api/sessions/imported/maintenance-stats"
	case "billing.getStatus":
		return "/api/billing/status"
	case "billing.getProviderQuotas":
		return "/api/billing/provider-quotas"
	case "billing.getCostHistory":
		days := 30
		if input != nil {
			if dVal, ok := input["days"]; ok {
				if dNum, ok := dVal.(float64); ok {
					days = int(dNum)
				}
			}
		}
		return fmt.Sprintf("/api/billing/cost-history?days=%d", days)
	case "billing.getModelPricing":
		return "/api/billing/model-pricing"
	case "billing.getFallbackChain":
		taskType := ""
		if input != nil {
			if tVal, ok := input["taskType"]; ok {
				if tStr, ok := tVal.(string); ok {
					taskType = tStr
				}
			}
		}
		if taskType != "" {
			return fmt.Sprintf("/api/billing/fallback-chain?taskType=%s", url.QueryEscape(taskType))
		}
		return "/api/billing/fallback-chain"
	case "billing.getTaskRoutingRules":
		return "/api/billing/task-routing-rules"
	case "billing.getDepletedModels":
		return "/api/billing/depleted-models"
	case "billing.getFallbackHistory":
		limit := 20
		if input != nil {
			if lVal, ok := input["limit"]; ok {
				if lNum, ok := lVal.(float64); ok {
					limit = int(lNum)
				}
			}
		}
		return fmt.Sprintf("/api/billing/fallback-history?limit=%d", limit)
	case "director.status":
		return "/api/director/status"
	case "directorConfig.get":
		return "/api/director-config"
	case "llm.generate":
		return "/api/llm/generate"
	case "memory.getRecentObservations":
		limit := 6
		namespace := ""
		mType := ""
		if input != nil {
			if lVal, ok := input["limit"]; ok {
				if lNum, ok := lVal.(float64); ok {
					limit = int(lNum)
				}
			}
			if nsVal, ok := input["namespace"]; ok {
				if nsStr, ok := nsVal.(string); ok {
					namespace = nsStr
				}
			}
			if tVal, ok := input["type"]; ok {
				if tStr, ok := tVal.(string); ok {
					mType = tStr
				}
			}
		}
		q := url.Values{}
		q.Set("limit", strconv.Itoa(limit))
		if namespace != "" {
			q.Set("namespace", namespace)
		}
		if mType != "" {
			q.Set("type", mType)
		}
		return "/api/memory/observations/recent?" + q.Encode()
	case "memory.getRecentUserPrompts":
		limit := 5
		role := ""
		if input != nil {
			if lVal, ok := input["limit"]; ok {
				if lNum, ok := lVal.(float64); ok {
					limit = int(lNum)
				}
			}
			if rVal, ok := input["role"]; ok {
				if rStr, ok := rVal.(string); ok {
					role = rStr
				}
			}
		}
		q := url.Values{}
		q.Set("limit", strconv.Itoa(limit))
		if role != "" {
			q.Set("role", role)
		}
		return "/api/memory/user-prompts/recent?" + q.Encode()
	case "memory.getRecentSessionSummaries":
		limit := 4
		if input != nil {
			if lVal, ok := input["limit"]; ok {
				if lNum, ok := lVal.(float64); ok {
					limit = int(lNum)
				}
			}
		}
		return fmt.Sprintf("/api/memory/session-summaries/recent?limit=%d", limit)
	case "billing.stripe.plans":
		return "/api/billing/stripe/plans"
	case "billing.stripe.checkout":
		return "/api/billing/stripe/checkout"
	case "billing.stripe.portal":
		return "/api/billing/stripe/portal"
	case "billing.stripe.subscription":
		return "/api/billing/stripe/subscription"
	case "health":
		return "/api/health"
	case "git.getLog":
		limit := 10
		if input != nil {
			if lVal, ok := input["limit"]; ok {
				if lNum, ok := lVal.(float64); ok {
					limit = int(lNum)
				}
			}
		}
		return fmt.Sprintf("/api/git/log?limit=%d", limit)
	case "git.getStatus":
		return "/api/git/status"
	case "git.getModules":
		return "/api/git/modules"
	case "git.revert":
		return "/api/git/revert"
	case "graph.getSymbolsGraph", "graph.getSymbols":
		return "/api/graph/symbols"
	case "graph.get":
		return "/api/graph"
	case "knowledge.ingest":
		return "/api/knowledge/ingest"
	case "knowledge.getResources":
		return "/api/knowledge/resources"
	case "knowledge.graph":
		return "/api/knowledge/graph"
	case "knowledge.stats":
		return "/api/knowledge/stats"
	case "lsp.getSymbols":
		filePath := ""
		if input != nil {
			if fpVal, ok := input["filePath"]; ok {
				if fpStr, ok := fpVal.(string); ok {
					filePath = fpStr
				}
			}
		}
		return "/api/lsp/symbols?filePath=" + url.QueryEscape(filePath)
	case "lsp.findSymbol":
		filePath := ""
		symbolName := ""
		if input != nil {
			if fpVal, ok := input["filePath"]; ok {
				if fpStr, ok := fpVal.(string); ok {
					filePath = fpStr
				}
			}
			if snVal, ok := input["symbolName"]; ok {
				if snStr, ok := snVal.(string); ok {
					symbolName = snStr
				}
			}
		}
		return fmt.Sprintf("/api/lsp/find-symbol?filePath=%s&symbolName=%s", url.QueryEscape(filePath), url.QueryEscape(symbolName))
	case "lsp.findReferences":
		filePath := ""
		line := 0
		char := 0
		if input != nil {
			if fpVal, ok := input["filePath"]; ok {
				if fpStr, ok := fpVal.(string); ok {
					filePath = fpStr
				}
			}
			if lVal, ok := input["line"]; ok {
				if lNum, ok := lVal.(float64); ok {
					line = int(lNum)
				}
			}
			if cVal, ok := input["character"]; ok {
				if cNum, ok := cVal.(float64); ok {
					char = int(cNum)
				}
			}
		}
		return fmt.Sprintf("/api/lsp/find-references?filePath=%s&line=%d&character=%d", url.QueryEscape(filePath), line, char)
	case "lsp.searchSymbols":
		query := ""
		if input != nil {
			if qVal, ok := input["query"]; ok {
				if qStr, ok := qVal.(string); ok {
					query = qStr
				}
			}
		}
		return "/api/lsp/search?query=" + url.QueryEscape(query)
	case "lsp.indexProject":
		return "/api/lsp/index"
	case "memory.query", "memory.searchAgentMemory":
		query := ""
		limit := 10
		if input != nil {
			if qVal, ok := input["query"]; ok {
				if qStr, ok := qVal.(string); ok {
					query = qStr
				}
			}
			if lVal, ok := input["limit"]; ok {
				if lNum, ok := lVal.(float64); ok {
					limit = int(lNum)
				}
			}
		}
		return fmt.Sprintf("/api/memory/search?query=%s&limit=%d", url.QueryEscape(query), limit)
	case "memory.searchObservations":
		return "/api/memory/observations/search"
	case "memory.searchUserPrompts":
		return "/api/memory/user-prompts/search"
	case "memory.searchMemoryPivot":
		return "/api/memory/pivot/search"
	case "memory.searchSessionSummaries":
		return "/api/memory/session-summaries/search"
	}
	return ""
}

func (s *Server) handleTRPC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, trpc-accept")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	procedurePath := strings.TrimPrefix(r.URL.Path, "/trpc/")
	if procedurePath == "" {
		http.Error(w, "missing procedure name", http.StatusBadRequest)
		return
	}

	isBatch := strings.Contains(procedurePath, ",") || r.URL.Query().Get("batch") == "1"
	var procedures []string
	if isBatch {
		procedures = strings.Split(procedurePath, ",")
	} else {
		procedures = []string{procedurePath}
	}

	// Parse input parameters
	var inputs map[string]any
	var inputStr string
	if r.Method == http.MethodGet {
		inputStr = r.URL.Query().Get("input")
		if inputStr != "" {
			_ = json.Unmarshal([]byte(inputStr), &inputs)
		}
	} else if r.Method == http.MethodPost {
		bodyBytes, _ := io.ReadAll(r.Body)
		if len(bodyBytes) > 0 {
			inputStr = string(bodyBytes)
			_ = json.Unmarshal(bodyBytes, &inputs)
		}
	}

	// Prepare result wrapper
	var batchResults []any

	for idx, proc := range procedures {
		proc = strings.TrimSpace(proc)
		var singleInput map[string]any

		if isBatch {
			if inputs != nil {
				if keyVal, ok := inputs[strconv.Itoa(idx)]; ok {
					if keyMap, ok := keyVal.(map[string]any); ok {
						if jsonVal, exists := keyMap["json"]; exists {
							if jsonMap, ok := jsonVal.(map[string]any); ok {
								singleInput = jsonMap
							} else {
								singleInput = keyMap
							}
						} else {
							singleInput = keyMap
						}
					}
				}
			}
		} else {
			if inputs != nil {
				if jsonVal, exists := inputs["json"]; exists {
					if jsonMap, ok := jsonVal.(map[string]any); ok {
						singleInput = jsonMap
					} else {
						singleInput = inputs
					}
				} else {
					singleInput = inputs
				}
			}
		}

		// Low-Latency In-Memory Caching (GET requests only)
		var cacheKey string
		var ttl time.Duration
		var hasCache bool

		if r.Method == http.MethodGet {
			if t, ok := cacheableProcedures[proc]; ok {
				ttl = t
				singleInputStr := ""
				if isBatch {
					singleInputStr = strconv.Itoa(idx)
				} else {
					singleInputStr = inputStr
				}
				cacheKey = fmt.Sprintf("%s:%s", proc, singleInputStr)

				trpcCacheMu.RLock()
				entry, ok := trpcCache[cacheKey]
				trpcCacheMu.RUnlock()

				if ok && time.Since(entry.createdAt) < entry.ttl {
					batchResults = append(batchResults, entry.data)
					hasCache = true
				}
			}
		}

		if hasCache {
			continue
		}

		compatRoute := getCompatRoute(proc, singleInput)
		if compatRoute == "" {
			batchResults = append(batchResults, map[string]any{
				"error": map[string]any{
					"message": fmt.Sprintf("Procedure %s not supported natively", proc),
					"code":    -32601,
				},
			})
			continue
		}

		targetURL, err := url.Parse(compatRoute)
		if err != nil {
			batchResults = append(batchResults, map[string]any{
				"error": map[string]any{"message": err.Error(), "code": -32603},
			})
			continue
		}

		internalQuery := targetURL.Query()
		var internalBody io.Reader
		if r.Method == http.MethodPost && len(singleInput) > 0 {
			bodyBytes, _ := json.Marshal(singleInput)
			internalBody = bytes.NewReader(bodyBytes)
		}

		internalReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL.Path, internalBody)
		if err != nil {
			batchResults = append(batchResults, map[string]any{
				"error": map[string]any{"message": err.Error(), "code": -32603},
			})
			continue
		}

		internalReq.URL.RawQuery = internalQuery.Encode()
		internalReq.Header = r.Header.Clone()

		rec := newResponseRecorder()
		s.mux.ServeHTTP(rec, internalReq)

		var payloadData any
		bodyBytes := rec.body.Bytes()

		if rec.status >= 400 {
			batchResults = append(batchResults, map[string]any{
				"error": map[string]any{
					"message": string(bodyBytes),
					"code":    -32603,
				},
			})
			continue
		}

		if len(bodyBytes) > 0 {
			var parsedJSON any
			if err := json.Unmarshal(bodyBytes, &parsedJSON); err == nil {
				if m, ok := parsedJSON.(map[string]any); ok {
					if dVal, ok := m["data"]; ok {
						payloadData = dVal
					} else {
						payloadData = parsedJSON
					}
				} else {
					payloadData = parsedJSON
				}
			} else {
				payloadData = string(bodyBytes)
			}
		}

		resultWrapper := map[string]any{
			"result": map[string]any{
				"data": payloadData,
			},
		}

		// Store in cache
		if cacheKey != "" {
			trpcCacheMu.Lock()
			trpcCache[cacheKey] = cacheEntry{
				data:      resultWrapper,
				createdAt: time.Now(),
				ttl:       ttl,
			}
			trpcCacheMu.Unlock()
		}

		batchResults = append(batchResults, resultWrapper)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if isBatch {
		_ = json.NewEncoder(w).Encode(batchResults)
	} else {
		if len(batchResults) > 0 {
			_ = json.NewEncoder(w).Encode(batchResults[0])
		} else {
			_ = json.NewEncoder(w).Encode(map[string]any{})
		}
	}
}
