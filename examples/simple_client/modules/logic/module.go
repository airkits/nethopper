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

package logic

import (
	"time"

	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/server"
)

// Module struct to define module
type Module struct {
	server.BaseContext
}

// ModuleCreate  module create function
func ModuleCreate() (server.Module, error) {
	return &Module{}, nil
}

// UserData module custom option, can you store you data and you must keep goruntine safe
// func (s *Module) UserData() int32 {
// 	return 0
// }

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *Module) Setup(m map[string]interface{}) (server.Module, error) {

	return s, nil
}

//Reload reload config
// func (s *Module) Reload(m map[string]interface{}) error {
// 	return nil
// }

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
	time.Sleep(1 * time.Second)
	//server.Call(server.ModuleIDWSClient, common.CSLoginCmd, 1, "1", "game")
	//server.Call(server.ModuleIDGRPCClient, common.SSLoginCmd, 1, "1", "game")

	server.Call(server.ModuleIDTCPClient, common.SSLoginCmd, 1, "1", "game")

	server.Info("logic start call")
}

// func (s *Module) request() {
// 	server.TraceCost("request cost")
// 	for i := 0; i < 1000; i++ {
// 	}
// }

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
