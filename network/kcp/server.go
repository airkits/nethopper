package kcp

import (
	"net"
	"sync"
	"time"

	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/server"
	"github.com/xtaci/kcp-go"
)

//NewServer create kcp server
func NewServer(conf server.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	s.Conf = conf.(*ServerConfig)
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc
	return s
}

// Server kcp server define
type Server struct {
	Conf        *ServerConfig
	NewAgent    network.AgentCreateFunc
	CloseAgent  network.AgentCloseFunc
	kcpListener *kcp.Listener
	conns       ConnSet
	mutexConns  sync.Mutex
	wg          sync.WaitGroup
}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {

	// Listen and bind local ip
	listener, err := kcp.Listen(s.Conf.Address)
	if err != nil {
		panic(err)
	}
	server.Info("KCP listening on: %s", listener.Addr())
	s.kcpListener = listener.(*kcp.Listener)

	if err := s.kcpListener.SetReadBuffer(s.Conf.UDPSocketBufferSize); err != nil {
		server.Error("SetReadBuffer", err)
	}
	if err := s.kcpListener.SetWriteBuffer(s.Conf.UDPSocketBufferSize); err != nil {
		server.Error("SetWriteBuffer", err)
	}
	if err := s.kcpListener.SetDSCP(s.Conf.Dscp); err != nil {
		server.Error("SetDSCP", err)
	}
	// loop accepting
	for {
		conn, err := s.kcpListener.AcceptKCP()
		if err != nil {
			server.Warning("accept failed: %s", err.Error())
			continue
		}
		// set kcp parameters
		conn.SetWindowSize(s.Conf.Sndwnd, s.Conf.Rcvwnd)
		conn.SetNoDelay(s.Conf.Nodelay, s.Conf.Interval, s.Conf.Resend, s.Conf.Nc)
		conn.SetStreamMode(true)
		conn.SetMtu(s.Conf.Mtu)
		// start a goroutine for every incoming connection for reading
		//go conn
		go s.Transport(conn)
	}
}

//Transport kcp connection
func (s *Server) Transport(conn net.Conn) error {

	s.wg.Add(1)
	defer s.wg.Done()

	// s.conns[stream] = struct{}{}
	// s.mutexConns.Unlock()

	var agent network.IAgent
	c := NewConn(conn, s.Conf.SocketQueueSize, s.Conf.MaxMessageSize, s.Conf.ReadDeadline*time.Second)
	agent = s.NewAgent(c, 0, "")
	agent.Run()

	// cleanup
	conn.Close()
	// s.mutexConns.Lock()
	// delete(s.conns, stream)
	// s.mutexConns.Unlock()
	s.CloseAgent(agent)
	agent.OnClose()
	return nil
}

//Close kcp server
func (s *Server) Close() {
	s.kcpListener.Close()

	s.mutexConns.Lock()
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
	s.mutexConns.Unlock()

	s.wg.Wait()
}
