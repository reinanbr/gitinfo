package gitinfo_test

import (
	"fmt"
	"os"

	"github.com/reinanbr/gitinfo"
)

func ExampleGetReposInfo() {
	token := os.Getenv("GITHUB_TOKEN")

	repos, err := gitinfo.GetReposInfo("reinanbr", token)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(len(repos) > 0)
	// Output: true
}

func ExampleGetReposName() {
	token := os.Getenv("GITHUB_TOKEN")

	repos, err := gitinfo.GetReposName("reinanbr", token)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(len(repos) > 0)
	// Output: true
}

func ExampleGetLangPercents() {
	token := os.Getenv("GITHUB_TOKEN")

	ignoreLangs := []string{"Jupyter Notebook", "TeX"}

	result, err := gitinfo.GetLangPercents("reinanbr", token, ignoreLangs)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(result.TotalRepos > 0)
	fmt.Println(result.TotalBytes > 0)
	fmt.Println(len(result.LangPercentages) > 0)
	// Output:
	// true
	// true
	// true
}

func ExampleGetCommits() {
	token := os.Getenv("GITHUB_TOKEN")

	result, err := gitinfo.GetCommits("reinanbr", token)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(result.User)
	fmt.Println(result.TotalCommits > 0)
	fmt.Println(len(result.CommitsByYear) > 0)
	fmt.Println(len(result.CommitsByDay) > 0)
	// Output:
	// reinanbr
	// true
	// true
	// true
}

func ExampleGetStreaks() {
	token := os.Getenv("GITHUB_TOKEN")

	result, err := gitinfo.GetStreaks("reinanbr", token)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	_, hasUser := result["user"]
	_, hasStreak := result["streak"]

	fmt.Println(hasUser)
	fmt.Println(hasStreak)
	// Output:
	// true
	// true
}