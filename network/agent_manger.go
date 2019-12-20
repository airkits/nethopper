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
// * @Date: 2019-12-20 19:39:03
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-20 19:39:03

package network

import (
	"reflect"
	"sync"

	"github.com/gonethopper/nethopper/base/set"
)

//AgentManager manager agent
type AgentManager struct {
	agents     *set.HashSet
	authAgents map[string]Agent
}

var instance *AgentManager
var once sync.Once

//GetInstance agent manager instance
func GetInstance() *AgentManager {
	once.Do(func() {
		instance = &AgentManager{
			agents:     set.NewHashSet(),
			authAgents: make(map[string]Agent),
		}
	})
	return instance
}

//AddAgent add agent to manager
func (am *AgentManager) AddAgent(a Agent) {
	if a.IsAuth() {
		_, ok := am.authAgents[a.Token()]
		if !ok {
			am.authAgents[a.Token()] = a
		}
	} else {
		am.agents.Add(a)
	}
}

//RemoveAgent remove agent from manager
func (am *AgentManager) RemoveAgent(a Agent) {
	a.OnClose()
	if a.IsAuth() {
		storeAgent, ok := am.authAgents[a.Token()]
		if ok && reflect.DeepEqual(a, storeAgent) {
			delete(am.authAgents, a.Token())
		}
	} else {
		am.agents.Remove(a)
	}
}
