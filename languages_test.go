package gitinfo

import (
    "testing"
    "github.com/joho/godotenv"
)


func TestLangPercent(t *testing.T) {
	tokens, err := godotenv.Read(".env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v\n", err)
	}

	githubToken, ok := tokens["GITHUB_TOKEN"]
	if !ok {
		t.Fatal("GITHUB_TOKEN not found in .env file")
	}
	ignoreLangs := []string{"Jupyter Notebook","TeX"}
	langPercentages, err := GetLangPercents("reinanbr", githubToken, ignoreLangs)
	if err != nil {
		t.Fatalf("Error calculating language percentages: %v\n", err)
	}
	t.Logf("Total bytes: %d\n", langPercentages.TotalBytes)
	t.Logf("Total repos: %d\n", langPercentages.TotalRepos)
	t.Logf("Language percentages:")
	for _, lp := range langPercentages.LangPercentages {
		t.Logf("Language: %s, Percentage: %.2f%%\n", lp.Lang, lp.Percentage)
	}
	
}