package json

//Message json struct
type Message struct {
	ID      int32       `form:"id" json:"id"`
	Cmd     string      `form:"cmd" json:"cmd"`
	MsgType int32       `form:"msgType" json:"msgType"`
	Seq     int32       `form:"seq" json:"seq"`
	Options string      `form:"options" json:"options"`
	Body    interface{} `form:"body" json:"body"`
}

//GetMsgType >
func (m *Message) GetMsgType() int32 {
	return m.MsgType
}

//GetCmd >
func (m *Message) GetCmd() string {
	return m.Cmd
}

//GetID >
func (m *Message) GetID() int32 {
	return m.ID
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
func (m *Message) GetSeq() int32 {
	return m.Seq
}
