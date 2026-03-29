# gitinfo

> A Go library to fetch GitHub profile insights via the GitHub GraphQL API.

[![Go Reference](https://pkg.go.dev/badge/github.com/reinanbr/gitinfo.svg)](https://pkg.go.dev/github.com/reinanbr/gitinfo)
[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## Overview

`gitinfo` gives you a clean Go API over GitHub's GraphQL endpoint to pull user and repository data without dealing with pagination, query building, or response parsing yourself.

| Function | What it returns |
|---|---|
| `GetReposInfo` | Full repository data for a user |
| `GetReposName` | Repository names only |
| `GetLangPercents` | Language usage percentages across all repos |
| `GetCommits` | Total commits grouped by year and by day |
| `GetStreaks` | Max and current contribution streak with date ranges |

---

## Installation

```bash
go get github.com/reinanbr/gitinfo
```

**Requirements:** Go 1.21+ · GitHub Personal Access Token · Internet access to GitHub GraphQL API

---

## Quick Start

```go
package main

import (
    "fmt"
    "os"

    "github.com/reinanbr/gitinfo"
)

func main() {
    user  := "reinanbr"
    token := os.Getenv("GITHUB_TOKEN")

    repos, _ := gitinfo.GetReposInfo(user, token)
    fmt.Println("repos:", len(repos))

    langs, _ := gitinfo.GetLangPercents(user, token, []string{"Jupyter Notebook", "TeX"})
    fmt.Printf("top language: %s (%.1f%%)\n", langs.LangPercentages[0].Lang, langs.LangPercentages[0].Percentage)

    commits, _ := gitinfo.GetCommits(user, token)
    fmt.Println("total commits:", commits.TotalCommits)

    streaks, _ := gitinfo.GetStreaks(user, token)
    fmt.Println("streaks:", streaks)
}
```

---

## API Reference

### `GetReposInfo(user, token string) ([]RepoNode, error)`

Returns full repository metadata for a user.

```go
repos, err := gitinfo.GetReposInfo("reinanbr", token)
if err != nil {
    log.Fatal(err)
}
fmt.Println(len(repos), "repositories found")
```

---

### `GetReposName(user, token string) ([]RepoNode, error)`

Returns only repository names — lighter call for when you just need the list.

```go
names, err := gitinfo.GetReposName("reinanbr", token)
```

---

### `GetLangPercents(username, token string, ignoreLangs []string) (ResponseLangs, error)`

Returns language distribution across all repositories, sorted by usage.

```go
result, err := gitinfo.GetLangPercents("reinanbr", token, []string{"TeX", "Jupyter Notebook"})

fmt.Println("repos analyzed:", result.TotalRepos)
fmt.Println("total bytes:   ", result.TotalBytes)

for _, lp := range result.LangPercentages {
    fmt.Printf("  %-20s %.2f%%\n", lp.Lang, lp.Percentage)
}
```

**Response type:**
```go
type ResponseLangs struct {
    LangPercentages []LangPercentage
    TotalBytes      int
    TotalRepos      int
}

type LangPercentage struct {
    Lang       string
    Percentage float64
}
```

---

### `GetCommits(username, token string) (CommitsResponse, error)`

Returns total commits with two views: grouped by year and flat daily list.  
Future calendar days are automatically filtered out.

```go
commits, err := gitinfo.GetCommits("reinanbr", token)

fmt.Println("total:", commits.TotalCommits)
for _, year := range commits.CommitsByYear {
    fmt.Printf("  %d: %d commits\n", year.Year, len(year.Commits))
}
```

**Response type:**
```go
type CommitsResponse struct {
    User          string
    TotalCommits  int
    CommitsByYear []CommitsYear
    CommitsByDay  []CommitByDate
}
```

---

### `GetStreaks(username, token string) (map[string]interface{}, error)`

Returns max and current contribution streaks with start/end dates.

```go
result, err := gitinfo.GetStreaks("reinanbr", token)
```

**Example response:**
```json
{
  "user": "reinanbr",
  "streak": {
    "max_streak": 46,
    "current_streak": 3,
    "max_streak_period":     { "start": "2022-12-25", "end": "2023-02-08" },
    "current_streak_period": { "start": "2025-03-27", "end": "2025-03-29" }
  }
}
```

---

## Testing

Tests are integration tests and require a valid token. Create a `.env` file first:

```env
GITHUB_TOKEN=your_token_here
```

```bash
# Run all tests
go test -v

# Run a specific test
go test -v -run TestGetCommits
```

---

## Smoke Test CLI

A CLI at `cmd/smoke` validates all endpoints end-to-end before you publish or deploy.

```bash
# Run with env token
go run ./cmd/smoke -user reinanbr

# Run with explicit token
go run ./cmd/smoke -user reinanbr -token YOUR_TOKEN

# All flags
go run ./cmd/smoke -user reinanbr -ignore "Jupyter Notebook,TeX" -show-days
```

| Flag | Description | Default |
|---|---|---|
| `-user` | GitHub username **(required)** | — |
| `-token` | GitHub token (falls back to `GITHUB_TOKEN`) | — |
| `-ignore` | Comma-separated languages to ignore | `Jupyter Notebook,TeX` |
| `-show-days` | Print all `commitsByDay` entries | `false` |

---

## Notes

- Responses depend on GitHub API availability and your token's rate limits.
- Never commit your `.env` file — add it to `.gitignore`.
- Data is fetched from 2015 onwards by default.

---

## License

MIT © [reinanbr](https://github.com/reinanbr)