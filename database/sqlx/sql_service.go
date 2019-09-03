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

package sqlx

import (
	"time"

	"github.com/gonethopper/nethopper/server"
)

// SQLService struct to define service
type SQLService struct {
	server.BaseContext
	conn *SQLConnection
}

// SQLServiceCreate  service create function
func SQLServiceCreate() (server.Service, error) {

	return &SQLService{}, nil
}

// UserData service custom option, can you store you data and you must keep goruntine safe
func (s *SQLService) UserData() int32 {
	return 0
}

// Setup init custom service and pass config map to service
// config
// m := map[string]interface{}{
//  "queueSize":1000,
//  "driver:"mysql",
//  "dsn":"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
// }
func (s *SQLService) Setup(m map[string]interface{}) (server.Service, error) {

	conn, err := NewSQLConnection(m)
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
func (s *SQLService) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ServiceRun to call this function
func (s *SQLService) OnRun(dt time.Duration) {
	for i := 0; i < 128; i++ {
		m, err := s.MQ().AsyncPop()
		if err != nil {
			break
		}
		message := m.(*server.Message)
		s.ProcessMessage(message)
	}
}

// ProcessMessage receive message from mq and process message
func (s *SQLService) ProcessMessage(message *server.Message) {

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
func (s *SQLService) Stop() error {
	return nil
}

// PushMessage async send message to service
func (s *SQLService) PushMessage(option int32, msg *server.Message) error {
	return nil
}

// PushBytes async send string or bytes to queue
func (s *SQLService) PushBytes(option int32, buf []byte) error {
	return nil
}
