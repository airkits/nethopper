package natsrpc

import (
	"time"

	"github.com/airkits/nethopper/network/common"
)

// ServiceInfo service node info
type ServiceInfo struct {
	ID      int    `mapstructure:"id"`
	Name    string `mapstructure:"name"`
	Subject string `mapstructure:"subject"`
}

// ServiceGroup service group info
type ServiceGroup struct {
	ID    int           `mapstructure:"id"`
	Name  string        `mapstructure:"name"`
	Nodes []ServiceInfo `mapstructure:"nodes"`
	Hash  []int         `mapstructure:"hash"`
}

// NatsConfig grpc client config
type NatsConfig struct {
	Nodes               []common.NodeInfo `mapstructure:"nodes"`
	ServiceGroup        ServiceGroup      `mapstructure:"service_group"`
	TargetGroup         ServiceGroup      `mapstructure:"target_group"`
	PingInterval        time.Duration     `mapstructure:"ping_interval"`
	MaxPingsOutstanding int               `mapstructure:"max_ping_outstanding"`
	MaxReconnects       int               `mapstructure:"max_reconnects"`
	QueueSize           int               `mapstructure:"queue_size"`
	SocketQueueSize     int               `mapstructure:"socket_queue_size"`
	MaxMessageSize      uint32            `mapstructure:"max_message_size"`
	WorkerPoolCapacity  int               `mapstructure:"worker_pool_capacity"`
	WorkerPoolQueueSize int               `mapstructure:"worker_pool_queue_size"`
}

// GetQueueSize get module queue size
func (c *NatsConfig) GetQueueSize() int {
	return c.QueueSize
}
