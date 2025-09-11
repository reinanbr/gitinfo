package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
)

func GetGitHubTokens() []string {
	var tokens []string
	if token, exists := os.LookupEnv("TOKEN"); exists {
		tokens = append(tokens, token)
	}
	for i := 2; ; i++ {
		env := fmt.Sprintf("TOKEN%d", i)
		if token, exists := os.LookupEnv(env); exists {
			tokens = append(tokens, token)
		} else {
			break
		}
	}
	return tokens
}

func GetGitHubToken(tokens []string) (string, error) {
	if len(tokens) == 0 {
		return "", errors.New("no GitHub token available")
	}
	return tokens[rand.Intn(len(tokens))], nil
}


func GetGitHubTokenNative() (string, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return "", errors.New("GitHub token não encontrado")
	}
	return token, nil
}

func GetRandomGitHubToken(tokens []string) (string, error) {
	if len(tokens) == 0 {
		return "", errors.New("nenhum token disponível")
	}
	return tokens[rand.Intn(len(tokens))], nil
}
