package redis

import (
	"fmt"

	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/nethopper/utils"
)

// GetUserInfoHander 获取用户信息
func GetUserInfoHander(s *RedisService, obj *server.CallObject, uid string) (string, error) {
	defer utils.Trace("GetUserInfoHander")()
	password, err := s.rdb.GetString(s.Context(), fmt.Sprintf("uid_%d", uid))
	return password, err

}

// UpdateUserInfoHandler update user info
func UpdateUserInfoHandler(s *RedisService, obj *server.CallObject, uid string, pwd string) (bool, error) {

	var key = fmt.Sprintf("uid_%d", uid)
	err := s.rdb.Set(s.Context(), key, pwd, 0)
	if err != nil {
		return false, err
	}
	return true, nil
}
