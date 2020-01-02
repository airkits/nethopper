package logic

import (
	"strconv"

	"github.com/gonethopper/nethopper/examples/model/common"
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

// LoginHandler user to login
// @Summary LoginHandler
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param uid query string true "UserID"
// @Param pwd query string true "Password"
// @Success 200 {string} string 成功后返回值
// @Router /call/LoginHandler [put]
func LoginHandler(s *Module, obj *server.CallObject, uid string, pwd string) (string, error) {
	defer server.TraceCost("LoginHandler")()
	opt, err := strconv.Atoi(uid)
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
