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

package gclient

import (
	"time"

	"github.com/airkits/nethopper/examples/micro_server/gate/protocol"
	"github.com/airkits/nethopper/libs/skiplist"
	"github.com/airkits/nethopper/mediator"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/network/grpc"
	"github.com/airkits/nethopper/server"
)

// ModuleCreate  module create function
func ModuleCreate() (mediator.IModule, error) {
	return &Module{}, nil
}

// Module struct to define module
type Module struct {
	server.BaseContext
	grpcClient *grpc.Client
	Clients    *skiplist.SkipList
}

// UserData module custom option, can you store you data and you must keep goruntine safe
func (s *Module) UserData() int32 {
	return 0
}

//Handlers set moudle handlers
func (s *Module) Handlers() map[string]interface{} {
	return map[string]interface{}{}
}

//ReflectHandlers set moudle reflect handlers
func (s *Module) ReflectHandlers() map[string]interface{} {
	return map[string]interface{}{
		protocol.GClientGetUser: GetUser,
	}
}

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
//  "grpcAddress":14000,
// }
func (s *Module) Setup(conf server.IConfig) (mediator.IModule, error) {
	if err := s.ReadConfig(conf); err != nil {
		panic(err)
	}
	s.CreateWorkerPool(s, 128, 10*time.Second, true)

	s.Clients = skiplist.New()
	s.grpcClient = grpc.NewClient(conf, func(conn network.IConn, uid uint64, token string) network.IAgent {
		a := network.NewAgent(NewAgentAdapter(conn), uid, token)
		s.Clients.Set(float64(uid), a)
		return a
	}, func(agent network.IAgent) {
		s.Clients.Remove(float64(agent.UID()))
	})
	s.grpcClient.Run()

	return s, nil
}

//GetAgent get agent by option
func (s *Module) GetAgent(option uint32) network.IAgent {
	v := s.Clients.Get(float64(0))
	if v != nil {
		return v.Value().(network.IAgent)
	}
	return nil
}

// ReadConfig config map
// address default :80
func (s *Module) ReadConfig(conf server.IConfig) error {
	return nil
}

//Reload reload config
func (s *Module) Reload(conf server.IConfig) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
	server.RunSimpleFrame(s)
}

// Stop goruntine
func (s *Module) Stop() error {

	return nil
}

// PushBytes async send string or bytes to queue
func (s *Module) PushBytes(option int32, buf []byte) error {
	return nil
}
