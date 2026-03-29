# gitinfo

A Go library to fetch GitHub profile and repository insights using the GitHub GraphQL API.

It provides helpers for:
- Repositories (full info and names)
- Language percentages
- Contribution streaks (with periods)
- Commits (by year and by day)

## Features

- Fetch all repositories for a user
- Fetch only repository names
- Calculate language usage percentages across repositories
- Return max/current contribution streak plus date ranges
- Return total commits plus:
  - `commitsByYear` (grouped)
  - `commitsByDay` (flat daily list, up to current date)

## Installation

```bash
go get github.com/reinanbr/gitinfo
```

## Requirements

- Go 1.21+
- GitHub Personal Access Token in `GITHUB_TOKEN`
- Internet access to GitHub GraphQL API

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/reinanbr/gitinfo"
)

func main() {
    token := "YOUR_GITHUB_TOKEN"
    user := "reinanbr"

    repos, err := gitinfo.GetReposInfo(user, token)
    if err != nil {
        panic(err)
    }
    fmt.Println("repos:", len(repos))

    langs, err := gitinfo.GetLangPercents(user, token, []string{"Jupyter Notebook", "TeX"})
    if err != nil {
        panic(err)
    }
    fmt.Println("total bytes:", langs.TotalBytes)

    streaks, err := gitinfo.GetStreaks(user, token)
    if err != nil {
        panic(err)
    }
    fmt.Println("streak response:", streaks)

    commits, err := gitinfo.GetCommits(user, token)
    if err != nil {
        panic(err)
    }
    fmt.Println("total commits:", commits.TotalCommits)
}
```

## Public API

### `GetReposInfo(user, token)`

Returns detailed repository data.

```go
repos, err := gitinfo.GetReposInfo("reinanbr", token)
```

### `GetReposName(user, token)`

Returns repository names.

```go
repoNames, err := gitinfo.GetReposName("reinanbr", token)
```

### `GetLangPercents(username, token, ignoreLangs)`

Returns percentages and totals:
- `LangPercentages []LangPercentage`
- `TotalBytes int`
- `TotalRepos int`

```go
result, err := gitinfo.GetLangPercents("reinanbr", token, []string{"TeX"})
```

### `GetStreaks(username, token)`

Returns:

```json
{
  "user": "reinanbr",
  "streak": {
    "max_streak": 46,
    "current_streak": 0,
    "max_streak_period": { "start": "2022-12-25", "end": "2023-02-08" },
    "current_streak_period": { "start": "", "end": "" }
  }
}
```

### `GetCommits(username, token)`

Returns:
- `User`
- `TotalCommits`
- `CommitsByYear []CommitsYear`
- `CommitsByDay []CommitByDate`

`CommitsByDay` includes daily commit counts up to the current date (future calendar days are filtered out).

## Run Tests

Tests use `.env` and expect `GITHUB_TOKEN`:

```bash
go test -v
```

Run a specific test:

```bash
go test -v -run TestGetCommits
```

## Smoke Test CLI

A CLI is included at `cmd/smoke/main.go` to quickly validate all endpoints before publishing.

### 1) Create `.env`

```env
GITHUB_TOKEN=your_token_here
```

### 2) Run smoke test

```bash
go run ./cmd/smoke -user reinanbr
```

### 3) Optional flags

```bash
go run ./cmd/smoke -user reinanbr -show-days=false -ignore "Jupyter Notebook,TeX"
```

Flags:
- `-user` GitHub username (required)
- `-token` GitHub token (optional if `GITHUB_TOKEN` exists in env/.env)
- `-ignore` comma-separated ignored languages
- `-show-days` print all `commitsByDay` entries

## Notes

- API responses depend on GitHub availability and rate limits.
- Some tests are integration tests and require a valid token.
- Keep your token private. Do not commit `.env`.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE).
