package utils

type GraphQLQuery struct {
	Query string `json:"query"`
}

type RepoNode struct {
	Name             string `json:"name"`
	CreatedAt        string `json:"createdAt"`
	Description      string `json:"description"`
	Url              string `json:"url"`
	IsPrivate        bool   `json:"isPrivate"`
	DefaultBranchRef *struct {
		Target struct {
			CommittedDate string `json:"committedDate"`
		} `json:"target"`
	} `json:"defaultBranchRef"`
	Languages struct {
		Edges []struct {
			Size int `json:"size"`
			Node struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"languages"`
	LastCommitDate string `json:"lastCommitDate"`
}

type UserInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatarUrl"`
	CreatedAt string `json:"createdAt"`
	URL       string `json:"url"`
	Followers struct {
		TotalCount int `json:"totalCount"`
	} `json:"followers"`
	Following struct {
		TotalCount int `json:"totalCount"`
	} `json:"following"`
	Repositories struct {
		TotalCount int `json:"totalCount"`
	} `json:"repositories"`
}

type RepoResponse struct {
	Data struct {
		User struct {
			Repositories struct {
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
				Nodes []RepoNode `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type UserResponse struct {
	Data struct {
		User UserInfo `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type ContributionGraphQuery struct {
	Query string `json:"query"`
}

type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
}

type Week struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

type ContributionCalendar struct {
	Weeks []Week `json:"weeks"`
}

type ContributionsCollection struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
}

type User struct {
	CreatedAt               string                  `json:"createdAt"`
	ContributionsCollection ContributionsCollection `json:"contributionsCollection"`
}

type Data struct {
	User User `json:"user"`
}

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	Data   Data    `json:"data"`
	Errors []Error `json:"errors"`
}

// Language represents a programming language.
type Language struct {
	Name string `json:"name"`
}

// LanguageEdge represents the size of a language in a repository.
type LanguageEdge struct {
	Size int      `json:"size"`
	Node Language `json:"node"`
}

// Repository represents a GitHub repository with language and metadata.
type Repository struct {
	Name       string `json:"name"`
	DateCreate string `json:"createdAt"`
	Languages  struct {
		Edges []LanguageEdge `json:"edges"`
	} `json:"languages"`
	DefaultBranchRef struct {
		Target struct {
			CommittedDate string `json:"committedDate"`
		} `json:"target"`
	} `json:"defaultBranchRef"`
	LastCommitDate string `json:"lastCommitDate"`
}

// Repo represents a collection of repositories.
type Repo struct {
	Repositories struct {
		PageInfo struct {
			HasNextPage bool   `json:"hasNextPage"`
			EndCursor   string `json:"endCursor"`
		} `json:"pageInfo"`
		Nodes []Repository `json:"nodes"`
	} `json:"repositories"`
}

// ResponseLangs represents the GraphQL response for language data.
type ResponseLangs struct {
	Data struct {
		Repo Repo `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// RepositoryLite represents a lightweight repository structure.
type RepositoryLite struct {
	Name string `json:"name"`
}

// RepoName represents a collection of lightweight repositories.
type RepoName struct {
	Repositories struct {
		Nodes []RepositoryLite `json:"nodes"`
	} `json:"repositories"`
}

// ResponseLite represents the GraphQL response for lightweight repository data.
type ResponseLite struct {
	Data struct {
		Repo RepoName `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}
