package model

import (
	"errors"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/codec/common"
	"github.com/gonethopper/nethopper/examples/model/json"
	"github.com/gonethopper/nethopper/examples/model/pb"
	"github.com/gonethopper/nethopper/examples/model/pb/cs"
	"github.com/gonethopper/nethopper/server"
)

//NewWSMessage create new websocket message
func NewWSMessage(uid string, cmd string, seq int32, msgType int32, userdata int32, codec codec.Codec) *WSMessage {
	m := &WSMessage{codec: codec}
	m.Head = m.CreateHeader(uid, cmd, seq, msgType, userdata)
	return m
}

//NewEmptyWSMessage create empty object to receive websocket message data
func NewEmptyWSMessage(codec codec.Codec) *WSMessage {
	m := &WSMessage{codec: codec}
	return m
}

//IHeader message header interface

//WSMessage message struct
type WSMessage struct {
	Head  interface{}
	Body  interface{}
	codec codec.Codec
}

//DecodeHead head decode
func (m *WSMessage) DecodeHead(payload []byte) error {
	head := &json.WSHeader{}
	if err := m.codec.Unmarshal(payload, &head, nil); err != nil {
		return err
	}
	m.Head = head
	return nil
}

//CreateHeader create message header
func (m *WSMessage) CreateHeader(uid string, cmd string, seq int32, msgType int32, userdata int32) interface{} {
	switch m.codec.Type() {
	case common.CodecTypeJSON:
		{
			return &json.WSHeader{UID: uid, Cmd: cmd, Seq: seq, UserData: userdata, MsgType: msgType}
		}
	case common.CodecTypePB:
		{
			return &cs.WSHeader{Uid: uid, Cmd: cmd, Seq: seq, Userdata: userdata, MsgType: msgType}

		}
	}
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

func (m *WSMessage) decodePBBody() error {
	head := m.Head.(*cs.WSHeader)
	var body proto.Message
	var err error
	if body, err = pb.CreateBody(head.MsgType, head.Cmd); err != nil {
		return err
	}
	if err = m.codec.Unmarshal(head.Payload, body, nil); err != nil {
		return err
	}

	m.Body = body
	return nil
}
func (m *WSMessage) decodeJSONBody() error {
	head := m.Head.(*json.WSHeader)
	var body json.IWSBody
	var err error
	if body, err = json.CreateBody(head.MsgType, head.Cmd); err != nil {
		return err
	}
	server.Info("type %s", reflect.TypeOf(head.Payload))
	switch head.Payload.(type) {
	case string:
		{
			if err = m.codec.Unmarshal([]byte((head.Payload).(string)), body, nil); err != nil {
				return err
			}
		}
	case []byte:
		{
			if err = m.codec.Unmarshal((head.Payload).([]byte), body, nil); err != nil {
				return err
			}
		}

	default:
		server.Error("receive unknown message %x", head.Payload)
	}

	m.Body = body
	return nil
}

//DecodeBody decode body,first you should decode head first
func (m *WSMessage) DecodeBody() error {
	if m.Head == nil {
		return errors.New("message head is null")
	}
	switch m.codec.Type() {
	case common.CodecTypeJSON:
		{
			return m.decodeJSONBody()
		}
	case common.CodecTypePB:
		{
			return m.decodePBBody()
		}
	}

	return nil
}

//encodeJSONBody encode body
func (m *WSMessage) encodeJSONBody(payload []byte) error {
	(m.Head).(*json.WSHeader).Payload = string(payload)
	return nil
}

//encodePBBody encode body
func (m *WSMessage) encodePBBody(payload []byte) error {
	(m.Head).(*cs.WSHeader).Payload = payload
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
	switch m.codec.Type() {
	case common.CodecTypeJSON:
		{
			return m.encodeJSONBody(payload)
		}
	case common.CodecTypePB:
		{
			return m.encodePBBody(payload)
		}
	default:
		{
			return errors.New("unknown codec type")
		}
	}
}
