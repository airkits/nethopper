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
// * @Date: 2019-12-26 09:21:53
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-26 09:21:53

package network

import (
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/server"
)

//AgentAdapter agent adapter
type AgentAdapter struct {
	codec codec.Codec
	conn  IConn
}

//Setup AgentAdapter
func (a *AgentAdapter) Setup(conn IConn, codec codec.Codec) {
	a.conn = conn
	a.codec = codec
}

//WriteMessage to connection
func (a *AgentAdapter) WriteMessage(payload interface{}) error {
	if err := a.conn.WriteMessage(payload); err != nil {
		server.Error("write message %x error: %v", payload, err)
		return err
	}
	return nil
}

//ReadMessage goroutine not safe
func (a *AgentAdapter) ReadMessage() (interface{}, error) {
	b, err := a.conn.ReadMessage()
	return b, err
}

// Codec get codec
func (a *AgentAdapter) Codec() codec.Codec {
	return a.codec
}

//SetCodec set codec
func (a *AgentAdapter) SetCodec(c codec.Codec) {
	a.codec = c
}

//Conn get conn
func (a *AgentAdapter) Conn() IConn {
	return a.conn
}

// SetConn set conn
func (a *AgentAdapter) SetConn(conn IConn) {
	a.conn = conn
}
