// MIT License

// Copyright (c) 2019 gonethopper

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * @Author: ankye
// * @Date: 2019-12-20 19:39:11
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-20 19:39:11

package network

import (
	"net"

	"github.com/gonethopper/nethopper/server"
)

//IAgent agent interface define
type IAgent interface {
	Run()
	OnClose()

	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})

	Token() string
	SetToken(string)
	IsAuth() bool

	GetAdapter() IAgentAdapter

	SendMessage(payload []byte) error
}

//NewAgent create new agent
func NewAgent(userData interface{}, adapter IAgentAdapter) IAgent {
	return &Agent{userData: userData, adapter: adapter}
}

//Agent base agent struct
type Agent struct {
	userData interface{}
	adapter  IAgentAdapter
	token    string
}

//Token get token
func (a *Agent) Token() string {
	return a.token
}

//GetAdapter get agent adapter
func (a *Agent) GetAdapter() IAgentAdapter {
	return a.adapter
}

//IsAuth if set token return true else return false
func (a *Agent) IsAuth() bool {
	return len(a.token) > 0
}

//SetToken set token
func (a *Agent) SetToken(token string) {
	a.token = token
}

//Run agent start run
//usage
//func (a *Agent) Run (
//  for {
// 	data, err := a.ReadMessage()
// 	if err != nil {
// 		server.Debug("read message: %v", err)
// 		break
// 	}
// 	out := make(map[string]interface{})
// 	if err := a.Codec().Unmarshal(data, &out, nil); err == nil {
// 		server.Info("receive message %v", out)
// 		out["seq"] = out["seq"].(float64) + 1
// 	} else {
// 		server.Error(err)
// 	}
// 	a.WriteMessage(out)
// }
// }
func (a *Agent) Run() {
	for {
		data, err := a.adapter.ReadMessage()
		if err != nil {
			server.Debug("read message: %v", err)
			break
		}
		a.adapter.ProcessMessage(data)
	}
}

// OnClose agent close
func (a *Agent) OnClose() {

}

// SendMessage send message to conn
func (a *Agent) SendMessage(payload []byte) error {
	return a.GetAdapter().WriteMessage(payload)
}

//LocalAddr get local addr
func (a *Agent) LocalAddr() net.Addr {
	return a.adapter.Conn().LocalAddr()
}

//RemoteAddr get remote addr
func (a *Agent) RemoteAddr() net.Addr {
	return a.adapter.Conn().RemoteAddr()
}

//Close agent close
func (a *Agent) Close() {
	a.adapter.Conn().Close()
}

//Destroy agent destory
func (a *Agent) Destroy() {
	a.adapter.Conn().Destroy()
}

//UserData get userdata
func (a *Agent) UserData() interface{} {
	return a.userData
}

//SetUserData set userdata
func (a *Agent) SetUserData(data interface{}) {
	a.userData = data
}
