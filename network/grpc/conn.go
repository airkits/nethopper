package grpc

import (
	"net"
	"sync"

	"github.com/gonethopper/nethopper/examples/model/pb/ss"
	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"google.golang.org/grpc/peer"
)

//ConnSet websocket conn set
type ConnSet map[ss.RPC_TransportServer]struct{}

//Conn websocket conn define
type Conn struct {
	sync.Mutex
	stream         ss.RPC_TransportServer
	writeChan      chan *ss.SSMessage
	maxMessageSize uint32
	closeFlag      bool
}

//NewConn create websocket conn
func NewConn(stream ss.RPC_TransportServer, rwQueueSize int, maxMessageSize uint32) network.Conn {
	grpcConn := new(Conn)
	grpcConn.stream = stream
	grpcConn.writeChan = make(chan *ss.SSMessage, rwQueueSize)
	grpcConn.maxMessageSize = maxMessageSize

	go func() {
		defer func() {
			if err := recover(); err != nil {
				server.PrintStack(false)
			}
		}()
		for b := range grpcConn.writeChan {
			if b == nil {
				break
			}

			err := stream.Send(b)
			if err != nil {
				break
			}
		}

		//	conn.Close()
		grpcConn.Lock()
		grpcConn.closeFlag = true
		grpcConn.Unlock()
	}()

	return grpcConn
}

func (c *Conn) doDestroy() {
	// c.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	// c.conn.Close()

	if !c.closeFlag {
		close(c.writeChan)
		c.closeFlag = true
	}
}

//Destroy websocket conn destory
func (c *Conn) Destroy() {
	c.Lock()
	defer c.Unlock()

	c.doDestroy()
}

//Close websocket conn close
func (c *Conn) Close() {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return
	}

	c.doWrite(nil)
	c.closeFlag = true
}

func (c *Conn) doWrite(b *ss.SSMessage) {
	if len(c.writeChan) == cap(c.writeChan) {
		server.Debug("close conn: channel full")
		c.doDestroy()
		return
	}

	c.writeChan <- b
}

//LocalAddr get local addr
func (c *Conn) LocalAddr() net.Addr {
	server.Error("[LocalAddr] invoke LocalAddr() unsupport")
	return nil
}

//RemoteAddr get remote addr
func (c *Conn) RemoteAddr() net.Addr {
	pr, ok := peer.FromContext(c.stream.Context())
	if !ok {
		server.Error("[RemoteAddr] invoke FromContext() failed")
		return nil
	}
	if pr.Addr == net.Addr(nil) {
		server.Error("[RemoteAddr] peer.Addr is nil")
		return nil
	}

	return pr.Addr
}

//ReadMessage goroutine not safe
func (c *Conn) ReadMessage() (interface{}, error) {
	return c.stream.Recv()
}

//WriteMessage args must not be modified by the others goroutines
func (c *Conn) WriteMessage(args ...interface{}) error {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return nil
	}

	for i := 0; i < len(args); i++ {
		c.doWrite(args[i].(*ss.SSMessage))
	}
	return nil
}
