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

package gclient

import (
	"errors"

	"github.com/gonethopper/nethopper/base/skiplist"
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"

	"github.com/gonethopper/nethopper/server"
)

//NewAgentAdapter create agent adapter
func NewAgentAdapter(conn network.IConn) network.IAgentAdapter {
	a := new(AgentAdapter)
	a.Setup(conn, codec.PBCodec)
	a.Cache = skiplist.New()
	return a
}

//AgentAdapter do agent hander
type AgentAdapter struct {
	network.AgentAdapter
	Cache *skiplist.SkipList
}

// func (a *AgentAdapter) decodePBBody(m transport.IMessage) error {
// 	message := m.(*ss.Message)
// 	var body proto.Message
// 	var err error
// 	if body, err = pb.CreateBody(head.MsgType, head.Cmd); err != nil {
// 		return err
// 	}
// 	if err = m.Codec().Unmarshal(head.Payload, body, nil); err != nil {
// 		return err
// 	}

// 	m.Body = body
// 	return nil
// }

//ProcessMessage process request and notify message
func (a *AgentAdapter) ProcessMessage(payload interface{}) error {
	m := payload.(*ss.Message)
	switch m.MsgType {
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

func (a *AgentAdapter) processRequestMessage(m *ss.Message) error {

	switch m.Cmd {
	default:
		return errors.New("unknown message")
	}

}
func (a *AgentAdapter) processResponseMessage(m *ss.Message) error {
	v := a.Cache.Remove(float64(m.GetID()))
	if v != nil {
		obj := v.Value().(*server.CallObject)
		obj.ChanRet <- server.RetObject{
			Ret: m,
			Err: nil,
		}
		return nil
	}
	return errors.New("cant find request object")
}
func (a *AgentAdapter) processNotifyMessage(m *ss.Message) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processBroadcastMessage(m *ss.Message) error {
	return errors.New("unknown message")
}

//RPCCall remote call
func (a *AgentAdapter) RPCCall(msg *ss.Message) (*ss.Message, error) {
	var obj = server.NewCallObject(msg.GetCmd(), 0, msg)
	a.Cache.Set(float64(msg.GetID()), obj)
	a.WriteMessage(msg)
	result := <-obj.ChanRet
	return (result.Ret).(*ss.Message), result.Err
}

//WriteMessage to connection
func (a *AgentAdapter) WriteMessage(payload interface{}) error {
	if err := a.Conn().WriteMessage(payload); err != nil {
		server.Error("write message %x error: %v", payload, err)
		return err
	}
	return nil
}

//ReadMessage goroutine not safe
func (a *AgentAdapter) ReadMessage() (interface{}, error) {
	b, err := a.Conn().ReadMessage()
	return b, err
}

//OnClose agent close and clear
func (a *AgentAdapter) OnClose() {

}
