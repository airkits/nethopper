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
