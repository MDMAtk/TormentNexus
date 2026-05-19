package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) handleHypercodeProtocol(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uri := r.URL.Query().Get("uri")
	if uri == "" || !strings.HasPrefix(uri, "hypercode://") {
		http.Error(w, "Invalid or missing hypercode:// URI", http.StatusBadRequest)
		return
	}

	// Basic parsing: hypercode://attach?session=xyz
	parts := strings.SplitN(strings.TrimPrefix(uri, "hypercode://"), "?", 2)
	action := parts[0]

	params := make(map[string]string)
	if len(parts) > 1 {
		pairs := strings.Split(parts[1], "&")
		for _, pair := range pairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				params[kv[0]] = kv[1]
			}
		}
	}

	response := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"action": action,
			"params": params,
			"status": "attached",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
