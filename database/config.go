package database

import "time"

//NodeInfo server node info
type NodeInfo struct {
	ID     int    `mapstructure:"id"`
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

//Config db client config
type Config struct {
	Nodes           []NodeInfo    `mapstructure:"nodes"`
	ConnectInterval time.Duration `mapstructure:"connect_interval"`
	QueueSize       int           `mapstructure:"queue_size"`
}

//GetQueueSize get module queue size
func (c *Config) GetQueueSize() int {
	return c.QueueSize
}
