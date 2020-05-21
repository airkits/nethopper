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
func NewClient(conf server.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	c.Conf = conf.(*ClientConfig)
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

func (c *Client) dial(serverID int, address string) (quic.Session, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(address, tlsConf, nil)
	if err == nil {
		return session, nil
	}
	return nil, err

}

func (c *Client) connect(serverID int, name string, address string) {
	defer c.wg.Done()

reconnect:
	sess, err := c.dial(serverID, address)
	if err != nil {
		server.Fatal("quic client connect to id:[%d] %s %s failed, reason: %v", serverID, name, address, err)
		if c.Conf.AutoReconnect {
			time.Sleep(c.Conf.ConnectInterval * time.Second)
			server.Warning("quic client try reconnect to id:[%d] %s %s", serverID, name, address)
			goto reconnect
		}
	}

	stream, err := sess.OpenStreamSync(context.Background())
	if err != nil {
		server.Info("quic client connect to id:[%d] %s %s transport failed, reason %v", serverID, name, address, err.Error())
		if c.Conf.AutoReconnect {
			time.Sleep(c.Conf.ConnectInterval * time.Second)
			server.Warning("quic client try reconnect to id:[%d] %s %s", serverID, name, address)
			goto reconnect
		}
	}
	server.Info("quic client create new connection to id:[%d] %s %s.", serverID, name, address)
	c.Lock()
	c.conns[sess] = struct{}{}
	c.Unlock()

	quicConn := NewConn(sess, stream, c.Conf.SocketQueueSize, c.Conf.MaxMessageSize, c.Conf.ReadDeadline*time.Second)
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
		time.Sleep(c.Conf.ConnectInterval * time.Second)
		server.Warning("quic client try reconnect to id:[%d] %s %s", serverID, name, address)
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
