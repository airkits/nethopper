package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gonethopper/libs/logs"
	"github.com/gonethopper/libs/utils"
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/examples/simple_server/common"
	"github.com/gonethopper/nethopper/server"
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

	group.POST("/", Index)
	group.POST("/call/:mid/:cmd/:opt", Call)
}

//LoginReq request body
type LoginReq struct {
	UID    int64  `form:"uid" json:"uid"`
	Passwd string `form:"passwd" json:"passwd"`
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
		logs.Error("gen uuid failed")
		return nil
	}

	return sess
}

type Response struct {
	Code int
	Msg  string
	Data interface{}
}

//ResponseError 返回错误信息
func ResponseError(session *HTTPSession, code int, msg error) {
	if session.Context != nil {

		session.Context.JSON(http.StatusOK, &Response{
			Code: code,
			Msg:  msg.Error(),
			Data: nil,
		})
		logs.Error("request [%s] response error. client address:[%s] errCode:[%d] msg:[%s]", session.Context.Request.URL.Path, session.Context.ClientIP(), code, msg)
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
		logs.Error("request [%s] response success. client address:[%s] ", session.Context.Request.URL.Path, session.Context.ClientIP())

	}

}

// Index api index
// @Summary 登录
// @Tags http web 模块
// @version 1.0
// @Accept  json
// @Produce  json
// @Param   account body  http.LoginReq    true        "LoginReq"
// @Success 200 object Response 成功后返回值
// @Router /v1/ [post]
func Index(c *gin.Context) {
	defer server.TraceCost("Index")()
	session := NewHTTPSession(c)
	model := &LoginReq{}
	if err := c.BindJSON(model); err != nil {
		ResponseError(session, CSErrorCodeClientError, err)
		return
	}

	result, err2 := server.Call(server.ModuleIDLogic, common.CallIDLoginCmd, int32(model.UID), strconv.FormatInt(model.UID, 10), model.Passwd)
	if err2 != nil {
		server.Info("message done,get pwd  %v ,err %s", result.(string), err2.Error())
		ResponseError(session, CSErrorCodeClientError, err2)
	} else {
		server.Info("message done,get pwd  %v", result.(string))

		ResponseSuccess(session, result)
	}

}

// Call api call tool
// @Summary 登录
// @Tags http web 模块
// @version 1.0
// @Accept  multipart/form-data
// @Produce  json
// @Param  mid query int true "mid"
// @Param cmd query string true "cmd"
// @Param opt query int true "opt"
// @Param   data   formData string    true        "data"
// @Success 200 object Response 成功后返回值
// @Router /v1/call/:mid/:cmd/:opt [post]
func Call(c *gin.Context) {
	defer server.TraceCost("Call")()
	session := NewHTTPSession(c)
	var data string
	var ok bool
	var err error
	mid, ok := c.GetQuery("mid")
	midInt, err := strconv.Atoi(mid)
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
	if err := codec.JSONCodec.Unmarshal([]byte(data), &model, nil); err != nil {
		ResponseError(session, CSErrorCodeClientError, err)
		return
	}
	for _, col := range model {
		args = append(args, col)
	}

	result, err2 := server.Call(int32(midInt), cmd, int32(optionInt), args...)
	if err2 != nil {
		server.Info("message done,get pwd  %v ,err %s", result.(string), err2.Error())
		ResponseError(session, CSErrorCodeClientError, err2)
	} else {
		server.Info("message done,get pwd  %v", result.(string))

		ResponseSuccess(session, result)
	}
}
