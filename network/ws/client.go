package ws

import (
	"net/http"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"github.com/gorilla/websocket"
	"github.com/prometheus/common/log"
)

// NewClient create websocket client
func NewClient(m map[string]interface{}, agentFunc network.AgentCreateFunc) *Client {
	c := new(Client)
	if err := c.ReadConfig(m); err != nil {
		panic(err)
	}
	c.NewAgent = agentFunc
	return c
}

//Client websocket client
type Client struct {
	sync.Mutex
	Address          string
	ConnNum          int
	ConnectInterval  time.Duration
	RWQueueSize      int
	MaxMessageSize   uint32
	HandshakeTimeout time.Duration
	AutoReconnect    bool
	NewAgent         network.AgentCreateFunc
	dialer           websocket.Dialer
	conns            ConnSet
	wg               sync.WaitGroup
	closeFlag        bool
	headers          http.Header
	Token            string
}

// Run client start run
func (c *Client) Run() {
	c.init()
	for i := 0; i < c.ConnNum; i++ {
		c.wg.Add(1)
		go c.connect()
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

	address, err := server.ParseValue(m, "address", "ws://127.0.0.1:12080")
	if err != nil {
		return err
	}
	c.Address = address.(string)

	connNum, err := server.ParseValue(m, "connNum", 1)
	if err != nil {
		return err
	}
	c.ConnNum = connNum.(int)

	rwQueueSize, err := server.ParseValue(m, "socketQueueSize", 100)
	if err != nil {
		return err
	}
	c.RWQueueSize = rwQueueSize.(int)

	maxMessageSize, err := server.ParseValue(m, "maxMessageSize", 4096)
	if err != nil {
		return err
	}
	c.MaxMessageSize = uint32(maxMessageSize.(int))

	connectInterval, err := server.ParseValue(m, "connectInterval", 3)
	if err != nil {
		return err
	}
	c.ConnectInterval = time.Duration(connectInterval.(int)) * time.Second

	timeout, err := server.ParseValue(m, "handshakeTimeout", 10)
	if err != nil {
		return err
	}
	c.HandshakeTimeout = time.Duration(timeout.(int)) * time.Second

	token, err := server.ParseValue(m, "token", "12345678")
	if err != nil {
		return err
	}
	c.Token = token.(string)
	c.headers = make(http.Header)
	c.headers.Set("token", c.Token)

	return nil
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
		HandshakeTimeout: c.HandshakeTimeout,
	}

}

func (c *Client) dial() *websocket.Conn {
	for {
		conn, _, err := c.dialer.Dial(c.Address, c.headers)
		if err == nil || c.closeFlag {
			return conn
		}

		server.Warning("connect to %v error: %v", c.Address, err)
		time.Sleep(c.ConnectInterval)
		continue
	}
}

func (c *Client) connect() {
	defer c.wg.Done()

reconnect:
	conn := c.dial()
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
	agent := c.NewAgent(wsConn)
	agent.Run()

	// cleanup
	wsConn.Close()
	c.Lock()
	delete(c.conns, conn)
	c.Unlock()
	agent.OnClose()

	if c.AutoReconnect {
		time.Sleep(c.ConnectInterval)
		goto reconnect
	}
}

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
