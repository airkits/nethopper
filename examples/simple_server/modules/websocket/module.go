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

package websocket

import (
	"time"

	"github.com/gonethopper/nethopper/server"
)

// HTTPTimeout http timeout (second)
const HTTPTimeout = 10

// ModuleCreate  module create function
func ModuleCreate() (server.Module, error) {
	return &Module{}, nil
}

// Module struct to define module
type Module struct {
	server.BaseContext
	Address  string
	CertFile string
	KeyFile  string
}

// UserData module custom option, can you store you data and you must keep goruntine safe
func (s *Module) UserData() int32 {
	return 0
}

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *Module) Setup(m map[string]interface{}) (server.Module, error) {
	if err := s.readConfig(m); err != nil {
		panic(err)
	}

	var wsServer *WSServer
	if s.Address != "" {
		wsServer = new(WSServer)
		wsServer.WsConfig = WsConfig{
			Address:         s.Address,
			MaxConnNum:      1024,
			PendingWriteNum: 1024,
			MaxMsgLen:       4096,
			HTTPTimeout:     HTTPTimeout,
			CertFile:        s.CertFile,
			KeyFile:         s.KeyFile,
		}

		wsServer.NewAgent = func(conn *WSConn) Agent {
			a := &agent{conn: conn, userData: s}

			return a
		}
	}
	wsServer.Start()
	return s, nil
}

// config map
// address default :80
func (s *Module) readConfig(m map[string]interface{}) error {

	address, err := server.ParseValue(m, "address", ":12080")
	if err != nil {
		return err
	}
	s.Address = address.(string)

	certFile, err := server.ParseValue(m, "certFile", "")
	if err != nil {
		return err
	}
	s.CertFile = certFile.(string)

	keyFile, err := server.ParseValue(m, "keyFile", "")
	if err != nil {
		return err
	}
	s.KeyFile = keyFile.(string)
	return nil
}

//Reload reload config
func (s *Module) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
}

// Stop goruntine
func (s *Module) Stop() error {

	return nil
}

// // Call async send message to module
// func (s *Module) Call(option int32, obj *server.CallObject) error {
// 	return nil
// }

// PushBytes async send string or bytes to queue
func (s *Module) PushBytes(option int32, buf []byte) error {
	return nil
}
