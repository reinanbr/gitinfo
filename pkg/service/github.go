package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/reinanbr/gitinfo/pkg/query"
)

const githubURL = "https://api.github.com/graphql"

type GraphQLQuery struct {
	Query string `json:"query"`
}
type RepoNode struct {
	Name            string    `json:"name"`
	CreatedAt       string    `json:"createdAt"`
	DefaultBranchRef *struct {
		Target struct {
			CommittedDate string `json:"committedDate"`
		} `json:"target"`
	} `json:"defaultBranchRef"`
	Languages struct {
		Edges []struct {
			Size int `json:"size"`
			Node struct {
				Name string `json:"name"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"languages"`
}


type UserInfo struct {
	Name      string `json:"name"`
	Login     string `json:"login"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatarUrl"`
	CreatedAt string `json:"createdAt"`
}

type RepoResponse struct {
	Data struct {
		User struct {
			Repositories struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []RepoNode `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type RepoNodeName struct {
	Name            string    `json:"name"`
}

type RepoResponseName struct {
	Data struct {
		User struct {
			Repositories struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []RepoNode `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}



type UserResponse struct {
	Data struct {
		User UserInfo `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func FetchUserInfo(username, token string) (UserInfo, error) {
	q := query.BuildUserQuery(username)
	body, _ := json.Marshal(GraphQLQuery{Query: q})

	req, err := http.NewRequest("POST", githubURL, bytes.NewBuffer(body))
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

	var response UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return UserInfo{}, err
	}
	if len(response.Errors) > 0 {
		return UserInfo{}, errors.New(response.Errors[0].Message)
	}

	return response.Data.User, nil
}

func FetchAllRepos(username, token string, cursor *string) ([]RepoNode, error) {
	q := query.BuildRepoQuery(username, cursor)
	body, _ := json.Marshal(GraphQLQuery{Query: q})

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

	var response RepoResponse
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

	return nodes, nil
}


func FetchAllReposName(username, token string, cursor *string,page int) ([]RepoNode, error) {
	q := query.BuildRepoNameQuery(username, cursor)
	body, _ := json.Marshal(GraphQLQuery{Query: q})
	
	req, err := http.NewRequest("POST", githubURL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error :%s\n",err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response RepoResponseName
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if len(response.Errors) > 0 {
		return nil, errors.New(response.Errors[0].Message)
	}
	nodes := response.Data.User.Repositories.Nodes
	//fmt.Printf("page: %d | len repos: %d\n",page,len(nodes))
	//recursion
	if response.Data.User.Repositories.PageInfo.HasNextPage{
		page += 1
		nextCursor := response.Data.User.Repositories.PageInfo.EndCursor
		nextNodes, err := FetchAllReposName(username, token, &nextCursor,page)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, nextNodes...);
	}
	return nodes, nil
}
		