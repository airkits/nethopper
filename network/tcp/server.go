package tcp

import (
	"net"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

//NewServer create tcp server
func NewServer(m map[string]interface{}, agentFunc network.AgentCreateFunc) network.IServer {
	s := new(Server)
	if err := s.ReadConfig(m); err != nil {
		panic(err)
	}
	s.NewAgent = agentFunc
	return s
}

// Server websocket server define
type Server struct {
	Config
	NewAgent    network.AgentCreateFunc
	tcpListener *net.TCPListener

	conns ConnSet

	mutexConns sync.Mutex
	wg         sync.WaitGroup
}

// ReadConfig config map
// m := map[string]interface{}{
// readBufferSize default 32767
// writeBufferSize default 32767
// address default :8888
// network default "tcp4"  use "tcp4/tcp6"
// readDeadline default 15
//	"maxConnNum":1024,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
// //tls support
//  "certFile":"",
//  "keyFile":"",
// }
func (s *Server) ReadConfig(m map[string]interface{}) error {

	if err := server.ParseConfigValue(m, "address", ":8888", &s.Address); err != nil {
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

	if err := server.ParseConfigValue(m, "readBufferSize", 32767, &s.ReadBufferSize); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "writeBufferSize", 32767, &s.WriteBufferSize); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "network", "tcp4", &s.Network); err != nil {
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
	tcpAddr, err := net.ResolveTCPAddr(s.Network, s.Address)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	s.tcpListener = listener
	server.Info("listening on: %s %s", s.Network, listener.Addr())
	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			server.Warning("accept failed: %s", err.Error())
			continue
		}
		// set socket read buffer
		conn.SetReadBuffer(s.ReadBufferSize)
		// set socket write buffer
		conn.SetWriteBuffer(s.WriteBufferSize)
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
	c := NewConn(conn, s.RWQueueSize, s.MaxMessageSize, s.ReadDeadline)
	agent = s.NewAgent(c)
	agent.SetToken("token")
	network.GetInstance().AddAgent(agent)
	agent.Run()

	// cleanup
	conn.Close()
	// s.mutexConns.Lock()
	// delete(s.conns, stream)
	// s.mutexConns.Unlock()
	agent.OnClose()
	return nil
}

//Close websocket server
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
