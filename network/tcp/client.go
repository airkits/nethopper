package tcp

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
)

// NewClient create tcp client
func NewClient(conf *ClientConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	c.Conf = conf
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc

	return c
}

//Client tcp client
type Client struct {
	sync.Mutex
	Conf       *ClientConfig
	NewAgent   network.AgentCreateFunc
	CloseAgent network.AgentCloseFunc
	conns      ConnSet
	wg         sync.WaitGroup
}

// Run client start run
func (c *Client) Run() {
	c.init()
	for _, info := range c.Conf.Nodes {
		for i := 0; i < c.Conf.ConnNum; i++ {
			c.wg.Add(1)
			go c.connect(info.ID, info.Name, info.Address)
		}
	}
}

func (c *Client) init() {
	c.Lock()
	defer c.Unlock()

	if c.NewAgent == nil {
		log.Fatal("NewAgent must not be nil")
	}
	if c.conns != nil {
		log.Fatal("client is running")
	}

	c.conns = make(ConnSet)

}

func (c *Client) dial(serverID int, address string) net.Conn {
	for {
		//conn, err := net.DialTimeout(c.Network, c.Address, time.Second*30)
		conn, err := net.Dial(c.Conf.Network, address)
		if err == nil {
			return conn
		}

		server.Warning("connect to %v error: %v", address, err)
		time.Sleep(c.Conf.ConnectInterval)
		continue
	}
}

func (c *Client) connect(serverID int, name string, address string) {
	defer c.wg.Done()

reconnect:
	conn := c.dial(serverID, address)
	if conn == nil {
		return
	}
	c.Lock()
	c.conns[conn] = struct{}{}
	c.Unlock()

	tcpConn := NewConn(conn, c.Conf.SocketQueueSize, c.Conf.MaxMessageSize, c.Conf.ReadDeadline)
	agent := c.NewAgent(tcpConn, c.Conf.UID, c.Conf.Token)
	agent.Run()

	// cleanup
	tcpConn.Close()
	c.Lock()
	delete(c.conns, conn)
	c.Unlock()
	c.CloseAgent(agent)
	agent.OnClose()

	if c.Conf.AutoReconnect {
		time.Sleep(c.Conf.ConnectInterval)
		goto reconnect
	}
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
