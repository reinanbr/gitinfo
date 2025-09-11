# GitInfo

Uma biblioteca Go para interagir com a API do GitHub e obter informações sobre usuários, repositórios, contribuições e linguagens.

## Estrutura do Projeto

```
gitinfo/
├── pkg/
│   ├── auth/           # Gerenciamento de tokens do GitHub
│   ├── github/         # Informações de usuários e repositórios
│   ├── graphql/        # Queries e respostas GraphQL
│   ├── languages/      # Análise de linguagens de programação
│   └── utils/          # Funções utilitárias
├── example/            # Exemplos de uso
├── gitinfo.go          # Arquivo principal da biblioteca
├── go.mod              # Módulo Go
└── README.md           # Este arquivo
```

## Instalação

```bash
go get github.com/reinanbr/gitinfo
```

## Uso

### Autenticação

```go
import "github.com/reinanbr/gitinfo/pkg/auth"

// Obter token do ambiente
token, err := auth.GetGitHubTokenNative()

// Obter todos os tokens disponíveis
tokens := auth.GetGitHubTokens()

// Obter um token aleatório da lista
token, err := auth.GetGitHubToken(tokens)
```

### Informações do Usuário

```go
import "github.com/reinanbr/gitinfo/pkg/github"

// Obter informações do usuário
userInfo, err := github.GetUserInfo("username", token)
```

### Repositórios

```go
import "github.com/reinanbr/gitinfo/pkg/utils"

// Buscar todos os repositórios de um usuário
repos, err := utils.FetchAllRepos("username", token, nil)
```

### Linguagens

```go
import "github.com/reinanbr/gitinfo/pkg/languages"

// Calcular percentuais de linguagens
percentages, total, err := languages.CalculateLanguagePercentages("username", tokens)
```

## Variáveis de Ambiente

- `TOKEN`: Token principal do GitHub
- `TOKEN2`, `TOKEN3`, etc.: Tokens adicionais para balanceamento de carga

## Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
