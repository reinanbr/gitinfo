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
			location
			company
			websiteUrl
			twitterUsername
			isHireable
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
