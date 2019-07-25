package http

import (
	"net"
	"net/http"

	"github.com/gonethopper/nethopper/server"
	"github.com/julienschmidt/httprouter"
)

// HTTPConnect http connect warpper
type HTTPConnect struct {
	Address string
	router  *httprouter.Router
}

// Setup init Connect with config
func (c *HTTPConnect) Setup(m map[string]interface{}) (*HTTPConnect, error) {
	if err := c.readConfig(m); err != nil {
		panic(err)
	}
	return c, nil
}

// config map
// address default :80
func (c *HTTPConnect) readConfig(m map[string]interface{}) error {

	address, err := server.ParseValue(m, "address", ":80")
	if err != nil {
		return err
	}
	c.Address = address.(string)

	return nil
}

// Listen and bind local ip
func (c *HTTPConnect) Listen() {

	router := httprouter.New()
	c.router = router
	if err := http.ListenAndServe(c.Address, router); err != nil {
		panic(err)
	}
	server.Info("listening on: %s %s", c.Address)

}

// Accept accepts the next incoming call and returns the new connection.
func (c *HTTPConnect) Accept() (net.Conn, error) {
	return nil, nil
}
