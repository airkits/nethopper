package registry_test

import (
	"sync"
	"testing"
	"time"

	"github.com/gonethopper/nethopper/registry/etcd"
)

func TestEtcdClient(t *testing.T) {

	conf := &etcd.Config{
		UserName:    "",
		Password:    "",
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Duration(10) * time.Second,
		Interval:    time.Duration(5) * time.Second,
		TTL:         time.Duration(10) * time.Second,
		QueueSize:   100,
	}
	c, err := etcd.NewEtcdClient(conf)
	if err != nil {
		t.Error(err)
	}
	c.Set("/nethopper/nethopper", "adn")

	v, e := c.Get("/nethopper/nethopper", false)
	for _, b := range v {
		t.Logf("etcd value:[%s] len[%d]", b, len(v))

	}
	t.Logf("etcd get :value [%v]len [%d]", v, len(v))
	if e != nil {
		t.Logf("etcd get value error:[%v] ", e)
	}
	// //去注册
	c.Register("/nethopper/nethopper1", "127.0.0.1:1234", time.Duration(10)*time.Second, time.Duration(5)*time.Second)
	//	etcd.Register("/nethopper/nethopper1", "192.168.1.178:1234", time.Duration(10)*time.Second, time.Duration(5)*time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	go c.Watcher("/nethopper/nethopper1", func(action string, key, val []byte) {
		t.Logf("etcd callback:action[%s],key[%s],value[%s]", action, string(key), string(val))
		wg.Done()
	})
	wg.Wait()
	t.Logf("main exit...")
}
