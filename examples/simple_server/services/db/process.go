package db

import (
	"github.com/gonethopper/nethopper/server"
)

// GetUserInfoHander 获取用户信息
func GetUserInfoHander(s *DBService, obj *server.CallObject, u string) (string, error) {
	//var uid = (obj.Args[0]).(string)
	//uid := 1
	sql := "select password from user where uid= ?"
	row := s.conn.QueryRow(sql, u)
	var password string
	if err := row.Scan(&password); err == nil {
		server.Info(password)
	}
	return password, nil
}

// InsertUserInfoHander 获取用户信息
// func InsertUserInfoHander(s *DBService, obj *server.CallObject) {
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
// func UpdateUserInfoHander(s *DBService, obj *server.CallObject) {
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
