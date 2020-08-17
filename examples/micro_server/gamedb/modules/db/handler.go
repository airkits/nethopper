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
	"github.com/gonethopper/nethopper/server"
)

// GetUser 获取用户信息
func GetUser(s *Module, obj *server.CallObject, u string) (string, server.Result) {

	//var uid = (obj.Args[0]).(string)
	//uid := 1
	sql := "select password from user where uid= ?"
	row := s.conn.QueryRow(sql, u)
	var password string
	var err error
	if err = row.Scan(&password); err == nil {
		return password, server.Result{Code: 0, Err: nil}
	}

	return "", server.Result{Code: -1, Err: nil}

}

// InsertUserInfoHander 获取用户信息
// func InsertUserInfoHander(s *Module, obj *server.CallObject) {
// 	var uid = (obj.Args[0]).(string)
// 	var password = (obj.Args[1]).(string)
// 	sql := "insert into user.user(uid,password) value(?,?)"
// 	_, err := s.conn.Exec(sql, password, uid)

// 	if err == nil {
// 		var ret = server.RetObject{
// 			Ret: password,
// 			Err: nil,
// 		}
// 		obj.ChanRet <- ret
// 	} else {
// 		var ret = server.RetObject{
// 			Ret: nil,
// 			Err: err,
// 		}
// 		obj.ChanRet <- ret
// 	}

// }

// UpdateUserInfoHander 获取用户信息
// func UpdateUserInfoHander(s *Module, obj *server.CallObject) {
// 	var uid = (obj.Args[0]).(string)
// 	var password = (obj.Args[1]).(string)
// 	sql := "update user.user set password=? where uid=?"
// 	_, err := s.conn.Exec(sql, password, uid)

// 	if err == nil {
// 		var ret = server.RetObject{
// 			Ret: password,
// 			Err: nil,
// 		}
// 		obj.ChanRet <- ret
// 	} else {
// 		var ret = server.RetObject{
// 			Ret: nil,
// 			Err: err,
// 		}
// 		obj.ChanRet <- ret
// 	}

// }
