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

	//links
	fmt.Printf("Repos link of %s:\n", user)
	reposName, err := gitinfo.GetRepos(user, token)
	if err != nil {
		fmt.Printf("Error fetching repos name: %v\n", err)
		return
	}
	for _, repo := range reposName {
		fmt.Printf("https://github.com/%s/%s\n", user, repo.Name)
	}

	// repos with info
	fmt.Printf("\nRepos with info of %s:\n", user)
	reposInfo, err := gitinfo.GetReposInfo(user, token)
	if err != nil {
		fmt.Printf("Error fetching repos info: %v\n", err)
		return
	}
	for i, repo := range reposInfo {
		totalSize := 0
		for _, language := range repo.Languages.Edges {
			totalSize += language.Size
		}
		
		fmt.Printf("[%d]\nRepo: %s\n", i, repo.Name)
		fmt.Printf("-> Description: %s\n", repo.Description)
		fmt.Printf("-> URL: %s\n", repo.Url)
		fmt.Printf("-> Private: %t\n", repo.IsPrivate)
		fmt.Printf("-> Created At: %s\n", repo.CreatedAt)
		fmt.Printf("-> Last Commit Date: %s\n", repo.LastCommitDate)
		fmt.Printf("-> Total Size: %d bytes\n", totalSize)
		fmt.Printf("-> Languages:%d\n", len(repo.Languages.Edges))
		for _, language := range repo.Languages.Edges {
			percentLang := func(size int) float64 {
				if totalSize == 0 {
					return 0.0
				}
				return (float64(size) / float64(totalSize)) * 100
			}
			fmt.Printf("--> Language: %s | size: %d | percentage: %.2f%%\n", language.Node.Name, language.Size, percentLang(language.Size))
		}
		fmt.Println()
	}
}
