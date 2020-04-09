package wspb

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/c2s"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport/pb/cs"
	"github.com/gonethopper/nethopper/server"
)

// NotifyLogin user to login
func NotifyLogin(s *Module, obj *server.CallObject, uid string, pwd string) (string, error) {

	if agent, ok := network.GetInstance().GetAuthAgent("user"); ok {
		req := &c2s.LoginReq{
			Uid:    uid,
			Passwd: pwd,
		}

		var body []byte
		var err error
		if body, err = agent.GetAdapter().Codec().Marshal(req, nil); err != nil {
			return "", err
		}
		msg := &cs.Message{
			ID:      1,
			Cmd:     common.CSLoginCmd,
			MsgType: server.MTRequest,
			Body:    &any.Any{TypeUrl: "./" + common.CSLoginCmd, Value: body},
		}
		var payload []byte
		if payload, err = agent.GetAdapter().Codec().Marshal(msg, nil); err != nil {
			return "", err
		}
		if err := agent.SendMessage(payload); err != nil {
			server.Error("Notify login send failed %s ", err.Error())
			time.Sleep(1 * time.Second)
		} else {
			server.Info("Notify login send success")
		}
	}
	return "ok", nil
}

//LoginResponse request login
func LoginResponse(agent network.IAgentAdapter, m *cs.Message) error {
	server.Info("LoginResponse get result %v", *m)
	resp := &c2s.LoginResp{}
	if err := ptypes.UnmarshalAny(m.Body, resp); err != nil {
		return nil
	}
	server.Info("LoginResponse get body %v", resp)
	return nil
}
