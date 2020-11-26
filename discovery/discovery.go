package discovery

import "google.golang.org/grpc/resolver"

//IDiscovery services discovery interface
type IDiscovery interface {
	// Register 注册地址到ETCD组件中 使用 ; 分割
	Register(etcdAddr, name string, addr string, ttl int64) error
	// WithAlive 创建租约
	WithAlive(name string, addr string, ttl int64) error
	// UnRegister remove service from etcd
	UnRegister(name string, addr string)
	// NewResolver initialize an etcd client
	NewResolver(etcdAddr string) resolver.Builder
}
