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
// * @Date: 2019-12-11 10:13:10
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-11 10:13:10

package rpc_test

import (
	"fmt"

	"github.com/airkits/nethopper/mq"
	"github.com/airkits/proto/s2s"
	"github.com/airkits/proto/ss"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func HandlerHeartBeat(a *AgentAdapter, msg *ss.Message) error {
	req := &s2s.HeartBeatReq{}
	err := anypb.UnmarshalTo(msg.Body, req, proto.UnmarshalOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
	msg.MsgType = mq.MTResponse
	any, _ := anypb.New(&s2s.HeartBeatResp{Result: &s2s.Result{Code: 0, Msg: "ok"}, Time: req.Time})
	msg.Body = any
	return a.WriteMessage(msg)

}
