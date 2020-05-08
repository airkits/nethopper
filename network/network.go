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
// * @Date: 2019-12-20 19:39:19
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-20 19:39:19

package network

import (
	"net"
	"sync"

	"github.com/gonethopper/nethopper/base/set"
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/libs/skiplist"
)

// IClient network client interface
type IClient interface {
	ReadConfig(m map[string]interface{}) error
	Run()
	Close()
}

//IServer network server interface
type IServer interface {
	ReadConfig(m map[string]interface{}) error
	ListenAndServe()
	Close()
}

// IConn define network conn interface
type IConn interface {
	//ReadMessage read message/[]byte from conn
	ReadMessage() (interface{}, error)
	//WriteMessage write message/[]byte to conn
	WriteMessage(args ...interface{}) error
	//LocalAddr get local addr
	LocalAddr() net.Addr
	//RemoteAddr get remote addr
	RemoteAddr() net.Addr
	//Close conn
	Close()
	//Destory conn
	Destroy()
}

//IAgent agent interface define
type IAgent interface {
	Run()
	OnClose()
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UID() uint64
	SetUID(uint64)
	Token() string
	SetToken(string)
	IsAuth() bool
	GetAdapter() IAgentAdapter
	SendMessage(payload []byte) error
}

// IAgentAdapter agent adapter interface
type IAgentAdapter interface {
	//Setup AgentAdapter
	Setup(conn IConn, codec codec.Codec)
	//ProcessMessage process request and notify message
	ProcessMessage(payload interface{}) error

	//WriteMessage to connection
	WriteMessage(payload interface{}) error
	//ReadMessage goroutine not safe
	ReadMessage() (interface{}, error)
	// Codec get codec
	Codec() codec.Codec
	//SetCodec set codec
	SetCodec(c codec.Codec)
	//Conn get conn
	Conn() IConn
	// SetConn set conn
	SetConn(conn IConn)
	//GetSequence get inc id
	GetSequence() uint32
	//OnClose agent close and clear
	OnClose()
}

//AgentCreateFunc create agent func
type AgentCreateFunc func(conn IConn, uid uint64, token string) IAgent

//AgentCloseFunc close agent func
type AgentCloseFunc func(IAgent)

var instance *AgentManager
var once sync.Once

//GetInstance agent manager instance
func GetInstance() *AgentManager {
	once.Do(func() {
		instance = &AgentManager{
			agents:     set.NewHashSet(),
			authAgents: skiplist.New(),
		}
	})
	return instance
}

//AgentManager manager agent
type AgentManager struct {
	agents     *set.HashSet
	authAgents *skiplist.SkipList
}

//AddAgent add agent to manager
func (am *AgentManager) AddAgent(a IAgent) {
	if a.IsAuth() {
		v := am.authAgents.Get(float64(a.UID()))
		if v != nil {
			am.authAgents.Set(float64(a.UID()), a)
		}

	} else {
		am.agents.Add(a)
	}
}

//GetAuthAgent get auth agent,if exist return agent and true,else return false
func (am *AgentManager) GetAuthAgent(uid uint64) (IAgent, bool) {
	v := am.authAgents.Get(float64(uid))
	if v != nil {
		return v.Value().(IAgent), true
	}
	return nil, false
}

//RemoveAgent remove agent from manager
func (am *AgentManager) RemoveAgent(a IAgent) {
	a.OnClose()
	if a.IsAuth() {
		am.authAgents.Remove(float64(a.UID()))
	} else {
		am.agents.Remove(a)
	}
}
