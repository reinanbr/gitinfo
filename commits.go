package gitinfo

import (
	"sort"
	"time"

	"github.com/reinanbr/gitinfo/internal/graphql/fetch"
)

type CommitByDate struct {
	Date         string `json:"date"`
	CountCommits int    `json:"countCommits"`
}

type CommitsYear struct {
	Year    int            `json:"year"`
	Commits []CommitByDate `json:"commits"`
}

type CommitsResponse struct {
	User          string         `json:"user"`
	TotalCommits  int            `json:"totalCommits"`
	CommitsByYear []CommitsYear  `json:"commitsByYear"`
	CommitsByDay  []CommitByDate `json:"commitsByDay"`
}

func GetCommits(username string, token string) (CommitsResponse, error) {
	startingYear := 2015
	graphs, err := fetch.GetContributionGraphs(username, startingYear, token)
	if err != nil {
		return CommitsResponse{}, err
	}

	sortYears := []int{}
	for year := range graphs {
		sortYears = append(sortYears, year)
	}
	sort.Ints(sortYears)
	today := time.Now().Format("2006-01-02")

	total := 0
	commitsByYear := make([]CommitsYear, 0, len(sortYears))
	commitsByDay := make([]CommitByDate, 0)
	for _, year := range sortYears {
		yearly := CommitsYear{Year: year, Commits: []CommitByDate{}}
		for _, week := range graphs[year].Data.User.ContributionsCollection.ContributionCalendar.Weeks {
			for _, day := range week.ContributionDays {
				if day.Date > today {
					continue
				}
				total += day.ContributionCount
				commitsByDay = append(commitsByDay, CommitByDate{
					Date:         day.Date,
					CountCommits: day.ContributionCount,
				})
				if day.ContributionCount > 0 {
					yearly.Commits = append(yearly.Commits, CommitByDate{
						Date:         day.Date,
						CountCommits: day.ContributionCount,
					})
				}
			}
		}
		commitsByYear = append(commitsByYear, yearly)
	}

	response := CommitsResponse{
		User:          username,
		TotalCommits:  total,
		CommitsByYear: commitsByYear,
		CommitsByDay:  commitsByDay,
	}

	return response, nil
}
