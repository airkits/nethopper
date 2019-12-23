package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gonethopper/libs/logs"
	"github.com/gonethopper/libs/utils"
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

//ResponseError 返回错误信息
func ResponseError(session *HTTPSession, code int, msg error) {
	if session.Context != nil {

		session.Context.JSON(http.StatusOK, gin.H{
			"code":    code,
			"message": msg.Error(),
		})
		logs.Error("request [%s] response error. client address:[%s] errCode:[%d] msg:[%s]", session.Context.Request.URL.Path, session.Context.ClientIP(), code, msg)
	}

}

//ResponseSuccess 返回成功结果
func ResponseSuccess(session *HTTPSession, data interface{}) {

	if session.Context != nil {
		session.Context.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    data,
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
// - name: uid
//   in: body
//   description: 用户id
//   type: int
//   required: true
// - name: passwd
//   in: body
//   description: 用户密码
//   type: string
//   required: true
// @Success 200 object gin.H 成功后返回值
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

// Index api index
// func Index(w http.ResponseWriter, r *http.Request) {
// 	defer server.TraceCost("Index")()
// 	fmt.Fprint(w, "Welcome!\n")
// 	token := context.Get(r, "token").(string)
// 	fmt.Fprint(w, token+"\n")

// 	sbody, e := ioutil.ReadAll(r.Body)
// 	defer r.Body.Close()
// 	if e != nil {
// 		w.WriteHeader(500)
// 		return
// 	}
// 	server.Info(string(sbody))
// 	var v = make(map[string]interface{})
// 	if err := codec.JSONCodec.Unmarshal(sbody, &v, nil); err != nil {
// 		server.Info(err)
// 		return
// 	}
// 	uid := v["uid"].(float64)
// 	pwd := v["passwd"].(string)
// 	result, err2 := server.Call(server.ModuleIDLogic, common.CallIDLoginCmd, int32(uid), strconv.FormatFloat(uid, 'f', -1, 64), pwd)
// 	if err2 != nil {
// 		server.Info("message done,get pwd  %v ,err %s", result.(string), err2.Error())
// 		fmt.Fprint(w, "login failed")
// 	} else {
// 		server.Info("message done,get pwd  %v", result.(string))
// 		fmt.Fprint(w, "login success")
// 	}

// }

// // Hello api hello
// func Hello(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Category: %v\n", vars["category"])
// }
