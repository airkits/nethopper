package grpc

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/base/queue"
	"github.com/gonethopper/nethopper/examples/model/pb/ss"
	"github.com/gonethopper/nethopper/examples/simple_server/modules/grpc"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"github.com/gorilla/websocket"
)

//Config grpc conn config
type Config struct {
	Address        string
	MaxConnNum     int
	RWQueueSize    int
	MaxMessageSize uint32
}

//NewServer create grpc server
func NewServer(m map[string]interface{}, agentFunc network.AgentCreateFunc) *Server {
	s := new(Server)
	if err := s.ReadConfig(m); err != nil {
		panic(err)
	}
	s.NewAgent = agentFunc
	return s
}

// Server grpc server define
type Server struct {
	ss.UnimplementedRPCServer
	Address        string
	MaxConnNum     int
	RWQueueSize    int
	MaxMessageSize uint32
	NewAgent       network.AgentCreateFunc
	gs             *grpc.Server
	listener       net.Listener
	conns          ConnSet
	mutexConns     sync.Mutex
	wg             sync.WaitGroup
	q              queue.Queue
}

// ReadConfig config map
// m := map[string]interface{}{
//  "address":":12080",
//	"maxConnNum":1024,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
// //tls support
//  "certFile":"",
//  "keyFile":"",
// }
func (s *Server) ReadConfig(m map[string]interface{}) error {

	if err := server.ParseConfigValue(m, "address", ":12080", &s.Address); err != nil {
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
	ln, err := net.Listen("tcp", s.Address)
	if err != nil {
		server.Fatal("%v", err)
	}
	server.Info("websocket start listen:%s", s.Address)
	if s.CertFile != "" || s.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err != nil {
			server.Fatal("%v", err)
		}

		ln = tls.NewListener(ln, config)
	}

	s.ln = ln

	s.upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Duration(s.HTTPTimeout) * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			server.Info("connection header:%v", r.Header)
			return true
		}}

	s.httpServer = &http.Server{
		Addr:           s.Address,
		Handler:        s,
		ReadTimeout:    time.Duration(s.HTTPTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.HTTPTimeout) * time.Second,
		MaxHeaderBytes: 1024,
	}
	s.httpServer.Serve(s.ln)
}

//Close websocket server
func (s *Server) Close() {
	s.ln.Close()

	s.mutexConns.Lock()
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
	s.mutexConns.Unlock()

	s.wg.Wait()
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
