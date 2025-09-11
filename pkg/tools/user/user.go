package user

import (
	"github.com/reinanbr/gitinfo/pkg/tools/auth"
	"github.com/reinanbr/gitinfo/pkg/tools/graphql"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserInfo struct {
	Name      string `json:"name"`
	Login     string `json:"login"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatarUrl"`
	CreatedAt string `json:"createdAt"`
}

type ResponseInfo struct {
	Data struct {
		UserInfo UserInfo `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func FetchUserData(user string) (UserInfo, error) {
	token, err := auth.GetGitHubTokenNative()
	if err != nil {
		return UserInfo{}, err
	}

	query := graphql.BuildGraphQLQueryUser(user)
	body, _ := json.Marshal(graphql.GraphQLQuery{Query: query})

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	if err != nil {
		return UserInfo{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UserInfo{}, fmt.Errorf("erro na requisição: status %d", resp.StatusCode)
	}

	var response ResponseInfo
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return UserInfo{}, err
	}

	if len(response.Errors) > 0 {
		return UserInfo{}, errors.New(response.Errors[0].Message)
	}

	return response.Data.UserInfo, nil
}
