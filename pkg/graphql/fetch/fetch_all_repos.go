package fetch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/reinanbr/gitinfo/pkg/graphql/query"
	"github.com/reinanbr/gitinfo/pkg/utils"
)

func FetchAllRepos(username, token string, cursor *string) ([]utils.RepoNode, error) {
	if username == "" || token == "" {
		return nil, errors.New("username and token are required")
	}

	q := query.BuildRepoQuery(username, cursor)
	body, err := json.Marshal(utils.GraphQLQuery{Query: q})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", githubURL, bytes.NewBuffer(body))
	if err != nil {
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

	if response.Data.User.Repositories.PageInfo.HasNextPage {
		nextCursor := response.Data.User.Repositories.PageInfo.EndCursor
		nextNodes, err := FetchAllRepos(username, token, &nextCursor)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, nextNodes...)
	}
for i := range nodes {
    if nodes[i].DefaultBranchRef != nil {
        nodes[i].LastCommitDate = nodes[i].DefaultBranchRef.Target.CommittedDate
    } else {
        nodes[i].LastCommitDate = "N/A"
    }
}

return nodes, nil
}
