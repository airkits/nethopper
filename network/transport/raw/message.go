package raw

import (
	"github.com/gonethopper/nethopper/codec/raw"
)

//Message raw struct
type Message struct {
	ID      uint32
	Cmd     string
	MsgType uint32
	Seq     uint32
	Options string
	Body    interface{}
}

//Pack raw message
func (m *Message) Pack() []byte {

	coder := raw.NewCoder(nil, true)
	coder.SkipUint16()
	coder.WriteUint32(uint32(m.ID))
	coder.WriteString(m.Cmd)
	coder.WriteUint32(uint32(m.MsgType))
	coder.WriteUint32(uint32(m.Seq))
	coder.WriteString(m.Options)
	coder.WriteRaw(m.Body.([]byte))
	coder.SeekWriteUint16(0, uint16(coder.Length()))
	return coder.RawData()
}

//Unpack raw message
func (m *Message) Unpack(buffer []byte) error {
	coder := raw.NewCoder(buffer, true)
	m.ID = coder.ReadUint32()
	m.Cmd = coder.ReadString()
	m.MsgType = coder.ReadUint32()
	m.Seq = coder.ReadUint32()
	m.Options = coder.ReadString()
	m.Body = coder.ReadRaw()
	return nil
}

//GetID >
func (m *Message) GetID() uint32 {
	return m.ID
}

//GetCmd >
func (m *Message) GetCmd() string {
	return m.Cmd
}

//GetMsgType >
func (m *Message) GetMsgType() uint32 {
	return m.MsgType
}

//GetSeq >
func (m *Message) GetSeq() uint32 {
	return m.Seq
}
