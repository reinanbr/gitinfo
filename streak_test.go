package gitinfo

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestGetStreaks(t *testing.T) {
	tokens, err := godotenv.Read(".env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v\n", err)
	}

	githubToken, ok := tokens["GITHUB_TOKEN"]
	if !ok {
		t.Fatal("GITHUB_TOKEN not found in .env file")
	}

	streaks, err := GetStreaks("reinanbr", githubToken)
	if err != nil {
		t.Fatalf("Error fetching streaks: %v\n", err)
	}

	if streaks.User == "" {
		t.Fatal("user should not be empty")
	}

	if streaks.Streak.MaxStreak > 0 && (streaks.Streak.MaxStreakPeriod.Start == "" || streaks.Streak.MaxStreakPeriod.End == "") {
		t.Fatal("max_streak_period should include start and end when max_streak > 0")
	}

	if streaks.Streak.CurrentStreak > 0 && (streaks.Streak.CurrentStreakPeriod.Start == "" || streaks.Streak.CurrentStreakPeriod.End == "") {
		t.Fatal("current_streak_period should include start and end when current_streak > 0")
	}

	t.Logf("User: %s\n", streaks.User)
	t.Logf("Max Streak: %d\n", streaks.Streak.MaxStreak)
	t.Logf("Current Streak: %d\n", streaks.Streak.CurrentStreak)
	t.Logf("Max Streak Period: %s -> %s\n", streaks.Streak.MaxStreakPeriod.Start, streaks.Streak.MaxStreakPeriod.End)
	t.Logf("Current Streak Period: %s -> %s\n", streaks.Streak.CurrentStreakPeriod.Start, streaks.Streak.CurrentStreakPeriod.End)
}