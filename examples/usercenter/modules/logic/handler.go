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
	"github.com/gonethopper/nethopper/examples/usercenter/cmd"
	"github.com/gonethopper/nethopper/examples/usercenter/global"
	"github.com/gonethopper/nethopper/examples/usercenter/model"
	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/nethopper/utils"
)

// WXLoginHandler weixin user login
// @Summary WXLoginHandler
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param appID query string true "appID"
// @Param code query string true "code"
// @Success 200 {string} string 成功后返回值
// @Router /call/WXLoginHandler [put]
func WXLoginHandler(s *Module, obj *server.CallObject, appID string, code string) (*model.User, error) {
	defer server.TraceCost("LoginHandler")()
	info, err := server.Call(global.ModuleIDWechatClient, cmd.MCWXLogin, utils.RandomInt32(0, 1024), appID, code)
	if err != nil {
		return nil, err
	}
	server.Info("%v", info)
	return &model.User{
		UID: 1111,
	}, nil
	// user, err := server.Call(server.ModuleIDRedis, common.CallIDGetUserInfoCmd, appID, code)
	// if err == nil {
	// 	server.Info("get from redis")
	// 	return password.(string), err
	// }
	// password, err = server.Call(server.ModuleIDDB, common.CallIDGetUserInfoCmd, int32(opt), uid)
	// if err != nil {
	// 	return "", err
	// }
	// updated, err := server.Call(server.ModuleIDRedis, common.CallIDUpdateUserInfoCmd, int32(opt), uid, password)
	// if updated == false {
	// 	server.Info("update redis failed %s %s", uid, password.(string))
	// }
	// server.Info("get from mysql")
	//return nil, errors.New("no user")
}
