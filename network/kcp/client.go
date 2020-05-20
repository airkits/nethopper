package kcp

import (
	"log"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"github.com/xtaci/kcp-go"
)

// NewClient create kcp client
func NewClient(conf *ClientConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	c.Conf = conf
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc

	return c
}

//Client kcp client
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

func (c *Client) dial(serverID int, address string) (*kcp.UDPSession, error) {
	conn, err := kcp.DialWithOptions(address, nil, 0, 0)
	if err == nil {
		return conn, nil
	}
	return nil, err
}

func (c *Client) connect(serverID int, name string, address string) {
	defer c.wg.Done()

reconnect:
	conn, err := c.dial(serverID, address)
	if err != nil {
		server.Fatal("kcp client connect to id:[%d] %s %s failed, reason: %v", serverID, name, address, err)
		if c.Conf.AutoReconnect {
			time.Sleep(c.Conf.ConnectInterval * time.Second)
			server.Warning("kcp client try reconnect to id:[%d] %s %s", serverID, name, address)
			goto reconnect
		}
	}
	conn.SetDSCP(c.Conf.Dscp)
	conn.SetWindowSize(c.Conf.Sndwnd, c.Conf.Rcvwnd)
	conn.SetNoDelay(c.Conf.Nodelay, c.Conf.Interval, c.Conf.Resend, c.Conf.Nc)
	conn.SetStreamMode(true)
	conn.SetMtu(c.Conf.Mtu)

	c.Lock()
	c.conns[conn] = struct{}{}
	c.Unlock()
	kcpConn := NewConn(conn, c.Conf.SocketQueueSize, c.Conf.MaxMessageSize, c.Conf.ReadDeadline*time.Second)
	agent := c.NewAgent(kcpConn, c.Conf.UID, c.Conf.Token)
	agent.Run()

	// cleanup
	kcpConn.Close()
	c.Lock()
	delete(c.conns, conn)
	c.Unlock()
	c.CloseAgent(agent)
	agent.OnClose()

	if c.Conf.AutoReconnect {
		time.Sleep(c.Conf.ConnectInterval * time.Second)
		server.Warning("kcp client try reconnect to id:[%d] %s %s", serverID, name, address)
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
