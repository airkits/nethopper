package grpc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/network/transport/pb/ss"
	"github.com/gonethopper/nethopper/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewClient create grpc client
func NewClient(m map[string]interface{}, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	c.GClientInfo = make([]*common.ClientInfo, 0)
	if err := c.ReadConfig(m); err != nil {
		panic(err)
	}
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc

	return c
}

//Client grpc client
type Client struct {
	sync.Mutex
	GClientInfo      []*common.ClientInfo
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
}

// Run client start run
func (c *Client) Run() {
	c.init()
	for _, info := range c.GClientInfo {
		for i := 0; i < c.ConnNum; i++ {
			c.wg.Add(1)
			go c.connect(info.ServerID, info.Name, info.Address)
		}
	}
}

// ReadConfig config map
// m := map[string]interface{}{
//  "grpcServerID_0":0,
//  "grpcAddress_0":":14000",
//  "grpcName_0":"grpcclient_0",
//	"connNum":1,
//  "socketQueueSize":100,
//  "maxMessageSize":4096
//  "connectInterval":3,
//  "handshakeTimeout":10,
//  "token":"12345678",
// }
func (c *Client) ReadConfig(m map[string]interface{}) error {

	for i := 0; i < 16; i++ {

		if !server.HasConfigKey(m, fmt.Sprintf("grpcServerID_%d", i)) {
			break
		}
		info := new(common.ClientInfo)
		if err := server.ParseConfigValue(m, fmt.Sprintf("grpcServerID_%d", i), i, &info.ServerID); err != nil {
			return err
		}
		if err := server.ParseConfigValue(m, fmt.Sprintf("grpcAddress_%d", i), "127.0.0.1:14000", &info.Address); err != nil {
			return err
		}
		if err := server.ParseConfigValue(m, fmt.Sprintf("grpcName_%d", i), fmt.Sprintf("grpcName_%d", i), &info.Name); err != nil {
			return err
		}
		c.GClientInfo = append(c.GClientInfo, info)
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

	if err := server.ParseConfigValue(m, "autoReconnect", true, &c.AutoReconnect); err != nil {
		return err
	}
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

func (c *Client) connect(serverID int, name string, address string) {
	defer c.wg.Done()

reconnect:
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		server.Fatal("grpc client connect to id:[%d] %s %s failed, reason: %v", serverID, name, address, err)
		if c.AutoReconnect {
			time.Sleep(c.ConnectInterval)
			server.Warning("grpc client try reconnect to id:[%d] %s %s", serverID, name, address)
			goto reconnect
		}
	}

	client := ss.NewRPCClient(conn)
	md := metadata.New(map[string]string{"token": name, "UID": strconv.Itoa(serverID)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ctx, cancel := context.WithCancel(ctx) // context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := client.Transport(ctx)
	if err != nil {
		server.Info("grpc client connect to id:[%d] %s %s transport failed, reason %v", serverID, name, address, err.Error())
		if c.AutoReconnect {
			time.Sleep(c.ConnectInterval)
			server.Warning("grpc client try reconnect to id:[%d] %s %s", serverID, name, address)
			goto reconnect
		}
	}
	server.Info("grpc client create new connection to id:[%d] %s %s.", serverID, name, address)
	c.Lock()
	c.conns[stream] = struct{}{}
	c.Unlock()

	grpcConn := NewConn(stream, c.RWQueueSize, c.MaxMessageSize)
	agent := c.NewAgent(grpcConn, uint64(serverID), name)

	agent.Run()

	// cleanup
	cancel()
	grpcConn.Close()
	c.Lock()
	delete(c.conns, stream)
	c.Unlock()
	c.CloseAgent(agent)
	agent.OnClose()

	if c.AutoReconnect {
		time.Sleep(c.ConnectInterval)
		server.Warning("grpc client try reconnect to id:[%d] %s %s", serverID, name, address)
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
