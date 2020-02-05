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

package wspb

import (
	"errors"

	"github.com/gogo/protobuf/proto"
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/pb/cs"

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

func (a *AgentAdapter) decodePBBody(m *transport.Message) error {
	head := m.Header.(*cs.Header)
	var body proto.Message
	var err error
	if body, err = pb.CreateBody(head.MsgType, head.Cmd); err != nil {
		return err
	}
	if err = m.Codec().Unmarshal(head.Payload, body, nil); err != nil {
		return err
	}

	m.Body = body
	return nil
}

//ProcessMessage process request and notify message
func (a *AgentAdapter) ProcessMessage(payload interface{}) error {
	m := transport.NewMessage(transport.HeaderTypeWSPB, a.Codec())
	if err := m.DecodeHeader(payload.([]byte)); err != nil {
		server.Error("decode head failed ,err :%s", err.Error())
		return err
	}
	if err := a.decodePBBody(m); err != nil {
		server.Error("decode body failed ,err :%s", err.Error())
		return err
	}
	head := m.Header.(*cs.Header)
	switch head.MsgType {
	case server.MTRequest:
		return a.processRequestMessage(m)
	case server.MTResponse:
		return a.processResponseMessage(m)
	case server.MTNotify:
		return a.processNotifyMessage(m)
	case server.MTBroadcast:
		return a.processResponseMessage(m)
	default:
		return errors.New("unknown message type")
	}
}

func (a *AgentAdapter) processRequestMessage(m *transport.Message) error {

	head := m.Header.(*cs.Header)
	switch head.Cmd {
	case common.CSLoginCmd:
		return LoginHandler(a, m)
	default:
		return errors.New("unknown message")
	}

}
func (a *AgentAdapter) processResponseMessage(m *transport.Message) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processNotifyMessage(m *transport.Message) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processBroadcastMessage(m *transport.Message) error {
	return errors.New("unknown message")
}
