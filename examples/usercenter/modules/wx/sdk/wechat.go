package sdk

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gonethopper/nethopper/network/http"
)

//API api server
const API = "https://api.weixin.qq.com"

//LoginURL login account balance url
const LoginURL = API + "/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"

//Login 微信登陆
func Login(appID string, appSecret string, code string) (map[string]interface{}, error) {
	url := fmt.Sprintf(LoginURL, appID, appSecret, code)
	var content string
	if err := http.Request(url, http.GET, http.RequestTypeText, nil, nil, http.ResponseTypeText, &content, http.ConnTimeoutMS, http.ServeTimeoutMS); err != nil {
		return nil, err
	}
	values := make(map[string]interface{})
	if err := json.Unmarshal([]byte(content), &values); err != nil {
		return nil, err
	}
	if _, ok := values["errcode"]; ok {
		return nil, errors.New(values["errmsg"].(string))
	}
	openID := values["openid"].(string)
	sessionKey := values["session_key"].(string)
	if len(openID) == 0 || len(sessionKey) == 0 {
		return nil, errors.New("openID failed")
	}
	return values, nil
}
