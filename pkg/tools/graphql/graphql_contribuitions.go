package graphql

import (
	"encoding/json"
	"net/http"
	"errors"
	"io"
	"bytes"
)


//////////// structs //////////

type ContributionGraphQuery struct {
	Query string `json:"query"`
}

type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
}

type Week struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

type ContributionCalendar struct {
	Weeks []Week `json:"weeks"`
}

type ContributionsCollection struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
}

type User struct {
	CreatedAt              string                  `json:"createdAt"`
	ContributionsCollection ContributionsCollection `json:"contributionsCollection"`
}

type Data struct {
	User User `json:"user"`
}

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	Data   Data    `json:"data"`
	Errors []Error `json:"errors"`
}





func ExecuteGraphQLQuery(query, token string) (Response, error) {
	var response Response
	body, _ := json.Marshal(ContributionGraphQuery{Query: query})
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