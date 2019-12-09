package redis

import (
	"fmt"

	"github.com/gonethopper/nethopper/server"
)

// GetUserInfoHander 获取用户信息
func GetUserInfoHander(s *RedisModule, obj *server.CallObject, uid string) (string, error) {
	defer server.TraceCost("GetUserInfoHander")()
	password, err := s.rdb.GetString(s.Context(), fmt.Sprintf("uid_%d", uid))
	return password, err

}

// UpdateUserInfoHandler update user info
func UpdateUserInfoHandler(s *RedisModule, obj *server.CallObject, uid string, pwd string) (bool, error) {

	var key = fmt.Sprintf("uid_%d", uid)
	err := s.rdb.Set(s.Context(), key, pwd, 0)
	if err != nil {
		return false, err
	}
	return true, nil
}
