package handler

import (
	"github.com/reinanbr/gitinfo/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

func GitCommit(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("user")
	if username == "" {
		http.Error(w, "Missing 'user' parameter", http.StatusBadRequest)
		return
	}

	startingYear := 2015
	graphs, err := utils.GetContributionGraphs(username, startingYear)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving graphsContribuitions: %v", err), http.StatusInternalServerError)
		return
	}

	sortYears := []int{}
	for year := range graphs {
		sortYears = append(sortYears, year)
	}
	sort.Ints(sortYears)


	// Monta a resposta JSON
	response := make(map[string]interface{})
	response["user"] = username
	response["commit"] = graphs

	// Calcula o total de commits
	total := 0
	for _, year := range sortYears {
		for _, week := range graphs[year].Data.User.ContributionsCollection.ContributionCalendar.Weeks {
			for _, day := range week.ContributionDays {
				total += day.ContributionCount
			}
		}
	}
	response["total"] = total

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
