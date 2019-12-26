package websocket

import (
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/network"
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
func (a *AgentAdapter) ProcessMessage(payload []byte) {

}

//ProcessNotify process notify to client
func (a *AgentAdapter) ProcessNotify(obj interface{}) {

}
