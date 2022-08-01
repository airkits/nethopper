package natsrpc

import (
	"time"

	"github.com/airkits/nethopper/network/common"
)

//ClientConfig grpc client config
type ClientConfig struct {
	Nodes               []common.NodeInfo `mapstructure:"nodes"`
	PingInterval        time.Duration     `mapstructure:"ping_interval"`
	MaxPingsOutstanding int               `mapstructure:"max_ping_outstanding"`
	MaxReconnects       int               `mapstructure:"max_reconnects"`
	QueueSize           int               `mapstructure:"queue_size"`
	SocketQueueSize     int               `mapstructure:"socket_queue_size"`
	MaxMessageSize      uint32            `mapstructure:"max_message_size"`
}

//GetQueueSize get module queue size
func (c *ClientConfig) GetQueueSize() int {
	return c.QueueSize
}
