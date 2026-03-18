package fetch


import (

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/reinanbr/gitinfo/pkg/graphql/query"
	"github.com/reinanbr/gitinfo/pkg/utils"
	"github.com/reinanbr/gitinfo/pkg/graphql"

)



// FetchUserLangsFull fetches detailed language data for a user's repositories.
func FetchUserLangsFull(user string, token string) (utils.Repo, error) {


	query, err := query.BuildGraphQLQueryLangFull(user)
	if err != nil {
		return utils.Repo{}, fmt.Errorf("failed to build GraphQL query: %w", err)
	}
	body, _ := json.Marshal(utils.GraphQLQuery{Query: query})

	resp, err := graphql.MakeGraphQLRequest(token, body)
	if err != nil {
		return utils.Repo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return utils.Repo{}, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var response utils.ResponseLangs
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return utils.Repo{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Errors) > 0 {
		return utils.Repo{}, errors.New(response.Errors[0].Message)
	}

	return response.Data.Repo, nil
}