package main

import (
	"fmt"
	"log"
	"os"
	"github.com/reinanbr/chronus"

	"github.com/joho/godotenv"
	gitinfo "github.com/reinanbr/gitinfo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	TOKEN := os.Getenv("GITHUB_TOKEN")
	user := "reinanbr"

	// init ping
	timeInit := chronus.Now()
	//user info
	userInfo, err := gitinfo.FetchUserInfo(user, TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("+=+=+=+=+=+= User Info =+=+=+=+=+=+")
	for key, value := range userInfo.(map[string]interface{}) {
		fmt.Printf("%s: %v\n", key, value)
	}
	timeUserInfo := chronus.Now() - timeInit;
	fmt.Printf("[Time to fetch user info: %vms]\n", timeUserInfo)
	fmt.Println("--------------------------------------------------")

	//repos
	repos, err := gitinfo.FetchAllRepos(user, TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("+=+=+=+=+=+= Repositories =+=+=+=+=+=+")
	for key, repo := range repos {
		fmt.Printf("Repository %d:\n", key+1)
		for k, v := range repo.Languages.Edges {
			fmt.Printf("  Language %d: %s (%d)\n", k+1, v.Node.Name, v.Size)
		}
		fmt.Println("----")
	}
	timeAllRepos := chronus.Now() - (timeInit + timeUserInfo);
	fmt.Printf("[Time to fetch repositories: %vms]\n", timeAllRepos)
	fmt.Println("--------------------------------------------------")
}
