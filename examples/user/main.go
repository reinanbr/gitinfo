package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/reinanbr/gitinfo"
)

func main() {
	tokens, err := godotenv.Read(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}

	user := "reinanbr"
	token := tokens["GITHUB_TOKEN"]

	userInfo, err := gitinfo.GetUserInfo(user, token)
	if err != nil {
		fmt.Printf("Error fetching user info: %v\n", err)
		return
	}

	fmt.Printf("User: %s (%s)\n", userInfo.Name, userInfo.Login)
	fmt.Printf("Bio: %s\n", userInfo.Bio)
	fmt.Printf("Created At: %s\n", userInfo.CreatedAt)
	fmt.Printf("Avatar: %s\n", userInfo.AvatarUrl)
	fmt.Printf("URL: %s\n", userInfo.URL)
	fmt.Printf("Repos: %d\n", userInfo.Repositories.TotalCount)
	fmt.Printf("Followers: %d\n", userInfo.Followers.TotalCount)
	fmt.Printf("Following: %d\n", userInfo.Following.TotalCount)
}
