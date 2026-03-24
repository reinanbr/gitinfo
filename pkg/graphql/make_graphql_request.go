package graphql

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// MakeGraphQLRequest creates and executes a GraphQL HTTP request.
func MakeGraphQLRequest(token string, body []byte) (*http.Response, error) {
	const maxAttempts = 4
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "gitinfo")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to execute request: %w", err)
			if attempt < maxAttempts {
				time.Sleep(time.Duration(attempt*200) * time.Millisecond)
				continue
			}
			return nil, lastErr
		}

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusBadGateway || resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusGatewayTimeout {
			if attempt < maxAttempts {
				resp.Body.Close()
				time.Sleep(time.Duration(attempt*200) * time.Millisecond)
				continue
			}
		}

		return resp, nil
	}

	return nil, lastErr
}
