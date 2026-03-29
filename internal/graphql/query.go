package graphql


import "fmt"

func BuildUserQuery(username string) string {
	return fmt.Sprintf(`
	{
		user(login: "%s") {
			name
			login
			bio
			avatarUrl
			createdAt
		}
	}
	`, username)
}

func BuildRepoQuery(username string, cursor *string) string {
	after := ""
	if cursor != nil {
		after = fmt.Sprintf(`, after: "%s"`, *cursor)
	}
	return fmt.Sprintf(`
{
  user(login: "%s") {
    repositories(first: 100, privacy: PUBLIC%s, orderBy: {field: UPDATED_AT, direction: DESC}) {
      pageInfo {
        hasNextPage
        endCursor
      }
      nodes {
        name
        createdAt
        defaultBranchRef {
          target {
            ... on Commit {
              committedDate
            }
          }
        }
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

	`, username, after)
}



func BuildRepoNameQuery(username string, cursor *string) string {
	after := ""
	if cursor != nil {
		after = fmt.Sprintf(`, after: "%s"`, *cursor)
	}
	return fmt.Sprintf(`
{
  user(login: "%s") {
    repositories(first: 100, privacy: PUBLIC%s, orderBy: {field: UPDATED_AT, direction: DESC}) {
      pageInfo {
        hasNextPage
        endCursor
      }
      nodes {
        name
      }
    }
  }
}
	`, username, after)
}









func BuildGraphQLQueryUser(user string) string {
	// Esta consulta irá pegar o histórico de contribuições completo desde que o usuário entrou no GitHub
	return fmt.Sprintf(`
{
user(login: "%s") {
    name
    login
    bio
    avatarUrl
    createdAt
	}
}

	`, user)
}



