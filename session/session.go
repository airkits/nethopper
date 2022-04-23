// MIT License

// Copyright (c) 2019 gonethopper

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * @Author: ankye
// * @Date: 2019-12-11 10:13:40
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-11 10:13:40

package session

import (
	"net"
	"sync"

	"github.com/airkits/nethopper/base/queue"
	uuid "github.com/satori/go.uuid"
	proto "google.golang.org/protobuf/proto"
)

// MaxIDSequence max srcid sequence
const MaxIDSequence = 10

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

// Session connection identify
type Session struct {
	IP        net.IP
	Port      string
	SrcID     int32 //module id
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
