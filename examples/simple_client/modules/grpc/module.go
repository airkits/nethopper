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
	"context"
	"io"
	"time"

	"github.com/gonethopper/nethopper/examples/model/pb/ss"
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
	Address string
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
	if err := s.ReadConfig(m); err != nil {
		panic(err)
	}

	conn, err := grpc.Dial(s.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		server.Fatal("did not connect: %v", err)
	}
	defer conn.Close()
	client := ss.NewRPCClient(conn)

	Transport(client)
	return s, nil
}

//Transport send and receive message by grpc connection
func Transport(c ss.RPCClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.Transport(ctx)
	if err != nil {
		server.Info("transport %v", err.Error())
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				server.Error("stream get error %s", err.Error())
				return
			}
			server.Info("receive message from server %v", in.GetCmd())
		}
	}()

	i := 1
	for {
		i++
		stream.Send(&ss.SSMessage{Cmd: "ss"})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}

	}
	stream.CloseSend()
	return nil
}

// ReadConfig config map
// address default :14000
func (s *Module) ReadConfig(m map[string]interface{}) error {
	if err := server.ParseConfigValue(m, "address", ":14000", &s.Address); err != nil {
		return err
	}
	return nil
}

//Reload reload config
func (s *Module) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
	server.RunSimpleFrame(s, 128)
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
