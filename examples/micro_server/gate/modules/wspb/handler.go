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

package wspb

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/c2s"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/pb/cs"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
)

//LoginHandler request login
func LoginHandler(agent network.IAgentAdapter, m transport.IMessage) error {

	message := m.(*cs.Message)
	req := c2s.LoginReq{}
	if err := ptypes.UnmarshalAny(message.Body, &req); err != nil {
		fmt.Println(err)
		return nil
	}
	server.Info("receive message %v", message)
	userID := server.StringToInt64(req.Uid)
	result, err := server.Call(server.ModuleIDLogic, common.CallIDLoginCmd, int32(userID), req.Uid, req.Passwd)

	resp := &c2s.LoginResp{
		Result: &c2s.Result{
			Code: 0,
			Msg:  "ok",
		},
		Uid:    req.Uid,
		Passwd: result.(string),
		Name:   result.(string),
	}
	if err != nil {
		resp.Result.Code = 500
		resp.Result.Msg = err.Error()
	}
	body, err := proto.Marshal(resp)
	if err != nil {
		return nil
	}

	respMsg := &ss.Message{
		ID:      message.GetID(),
		Cmd:     message.GetCmd(),
		MsgType: server.MTResponse,
		Body:    &any.Any{TypeUrl: "./c2s.LoginResp", Value: body},
	}

	payload, err := proto.Marshal(respMsg)
	if err != nil {
		return nil
	}
	agent.WriteMessage(payload)
	server.Info("send message %v", payload)
	return nil
}
