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
	"time"

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
// @Param channel query string true "channel"
// @Param nickname query string true "nickname"
// @Param gender query int true "gender"
// @Param avatar query string true "avatar"
// @Success 200 {string} string 成功后返回值
// @Router /call/WXLogin [put]
func WXLogin(s *Module, obj *server.CallObject, appID string, code string, channel string, nickname string, gender int, avatar string) (*model.User, server.Result) {
	defer server.TraceCost(server.RunModuleFuncName(s))()
	user := &model.User{
		UID:     0,
		AppID:   appID,
		Channel: channel,
		Name:    nickname,
		Gender:  gender,
		Avatar:  avatar,
		Gold:    10000,
		Coin:    10000,
	}
	wxuser, result := server.Call(global.MIDWechatClient, cmd.WXLogin, utils.RandomInt32(0, 1024), appID, code)
	if result.Err != nil {
		return nil, result
	}
	wxu := wxuser.(*model.WXUser)
	server.Info("[%s] get openID %s", s.Name(), wxu.OpenID)
	user.OpenID = wxu.OpenID

	uid, result := server.Call(server.MIDLogic, cmd.LogicGetUIDByOpenID, utils.RandomInt32(0, 1024), wxu.OpenID)
	if result.Err != nil {
		server.Info("[%s] get uid error %s", s.Name(), result.Err.Error())

		user.UID = uid.(uint64)
		u, result := server.Call(server.MIDLogic, cmd.LogicGetUser, int32(user.UID), user)
		if result.Err != nil {
			server.Error(result.Err.Error())
			return nil, result
		}
		server.Info("[%s] get from snowflake success, uid %d", s.Name(), uid)
		return u.(*model.User), result
	}
	user.UID = uid.(uint64)

	return user, result

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
func GetUIDByOpenID(s *Module, obj *server.CallObject, openID string) (uint64, server.Result) {
	defer server.TraceCost(server.RunModuleFuncName(s))()

	uid, result := server.Call(server.MIDRedis, cmd.RedisGetUIDByOpenID, utils.RandomInt32(0, 1024), openID)
	if result.Err == nil {
		server.Info("[%s] get from redis uid=%d", s.Name(), uid)
		return uid.(uint64), result
	}

	uid, result = server.Call(server.MIDDB, cmd.DBGetUIDByOpenID, utils.RandomInt32(0, 1024), openID)
	if result.Err == nil {
		server.Info("[%s] get uid from db uid=%d", s.Name(), uid)
		return uid.(uint64), result
	}

	uid, result = server.Call(global.MIDSnowflakeClient, cmd.SFGetUID, utils.RandomInt32(0, 1024), int32(0))
	if result.Err != nil {
		server.Info("[%s] get from snowflake error %s", s.Name(), result.Err.Error())
		return 0, result
	}

	updated, result := server.Call(server.MIDDB, cmd.DBInsertOID2UID, int32(uid.(uint64)), openID, uid)
	if updated == false {
		server.Info("[%s] set db failed %s %d", s.Name(), openID, uid)
		return 0, result
	}
	updated, result = server.Call(server.MIDRedis, cmd.RedisSetUIDByOpenID, int32(uid.(uint64)), openID, uid)
	if updated == false {
		server.Info("[%s] update redis failed %s %d", s.Name(), openID, uid)
	}
	return uid.(uint64), result
}

// GetUser get user
// @Summary GetUser
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param user query string true "user"
// @Success 200 {string} string 成功后返回值
// @Router /call/v [put]
func GetUser(s *Module, obj *server.CallObject, uid uint64) (*model.User, server.Result) {
	defer server.TraceCost(server.RunModuleFuncName(s))()
	user, result := server.Call(server.MIDRedis, cmd.DBGetUserByUID, int32(uid), uid)
	if result.Err == nil {
		return user.(*model.User), result
	}
	user, result = server.Call(server.MIDDB, cmd.DBGetUserByUID, int32(uid), uid)
	if result.Err == nil {
		server.Call(server.MIDRedis, cmd.RedisUpdateUserInfo, int32(uid), user)
		return user.(*model.User), result
	}
	return nil, result
}

// CreateUser create user
// @Summary CreateUser
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param user query string true "user"
// @Success 200 {string} string 成功后返回值
// @Router /call/CreateUser [put]
func CreateUser(s *Module, obj *server.CallObject, user *model.User) (*model.User, server.Result) {
	defer server.TraceCost(server.RunModuleFuncName(s))()
	user.CreateAt = time.Now()
	user.LoginAt = time.Now()
	v, result := server.Call(server.MIDDB, cmd.DBCreateUser, int32(user.UID), user)
	if result.Err != nil {
		return nil, result
	}
	updated, result := server.Call(server.MIDRedis, cmd.RedisUpdateUserInfo, int32(user.UID), v)
	if updated == false {
		server.Info("update redis failed %s %d  err:%s", user.OpenID, user.UID, result.Err.Error())
	}
	return v.(*model.User), result
}
