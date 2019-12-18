package network

import "net"

// Conn define network conn interface
type Conn interface {
	//ReadMessage read message from conn
	ReadMessage() ([]byte, error)
	//WriteMessage write message to conn
	WriteMessage(args ...[]byte) error
	//LocalAddr get local addr
	LocalAddr() net.Addr
	//RemoteAddr get remote addr
	RemoteAddr() net.Addr
	//Close conn
	Close()
	//Destory conn
	Destroy()
}

//AgentCreateFunc create agent func
type AgentCreateFunc func(Conn) IAgent
