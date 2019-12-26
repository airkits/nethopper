package websocket

import (
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

//NewAgentAdapter create agent adapter
func NewAgentAdapter(conn network.Conn) network.IAgentAdapter {
	a := new(AgentAdapter)
	a.Setup(conn, codec.JSONCodec)
	return a
}

//AgentAdapter do agent hander
type AgentAdapter struct {
	network.AgentAdapter
}

//ProcessMessage process request and notify message
//ProcessMessage process request and notify message
func (a *AgentAdapter) ProcessMessage(payload []byte) {
	m := new(LoginCmd)
	if err := a.Codec().Unmarshal(payload, m, nil); err != nil {
		return
	}
	m.Seq++

	server.Info("recevie message %v", m)

	if data, err := a.Codec().Marshal(m, nil); err == nil {
		if err := a.WriteMessage(data); err != nil {
			server.Info("send message to client failed , error %s", err.Error())
			return
		}
		server.Info("send message to client %v ", m)
	}
}

//ProcessNotify process notify to client
func (a *AgentAdapter) ProcessNotify(obj interface{}) {

}
