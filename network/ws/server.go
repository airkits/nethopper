package ws

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/network/common"
	"github.com/airkits/nethopper/server"
	"github.com/airkits/nethopper/utils/conv"
	"github.com/gorilla/websocket"
)

//NewServer create ws server
func NewServer(conf server.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	s.Conf = conf.(*ServerConfig)
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc

	return s
}

// Server websocket server define
type Server struct {
	Conf       *ServerConfig
	NewAgent   network.AgentCreateFunc
	ln         net.Listener
	upgrader   websocket.Upgrader
	conns      ConnSet
	CloseAgent network.AgentCloseFunc
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
	conn.SetReadLimit(int64(s.Conf.MaxMessageSize))

	s.wg.Add(1)
	defer s.wg.Done()

	s.mutexConns.Lock()
	if s.conns == nil {
		s.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(s.conns) >= s.Conf.MaxConnNum {
		s.mutexConns.Unlock()
		conn.Close()
		server.Debug("too many connections")
		return
	}
	s.conns[conn] = struct{}{}
	s.mutexConns.Unlock()
	query := r.URL.Query()
	token := query.Get(common.HeaderToken)
	uid := query.Get(common.HeaderUID)
	if len(token) <= 0 {
		token = r.Header.Get(common.HeaderToken)
		uid = r.Header.Get(common.HeaderUID)
	}
	userID := conv.Str2Uint64(uid)
	if len(token) > 0 && userID > 0 {

		var agent network.IAgent

		wsConn := NewConn(conn, s.Conf.SocketQueueSize, s.Conf.MaxMessageSize)
		agent = s.NewAgent(wsConn, userID, token)

		agent.Run()

		// cleanup
		wsConn.Close()
		s.mutexConns.Lock()
		delete(s.conns, conn)
		s.mutexConns.Unlock()
		s.CloseAgent(agent)
		agent.OnClose()
	} else {
		s.mutexConns.Lock()
		delete(s.conns, conn)
		s.mutexConns.Unlock()
		conn.Close()
	}

}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {
	if s.NewAgent == nil {
		server.Fatal("NewAgent must not be nil")
	}
	s.conns = make(ConnSet)
	ln, err := net.Listen("tcp", s.Conf.Address)
	if err != nil {
		server.Fatal("%v", err)
	}
	server.Info("websocket start listen:%s", s.Conf.Address)
	if s.Conf.CertFile != "" || s.Conf.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(s.Conf.CertFile, s.Conf.KeyFile)
		if err != nil {
			server.Fatal("%v", err)
		}

		ln = tls.NewListener(ln, config)
	}

	s.ln = ln

	s.upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Duration(s.Conf.HTTPTimeout) * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			server.Info("connection header:%v", r.Header)
			return true
		}}

	s.httpServer = &http.Server{
		Addr:           s.Conf.Address,
		Handler:        s,
		ReadTimeout:    time.Duration(s.Conf.HTTPTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.Conf.HTTPTimeout) * time.Second,
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
