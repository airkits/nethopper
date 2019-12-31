package wspb

import (
	"errors"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/examples/model/json"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

//NewAgentAdapter create agent adapter
func NewAgentAdapter(conn network.Conn) network.IAgentAdapter {
	a := new(AgentAdapter)
	a.Setup(conn, codec.PBCodec)
	return a
}

//AgentAdapter do agent hander
type AgentAdapter struct {
	network.AgentAdapter
}

//ProcessMessage process request and notify message
func (a *AgentAdapter) ProcessMessage(payload []byte) error {
	m := json.NewEmptyWSMessage(a.Codec())
	if err := m.DecodeHead(payload); err != nil {
		server.Error("decode head failed ,err :%s", err.Error())
		return err
	}
	if err := m.DecodeBody(); err != nil {
		server.Error("decode body failed ,err :%s", err.Error())
		return err
	}

	switch m.Head.MsgType {
	case server.MTRequest:
		return a.processRequestMessage(m)
	case server.MTResponse:
		return a.processResponseMessage(m)
	case server.MTNotify:
		return a.processNotifyMessage(m)
	case server.MTBroadcast:
		return a.processResponseMessage(m)
	default:
		return errors.New("unknown message type")
	}
}

func (a *AgentAdapter) processRequestMessage(m *json.WSMessage) error {

	switch m.Head.CMD {
	case json.CSLoginCmd:
		return LoginHandler(a, m)
	default:
		return errors.New("unknown message")
	}

}
func (a *AgentAdapter) processResponseMessage(m *json.WSMessage) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processNotifyMessage(m *json.WSMessage) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processBroadcastMessage(m *json.WSMessage) error {
	return errors.New("unknown message")
}
