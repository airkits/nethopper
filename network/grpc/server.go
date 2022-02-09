package grpc

import (
	"context"
	"net"
	"sync"

	"github.com/airkits/nethopper/base/queue"
	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/utils"
	"github.com/airkits/proto/ss"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//NewServer create grpc server
func NewServer(conf config.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	s.Conf = conf.(*ServerConfig)
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc

	return s
}

// Server grpc server define
type Server struct {
	ss.UnimplementedRPCServer
	Conf       *ServerConfig
	NewAgent   network.AgentCreateFunc
	CloseAgent network.AgentCloseFunc
	gs         *grpc.Server
	listener   net.Listener
	conns      ConnSet
	mutexConns sync.Mutex
	wg         sync.WaitGroup
	q          queue.Queue
}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {
	if s.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}
	s.conns = make(ConnSet)

	s.gs = grpc.NewServer()
	ss.RegisterRPCServer(s.gs, s)

	lis, err := net.Listen("tcp", s.Conf.Address)

	if err != nil {
		log.Error("failed to listen: %v", err)
		return
	}
	log.Info("grpc start listen:%s", s.Conf.Address)
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
	uid := uint64(0)
	// get context from stream
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		if md.Get("token") != nil {
			token = md.Get("token")[0]
			log.Info("token from header: %s", token)
		}
		if md.Get("UID") != nil {
			uidStr := md.Get("UID")[0]
			uid = utils.Str2Uint64(uidStr)
			log.Info("UID from header: %d", uid)
		}
	}

	log.Info("one client connection opened.")
	s.mutexConns.Lock()
	s.conns[stream] = struct{}{}
	s.mutexConns.Unlock()

	var agent network.IAgent
	conn := NewConn(stream, s.Conf.SocketQueueSize, s.Conf.MaxMessageSize)
	agent = s.NewAgent(conn, uid, token)

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

// Call implements call function
func (s *Server) Call(ctx context.Context, in *ss.Message) (*ss.Message, error) {
	return &ss.Message{}, nil
}
