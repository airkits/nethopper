package model

//SFConfig snowflake config list
type SFConfig struct {
	Hosts     []string `mapstructure:"hosts"`
	QueueSize int      `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (c *SFConfig) GetQueueSize() int {
	return c.QueueSize
}

//GenUIDReq request body
type GenUIDReq struct {
	Channel int32 `form:"channel" json:"channel"`
}

//GenUIDResp response body
type GenUIDResp struct {
	UID uint64 `form:"uid" json:"uid"`
}

//Response response root struct
type Response struct {
	Code int         `form:"code" json:"code"`
	Msg  string      `form:"msg" json:"msg"`
	Data interface{} `form:"data" json:"data"`
}
