package graphql

import (
	"fmt"
)


type GraphQLQuery struct {
	Query string `json:"query"`
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




