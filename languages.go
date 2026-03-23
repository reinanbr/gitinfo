package gitinfo

import (
	"errors"
	"sort"
	"github.com/reinanbr/gitinfo/pkg/graphql/fetch"
)

type LangPercentage struct {
	Lang       string
	Percentage float64
}

func CalculateLanguagePercentages(username string, token string) ([]LangPercentage, int, error) {

	repos, err := fetch.FetchAllRepos(username, token, nil)
	if err != nil {
		return nil, 0, err
	}

	langBytes := make(map[string]int)

	for _, repo := range repos {
		for _, edge := range repo.Languages.Edges {
			lang := edge.Node.Name
			if lang == "Jupyter Notebook" {
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
		return nil, 0, errors.New("no language data found")
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

	return langPercentages, totalBytes, nil
}