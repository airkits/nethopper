package connect

import "net"

// Connect define network interface
type Connect interface {
	// Setup init Connect with config
	Setup(m map[string]interface{}) (Connect, error)
	// Listen on local port
	Listen()
	// Accept accepts the next incoming call and returns the new connection.
	Accept() (net.Conn, error)
}
