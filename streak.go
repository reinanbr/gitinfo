package gitinfo

import (
	"github.com/reinanbr/gitinfo/internal/graphql/fetch"
	"github.com/reinanbr/gitinfo/internal/utils"
)

type StreakPeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type StreakData struct {
	MaxStreak           int          `json:"max_streak"`
	CurrentStreak       int          `json:"current_streak"`
	MaxStreakPeriod     StreakPeriod `json:"max_streak_period"`
	CurrentStreakPeriod StreakPeriod `json:"current_streak_period"`
}

type StreakResponse struct {
	User   string    `json:"user"`
	Streak StreakData `json:"streak"`
}

func GetStreaks(username string, token string) (StreakResponse, error) {
	startingYear := 2015
	graphs, err := fetch.GetContributionGraphs(username, startingYear, token)
	if err != nil {
		return StreakResponse{}, err
	}

	maxStreak, currentStreak, maxStart, maxEnd, currentStart, currentEnd := utils.GetContributionStreaks(graphs)

	return StreakResponse{
		User: username,
		Streak: StreakData{
			MaxStreak:     maxStreak,
			CurrentStreak: currentStreak,
			MaxStreakPeriod: StreakPeriod{
				Start: maxStart,
				End:   maxEnd,
			},
			CurrentStreakPeriod: StreakPeriod{
				Start: currentStart,
				End:   currentEnd,
			},
		},
	}, nil
}