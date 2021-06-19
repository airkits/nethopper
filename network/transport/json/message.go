package json

//Message json struct
type Message struct {
	ID      uint32      `form:"id" json:"id"`
	UID     uint64      `form:"uid" json:"uid"`
	MsgID   uint32      `form:"msgID" json:"msgID"`
	MsgType uint32      `form:"msgType" json:"msgType"`
	Seq     uint32      `form:"seq" json:"seq"`
	Options string      `form:"options" json:"options"`
	Body    interface{} `form:"body" json:"body"`
}

//GetMsgType >
func (m *Message) GetMsgType() uint32 {
	return m.MsgType
}

//GetMsgID >
func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

//GetID >
func (m *Message) GetID() uint32 {
	return m.ID
}

//GetUID >
func (m *Message) GetUID() uint64 {
	return m.UID
}

//GetOptions >
func (m *Message) GetOptions() interface{} {
	return m.Options
}

//GetBody >
func (m *Message) GetBody() interface{} {
	return m.Body
}

//GetSeq >
func (m *Message) GetSeq() uint32 {
	return m.Seq
}
