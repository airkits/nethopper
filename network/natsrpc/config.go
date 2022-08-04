package natsrpc

import (
	"time"

	"github.com/airkits/nethopper/network/common"
)

//ServiceInfo service node info
type ServiceInfo struct {
	ID      int    `mapstructure:"id"`
	Name    string `mapstructure:"name"`
	Subject string `mapstructure:"subject"`
}

//ServiceGroup service group info
type ServiceGroup struct {
	ID    int           `mapstructure:"id"`
	Name  string        `mapstructure:"name"`
	Nodes []ServiceInfo `mapstructure:"nodes"`
	Hash  []int         `mapstructure:"hash"`
}

//ClientConfig grpc client config
type ClientConfig struct {
	Nodes               []common.NodeInfo `mapstructure:"nodes"`
	ServiceGroup        ServiceGroup      `mapstructure:"service_group"`
	TargetGroup         ServiceGroup      `mapstructure:"target_group"`
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

//ClientConfig grpc client config
type ServerConfig struct {
	Nodes               []common.NodeInfo `mapstructure:"nodes"`
	ServiceGroup        ServiceGroup      `mapstructure:"service_group"`
	TargetGroup         ServiceGroup      `mapstructure:"target_group"`
	PingInterval        time.Duration     `mapstructure:"ping_interval"`
	MaxPingsOutstanding int               `mapstructure:"max_ping_outstanding"`
	MaxReconnects       int               `mapstructure:"max_reconnects"`
	QueueSize           int               `mapstructure:"queue_size"`
	SocketQueueSize     int               `mapstructure:"socket_queue_size"`
	MaxMessageSize      uint32            `mapstructure:"max_message_size"`
}

//GetQueueSize get module queue size
func (c *ServerConfig) GetQueueSize() int {
	return c.QueueSize
}
