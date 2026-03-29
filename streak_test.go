package gitinfo

import (
	"testing"

	"github.com/joho/godotenv"
)

func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int32:
		return int(n), true
	case int64:
		return int(n), true
	case float32:
		return int(n), true
	case float64:
		return int(n), true
	default:
		return 0, false
	}
}

func toStringMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}

func toString(v interface{}) (string, bool) {
	s, ok := v.(string)
	return s, ok
}

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
	t.Logf("User: %s\n", streaks["user"])
	streakData := streaks["streak"].(map[string]interface{})
	maxStreak, ok := toInt(streakData["max_streak"])
	if !ok {
		t.Fatalf("max_streak has unexpected type: %T", streakData["max_streak"])
	}
	currentStreak, ok := toInt(streakData["current_streak"])
	if !ok {
		t.Fatalf("current_streak has unexpected type: %T", streakData["current_streak"])
	}
	maxPeriod, ok := toStringMap(streakData["max_streak_period"])
	if !ok {
		t.Fatalf("max_streak_period has unexpected type: %T", streakData["max_streak_period"])
	}
	currentPeriod, ok := toStringMap(streakData["current_streak_period"])
	if !ok {
		t.Fatalf("current_streak_period has unexpected type: %T", streakData["current_streak_period"])
	}

	maxStart, ok := toString(maxPeriod["start"])
	if !ok {
		t.Fatalf("max_streak_period.start has unexpected type: %T", maxPeriod["start"])
	}
	maxEnd, ok := toString(maxPeriod["end"])
	if !ok {
		t.Fatalf("max_streak_period.end has unexpected type: %T", maxPeriod["end"])
	}
	currentStart, ok := toString(currentPeriod["start"])
	if !ok {
		t.Fatalf("current_streak_period.start has unexpected type: %T", currentPeriod["start"])
	}
	currentEnd, ok := toString(currentPeriod["end"])
	if !ok {
		t.Fatalf("current_streak_period.end has unexpected type: %T", currentPeriod["end"])
	}

	if maxStreak > 0 && (maxStart == "" || maxEnd == "") {
		t.Fatal("max_streak_period should include start and end when max_streak > 0")
	}
	if currentStreak > 0 && (currentStart == "" || currentEnd == "") {
		t.Fatal("current_streak_period should include start and end when current_streak > 0")
	}
	t.Logf("Max Streak: %d\n", maxStreak)
	t.Logf("Current Streak: %d\n", currentStreak)
	t.Logf("Max Streak Period: %s -> %s\n", maxStart, maxEnd)
	t.Logf("Current Streak Period: %s -> %s\n", currentStart, currentEnd)
}
