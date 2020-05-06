package grpc

import (
	"net"
	"sync"

	"github.com/gonethopper/nethopper/base/queue"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//Config grpc conn config
type Config struct {
	Address        string
	MaxConnNum     int
	RWQueueSize    int
	MaxMessageSize uint32
}

//NewServer create grpc server
func NewServer(m map[string]interface{}, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	if err := s.ReadConfig(m); err != nil {
		panic(err)
	}
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc

	return s
}

// Server grpc server define
type Server struct {
	ss.UnimplementedRPCServer
	Config
	NewAgent   network.AgentCreateFunc
	CloseAgent network.AgentCloseFunc
	gs         *grpc.Server
	listener   net.Listener
	conns      ConnSet
	mutexConns sync.Mutex
	wg         sync.WaitGroup
	q          queue.Queue
}

// ReadConfig config map
// m := map[string]interface{}{
//  "address":":14000",
//	"maxConnNum":1024,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
// //tls support
//  "certFile":"",
//  "keyFile":"",
// }
func (s *Server) ReadConfig(m map[string]interface{}) error {

	if err := server.ParseConfigValue(m, "address", ":14000", &s.Address); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "maxConnNum", 1024, &s.MaxConnNum); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "socketQueueSize", 100, &s.RWQueueSize); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "maxMessageSize", 4096, &s.MaxMessageSize); err != nil {
		return err
	}
	return nil
}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {
	if s.NewAgent == nil {
		server.Fatal("NewAgent must not be nil")
	}
	s.conns = make(ConnSet)

	s.gs = grpc.NewServer()
	ss.RegisterRPCServer(s.gs, s)

	lis, err := net.Listen("tcp", s.Address)

	if err != nil {
		server.Error("failed to listen: %v", err)
		return
	}
	server.Info("grpc start listen:%s", s.Address)
	s.listener = lis
	s.gs.Serve(lis)
}

//Close grpc server
func (s *Server) Close() {
	s.listener.Close()

	s.mutexConns.Lock()
	for conn := range s.conns {
		conn.Context().Done()
	}
	s.conns = nil
	s.mutexConns.Unlock()

	s.wg.Wait()
}

//Transport grpc connection
func (s *Server) Transport(stream ss.RPC_TransportServer) error {

	s.wg.Add(1)
	defer s.wg.Done()
	token := ""
	// get context from stream
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		if md.Get("token") != nil {
			token = md.Get("token")[0]
			server.Info("token from header: %s", token)
		}
	}

	server.Info("one client connection opened.")
	s.mutexConns.Lock()
	s.conns[stream] = struct{}{}
	s.mutexConns.Unlock()

	var agent network.IAgent
	conn := NewConn(stream, s.RWQueueSize, s.MaxMessageSize)

	agent = s.NewAgent(conn, token)

	agent.Run()

	// cleanup
	conn.Close()
	s.mutexConns.Lock()
	delete(s.conns, stream)
	s.mutexConns.Unlock()
	s.CloseAgent(agent)
	agent.OnClose()
	return nil
}
