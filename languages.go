package gitinfo

import (
	"errors"
	"sort"

	"github.com/reinanbr/gitinfo/internal/graphql/fetch"
)

type LangPercentage struct {
	Lang       string
	Percentage float64
}

type ResponseLangs struct {
	LangPercentages []LangPercentage
	TotalBytes     int
	TotalRepos       int
}
func GetLangPercents(username string, token string, ignoreLangs []string) (ResponseLangs, error) {

	repos, err := fetch.FetchUserLangsFull(username, token)
	if err != nil {
		return ResponseLangs{}, err
	}

	langBytes := make(map[string]int)

	for _, repo := range repos.Repositories.Nodes {
		for _, edge := range repo.Languages.Edges {
			lang := edge.Node.Name
			shouldIgnore := false
			for _, ignoreLang := range ignoreLangs {
				if lang == ignoreLang {
					shouldIgnore = true
					break
				}
			}
			if shouldIgnore {
				continue
			}
			langBytes[lang] += edge.Size
		}
	}

	totalBytes := 0
	for _, size := range langBytes {
		totalBytes += size
	}

	if totalBytes == 0 {
		return ResponseLangs{}, errors.New("no language data found")
	}

	// Cria slice de LangPercentage
	var langPercentages []LangPercentage
	for lang, size := range langBytes {
		percent := (float64(size) / float64(totalBytes)) * 100
		langPercentages = append(langPercentages, LangPercentage{Lang: lang, Percentage: percent})
	}

	// Ordena do maior para o menor
	sort.Slice(langPercentages, func(i, j int) bool {
		return langPercentages[i].Percentage > langPercentages[j].Percentage
	})

	return ResponseLangs{
		LangPercentages: langPercentages,
		TotalBytes:     totalBytes,
		TotalRepos:       len(repos.Repositories.Nodes),
	},nil
}
