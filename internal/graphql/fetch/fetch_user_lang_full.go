package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/reinanbr/gitinfo/internal/graphql"
	"github.com/reinanbr/gitinfo/internal/graphql/query"
	"github.com/reinanbr/gitinfo/internal/utils"
)

// FetchUserLangsFull fetches detailed language data for ALL of a user's repositories (handles pagination).
func FetchUserLangsFull(user string, token string) (utils.Repo, error) {
	if user == "" || token == "" {
		return utils.Repo{}, errors.New("username and token are required")
	}

	var allRepos utils.Repo
	var cursor *string

	for {
		query, err := query.BuildGraphQLQueryLangFull(user, cursor)
		if err != nil {
			return utils.Repo{}, fmt.Errorf("failed to build GraphQL query: %w", err)
		}

		body, err := json.Marshal(utils.GraphQLQuery{Query: query})
		if err != nil {
			return utils.Repo{}, fmt.Errorf("failed to marshal query: %w", err)
		}

		resp, err := graphql.MakeGraphQLRequest(token, body)
		if err != nil {
			return utils.Repo{}, err
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return utils.Repo{}, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}

		var response utils.ResponseLangs
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			resp.Body.Close()
			return utils.Repo{}, fmt.Errorf("failed to decode response: %w", err)
		}
		resp.Body.Close()

		if len(response.Errors) > 0 {
			return utils.Repo{}, errors.New(response.Errors[0].Message)
		}

		// Merge repos from this page into allRepos
		allRepos.Repositories.Nodes = append(allRepos.Repositories.Nodes, response.Data.Repo.Repositories.Nodes...)

		// Stop if there are no more pages
		if !response.Data.Repo.Repositories.PageInfo.HasNextPage {
			break
		}

		// Advance cursor to next page
		nextCursor := response.Data.Repo.Repositories.PageInfo.EndCursor
		cursor = &nextCursor
	}

	return allRepos, nil
}
