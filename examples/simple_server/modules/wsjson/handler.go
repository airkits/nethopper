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
// * @Date: 2019-12-25 23:19:18
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-25 23:19:18

package wsjson

import (
	csjson "github.com/gonethopper/nethopper/examples/model/json"
	"github.com/gonethopper/nethopper/examples/simple_server/cmd"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/json"
	"github.com/gonethopper/nethopper/server"
)

//LoginHandler request login
func LoginHandler(agent network.IAgentAdapter, m transport.IMessage) error {
	message := m.(*json.Message)
	req := message.Body.(*csjson.LoginReq)
	server.Info("receive message %v", m)
	userID := server.StringToInt64(req.UID)
	v, result := server.Call(server.MIDLogic, cmd.LogicLogin, int32(userID), req.UID, req.Passwd)
	resp := &csjson.LoginResp{
		Data: v.(string),
	}
	if result.Err != nil {
		resp.Error(500, result.Err.Error())
	} else {
		resp.OK()
	}

	var body []byte
	var err error
	if body, err = agent.Codec().Marshal(resp); err != nil {
		return err
	}
	respMsg := &json.Message{
		Cmd:     message.GetCmd(),
		UID:     uint64(userID),
		MsgType: server.MTResponse,
		ID:      message.GetID(),
		Body:    string(body),
	}

	var payload []byte
	if payload, err = agent.Codec().Marshal(respMsg); err != nil {
		return err
	}

	agent.WriteMessage(payload)
	server.Info("send message %v", payload)
	return nil
}
