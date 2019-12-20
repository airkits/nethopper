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
// * @Date: 2019-12-20 19:39:11
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-20 19:39:11

package network

import (
	"net"
	"reflect"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/server"
)

//IAgent agent interface define
type IAgent interface {
	Run()
	OnClose()
	WriteMessage(msg interface{})
	ReadMessage() ([]byte, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	Conn() Conn
	SetConn(conn Conn)
	Codec() codec.Codec
	SetCodec(c codec.Codec)
	Init(Conn, interface{}, codec.Codec)
	Token() string
	SetToken(string)
	IsAuth() bool
}

//NewAgent create new agent
func NewAgent(conn Conn, userData interface{}, codec codec.Codec) IAgent {
	return &Agent{conn: conn, userData: userData, codec: codec}
}

//Agent base agent struct
type Agent struct {
	conn     Conn
	userData interface{}
	codec    codec.Codec
	token    string
}

//Init agent
func (a *Agent) Init(conn Conn, userData interface{}, codec codec.Codec) {
	a.conn = conn
	a.userData = userData
	a.codec = codec
}

//Token get token
func (a *Agent) Token() string {
	return a.token
}

//IsAuth if set token return true else return false
func (a *Agent) IsAuth() bool {
	return len(a.token) > 0
}

//SetToken set token
func (a *Agent) SetToken(token string) {
	a.token = token
}

//Conn get conn
func (a *Agent) Conn() Conn {
	return a.conn
}

// SetConn set conn
func (a *Agent) SetConn(conn Conn) {
	a.conn = conn
}

// Codec get codec
func (a *Agent) Codec() codec.Codec {
	return a.codec
}

//SetCodec set codec
func (a *Agent) SetCodec(c codec.Codec) {
	a.codec = c
}

//Run agent start run
func (a *Agent) Run() {
	// for {
	// 	data, err := a.ReadMessage()
	// 	if err != nil {
	// 		server.Debug("read message: %v", err)
	// 		break
	// 	}
	// 	out := make(map[string]interface{})
	// 	if err := a.Codec().Unmarshal(data, &out, nil); err == nil {
	// 		server.Info("receive message %v", out)
	// 		out["seq"] = out["seq"].(float64) + 1
	// 	} else {
	// 		server.Error(err)
	// 	}
	// 	a.WriteMessage(out)
	// }
}

func (a *Agent) OnClose() {

}

func (a *Agent) WriteMessage(msg interface{}) {
	data, err := a.Codec().Marshal(msg, nil)
	if err != nil {
		server.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
		return
	}
	err = a.conn.WriteMessage(data)
	if err != nil {
		server.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}

}

//ReadMessage goroutine not safe
func (a *Agent) ReadMessage() ([]byte, error) {
	b, err := a.conn.ReadMessage()
	return b, err
}

func (a *Agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *Agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *Agent) Close() {
	a.conn.Close()
}

func (a *Agent) Destroy() {
	a.conn.Destroy()
}

func (a *Agent) UserData() interface{} {
	return a.userData
}

func (a *Agent) SetUserData(data interface{}) {
	a.userData = data
}
