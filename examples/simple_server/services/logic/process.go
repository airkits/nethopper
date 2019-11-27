package logic

import "github.com/gonethopper/nethopper/server"

import "github.com/gonethopper/nethopper/examples/simple_server/common"

func LoginHandler(obj *server.CallObject) {

	var uid = (obj.Args[0]).(string)
	var retChan = make(chan *server.RetObject, 1)
	var callObj = server.NewCallObject(common.CallIDGetUserInfoCmd, retChan, uid)
	server.Call(server.ServiceIDDB, 0, callObj)

	result := <-retChan
	if result.Err != nil {
		var ret = &server.RetObject{
			Ret: result.Ret,
			Err: nil,
		}

		obj.ChanRet <- ret
	}
}
