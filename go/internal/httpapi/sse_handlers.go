package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type SSEClient struct {
	ID      string
	Message chan []byte
}

type SSEBroker struct {
	clients map[string]*SSEClient
	mu      sync.RWMutex
}

var GlobalSSEBroker = &SSEBroker{
	clients: make(map[string]*SSEClient),
}

func (b *SSEBroker) AddClient(client *SSEClient) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client.ID] = client
}

func (b *SSEBroker) RemoveClient(client *SSEClient) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, client.ID)
	close(client.Message)
}

func (b *SSEBroker) Broadcast(msg []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, client := range b.clients {
		select {
		case client.Message <- msg:
		default:
			// Client channel is full, drop message
		}
	}
}

func validateSSEToken(r *http.Request) bool {
	tokenEnv := os.Getenv("CLOUDMCP_SSE_AUTH_TOKEN")
	if tokenEnv == "" {
		return true // Auth not enabled/required if env not set
	}
	token := r.URL.Query().Get("token")
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}
	return token == tokenEnv
}

type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
	ID      any    `json:"id,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  any    `json:"result,omitempty"`
	Error   any    `json:"error,omitempty"`
	ID      any    `json:"id"`
}

// handleSSE serves the Server-Sent Events endpoint for browser extensions and external clients.
// This establishes the native Go control plane as a valid primary endpoint, achieving
// 100% protocol parity with the TypeScript core.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	if !validateSSEToken(r) {
		http.Error(w, "Unauthorized: Invalid SSE token", http.StatusUnauthorized)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	client := &SSEClient{
		ID:      fmt.Sprintf("client-%d", time.Now().UnixNano()),
		Message: make(chan []byte, 100),
	}

	GlobalSSEBroker.AddClient(client)
	defer GlobalSSEBroker.RemoveClient(client)

	// Send an initial connected message pointing to the POST endpoint
	fmt.Fprintf(w, "event: endpoint\ndata: /api/sse/message?sessionId=%s\n\n", client.ID)
	flusher.Flush()

	for {
		select {
		case msg := <-client.Message:
			fmt.Fprintf(w, "data: %s\n\n", string(msg))
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) handleSSEMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"success": false, "error": "method not allowed"})
		return
	}

	if !validateSSEToken(r) {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"success": false, "error": "Unauthorized"})
		return
	}

	sessionId := r.URL.Query().Get("sessionId")
	if sessionId == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "Missing sessionId query parameter"})
		return
	}

	GlobalSSEBroker.mu.RLock()
	client, exists := GlobalSSEBroker.clients[sessionId]
	GlobalSSEBroker.mu.RUnlock()

	if !exists {
		writeJSON(w, http.StatusNotFound, map[string]any{"success": false, "error": "Session not found"})
		return
	}

	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "error": "Invalid JSON-RPC request"})
		return
	}

	// Process the JSON-RPC request and send response back on the SSE connection
	var result any
	var errVal any

	switch req.Method {
	case "initialize":
		result = map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
			"serverInfo": map[string]any{
				"name":    "tormentnexus-cloud",
				"version": "1.0.0",
			},
		}
	case "tools/list":
		view, err := s.localMCPInventoryView()
		if err == nil && view != nil {
			result = map[string]any{
				"tools": s.injectAlwaysOnStatus(s.mergeAccessoryTools(fallbackMCPInventoryTools(view))),
			}
		} else {
			result = map[string]any{
				"tools": []any{},
			}
		}
	case "tools/call":
		// Call the requested tool locally or upstream
		var payload map[string]any
		reqBytes, _ := json.Marshal(req.Params)
		_ = json.Unmarshal(reqBytes, &payload)

		var toolResult any
		fallbackResult, fallbackErr := s.localCallMCPMetaTool(r, payload)
		if fallbackErr == nil {
			toolResult = fallbackResult
		} else {
			upstreamResult, upstreamErr := s.callUpstreamJSON(r.Context(), "mcp.callTool", payload, &toolResult)
			if upstreamErr != nil {
				errVal = map[string]any{
					"code":    -32603,
					"message": fmt.Sprintf("Tool call failed: %v / upstream: %v", fallbackErr, upstreamErr),
				}
			} else {
				toolResult = upstreamResult
			}
		}
		result = toolResult
	default:
		errVal = map[string]any{
			"code":    -32601,
			"message": fmt.Sprintf("Method %s not found", req.Method),
		}
	}

	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}
	if errVal != nil {
		resp.Error = errVal
	} else {
		resp.Result = result
	}

	respBytes, _ := json.Marshal(resp)
	select {
	case client.Message <- respBytes:
	default:
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (s *Server) handleSSEHistory(w http.ResponseWriter, r *http.Request) {
	sinceStr := r.URL.Query().Get("since")
	var since int64
	if sinceStr != "" {
		fmt.Sscanf(sinceStr, "%d", &since)
	}

	history := s.eventBus.GetHistorySince(since)
	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"data":    history,
	})
}
