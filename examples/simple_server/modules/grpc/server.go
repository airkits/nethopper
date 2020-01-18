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
// * @Date: 2020-01-09 11:20:09
// * @Last Modified by:   ankye
// * @Last Modified time: 2020-01-09 11:20:09

package grpc

import (
	"io"

	"github.com/gonethopper/nethopper/base/queue"
	"github.com/gonethopper/nethopper/examples/model/pb/ss"
	"github.com/gonethopper/nethopper/server"
)

// Server is used to implement ss.UnimplementedRPCServer
type Server struct {
	ss.UnimplementedRPCServer
	q queue.Queue
}

//Transport grpc connection
func (s *Server) Transport(stream ss.RPC_TransportServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := s.q.AsyncPush(msg); err != nil {
			server.Error("%s", err.Error())
		}
		m, err := s.q.AsyncPop()
		if err != nil {
			server.Error("%s", err.Error())
		} else {
			if err := stream.Send(m.(*ss.SSMessage)); err != nil {
				return err
			}
		}
	}
}
