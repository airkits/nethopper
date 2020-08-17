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

package grpc

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/micro_server/gamedb/cmd"
	"github.com/gonethopper/nethopper/examples/model/pb/s2s"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
)

//Login request login
func Login(agent network.IAgentAdapter, m transport.IMessage) error {
	message := m.(*ss.Message)
	req := s2s.LoginReq{}
	if err := ptypes.UnmarshalAny(message.Body, &req); err != nil {
		fmt.Println(err)
		return nil
	}
	server.Info("receive message %v", req)
	userID := server.StringToInt64(req.Uid)
	v, result := server.Call(server.MIDLogic, cmd.GRPCLogin, int32(userID), req.Uid, req.Passwd)
	// header := m.(*ss.Header)
	// outM := transport.NewMessage(transport.HeaderTypeGRPCPB, agent.Codec())
	// outM.Header = outM.NewHeader(header.GetID(), header.GetCmd(), server.MTResponse)

	resp := &s2s.LoginResp{
		Result: &s2s.Result{
			Code: 0,
			Msg:  "ok",
		},
		Uid:    req.Uid,
		Passwd: v.(string),
		Name:   v.(string),
	}
	if result.Err != nil {
		resp.Result.Code = 500
		resp.Result.Msg = result.Err.Error()
	}
	body, err := proto.Marshal(resp)
	if err != nil {
		return nil
	}

	respMsg := &ss.Message{
		ID:      message.GetID(),
		UID:     uint64(userID),
		Cmd:     message.GetCmd(),
		MsgType: server.MTResponse,
		Body:    &any.Any{TypeUrl: "./s2s.LoginResp", Value: body},
	}

	agent.WriteMessage(respMsg)
	server.Info("send message %v", respMsg)
	return nil
}
