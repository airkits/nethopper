package global

import (
	"sync"

	"github.com/gonethopper/nethopper/cache"
	"github.com/gonethopper/nethopper/database"
	"github.com/gonethopper/nethopper/examples/usercenter/model"
	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/network/grpc"
	"github.com/gonethopper/nethopper/network/http"
)

const (
	// ModuleIDUserCustom User custom define named modules from 64-128
	ModuleIDUserCustom = 64
	// ModuleIDWechatClient module id wechat
	ModuleIDWechatClient
)

// Config server config
type Config struct {
	Env   string             `default:"env"`
	SID   int                `default:"id"`
	Log   log.Config         `mapstructure:"log"`
	GPRC  grpc.ServerConfig  `mapstructure:"grpc"`
	Logic common.LogicConfig `mapstructure:"logic"`
	Redis cache.Config       `mapstructure:"redis"`
	Mysql database.Config    `mapstructure:"mysql"`
	HTTP  http.ServerConfig  `mapstructure:"http"`
	WX    model.WXConfig     `mapstructure:"wx"`
}

//ConfigManager define
type ConfigManager struct {
	cfg Config
}

var instance *ConfigManager
var once sync.Once

//GetInstance agent manager instance
func GetInstance() *ConfigManager {
	once.Do(func() {
		instance = &ConfigManager{}
	})
	return instance
}

//GetConfig local config
func (c *ConfigManager) GetConfig() *Config {
	return &c.cfg
}
