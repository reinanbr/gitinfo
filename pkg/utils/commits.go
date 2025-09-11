package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"os"
	"errors"
	"math/rand"
	"encoding/json"
	"net/http"
	"bytes"
	"io"
	


)

type ContributionGraphQuery struct {
	Query string `json:"query"`
}

type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
}

type Week struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

type ContributionCalendar struct {
	Weeks []Week `json:"weeks"`
}

type ContributionsCollection struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
}

type User struct {
	CreatedAt              string                  `json:"createdAt"`
	ContributionsCollection ContributionsCollection `json:"contributionsCollection"`
}

type Data struct {
	User User `json:"user"`
}

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	Data   Data    `json:"data"`
	Errors []Error `json:"errors"`
}



func executeGraphQLQuery(query, token string) (Response, error) {
	var response Response
	body, _ := json.Marshal(ContributionGraphQuery{Query: query})
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GitHub-Readme-Streak-Stats")

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	json.Unmarshal(data, &response)

	if len(response.Errors) > 0 {
		return response, errors.New(response.Errors[0].Message)
	}

	return response, nil
}


func getGitHubTokens() []string {
	var tokens []string
	if token, exists := os.LookupEnv("TOKEN"); exists {
		tokens = append(tokens, token)
	}
	for i := 2; ; i++ {
		envVar := fmt.Sprintf("TOKEN%d", i)
		if token, exists := os.LookupEnv(envVar); exists {
			tokens = append(tokens, token)
		} else {
			break
		}
	}
	return tokens
}

func getGitHubToken(tokens []string) (string, error) {
	if len(tokens) == 0 {
		return "", errors.New("no GitHub token available")
	}
	return tokens[rand.Intn(len(tokens))], nil
}


// buildContributionGraphQuery constructs the GraphQL query for fetching contribution data.
func buildContributionGraphQuery(user string, year int) string {
	start := fmt.Sprintf("%d-01-01T00:00:00Z", year)
	end := fmt.Sprintf("%d-12-31T23:59:59Z", year)
	return fmt.Sprintf(`
		query {
			user(login: "%s") {
				createdAt
				contributionsCollection(from: "%s", to: "%s") {
					contributionCalendar {
						totalContributions
						weeks {
							contributionDays {
								contributionCount
								date
							}
						}
					}
					restrictedContributionsCount
				}
			}
		}
	`, user, start, end)
}


// ExecuteContributionGraphRequests executes GraphQL queries for multiple years.
func ExecuteContributionGraphRequests(user string, years []int, tokens []string) (map[int]Response, error) {
	responses := make(map[int]Response)

	for _, year := range years {
		token, err := getGitHubToken(tokens)
		if err != nil {
			return nil, err
		}

		query := buildContributionGraphQuery(user, year)
		response, err := executeGraphQLQuery(query, token)
		if err != nil {
			return nil, err
		}

		responses[year] = response
	}

	return responses, nil
}

// GetContributionGraphs retrieves contribution data for a user starting from a specific year.
func GetContributionGraphs(user string, startingYear int) (map[int]Response, error) {
	currentYear := time.Now().Year()
	tokens := getGitHubTokens()

	// Fetch the user's creation year
	initialResponses, err := ExecuteContributionGraphRequests(user, []int{currentYear}, tokens)
	if err != nil {
		return nil, err
	}

	userCreatedYear := extractUserCreatedYear(initialResponses, currentYear)
	minYear := max(startingYear, userCreatedYear)

	yearsToRequest := generateYearRange(minYear, currentYear)
	moreResponses, err := ExecuteContributionGraphRequests(user, yearsToRequest, tokens)
	if err != nil {
		return nil, err
	}

	// Combine responses
	for year, resp := range moreResponses {
		initialResponses[year] = resp
	}

	return initialResponses, nil
}

// extractUserCreatedYear extracts the year the user was created from the GraphQL response.
func extractUserCreatedYear(responses map[int]Response, currentYear int) int {
	if response, exists := responses[currentYear]; exists && response.Data.User.CreatedAt != "" {
		createdAt := response.Data.User.CreatedAt
		if year, err := strconv.Atoi(strings.Split(createdAt, "-")[0]); err == nil {
			return year
		}
	}
	return 2005 // Default fallback year
}

// generateYearRange generates a range of years from start to end (exclusive).
func generateYearRange(start, end int) []int {
	years := []int{}
	for y := start; y < end; y++ {
		years = append(years, y)
	}
	return years
}

// GetContributionStreaks calculates the maximum and current contribution streaks.
func GetContributionStreaks(responses map[int]Response) (int, int) {
	var maxStreak, currentStreak, tempStreak int
	var lastContributionDate time.Time
	streakActive := false

	for _, response := range responses {
		for _, week := range response.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
			for _, day := range week.ContributionDays {
				date, _ := time.Parse("2006-01-02", day.Date)
				if day.ContributionCount > 0 {
					if streakActive && date.Sub(lastContributionDate).Hours() == 24 {
						tempStreak++
					} else {
						tempStreak = 1
						streakActive = true
					}
					lastContributionDate = date
					maxStreak = max(maxStreak, tempStreak)
				} else {
					streakActive = false
					tempStreak = 0
				}
			}
		}
	}

	if streakActive {
		currentStreak = tempStreak
	}

	return maxStreak, currentStreak
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


//total contributions by year
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
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
