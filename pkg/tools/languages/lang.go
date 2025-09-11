package languages

import (
	"github.com/reinanbr/gitinfo/pkg/auth"
	"github.com/reinanbr/gitinfo/pkg/graphql"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
)

// Language represents a programming language.
type Language struct {
	Name string `json:"name"`
}

// LanguageEdge represents the size of a language in a repository.
type LanguageEdge struct {
	Size int      `json:"size"`
	Node Language `json:"node"`
}

// Repository represents a GitHub repository with language and metadata.
type Repository struct {
	Name       string `json:"name"`
	DateCreate string `json:"createdAt"`
	Languages  struct {
		Edges []LanguageEdge `json:"edges"`
	} `json:"languages"`
	DefaultBranchRef struct {
		Target struct {
			CommittedDate string `json:"committedDate"`
		} `json:"target"`
	} `json:"defaultBranchRef"`
}

// Repo represents a collection of repositories.
type Repo struct {
	Repositories struct {
		Nodes []Repository `json:"nodes"`
	} `json:"repositories"`
}

// ResponseLangs represents the GraphQL response for language data.
type ResponseLangs struct {
	Data struct {
		Repo Repo `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// RepositoryLite represents a lightweight repository structure.
type RepositoryLite struct {
	Name string `json:"name"`
}

// RepoName represents a collection of lightweight repositories.
type RepoName struct {
	Repositories struct {
		Nodes []RepositoryLite `json:"nodes"`
	} `json:"repositories"`
}

// ResponseLite represents the GraphQL response for lightweight repository data.
type ResponseLite struct {
	Data struct {
		Repo RepoName `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// FetchUserLangsFull fetches detailed language data for a user's repositories.
func FetchUserLangsFull(user string) (Repo, error) {
	token, err := auth.GetGitHubTokenNative()
	if err != nil {
		return Repo{}, fmt.Errorf("failed to get GitHub token: %w", err)
	}

	query, err := graphql.BuildGraphQLQueryLangFull(user)
	if err != nil {
		return Repo{}, fmt.Errorf("failed to build GraphQL query: %w", err)
	}
	body, _ := json.Marshal(graphql.GraphQLQuery{Query: query})

	resp, err := makeGraphQLRequest(token, body)
	if err != nil {
		return Repo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Repo{}, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var response ResponseLangs
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Repo{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Errors) > 0 {
		return Repo{}, errors.New(response.Errors[0].Message)
	}

	return response.Data.Repo, nil
}

// FetchUserLite fetches lightweight repository data for a user.
func FetchUserLite(user string) (RepoName, error) {
	token, err := auth.GetGitHubTokenNative()
	if err != nil {
		return RepoName{}, fmt.Errorf("failed to get GitHub token: %w", err)
	}

	query, err := graphql.BuildGraphQLQueryLite(user)
	if err != nil {
		return RepoName{}, fmt.Errorf("failed to build GraphQL query: %w", err)
	}
	body, _ := json.Marshal(graphql.GraphQLQuery{Query: query})

	resp, err := makeGraphQLRequest(token, body)
	if err != nil {
		return RepoName{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return RepoName{}, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var response ResponseLite
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return RepoName{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Errors) > 0 {
		return RepoName{}, errors.New(response.Errors[0].Message)
	}

	return response.Data.Repo, nil
}

// makeGraphQLRequest creates and executes a GraphQL HTTP request.
func makeGraphQLRequest(token string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// LanguagePercentage represents the percentage of a language in a repository.
type LanguagePercentage struct {
	Name  string
	Value float64
}

// CalculateLanguagePercentage calculates the percentage of each language in a user's repositories.
func CalculateLanguagePercentage(repo Repo) ([]LanguagePercentage, float64) {
	languageSizes := make(map[string]int)
	totalSize := 0

	for _, repository := range repo.Repositories.Nodes {
		for _, edge := range repository.Languages.Edges {
			if edge.Node.Name == "Jupyter Notebook" {
				continue
			}
			languageSizes[edge.Node.Name] += edge.Size
			totalSize += edge.Size
		}
	}

	var sortedLanguages []LanguagePercentage
	for lang, size := range languageSizes {
		sortedLanguages = append(sortedLanguages, LanguagePercentage{
			Name:  lang,
			Value: (float64(size) / float64(totalSize)) * 100,
		})
	}

	sort.Slice(sortedLanguages, func(i, j int) bool {
		return sortedLanguages[i].Value > sortedLanguages[j].Value
	})

	totalSizeMb := float64(totalSize) / (1024 * 1024)
	return sortedLanguages, totalSizeMb
}
