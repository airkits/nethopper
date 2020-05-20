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

package orm

import (
	"time"

	"github.com/gonethopper/nethopper/server"
)

// OrmModule struct to define module
type OrmModule struct {
	server.BaseContext
}

// OrmModuleCreate  module create function
func OrmModuleCreate() (server.Module, error) {
	return &OrmModule{}, nil
}

// UserData module custom option, can you store you data and you must keep goruntine safe
func (s *OrmModule) UserData() int32 {
	return 0
}

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *OrmModule) Setup(conf server.IConfig) (server.Module, error) {
	return s, nil
}

//Reload reload config
func (s *OrmModule) Reload(conf server.IConfig) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *OrmModule) OnRun(dt time.Duration) {

}

// Stop goruntine
func (s *OrmModule) Stop() error {
	return nil
}

// Call async send message to module
func (s *OrmModule) Call(option int32, obj *server.CallObject) error {
	return nil
}

// PushBytes async send string or bytes to queue
func (s *OrmModule) PushBytes(option int32, buf []byte) error {
	return nil
}
