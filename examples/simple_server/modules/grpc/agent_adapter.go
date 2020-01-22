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
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/ss"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

//NewAgentAdapter create agent adapter
func NewAgentAdapter(conn network.Conn) network.IAgentAdapter {
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
	m := payload.(*ss.SSMessage)
	var body interface{}
	if err := a.Codec().Unmarshal(m.Payload, &body, nil); err != nil {
		server.Error("decode body failed ,err :%s", err.Error())
		return err
	}
	switch m.GetMsgType() {
	case server.MTRequest:
		return a.processRequestMessage(m, body)
	case server.MTResponse:
		return a.processResponseMessage(m, body)
	case server.MTNotify:
		return a.processNotifyMessage(m, body)
	case server.MTBroadcast:
		return a.processResponseMessage(m, body)
	default:
		return errors.New("unknown message type")
	}
}

func (a *AgentAdapter) processRequestMessage(m *ss.SSMessage, body interface{}) error {

	switch m.GetCmd() {
	case common.CSLoginCmd:
		return LoginHandler(a, m, body)
	default:
		return errors.New("unknown message")
	}

}
func (a *AgentAdapter) processResponseMessage(m *ss.SSMessage, body interface{}) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processNotifyMessage(m *ss.SSMessage, body interface{}) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processBroadcastMessage(m *ss.SSMessage, body interface{}) error {
	return errors.New("unknown message")
}
