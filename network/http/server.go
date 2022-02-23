package http

import (
	"sync"

	"github.com/airkits/nethopper/base"
	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/gin-gonic/gin"
)

//NewServer create http server
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
	conns      ConnSet
	CloseAgent network.AgentCloseFunc
	mutexConns sync.Mutex
	wg         *sync.WaitGroup
	gs         *gin.Engine
}

func (s *Server) web() {
	s.gs.Run(s.Conf.Address)

}

//ListenAndServe start serve
func (s *Server) ListenAndServe() {
	if s.NewAgent == nil {
		log.Error("NewAgent must not be nil")
	}
	s.gs = gin.New()
	base.GO(s.web)
}

//Close websocket server
func (s *Server) Close() {

}
