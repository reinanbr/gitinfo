package gitinfo

import (
	"github.com/reinanbr/gitinfo/pkg/graphql/fetch"
	"github.com/reinanbr/gitinfo/pkg/utils"
)

func GetReposInfo(user string, token string) ([]utils.RepoNode, error) {
	repos, err := fetch.FetchAllRepos(user, token, nil)
	if err != nil {
		return nil, err
	}
	return repos, err
}

func GetRepos(user string, token string) ([]utils.RepoNode, error) {
	repos, err := fetch.FetchAllReposName(user, token, nil, 0)
	if err != nil {
		return nil, err
	}
	return repos, err
}
