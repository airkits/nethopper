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
// * @Date: 2019-06-24 11:07:19
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-24 11:07:19

package services

import (
	"time"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/database/sqlx"
	"github.com/gonethopper/nethopper/server"
)

// DBService struct to define service
type DBService struct {
	server.BaseContext
	conn *sqlx.SQLConnection
}

// DBServiceCreate  service create function
func DBServiceCreate() (server.Service, error) {

	return &DBService{}, nil
}

// UserData service custom option, can you store you data and you must keep goruntine safe
func (s *DBService) UserData() int32 {
	return 0
}

// Setup init custom service and pass config map to service
// config
// m := map[string]interface{}{
//  "queueSize":1000,
//  "driver:"mysql",
//  "dsn":"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
// }
func (s *DBService) Setup(m map[string]interface{}) (server.Service, error) {

	conn, err := sqlx.NewSQLConnection(m)
	if err != nil {
		return nil, err
	}
	s.conn = conn
	if err := s.conn.Open(); err != nil {
		panic(err)
	}
	return s, nil
}

//Reload reload config
func (s *DBService) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ServiceRun to call this function
func (s *DBService) OnRun(dt time.Duration) {
	for i := 0; i < 128; i++ {
		m, err := s.MQ().AsyncPop()
		if err != nil {
			break
		}
		message := m.(*server.Message)
		server.Info("%s receive one request message from mq,cmd = %s", s.Name(), message.Cmd)

		s.ProcessMessage(message)
	}
}

// ProcessMessage receive message from mq and process message
func (s *DBService) ProcessMessage(message *server.Message) {
	cmd := message.Cmd
	if cmd == "login" {
		var v = make(map[string]interface{})
		server.Info("%s", string(message.Payload))
		if err := codec.JSONCodec.Unmarshal(message.Payload, &v, nil); err != nil {
			server.Info(err)
			return
		}
		sql := "select password from user.user where uid= ?"
		row := s.conn.QueryRow(sql, v["uid"])
		var pasword string
		if err := row.Scan(&pasword); err == nil {
			server.Info(pasword)
		}
		message.DestID = message.PopSrcID()
		message.MsgType = server.MTResponse
		message.Payload = []byte(pasword)
		server.SendMessage(message.DestID, 0, message)

	}
	// msgType := message.MsgType
	// switch msgType {
	// case server.MTRequest:
	// 	{
	// 		server.Info("receive message %s", message.Cmd)
	// 		message.SrcID = s.ID()

	// 		server.SendMessage(message.DestID, 0, message)
	// 		break
	// 	}
	// }
}

// Stop goruntine
func (s *DBService) Stop() error {
	return nil
}

// PushMessage async send message to service
func (s *DBService) PushMessage(option int32, msg *server.Message) error {
	if err := s.MQ().AsyncPush(msg); err != nil {
		server.Error(err.Error())
	}
	return nil
}

// PushBytes async send string or bytes to queue
func (s *DBService) PushBytes(option int32, buf []byte) error {
	return nil
}
