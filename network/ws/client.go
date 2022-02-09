package ws

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/network/common"
	"github.com/gorilla/websocket"
)

// NewClient create websocket client
func NewClient(conf config.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	c.Conf = conf.(*ClientConfig)
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc

	return c
}

//Client websocket client
type Client struct {
	sync.Mutex
	Conf       *ClientConfig
	NewAgent   network.AgentCreateFunc
	CloseAgent network.AgentCloseFunc
	dialer     websocket.Dialer
	conns      ConnSet
	wg         sync.WaitGroup
	closeFlag  bool
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
	c.closeFlag = false
	c.dialer = websocket.Dialer{
		HandshakeTimeout: c.Conf.HandshakeTimeout * time.Second,
	}

}

func (c *Client) dial(uid uint32, address string) (*websocket.Conn, error) {
	headers := make(http.Header)
	headers.Set(common.HeaderToken, c.Conf.Token)
	headers.Set(common.HeaderUID, fmt.Sprintf("%d", uid))

	conn, _, err := c.dialer.Dial(address, headers)
	if err == nil || c.closeFlag {
		return conn, nil
	}
	return nil, err
}

func (c *Client) connect(serverID int, name string, address string) {
	defer c.wg.Done()

reconnect:
	conn, err := c.dial(c.Conf.UID, address)
	if err != nil {
		log.Fatal("websocket client connect to id:[%d] %s %s failed, reason: %v", serverID, name, address, err)
		if c.Conf.AutoReconnect {
			time.Sleep(c.Conf.ConnectInterval * time.Second)
			log.Warning("websocket client try reconnect to id:[%d] %s %s", serverID, name, address)
			goto reconnect
		}
	}
	conn.SetReadLimit(int64(c.Conf.MaxMessageSize))

	c.Lock()
	if c.closeFlag {
		c.Unlock()
		conn.Close()
		return
	}
	c.conns[conn] = struct{}{}
	c.Unlock()

	wsConn := NewConn(conn, c.Conf.SocketQueueSize, c.Conf.MaxMessageSize)
	agent := c.NewAgent(wsConn, uint64(c.Conf.UID), c.Conf.Token)
	agent.Run()

	// cleanup
	wsConn.Close()
	c.Lock()
	delete(c.conns, conn)
	c.Unlock()
	c.CloseAgent(agent)
	agent.OnClose()

	if c.Conf.AutoReconnect {
		time.Sleep(c.Conf.ConnectInterval * time.Second)
		log.Warning("websocket client try reconnect to id:[%d] %s %s", serverID, name, address)
		goto reconnect
	}
}

// Close client connections
func (c *Client) Close() {
	c.Lock()
	c.closeFlag = true
	for conn := range c.conns {
		conn.Close()
	}
	c.conns = nil
	c.Unlock()
	c.wg.Wait()
}
