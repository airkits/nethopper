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
// * @Date: 2019-06-12 15:53:22
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-12 15:53:22

package server

import "sync"

const (
	// MessageType message type enum
	MessageType = iota
	// MTRequest request =1
	MTRequest
	// MTResponse response = 2
	MTResponse
	// MTNotify notify = 3
	MTNotify
	// MTBroadcast broadcast = 4
	MTBroadcast
)
const (
	// InvalidInt32 Invalid values set to -1
	InvalidInt32 = -1
)

// NewMessagePool new message pool
func NewMessagePool() *MessagePool {
	mp := &MessagePool{}
	mp.Pool = &sync.Pool{
		New: func() interface{} {
			m := &Message{}
			return m
		}}
	return mp
}

// MessagePool mamager message objests
type MessagePool struct {
	Pool *sync.Pool
}

// Alloc borrow message from pool
func (p *MessagePool) Alloc(srcID int32, destID int32, msgType int8, cmd string, payLoad []byte) *Message {
	m := p.Pool.Get().(*Message)
	m.Reset()
	m.SrcIDs.Push(srcID)
	m.DestID = destID
	m.MsgType = msgType
	m.Cmd = cmd
	m.Payload = payLoad
	return m
}

// Free retrun message to pool
func (p *MessagePool) Free(m *Message) {
	m.Reset()
	p.Pool.Put(m)
}

// CreateMessage get message from pool
func CreateMessage(srcID int32, destID int32, msgType int8, cmd string, payLoad []byte) *Message {
	return GMessagePool.Alloc(srcID, destID, msgType, cmd, payLoad)
}

// RemoveMessage return message to pool
func RemoveMessage(m *Message) {
	GMessagePool.Free(m)
}

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

//Message mq Message
type Message struct {
	SrcIDs    *IDStack
	DestID    int32
	SessionID string
	MsgType   int8
	Cmd       string
	Payload   []byte
}

// Reset message set to default value
func (m *Message) Reset() {
	if m.SrcIDs == nil {
		m.SrcIDs = &IDStack{}
	} else {
		m.SrcIDs.Reset()
	}
	m.DestID = InvalidInt32
	m.SessionID = ""
	m.MsgType = MessageType
	m.Cmd = ""
	if len(m.Payload) > 0 {
		GBytesPool.Free(m.Payload)
		m.Payload = nil
	}
}

// PushSrcID add SrcID to message src seq
func (m *Message) PushSrcID(srcID int32) {
	m.SrcIDs.Push(srcID)
}

// PopSrcID get last srcID
func (m *Message) PopSrcID() int32 {
	return m.SrcIDs.Pop()
}
