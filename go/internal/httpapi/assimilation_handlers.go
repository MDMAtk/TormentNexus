package httpapi

import (
	"net/http"
	"os/exec"
	"path/filepath"
)

func (s *Server) handleAssimilationTriggerResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"success": false, "error": "method not allowed"})
		return
	}

	scriptPath := filepath.Join(s.cfg.WorkspaceRoot, "scripts", "assimilate_all_resources.py")
	cmd := exec.CommandContext(r.Context(), "python", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
			"output":  string(output),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"output":  string(output),
	})
}

func (s *Server) handleAssimilationTriggerServers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"success": false, "error": "method not allowed"})
		return
	}

	scriptPath := filepath.Join(s.cfg.WorkspaceRoot, "scripts", "assimilate_mcp_servers.py")
	cmd := exec.CommandContext(r.Context(), "python", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   err.Error(),
			"output":  string(output),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
		"output":  string(output),
	})
}
