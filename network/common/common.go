package common

const (
	// HeaderToken token key define
	HeaderToken = "token"
	// HeaderUID UID key define
	HeaderUID = "UID"
)

const (
	//PackageLengthSize package length size
	PackageLengthSize = 2
)

// NodeInfo server node info
type NodeInfo struct {
	ID      int    `mapstructure:"id"`
	Name    string `mapstructure:"name"`
	Address string `mapstructure:"address"`
}

// LogicConfig logic common config
type LogicConfig struct {
	QueueSize          int `mapstructure:"queue_size"`
	WorkerPoolCapacity int `mapstructure:"worker_pool_capacity"`
}

// GetQueueSize get module queue size
func (s *LogicConfig) GetQueueSize() int {
	return s.QueueSize
}
