package json

//WSHeader request header
type WSHeader struct {
	UID      string      `form:"uid" json:"uid"`
	Cmd      string      `form:"cmd" json:"cmd"`
	Seq      int32       `form:"seq" json:"seq"`
	MsgType  int32       `form:"msgType" json:"msgType"`
	UserData int32       `form:"userdata" json:"userdata"`
	Payload  interface{} `form:"payload" json:"payload"`
}

//GetMsgType >
func (h *WSHeader) GetMsgType() int32 {
	return h.MsgType
}

//GetPayload >
func (h *WSHeader) GetPayload() interface{} {
	return h.Payload
}
