package natsrpc

import (
	"time"
)

const NatsServiceKey = "gameconfig"

// ServiceGroup service group info
type ServiceGroup struct {
	Key     string `mapstructure:"key" json:"key"`
	Type    int    `mapstructure:"type" json:"type"`
	Version int    `mapstructure:"version" json:"version"`
	Mode    int    `mapstructure:"mode" json:"mode"`
	Hash    []int  `mapstructure:"hash"  json:"hash"`
}

// NatsConfig grpc client config
type NatsConfig struct {
	ServiceType         int            `mapstructure:"service_type"`
	ServiceID           int            `mapstructure:"service_id"`
	Nats                []string       `mapstructure:"nats"`
	Services            []ServiceGroup `mapstructure:"services"`
	PingInterval        time.Duration  `mapstructure:"ping_interval"`
	MaxPingsOutstanding int            `mapstructure:"max_ping_outstanding"`
	MaxReconnects       int            `mapstructure:"max_reconnects"`
	QueueSize           int            `mapstructure:"queue_size"`
	SocketQueueSize     int            `mapstructure:"socket_queue_size"`
	AsyncMaxPending     uint32         `mapstructure:"async_max_pending"`
	WorkerPoolCapacity  int            `mapstructure:"worker_pool_capacity"`
	WorkerPoolQueueSize int            `mapstructure:"worker_pool_queue_size"`
}

// GetQueueSize get module queue size
func (c *NatsConfig) GetQueueSize() int {
	return c.QueueSize
}
