package utils

import (
	"testing"
	"time"
)

const testDateLayout = "2006-01-02"

func buildResponses(year int, days []ContributionDay) map[int]Response {
	responses := map[int]Response{year: {}}
	resp := responses[year]
	resp.Data.User.ContributionsCollection.ContributionCalendar.Weeks = []Week{
		{ContributionDays: days},
	}
	responses[year] = resp
	return responses
}

// A streak that ended yesterday should still count as "current" today,
// since today isn't over yet and hasn't had a chance to break it.
func TestGetContributionStreaks_GraceForToday(t *testing.T) {
	today := time.Now().UTC()
	yesterday := today.AddDate(0, 0, -1)
	twoDaysAgo := today.AddDate(0, 0, -2)

	days := []ContributionDay{
		{Date: twoDaysAgo.Format(testDateLayout), ContributionCount: 1},
		{Date: yesterday.Format(testDateLayout), ContributionCount: 1},
		{Date: today.Format(testDateLayout), ContributionCount: 0}, // no commit yet today
	}

	maxStreak, currentStreak, _, _, currentStart, currentEnd := GetContributionStreaks(buildResponses(today.Year(), days))

	if maxStreak != 2 {
		t.Fatalf("expected maxStreak=2, got %d", maxStreak)
	}
	if currentStreak != 2 {
		t.Fatalf("expected currentStreak=2 (grace period for today), got %d", currentStreak)
	}
	if currentStart != twoDaysAgo.Format(testDateLayout) || currentEnd != yesterday.Format(testDateLayout) {
		t.Fatalf("unexpected current streak period: %s -> %s", currentStart, currentEnd)
	}
}

// A streak with no contribution for a full missed day should be reported as broken.
func TestGetContributionStreaks_BrokenStreak(t *testing.T) {
	today := time.Now().UTC()
	threeDaysAgo := today.AddDate(0, 0, -3)

	days := []ContributionDay{
		{Date: threeDaysAgo.Format(testDateLayout), ContributionCount: 1},
		{Date: today.Format(testDateLayout), ContributionCount: 0},
	}

	_, currentStreak, _, _, _, _ := GetContributionStreaks(buildResponses(today.Year(), days))
	if currentStreak != 0 {
		t.Fatalf("expected currentStreak=0 for a streak broken multiple days ago, got %d", currentStreak)
	}
}

// A contribution made today should extend the current streak as before.
func TestGetContributionStreaks_ContributionToday(t *testing.T) {
	today := time.Now().UTC()
	yesterday := today.AddDate(0, 0, -1)

	days := []ContributionDay{
		{Date: yesterday.Format(testDateLayout), ContributionCount: 1},
		{Date: today.Format(testDateLayout), ContributionCount: 3},
	}

	_, currentStreak, _, _, _, currentEnd := GetContributionStreaks(buildResponses(today.Year(), days))
	if currentStreak != 2 {
		t.Fatalf("expected currentStreak=2, got %d", currentStreak)
	}
	if currentEnd != today.Format(testDateLayout) {
		t.Fatalf("expected current streak to end today, got %s", currentEnd)
	}
}
