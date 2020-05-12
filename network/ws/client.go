package ws

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/server"
	"github.com/gorilla/websocket"
)

// NewClient create websocket client
func NewClient(m map[string]interface{}, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)

	c.WSClientInfo = make([]*common.ClientInfo, 0)
	if err := c.ReadConfig(m); err != nil {
		panic(err)
	}
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc

	return c
}

//Client websocket client
type Client struct {
	sync.Mutex
	WSClientInfo     []*common.ClientInfo
	ConnNum          int
	ConnectInterval  time.Duration
	RWQueueSize      int
	MaxMessageSize   uint32
	HandshakeTimeout time.Duration
	AutoReconnect    bool
	NewAgent         network.AgentCreateFunc
	CloseAgent       network.AgentCloseFunc
	dialer           websocket.Dialer
	conns            ConnSet
	wg               sync.WaitGroup
	closeFlag        bool

	Token string
	UID   uint64
}

// Run client start run
func (c *Client) Run() {
	c.init()
	for _, info := range c.WSClientInfo {
		for i := 0; i < c.ConnNum; i++ {
			c.wg.Add(1)
			go c.connect(info.ServerID, info.Name, info.Address)
		}
	}
}

// ReadConfig config map
// m := map[string]interface{}{
//  "address":":12080",
//	"connNum":1,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
//  "connectInterval":3,
//  "handshakeTimeout":10,
//  "token":"12345678",
// }
func (c *Client) ReadConfig(m map[string]interface{}) error {

	for i := 0; i < 16; i++ {

		if !server.HasConfigKey(m, fmt.Sprintf("wsServerID_%d", i)) {
			break
		}
		info := new(common.ClientInfo)
		if err := server.ParseConfigValue(m, fmt.Sprintf("wsServerID_%d", i), i, &info.ServerID); err != nil {
			return err
		}
		if err := server.ParseConfigValue(m, fmt.Sprintf("wsAddress_%d", i), "ws://127.0.0.1:12080", &info.Address); err != nil {
			return err
		}
		if err := server.ParseConfigValue(m, fmt.Sprintf("wsName_%d", i), fmt.Sprintf("grpcName_%d", i), &info.Name); err != nil {
			return err
		}
		c.WSClientInfo = append(c.WSClientInfo, info)
	}

	if err := server.ParseConfigValue(m, "connNum", 1, &c.ConnNum); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "socketQueueSize", 100, &c.RWQueueSize); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "maxMessageSize", 4096, &c.MaxMessageSize); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "connectInterval", 3, &c.ConnectInterval); err != nil {
		return err
	}
	c.ConnectInterval = c.ConnectInterval * time.Second
	if err := server.ParseConfigValue(m, "handshakeTimeout", 10, &c.HandshakeTimeout); err != nil {
		return err
	}
	c.HandshakeTimeout = c.HandshakeTimeout * time.Second

	if err := server.ParseConfigValue(m, "token", "12345678", &c.Token); err != nil {
		return err
	}

	return nil
}

func (c *Client) init() {
	c.Lock()
	defer c.Unlock()

	if c.NewAgent == nil {
		server.Fatal("NewAgent must not be nil")
	}
	if c.conns != nil {
		server.Fatal("client is running")
	}

	c.conns = make(ConnSet)
	c.closeFlag = false
	c.dialer = websocket.Dialer{
		HandshakeTimeout: c.HandshakeTimeout,
	}

}

func (c *Client) dial(serverID int, address string) *websocket.Conn {
	headers := make(http.Header)
	headers.Set(common.HeaderToken, c.Token)
	headers.Set(common.HeaderUID, fmt.Sprintf("%d", serverID))
	for {
		conn, _, err := c.dialer.Dial(address, headers)
		if err == nil || c.closeFlag {
			return conn
		}

		server.Warning("connect to %v error: %v", address, err)
		time.Sleep(c.ConnectInterval)
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
	conn.SetReadLimit(int64(c.MaxMessageSize))

	c.Lock()
	if c.closeFlag {
		c.Unlock()
		conn.Close()
		return
	}
	c.conns[conn] = struct{}{}
	c.Unlock()

	wsConn := NewConn(conn, c.RWQueueSize, c.MaxMessageSize)
	agent := c.NewAgent(wsConn, uint64(serverID), c.Token)
	agent.Run()

	// cleanup
	wsConn.Close()
	c.Lock()
	delete(c.conns, conn)
	c.Unlock()
	c.CloseAgent(agent)
	agent.OnClose()

	if c.AutoReconnect {
		time.Sleep(c.ConnectInterval)
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
