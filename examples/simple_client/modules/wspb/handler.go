package wspb

import (
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/c2s"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport/pb/cs"
	"github.com/gonethopper/nethopper/server"
)

// Login user to login
func Login(s *Module, obj *server.CallObject, uid string, pwd string) (string, server.Ret) {

	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return "", server.Ret{Code: -1, Err: err}
	}
	if agent := s.GetAgent(uint32(uidInt)); agent != nil {
		req := &c2s.LoginReq{
			Uid:    uid,
			Passwd: pwd,
		}

		var body []byte
		var err error
		if body, err = agent.GetAdapter().Codec().Marshal(req); err != nil {
			return "", server.Ret{Code: -1, Err: err}
		}
		msg := &cs.Message{
			ID:      1,
			UID:     uint64(uidInt),
			Cmd:     common.CSLoginCmd,
			MsgType: server.MTRequest,
			Body:    &any.Any{TypeUrl: "./c2s.LoginReq", Value: body},
		}
		var payload []byte
		if payload, err = agent.GetAdapter().Codec().Marshal(msg); err != nil {
			return "", server.Ret{Code: -1, Err: err}
		}
		if err := agent.SendMessage(payload); err != nil {
			server.Error("Notify login send failed %s ", err.Error())
			time.Sleep(1 * time.Second)
		} else {
			server.Info("Notify login send success")
		}

	}
	return "ok", server.Ret{Code: 0, Err: nil}
}

//LoginResponse request login
func LoginResponse(agent network.IAgentAdapter, m *cs.Message) error {
	server.Info("LoginResponse get result %v", m)
	resp := &c2s.LoginResp{}
	if err := ptypes.UnmarshalAny(m.Body, resp); err != nil {
		return nil
	}
	server.Info("LoginResponse get body %v", resp)
	return nil
}
