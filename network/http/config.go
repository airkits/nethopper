package http

import (
	"time"

	"github.com/airkits/nethopper/network/common"
)

//ClientConfig http client config
type ClientConfig struct {
	Nodes     []common.NodeInfo `mapstructure:"nodes"`
	Timeout   time.Duration     `mapstructure:"timeout"`
	QueueSize int               `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (c *ClientConfig) GetQueueSize() int {
	return c.QueueSize
}

//ServerConfig http server config
type ServerConfig struct {
	Address   string `mapstructure:"address"`
	QueueSize int    `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (s *ServerConfig) GetQueueSize() int {
	return s.QueueSize
}
