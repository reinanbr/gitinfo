package query

import (
  "fmt"
)


// BuildGraphQLQueryLangFull generates a GraphQL query to retrieve the complete contribution history.
func BuildGraphQLQueryLangFull(user string) (string, error) {
  if err := validateUser(user); err != nil {
    return "", err
  }

  query := `
{
  user(login: "%s") {
  repositories(first: 100,privacy: PUBLIC) {
     nodes {
    name
    url
    isPrivate
    createdAt
    languages(first: 100, orderBy: {field: SIZE, direction: DESC}) {
          totalSize
          edges {
            size
            node {
              name
              color
            }
          }
        }
    defaultBranchRef {
      target {
      ... on Commit {
        committedDate
      }
      }
    }
    }
  }
  }
}
`
  return fmt.Sprintf(query, user), nil
}