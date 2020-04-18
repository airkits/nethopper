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
func NewClient(m map[string]interface{}, agentFunc network.AgentCreateFunc) *Client {
	c := new(Client)
	if err := c.ReadConfig(m); err != nil {
		panic(err)
	}
	c.NewAgent = agentFunc
	return c
}

//Client kcp client
type Client struct {
	sync.Mutex
	Address             string
	ConnNum             int
	ConnectInterval     time.Duration
	RWQueueSize         int
	MaxMessageSize      uint32
	HandshakeTimeout    time.Duration
	AutoReconnect       bool
	NewAgent            network.AgentCreateFunc
	conns               ConnSet
	wg                  sync.WaitGroup
	Token               string
	UDPSocketBufferSize int //UDP listener socket buffer
	dscp                int //set DSCP(6bit)
	sndwnd              int //per connection UDP send window
	rcvwnd              int //per connection UDP recv window
	mtu                 int //MTU of UDP packets, without IP(20) + UDP(8)
	nodelay             int //ikcp_nodelay()
	interval            int //ikcp_nodelay()
	resend              int //ikcp_nodelay()
	nc                  int //ikcp_nodelay()
	ReadDeadline        time.Duration
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
//  "address":":8888",
//	"connNum":1,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
//  "connectInterval":3,
//  "handshakeTimeout":10,
//  "token":"12345678",
// }
func (c *Client) ReadConfig(m map[string]interface{}) error {

	if err := server.ParseConfigValue(m, "address", "127.0.0.1:8888", &c.Address); err != nil {
		return err
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

	if err := server.ParseConfigValue(m, "udpSocketBuf", 4194304, &c.UDPSocketBufferSize); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "udpSndWnd", 32, &c.sndwnd); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "udpRcvWnd", 32, &c.rcvwnd); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "udpMtu", 1280, &c.mtu); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "dscp", 46, &c.dscp); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "nodelay", 1, &c.nodelay); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "interval", 20, &c.interval); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "resend", 1, &c.resend); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "nc", 1, &c.nc); err != nil {
		return err
	}
	if err := server.ParseConfigValue(m, "readDeadline", 15, &c.ReadDeadline); err != nil {
		return err
	}
	c.ReadDeadline = c.ReadDeadline * time.Second

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

}

func (c *Client) dial() *kcp.UDPSession {
	for {
		conn, err := kcp.DialWithOptions(c.Address, nil, 0, 0)
		if err == nil {
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
	conn.SetDSCP(c.dscp)
	conn.SetWindowSize(c.sndwnd, c.rcvwnd)
	conn.SetNoDelay(c.nodelay, c.interval, c.resend, c.nc)
	conn.SetStreamMode(true)
	conn.SetMtu(c.mtu)

	c.Lock()
	c.conns[conn] = struct{}{}
	c.Unlock()
	kcpConn := NewConn(conn, c.RWQueueSize, c.MaxMessageSize, c.ReadDeadline)
	agent := c.NewAgent(kcpConn)
	agent.Run()

	// cleanup
	kcpConn.Close()
	c.Lock()
	delete(c.conns, conn)
	c.Unlock()
	agent.OnClose()

	if c.AutoReconnect {
		time.Sleep(c.ConnectInterval)
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
