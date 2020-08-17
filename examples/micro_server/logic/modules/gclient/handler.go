package gclient

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/s2s"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
)

// GetUser user to login
func GetUser(s *Module, obj *server.CallObject, uid string, pwd string) (string, server.Result) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return "", server.Result{Code: -1, Err: errors.New("convert uid failed")}
	}
	if agent := s.GetAgent(uint32(uidInt)); agent != nil {

		req := &s2s.LoginReq{
			Uid:    uid,
			Passwd: pwd,
		}

		body, err := proto.Marshal(req)
		if err != nil {
			server.Error("Notify login send failed")
			return "error", server.Result{Code: 0, Err: nil}
		}

		m := &ss.Message{
			ID:      agent.GetAdapter().GetSequence(),
			UID:     uint64(uidInt),
			Cmd:     common.SSLoginCmd,
			MsgType: server.MTRequest,
			Body:    &any.Any{TypeUrl: "./s2s.LoginReq", Value: body},
		}

		if v, result := (agent.GetAdapter().(*AgentAdapter)).RPCCall(m); result.Err == nil {

			server.Info("LoginResponse get result %v", result)
			resp := &s2s.LoginResp{}
			if err := ptypes.UnmarshalAny(v.Body, resp); err != nil {
				fmt.Println(err)
				return "", result
			}
			server.Info("LoginResponse get body %v", resp)
			if resp.Result.Code != 0 {
				return "", result
			}
			return resp.GetPasswd(), result
		}

	}
	return "", server.Result{Code: 0, Err: errors.New("cant get agent")}
}
