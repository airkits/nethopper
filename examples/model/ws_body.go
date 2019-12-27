package model

import (
	"errors"

	"github.com/gonethopper/nethopper/server"
)

//CreateBody create message body
func CreateBody(msgType int, cmd string) (IWSBody, error) {
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
func CreateRequestBody(cmd string) (IWSBody, error) {
	switch cmd {
	case CSLoginCmd:
		return &LoginReq{}, nil
	}
	return nil, errors.New("create body failed,can't find request body")
}

//CreateResponseBody create response body
func CreateResponseBody(cmd string) (IWSBody, error) {
	switch cmd {
	case CSLoginCmd:
		return &LoginResp{}, nil
	}
	return nil, errors.New("create body failed,can't find body response body")
}

//CreateNotifyBody create notify body
func CreateNotifyBody(cmd string) (IWSBody, error) {
	return nil, errors.New("create body failed,can't find body notify body")
}

//CreateBroadcastBody create broadcast body
func CreateBroadcastBody(cmd string) (IWSBody, error) {
	return nil, errors.New("create body failed,can't find body broadcast body")
}

//WSBody body base
type WSBody struct {
}

//Setup init
func (b *WSBody) Setup() {

}

//LoginReq login request
type LoginReq struct {
	WSBody
	UID    int64  `form:"uid" json:"uid"`
	Passwd string `form:"passwd" json:"passwd"`
}

//LoginResp login response
type LoginResp struct {
	WSBody
	Result string `form:"result" json:"result"`
	Msg    string `form:"msg" json:"msg"`
	Code   int    `form:"code" json:"code"`
}
