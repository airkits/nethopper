package global

import (
	"sync"

	"github.com/gonethopper/nethopper/cache"
	"github.com/gonethopper/nethopper/database"
	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/network/grpc"
	"github.com/gonethopper/nethopper/network/http"
	"github.com/gonethopper/nethopper/network/kcp"
	"github.com/gonethopper/nethopper/network/quic"
	"github.com/gonethopper/nethopper/network/tcp"
	"github.com/gonethopper/nethopper/network/ws"
)

// Config server config
type Config struct {
	Env   string             `default:"env"`
	Log   log.Config         `mapstructure:"log"`
	GPRC  grpc.ServerConfig  `mapstructure:"grpc"`
	KCP   kcp.ServerConfig   `mapstructure:"kcp"`
	QUIC  quic.ServerConfig  `mapstructure:"quic"`
	TCP   tcp.ServerConfig   `mapstructure:"tcp"`
	WS    ws.ServerConfig    `mapstructure:"wsjson"`
	Logic common.LogicConfig `mapstructure:"logic"`
	Mysql database.Config    `mapstructure:"mysql"`
	Redis cache.Config       `mapstructure:"redis"`
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
