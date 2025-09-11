package auth

import (
	"fmt"
	"os"
	"errors"
	"math/rand"
)

func GetGitHubTokens() []string {
	var tokens []string
	if token, exists := os.LookupEnv("TOKEN"); exists {
		tokens = append(tokens, token)
	}
	for i := 2; ; i++ {
		envVar := fmt.Sprintf("TOKEN%d", i)
		if token, exists := os.LookupEnv(envVar); exists {
			tokens = append(tokens, token)
		} else {
			break
		}
	}
	return tokens
}


func GetGitHubTokenNative() (string, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return "", errors.New("GitHub token não encontrado")
	}
	return token, nil
}

func GetGitHubToken(tokens []string) (string, error) {
	if len(tokens) == 0 {
		return "", errors.New("no GitHub token available")
	}
	return tokens[rand.Intn(len(tokens))], nil
}
