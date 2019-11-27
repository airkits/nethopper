package db

import (
	"github.com/gonethopper/nethopper/server"
)

// GetUserInfoHander 获取用户信息
func GetUserInfoHander(s *DBService, obj *server.CallObject) {
	var uid = (obj.Args[0]).(string)

	// body := obj.Args
	sql := "select password from user.user where uid= ?"
	row := s.conn.QueryRow(sql, uid)
	var password string
	if err := row.Scan(&password); err == nil {
		server.Info(password)
	}
	var ret = &server.RetObject{
		Ret: password,
		Err: nil,
	}
	obj.ChanRet <- ret
	// m := server.CreateMessage(common.MessageIDLogin, s.ID(), req.SrcID, server.MTResponse, req.Cmd, req.SessionID)
	// body.Passwd = password
	// m.SetBody(body)
	// server.Call(m.DestID, 0, m)
}
