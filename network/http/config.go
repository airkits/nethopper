package http

import (
	"time"

	"github.com/gonethopper/nethopper/network/common"
)

//ClientConfig http client config
type ClientConfig struct {
	Nodes   []common.NodeInfo `yarm:"nodes"`
	Timeout time.Duration     `yarm:"timeout"`
}

//ServerConfig http server config
type ServerConfig struct {
	Address string `yaml:"address"`
}
