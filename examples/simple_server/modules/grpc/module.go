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

package grpc

import (
	"net"
	"time"

	"github.com/gonethopper/nethopper/server"
	"google.golang.org/grpc"
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
	gs             *grpc.Server
	Address        string
	MaxConnNum     int
	RWQueueSize    int
	MaxMessageSize int
	listener       net.Listener
}

//ReadConfig read config
// config
// m := map[string]interface{}{
//  "queueSize":1000,
//  "grpcAddress":":14000",
//	"maxConnNum":1024,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
// }
func (s *Module) ReadConfig(m map[string]interface{}) error {
	if err := server.ParseConfigValue(m, "grpcAddress", ":14000", &s.Address); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "maxConnNum", 1024, &s.MaxConnNum); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "socketQueueSize", 100, &s.RWQueueSize); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "maxMessageSize", 4096, &s.MaxMessageSize); err != nil {
		return err
	}
	return nil
}

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
//  "grpcAddress":":14000",
//	"maxConnNum":1024,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
// }
func (s *Module) Setup(m map[string]interface{}) (server.Module, error) {
	if err := s.ReadConfig(m); err != nil {
		panic(err)
	}

	s.gs = grpc.NewServer()
	lis, err := net.Listen("tcp", s.Address)

	if err != nil {
		server.Error("failed to listen: %v", err)
		return nil, err
	}
	s.listener = lis

	server.GO(s.web)

	return s, nil
}
func (s *Module) web() {
	_ = s.gs.Serve(s.listener)
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
	server.RunSimpleFrame(s, 128)
}
