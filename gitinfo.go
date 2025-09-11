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
// Example usage:
//
//	import "github.com/reinanbr/gitinfo/pkg/auth"
//	import "github.com/reinanbr/gitinfo/pkg/github"
//
//	token, err := auth.GetGitHubTokenNative()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	userInfo, err := github.GetUserInfo("username", token)
//	if err != nil {
//		log.Fatal(err)
//	}
package gitinfo
