package global

import (
	"sync"

	"github.com/airkits/nethopper/examples/simple_client/modules/logic"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network/grpc"
	"github.com/airkits/nethopper/network/kcp"
	"github.com/airkits/nethopper/network/quic"
	"github.com/airkits/nethopper/network/tcp"
	"github.com/airkits/nethopper/network/ws"
)

//UserConfig use config user info
type UserConfig struct {
	Token string `mapstructure:"token"`
	UID   uint64 `mapstructure:"uid"`
}

// Config server config
type Config struct {
	Env   string            `default:"env"`
	User  UserConfig        `mapstructure:"user"`
	Log   log.Config        `mapstructure:"log"`
	GPRC  grpc.ClientConfig `mapstructure:"grpc_client"`
	KCP   kcp.ClientConfig  `mapstructure:"kcp_client"`
	QUIC  quic.ClientConfig `mapstructure:"quic_client"`
	TCP   tcp.ClientConfig  `mapstructure:"tcp_client"`
	WS    ws.ClientConfig   `mapstructure:"ws_client"`
	Logic logic.Config      `mapstructure:"logic"`
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

//GetUser config user id and token
func (c *ConfigManager) GetUser() *UserConfig {
	return &c.cfg.User
}
