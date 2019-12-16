package logic

import (
	"strconv"

	"github.com/gonethopper/nethopper/examples/simple_server/common"
	"github.com/gonethopper/nethopper/server"
)

// func CreateUserHander(s *Module, obj *server.CallObject) {
// 	var uid = (obj.Args[0]).(string)

// 	var redisObj = server.NewCallObject(common.CallIDGetUserInfoCmd, uid)
// 	server.Call(server.ModuleIDRedis, 0, redisObj)
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
func LoginHandler(s *Module, obj *server.CallObject, uid string, pwd string) (string, error) {
	defer server.TraceCost("LoginHandler")()
	opt, err := strconv.Atoi(uid)
	server.Info("get opt %d", opt)
	password, err := server.Call(server.ModuleIDRedis, common.CallIDGetUserInfoCmd, int32(opt), uid)
	if err == nil {
		server.Info("get from redis")
		return password.(string), err
	}
	password, err = server.Call(server.ModuleIDDB, common.CallIDGetUserInfoCmd, int32(opt), uid)
	if err != nil {
		return "", err
	}
	updated, err := server.Call(server.ModuleIDRedis, common.CallIDUpdateUserInfoCmd, int32(opt), uid, password)
	if updated == false {
		server.Info("update redis failed %s %s", uid, password.(string))
	}
	server.Info("get from mysql")
	return password.(string), err
}
