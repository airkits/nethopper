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

// WXLogin weixin user login
// @Summary WXLogin
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param appID query string true "appID"
// @Param code query string true "code"
// @Success 200 {string} string 成功后返回值
// @Router /call/WXLogin [put]
func WXLogin(s *Module, obj *server.CallObject, appID string, code string) (*model.User, error) {
	defer server.TraceCost("LoginHandler")()
	wxuser, err := server.Call(global.ModuleIDWechatClient, cmd.MCWXLogin, utils.RandomInt32(0, 1024), appID, code)
	if err != nil {
		return nil, err
	}
	server.Info("%v", wxuser)
	uid, err := GetUIDByOpenID(s, obj, wxuser.(*model.WXUser).OpenID)
	if err != nil {
		server.Info("get uid error %s", err.Error())
		
	}

	return &model.User{
		UID: uid,
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

// GetUIDByOpenID convert openid to uid
// @Summary GetUIDByOpenID
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param openID query string true "appID"
// @Success 200 {string} string 成功后返回值
// @Router /call/GetUIDByOpenID [put]
func GetUIDByOpenID(s *Module, obj *server.CallObject, openID string) (uint64, error) {
	defer server.TraceCost("GetUIDByOpenID")()

	uid, err := server.Call(server.ModuleIDRedis, cmd.MCRedisGetUIDByOpenID, utils.RandomInt32(0, 1024), openID)
	if err == nil {
		server.Info("get from redis uid=%ld", uid)
		return uid.(uint64), err
	}

	uid, err = server.Call(server.ModuleIDDB, cmd.MCDBGetUIDByOpenID, utils.RandomInt32(0, 1024), openID)
	if err != nil {
		return 0, err
	}

	updated, err := server.Call(server.ModuleIDRedis, cmd.MCRedisSetUIDByOpenID, utils.RandomInt32(0, 1024), openID, uid)
	if updated == false {
		server.Info("update redis failed %s %ld", openID, uid)
	}
	return uid.(uint64), err
}
