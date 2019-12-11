package ws

import (
	"net"

	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/queue"
	"github.com/gorilla/websocket"
)

type WebsocketConnSet map[*websocket.Conn]struct{}

type WSConnection struct {
	conn *websocket.Conn
	q    queue.Queue
}

func newWSConn(conn *websocket.Conn, queueSize int32) *WSConnection {
	wsConn := new(WSConnection)
	wsConn.conn = conn
	wsConn.q = queue.NewChanQueue(queueSize)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				server.PrintStack(false)
			}
		}()
		for {
			buf, err := wsConn.q.Pop()
			if err != nil {
				break
			}
			err = conn.WriteMessage(websocket.BinaryMessage, buf.([]byte))
			if err != nil {
				break
			}
		}

		conn.Close()
	}()

	return wsConn
}

func (wsConn *WSConnection) Close() {

	wsConn.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	wsConn.conn.Close()
	wsConn.q.Close()

}

func (wsConn *WSConnection) doWrite(b []byte) error {
	return wsConn.q.AsyncPush(b)
}

func (wsConn *WSConnection) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

func (wsConn *WSConnection) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

// goroutine not safe
func (wsConn *WSConnection) ReadMsg() ([]byte, error) {
	_, b, err := wsConn.conn.ReadMessage()
	return b, err
}

// args must not be modified by the others goroutines
func (wsConn *WSConnection) WriteMsg(args ...[]byte) error {

	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// don't copy
	if len(args) == 1 {
		wsConn.doWrite(args[0])
		return nil
	}

	// merge the args
	msg := make([]byte, msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}
	wsConn.doWrite(msg)

	return nil
}
