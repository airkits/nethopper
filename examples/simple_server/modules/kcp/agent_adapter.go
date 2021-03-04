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

package kcp

import (
	"errors"

	"github.com/airkits/nethopper/codec"
	"github.com/airkits/nethopper/examples/model/common"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/network/transport/raw"

	"github.com/airkits/nethopper/server"
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
	message := payload.(*raw.Message)
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

//WriteMessage to connection
func (a *AgentAdapter) WriteMessage(msg interface{}) (err error) {
	msgBytes := msg.(*raw.Message).Pack()
	if err := a.Conn().WriteMessage(msgBytes); err != nil {
		server.Error("write message %x error: %v", msgBytes, err)
		return err
	}
	server.Info("send message success, length:%d", len(msgBytes))
	return nil
}

//ReadMessage goroutine not safe
func (a *AgentAdapter) ReadMessage() (interface{}, error) {
	var err error
	var b interface{}
	if b, err = a.Conn().ReadMessage(); err == nil {
		if b == nil {
			return b, err
		}
		msg := &raw.Message{}
		if err := msg.Unpack(b.([]byte)); err != nil {
			return nil, err
		}
		return msg, nil
	}
	return nil, err
}

func (a *AgentAdapter) processRequestMessage(message *raw.Message) error {

	switch message.Cmd {
	case common.SSLoginCmd:
		return LoginHandler(a, message)
	default:
		return errors.New("unknown message")
	}

}
func (a *AgentAdapter) processResponseMessage(message *raw.Message) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processNotifyMessage(message *raw.Message) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processBroadcastMessage(message *raw.Message) error {
	return errors.New("unknown message")
}

//OnClose agent close and clear
func (a *AgentAdapter) OnClose() {

}
