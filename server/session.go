package server

import (
	"net"
	"sync"

	"github.com/gonethopper/queue"
	uuid "github.com/satori/go.uuid"
)

// NewSessionPool new session pool
func NewSessionPool() *SessionPool {
	sp := &SessionPool{}
	sp.Pool = &sync.Pool{
		New: func() interface{} {
			m := &Session{}
			return m
		}}
	return sp
}

// SessionPool mamager session objests
type SessionPool struct {
	Pool *sync.Pool
	Objs sync.Map
}

// Alloc borrow session from pool
func (p *SessionPool) Alloc(srcID int32, host string, port string) *Session {
	sess := p.Pool.Get().(*Session)

	sess.IP = net.ParseIP(host)
	sess.Port = port
	sess.SrcID = srcID
	sess.Die = make(chan struct{})
	sess.MQ = queue.NewChanQueue(16)
	sess.Done = make(chan *Session)
	sess.SessionID = uuid.NewV4().String()
	p.Objs.Store(sess.SessionID, sess)
	return sess
}

// Free retrun session to pool
func (p *SessionPool) Free(sess *Session) {
	p.Objs.Delete(sess.SessionID)
	p.Pool.Put(sess)
}

// CreateSession from session pool
func CreateSession(srcID int32, host string, port string) *Session {
	return GSessionPool.Alloc(srcID, host, port)
}

// GetSession get Session By sessionID
func GetSession(sessionID string) *Session {
	if v, ok := GSessionPool.Objs.Load(sessionID); ok {
		return v.(*Session)
	}
	return nil

}

// RemoveSession remove from pool
func RemoveSession(sess *Session) {
	GSessionPool.Free(sess)
}

// Session connection identify
type Session struct {
	IP        net.IP
	Port      string
	SrcID     int32 //service id
	MQ        queue.Queue
	SessionID string
	Done      chan *Session
	Message   *Message
	Die       chan struct{} // session die signal, will be triggered by others
}

// Reset session set to default value
func (s *Session) Reset() {
	s.IP = nil
	s.Port = ""
	s.SrcID = 0
	s.SessionID = ""
	s.Done = make(chan *Session)
	s.Message = nil
	s.Die = make(chan struct{})
}

// NotifyDone tigger done notify
func (s *Session) NotifyDone() {
	select {
	case s.Done <- s:
	// ok
	default:
		// 阻塞情况处理,这里忽略
	}
}
