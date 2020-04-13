package tcp

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
	"google.golang.org/grpc"
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

//Client websocket client
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

func (c *Client) connect() {
	defer c.wg.Done()

reconnect:
	conn, err := grpc.Dial(c.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		server.Fatal("did not connect: %v", err)
	}

	client := ss.NewRPCClient(conn)

	ctx, cancel := context.WithCancel(context.Background()) // context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := client.Transport(ctx)
	if err != nil {
		server.Info("transport %v", err.Error())
	}

	c.Lock()
	c.conns[stream] = struct{}{}
	c.Unlock()

	grpcConn := NewConn(stream, c.RWQueueSize, c.MaxMessageSize)
	agent := c.NewAgent(grpcConn)
	agent.Run()

	// cleanup
	cancel()
	grpcConn.Close()
	c.Lock()
	delete(c.conns, stream)
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
