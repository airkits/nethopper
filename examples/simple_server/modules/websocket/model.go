package websocket

//LoginCmd do Login action
type LoginCmd struct {
	UID    int64       `form:"uid" json:"uid"`
	Passwd string      `form:"passwd" json:"passwd"`
	Seq    int64       `form:"seq" json:"seq"`
	Data   interface{} `form:"data" json:"data"`
}
