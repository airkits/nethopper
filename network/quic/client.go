package quic

import (
	"context"
	"crypto/tls"
	"log"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	quic "github.com/lucas-clemente/quic-go"
)

// NewClient create quic client
func NewClient(conf *ClientConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	c.Conf = conf
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc

	return c
}

//Client quic client
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

func (c *Client) dial(serverID int, address string) quic.Session {
	for {
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
		}
		session, err := quic.DialAddr(address, tlsConf, nil)
		if err == nil {
			return session
		}

		server.Warning("connect to %v error: %v", address, err)
		time.Sleep(c.Conf.ConnectInterval)
		continue
	}
}

func (c *Client) connect(serverID int, name string, address string) {
	defer c.wg.Done()

reconnect:
	sess := c.dial(serverID, address)
	if sess == nil {
		return
	}
	c.Lock()
	c.conns[sess] = struct{}{}
	c.Unlock()

	stream, err := sess.OpenStreamSync(context.Background())
	if err != nil {
		return
	}

	quicConn := NewConn(sess, stream, c.Conf.SocketQueueSize, c.Conf.MaxMessageSize, c.Conf.ReadDeadline)
	agent := c.NewAgent(quicConn, c.Conf.UID, c.Conf.Token)
	agent.Run()

	// cleanup
	quicConn.Close()
	c.Lock()
	delete(c.conns, sess)
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
		conn.Context().Done()
	}
	c.conns = nil
	c.Unlock()
	c.wg.Wait()
}
