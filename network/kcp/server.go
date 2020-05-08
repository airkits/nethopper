package kcp

import (
	"net"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"github.com/xtaci/kcp-go"
)

//NewServer create kcp server
func NewServer(m map[string]interface{}, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	if err := s.ReadConfig(m); err != nil {
		panic(err)
	}
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc
	return s
}

// Server kcp server define
type Server struct {
	Config
	NewAgent    network.AgentCreateFunc
	CloseAgent  network.AgentCloseFunc
	kcpListener *kcp.Listener
	conns       ConnSet
	mutexConns  sync.Mutex
	wg          sync.WaitGroup
}

// ReadConfig config map
// m := map[string]interface{}{
// udpSocketBuf default 4194304
// address default :14000
// readDeadline default 15
//	"maxConnNum":1024,
//  "socketQueueSize":100,
//  "maxMessageSize":4096,
//  "":0,
//  "":0,
//  "":0,
//  "":0,
//  "":0,
//  "":0,
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

	if err := server.ParseConfigValue(m, "udpSocketBuf", 4194304, &s.UDPSocketBufferSize); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "udpSndWnd", 32, &s.sndwnd); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "udpRcvWnd", 32, &s.rcvwnd); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "udpMtu", 1280, &s.mtu); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "dscp", 46, &s.dscp); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "nodelay", 1, &s.nodelay); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "interval", 20, &s.interval); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "resend", 1, &s.resend); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "nc", 1, &s.nc); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "readDeadline", 15, &s.ReadDeadline); err != nil {
		return err
	}
	s.ReadDeadline = s.ReadDeadline * time.Second
	return nil
}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {

	// Listen and bind local ip
	listener, err := kcp.Listen(s.Address)
	if err != nil {
		panic(err)
	}
	server.Info("KCP listening on: %s", listener.Addr())
	s.kcpListener = listener.(*kcp.Listener)

	if err := s.kcpListener.SetReadBuffer(s.UDPSocketBufferSize); err != nil {
		server.Error("SetReadBuffer", err)
	}
	if err := s.kcpListener.SetWriteBuffer(s.UDPSocketBufferSize); err != nil {
		server.Error("SetWriteBuffer", err)
	}
	if err := s.kcpListener.SetDSCP(s.dscp); err != nil {
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
		conn.SetWindowSize(s.sndwnd, s.rcvwnd)
		conn.SetNoDelay(s.nodelay, s.interval, s.resend, s.nc)
		conn.SetStreamMode(true)
		conn.SetMtu(s.mtu)
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
	c := NewConn(conn, s.RWQueueSize, s.MaxMessageSize, s.ReadDeadline)
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
