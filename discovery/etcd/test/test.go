package main

import (
	"time"

	"github.com/gonethopper/nethopper/discovery/etcd"
	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/server"
)

func main() {

	conf := log.Config{
		Filename:     "logs/server_log.log",
		Level:        7,
		MaxLines:     1000,
		MaxSize:      50,
		HourEnabled:  true,
		DailyEnabled: true,
		QueueSize:    1000,
	}
	server.NewNamedModule(server.MIDLog, "log", log.LogModuleCreate, nil, &conf)

	options := &etcd.Options{}
	options.Endpoints = []string{"127.0.0.1:2379"}
	options.DialTimeout = time.Duration(2) * time.Second

	etcd.NewEtcd(options)
	etcd.Set("/nethopper/nethopper", "adn")

	e, v := etcd.Get("/nethopper/nethopper", false)
	for _, b := range v {
		server.Debug("etcd value:[%s] len[%d]", b, len(v))

	}
	server.Debug("etcd get :value [%v]len [%d]", v, len(v))
	if e != nil {
		server.Debug("etcd get value error:[%v] ", e)
	}
	// //去注册
	etcd.Register("/nethopper/nethopper1", "127.0.0.1:1234", time.Duration(10)*time.Second, time.Duration(5)*time.Second)
	//	etcd.Register("/nethopper/nethopper1", "192.168.1.178:1234", time.Duration(10)*time.Second, time.Duration(5)*time.Second)
	go etcd.Watcher("/nethopper/nethopper1", watchBack)

	server.InitSignal()
}

func watchBack(action string, key, val []byte) {
	server.Debug("etcd callback:action[%s],key[%s],value[%s]", action, string(key), string(val))
}
