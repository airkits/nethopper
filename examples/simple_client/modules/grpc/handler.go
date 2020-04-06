package grpc

import (
	"time"

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
		m := transport.NewMessage(transport.HeaderTypeGRPCPB, agent.GetAdapter().Codec())
		m.Header = m.NewHeader(1, common.SSLoginCmd, server.MTRequest)

		body := &s2s.LoginReq{
			Uid:    uid,
			Passwd: pwd,
		}
		m.Body = body
		if err := m.EncodeBody(); err != nil {
			server.Error("Notify login send failed %s ", err.Error())
			return "error", nil
		}
		if err := agent.GetAdapter().WriteMessage(m.Header); err != nil {
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
	server.Info("LoginResponse get result %v", *(m.Header.(*ss.Header)))
	server.Info("LoginResponse get body %v", *(m.Body.(*s2s.LoginResp)))

	return nil
}
