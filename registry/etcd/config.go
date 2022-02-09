package etcd

// import (
// 	"time"
// )

// const (
// 	//ActionExpire expire action
// 	ActionExpire = "EXPIRE"
// 	//ActionPut put action
// 	ActionPut = "PUT"
// 	//ActionDel del action
// 	ActionDel = "DELETE"
// )

// //WatchCallback discover watch callback
// type WatchCallback func(action string, key, val []byte)

// //Config etcd client config
// type Config struct {
// 	UserName    string        `mapstructure:"username"`
// 	Password    string        `mapstructure:"password"`
// 	Endpoints   []string      `mapstructure:"endpoints"`
// 	DialTimeout time.Duration `mapstructure:"dial_timeout"`
// 	Interval    time.Duration `mapstructure:"interval"`
// 	TTL         time.Duration `mapstructure:"ttl"`
// 	QueueSize   int           `mapstructure:"queue_size"`
// }

// //GetQueueSize get module queue size
// func (conf *Config) GetQueueSize() int {
// 	return conf.QueueSize
// }
