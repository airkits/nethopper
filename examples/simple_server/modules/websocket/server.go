package websocket

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/server"
	"github.com/gorilla/websocket"
)

type WSServer struct {
	WsConfig
	NewAgent func(*WSConn) Agent
	ln       net.Listener
	handler  *WSHandler
}

type WSHandler struct {
	maxConnNum      int
	pendingWriteNum int
	maxMsgLen       uint32
	newAgent        func(*WSConn) Agent
	upgrader        websocket.Upgrader
	conns           WebsocketConnSet
	mutexConns      sync.Mutex
	wg              sync.WaitGroup
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		server.Debug("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(handler.maxMsgLen))

	handler.wg.Add(1)
	defer handler.wg.Done()

	handler.mutexConns.Lock()
	if handler.conns == nil {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(handler.conns) >= handler.maxConnNum {
		handler.mutexConns.Unlock()
		conn.Close()
		server.Debug("too many connections")
		return
	}
	handler.conns[conn] = struct{}{}
	handler.mutexConns.Unlock()

	wsConn := newWSConn(conn, handler.pendingWriteNum, handler.maxMsgLen)
	agent := handler.newAgent(wsConn)
	agent.Run()

	// cleanup
	wsConn.Close()
	handler.mutexConns.Lock()
	delete(handler.conns, conn)
	handler.mutexConns.Unlock()
	agent.OnClose()
}

func (s *WSServer) Start() {
	ln, err := net.Listen("tcp", s.Address)
	if err != nil {
		server.Fatal("%v", err)
	}

	if s.MaxConnNum <= 0 {
		s.MaxConnNum = 100
		server.Warning("invalid MaxConnNum, reset to %v", s.MaxConnNum)
	}
	if s.PendingWriteNum <= 0 {
		s.PendingWriteNum = 100
		server.Warning("invalid PendingWriteNum, reset to %v", s.PendingWriteNum)
	}
	if s.MaxMsgLen <= 0 {
		s.MaxMsgLen = 4096
		server.Warning("invalid MaxMsgLen, reset to %v", s.MaxMsgLen)
	}
	if s.HTTPTimeout <= 0 {
		s.HTTPTimeout = 10
		server.Warning("invalid HTTPTimeout, reset to %v", s.HTTPTimeout)
	}
	if s.NewAgent == nil {
		server.Fatal("NewAgent must not be nil")
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
	s.handler = &WSHandler{
		maxConnNum:      s.MaxConnNum,
		pendingWriteNum: s.PendingWriteNum,
		maxMsgLen:       s.MaxMsgLen,
		newAgent:        s.NewAgent,
		conns:           make(WebsocketConnSet),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: time.Duration(s.HTTPTimeout) * time.Second,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	httpServer := &http.Server{
		Addr:           s.Address,
		Handler:        s.handler,
		ReadTimeout:    time.Duration(s.HTTPTimeout) * time.Second,
		WriteTimeout:   time.Duration(s.HTTPTimeout) * time.Second,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
}

func (s *WSServer) Close() {
	s.ln.Close()

	s.handler.mutexConns.Lock()
	for conn := range s.handler.conns {
		conn.Close()
	}
	s.handler.conns = nil
	s.handler.mutexConns.Unlock()

	s.handler.wg.Wait()
}
