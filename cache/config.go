package cache

import (
	"time"
)

//NodeInfo server node info
type NodeInfo struct {
	ID       int    `mapstructure:"id"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
	DB       int    `mapstructure:"db"`
}

//Config grpc client config
type Config struct {
	Nodes           []NodeInfo    `mapstructure:"nodes"`
	ConnectInterval time.Duration `mapstructure:"connect_interval"`
	MaxActive       int           `mapstructure:"max_active"`
	MaxIdle         int           `mapstructure:"max_idle"`
	AutoReconnect   bool          `mapstructure:"auto_reconnect"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"`
	QueueSize       int           `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (c *Config) GetQueueSize() int {
	return c.QueueSize
}
