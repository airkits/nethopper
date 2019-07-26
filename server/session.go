package server

import (
	"net"

	"github.com/gonethopper/queue"
)

// NewSession create session
func NewSession(srcID int32, host string, port string) *Session {
	sess := &Session{
		IP:    net.ParseIP(host),
		Port:  port,
		SrcID: srcID,
		Q:     queue.NewChanQueue(16),
		Die:   make(chan struct{}),
	}
	return sess
}

// Session connection identify
type Session struct {
	IP        net.IP
	Port      string
	SrcID     int32 //service id
	Q         queue.Queue
	SessionID int64
	Die       chan struct{} // session die signal, will be triggered by others
}
