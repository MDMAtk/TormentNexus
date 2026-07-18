package cloud

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// MCPStreamableHTTP handles Streamable HTTP transport for MCP
type MCPStreamableHTTP struct {
	am          *AccountManager
	provisioner *Provisioner
	sessions    map[string]*MCPSession
	mu          sync.RWMutex
}

// MCPSession represents an MCP session
type MCPSession struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	ContainerID string    `json:"container_id"`
	CreatedAt   time.Time `json:"created_at"`
	LastPing    time.Time `json:"last_ping"`
}

// MCPRequest represents an MCP JSON-RPC request
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents an MCP JSON-RPC response
type MCPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// NewMCPStreamableHTTP creates a new MCP Streamable HTTP handler
func NewMCPStreamableHTTP(am *AccountManager, provisioner *Provisioner) *MCPStreamableHTTP {
	return &MCPStreamableHTTP{
		am:          am,
		provisioner: provisioner,
		sessions:    make(map[string]*MCPSession),
	}
}

// RegisterRoutes registers the MCP HTTP routes
func (h *MCPStreamableHTTP) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/mcp/v1", h.handleMCP)
	mux.HandleFunc("/mcp/v1/sse", h.handleMCPSSE)
	mux.HandleFunc("/mcp/v1/health", h.handleHealth)
}

// handleMCP handles MCP requests via Streamable HTTP
func (h *MCPStreamableHTTP) handleMCP(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	account, err := h.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Parse request
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Try parsing as SSE message
		h.handleSSEMessage(w, r, account)
		return
	}

	// Handle request
	resp := h.handleRequest(r.Context(), account, &req)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleMCPSSE handles SSE connections for MCP
func (h *MCPStreamableHTTP) handleMCPSSE(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	account, err := h.authenticate(r)
	if err != nil {
		http.Error(w, `{"error": "unauthorized"}`, http.StatusUnauthorized)
		return
	}

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	origin := r.Header.Get("Origin")
	if origin == "https://cloud.hypernexus.site" || origin == "https://hypernexus.site" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// Create session
	session := &MCPSession{
		ID:        generateID(),
		AccountID: account.ID,
		CreatedAt: time.Now(),
		LastPing:  time.Now(),
	}

	// Get container
	container, err := h.provisioner.GetContainerByAccount(account.ID)
	if err != nil {
		// Provision container if needed
		container, err = h.provisioner.Provision(r.Context(), account.ID, string(account.Plan))
		if err != nil {
			http.Error(w, `{"error": "failed to provision container"}`, http.StatusInternalServerError)
			return
		}
	}
	session.ContainerID = container.ID

	// Store session
	h.mu.Lock()
	h.sessions[session.ID] = session
	h.mu.Unlock()

	// Send session ID
	fmt.Fprintf(w, "data: {\"session_id\": \"%s\"}\n\n", session.ID)
	w.(http.Flusher).Flush()

	// Handle SSE messages
	ctx := r.Context()
	for {
		select {
		case <-ctx.Done():
			h.mu.Lock()
			delete(h.sessions, session.ID)
			h.mu.Unlock()
			return
		default:
			// Keep connection alive
			time.Sleep(30 * time.Second)
			fmt.Fprintf(w, "data: {\"type\": \"ping\"}\n\n")
			w.(http.Flusher).Flush()
		}
	}
}

// handleSSEMessage handles an SSE message
func (h *MCPStreamableHTTP) handleSSEMessage(w http.ResponseWriter, r *http.Request, account *Account) {
	// Read the message
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error": "failed to read body"}`, http.StatusBadRequest)
		return
	}

	// Parse SSE format
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			var req MCPRequest
			if err := json.Unmarshal([]byte(data), &req); err != nil {
				continue
			}
			resp := h.handleRequest(r.Context(), account, &req)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	http.Error(w, `{"error": "invalid message format"}`, http.StatusBadRequest)
}

// handleRequest handles an MCP request
func (h *MCPStreamableHTTP) handleRequest(r interface{}, account *Account, req *MCPRequest) *MCPResponse {
	switch req.Method {
	case "initialize":
		return h.handleInitialize(req)
	case "tools/list":
		return h.handleToolsList(req)
	case "tools/call":
		return h.handleToolsCall(req, account)
	case "resources/list":
		return h.handleResourcesList(req)
	case "resources/read":
		return h.handleResourcesRead(req)
	case "ping":
		return h.handlePing(req)
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}
}

// handleInitialize handles the initialize request
func (h *MCPStreamableHTTP) handleInitialize(req *MCPRequest) *MCPResponse {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true,
			},
			"resources": map[string]interface{}{
				"subscribe":   true,
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "HyperNexus Cloud",
			"version": "1.0.0",
		},
	}

	resultJSON, _ := json.Marshal(result)
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  resultJSON,
	}
}

// handleToolsList handles the tools/list request
func (h *MCPStreamableHTTP) handleToolsList(req *MCPRequest) *MCPResponse {
	tools := []map[string]interface{}{
		{
			"name":        "memory_store",
			"description": "Store a memory in the persistent memory system",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"content": map[string]interface{}{
						"type":        "string",
						"description": "The content to store",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Tags for the memory",
						"items":       map[string]string{"type": "string"},
					},
				},
				"required": []string{"content"},
			},
		},
		{
			"name":        "memory_search",
			"description": "Search memories by query",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query",
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			"name":        "tool_search",
			"description": "Search for MCP tools",
			"inputSchema": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The search query",
					},
				},
				"required": []string{"query"},
			},
		},
	}

	toolsJSON, _ := json.Marshal(tools)
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  toolsJSON,
	}
}

// handleToolsCall handles the tools/call request
func (h *MCPStreamableHTTP) handleToolsCall(req *MCPRequest, account *Account) *MCPResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	// Forward to container
	container, err := h.provisioner.GetContainerByAccount(account.ID)
	if err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32000,
				Message: "Container not available",
			},
		}
	}

	// Make request to container
	url := fmt.Sprintf("http://localhost:%d/api/mcp/tools/call", container.Port)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(params.Arguments)))
	if err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32000,
				Message: "Failed to call tool",
			},
		}
	}
	defer resp.Body.Close()

	// Read response
	body, _ := io.ReadAll(resp.Body)
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  body,
	}
}

// handleResourcesList handles the resources/list request
func (h *MCPStreamableHTTP) handleResourcesList(req *MCPRequest) *MCPResponse {
	resources := []map[string]interface{}{
		{
			"uri":         "hypernexus://memory",
			"name":        "Persistent Memory",
			"description": "Access to the persistent memory system",
			"mimeType":    "application/json",
		},
		{
			"uri":         "hypernexus://tools",
			"name":        "MCP Tools",
			"description": "Available MCP tools",
			"mimeType":    "application/json",
		},
	}

	resourcesJSON, _ := json.Marshal(resources)
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  resourcesJSON,
	}
}

// handleResourcesRead handles the resources/read request
func (h *MCPStreamableHTTP) handleResourcesRead(req *MCPRequest) *MCPResponse {
	var params struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	// Handle different resources
	switch params.URI {
	case "hypernexus://memory":
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"type": "memory", "status": "available"}`),
		}
	case "hypernexus://tools":
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  json.RawMessage(`{"type": "tools", "status": "available"}`),
		}
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Resource not found",
			},
		}
	}
}

// handlePing handles the ping request
func (h *MCPStreamableHTTP) handlePing(req *MCPRequest) *MCPResponse {
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  json.RawMessage(`{"pong": true}`),
	}
}

// handleHealth handles health check requests
func (h *MCPStreamableHTTP) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"version":   "1.0.0",
	})
}

// authenticate authenticates a request
func (h *MCPStreamableHTTP) authenticate(r *http.Request) (*Account, error) {
	// Check API key header
	apiKey := r.Header.Get("X-API-Key")
	if apiKey != "" {
		return h.am.GetByAPIKey(apiKey)
	}

	// Check Authorization header
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		apiKey = strings.TrimPrefix(auth, "Bearer ")
		return h.am.GetByAPIKey(apiKey)
	}

	// Check query parameter
	apiKey = r.URL.Query().Get("api_key")
	if apiKey != "" {
		return h.am.GetByAPIKey(apiKey)
	}

	return nil, fmt.Errorf("no authentication provided")
}

// readSSEMessages reads SSE messages from a reader
func readSSEMessages(reader *bufio.Reader) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			if strings.HasPrefix(line, "data: ") {
				ch <- strings.TrimPrefix(line, "data: ")
			}
		}
	}()
	return ch
}

func init() {
	log.Println("[MCP] Streamable HTTP transport initialized")
}
