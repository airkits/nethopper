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

package tcp

import (
	"net"
	"time"

	"github.com/gonethopper/nethopper/server"
)

// ClientSocketModule struct to define module
type ClientSocketModule struct {
	server.BaseContext
	Address         string
	Network         string
	ReadBufferSize  int
	WriteBufferSize int
	ReadDeadline    time.Duration
	conn            net.Conn
}

// ClientSocketModuleCreate  module create function
func ClientSocketModuleCreate() (server.Module, error) {

	return &ClientSocketModule{}, nil
}

// UserData module custom option, can you store you data and you must keep goruntine safe
func (s *ClientSocketModule) UserData() int32 {
	return 0
}

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//	"readBufferSize":32767,
//  "writeBufferSize":32767,
// 	"address":":8888",
//  "network":"tcp4",
//  "readDeadline":15,
//  "queueSize":1000,
// }
func (s *ClientSocketModule) Setup(m map[string]interface{}) (server.Module, error) {

	if err := s.readConfig(m); err != nil {
		panic(err)
	}
	// Connect to server

	conn, err := net.Dial("tcp", s.Address)

	if err != nil {
		panic(err)
	}
	s.conn = conn
	server.Info("connect to : %s %s", s.Network, conn.RemoteAddr().String())

	return s, nil
}

// config map
// readBufferSize default 32767
// writeBufferSize default 32767
// address default :8888
// network default "tcp4"  use "tcp4/tcp6"
// readDeadline default 15
func (s *ClientSocketModule) readConfig(m map[string]interface{}) error {
	readBufferSize, err := server.ParseValue(m, "readBufferSize", 32767)
	if err != nil {
		return err
	}
	s.ReadBufferSize = readBufferSize.(int)

	writeBufferSize, err := server.ParseValue(m, "writeBufferSize", 32767)
	if err != nil {
		return err
	}
	s.WriteBufferSize = writeBufferSize.(int)

	address, err := server.ParseValue(m, "address", ":8888")
	if err != nil {
		return err
	}
	s.Address = address.(string)
	network, err := server.ParseValue(m, "network", "tcp4")
	if err != nil {
		return err
	}
	s.Network = network.(string)

	readDeadline, err := server.ParseValue(m, "readDeadline", 15)
	if err != nil {
		return err
	}
	s.ReadDeadline = time.Duration(readDeadline.(int)) * time.Second

	return nil
}

//Reload reload config
func (s *ClientSocketModule) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
// loop accepting
func (s *ClientSocketModule) OnRun(dt time.Duration) {

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			server.Info("ticker run request timeout")
		}
	}

}

// Stop goruntine
func (s *ClientSocketModule) Stop() error {
	return nil
}

// Call async send message to module
func (s *ClientSocketModule) Call(option int32, obj *server.CallObject) error {
	return nil
}

// PushBytes async send string or bytes to queue
func (s *ClientSocketModule) PushBytes(option int32, buf []byte) error {
	return nil
}
