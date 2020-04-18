package kcp

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gonethopper/nethopper/network"
	"github.com/gonethopper/nethopper/network/common"
	"github.com/gonethopper/nethopper/server"
)

//ErrReadPackageLength read package length failed
var ErrReadPackageLength = errors.New("read package length failed")

//ErrReadMessage read message failed
var ErrReadMessage = errors.New("read message failed")

const (
	//PackageLengthSize package length size
	PackageLengthSize = 2
)

//Config kcp conn config
type Config struct {
	Address             string
	MaxConnNum          int
	RWQueueSize         int
	MaxMessageSize      uint32
	UDPSocketBufferSize int //UDP listener socket buffer
	ReadDeadline        time.Duration
	dscp                int //set DSCP(6bit)
	sndwnd              int //per connection UDP send window
	rcvwnd              int //per connection UDP recv window
	mtu                 int //MTU of UDP packets, without IP(20) + UDP(8)
	nodelay             int //ikcp_nodelay()
	interval            int //ikcp_nodelay()
	resend              int //ikcp_nodelay()
	nc                  int //ikcp_nodelay()
}

//ConnSet kcp conn set
type ConnSet map[net.Conn]struct{}

//Conn k c p conn define
type Conn struct {
	sync.Mutex
	conn           net.Conn
	writeChan      chan []byte
	maxMessageSize uint32
	readDeadline   time.Duration
	closeFlag      bool
}

//NewConn create kcp conn
func NewConn(conn net.Conn, rwQueueSize int, maxMessageSize uint32, readDeadline time.Duration) network.IConn {
	kcpConn := new(Conn)
	kcpConn.conn = conn
	kcpConn.writeChan = make(chan []byte, rwQueueSize)
	kcpConn.maxMessageSize = maxMessageSize
	kcpConn.readDeadline = readDeadline
	go func() {
		defer func() {
			if err := recover(); err != nil {
				server.PrintStack(false)
			}
		}()
		for b := range kcpConn.writeChan {
			if b == nil {
				break
			}
			// write data
			n, err := conn.Write(b)
			if err != nil {
				server.Warning("Error send reply data, bytes: %v reason: %v", n, err)
				break
			}

		}

		conn.Close()
		kcpConn.Lock()
		kcpConn.closeFlag = true
		kcpConn.Unlock()
	}()

	return kcpConn
}

func (c *Conn) doDestroy() {
	c.conn.Close()
	if !c.closeFlag {
		close(c.writeChan)
		c.closeFlag = true
	}
}

//Destroy kcp conn destory
func (c *Conn) Destroy() {
	c.Lock()
	defer c.Unlock()

	c.doDestroy()
}

//Close kcp conn close
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
func (c *Conn) ReadMessage() (interface{}, error) {
	// for reading the 2-Byte Package Length
	pkgLen := make([]byte, common.PackageLengthSize)
	// read loop
	for {
		// solve dead link problem:
		// physical disconnection without any communcation between client and server
		// will cause the read to block FOREVER, so a timeout is a rescue.
		c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))

		// read 2B Package Length
		n, err := io.ReadFull(c.conn, pkgLen)
		if err != nil {
			server.Warning("read package length failed, ip:%v reason:%v size:%v", c.RemoteAddr().String(), err, n)
			return nil, ErrReadPackageLength
		}
		size := binary.BigEndian.Uint16(pkgLen)

		// alloc a byte slice of the size defined in the package length for reading data
		payload := make([]byte, size)
		n, err = io.ReadFull(c.conn, payload)
		if err != nil {
			server.Warning("read payload failed, ip:%v reason:%v size:%v", c.RemoteAddr().String(), err, n)
			return nil, ErrReadMessage
		}

		return payload, nil
	}
}

//WriteMessage args must not be modified by the others goroutines
//buffer must packet with 2B package length
func (c *Conn) WriteMessage(args ...interface{}) error {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return nil
	}

	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i].([]byte)))
	}

	// check len
	if msgLen > c.maxMessageSize {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}

	// don't copy
	if len(args) == 1 {
		c.doWrite(args[0].([]byte))
		return nil
	}

	// merge the args
	msg := make([]byte, msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i].([]byte))
		l += len(args[i].([]byte))
	}

	c.doWrite(msg)

	return nil
}
