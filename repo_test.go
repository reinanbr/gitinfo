package gitinfo

import (
    "testing"
    "github.com/joho/godotenv"
)

func TestFetchReposWithInfo(t *testing.T) {
    tokens, err := godotenv.Read(".env")
    if err != nil {
        t.Fatalf("Error loading .env file: %v\n", err)
    }

	githubToken, ok := tokens["GITHUB_TOKEN"]
    if !ok {
        t.Fatal("GITHUB_TOKEN not found in .env file")
    }
    repos, err := GetReposInfo("reinanbr", githubToken)
    if err != nil {
        t.Fatalf("Error fetching all repos: %v\n", err)
    }
    for i,repo := range repos {
        t.Logf("[%d] Repo: %s| Langs: %d\n", i, repo.Name, len(repo.Languages.Edges))
    }
}


func TestRepos(t *testing.T) {
    tokens, err := godotenv.Read(".env")
    if err != nil {
        t.Fatalf("Error loading .env file: %v\n", err)
    }

    githubToken, ok := tokens["GITHUB_TOKEN"]
    if !ok {
        t.Fatal("GITHUB_TOKEN not found in .env file")
    }
    repos, err := GetRepos("reinanbr", githubToken)
    if err != nil {
        t.Fatalf("Error fetching repos name: %v\n", err)
    }
    for i,repo := range repos {
        t.Logf("[%d] Repo: %s\n", i, repo.Name)
    }
}