package fetch

import (
	
	"time"

	"github.com/reinanbr/gitinfo/internal/graphql/query"
	"github.com/reinanbr/gitinfo/internal/graphql"
	"github.com/reinanbr/gitinfo/internal/utils"
)

// buildContributionGraphQuery constructs the GraphQL query for fetching contribution data.

// ExecuteContributionGraphRequests executes GraphQL queries for multiple years.
func ExecuteContributionGraphRequests(user string, years []int, token string) (map[int]utils.Response, error) {
	responses := make(map[int]utils.Response)

	for _, year := range years {
	

		quer := query.BuildContributionGraphQuery(user, year)
		response, err := graphql.ExecuteGraphQLQuery(quer, token)
		if err != nil {
			return nil, err
		}

		responses[year] = response
	}

	return responses, nil
}

// GetContributionGraphs retrieves contribution data for a user starting from a specific year.
func GetContributionGraphs(user string, startingYear int, token string) (map[int]utils.Response, error) {
	currentYear := time.Now().Year()


	// Fetch the user's creation year
	initialResponses, err := ExecuteContributionGraphRequests(user, []int{currentYear}, token)
	if err != nil {
		return nil, err
	}

	userCreatedYear := utils.ExtractUserCreatedYear(initialResponses, currentYear)
	minYear := utils.Max(startingYear, userCreatedYear)

	yearsToRequest := utils.GenerateYearRange(minYear, currentYear)
	moreResponses, err := ExecuteContributionGraphRequests(user, yearsToRequest, token)
	if err != nil {
		return nil, err
	}

	// Combine responses
	for year, resp := range moreResponses {
		initialResponses[year] = resp
	}

	return initialResponses, nil
}
