package utils

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

// extractUserCreatedYear extracts the year the user was created from the GraphQL response.
func ExtractUserCreatedYear(responses map[int]Response, currentYear int) int {
	if response, exists := responses[currentYear]; exists && response.Data.User.CreatedAt != "" {
		createdAt := response.Data.User.CreatedAt
		if year, err := strconv.Atoi(strings.Split(createdAt, "-")[0]); err == nil {
			return year
		}
	}
	return 2005 // Default fallback year
}

// generateYearRange generates a range of years from start to end (exclusive).
func GenerateYearRange(start, end int) []int {
	years := []int{}
	for y := start; y < end; y++ {
		years = append(years, y)
	}
	return years
}

// GetContributionStreaks calculates maximum/current streaks and their date periods.
func GetContributionStreaks(responses map[int]Response) (int, int, string, string, string, string) {
	const dateLayout = "2006-01-02"

	var maxStreak, currentStreak, tempStreak int
	var maxStart, maxEnd, currentStart, currentEnd string
	var tempStartDate, tempEndDate time.Time
	var lastContributionDate time.Time
	lastDayHadContribution := false

	sortedYears := make([]int, 0, len(responses))
	for year := range responses {
		sortedYears = append(sortedYears, year)
	}
	sort.Ints(sortedYears)

	for _, year := range sortedYears {
		response := responses[year]
		for _, week := range response.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
			for _, day := range week.ContributionDays {
				date, err := time.Parse(dateLayout, day.Date)
				if err != nil {
					continue
				}

				if day.ContributionCount > 0 {
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
						maxStart = tempStartDate.Format(dateLayout)
						maxEnd = tempEndDate.Format(dateLayout)
					}
				} else {
					tempStreak = 0
					lastDayHadContribution = false
				}
			}
		}
	}

	if lastDayHadContribution && tempStreak > 0 {
		currentStreak = tempStreak
		currentStart = tempStartDate.Format(dateLayout)
		currentEnd = tempEndDate.Format(dateLayout)
	}

	return maxStreak, currentStreak, maxStart, maxEnd, currentStart, currentEnd
}

// GetTotalContributions calculates the total number of contributions.
func GetTotalContributions(responses map[int]Response) int {
	totalContributions := 0

	for _, response := range responses {
		for _, week := range response.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
			for _, day := range week.ContributionDays {
				totalContributions += day.ContributionCount
			}
		}
	}

	return totalContributions
}

// total contributions by year
// GetContributionsByYear calculates the total contributions for each year.
func GetContributionsByYear(responses map[int]Response) map[int]int {
	contributionsByYear := make(map[int]int)

	for year, response := range responses {
		totalContributions := 0
		for _, week := range response.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
			for _, day := range week.ContributionDays {
				totalContributions += day.ContributionCount
			}
		}
		contributionsByYear[year] = totalContributions
	}

	return contributionsByYear
}

// max returns the maximum of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
