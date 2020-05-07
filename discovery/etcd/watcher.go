package etcd

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	log "github.com/gonethopper/libs/logs"
)

const (
	Expire = "EXPIRE"
	Put    = "PUT"
	Del    = "DELETE"
)

type EtcdWatchCallback func(action string, key, val []byte)

/**
 * callback的处理必须是非阻塞的
 */
func Watcher(serviceKey string, callback EtcdWatchCallback) {
	watchChan := etcdClient.Watch(context.Background(), serviceKey, clientv3.WithPrefix())
	log.Debug("start watch: %s\n", serviceKey)
	for wresp := range watchChan {
		for _, ev := range wresp.Events {
			log.Debug("etcd: watch[%s],ev[%q]", serviceKey, ev)
			if ev.Type.String() == Put {
				callback(Put, ev.Kv.Key, ev.Kv.Value)
			} else if ev.Type.String() == Del {
				callback(Del, ev.Kv.Key, ev.Kv.Value)
			} else if ev.Type.String() == Put {
				callback(Put, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}
