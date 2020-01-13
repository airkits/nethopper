package grpc

import (
	"errors"
	"net"
	"sync"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/server"
	"github.com/gorilla/websocket"
)

//Config websocket conn config
type Config struct {
	Address        string
	MaxConnNum     int
	RWQueueSize    int
	MaxMessageSize uint32
	HTTPTimeout    uint32
	CertFile       string
	KeyFile        string
}

//ConnSet websocket conn set
type ConnSet map[*websocket.Conn]struct{}

//Conn websocket conn define
type Conn struct {
	sync.Mutex
	conn           *websocket.Conn
	writeChan      chan []byte
	maxMessageSize uint32
	closeFlag      bool
}

//NewConn create websocket conn
func NewConn(conn *websocket.Conn, rwQueueSize int, maxMessageSize uint32) network.Conn {
	wsConn := new(Conn)
	wsConn.conn = conn
	wsConn.writeChan = make(chan []byte, rwQueueSize)
	wsConn.maxMessageSize = maxMessageSize

	go func() {
		defer func() {
			if err := recover(); err != nil {
				server.PrintStack(false)
			}
		}()
		for b := range wsConn.writeChan {
			if b == nil {
				break
			}

			err := conn.WriteMessage(websocket.BinaryMessage, b)
			if err != nil {
				break
			}
		}

		conn.Close()
		wsConn.Lock()
		wsConn.closeFlag = true
		wsConn.Unlock()
	}()

	return wsConn
}

func (c *Conn) doDestroy() {
	c.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	c.conn.Close()

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

func (c *Conn) doWrite(b []byte) {
	if len(c.writeChan) == cap(c.writeChan) {
		server.Debug("close conn: channel full")
		c.doDestroy()
		return
	}

	c.writeChan <- b
}

//LocalAddr get local addr
func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

//RemoteAddr get remote addr
func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

//ReadMessage goroutine not safe
func (c *Conn) ReadMessage() ([]byte, error) {
	_, b, err := c.conn.ReadMessage()
	return b, err
}

//WriteMessage args must not be modified by the others goroutines
func (c *Conn) WriteMessage(args ...[]byte) error {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return nil
	}

	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// check len
	if msgLen > c.maxMessageSize {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}

	// don't copy
	if len(args) == 1 {
		c.doWrite(args[0])
		return nil
	}

	// merge the args
	msg := make([]byte, msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}

	c.doWrite(msg)

	return nil
}
