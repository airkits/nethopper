package natsrpc

import (
	"sync"
	"time"

	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/proto/ss"
	"github.com/nats-io/nats.go"
)

// NewClient create grpc client
func NewClient(conf config.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	c.Conf = conf.(*ClientConfig)
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc
	c.wg = &sync.WaitGroup{}
	return c
}

//Client nats client
type Client struct {
	sync.Mutex
	Conf       *ClientConfig
	NewAgent   network.AgentCreateFunc
	CloseAgent network.AgentCloseFunc
	conns      ConnSet
	wg         *sync.WaitGroup
}

// Run client start run
func (c *Client) Run() {
	c.init()
	for _, info := range c.Conf.Nodes {

		c.wg.Add(1)
		go c.connect(info.ID, info.Name, info.Address)

	}
}

func (c *Client) init() {
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
func (c *Client) Reconnect(natsConn *nats.Conn) {

}
func (c *Client) Disconnect(natsConn *nats.Conn) {

}
func (c *Client) connect(serverID int, name string, address string) error {
	defer c.wg.Done()

	nc, err := nats.Connect(c.Conf.Nodes[0].Address,
		nats.PingInterval(c.Conf.PingInterval*time.Second),
		nats.MaxPingsOutstanding(c.Conf.MaxPingsOutstanding),
		nats.MaxReconnects(c.Conf.MaxReconnects),
		nats.ReconnectHandler(c.Reconnect),
		nats.DisconnectHandler(c.Disconnect),
	)
	if err != nil {
		return err
	}
	stream := NewStream(nc)
	log.Info("[GRPCClient] grpc client create new connection to id:[%d] %s %s.", serverID, name, address)
	c.Lock()
	c.conns[stream] = struct{}{}
	c.Unlock()

	grpcConn := NewConn(stream, c.Conf.SocketQueueSize, c.Conf.MaxMessageSize)
	agent := c.NewAgent(grpcConn, uint64(serverID), name)

	agent.Run()

	grpcConn.Close()
	c.Lock()
	delete(c.conns, stream)
	c.Unlock()
	c.CloseAgent(agent)
	agent.OnClose()
	return nil
}

// Close client connections
func (c *Client) Close() {
	c.Lock()
	for conn := range c.conns {
		conn.Close()
	}
	c.conns = nil
	c.Unlock()
	c.wg.Wait()
}

//Call sync get response
func (c *Client) Call(serverID int, name string, address string) *ss.Message {

	return nil
}
