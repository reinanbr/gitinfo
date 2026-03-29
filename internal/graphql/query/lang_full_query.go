package query

import (
	"fmt"
)

// BuildGraphQLQueryLangFull generates a GraphQL query to retrieve the complete contribution history.
func BuildGraphQLQueryLangFull(user string, cursor *string) (string, error) {
	if err := validateUser(user); err != nil {
		return "", err
	}

	after := ""
	if cursor != nil {
		after = fmt.Sprintf(`, after: "%s"`, *cursor)
	}

	query := `
{
  user(login: "%s") {
    repositories(first: 100, privacy: PUBLIC%s, orderBy: {field: UPDATED_AT, direction: DESC}) {
      pageInfo {
        hasNextPage
        endCursor
      }
      nodes {
        languages(first: 100) {
          edges {
            size
            node {
              name
            }
          }
        }
      }
    }
  }
}
`
	return fmt.Sprintf(query, user, after), nil
}
