package query

import ("fmt")


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
