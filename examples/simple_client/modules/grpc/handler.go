package grpc

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/s2s"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
)

// NotifyLogin user to login
func NotifyLogin(s *Module, obj *server.CallObject, uid string, pwd string) (string, error) {

	if agent, ok := network.GetInstance().GetAuthAgent("user"); ok {

		req := &s2s.LoginReq{
			Uid:    uid,
			Passwd: pwd,
		}

		body, err := proto.Marshal(req)
		if err != nil {
			server.Error("Notify login send failed")
			return "error", nil
		}

		m := &ss.Message{
			ID:      agent.GetAdapter().GetSequence(),
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
	return "ok", nil
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
