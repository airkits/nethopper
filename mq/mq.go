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

package mq

import (
	"sync"

	proto "google.golang.org/protobuf/proto"
)

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
	// ErrorCodeOK Define error code = 0 if success
	ErrorCodeOK = 0
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
func (p *MessagePool) Alloc(msgID int32, srcID int32, destID int32, msgType int8, cmd string, sessionID string) *Message {
	m := p.Pool.Get().(*Message)
	m.Reset()
	m.MsgID = msgID
	m.SrcID = srcID
	m.DestID = destID
	m.MsgType = msgType
	m.Cmd = cmd
	m.ErrCode = ErrorCodeOK
	m.SessionID = sessionID
	return m
}

// Free retrun message to pool
func (p *MessagePool) Free(m *Message) {
	m.Reset()
	p.Pool.Put(m)
}

// CreateMessage get message from pool
func CreateMessage(msgID int32, srcID int32, destID int32, msgType int8, cmd string, sessionID string) *Message {
	return GMessagePool.Alloc(msgID, srcID, destID, msgType, cmd, sessionID)
}

// RemoveMessage return message to pool
func RemoveMessage(m *Message) {
	GMessagePool.Free(m)
}

// MaxIDSequence max srcid sequence
const MaxIDSequence = 10

//Message mq Message
type Message struct {
	MsgID     int32
	SrcID     int32
	DestID    int32
	SessionID string
	MsgType   int8
	Cmd       string
	Body      proto.Message
	ErrCode   int32
}

// SetBody set message body
func (m *Message) SetBody(body proto.Message) {
	m.Body = body
}

// Reset message set to default value
func (m *Message) Reset() {
	m.MsgID = InvalidInt32
	m.SrcID = InvalidInt32
	m.DestID = InvalidInt32
	m.SessionID = ""
	m.ErrCode = ErrorCodeOK
	m.MsgType = MessageType
	m.Cmd = ""
	m.Body = nil
}
