package websocket

import (
	"net"
	"reflect"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/server"
)

type Agent interface {
	Run()
	OnClose()
}

type agent struct {
	conn     *WSConn
	userData interface{}
}

func (a *agent) Run() {
	req := map[string]interface{}{
		"cmd":    "login",
		"uid":    1,
		"passwd": "game",
		"seq":    1,
	}
	for {

		a.WriteMsg(req)
		server.Info("send message %v", req)

		data, err := a.conn.ReadMsg()
		if err != nil {
			server.Debug("read message: %v", err)
			break
		}
		out := make(map[string]interface{})
		if err := codec.JSONCodec.Unmarshal(data, &out, nil); err == nil {
			req["seq"] = out["seq"].(float64) + 1
		}
		server.Info(string(data))

		// if a.gate.Processor != nil {
		// 	msg, err := a.gate.Processor.Unmarshal(data)
		// 	if err != nil {
		// 		server.Debug("unmarshal message error: %v", err)
		// 		break
		// 	}
		// 	//msgType := reflect.TypeOf(msg)
		// 	a.gate.AgentChanRPC.Go(CommandAgentMsg, msg, a)
		// 	//err = a.gate.Processor.Route(msg, a)
		// 	//if err != nil {
		// 	//	log.Debug("route message error: %v", err)
		// 	//	break
		// 	//}
		// }
	}
}

func (a *agent) OnClose() {
	// if a.gate.AgentChanRPC != nil {
	// 	err := a.gate.AgentChanRPC.Call0(CommandAgentClose, a)
	// 	if err != nil {
	// 		server.Error("chanrpc error: %v", err)
	// 	}
	// }
}

func (a *agent) WriteMsg(msg interface{}) {
	data, err := codec.JSONCodec.Marshal(msg, nil)
	// if a.gate.Processor != nil {
	// 	data, err := a.gate.Processor.Marshal(msg)
	if err != nil {
		server.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
		return
	}
	err = a.conn.WriteMsg(data)
	if err != nil {
		server.Error("write message %v error: %v", reflect.TypeOf(msg), err)
	}
	// }
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
