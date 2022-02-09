package etcd

// import (
// 	"context"
// 	"time"

// 	"github.com/coreos/etcd/clientv3"
// 	"go.etcd.io/etcd/clientv3/concurrency"
// )

// //EtcdLock etcd 分布式锁
// type EtcdLock struct {
// 	cli     *clientv3.Client
// 	key     string
// 	timeout time.Duration
// 	session *concurrency.Session
// 	mu      *concurrency.Mutex
// 	ctx     context.Context
// 	cancel  context.CancelFunc
// }

// //NewLocker NewLocker
// func NewLocker(cli *clientv3.Client, key string, timeout time.Duration) EtcdLock {
// 	if timeout <= 0 {
// 		timeout = 10 * time.Second
// 	}
// 	return EtcdLock{
// 		cli:     cli,
// 		key:     key,
// 		timeout: timeout,
// 	}
// }

// //Lock Lock
// func (c *EtcdLock) Lock() error {
// 	s, err := concurrency.NewSession(c.cli)
// 	if err != nil {
// 		return err
// 	}
// 	c.session = s
// 	c.mu = concurrency.NewMutex(s, c.key)

// 	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
// 	c.ctx = ctx
// 	c.cancel = cancel

// 	err = c.mu.Lock(c.ctx)
// 	return err
// }

// //Unlock Unlock
// func (c *EtcdLock) Unlock() error {
// 	err := c.mu.Unlock(c.ctx)
// 	c.cancel()
// 	c.session.Close()
// 	return err
// }
