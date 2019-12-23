package ws

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"github.com/gorilla/websocket"
)

//NewServer create ws server
func NewServer(m map[string]interface{}, agentFunc network.AgentCreateFunc) *Server {
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
	NewAgent network.AgentCreateFunc
	ln       net.Listener
	upgrader websocket.Upgrader
	conns    ConnSet

	mutexConns sync.Mutex
	wg         sync.WaitGroup
	httpServer *http.Server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.Debug("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(s.MaxMessageSize))

	s.wg.Add(1)
	defer s.wg.Done()

	s.mutexConns.Lock()
	if s.conns == nil {
		s.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(s.conns) >= s.MaxConnNum {
		s.mutexConns.Unlock()
		conn.Close()
		server.Debug("too many connections")
		return
	}
	s.conns[conn] = struct{}{}
	s.mutexConns.Unlock()
	var token = r.Header.Get("token")
	var agent network.IAgent
	var ok bool

	if len(token) > 0 {
		agent, ok = network.GetInstance().GetAuthAgent(token)
		if ok { //exist agent,kick out old connection
			network.GetInstance().RemoveAgent(agent)
		}
	}
	wsConn := NewConn(conn, s.RWQueueSize, s.MaxMessageSize)
	agent = s.NewAgent(wsConn)
	agent.SetToken(r.Header.Get("token"))
	network.GetInstance().AddAgent(agent)
	agent.Run()

	// cleanup
	wsConn.Close()
	s.mutexConns.Lock()
	delete(s.conns, conn)
	s.mutexConns.Unlock()
	agent.OnClose()
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
	if err := server.ParseConfigValue(m, "httpTimeout", 10, &s.HTTPTimeout); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "certFile", "", &s.CertFile); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "keyFile", "", &s.KeyFile); err != nil {
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
