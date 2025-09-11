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
import (
	"github.com/reinanbr/gitinfo/pkg/service"
)

func FetchUserInfo(user, token string) (interface{}, error) {
	UserInfoData, err := service.FetchUserInfo(user, token);
	if err != nil {
		return nil, err
	}
	UserInfo := map[string]interface{}{
		"login":             UserInfoData.Login,
		"name":              UserInfoData.Name,
		"Bio":               UserInfoData.Bio,
		"avatarUrl":         UserInfoData.AvatarUrl,
		"createdAt":         UserInfoData.CreatedAt,
	}
	return UserInfo, nil
}