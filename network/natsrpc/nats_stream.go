package natsrpc

import (
	"fmt"
	"runtime"

	"github.com/airkits/proto/ss"
	"github.com/nats-io/nats.go"
)

//NewStream create nats stream
func NewStream(conn *nats.Conn) INatsStream {
	s := new(NatsStream)
	s.conn = conn
	js, _ := conn.JetStream(nats.PublishAsyncMaxPending(256))
	s.js = js
	s.createStream()

	return s
}

//NatsStream define nat stream interface
type NatsStream struct {
	js   nats.JetStreamContext
	conn *nats.Conn
}

func (s *NatsStream) createStream() error {
	err := s.js.DeleteStream("query.*")
	if err != nil {
		return err
	}
	_, err = s.js.AddStream(&nats.StreamConfig{
		Name:     "query.*",
		Subjects: []string{"query.*"},
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *NatsStream) SubscribeToStream(f func([]byte) []byte) {
	fmt.Printf("Subscribing to query.unserialized")
	s.js.Subscribe("query.unserialized", func(msg *nats.Msg) {
		fmt.Printf("Msg recieved")
		msg.Ack()
		fmt.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		res := f(msg.Data)
		fmt.Printf("Data recieved: ")
		s.PublishToStream(res)
	}, nats.Durable("monitor"), nats.ManualAck())
	runtime.Goexit()
}

func (s *NatsStream) PublishToStream(data []byte) {
	_, err := s.js.Publish("query.serialized", data)

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Published queryJSON:%s to subjectName:query.serialized", string(data))

}

func (s *NatsStream) Send(msg *ss.Message) error {

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
