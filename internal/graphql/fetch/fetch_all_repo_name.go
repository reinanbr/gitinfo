package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"github.com/reinanbr/gitinfo/internal/utils"
	"github.com/reinanbr/gitinfo/internal/graphql/query"
	
)

const githubURL = "https://api.github.com/graphql"



func FetchAllReposName(username, token string, cursor *string, page int) ([]utils.RepoNode, error) {
	if username == "" || token == "" {
		return nil, errors.New("username and token are required")
	}

	q := query.BuildRepoNameQuery(username, cursor)
	body, err := json.Marshal(utils.GraphQLQuery{Query: q})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", githubURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error :%s\n", err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var response utils.RepoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Errors) > 0 {
		return nil, errors.New(response.Errors[0].Message)
	}

	nodes := response.Data.User.Repositories.Nodes

	// recursion with nil check
	if response.Data.User.Repositories.PageInfo.HasNextPage {
		page += 1
		nextCursor := response.Data.User.Repositories.PageInfo.EndCursor
		nextNodes, err := FetchAllReposName(username, token, &nextCursor, page)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, nextNodes...)
	}
	return nodes, nil
}