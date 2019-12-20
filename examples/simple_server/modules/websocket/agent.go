package websocket

import (
	"strconv"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/examples/simple_server/common"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

func NewAgent(conn network.Conn, userData interface{}, codec codec.Codec) network.IAgent {
	a := new(Agent)
	a.Init(conn, userData, codec)
	return a
}

type Agent struct {
	network.Agent
}

func (a *Agent) Run() {
	for {
		data, err := a.ReadMessage()
		if err != nil {
			server.Debug("read message: %v", err)
			break
		}
		out := make(map[string]interface{})

		if err := a.Codec().Unmarshal(data, &out, nil); err == nil {
			server.Info("receive message %v", out)
			uid := out["uid"].(float64)
			pwd := out["passwd"].(string)
			result, err2 := server.Call(server.ModuleIDLogic, common.CallIDLoginCmd, int32(uid), strconv.FormatFloat(uid, 'f', -1, 64), pwd)
			if err2 != nil {
				out["result"] = result.(string)
				out["code"] = 1
			} else {
				out["result"] = result.(string)
				out["code"] = 0
			}

			out["seq"] = out["seq"].(float64) + 1
		} else {
			server.Error(err)
		}
		a.WriteMessage(out)

	}
}
