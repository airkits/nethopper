package gclient

import (
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/s2s"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
)

// RequestGetUserInfo user to login
func RequestGetUserInfo(s *Module, obj *server.CallObject, uid string, pwd string) (string, error) {

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

		if result, err := (agent.GetAdapter().(*AgentAdapter)).RPCCall(m); err == nil {

			server.Info("LoginResponse get result %v", result)
			resp := &s2s.LoginResp{}
			if err := ptypes.UnmarshalAny(result.Body, resp); err != nil {
				fmt.Println(err)
				return "", err
			}
			server.Info("LoginResponse get body %v", resp)
			if resp.Result.Code != 0 {
				return "", errors.New(resp.Result.Msg)
			}
			return resp.GetPasswd(), nil
		}

	}
	return "", errors.New("cant get agent")
}
