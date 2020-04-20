package raw

import (
	"github.com/gonethopper/nethopper/codec/raw"
)

//Message raw struct
type Message struct {
	ID      int32
	Cmd     string
	MsgType int32
	Seq     int32
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
	m.ID = int32(coder.ReadUint32())
	m.Cmd = coder.ReadString()
	m.MsgType = int32(coder.ReadUint32())
	m.Seq = int32(coder.ReadUint32())
	m.Options = coder.ReadString()
	m.Body = coder.ReadRaw()
	return nil
}

//GetID >
func (m *Message) GetID() int32 {
	return m.ID
}

//GetCmd >
func (m *Message) GetCmd() string {
	return m.Cmd
}

//GetMsgType >
func (m *Message) GetMsgType() int32 {
	return m.MsgType
}

//GetSeq >
func (m *Message) GetSeq() int32 {
	return m.Seq
}
