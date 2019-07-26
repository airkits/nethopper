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

// TCPService struct to define service
type TCPService struct {
	server.BaseContext
	Address         string
	Network         string
	ReadBufferSize  int
	WriteBufferSize int
	ReadDeadline    time.Duration
	tcpListener     *net.TCPListener
}

// TCPServiceCreate  service create function
func TCPServiceCreate() (server.Service, error) {

	return &TCPService{}, nil
}

// UserData service custom option, can you store you data and you must keep goruntine safe
func (s *TCPService) UserData() int32 {
	return 0
}

// Setup init custom service and pass config map to service
// config
// m := map[string]interface{}{
//	"readBufferSize":32767,
//  "writeBufferSize":32767,
// 	"address":":8080",
//  "network":"tcp4",
//  "readDeadline":15,
//  "queueSize":1000,
// }
func (s *TCPService) Setup(m map[string]interface{}) (server.Service, error) {

	if err := s.readConfig(m); err != nil {
		panic(err)
	}
	return s, nil
}

// config map
// readBufferSize default 32767
// writeBufferSize default 32767
// address default :8080
// network default "tcp4"  use "tcp4/tcp6"
// readDeadline default 15
func (s *TCPService) readConfig(m map[string]interface{}) error {
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

	address, err := server.ParseValue(m, "address", ":8080")
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
func (s *TCPService) Reload(m map[string]interface{}) error {
	return nil
}

// Run create goruntine and run, always use ServiceRun to call this function
// Listen and bind local ip
func (s *TCPService) Run() {

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

	// loop accepting
	for {
		conn, err := s.accept()
		if err != nil {
			continue
		}
		go s.handler(conn, s.ReadDeadline)
	}
}

// accept the next incoming call and returns the new connection.
func (s *TCPService) accept() (net.Conn, error) {
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

func (s *TCPService) handler(conn net.Conn, readDeadline time.Duration) {
	defer conn.Close()
	// for reading the 2-Byte header
	header := make([]byte, 2)
	
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		server.Error("cannot get remote address:", err)
		return
	}
	// create a new session object for the connection
	// and record it's IP address
	sess := server.NewSession(s.ID(), host, port)
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

		// alloc a byte slice of the size defined in the header for reading data
		payload := make([]byte, size)
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			server.Warning("read payload failed, ip:%v reason:%v size:%v", sess.IP, err, n)
			return
		}

		// deliver the data to the input queue of agent()
		select {
		case in <- payload: // payload queued
		case <-sess.Die:
			server.Warning("connection closed by logic ip:%v", sess.IP)
			return
		}
	}
}

// Stop goruntine
func (s *TCPService) Stop() error {
	return nil
}

// SendMessage async send message to service
func (s *TCPService) SendMessage(option int32, msg *server.Message) error {
	return nil
}

// SendBytes async send string or bytes to queue
func (s *TCPService) SendBytes(option int32, buf []byte) error {
	return nil
}
