package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/reinanbr/gitinfo"
)

func main() {
	user := flag.String("user", "", "GitHub username (required)")
	token := flag.String("token", "", "GitHub token (optional, falls back to GITHUB_TOKEN)")
	ignore := flag.String("ignore", "Jupyter Notebook,TeX", "comma-separated languages to ignore")
	showDays := flag.Bool("show-days", false, "print commitsByDay entries")
	flag.Parse()

	if *user == "" {
		exitWithError(errors.New("-user is required"))
	}

	_ = godotenv.Load(".env")
	resolvedToken := strings.TrimSpace(*token)
	if resolvedToken == "" {
		resolvedToken = strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	}
	if resolvedToken == "" {
		exitWithError(errors.New("missing GitHub token: provide -token or set GITHUB_TOKEN in env/.env"))
	}

	ignoreLangs := splitCSV(*ignore)

	fmt.Printf("Running gitinfo smoke test for user: %s\n", *user)
	fmt.Println(strings.Repeat("=", 72))

	if err := runReposChecks(*user, resolvedToken); err != nil {
		exitWithError(err)
	}
	if err := runLanguageChecks(*user, resolvedToken, ignoreLangs); err != nil {
		exitWithError(err)
	}
	if err := runCommitsAndStreakChecks(*user, resolvedToken, *showDays); err != nil {
		exitWithError(err)
	}

	fmt.Println(strings.Repeat("=", 72))
	fmt.Println("Smoke test completed successfully.")
}

func runReposChecks(user, token string) error {
	reposInfo, err := gitinfo.GetReposInfo(user, token)
	if err != nil {
		return fmt.Errorf("GetReposInfo failed: %w", err)
	}
	reposName, err := gitinfo.GetReposName(user, token)
	if err != nil {
		return fmt.Errorf("GetReposName failed: %w", err)
	}

	fmt.Printf("Repos: info=%d, names=%d\n", len(reposInfo), len(reposName))
	return nil
}

func runLanguageChecks(user, token string, ignoreLangs []string) error {
	langs, err := gitinfo.GetLangPercents(user, token, ignoreLangs)
	if err != nil {
		return fmt.Errorf("GetLangPercents failed: %w", err)
	}

	fmt.Printf("Languages: repos=%d, totalBytes=%d\n", langs.TotalRepos, langs.TotalBytes)
	for i, lp := range langs.LangPercentages {
		if i == 5 {
			break
		}
		fmt.Printf("  - %s: %.2f%%\n", lp.Lang, lp.Percentage)
	}
	return nil
}

func runCommitsAndStreakChecks(user, token string, showDays bool) error {
	commits, err := gitinfo.GetCommits(user, token)
	if err != nil {
		return fmt.Errorf("GetCommits failed: %w", err)
	}
	streakResp, err := gitinfo.GetStreaks(user, token)
	if err != nil {
		return fmt.Errorf("GetStreaks failed: %w", err)
	}

	today := time.Now().Format("2006-01-02")
	totalFromDays := 0
	for _, d := range commits.CommitsByDay {
		if d.Date > today {
			return fmt.Errorf("future date in commitsByDay: %s (today=%s)", d.Date, today)
		}
		if d.CountCommits < 0 {
			return fmt.Errorf("negative countCommits for date %s", d.Date)
		}
		totalFromDays += d.CountCommits
	}
	if totalFromDays != commits.TotalCommits {
		return fmt.Errorf("total mismatch: commits.TotalCommits=%d, sum(commitsByDay)=%d", commits.TotalCommits, totalFromDays)
	}

	calcMax, calcCurrent, calcMaxStart, calcMaxEnd, calcCurrentStart, calcCurrentEnd := calcStreakFromCommitsByDay(commits.CommitsByDay)

	streakData, ok := streakResp["streak"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected streak type: %T", streakResp["streak"])
	}

	realMax, ok := intFromAny(streakData["max_streak"])
	if !ok {
		return fmt.Errorf("unexpected max_streak type: %T", streakData["max_streak"])
	}
	realCurrent, ok := intFromAny(streakData["current_streak"])
	if !ok {
		return fmt.Errorf("unexpected current_streak type: %T", streakData["current_streak"])
	}

	realMaxPeriod, ok := streakData["max_streak_period"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected max_streak_period type: %T", streakData["max_streak_period"])
	}
	realCurrentPeriod, ok := streakData["current_streak_period"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected current_streak_period type: %T", streakData["current_streak_period"])
	}

	realMaxStart, _ := realMaxPeriod["start"].(string)
	realMaxEnd, _ := realMaxPeriod["end"].(string)
	realCurrentStart, _ := realCurrentPeriod["start"].(string)
	realCurrentEnd, _ := realCurrentPeriod["end"].(string)

	if calcMax != realMax || calcCurrent != realCurrent {
		return fmt.Errorf("streak mismatch: calc=(%d,%d) real=(%d,%d)", calcMax, calcCurrent, realMax, realCurrent)
	}
	if calcMaxStart != realMaxStart || calcMaxEnd != realMaxEnd {
		return fmt.Errorf("max streak period mismatch: calc=(%s -> %s) real=(%s -> %s)", calcMaxStart, calcMaxEnd, realMaxStart, realMaxEnd)
	}
	if calcCurrentStart != realCurrentStart || calcCurrentEnd != realCurrentEnd {
		return fmt.Errorf("current streak period mismatch: calc=(%s -> %s) real=(%s -> %s)", calcCurrentStart, calcCurrentEnd, realCurrentStart, realCurrentEnd)
	}

	fmt.Printf("Commits: total=%d, years=%d, days=%d\n", commits.TotalCommits, len(commits.CommitsByYear), len(commits.CommitsByDay))
	fmt.Printf("Streak: max=%d (%s -> %s), current=%d (%s -> %s)\n", realMax, realMaxStart, realMaxEnd, realCurrent, realCurrentStart, realCurrentEnd)

	if showDays {
		fmt.Println("commitsByDay:")
		for _, d := range commits.CommitsByDay {
			fmt.Printf("  %s: %d\n", d.Date, d.CountCommits)
		}
	}

	return nil
}

func calcStreakFromCommitsByDay(days []gitinfo.CommitByDate) (int, int, string, string, string, string) {
	const layout = "2006-01-02"

	if len(days) == 0 {
		return 0, 0, "", "", "", ""
	}

	sortedDays := make([]gitinfo.CommitByDate, len(days))
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

func splitCSV(input string) []string {
	if strings.TrimSpace(input) == "" {
		return []string{}
	}
	parts := strings.Split(input, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
