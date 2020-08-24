package grpc

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/s2s"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
)

// Login user to login
func Login(s *Module, obj *server.CallObject, uid string, pwd string) (string, server.Ret) {

	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return "", server.Ret{Code: -1, Err: err}
	}
	if agent := s.GetAgent(uint32(uidInt)); agent != nil {
		req := &s2s.LoginReq{
			Uid:    uid,
			Passwd: pwd,
		}

		body, err := proto.Marshal(req)
		if err != nil {
			server.Error("Notify login send failed")
			return "error", server.Ret{Code: -1, Err: err}
		}

		m := &ss.Message{
			ID:      agent.GetAdapter().GetSequence(),
			UID:     uint64(uidInt),
			Cmd:     common.SSLoginCmd,
			MsgType: server.MTRequest,
			Body:    &any.Any{TypeUrl: "./s2s.LoginReq", Value: body},
		}

		if err := agent.GetAdapter().WriteMessage(m); err != nil {
			server.Error("Notify login send failed %s ", err.Error())
			time.Sleep(1 * time.Second)
		} else {
			server.Info("Notify login send success")
		}

	}
	return "ok", server.Ret{Code: 0, Err: nil}
}

//LoginResponse request login
func LoginResponse(agent network.IAgentAdapter, m transport.IMessage) error {
	msg := m.(*ss.Message)
	server.Info("LoginResponse get result %v", msg)
	resp := &s2s.LoginResp{}
	if err := ptypes.UnmarshalAny(msg.Body, resp); err != nil {
		fmt.Println(err)
		return nil
	}
	server.Info("LoginResponse get body %v", resp)

	return nil
}
