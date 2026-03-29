package query

import ("fmt")


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