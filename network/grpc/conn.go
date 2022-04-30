package grpc

import (
	"context"
	"net"
	"sync"

	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/proto/ss"
	"google.golang.org/grpc/peer"
)

//ConnSet grpc conn set
type ConnSet map[IRPCStream]struct{}

//IRPCStream define rpc stream interface
type IRPCStream interface {
	Send(*ss.Message) error
	Recv() (*ss.Message, error)
	// Context returns the context for this stream.
	Context() context.Context
}

//Conn grpc conn define
type Conn struct {
	sync.Mutex
	stream         IRPCStream
	writeChan      chan *ss.Message
	maxMessageSize uint32
	closeFlag      bool
}

//NewConn create websocket conn
func NewConn(stream IRPCStream, rwQueueSize int, maxMessageSize uint32) network.IConn {
	grpcConn := new(Conn)
	grpcConn.stream = stream
	grpcConn.writeChan = make(chan *ss.Message, rwQueueSize)
	grpcConn.maxMessageSize = maxMessageSize

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.PrintStack(false)
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

//Destroy grpc conn destory
func (c *Conn) Destroy() {
	c.Lock()
	defer c.Unlock()

	c.doDestroy()
}

//Close grpc conn close
func (c *Conn) Close() {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return
	}

	c.doWrite(nil)
	c.closeFlag = true
}

func (c *Conn) doWrite(b *ss.Message) {
	if len(c.writeChan) == cap(c.writeChan) {
		log.Error("close conn: channel full")
		c.doDestroy()
		return
	}

	c.writeChan <- b
}

//LocalAddr get local addr
func (c *Conn) LocalAddr() net.Addr {
	log.Error("[LocalAddr] invoke LocalAddr() unsupport")
	return nil
}

//RemoteAddr get remote addr
func (c *Conn) RemoteAddr() net.Addr {
	pr, ok := peer.FromContext(c.stream.Context())
	if !ok {
		log.Error("[RemoteAddr] invoke FromContext() failed")
		return nil
	}
	if pr.Addr == net.Addr(nil) {
		log.Error("[RemoteAddr] peer.Addr is nil")
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
		c.doWrite(args[i].(*ss.Message))
	}
	return nil
}
