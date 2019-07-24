package tcp

import (
	"net"

	"github.com/gonethopper/nethopper/server"
)

// TCPConnect listen port and accept tcp connection
type TCPConnect struct {
	Address         string
	Network         string
	ReadBufferSize  int
	WriteBufferSize int
	tcpListener     *net.TCPListener
}

// Setup init Connect with config
func (c *TCPConnect) Setup(m map[string]interface{}) (*TCPConnect, error) {
	if err := c.readConfig(m); err != nil {
		panic(err)
	}
	return c, nil
}

// config map
// readBufferSize default 32767
// writeBufferSize default 32767
// address default :8080
// network default "tcp4"  use "tcp4/tcp6"
func (c *TCPConnect) readConfig(m map[string]interface{}) error {
	readBufferSize, err := server.ParseValue(m, "readBufferSize", 32767)
	if err != nil {
		return err
	}
	c.ReadBufferSize = readBufferSize.(int)

	writeBufferSize, err := server.ParseValue(m, "writeBufferSize", 32767)
	if err != nil {
		return err
	}
	c.WriteBufferSize = writeBufferSize.(int)

	address, err := server.ParseValue(m, "address", ":8080")
	if err != nil {
		return err
	}
	c.Address = address.(string)
	network, err := server.ParseValue(m, "network", "tcp4")
	if err != nil {
		return err
	}
	c.Network = network.(string)

	return nil
}

// Listen and bind local ip
func (c *TCPConnect) Listen() {

	tcpAddr, err := net.ResolveTCPAddr(c.Network, c.Address)
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	c.tcpListener = listener
	server.Info("listening on: %s %s", c.Network, listener.Addr())

	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			server.Warning("accept failed: %s", err.Error())
			continue
		}
		// set socket read buffer
		conn.SetReadBuffer(c.ReadBufferSize)
		// set socket write buffer
		conn.SetWriteBuffer(c.WriteBufferSize)
		// start a goroutine for every incoming connection for reading
		//	go handleClient(conn, config)
	}
}
