package logic

import (
	"github.com/gonethopper/nethopper/examples/simple_server/common"
	"github.com/gonethopper/nethopper/server"
)

// func CreateUserHander(s *LogicService, obj *server.CallObject) {
// 	var uid = (obj.Args[0]).(string)

// 	var redisObj = server.NewCallObject(common.CallIDGetUserInfoCmd, uid)
// 	server.Call(server.ServiceIDRedis, 0, redisObj)
// 	result := <-redisObj.ChanRet
// 	if result.Err == nil {
// 		var ret = server.RetObject{
// 			Ret: result.Ret,
// 			Err: nil,
// 		}
// 		obj.ChanRet <- ret
// 		return
// 	}
// }
func LoginHandler(s *LogicService, obj *server.CallObject, uid string, pwd string) (string, error) {
	defer server.TraceCost("LoginHandler")()
	password, err := server.Call(server.ServiceIDRedis, common.CallIDGetUserInfoCmd, 0, uid)
	if err == nil {
		server.Info("get from redis")
		return password.(string), err
	}
	password, err = server.Call(server.ServiceIDDB, common.CallIDGetUserInfoCmd, 0, uid)
	if err != nil {
		return "", err
	}
	updated, err := server.Call(server.ServiceIDRedis, common.CallIDUpdateUserInfoCmd, 0, uid, password)
	if updated == false {
		server.Info("update redis failed %s %s", uid, password.(string))
	}
	server.Info("get from mysql")
	return password.(string), err
}
