package etcd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gonethopper/nethopper/server"
	"go.etcd.io/etcd/clientv3"
)

var (
	EtcdNoAuthClientError = errors.New("etcd auth client is nil")
)

//NewEtcdClient 生成etcd client 实例，在一个实例服务器中，有且只有一个etcd client实例
func NewEtcdClient(conf server.IConfig) (*Client, error) {
	c := &Client{}
	return c.Setup(conf)

}

//Client etcd client
type Client struct {
	etcdClient *clientv3.Client
	auth       clientv3.Auth
	Conf       *Config
	stopSignal chan bool
}

//Setup init
func (c *Client) Setup(conf server.IConfig) (*Client, error) {
	c.Conf = conf.(*Config)
	options := clientv3.Config{}
	options.Endpoints = c.Conf.Endpoints
	options.DialTimeout = c.Conf.DialTimeout
	options.Username = c.Conf.UserName
	options.Password = c.Conf.Password
	options.Context = context.Background()

	client, err := clientv3.New(options)
	if err != nil {
		server.Error("ETCD:create etcd client failed,error(%v)", err)
		return nil, err
	}

	if options.Username != "" {
		c.auth = clientv3.NewAuth(client)
	}
	c.stopSignal = make(chan bool, 1)
	c.etcdClient = client
	return c, nil
}

//Register 注册自己到etcd服务
func (c *Client) Register(serviceKey string, val string, interval time.Duration, ttl time.Duration) {
	//因为有ttl机制，所以需要设置ticker，来保证注册
	if serviceKey[0] != '/' {
		serviceKey = "/" + serviceKey
	}
	go c.register(serviceKey, val, interval, ttl)
}

func (c *Client) register(serviceKey string, val string, interval, ttl time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		//先查询这个key是否存在
		vals, err := c.Get(serviceKey, false)
		if err != nil {

			//申请租约
			resp, err := c.etcdClient.Grant(context.Background(), int64(ttl/time.Second))
			if err != nil {
				server.Error("etcd:grant failed[%v],key[%s],resp[%q]", err, serviceKey, resp)
			}
			//创建key-value 进去
			putResp, err := c.etcdClient.Put(context.Background(), serviceKey, val, clientv3.WithLease(resp.ID))
			if err != nil {
				server.Error("etcd:put failed[%v],key[%s],resp[%q]", err, serviceKey, putResp)
			}
		} else {
			if len(vals) == 0 {
				//申请租约
				resp, err := c.etcdClient.Grant(context.Background(), int64(ttl/time.Second))
				if err != nil {
					server.Error("etcd:grant failed[%v],key[%s],resp[%q]", err, serviceKey, resp)
				}
				//创建key-value 进去
				putResp, err := c.etcdClient.Put(context.Background(), serviceKey, val, clientv3.WithLease(resp.ID))
				if err != nil {
					fmt.Println("get error")
					server.Error("etcd:put failed[%v],key[%s],resp[%q]", err, serviceKey, putResp)

				}
			} else {
				server.Debug("etcd:registed key[%s] val[%v]", serviceKey, vals)
			}
		}
		//这里
		select {
		case <-c.stopSignal:
			return
		case <-ticker.C:
			//不做任何处理，一直循环
		}
	}
}

//UnRegister unregister service key
func (c *Client) UnRegister(serviceKey string) {

	//停止这个注册服务
	c.stopSignal <- true
	//重置
	c.stopSignal = make(chan bool, 1)
	resp, err := c.etcdClient.Delete(context.Background(), serviceKey)
	if err != nil {
		server.Error("etcd:delete failed [%v],key[%s] resp[%q]", err, serviceKey, resp)
	}
}

//Watcher callback的处理必须是非阻塞的
func (c *Client) Watcher(serviceKey string, callback WatchCallback) {
	watchChan := c.etcdClient.Watch(context.Background(), serviceKey, clientv3.WithPrefix())
	server.Debug("start watch: %s\n", serviceKey)
	for wresp := range watchChan {
		for _, ev := range wresp.Events {
			server.Debug("etcd: watch[%s],ev[%q]", serviceKey, ev)
			if ev.Type.String() == ActionPut {
				callback(ActionPut, ev.Kv.Key, ev.Kv.Value)
			} else if ev.Type.String() == ActionDel {
				callback(ActionDel, ev.Kv.Key, ev.Kv.Value)
			} else if ev.Type.String() == ActionExpire {
				callback(ActionExpire, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}

//Get from etcd
func (c *Client) Get(serviceKey string, withPrefix bool) (vals []string, err error) {
	vals = make([]string, 0)
	var (
		resp *clientv3.GetResponse
	)
	if withPrefix {
		resp, err = c.etcdClient.Get(context.Background(), serviceKey, clientv3.WithPrefix())

	} else {
		resp, err = c.etcdClient.Get(context.Background(), serviceKey)
	}
	if err != nil {
		return nil, err
	}
	for _, v := range resp.Kvs {
		vals = append(vals, string(v.Value))
	}
	return vals, err
}

//Set to etcd opt
func (c *Client) Set(key, val string) error {
	putResp, err := c.etcdClient.Put(context.Background(), key, val)
	if err != nil {

		server.Error("etcd:KV put failed[%v],key[%s],resp[%q]", err, key, putResp)
		return err
	}
	return nil
}

//Close 关闭
func (c *Client) Close() error {
	return c.etcdClient.Close()
}
