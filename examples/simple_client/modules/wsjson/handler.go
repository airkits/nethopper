package wsjson

import (
	"time"

	csjson "github.com/gonethopper/nethopper/examples/model/json"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport"
	"github.com/gonethopper/nethopper/network/transport/json"
	"github.com/gonethopper/nethopper/server"
)

// NotifyLogin user to login
func NotifyLogin(s *Module, obj *server.CallObject, uid string, pwd string) (string, error) {

	if agent, ok := network.GetInstance().GetAuthAgent("user"); ok {
		m := transport.NewMessage(transport.HeaderTypeWSJSON, agent.GetAdapter().Codec())
		body := &csjson.LoginReq{
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
func LoginResponse(agent network.IAgentAdapter, m *transport.Message) error {
	server.Info("LoginResponse get result %v", *(m.Header.(*json.Header)))
	return nil
}
