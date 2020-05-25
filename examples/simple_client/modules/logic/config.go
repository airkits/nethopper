package logic

//Config use logic config
type Config struct {
	UID                uint32 `mapstructure:"uid"`
	Password           string `mapstructure:"password"`
	QueueSize          int    `mapstructure:"queue_size"`
	WorkerPoolCapacity int    `mapstructure:"worker_pool_capacity"`
}

//GetQueueSize get queue size
func (c *Config) GetQueueSize() int {
	return c.QueueSize
}
