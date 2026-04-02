package gitinfo

import (
	"github.com/reinanbr/gitinfo/internal/graphql/fetch"
	"github.com/reinanbr/gitinfo/internal/utils"
)


func GetUserInfo(username string, token string) (utils.UserInfo, error) {
	return fetch.FetchUserInfo(username, token)
}