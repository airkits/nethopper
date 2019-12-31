package json

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/server"
)

//IWSBody websocket body interface
type IWSBody interface {
	Setup()
}

//WSHeader request header
type WSHeader struct {
	UID      string      `form:"uid" json:"uid"`
	CMD      string      `form:"cmd" json:"cmd"`
	Seq      int64       `form:"seq" json:"seq"`
	MsgType  int         `form:"msgType" json:"msgType"`
	UserData int         `form:"userdata" json:"userdata"`
	Payload  interface{} `form:"payload" json:"payload"`
}

//NewWSMessage create new websocket message
func NewWSMessage(uid string, cmd string, seq int64, msgType int, userdata int, codec codec.Codec) *WSMessage {
	m := &WSMessage{Head: &WSHeader{UID: uid, CMD: cmd, Seq: seq, UserData: userdata, MsgType: msgType}, codec: codec}
	return m
}

//NewEmptyWSMessage create empty object to receive websocket message data
func NewEmptyWSMessage(codec codec.Codec) *WSMessage {
	m := &WSMessage{codec: codec}
	return m
}

//WSMessage message struct
type WSMessage struct {
	Head  *WSHeader
	Body  IWSBody
	codec codec.Codec
}

//DecodeHead head decode
func (m *WSMessage) DecodeHead(payload []byte) error {
	head := &WSHeader{}
	if err := m.codec.Unmarshal(payload, &head, nil); err != nil {
		return err
	}
	m.Head = head
	return nil
}

//Encode encode message
func (m *WSMessage) Encode() ([]byte, error) {
	if m.Head == nil {
		return nil, errors.New("message head is null")
	}

	var err error
	if err = m.encodeBody(); err != nil {
		return nil, err
	}
	var payload []byte
	if payload, err = m.codec.Marshal(m.Head, nil); err != nil {
		return nil, err
	}
	return payload, nil
}

//DecodeBody decode body,first you should decode head first
func (m *WSMessage) DecodeBody() error {
	if m.Head == nil {
		return errors.New("message head is null")
	}
	var body IWSBody
	var err error
	if body, err = CreateBody(m.Head.MsgType, m.Head.CMD); err != nil {
		return err
	}
	server.Info("type %s", reflect.TypeOf(m.Head.Payload))
	switch m.Head.Payload.(type) {
	case string:
		{
			if err = m.codec.Unmarshal([]byte((m.Head.Payload).(string)), body, nil); err != nil {
				return err
			}
		}
	case []byte:
		{
			if err = m.codec.Unmarshal((m.Head.Payload).([]byte), body, nil); err != nil {
				return err
			}
		}

	default:
		server.Error("receive unknown message %x", m.Head.Payload)
	}

	m.Body = body
	return nil
}

//encodeBody encode body
func (m *WSMessage) encodeBody() error {
	if m.Head == nil {
		return errors.New("message head is null")
	}
	var payload []byte
	var err error
	if payload, err = m.codec.Marshal(m.Body, nil); err != nil {
		return err
	}
	if strings.Compare(m.codec.Name(), "JSONCodec") == 0 {
		m.Head.Payload = string(payload)
	} else {
		m.Head.Payload = payload
	}

	return nil
}
