package gitinfo

import (
    "testing"
    "github.com/joho/godotenv"
)

func TestGetReposInfo(t *testing.T) {
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
        t.Logf("[%d Repo: %s]\n", i, repo.Name)
        totalSize := 0
        for _,language := range repo.Languages.Edges {
            totalSize += language.Size
        }
        percentLang := func(size int) float64 {
            if totalSize == 0 {
                return 0.0
            }
            return (float64(size) / float64(totalSize)) * 100
        }
        for _,language := range repo.Languages.Edges {
            t.Logf("Language: %s\n | size: %d | percentage: %.2f%%", language.Node.Name, language.Size, percentLang(language.Size))
        }
    }
}


func TestGetReposName(t *testing.T) {
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
        t.Logf("[%d] https://github.com/%s/%s\n", i, "reinanbr", repo.Name)
    }
}