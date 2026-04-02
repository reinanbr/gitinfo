package query

import ("fmt")


func BuildUserQuery(username string) string {
	return fmt.Sprintf(`
	{
		user(login: "%s") {
			id
			name
			login
			bio
			avatarUrl
			createdAt
			url
			followers {
				totalCount
			}
			following {
				totalCount
			}
			repositories(privacy: PUBLIC) {
				totalCount
			}
		}
	}
	`, username)
}
