package wspb

import (
	"errors"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/examples/model"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/cs"
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
func (a *AgentAdapter) ProcessMessage(payload interface{}) error {
	m := model.NewEmptyWSMessage(a.Codec())
	if err := m.DecodeHead(payload.([]byte)); err != nil {
		server.Error("recevie message failed ,err :%s", err.Error())
		return err
	}
	if err := m.DecodeBody(); err != nil {
		server.Error("decode body failed ,err :%s", err.Error())
		return err
	}
	head := m.Head.(*cs.WSHeader)
	switch head.MsgType {
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

func (a *AgentAdapter) processRequestMessage(m *model.WSMessage) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processResponseMessage(m *model.WSMessage) error {
	head := m.Head.(*cs.WSHeader)
	switch head.Cmd {
	case common.CSLoginCmd:
		return LoginResponse(a, m)
	default:
		return errors.New("unknown message")
	}

}
func (a *AgentAdapter) processNotifyMessage(m *model.WSMessage) error {
	return errors.New("unknown message")
}
func (a *AgentAdapter) processBroadcastMessage(m *model.WSMessage) error {
	return errors.New("unknown message")
}
