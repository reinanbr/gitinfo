package gitinfo

import (
	"sort"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func calcStreakFromCommitsByDay(days []CommitByDate) (int, int, string, string, string, string) {
	const layout = "2006-01-02"

	if len(days) == 0 {
		return 0, 0, "", "", "", ""
	}

	sortedDays := make([]CommitByDate, len(days))
	copy(sortedDays, days)
	sort.Slice(sortedDays, func(i, j int) bool {
		return sortedDays[i].Date < sortedDays[j].Date
	})

	var maxStreak, currentStreak, tempStreak int
	var maxStart, maxEnd, currentStart, currentEnd string
	var tempStartDate, tempEndDate time.Time
	var lastContributionDate time.Time
	lastDayHadContribution := false

	for _, day := range sortedDays {
		date, err := time.Parse(layout, day.Date)
		if err != nil {
			continue
		}

		if day.CountCommits > 0 {
			if tempStreak > 0 && date.Sub(lastContributionDate).Hours() == 24 {
				tempStreak++
				tempEndDate = date
			} else {
				tempStreak = 1
				tempStartDate = date
				tempEndDate = date
			}

			lastContributionDate = date
			lastDayHadContribution = true

			if tempStreak > maxStreak {
				maxStreak = tempStreak
				maxStart = tempStartDate.Format(layout)
				maxEnd = tempEndDate.Format(layout)
			}
		} else {
			tempStreak = 0
			lastDayHadContribution = false
		}
	}

	if lastDayHadContribution && tempStreak > 0 {
		currentStreak = tempStreak
		currentStart = tempStartDate.Format(layout)
		currentEnd = tempEndDate.Format(layout)
	}

	return maxStreak, currentStreak, maxStart, maxEnd, currentStart, currentEnd
}

func intFromAny(v interface{}) (int, bool) {
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

func mapFromAny(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}

func stringFromAny(v interface{}) (string, bool) {
	s, ok := v.(string)
	return s, ok
}

func TestGetCommits(t *testing.T) {
	tokens, err := godotenv.Read(".env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v\n", err)
	}

	githubToken, ok := tokens["GITHUB_TOKEN"]
	if !ok {
		t.Fatal("GITHUB_TOKEN not found in .env file")
	}
	commits, err := GetCommits("reinanbr", githubToken)
	if err != nil {
		t.Fatalf("Error fetching commits: %v\n", err)
	}
	t.Logf("User: %s\n", commits.User)
	t.Logf("Total commits: %d\n", commits.TotalCommits)
	today := time.Now().Format("2006-01-02")

	if len(commits.CommitsByDay) == 0 {
		t.Fatal("commitsByDay should not be empty")
	}

	totalFromDays := 0
	dayCounts := make(map[string]int, len(commits.CommitsByDay))
	for _, day := range commits.CommitsByDay {
		if day.Date == "" {
			t.Fatal("commitsByDay.date should not be empty")
		}
		if day.Date > today {
			t.Fatalf("commitsByDay contains future date: %s (today: %s)", day.Date, today)
		}
		if day.CountCommits < 0 {
			t.Fatalf("commitsByDay.countCommits should be >= 0, got %d", day.CountCommits)
		}
		dayCounts[day.Date] = day.CountCommits
		totalFromDays += day.CountCommits
		t.Logf("Day: %s | countCommits: %d\n", day.Date, day.CountCommits)
	}

	if totalFromDays != commits.TotalCommits {
		t.Fatalf("totalCommits mismatch with commitsByDay (expected: %d, got: %d)", commits.TotalCommits, totalFromDays)
	}

	totalFromYears := 0
	for _, year := range commits.CommitsByYear {
		t.Logf("Year: %d, days with commits: %d\n", year.Year, len(year.Commits))
		for _, day := range year.Commits {
			if day.Date == "" {
				t.Fatal("date should not be empty")
			}
			if day.CountCommits <= 0 {
				t.Fatalf("countCommits should be > 0, got %d", day.CountCommits)
			}
			totalFromYears += day.CountCommits

			countByDay, exists := dayCounts[day.Date]
			if !exists {
				t.Fatalf("commitsByYear day not found in commitsByDay: %s", day.Date)
			}
			if countByDay != day.CountCommits {
				t.Fatalf("count mismatch for day %s (commitsByDay: %d, commitsByYear: %d)", day.Date, countByDay, day.CountCommits)
			}
		}
	}

	if totalFromYears > commits.TotalCommits {
		t.Fatalf("sum of commitsByYear cannot exceed totalCommits (years: %d, total: %d)", totalFromYears, commits.TotalCommits)
	}

	calcMax, calcCurrent, calcMaxStart, calcMaxEnd, calcCurrentStart, calcCurrentEnd := calcStreakFromCommitsByDay(commits.CommitsByDay)

	streakResp, err := GetStreaks("reinanbr", githubToken)
	if err != nil {
		t.Fatalf("Error fetching streaks: %v\n", err)
	}

	streakData, ok := mapFromAny(streakResp["streak"])
	if !ok {
		t.Fatalf("streak response has unexpected type: %T", streakResp["streak"])
	}

	realMax, ok := intFromAny(streakData["max_streak"])
	if !ok {
		t.Fatalf("max_streak has unexpected type: %T", streakData["max_streak"])
	}
	realCurrent, ok := intFromAny(streakData["current_streak"])
	if !ok {
		t.Fatalf("current_streak has unexpected type: %T", streakData["current_streak"])
	}

	realMaxPeriod, ok := mapFromAny(streakData["max_streak_period"])
	if !ok {
		t.Fatalf("max_streak_period has unexpected type: %T", streakData["max_streak_period"])
	}
	realCurrentPeriod, ok := mapFromAny(streakData["current_streak_period"])
	if !ok {
		t.Fatalf("current_streak_period has unexpected type: %T", streakData["current_streak_period"])
	}

	realMaxStart, ok := stringFromAny(realMaxPeriod["start"])
	if !ok {
		t.Fatalf("max_streak_period.start has unexpected type: %T", realMaxPeriod["start"])
	}
	realMaxEnd, ok := stringFromAny(realMaxPeriod["end"])
	if !ok {
		t.Fatalf("max_streak_period.end has unexpected type: %T", realMaxPeriod["end"])
	}
	realCurrentStart, ok := stringFromAny(realCurrentPeriod["start"])
	if !ok {
		t.Fatalf("current_streak_period.start has unexpected type: %T", realCurrentPeriod["start"])
	}
	realCurrentEnd, ok := stringFromAny(realCurrentPeriod["end"])
	if !ok {
		t.Fatalf("current_streak_period.end has unexpected type: %T", realCurrentPeriod["end"])
	}

	if calcMax != realMax {
		t.Fatalf("max streak mismatch (calculated: %d, real: %d)", calcMax, realMax)
	}
	if calcCurrent != realCurrent {
		t.Fatalf("current streak mismatch (calculated: %d, real: %d)", calcCurrent, realCurrent)
	}
	if calcMaxStart != realMaxStart || calcMaxEnd != realMaxEnd {
		t.Fatalf("max streak period mismatch (calculated: %s -> %s, real: %s -> %s)", calcMaxStart, calcMaxEnd, realMaxStart, realMaxEnd)
	}
	if calcCurrentStart != realCurrentStart || calcCurrentEnd != realCurrentEnd {
		t.Fatalf("current streak period mismatch (calculated: %s -> %s, real: %s -> %s)", calcCurrentStart, calcCurrentEnd, realCurrentStart, realCurrentEnd)
	}

	t.Logf("Streak cross-check OK | max: %d (%s -> %s) | current: %d (%s -> %s)", realMax, realMaxStart, realMaxEnd, realCurrent, realCurrentStart, realCurrentEnd)
}
