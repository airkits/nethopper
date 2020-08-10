package model

//AppInfo info struct
type AppInfo struct {
	AppID     string `mapstructure:"appid"`
	AppSecret string `mapstructure:"app_secret"`
}

//WXConfig wechat app list
type WXConfig struct {
	Apps      []AppInfo `mapstructure:"apps"`
	QueueSize int       `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (c *WXConfig) GetQueueSize() int {
	return c.QueueSize
}

//WXUser info struct
type WXUser struct {
	OpenID     string `json:"openid"`
	UUID       string `json:"uuid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errMsg"`
}
