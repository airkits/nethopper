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

package snowflake

import (
	"time"

	"github.com/gonethopper/nethopper/examples/usercenter/cmd"
	"github.com/gonethopper/nethopper/examples/usercenter/model"
	"github.com/gonethopper/nethopper/server"
)

// Module struct to define module
type Module struct {
	server.BaseContext
	Conf *model.SFConfig
}

// ModuleCreate  module create function
func ModuleCreate() (server.Module, error) {
	return &Module{}, nil
}

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *Module) Setup(conf server.IConfig) (server.Module, error) {
	s.Conf = conf.(*model.SFConfig)

	s.RegisterHandler(cmd.MCSFGetUID, GetUID)
	s.CreateWorkerPool(s, 128, 10*time.Second, true)
	return s, nil
}

//GetHost get snowflake server host
func (s *Module) GetHost() string {
	return s.Conf.Hosts[0]
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
	server.RunSimpleFrame(s, 128)
}

// Stop goruntine
func (s *Module) Stop() error {
	return nil
}