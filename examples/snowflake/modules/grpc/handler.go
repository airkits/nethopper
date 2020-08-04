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

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/model/pb/s2s"
	"github.com/gonethopper/nethopper/examples/snowflake/cmd"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/nethopper/utils"
)

//GenUIDHandler request uid
func GenUIDHandler(agent network.IAgentAdapter, m transport.IMessage) error {
	message := m.(*ss.Message)
	req := s2s.GenUIDReq{}
	if err := ptypes.UnmarshalAny(message.Body, &req); err != nil {
		fmt.Println(err)
		return nil
	}
	server.Info("receive message %v", req)
	result, err := server.Call(server.ModuleIDLogic, cmd.CallIDGenUIDCmd, utils.RandomInt32(0, 1000), req.Channel)
	// header := m.(*ss.Header)
	// outM := transport.NewMessage(transport.HeaderTypeGRPCPB, agent.Codec())
	// outM.Header = outM.NewHeader(header.GetID(), header.GetCmd(), server.MTResponse)

	resp := &s2s.GenUIDResp{
		Result: &s2s.Result{
			Code: 0,
			Msg:  "ok",
		},
		Uid: result.(uint64),
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
		UID:     uint64(req.Channel),
		Cmd:     message.GetCmd(),
		MsgType: server.MTResponse,
		Body:    &any.Any{TypeUrl: "./s2s.GenUIDResp", Value: body},
	}

	agent.WriteMessage(respMsg)
	server.Info("send message %v", respMsg)
	return nil
}

//GenUIDsHandler request uid
func GenUIDsHandler(agent network.IAgentAdapter, m transport.IMessage) error {
	message := m.(*ss.Message)
	req := s2s.GenUIDsReq{}
	if err := ptypes.UnmarshalAny(message.Body, &req); err != nil {
		fmt.Println(err)
		return nil
	}
	server.Info("receive message %v", req)
	result, err2 := server.Call(server.ModuleIDLogic, cmd.CallIDGenUIDsCmd, utils.RandomInt32(0, 1024), req.Channel, req.Num)
	// header := m.(*ss.Header)
	// outM := transport.NewMessage(transport.HeaderTypeGRPCPB, agent.Codec())
	// outM.Header = outM.NewHeader(header.GetID(), header.GetCmd(), server.MTResponse)

	resp := &s2s.GenUIDsResp{
		Result: &s2s.Result{
			Code: 0,
			Msg:  "ok",
		},
		Uid: result.([]uint64),
	}
	if err2 != nil {
		resp.Result.Code = 500
		resp.Result.Msg = err2.Error()
	}
	body, err := proto.Marshal(resp)
	if err != nil {
		return nil
	}

	respMsg := &ss.Message{
		ID:      message.GetID(),
		UID:     uint64(req.Channel),
		Cmd:     message.GetCmd(),
		MsgType: server.MTResponse,
		Body:    &any.Any{TypeUrl: "./s2s.GenUIDsResp", Value: body},
	}

	agent.WriteMessage(respMsg)
	server.Info("send message %v", respMsg)
	return nil
}
