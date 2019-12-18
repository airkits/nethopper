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
	address, err := server.ParseValue(m, "address", ":8888")
	if err != nil {
		return err
	}
	c.Address = address.(string)

	bufferSize, err := server.ParseValue(m, "bufferSize", 4194304)
	if err != nil {
		return err
	}
	c.BufferSize = bufferSize.(int)

	// sndwnd default 32, Usage: "per connection UDP send window",
	sndwnd, err := server.ParseValue(m, "sndwnd", 32)
	if err != nil {
		return err
	}
	c.Sndwnd = sndwnd.(int)
	// rcvwnd defualt 32, Usage: "per connection UDP recv window",
	rcvwnd, err := server.ParseValue(m, "rcvwnd", 32)
	if err != nil {
		return err
	}
	c.Rcvwnd = rcvwnd.(int)
	// mtu default 1280, Usage: "MTU of UDP packets, without IP(20) + UDP(8)",
	mtu, err := server.ParseValue(m, "mtu", 1280)
	if err != nil {
		return err
	}
	c.Mtu = mtu.(int)
	// dscp default 46, Usage: "set DSCP(6bit)",
	dscp, err := server.ParseValue(m, "dscp", 46)
	if err != nil {
		return err
	}
	c.Dscp = dscp.(int)
	// nodelay default 1,Usage: "ikcp_nodelay()",
	nodelay, err := server.ParseValue(m, "nodelay", 1)
	if err != nil {
		return err
	}
	c.Nodelay = nodelay.(int)
	// interval default 20, Usage: "ikcp_nodelay()",
	interval, err := server.ParseValue(m, "interval", 20)
	if err != nil {
		return err
	}
	c.Interval = interval.(int)
	// resend default 1, Usage: "ikcp_nodelay()",
	resend, err := server.ParseValue(m, "resend", 1)
	if err != nil {
		return err
	}
	c.Resend = resend.(int)
	// nc default 1,Usage: "ikcp_nodelay()"
	nc, err := server.ParseValue(m, "nc", 1)
	if err != nil {
		return err
	}
	c.Nc = nc.(int)
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
