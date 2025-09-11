package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type LanguageEdge struct {
	Size int `json:"size"`
	Node struct {
		Name string `json:"name"`
	} `json:"node"`
}

type RepoNode struct {
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	Languages struct {
		Edges []LanguageEdge `json:"edges"`
	} `json:"languages"`
}

type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type Repositories struct {
	PageInfo PageInfo    `json:"pageInfo"`
	Nodes    []RepoNode  `json:"nodes"`
}

type UserRepositories struct {
	Repositories Repositories `json:"repositories"`
}

type GraphQLReposResponse struct {
	User UserRepositories `json:"user"`
}


// Estruturas necessárias


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



// Monta a query para GraphQL
func BuildGraphQLQueryRepos(user string, cursor *string) string {
	after := ""
	if cursor != nil {
		after = fmt.Sprintf(`, after: "%s"`, *cursor)
	}

	return fmt.Sprintf(`
	{
		user(login: "%s") {
			repositories(first: 100, isFork: false%s) {
				pageInfo {
					hasNextPage
					endCursor
				}
				nodes {
					name
					createdAt
					defaultBranchRef {
          				target {
            				... on Commit {
              					committedDate
            				}
          				}
        			}
					languages(first: 100) {
						edges {
							size
							node {
								name
							}
						}
					}
				}
			}
		}
	}
	`, user, after)
}


// Função principal para buscar todos os repositórios
func FetchAllRepos(user string, token string, cursor *string) ([]RepoNode, error) {
	query := BuildGraphQLQueryRepos(user, cursor)

	body, _ := json.Marshal(map[string]string{"query": query})
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
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
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	if len(response.Errors) > 0 {
		return nil, errors.New(response.Errors[0].Message)
	}

	nodes := response.Data.User.Repositories.Nodes

	// Verifica se há mais páginas
	if response.Data.User.Repositories.PageInfo.HasNextPage {
		nextCursor := response.Data.User.Repositories.PageInfo.EndCursor
		nextNodes, err := FetchAllRepos(user, token, &nextCursor)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, nextNodes...)
	}

	return nodes, nil
}
