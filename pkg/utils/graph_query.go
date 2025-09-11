package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)


// GraphQLError captura mensagens de erro da API GraphQL do GitHub
type GraphQLError struct {
	Message string `json:"message"`
}



// GraphQLResponse representa uma resposta genérica da API GraphQL do GitHub
type GraphQLResponse struct {
	Data   interface{}   `json:"data"`   // Conteúdo da resposta (qualquer tipo)
	Errors []GraphQLError `json:"errors"` // Lista de erros, se houver
}

// Método auxiliar para pegar os erros
func (r GraphQLResponse) GetErrors() []GraphQLError {
	return r.Errors
}

// GraphQLQuery define o payload da requisição
type GraphQLQuery struct {
	Query string `json:"query"`
}

// ExecuteGraphQLQuery executa qualquer query GraphQL no GitHub
func ExecuteGraphQLQuery(query string, token string, target interface{}) error {
	url := "https://api.github.com/graphql"

	// Monta o corpo da requisição
	body, _ := json.Marshal(GraphQLQuery{Query: query})

	// Cria a requisição HTTP
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Headers obrigatórios
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "GitHub-API-Client")

	// Executa
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verifica status da resposta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API error: status %d", resp.StatusCode)
	}

	// Faz o decode no target passado
	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		return err
	}

	// Checa erros específicos do GraphQL
	// Se o target possuir um campo Errors, tenta acessar
	if errResp, ok := target.(interface{ GetErrors() []struct{ Message string } }); ok {
		errorsList := errResp.GetErrors()
		if len(errorsList) > 0 {
			return errors.New(errorsList[0].Message)
		}
	}

	return nil
}
