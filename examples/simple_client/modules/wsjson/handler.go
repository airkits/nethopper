package wsjson

import (
	"strconv"
	"time"

	"github.com/gonethopper/nethopper/examples/model/common"
	csjson "github.com/gonethopper/nethopper/examples/model/json"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/json"
	"github.com/gonethopper/nethopper/server"
)

// Login user to login
func Login(s *Module, obj *server.CallObject, uid string, pwd string) (string, server.Result) {
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return "", server.Result{Code: -1, Err: err}
	}
	if agent := s.GetAgent(uint32(uidInt)); agent != nil {

		req := &csjson.LoginReq{
			UID:    uid,
			Passwd: pwd,
		}

		var payload []byte
		var err error
		if payload, err = agent.GetAdapter().Codec().Marshal(req); err != nil {
			return "", server.Result{Code: -1, Err: err}
		}
		m := &json.Message{
			ID:      1,
			UID:     uint64(uidInt),
			Cmd:     common.CSLoginCmd,
			MsgType: server.MTRequest,
			Body:    string(payload),
		}
		if payload, err = agent.GetAdapter().Codec().Marshal(m); err != nil {
			return "", server.Result{Code: -1, Err: err}
		}

		if err := agent.SendMessage(payload); err != nil {
			server.Error("Notify login send failed %s ", err.Error())
			time.Sleep(1 * time.Second)
		} else {
			server.Info("Notify login send success")
		}

	}
	return "ok", server.Result{Code: 0, Err: nil}
}

//LoginResponse request login
func LoginResponse(agent network.IAgentAdapter, m transport.IMessage) error {

	server.Info("LoginResponse get result %v,body %v", m.(*json.Message), m.(*json.Message).GetBody())
	return nil
}
