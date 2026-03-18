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








func FetchUserInfo(username, token string) (utils.UserInfo, error) {
	if username == "" || token == "" {
		return utils.UserInfo{}, errors.New("username and token are required")
	}

	q := query.BuildUserQuery(username)
	body, err := json.Marshal(utils.GraphQLQuery{Query: q})
	if err != nil {
		return utils.UserInfo{}, err
	}

	req, err := http.NewRequest("POST", githubURL, bytes.NewBuffer(body))
	if err != nil {
		return utils.UserInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return utils.UserInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return utils.UserInfo{}, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var response utils.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return utils.UserInfo{}, err
	}
	if len(response.Errors) > 0 {
		return utils.UserInfo{}, errors.New(response.Errors[0].Message)
	}

	return response.Data.User, nil
}
