package natsrpc

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/airkits/nethopper/codec/json"
	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/nats-io/nats.go"
)

// NewNatsRPC create nats client
func NewNatsRPC(conf config.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *NatsRPC {
	c := new(NatsRPC)
	c.Conf = conf.(*NatsConfig)
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc
	c.wg = &sync.WaitGroup{}

	return c
}

// NatsRPC nats PRC
type NatsRPC struct {
	sync.Mutex
	Conf        *NatsConfig
	NewAgent    network.AgentCreateFunc
	CloseAgent  network.AgentCloseFunc
	conn        network.IConn
	wg          *sync.WaitGroup
	agent       network.IAgent
	services    sync.Map
	watchClosed chan struct{}
}

func (c *NatsRPC) Wait() {
	c.wg.Wait()
}

// Run client start run
func (c *NatsRPC) Run() {
	c.init()
	c.wg.Add(1)
	c.connect()

}

func (c *NatsRPC) init() {
	c.Lock()
	defer c.Unlock()

	if c.NewAgent == nil {
		log.Fatal("[NatsRPC] type:[%d] id:[%d] NewAgent must not be nil", c.Conf.ServiceType, c.Conf.ServiceID)
	}
	if c.conn != nil {
		log.Fatal("[NatsRPC] type:[%d] id:[%d] client is running", c.Conf.ServiceType, c.Conf.ServiceID)
	}

}
func (c *NatsRPC) Reconnect(nc *nats.Conn) {
	log.Error("[NatsRPC] type:[%d] id:[%d] Reconnect %s", c.Conf.ServiceType, c.Conf.ServiceID, nc.ConnectedUrl())

	if err := c.conn.(*Conn).RegisterService(uint32(c.Conf.ServiceType), uint32(c.Conf.ServiceID)); err != nil {
		c.conn.(*Conn).ResetStream()
		time.Sleep(2 * time.Second)
		c.Reconnect(nc)
		return
	}
	if err := c.RegisterConfig(); err != nil {
		c.conn.(*Conn).ResetStream()
		time.Sleep(2 * time.Second)
		c.Reconnect(nc)
		return
	}
}

func (c *NatsRPC) DisconnectError(nc *nats.Conn, err error) {
	log.Error("[NatsRPC] type:[%d] id:[%d] Connect failed,DisconnectError %s", c.Conf.ServiceType, c.Conf.ServiceID, nc.Servers())
	c.conn.(*Conn).ResetStream()
	if c.watchClosed != nil {
		close(c.watchClosed)
		c.watchClosed = nil
	}
}
func (c *NatsRPC) ErrorHandler(nc *nats.Conn, sub *nats.Subscription, err error) {
	log.Error("[NatsRPC] type:[%d] id:[%d] ErrorHandler %s err:%s", c.Conf.ServiceType, c.Conf.ServiceID, sub.Subject, err.Error())

}
func (c *NatsRPC) natCloseHandler(nc *nats.Conn) {
	log.Error("[NatsRPC] type:[%d] id:[%d] natCloseHandler %s", c.Conf.ServiceType, c.Conf.ServiceID, nc.Servers())

}
func (c *NatsRPC) GetAgent() network.IAgent {
	return c.agent
}
func (c *NatsRPC) connect() error {
	defer c.wg.Done()
	// PingInterval ping间隔
	// MaxPingsOutstanding ping未响应次数
	// MaxReconnects 最大重连次数
	// RetryOnFailedConnect 失败重连
	// ReconnectWait 重连等待时间
	nc, err := nats.Connect(strings.Join(c.Conf.Nats, ","),
		nats.PingInterval(c.Conf.PingInterval*time.Second),
		nats.MaxPingsOutstanding(c.Conf.MaxPingsOutstanding),
		nats.MaxReconnects(-1),
		nats.RetryOnFailedConnect(true),
		nats.ReconnectWait(5*time.Second),
		nats.ReconnectHandler(c.Reconnect),
		nats.DisconnectErrHandler(c.DisconnectError),
		nats.ErrorHandler(c.ErrorHandler),
		nats.Timeout(10*time.Second),
		nats.ClosedHandler(c.natCloseHandler),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	log.Info("[NatsRPC] type:[%d] id:[%d] Connect to %s ID:[%s] VERSION:[%s].", c.Conf.ServiceType, c.Conf.ServiceID, nc.ConnectedServerName(), nc.ConnectedServerId(), nc.ConnectedServerVersion())
	mp := nc.MaxPayload()
	log.Info("[NatsRPC] type:[%d] id:[%d] Maximum payload is %d MB", c.Conf.ServiceType, c.Conf.ServiceID, mp/(1024*1024))

	c.conn = NewConn(nc, c.Conf)
	c.agent = c.NewAgent(c.conn, 0, nc.ConnectedServerId())
	c.conn.(*Conn).RegisterService(uint32(c.Conf.ServiceType), uint32(c.Conf.ServiceID))
	c.RegisterConfig()
	c.agent.Run()

	c.conn.Close()
	c.CloseAgent(c.agent)
	c.agent.OnClose()
	c.agent = nil
	return nil
}

func (c *NatsRPC) LoadServiceInfo(os nats.KeyValue, localInfo *ServiceGroup) error {

	result, err := os.Get(localInfo.Key)
	if err != nil {
		infoByte, err1 := json.Marshal(localInfo)
		if err1 != nil {
			return err1
		}
		os.PutString(localInfo.Key, string(infoByte))
		c.services.Store(localInfo.Type, localInfo)
		return err
	}
	remoteInfo := &ServiceGroup{}
	err = json.Unmarshal(result.Value(), remoteInfo)
	if err != nil {
		return err
	}
	if localInfo.Version > remoteInfo.Version {
		infoByte, err1 := json.Marshal(localInfo)
		if err1 != nil {
			return err1
		}
		os.PutString(localInfo.Key, string(infoByte))
		c.services.Store(localInfo.Type, localInfo)
	} else {
		c.services.Store(localInfo.Type, remoteInfo)
	}

	return nil
}
func (c *NatsRPC) RegisterConfig() error {
	kv, err := c.conn.(*Conn).GetKVBucket()
	if err != nil {
		return err
	}

	for i := 0; i < len(c.Conf.Services); i++ {
		if err := c.LoadServiceInfo(kv, &c.Conf.Services[i]); err != nil {
			return err
		}
	}
	c.watchClosed = make(chan struct{})
	go func() {

		// Create key watcher.
		wopts := []nats.WatchOpt{}
		watcher, err := kv.WatchAll(wopts...)
		if err != nil {
			log.Error("[NatsRPC] type:[%d] id:[%d]: nats.KeyValue.WatchAll failed, err: %v", c.Conf.ServiceType, c.Conf.ServiceID, err)
			return
		}
		for {
			select {
			case kve := <-watcher.Updates():
				if kve != nil {
					log.Info("[NatsRPC] type:[%d] id:[%d] RECV: key: %v", c.Conf.ServiceType, c.Conf.ServiceID, kve)
					for i := 0; i < len(c.Conf.Services); i++ {
						if c.Conf.Services[i].Key == kve.Key() {
							result := &ServiceGroup{}
							err := json.Unmarshal(kve.Value(), result)
							if err == nil && result.Version >= c.Conf.Services[i].Version {
								c.services.Store(c.Conf.Services[i].Type, result)
							}
						}
					}
				}
			case <-c.watchClosed:
				fmt.Println("watch close")
				return
			}
		}
	}()
	return nil
}
func (c *NatsRPC) GetHashValue(destType uint32, value uint64) uint32 {
	info, ok := c.services.Load(destType)

	if !ok {
		return 0
	}
	hashs := info.(ServiceGroup).Hash
	if hashs == nil {
		return 0
	}
	if info.(ServiceGroup).Mode == 1 {
		hashCode := int(value % uint64(len(hashs)))
		return uint32(hashs[hashCode])
	}
	return 0
}

// Close client connections
func (c *NatsRPC) Close() {
	c.Lock()
	c.conn.Close()
	c.Unlock()
	c.wg.Wait()
}
