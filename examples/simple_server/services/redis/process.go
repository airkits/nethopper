package redis

import (
	"fmt"

	"github.com/gonethopper/nethopper/server"
)

// GetUserInfoHander 获取用户信息
func GetUserInfoHander(s *RedisService, obj *server.CallObject) {
	var uid = (obj.Args[0]).(string)
	password, err := s.rdb.GetString(s.Context(), fmt.Sprintf("uid_%d", uid))
	var ret = server.RetObject{
		Ret: password,
		Err: err,
	}
	obj.ChanRet <- ret

}

func UpdateUserInfoHandler(s *RedisService, obj *server.CallObject) {
	var uid = (obj.Args[0]).(string)
	var password = (obj.Args[1]).(string)
	var key = fmt.Sprintf("uid_%d", uid)
	err := s.rdb.Set(s.Context(), key, password, 0)
	var ret = server.RetObject{
		Ret: password,
		Err: err,
	}
	obj.ChanRet <- ret
}
