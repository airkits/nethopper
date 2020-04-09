package json

//Message json struct
type Message struct {
	ID      int32       `form:"id" json:"id"`
	Cmd     string      `form:"cmd" json:"cmd"`
	MsgType int32       `form:"msgType" json:"msgType"`
	Options string      `form:"options" json:"options"`
	Body    interface{} `form:"body" json:"body"`
}

//GetMsgType >
func (m *Message) GetMsgType() int32 {
	return m.MsgType
}

//SetMsgType >
// func (m *Message) SetMsgType(v int32) {
// 	m.MsgType = v
// }

//GetCmd >
func (m *Message) GetCmd() string {
	return m.Cmd
}

//SetCmd >
// func (m *Message) SetCmd(v string) {
// 	m.Cmd = v
// }

//GetID >
func (m *Message) GetID() int32 {
	return m.ID
}

//SetID >
// func (m *Message) SetID(v int32) {
// 	m.ID = v
// }

//GetOptions >
func (m *Message) GetOptions() interface{} {
	return m.Options
}

//SetOptions >
// func (m *Message) SetOptions(v string) {
// 	m.Options = v
// }

//GetBody >
func (m *Message) GetBody() interface{} {
	return m.Body
}

//SetPayload >
// func (m *Message) SetPayload(v string) {
// 	m.Payload = v
// }
