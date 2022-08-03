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
	"testing"
	"time"

	"github.com/airkits/nethopper/mq"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/network/common"
	"github.com/airkits/nethopper/network/natsrpc"
	"github.com/airkits/proto/ss"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestNatsClientRequest(t *testing.T) {
	conf := &natsrpc.ClientConfig{
		Nodes: []common.NodeInfo{{
			ID:      0,
			Name:    "NAME",
			Address: "nats://192.168.1.178:4222",
		}},
		PingInterval:        30 * time.Second,
		MaxPingsOutstanding: 10,
		MaxReconnects:       10,
		QueueSize:           1000,
		SocketQueueSize:     1000,
		MaxMessageSize:      100,
	}
	client := natsrpc.NewClient(conf, func(conn network.IConn, uid uint64, token string) network.IAgent {
		a := network.NewAgent(NewAgentAdapter(conn), uid, token)

		any, _ := anypb.New(nil)
		msg := &ss.Message{
			ID:      1,
			UID:     uid,
			MsgID:   1,
			MsgType: mq.MTRequest,
			Seq:     1,
			Body:    any,
		}

		a.GetAdapter().WriteMessage(msg)
		return a
	}, func(agent network.IAgent) {
		fmt.Println("on error")
	})
	client.Run()

	client.Wait()
}
