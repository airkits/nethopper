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
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/gonethopper/nethopper/server"
)

// SocketService struct to define service
type SocketService struct {
	server.BaseContext
	Address         string
	Network         string
	ReadBufferSize  int
	WriteBufferSize int
	ReadDeadline    time.Duration
	tcpListener     *net.TCPListener
}

// SocketServiceCreate  service create function
func SocketServiceCreate() (server.Service, error) {

	return &SocketService{}, nil
}

// UserData service custom option, can you store you data and you must keep goruntine safe
func (s *SocketService) UserData() int32 {
	return 0
}

// Setup init custom service and pass config map to service
// config
// m := map[string]interface{}{
//	"readBufferSize":32767,
//  "writeBufferSize":32767,
// 	"address":":8888",
//  "network":"tcp4",
//  "readDeadline":15,
//  "queueSize":1000,
// }
func (s *SocketService) Setup(m map[string]interface{}) (server.Service, error) {

	if err := s.readConfig(m); err != nil {
		panic(err)
	}
	// Listen and bind local ip
	tcpAddr, err := net.ResolveTCPAddr(s.Network, s.Address)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	s.tcpListener = listener
	server.Info("listening on: %s %s", s.Network, listener.Addr())

	return s, nil
}

// config map
// readBufferSize default 32767
// writeBufferSize default 32767
// address default :8888
// network default "tcp4"  use "tcp4/tcp6"
// readDeadline default 15
func (s *SocketService) readConfig(m map[string]interface{}) error {
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
func (s *SocketService) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ServiceRun to call this function
// loop accepting
func (s *SocketService) OnRun(dt time.Duration) {

	conn, err := s.accept()
	if err != nil {
		return
	}
	go s.handler(conn, s.ReadDeadline)

}

// accept the next incoming call and returns the new connection.
func (s *SocketService) accept() (net.Conn, error) {
	conn, err := s.tcpListener.AcceptTCP()
	if err != nil {
		server.Warning("accept failed: %s", err.Error())
		return nil, err
	}
	// set socket read buffer
	conn.SetReadBuffer(s.ReadBufferSize)
	// set socket write buffer
	conn.SetWriteBuffer(s.WriteBufferSize)
	return conn, nil
}

func (s *SocketService) handler(conn net.Conn, readDeadline time.Duration) {
	defer conn.Close()
	// for reading the 2-Byte header
	header := make([]byte, 2)
	cmd := make([]byte, 2)

	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		server.Error("cannot get remote address:", err)
		return
	}

	// create a new session object for the connection
	// and record it's IP address
	sess := server.CreateSession(s.ID(), host, port)
	server.Info("new connection from:%v port:%v", host, port)

	// read loop
	for {
		// solve dead link problem:
		// physical disconnection without any communcation between client and server
		// will cause the read to block FOREVER, so a timeout is a rescue.
		conn.SetReadDeadline(time.Now().Add(readDeadline))

		// read 2B header
		n, err := io.ReadFull(conn, header)
		if err != nil {
			server.Warning("read header failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}
		size := binary.BigEndian.Uint16(header)

		n, err = io.ReadFull(conn, cmd)
		if err != nil {
			server.Warning("read cmd failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}
		cmdLength := binary.BigEndian.Uint16(cmd)

		cmdBuffer := make([]byte, cmdLength)
		n, err = io.ReadFull(conn, cmdBuffer)
		if err != nil {
			server.Warning("read cmdBuffer failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}

		// alloc a byte slice of the size defined in the header for reading data
		payload := make([]byte, size-cmdLength-2)
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			server.Warning("read payload failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}
		// message := server.CreateMessage(s.ID(), server.ServiceIDLogic, server.MTRequest, string(cmdBuffer), payload)
		// server.Call(message.DestID, 0, message)

		// for i := 0; i < 8; i++ {
		// 	m, err := sess.MQ.AsyncPop()
		// 	if err == nil {
		// 		payload2 := m.(*server.Message).Payload

		// 		conn.Write(payload2)
		// 	} else {
		// 		break
		// 	}
		// }
		// deliver the data to the input queue of agent()
		select {
		case <-sess.Die:
			server.Warning("connection closed by logic ip:%v", sess.IP)
			return
		}
	}
}

// Stop goruntine
func (s *SocketService) Stop() error {
	return nil
}

// Call async send message to service
func (s *SocketService) Call(option int32, obj *server.CallObject) error {
	return nil
}

// PushBytes async send string or bytes to queue
func (s *SocketService) PushBytes(option int32, buf []byte) error {
	return nil
}
