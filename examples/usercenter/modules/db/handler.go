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
	"github.com/gonethopper/nethopper/utils/crypto/md5"
)

func getUserTableByUID(uid uint64) string {
	return fmt.Sprintf("usercenter.user_%d", uid%8)
}
func getOID2UIDTable(openID string) string {
	num := md5.HashMod(openID, 8)
	return fmt.Sprintf("usercenter.oid2uid_%d", num)
}

// GetUIDByOpenID 获取uid
func GetUIDByOpenID(s *Module, obj *server.CallObject, openID string) (uint64, error) {
	sql := fmt.Sprintf("select uid from %s where openid=?", getOID2UIDTable(openID))
	row := s.conn.QueryRow(sql, openID)
	var uid uint64
	var err error
	if err = row.Scan(&uid); err == nil {
		return uid, nil
	}
	return 0, err
}

//InsertOID2UID insert oid and uid in mapping
func InsertOID2UID(s *Module, obj *server.CallObject, openID string, uid uint64) (bool, error) {
	sql := fmt.Sprintf("insert into %s (openid,uid) value(?,?)", getOID2UIDTable(openID))
	if _, err := s.conn.Exec(sql, openID, uid); err != nil {
		return false, err
	}
	return true, nil
}

// GetUserByUID 获取用户信息
func GetUserByUID(s *Module, obj *server.CallObject, uid uint64) (*model.User, error) {
	sql := fmt.Sprintf("select uid,appid,openid,uuid,avatar,name,password,phone,gender,age,gold,coin,loginat,createat,status,loginip,channel from %s where uid= ?", getUserTableByUID(uid))
	user := model.User{
		UID: uid,
	}
	var err error
	if err = s.conn.Select(&user, sql, uid); err == nil {
		return &user, nil
	}
	return nil, err
}

// CreateUser 创建用户信息
func CreateUser(s *Module, obj *server.CallObject, u *model.User) (*model.User, error) {
	sql := fmt.Sprintf("insert into %s(uid,appid,openid,uuid,avatar,name,password,phone,gender,age,gold,coin,status,channel,loginip,loginat,createat) value(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", getUserTableByUID(u.UID))
	_, err := s.conn.Exec(sql, u.UID, u.AppID, u.OpenID, u.UUID, u.Avatar, u.Name, u.Password, u.Phone, u.Gender, u.Age, u.Gold, u.Coin, u.Status, u.Channel, u.LoginIP, u.LoginAt, u.CreateAt)

	if err == nil {
		return u, nil
	}
	return nil, err
}
