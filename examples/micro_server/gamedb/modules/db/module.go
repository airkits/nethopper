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

package db

import (
	"time"

	"github.com/gonethopper/nethopper/database/sqlx"
	"github.com/gonethopper/nethopper/examples/micro_server/gamedb/cmd"
	"github.com/gonethopper/nethopper/server"
)

// Module struct to define module
type Module struct {
	server.BaseContext
	conn *sqlx.SQLConnection
}

// ModuleCreate  module create function
func ModuleCreate() (server.Module, error) {
	return &Module{}, nil
}

//Handlers set moudle handlers
func (s *Module) Handlers() map[string]interface{} {
	return map[string]interface{}{}
}

//ReflectHandlers set moudle reflect handlers
func (s *Module) ReflectHandlers() map[string]interface{} {
	return map[string]interface{}{
		cmd.DBGetUser: GetUser,
	}
}

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
//  "driver:"mysql",
//  "dsn":"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
// }
func (s *Module) Setup(conf server.IConfig) (server.Module, error) {
	conn, err := sqlx.NewSQLConnection(conf)
	if err != nil {
		return nil, err
	}
	s.conn = conn
	if err := s.conn.Open(); err != nil {
		panic(err)
	}
	s.CreateWorkerPool(s, 128, 10*time.Second, true)
	return s, nil
}

//Reload reload config
// func (s *Module) Reload(m map[string]interface{}) error {
// 	return nil
// }

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
	server.RunSimpleFrame(s, 128)
}

// Stop goruntine
func (s *Module) Stop() error {
	return nil
}

// Call async send message to module
// func (s *Module) Call(option int32, obj *server.CallObject) error {
// 	if err := s.MQ().AsyncPush(obj); err != nil {
// 		server.Error(err.Error())
// 	}
// 	return nil
// }

// PushBytes async send string or bytes to queue
// func (s *Module) PushBytes(option int32, buf []byte) error {
// 	return nil
// }
