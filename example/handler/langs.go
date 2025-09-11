package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/reinanbr/gitinfo/pkg/service"
	"github.com/reinanbr/gitinfo/pkg/utils"
)

func GitLangs(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	if username == "" {
		http.Error(w, "Missing 'user' parameter", http.StatusBadRequest)
		return
	}

	tokens := utils.GetGitHubTokens()
	if len(tokens) == 0 {
		http.Error(w, "No GitHub tokens available", http.StatusInternalServerError)
		return
	}

	// Processar linguagens
	langPercentage, totalBytes, err := service.CalculateLanguagePercentages(username, tokens)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calculating languages: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":        username,
		"languages":   langPercentage,
		"total_bytes": totalBytes,
	})
}
