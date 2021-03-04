package http

import (
	"sync"

	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/server"
	"github.com/gin-gonic/gin"
)

//NewServer create http server
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
	conns      ConnSet
	CloseAgent network.AgentCloseFunc
	mutexConns sync.Mutex
	wg         sync.WaitGroup
	gs         *gin.Engine
}

func (s *Server) web() {
	s.gs.Run(s.Conf.Address)

}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {
	if s.NewAgent == nil {
		server.Fatal("NewAgent must not be nil")
	}
	s.gs = gin.New()
	server.GO(s.web)
}

//Close websocket server
func (s *Server) Close() {

}
