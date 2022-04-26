package grpc

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/proto/ss"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// NewClient create grpc client
func NewClient(conf config.IConfig, agentFunc network.AgentCreateFunc, agentCloseFunc network.AgentCloseFunc) *Client {
	c := new(Client)
	c.Conf = conf.(*ClientConfig)
	c.NewAgent = agentFunc
	c.CloseAgent = agentCloseFunc
	c.wg = &sync.WaitGroup{}
	return c
}

//Client grpc client
type Client struct {
	sync.Mutex
	Conf       *ClientConfig
	NewAgent   network.AgentCreateFunc
	CloseAgent network.AgentCloseFunc
	conns      ConnSet
	wg         *sync.WaitGroup
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
		log.Fatal("[GRPCClient] NewAgent must not be nil")
	}
	if c.conns != nil {
		log.Fatal("[GRPCClient] client is running")
	}

	c.conns = make(ConnSet)
}

func (c *Client) connect(serverID int, name string, address string) {
	defer c.wg.Done()

reconnect:
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatal("[GRPCClient] grpc client connect to id:[%d] %s %s failed, reason: %v", serverID, name, address, err)
		if c.Conf.AutoReconnect {
			time.Sleep(c.Conf.ConnectInterval * time.Second)
			log.Warning("[GRPCClient] grpc client try reconnect to id:[%d] %s %s", serverID, name, address)
			goto reconnect
		}
	}

	client := ss.NewRPCClient(conn)
	md := metadata.New(map[string]string{"token": name, "UID": strconv.Itoa(serverID)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ctx, cancel := context.WithCancel(ctx) // context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := client.Transport(ctx)
	if err != nil {
		log.Info("[GRPCClient] grpc client connect to id:[%d] %s %s transport failed, reason %v", serverID, name, address, err.Error())
		if c.Conf.AutoReconnect {
			time.Sleep(c.Conf.ConnectInterval * time.Second)
			log.Warning("[GRPCClient] grpc client try reconnect to id:[%d] %s %s", serverID, name, address)
			goto reconnect
		}
	}
	log.Info("[GRPCClient] grpc client create new connection to id:[%d] %s %s.", serverID, name, address)
	c.Lock()
	c.conns[stream] = struct{}{}
	c.Unlock()

	grpcConn := NewConn(stream, c.Conf.SocketQueueSize, c.Conf.MaxMessageSize)
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

	if c.Conf.AutoReconnect {
		time.Sleep(c.Conf.ConnectInterval * time.Second)
		log.Warning("[GRPCClient] grpc client try reconnect to id:[%d] %s %s", serverID, name, address)
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

//Call sync get response
func (c *Client) Call(serverID int, name string, address string) *ss.Message {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		log.Fatal("[GRPCClient] grpc client connect to id:[%d] %s %s failed, reason: %v", serverID, name, address, err)
		if c.Conf.AutoReconnect {
			time.Sleep(c.Conf.ConnectInterval * time.Second)
			log.Warning("[GRPCClient] grpc client try reconnect to id:[%d] %s %s", serverID, name, address)
			return nil
		}
	}

	client := ss.NewRPCClient(conn)
	md := metadata.New(map[string]string{"token": "token", "UID": strconv.Itoa(serverID)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	ctx, _ = context.WithCancel(ctx) // context.WithTimeout(context.Background(), 10*time.Second)

	r, err := client.Call(ctx, &ss.Message{})
	if err != nil {
		log.Fatal("[GRPCClient] could not greet: %v", err)
	}
	return r
}
