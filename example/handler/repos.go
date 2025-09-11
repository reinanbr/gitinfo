package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/reinanbr/gitinfo/pkg/service"
	"github.com/reinanbr/gitinfo/pkg/utils"
)

func GitRepos(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	if username == "" {
		http.Error(w, "Missing 'user' parameter", http.StatusBadRequest)
		return
	}

	token, err := utils.GetGitHubTokenNative()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	repos, err := service.FetchAllRepos(username, token, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter repositórios: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"repositories": repos,
		"count":        len(repos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
