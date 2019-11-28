package logic

import (
	"errors"

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
func LoginHandler(s *LogicService, obj *server.CallObject) {

	var uid = (obj.Args[0]).(string)
	var pwd = (obj.Args[1]).(string)
	result, err := server.Call(server.ServiceIDDB, common.CallIDGetUserInfoCmd, 0, uid)

	var ret = server.RetObject{
		Ret: result,
		Err: err,
	}
	if result == pwd {
		ret.Err = nil
	} else {
		ret.Err = errors.New("no user")
	}

	obj.ChanRet <- ret
}
