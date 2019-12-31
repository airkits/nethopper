package wsjson

import (
	"github.com/gonethopper/nethopper/examples/model/json"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"time"
)

// NotifyLogin user to login
func NotifyLogin(s *Module, obj *server.CallObject, uid string, pwd string) (string, error) {

	if agent, ok := network.GetInstance().GetAuthAgent("user"); ok {
		m := json.NewWSMessage(uid, json.CSLoginCmd, 1, 1, server.MTRequest, agent.GetAdapter().Codec())
		body := &json.LoginReq{
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
func LoginResponse(agent network.IAgentAdapter, m *json.WSMessage) error {
	server.Info("LoginResponse get result %v", *m.Head)
	return nil
}
