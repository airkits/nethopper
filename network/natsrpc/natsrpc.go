package natsrpc

import (
	"fmt"
	"sync"
	"time"

	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/proto/ss"
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
}

func (c *NatsRPC) Wait() {
	c.wg.Wait()
}

// Run client start run
func (c *NatsRPC) Run() {
	c.init()
	info := c.Conf.Nodes[0]

	c.wg.Add(1)
	go c.connect(info.ID, info.Name, info.Address)

}

func (c *NatsRPC) init() {
	c.Lock()
	defer c.Unlock()

	if c.NewAgent == nil {
		log.Fatal("[GRPCClient] NewAgent must not be nil")
	}
	if c.conns != nil {
		log.Fatal("[GRPCClient] client is running")
	}

	c.conns = make(ConnSet)
}
func (c *NatsRPC) Reconnect(natsConn *nats.Conn) {

}
func (c *NatsRPC) Disconnect(natsConn *nats.Conn) {

}

func (c *NatsRPC) GetAgent() network.IAgent {
	return c.agent
}
func (c *NatsRPC) connect(serverID int, name string, address string) error {
	defer c.wg.Done()

	nc, err := nats.Connect(c.Conf.Nodes[0].Address,
		nats.PingInterval(c.Conf.PingInterval*time.Second),
		nats.MaxPingsOutstanding(c.Conf.MaxPingsOutstanding),
		nats.MaxReconnects(c.Conf.MaxReconnects),
		nats.RetryOnFailedConnect(true),
		nats.ReconnectWait(5*time.Second),
		nats.ReconnectHandler(c.Reconnect),
		nats.DisconnectHandler(c.Disconnect),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	log.Info("[Client] client create new connection to id:[%d] %s %s.", serverID, name, address)
	c.Lock()
	c.conns[nc] = struct{}{}
	c.Unlock()

	natsConn := NewConn(nc, c.Conf.SocketQueueSize, c.Conf.MaxMessageSize)
	c.agent = c.NewAgent(natsConn, uint64(serverID), name)

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

// Call sync get response
func (c *NatsRPC) Call(serverID int, name string, address string) *ss.Message {

	return nil
}
