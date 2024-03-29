package ws

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/network/common"
	"github.com/airkits/nethopper/utils"
	"github.com/gorilla/websocket"
)

//NewServer create ws server
func NewServer(conf config.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) network.IServer {
	s := new(Server)
	s.Conf = conf.(*ServerConfig)
	s.NewAgent = agentFunc
	s.CloseAgent = agentCloseFunc
	s.wg = &sync.WaitGroup{}
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
	wg         *sync.WaitGroup
	httpServer *http.Server
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debug("upgrade error: %v", err)
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
		log.Debug("too many connections")
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
	userID := utils.Str2Uint64(uid)
	if len(token) > 0 && userID > 0 {

		var agent network.IAgent

		wsConn := NewConn(conn, s.Conf.SocketQueueSize, s.Conf.MaxMessageSize)
		agent = s.NewAgent(wsConn, userID, token)
		log.Trace("[WS] one client connection opened, uid:[%d] token: %s success", userID, token)

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
		log.Fatal("[WS] NewAgent must not be nil")
	}
	s.conns = make(ConnSet)
	ln, err := net.Listen("tcp", s.Conf.Address)
	if err != nil {
		log.Fatal("%v", err)
	}
	log.Trace("[WS] websocket start listen %s", s.Conf.Address)
	if s.Conf.CertFile != "" || s.Conf.KeyFile != "" {
		config := &tls.Config{}
		config.NextProtos = []string{"http/1.1"}

		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(s.Conf.CertFile, s.Conf.KeyFile)
		if err != nil {
			log.Fatal("%v", err)
		}

		ln = tls.NewListener(ln, config)
	}

	s.ln = ln

	s.upgrader = websocket.Upgrader{
		HandshakeTimeout: time.Duration(s.Conf.HTTPTimeout) * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			log.Info("[WS] connection header:%v", r.Header)
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
