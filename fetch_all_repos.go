package gitinfo
// Package gitinfo provides utilities for fetching GitHub user information
// including repositories, contributions, languages, and user details.
//
// This library organizes GitHub API functionality into several packages:
//   - auth: GitHub token management
//   - github: User and repository information
//   - graphql: GraphQL queries and responses
//   - languages: Programming language analysis
//   - utils: Utility functions
//


import (	"github.com/reinanbr/gitinfo/pkg/utils"
)


func FetchAllRepos(user, token string) ([]utils.RepoNode, error) {
	ReposData, err := utils.FetchAllRepos(user, token, nil)
	if err != nil {
		return nil, err
	}
	repos := make([]utils.RepoNode, len(ReposData))
	for i, repo := range ReposData {
		repos[i] = utils.RepoNode{
			Name:      repo.Name,
			CreatedAt: repo.CreatedAt,
			Languages: repo.Languages,
		}
	}
	return repos, nil
}