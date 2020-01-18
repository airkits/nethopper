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
