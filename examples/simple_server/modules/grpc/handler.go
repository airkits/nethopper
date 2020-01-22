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
// * @Date: 2019-12-25 23:19:18
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-25 23:19:18

package grpc

import (
	"github.com/gonethopper/nethopper/examples/model/pb/ss"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

//LoginHandler request login
func LoginHandler(agent network.IAgentAdapter, m *ss.SSMessage, body interface{}) error {

	req := body.(*ss.LoginReq)
	server.Info("receive message %v", m)
	//userID := server.StringToInt64(req.Uid)
	//	result, err := server.Call(server.ModuleIDLogic, common.CallIDLoginCmd, int32(userID), req.Uid, req.Passwd)
	resp := &ss.LoginResp{
		Uid:  req.GetUid(),
		Name: req.GetPasswd(),
	}
	agent.WriteMessage(resp)
	server.Info("send message %v", resp)
	return nil
}
