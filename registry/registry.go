package registry

import (
	"time"

	"google.golang.org/grpc/resolver"
)

const (
	//ActionExpire expire action
	ActionExpire = "EXPIRE"
	//ActionPut put action
	ActionPut = "PUT"
	//ActionDel del action
	ActionDel = "DELETE"
)

//IRegistry services discovery interface
type IRegistry interface {
	// Register 注册service地址到ETCD组件中
	Register(serviceKey string, val string, interval time.Duration, ttl time.Duration)
	// WithAlive 创建租约
	WithAlive(name string, addr string, ttl int64) error
	// UnRegister remove service from etcd
	UnRegister(serviceKey string)
	//Watcher 接收数据信息
	Watcher(serviceKey string, callback func(action string, key, val []byte))

	// NewResolver initialize an etcd client
	NewResolver(etcdAddr string) resolver.Builder
}
