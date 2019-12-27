package wsclient

import (
	"github.com/gonethopper/nethopper/examples/model"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"time"
)

// NotifyLogin user to login
func NotifyLogin(s *Module, obj *server.CallObject, uid int64, pwd string) (string, error) {

	if agent, ok := network.GetInstance().GetAuthAgent("user"); ok {
		m := model.NewWSMessage(uid, model.CSLoginCmd, 1, server.MTRequest, agent.GetAdapter().Codec())
		body := &model.LoginReq{
			UID:    uid,
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
	server.Info("LoginResponse get result %v", *m.Head)
	return nil
}
