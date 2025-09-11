package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/reinanbr/gitinfo/pkg/service"
	"github.com/reinanbr/gitinfo/pkg/utils"
)

func GitUser(w http.ResponseWriter, r *http.Request) {
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

	userInfo, err := service.FetchUserInfo(username, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter dados do usuário: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"user": userInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
