package http

import (
	"time"

	"github.com/gonethopper/nethopper/network/common"
)

//ClientConfig http client config
type ClientConfig struct {
	Nodes   []common.NodeInfo `mapstructure:"nodes"`
	Timeout time.Duration     `mapstructure:"timeout"`
}

//ServerConfig http server config
type ServerConfig struct {
	Address string `mapstructure:"address"`
}
