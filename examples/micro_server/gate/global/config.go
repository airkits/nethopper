package global

import (
	"sync"

	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/network/grpc"
	"github.com/gonethopper/nethopper/network/http"
	"github.com/gonethopper/nethopper/network/ws"
)

// Config server config
type Config struct {
	Env        string             `default:"env"`
	Log        log.Config         `mapstructure:"log"`
	GPRC       grpc.ServerConfig  `mapstructure:"grpc"`
	GPRCClient grpc.ClientConfig  `mapstructure:"grpc_client"`
	Logic      common.LogicConfig `mapstructure:"logic"`
	WSJSON     ws.ServerConfig    `mapstructure:"wsjson"`
	WSPB       ws.ServerConfig    `mapstructure:"wspb"`
	HTTP       http.ServerConfig  `mapstructure:"http"`
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
