package etcd

import (
	"context"

	"github.com/coreos/etcd/clientv3"
	log "github.com/gonethopper/libs/logs"
)

func GetWithPrefix(serviceKey string) (err error, vals []string) {
	return Get(serviceKey, true)
}

func Get(serviceKey string, withPrefix bool) (err error, vals []string) {
	vals = make([]string, 0)
	var (
		resp *clientv3.GetResponse
	)
	if withPrefix {
		resp, err = etcdClient.Get(context.Background(), serviceKey, clientv3.WithPrefix())

	} else {
		resp, err = etcdClient.Get(context.Background(), serviceKey)
	}
	if err != nil {
		return
	}
	for _, v := range resp.Kvs {
		vals = append(vals, string(v.Value))
	}
	return
}

func Set(key, val string) error {
	putResp, err := etcdClient.Put(context.Background(), key, val)
	if err != nil {

		log.Error("etcd:KV put failed[%v],key[%s],resp[%q]", err, key, putResp)
		return err
	}
	return nil
}
