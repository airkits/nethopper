package json

//Header request header
type Header struct {
	ID      int32       `form:"id" json:"id"`
	Cmd     string      `form:"cmd" json:"cmd"`
	MsgType int32       `form:"msgType" json:"msgType"`
	Options string      `form:"options" json:"options"`
	Payload interface{} `form:"payload" json:"payload"`
}

//GetMsgType >
func (h *Header) GetMsgType() int32 {
	return h.MsgType
}

//GetCmd >
func (h *Header) GetCmd() string {
	return h.Cmd
}

//GetID >
func (h *Header) GetID() int32 {
	return h.ID
}

//GetOptions >
func (h *Header) GetOptions() interface{} {
	return h.Options
}

//GetPayload >
func (h *Header) GetPayload() interface{} {
	return h.Payload
}
