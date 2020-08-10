// MIT License

// Copyright (c) 2019 gonethopper

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * @Author: ankye
// * @Date: 2020-01-09 11:01:34
// * @Last Modified by:   ankye
// * @Last Modified time: 2020-01-09 11:01:34

package logic

import (
	"strconv"

	"github.com/gonethopper/nethopper/examples/micro_server/gamedb/cmd"
	"github.com/gonethopper/nethopper/server"
)

// func CreateUserHander(s *Module, obj *server.CallObject) {
// 	var uid = (obj.Args[0]).(string)

// 	var redisObj = server.NewCallObject(common.CallIDGetUserInfoCmd, uid)
// 	server.Call(server.MIDRedis, 0, redisObj)
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

	password, err := server.Call(server.MIDDB, cmd.CallIDGetUserInfoCmd, int32(opt), uid)
	if err != nil {
		return "", err
	}
	server.Info("get from mysql")
	return password.(string), err
}
