package fetch

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/reinanbr/gitinfo/pkg/graphql"
	"github.com/reinanbr/gitinfo/pkg/graphql/query"
	"github.com/reinanbr/gitinfo/pkg/utils"
)

// FetchUserLite fetches lightweight repository data for a user.
func FetchUserLite(user string, token string) (utils.RepoName, error) {


	query, err := query.BuildGraphQLQueryLite(user)
	if err != nil {
		return utils.RepoName{}, fmt.Errorf("failed to build GraphQL query: %w", err)
	}
	body, _ := json.Marshal(utils.GraphQLQuery{Query: query})

	resp, err := graphql.MakeGraphQLRequest(token, body)
	if err != nil {
		return utils.RepoName{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return utils.RepoName{}, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var response utils.ResponseLite
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return utils.RepoName{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Errors) > 0 {
		return utils.RepoName{}, errors.New(response.Errors[0].Message)
	}

	return response.Data.Repo, nil
}
