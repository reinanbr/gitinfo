package gitinfo

import (
	"sync"

	"github.com/reinanbr/gitinfo/internal/utils"
)

type ProfileInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Username        string `json:"username"`
	Bio             string `json:"bio"`
	AvatarUrl       string `json:"avatar_url"`
	URL             string `json:"url"`
	CreatedAt       string `json:"created_at"`
	Location        string `json:"location"`
	Company         string `json:"company"`
	WebsiteUrl      string `json:"website_url"`
	TwitterUsername string `json:"twitter_username"`
	IsHireable      bool   `json:"is_hireable"`
	TotalRepos      int    `json:"total_repos"`
	TotalStars      int    `json:"total_stars"`
	TotalCommits    int    `json:"total_commits"`
	TotalFollowers  int    `json:"total_followers"`
	TotalFollowing  int    `json:"total_following"`
	CurrentStreak   int    `json:"current_streak"`
	MaxStreak       int    `json:"max_streak"`
}

// GetTotalStars sums stargazerCount across every public repository owned by
// username. GitHub's GraphQL API has no direct "total stars received" field
// on a user, so this walks the same paginated repository list GetReposInfo
// uses and adds it up.
func GetTotalStars(username, token string) (int, error) {
	repos, err := GetReposInfo(username, token)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, r := range repos {
		total += r.StargazerCount
	}
	return total, nil
}

// GetProfile fetches the core "who is this user" snapshot in one call:
// identity fields, repo/follower/star counts, total commits, and streaks.
// The four underlying GraphQL calls run concurrently.
func GetProfile(username, token string) (ProfileInfo, error) {
	var (
		wg sync.WaitGroup

		userInfo utils.UserInfo
		userErr  error

		streaks    StreakResponse
		streaksErr error

		commits    CommitsResponse
		commitsErr error

		stars    int
		starsErr error
	)

	wg.Add(4)
	go func() {
		defer wg.Done()
		userInfo, userErr = GetUserInfo(username, token)
	}()
	go func() {
		defer wg.Done()
		streaks, streaksErr = GetStreaks(username, token)
	}()
	go func() {
		defer wg.Done()
		commits, commitsErr = GetCommits(username, token)
	}()
	go func() {
		defer wg.Done()
		stars, starsErr = GetTotalStars(username, token)
	}()
	wg.Wait()

	if userErr != nil {
		return ProfileInfo{}, userErr
	}
	if streaksErr != nil {
		return ProfileInfo{}, streaksErr
	}
	if commitsErr != nil {
		return ProfileInfo{}, commitsErr
	}
	if starsErr != nil {
		return ProfileInfo{}, starsErr
	}

	return ProfileInfo{
		ID:              userInfo.ID,
		Name:            userInfo.Name,
		Username:        userInfo.Login,
		Bio:             userInfo.Bio,
		AvatarUrl:       userInfo.AvatarUrl,
		URL:             userInfo.URL,
		CreatedAt:       userInfo.CreatedAt,
		Location:        userInfo.Location,
		Company:         userInfo.Company,
		WebsiteUrl:      userInfo.WebsiteUrl,
		TwitterUsername: userInfo.TwitterUsername,
		IsHireable:      userInfo.IsHireable,
		TotalRepos:      userInfo.Repositories.TotalCount,
		TotalStars:      stars,
		TotalCommits:    commits.TotalCommits,
		TotalFollowers:  userInfo.Followers.TotalCount,
		TotalFollowing:  userInfo.Following.TotalCount,
		CurrentStreak:   streaks.Streak.CurrentStreak,
		MaxStreak:       streaks.Streak.MaxStreak,
	}, nil
}
