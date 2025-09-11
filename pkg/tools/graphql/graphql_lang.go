package graphql

import (
  "fmt"
  "errors"
)

// validateUser ensures the user parameter is not empty.
func validateUser(user string) error {
  if user == "" {
    return errors.New("user parameter cannot be empty")
  }
  return nil
}

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
    createdAt
    languages(first: 100) {
      edges {
      size
      node {
        name
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

// BuildGraphQLQueryLite generates a GraphQL query to retrieve a lightweight contribution history.
func BuildGraphQLQueryLite(user string) (string, error) {
  if err := validateUser(user); err != nil {
    return "", err
  }

  query := `
{
  user(login: "%s") {
  repositories(first: 100,privacy: PUBLIC) {
    nodes {
    name
    }
  }
  }
}
`
  return fmt.Sprintf(query, user), nil
}
