package json

import (
	"errors"

	"github.com/airkits/nethopper/examples/model/common"
	"github.com/airkits/nethopper/server"
)

//CreateBody create message body
func CreateBody(msgType uint32, cmd string) (interface{}, error) {
	switch msgType {
	case server.MTRequest:
		return CreateRequestBody(cmd)
	case server.MTResponse:
		return CreateResponseBody(cmd)
	case server.MTNotify:
		return CreateNotifyBody(cmd)
	case server.MTBroadcast:
		return CreateBroadcastBody(cmd)
	default:
		return nil, errors.New("create body failed,unknow message type")
	}

}

//CreateRequestBody create request body
func CreateRequestBody(c string) (interface{}, error) {
	switch c {
	case common.CSLoginCmd:
		{
			return &LoginReq{}, nil
		}
	}
	return nil, errors.New("create body failed,can't find request body")
}

//CreateResponseBody create response body
func CreateResponseBody(c string) (interface{}, error) {
	switch c {
	case common.CSLoginCmd:
		{
			return &LoginResp{}, nil
		}
	}
	return nil, errors.New("create body failed,can't find body response body")
}

//CreateNotifyBody create notify body
func CreateNotifyBody(cmd string) (interface{}, error) {
	return nil, errors.New("create body failed,can't find body notify body")
}

//CreateBroadcastBody create broadcast body
func CreateBroadcastBody(cmd string) (interface{}, error) {
	return nil, errors.New("create body failed,can't find body broadcast body")
}

//WSBody body base

//LoginReq login request
type LoginReq struct {
	UID    string `form:"uid" json:"uid"`
	Passwd string `form:"passwd" json:"passwd"`
}

// Result response common struct
type Result struct {
	Code int32  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

//BaseResponse base response object
type BaseResponse struct {
	Result Result `form:"result" json:"result"`
}

//Setup init
func (b *BaseResponse) Setup() {

}

// OK set result ok
func (b *BaseResponse) OK() {
	b.Result.Code = 0
	b.Result.Msg = "ok"
}

// Error set result error code and msg
func (b *BaseResponse) Error(code int32, msg string) {
	b.Result.Code = code
	b.Result.Msg = msg
}

//LoginResp login response
type LoginResp struct {
	BaseResponse
	Data string `form:"data" json:"data"`
}
