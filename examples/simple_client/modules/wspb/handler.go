package wspb

import (
	"time"

	"github.com/gonethopper/nethopper/examples/model"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/cs"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

// NotifyLogin user to login
func NotifyLogin(s *Module, obj *server.CallObject, uid string, pwd string) (string, error) {

	if agent, ok := network.GetInstance().GetAuthAgent("user"); ok {
		m := model.NewWSMessage(uid, common.CSLoginCmd, 1, 1, server.MTRequest, agent.GetAdapter().Codec())
		body := &cs.LoginReq{
			Uid:    uid,
			Passwd: pwd,
		}
		m.Body = body
		var payload []byte
		var err error
		if payload, err = m.Encode(); err != nil {
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
func LoginResponse(agent network.IAgentAdapter, m *model.WSMessage) error {
	server.Info("LoginResponse get result %v", *(m.Head.(*cs.WSHeader)))
	server.Info("LoginResponse get body %v", *(m.Body.(*cs.LoginResp)))

	return nil
}
