package tcp

import (
	"fmt"
	"time"

	"github.com/airkits/nethopper/codec"
	"github.com/airkits/nethopper/examples/model/common"
	"github.com/airkits/nethopper/examples/model/pb/s2s"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/network/transport"
	"github.com/airkits/nethopper/network/transport/raw"
	"github.com/airkits/nethopper/server"
	"github.com/airkits/nethopper/utils/conv"
)

// Login user to login
func Login(s *Module, uid string, pwd string) (string, server.Ret) {

	if id := conv.Str2Uint64(uid); id > 0 {
		if agent, ok := network.GetInstance().GetAuthAgent(id); ok {

			req := &s2s.LoginReq{
				Uid:    uid,
				Passwd: pwd,
			}

			body, err := codec.PBCodec.Marshal(req)
			if err != nil {
				server.Error("Notify login send failed")
				return "error", server.Ret{Code: -1, Err: err}
			}

			m := &raw.Message{
				ID:      1,
				UID:     id,
				Cmd:     common.SSLoginCmd,
				MsgType: server.MTRequest,
				Seq:     0,
				Body:    body,
			}
			if err := agent.GetAdapter().WriteMessage(m); err != nil {
				server.Error("Notify login send failed %s ", err.Error())
				time.Sleep(1 * time.Second)
			} else {
				server.Info("Notify login send success")
			}
		}
	}
	return "ok", server.Ret{Code: 0, Err: nil}
}

//LoginResponse request login
func LoginResponse(agent network.IAgentAdapter, m transport.IMessage) error {
	msg := m.(*raw.Message)
	server.Info("LoginResponse get result %v", msg)
	resp := s2s.LoginResp{}
	if err := codec.PBCodec.Unmarshal((msg.Body).([]byte), &resp); err != nil {
		fmt.Println(err)
		return nil
	}
	server.Info("LoginResponse get body %v", resp)

	return nil
}
