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
	"github.com/gonethopper/nethopper/examples/simple_server/common"
	"github.com/gonethopper/nethopper/server"
)

// DBModule struct to define module
type DBModule struct {
	server.BaseContext
	conn *sqlx.SQLConnection
}

// DBModuleCreate  module create function
func DBModuleCreate() (server.Module, error) {

	return &DBModule{}, nil
}

// UserData module custom option, can you store you data and you must keep goruntine safe
func (s *DBModule) UserData() int32 {
	return 0
}

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
//  "driver:"mysql",
//  "dsn":"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
// }
func (s *DBModule) Setup(m map[string]interface{}) (server.Module, error) {
	s.RegisterHandler(common.CallIDGetUserInfoCmd, GetUserInfoHander)
	conn, err := sqlx.NewSQLConnection(m)
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
func (s *DBModule) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *DBModule) OnRun(dt time.Duration) {
	server.RunSimpleFrame(s)
}

// Stop goruntine
func (s *DBModule) Stop() error {
	return nil
}

// Call async send message to module
// func (s *DBModule) Call(option int32, obj *server.CallObject) error {
// 	if err := s.MQ().AsyncPush(obj); err != nil {
// 		server.Error(err.Error())
// 	}
// 	return nil
// }

// PushBytes async send string or bytes to queue
func (s *DBModule) PushBytes(option int32, buf []byte) error {
	return nil
}
