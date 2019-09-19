package server

import (
	"net"
	"sync"
	proto "github.com/golang/protobuf/proto"
	"github.com/gonethopper/queue"
	uuid "github.com/satori/go.uuid"
)

// IDStack store srcIDs,max
type IDStack struct {
	i    int
	data [MaxIDSequence]int32
}

// Push id to stack
func (s *IDStack) Push(v int32) {
	s.data[s.i] = v
	s.i++
}

// Pop id from stack
func (s *IDStack) Pop() (ret int32) {
	s.i--
	ret = s.data[s.i]
	return
}

// Reset data
func (s *IDStack) Reset() {
	s.i = 0
}

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
	sess.Reset()
	sess.IP = net.ParseIP(host)
	sess.Port = port
	sess.SrcID = srcID
	sess.PushSrcID(srcID)
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
	srcIDs    *IDStack
	Request   proto.Message
	Response  proto.Message
	Die       chan struct{} // session die signal, will be triggered by others
}

// Reset session set to default value
func (s *Session) Reset() {
	if s.srcIDs == nil {
		s.srcIDs = &IDStack{}
	} else {
		s.srcIDs.Reset()
	}
	s.IP = nil
	s.Port = ""
	s.SrcID = 0
	s.SessionID = ""
	s.Done = make(chan *Session)
	s.Request = nil
	s.Response = nil
	s.Die = make(chan struct{})
}

// PushSrcID add SrcID to message src seq
func (s *Session) PushSrcID(srcID int32) {
	s.srcIDs.Push(srcID)
}

// PopSrcID get last srcID
func (s *Session) PopSrcID() int32 {
	return s.srcIDs.Pop()
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
