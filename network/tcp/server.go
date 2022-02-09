package tcp

import (
	"net"
	"sync"
	"time"

	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
)

//NewServer create tcp server
func NewServer(conf config.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	s.Conf = conf.(*ServerConfig)
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc

	return s
}

// Server tcp server define
type Server struct {
	Conf        *ServerConfig
	NewAgent    network.AgentCreateFunc
	tcpListener *net.TCPListener
	CloseAgent  network.AgentCloseFunc
	conns       ConnSet
	mutexConns  sync.Mutex
	wg          sync.WaitGroup
}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {

	// Listen and bind local ip
	tcpAddr, err := net.ResolveTCPAddr(s.Conf.Network, s.Conf.Address)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	s.tcpListener = listener
	log.Info("listening on: %s %s", s.Conf.Network, listener.Addr())
	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Warning("accept failed: %s", err.Error())
			continue
		}
		log.Info("receive one client peer %s", conn.RemoteAddr().String())
		// set socket read buffer
		conn.SetReadBuffer(s.Conf.ReadBufferSize)
		// set socket write buffer
		conn.SetWriteBuffer(s.Conf.WriteBufferSize)
		// start a goroutine for every incoming connection for reading
		//go conn
		go s.Transport(conn)
	}
}

//Transport tcp connection
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

//Close tcp server
func (s *Server) Close() {
	s.tcpListener.Close()

	s.mutexConns.Lock()
	for conn := range s.conns {
		conn.Close()
	}
	s.conns = nil
	s.mutexConns.Unlock()

	s.wg.Wait()
}
