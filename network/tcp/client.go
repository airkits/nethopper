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
func NewClient(m map[string]interface{}, agentFunc network.AgentCreateFunc) *Client {
	c := new(Client)
	if err := c.ReadConfig(m); err != nil {
		panic(err)
	}
	c.NewAgent = agentFunc
	return c
}

//Client tcp client
type Client struct {
	sync.Mutex
	Address          string
	Network          string
	ConnNum          int
	ConnectInterval  time.Duration
	RWQueueSize      int
	MaxMessageSize   uint32
	HandshakeTimeout time.Duration
	AutoReconnect    bool
	NewAgent         network.AgentCreateFunc
	conns            ConnSet
	wg               sync.WaitGroup
	Token            string
	ReadBufferSize   int
	WriteBufferSize  int
	ReadDeadline     time.Duration
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
//  "address":":15000",
//	"connNum":1,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
//  "connectInterval":3,
//  "handshakeTimeout":10,
//  "token":"12345678",
// }
func (c *Client) ReadConfig(m map[string]interface{}) error {

	if err := server.ParseConfigValue(m, "address", "127.0.0.1:15000", &c.Address); err != nil {
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

	if err := server.ParseConfigValue(m, "readBufferSize", 32767, &c.ReadBufferSize); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "writeBufferSize", 32767, &c.WriteBufferSize); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "network", "tcp4", &c.Network); err != nil {
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

func (c *Client) dial() net.Conn {
	for {
		//conn, err := net.DialTimeout(c.Network, c.Address, time.Second*30)
		conn, err := net.Dial(c.Network, c.Address)
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
	c.Lock()
	c.conns[conn] = struct{}{}
	c.Unlock()

	tcpConn := NewConn(conn, c.RWQueueSize, c.MaxMessageSize, c.ReadDeadline)
	agent := c.NewAgent(tcpConn)
	agent.Run()

	// cleanup
	tcpConn.Close()
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
