package natsrpc

import (
	"github.com/airkits/proto/ss"
	"github.com/nats-io/nats.go"
)

//NewStream create nats stream
func NewStream(conn *nats.Conn) INatsStream {
	s := new(NatsStream)
	s.conn = conn
	js, _ := conn.JetStream(nats.PublishAsyncMaxPending(256))
	s.js = js
	return s
}

//NatsStream define nat stream interface
type NatsStream struct {
	js   nats.JetStreamContext
	conn *nats.Conn
}

func (s *NatsStream) Send(*ss.Message) error {

	return nil
}
func (s *NatsStream) Recv() (*ss.Message, error) {
	return nil, nil
}

// Context returns the context for this stream.
func (s *NatsStream) Context() nats.JetStreamContext {
	return s.js
}

func (s *NatsStream) Close() {
	s.conn.Close()
}
