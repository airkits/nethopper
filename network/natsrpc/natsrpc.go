package natsrpc

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/airkits/nethopper/base"
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
	Conf       *NatsConfig
	NewAgent   network.AgentCreateFunc
	CloseAgent network.AgentCloseFunc
	conns      ConnSet
	wg         *sync.WaitGroup
	agent      network.IAgent
	services   sync.Map
}

func (c *NatsRPC) Wait() {
	c.wg.Wait()
}

// Run client start run
func (c *NatsRPC) Run() {
	c.init()
	c.wg.Add(1)
	go c.connect()

}

func (c *NatsRPC) init() {
	c.Lock()
	defer c.Unlock()

	if c.NewAgent == nil {
		log.Fatal("[NatsRPC] NewAgent must not be nil")
	}
	if c.conns != nil {
		log.Fatal("[NatsRPC] client is running")
	}

	c.conns = make(ConnSet)
}
func (c *NatsRPC) Reconnect(natsConn *nats.Conn) {
	log.Error("[NatsRPC] Reconnect")

}

func (c *NatsRPC) DisconnectError(natsConn *nats.Conn, err error) {
	log.Error("[NatsRPC] Connect failed,DisconnectError")

}
func (c *NatsRPC) ErrorHandler(natsConn *nats.Conn, sub *nats.Subscription, err error) {
	log.Error("[NatsRPC] Connect failed,ErrorHandler")

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
		nats.MaxReconnects(c.Conf.MaxReconnects),
		nats.RetryOnFailedConnect(true),
		nats.ReconnectWait(5*time.Second),
		nats.ReconnectHandler(c.Reconnect),
		nats.DisconnectErrHandler(c.DisconnectError),
		nats.ErrorHandler(c.ErrorHandler),
		nats.Timeout(10*time.Second),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	log.Info("[NatsRPC] Connect to %s ID:[%s] VERSION:[%s].", nc.ConnectedServerName(), nc.ConnectedServerId(), nc.ConnectedServerVersion())
	mp := nc.MaxPayload()
	log.Info("[NatsRPC] Maximum payload is %d MB", mp/(1024*1024))
	c.Lock()
	c.conns[nc] = struct{}{}
	c.Unlock()

	natsConn := NewConn(nc, c.Conf)
	c.RegisterConfig(natsConn)
	c.agent = c.NewAgent(natsConn, 0, nc.ConnectedServerId())

	c.agent.Run()

	natsConn.Close()
	c.Lock()
	delete(c.conns, nc)
	c.Unlock()
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
func (c *NatsRPC) RegisterConfig(natsConn network.IConn) error {
	kv, err := natsConn.(*Conn).GetKVBucket()
	if err != nil {
		return err
	}

	for i := 0; i < len(c.Conf.Services); i++ {
		c.LoadServiceInfo(kv, &c.Conf.Services[i])
	}

	go func() {

		// Create key watcher.
		wopts := []nats.WatchOpt{}
		watcher, err := kv.WatchAll(wopts...)
		if err != nil {
			fmt.Printf("ERROR: nats.KeyValue.WatchAll failed, err: %v", err)
		}
		for {
			select {
			case kve := <-watcher.Updates():
				if kve != nil {
					fmt.Printf("RECV: key: %v", kve)
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
			case <-time.After(base.TimeoutChanTime):
				continue
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
	for conn := range c.conns {
		conn.Close()
	}
	c.conns = nil
	c.Unlock()
	c.wg.Wait()
}
