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
// * @Date: 2020-01-09 11:02:02
// * @Last Modified by:   ankye
// * @Last Modified time: 2020-01-09 11:02:02

package rpc_test

import (
	"errors"

	"github.com/airkits/nethopper/codec"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
)

//NewAgentAdapter create agent adapter
func NewAgentAdapter(conn network.IConn) network.IAgentAdapter {
	a := new(AgentAdapter)
	a.Setup(conn, codec.PBCodec)
	return a
}

//AgentAdapter do agent hander
type AgentAdapter struct {
	network.AgentAdapter
}

//DecodeMessage process request and notify message
func (a *AgentAdapter) DecodeMessage(payload interface{}) error {
	return errors.New("unknown message")
}

//WriteMessage to connection
func (a *AgentAdapter) WriteMessage(payload interface{}) error {
	if err := a.Conn().WriteMessage(payload); err != nil {
		log.Error("write message %x error: %v", payload, err)
		return err
	}
	return nil
}

//ReadMessage goroutine not safe
func (a *AgentAdapter) ReadMessage() (interface{}, error) {
	b, err := a.Conn().ReadMessage()
	return b, err
}

// func (a *AgentAdapter) processRequestMessage(m *ss.Message, body protoreflect.ProtoMessage) error {

// 	return errors.New("unknown message")

// }

func (a *AgentAdapter) OnClose() {

}
