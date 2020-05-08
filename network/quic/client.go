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
func NewClient(m map[string]interface{}, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	if err := c.ReadConfig(m); err != nil {
		panic(err)
	}
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc

	return c
}

//Client quic client
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
	CloseAgent       network.AgentCloseFunc
	conns            ConnSet
	wg               sync.WaitGroup
	Token            string
	UID              uint64
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
//  "address":":16000",
//	"connNum":1,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
//  "connectInterval":3,
//  "handshakeTimeout":10,
//  "token":"12345678",
// }
func (c *Client) ReadConfig(m map[string]interface{}) error {

	if err := server.ParseConfigValue(m, "address", "127.0.0.1:16000", &c.Address); err != nil {
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

	if err := server.ParseConfigValue(m, "network", "quic4", &c.Network); err != nil {
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

func (c *Client) dial() quic.Session {
	for {
		tlsConf := &tls.Config{
			InsecureSkipVerify: true,
			NextProtos:         []string{"quic-echo-example"},
		}
		session, err := quic.DialAddr(c.Address, tlsConf, nil)
		if err == nil {
			return session
		}

		server.Warning("connect to %v error: %v", c.Address, err)
		time.Sleep(c.ConnectInterval)
		continue
	}
}

func (c *Client) connect() {
	defer c.wg.Done()

reconnect:
	sess := c.dial()
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

	quicConn := NewConn(sess, stream, c.RWQueueSize, c.MaxMessageSize, c.ReadDeadline)
	agent := c.NewAgent(quicConn, c.UID, c.Token)
	agent.Run()

	// cleanup
	quicConn.Close()
	c.Lock()
	delete(c.conns, sess)
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
	for conn := range c.conns {
		conn.Context().Done()
	}
	c.conns = nil
	c.Unlock()
	c.wg.Wait()
}
