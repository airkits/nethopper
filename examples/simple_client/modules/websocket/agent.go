package websocket

import (
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

//NewAgent create new agent
func NewAgent(conn network.Conn, userData interface{}, codec codec.Codec) network.IAgent {
	a := new(ClientAgent)
	a.Init(conn, userData, codec)
	return a
}

type ClientAgent struct {
	network.Agent
}

func (a *ClientAgent) Run() {
	req := map[string]interface{}{
		"cmd":    "login",
		"uid":    1,
		"passwd": "game",
		"seq":    1,
	}
	for {

		a.WriteMessage(req)
		server.Info("send message %v", req)

		data, err := a.ReadMessage()
		if err != nil {
			server.Debug("read message: %v", err)
			break
		}
		out := make(map[string]interface{})
		if err := a.Codec().Unmarshal(data, &out, nil); err == nil {
			req["seq"] = out["seq"].(float64) + 1
		}
		server.Info(string(data))

	}
}
