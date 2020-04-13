package binary

//Message binary struct
type Message struct {
	ID      int32
	Cmd     string
	MsgType int32
	Seq     int32
	Options string
	Body    interface{}
}
