package sdk

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gonethopper/nethopper/examples/usercenter/model"
	"github.com/gonethopper/nethopper/network/http"
	"github.com/gonethopper/nethopper/server"
)

//API api server
const API = "https://api.weixin.qq.com"

//LoginURL login account balance url
const LoginURL = API + "/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

//Login 微信登陆
func Login(appID string, appSecret string, code string) (*model.WXUser, server.Result) {
	url := fmt.Sprintf(LoginURL, appID, appSecret, code)
	var content string
	if err := http.Request(url, http.GET, http.RequestTypeText, nil, nil, http.ResponseTypeText, &content, http.ConnTimeoutMS, http.ServeTimeoutMS); err != nil {
		return nil, server.Result{Code: -1, Err: err}
	}
	wxuser := &model.WXUser{}
	if err := json.Unmarshal([]byte(content), &wxuser); err != nil {
		return nil, server.Result{Code: -2, Err: err}
	}
	if wxuser.ErrCode == 0 {
		return wxuser, server.Result{Code: 0, Err: nil}
	}
	return nil, server.Result{Code: -3, Err: errors.New(wxuser.ErrMsg)}
}
