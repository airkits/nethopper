// MIT License

// Copyright (c) 2019 gonethopper

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * @Author: ankye
// * @Date: 2020-01-09 11:01:03
// * @Last Modified by:   ankye
// * @Last Modified time: 2020-01-09 11:01:03

package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/examples/snowflake/cmd"
	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/nethopper/utils"
)

const TIMEOUT = time.Second * 15

//Error Coce ID define
//Client ERROR Code ID 划分区域 0x0001 -- 0x4FFF
const (
	//客户端错误，一般是一些传参，校验错误等
	CSErrorCodeClientError = 0x0001
	//系统错误，一般表示系统未知错误，或者不希望用户了解具体情形的错误
	CSErrorCodeSystemError = 0x0002
	//系统超时
	CSErrorCodeTimeout = 0x0003
	CSErrorCodeDBError = 0x0004
)

// RegisterAPI register web api
// func RegisterAPI(router *mux.Router) {
// 	router.HandleFunc("/", Index)
// 	router.HandleFunc("/hello/:name", Hello)
// }

//NewAPIV1 v1版本api注册
func NewAPIV1(router *gin.RouterGroup) {
	RegisterCmdAPI(router)
}

//RegisterCmdAPI  api 初始化
func RegisterCmdAPI(group *gin.RouterGroup) {

	group.POST("/genuid", GenUID)
	group.POST("/genuids", GenUIDs)

	group.POST("/call/:mid/:cmd/:opt", Call)
}

//GenUIDReq request body
type GenUIDReq struct {
	Channel int32 `form:"channel" json:"channel"`
}

//GenUIDResp response body
type GenUIDResp struct {
	UID uint64 `form:"uid" json:"uid"`
}

//GenUIDsReq request body
type GenUIDsReq struct {
	Channel int32 `form:"channel" json:"channel"`
	Num     int32 `form:"num" json:"num"`
}

//GenUIDsResp response body
type GenUIDsResp struct {
	UIDs []uint64 `form:"uids" json:"uids"`
}

//HTTPSession 请求上下文，用于保存请求，用于reponse的上下文数据
type HTTPSession struct {
	SessionID string
	Context   *gin.Context //网络连接上下文

}

//NewHTTPSession 创建一个session
func NewHTTPSession(c *gin.Context) *HTTPSession {
	sess := new(HTTPSession)
	sess.Context = c

	if sess.SessionID = utils.GenUUID(); sess.SessionID == "" {
		server.Error("gen uuid failed")
		return nil
	}

	return sess
}

//Response response root struct
type Response struct {
	Code int         `form:"code" json:"code"`
	Msg  string      `form:"msg" json:"msg"`
	Data interface{} `form:"data" json:"data"`
}

//ResponseError 返回错误信息
func ResponseError(session *HTTPSession, code int, msg error) {
	if session.Context != nil {

		session.Context.JSON(http.StatusOK, &Response{
			Code: code,
			Msg:  msg.Error(),
			Data: nil,
		})
		server.Error("request [%s] response error. client address:[%s] errCode:[%d] msg:[%s]", session.Context.Request.URL.Path, session.Context.ClientIP(), code, msg)
	}

}

//ResponseSuccess 返回成功结果
func ResponseSuccess(session *HTTPSession, data interface{}) {

	if session.Context != nil {
		session.Context.JSON(http.StatusOK, &Response{
			Code: 0,
			Msg:  "ok",
			Data: data,
		})
		server.Error("request [%s] response success. client address:[%s] ", session.Context.Request.URL.Path, session.Context.ClientIP())

	}

}

// GenUID api index
// @Summary 登录
// @Tags http web 模块
// @version 1.0
// @Accept  json
// @Produce  json
// @Param   account body  http.GenUIDReq    true        "GenUIDReq"
// @Success 200 object Response 成功后返回值
// @Router /v1/genuid [post]
func GenUID(c *gin.Context) {
	defer server.TraceCost("GenUID")()
	session := NewHTTPSession(c)
	model := &GenUIDReq{}
	if err := c.BindJSON(model); err != nil {
		ResponseError(session, CSErrorCodeClientError, err)
		return
	}

	result, err2 := server.Call(server.ModuleIDLogic, cmd.CallIDGenUIDCmd, utils.RandomInt32(0, 1024), model.Channel)
	if err2 != nil {
		server.Info("message done, get err %s", err2.Error())
		ResponseError(session, CSErrorCodeClientError, err2)
	} else {
		server.Info("message done,get uid  %v", result.(uint64))

		ResponseSuccess(session, GenUIDResp{UID: result.(uint64)})
	}

}

// GenUIDs api index
// @Summary 登录
// @Tags http web 模块
// @version 1.0
// @Accept  json
// @Produce  json
// @Param   account body  http.GenUIDsReq    true        "GenUIDsReq"
// @Success 200 object Response 成功后返回值
// @Router /v1/genuids [post]
func GenUIDs(c *gin.Context) {
	defer server.TraceCost("GenUIDs")()
	session := NewHTTPSession(c)
	model := &GenUIDsReq{}
	if err := c.BindJSON(model); err != nil {
		ResponseError(session, CSErrorCodeClientError, err)
		return
	}

	result, err2 := server.Call(server.ModuleIDLogic, cmd.CallIDGenUIDsCmd, utils.RandomInt32(0, 1024), model.Channel, model.Num)
	if err2 != nil {
		server.Info("message done,get err %s", err2.Error())
		ResponseError(session, CSErrorCodeClientError, err2)
	} else {
		server.Info("message done,get uid  %v", result.([]uint64))

		ResponseSuccess(session, GenUIDsResp{UIDs: result.([]uint64)})
	}

}

// Call api call tool
// @Summary 登录
// @Tags http web 模块
// @version 1.0
// @Accept  multipart/form-data
// @Produce  json
// @Param  module query int true "module"
// @Param cmd query string true "cmd"
// @Param opt query int true "opt"
// @Param   data   formData string    true        "data"
// @Success 200 object Response 成功后返回值
// @Router /v1/call/:module/:cmd/:opt [post]
func Call(c *gin.Context) {
	defer server.TraceCost("Call")()
	session := NewHTTPSession(c)
	var data string
	var ok bool
	var err error
	module, ok := c.GetQuery("module")
	moduleInt, err := strconv.Atoi(module)
	cmd, ok := c.GetQuery("cmd")
	option, ok := c.GetQuery("opt")
	optionInt, err := strconv.Atoi(option)
	if err != nil || !ok {
		ResponseError(session, 500, errors.New("err data"))
		return
	}
	args := make([]interface{}, 0)

	data = c.PostForm("data")
	var model map[string]interface{}
	if err := codec.JSONCodec.Unmarshal([]byte(data), &model); err != nil {
		ResponseError(session, CSErrorCodeClientError, err)
		return
	}
	for _, col := range model {
		args = append(args, col)
	}

	result, err2 := server.Call(int32(moduleInt), cmd, int32(optionInt), args...)
	if err2 != nil {
		//server.Info("message done,get pwd  %v ,err %s", result.(string), err2.Error())
		ResponseError(session, CSErrorCodeClientError, err2)
	} else {
		//server.Info("message done,get pwd  %v", result.(string))
		ResponseSuccess(session, result)
	}
}
