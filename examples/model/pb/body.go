package pb

import (
	"errors"

	"github.com/gogo/protobuf/proto"
	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/model/pb/c2s"
	"github.com/gonethopper/nethopper/examples/model/pb/s2s"
	"github.com/gonethopper/nethopper/server"
)

//CreateBody create message body
func CreateBody(msgType int32, cmd string) (proto.Message, error) {
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
func CreateRequestBody(c string) (proto.Message, error) {
	switch c {
	case common.CSLoginCmd:
		return &c2s.LoginReq{}, nil
	case common.SSLoginCmd:
		return &s2s.LoginReq{}, nil
	}
	return nil, errors.New("create body failed,can't find request body")
}

//CreateResponseBody create response body
func CreateResponseBody(c string) (proto.Message, error) {
	switch c {
	case common.CSLoginCmd:
		return &c2s.LoginResp{}, nil
	case common.SSLoginCmd:
		return &s2s.LoginResp{}, nil
	}
	return nil, errors.New("create body failed,can't find body response body")
}

//CreateNotifyBody create notify body
func CreateNotifyBody(c string) (proto.Message, error) {
	return nil, errors.New("create body failed,can't find body notify body")
}

//CreateBroadcastBody create broadcast body
func CreateBroadcastBody(c string) (proto.Message, error) {
	return nil, errors.New("create body failed,can't find body broadcast body")
}
