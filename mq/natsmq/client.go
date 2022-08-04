package natsmq

//NatsClient define nat stream client interface
type NatsClient struct {
	conf *ClientConfig
}

func (c *NatsClient) Init(conf *ClientConfig) *NatsClient {
	c.conf = conf
	return c
}

func (c *NatsClient) Connect() {

}

func (c *NatsClient) ReConnect() {

}

func (c *NatsClient) Request() {

}

func (c *NatsClient) Subscription() {

}

func (c *NatsClient) SendMessage() {

}

func (c *NatsClient) RecvMessage() {

}
