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
// * @Date: 2020-01-09 11:01:18
// * @Last Modified by:   ankye
// * @Last Modified time: 2020-01-09 11:01:18

package db

import (
	"fmt"

	"github.com/gonethopper/nethopper/examples/usercenter/model"
	"github.com/gonethopper/nethopper/server"
)

// GetUserInfoByOpenIDHander 获取用户信息
func GetUserInfoByOpenIDHander(s *Module, obj *server.CallObject, appID string, openID string) (*model.User, error) {

	sql := "select uid,appid,openid,uuid,avatar,name,password,phone,gender,age,gold,coin,loginat,createat,status,loginip,channel from user where appid= ? and openid= ?"
	user := model.User{
		AppID:  appID,
		OpenID: openID,
	}
	var err error
	if err = s.conn.Select(&user, sql, appID, openID); err == nil {
		return &user, nil
	}
	return nil, err
}
func getTableByUID(uid uint64) string {
	return fmt.Sprintf("usercenter.user_%d", uid%8)
}

// CreateUserInfoHander 创建用户信息
func CreateUserInfoHander(s *Module, obj *server.CallObject, u *model.User) (*model.User, error) {
	sql := "insert into " + getTableByUID(u.UID) + "(uid,appid,openid,uuid,avatar,name,password,phone,gender,age,gold,coin,status,channel,loginip,loginat,createat) value(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := s.conn.Exec(sql, u.UID, u.AppID, u.OpenID, u.UUID, u.Avatar, u.Name, u.Password, u.Phone, u.Gender, u.Age, u.Gold, u.Coin, u.Status, u.Channel, u.LoginIP, u.LoginAt, u.CreateAt)

	if err == nil {
		return u, nil
	}
	return nil, err
}
