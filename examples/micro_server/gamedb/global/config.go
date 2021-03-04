package global

import (
	"sync"

	"github.com/airkits/nethopper/database"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network/common"
	"github.com/airkits/nethopper/network/grpc"
	"github.com/airkits/nethopper/network/http"
)

// Config server config
type Config struct {
	Env   string             `default:"env"`
	Log   log.Config         `mapstructure:"log"`
	GPRC  grpc.ServerConfig  `mapstructure:"grpc"`
	Logic common.LogicConfig `mapstructure:"logic"`
	Mysql database.Config    `mapstructure:"mysql"`
	HTTP  http.ServerConfig  `mapstructure:"http"`
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
