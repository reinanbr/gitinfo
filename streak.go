package gitinfo

import (
	"github.com/reinanbr/gitinfo/internal/graphql/fetch"
	"github.com/reinanbr/gitinfo/internal/utils"
)

func GetStreaks(username string, token string) (map[string]interface{}, error) {
	startingYear := 2015
	graphs, err := fetch.GetContributionGraphs(username, startingYear, token)
	if err != nil {
		return nil, err
	}
	maxStreak, currentStreak, maxStart, maxEnd, currentStart, currentEnd := utils.GetContributionStreaks(graphs)

	// Monta a resposta JSON
	response := make(map[string]interface{})
	response["user"] = username
	response["streak"] = map[string]interface{}{
		"max_streak":            maxStreak,
		"current_streak":        currentStreak,
		"max_streak_period":     map[string]interface{}{"start": maxStart, "end": maxEnd},
		"current_streak_period": map[string]interface{}{"start": currentStart, "end": currentEnd},
	}
	return response, nil
}
