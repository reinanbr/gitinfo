package graphql

import (


	"fmt"
	"net/http"
	"bytes"
)

// MakeGraphQLRequest creates and executes a GraphQL HTTP request.
func MakeGraphQLRequest(token string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

