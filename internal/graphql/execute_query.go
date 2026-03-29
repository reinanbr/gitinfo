package graphql



import (
	"bytes"
	"encoding/json"
	"errors"

	"net/http"
	"github.com/reinanbr/gitinfo/internal/utils"
	"io"
)

func ExecuteGraphQLQuery(query, token string) (utils.Response, error) {
	var response utils.Response
	body, _ := json.Marshal(utils.ContributionGraphQuery{Query: query})
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GitHub-Readme-Streak-Stats")

	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	json.Unmarshal(data, &response)

	if len(response.Errors) > 0 {
		return response, errors.New(response.Errors[0].Message)
	}

	return response, nil
}