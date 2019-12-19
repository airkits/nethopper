package udp

import (
	"net"

	"github.com/gonethopper/nethopper/server"
	"github.com/xtaci/kcp-go"
)

// KCPConnect listen port and accept tcp connection
type KCPConnect struct {
	Address    string
	BufferSize int
	Sndwnd     int
	Rcvwnd     int
	Mtu        int
	Dscp       int
	Nodelay    int
	Interval   int
	Resend     int
	Nc         int
	listener   *kcp.Listener
}

// Setup init Connect with config
func (c *KCPConnect) Setup(m map[string]interface{}) (*KCPConnect, error) {
	if err := c.ReadConfig(m); err != nil {
		panic(err)
	}
	return c, nil
}

// config map
// address default :8888
// bufferSize default 4194304, Usage: "UDP listener socket buffer",
// sndwnd default 32, Usage: "per connection UDP send window",
// rcvwnd defualt 32, Usage: "per connection UDP recv window",
// mtu default 1280, Usage: "MTU of UDP packets, without IP(20) + UDP(8)",
// dscp default 46, Usage: "set DSCP(6bit)",
// nodelay default 1,Usage: "ikcp_nodelay()",
// interval default 20, Usage: "ikcp_nodelay()",
// resend default 1, Usage: "ikcp_nodelay()",
// nc default 1,Usage: "ikcp_nodelay()"
func (c *KCPConnect) ReadConfig(m map[string]interface{}) error {
	if err := server.ParseConfigValue(m, "address", ":8888", &c.Address); err != nil {
		return err
	}

	if err := server.ParseConfigValue(m, "bufferSize", 4194304, &c.BufferSize); err != nil {
		return err
	}

	// sndwnd default 32, Usage: "per connection UDP send window",
	if err := server.ParseConfigValue(m, "sndwnd", 32, &c.Sndwnd); err != nil {
		return err
	}

	// rcvwnd defualt 32, Usage: "per connection UDP recv window",
	if err := server.ParseConfigValue(m, "rcvwnd", 32, &c.Rcvwnd); err != nil {
		return err
	}

	// mtu default 1280, Usage: "MTU of UDP packets, without IP(20) + UDP(8)",
	if err := server.ParseConfigValue(m, "mtu", 1280, &c.Mtu); err != nil {
		return err
	}

	// dscp default 46, Usage: "set DSCP(6bit)",
	if err := server.ParseConfigValue(m, "dscp", 46, &c.Dscp); err != nil {
		return err
	}

	// nodelay default 1,Usage: "ikcp_nodelay()",
	if err := server.ParseConfigValue(m, "nodelay", 1, &c.Nodelay); err != nil {
		return err
	}

	// interval default 20, Usage: "ikcp_nodelay()",
	if err := server.ParseConfigValue(m, "interval", 20, &c.Interval); err != nil {
		return err
	}

	// resend default 1, Usage: "ikcp_nodelay()",
	if err := server.ParseConfigValue(m, "resend", 1, &c.Resend); err != nil {
		return err
	}

	// nc default 1,Usage: "ikcp_nodelay()"
	if err := server.ParseConfigValue(m, "nc", 1, &c.Nc); err != nil {
		return err
	}

	return nil
}

// Listen and bind local ip
func (c *KCPConnect) Listen() {
	l, err := kcp.Listen(c.Address)
	if err != nil {
		panic(err)
	}
	server.Info("udp listening on: %s", l.Addr())
	listener := l.(*kcp.Listener)
	c.listener = listener
	if err := listener.SetReadBuffer(c.BufferSize); err != nil {
		server.Error("SetReadBuffer %s", err.Error())
	}
	if err := listener.SetWriteBuffer(c.BufferSize); err != nil {
		server.Error("SetWriteBuffer %s", err.Error())
	}
	if err := listener.SetDSCP(c.Dscp); err != nil {
		server.Error("SetDSCP %s", err.Error())
	}

}

// Accept accepts the next incoming call and returns the new connection.
func (c *KCPConnect) Accept() (net.Conn, error) {
	conn, err := c.listener.AcceptKCP()
	if err != nil {
		server.Warning("accept failed: %s", err.Error())
		return nil, err
	}
	// set kcp parameters
	conn.SetWindowSize(c.Sndwnd, c.Rcvwnd)
	conn.SetNoDelay(c.Nodelay, c.Interval, c.Resend, c.Nc)
	conn.SetStreamMode(true)
	conn.SetMtu(c.Mtu)
	return conn, nil
}
