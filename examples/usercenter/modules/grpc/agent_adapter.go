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

package grpc

import (
	"errors"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"

	"github.com/gonethopper/nethopper/server"
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

//ProcessMessage process request and notify message
func (a *AgentAdapter) ProcessMessage(payload interface{}) error {
	message := payload.(*ss.Message)
	switch message.MsgType {
	case server.MTRequest:
		return a.processRequestMessage(message)
	case server.MTResponse:
		return a.processResponseMessage(message)
	case server.MTNotify:
		return a.processNotifyMessage(message)
	case server.MTBroadcast:
		return a.processResponseMessage(message)
	default:
		return errors.New("unknown message type")
	}
}

func (a *AgentAdapter) processRequestMessage(message *ss.Message) error {

	switch message.Cmd {
	//case common.SSLoginCmd:
	//	return LoginHandler(a, message)
	default:
		return errors.New("unknown message")
	}

}
func (a *AgentAdapter) processResponseMessage(message *ss.Message) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processNotifyMessage(message *ss.Message) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processBroadcastMessage(message *ss.Message) error {
	return errors.New("unknown message")
}

//OnClose agent close and clear
func (a *AgentAdapter) OnClose() {

}
