package natsrpc

import (
	"fmt"

	"github.com/airkits/proto/ss"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

//NewStream create nats stream
func NewStream(conn *nats.Conn) INatsStream {
	s := new(NatsStream)
	s.conn = conn
	js, err := conn.JetStream(nats.PublishAsyncMaxPending(256))
	s.js = js
	if err != nil {
		fmt.Println(err)
	}
	s.createStream()
	s.SubscribeToStream(func(data []byte) []byte {
		return data
	})
	return s
}

//NatsStream define nat stream interface
type NatsStream struct {
	js         nats.JetStreamContext
	conn       *nats.Conn
	streamInfo *nats.StreamInfo
}

func (s *NatsStream) createStream() error {
	name := "query"
	if s.streamInfo != nil {
		err := s.js.DeleteStream(name)
		if err != nil {
			return err
		}
	}

	info, err := s.js.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: []string{"query.*"},
	})
	if err != nil {
		return err
	}
	s.streamInfo = info
	return nil
}

func (s *NatsStream) SubscribeToStream(f func([]byte) []byte) {
	fmt.Printf("Subscribing to query.serialized")
	result, err := s.js.Subscribe("query.test", func(msg *nats.Msg) {
		fmt.Printf("Msg recieved")
		msg.Ack()
		fmt.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		res := f(msg.Data)
		fmt.Printf("Data recieved: ")
		s.PublishToStream(res)
	}, nats.Durable("monitor"), nats.ManualAck())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

func (s *NatsStream) PublishToStream(data []byte) {

	_, err := s.js.Publish("query.test", data)

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Published queryJSON:%s to subjectName:query.serialized", string(data))

}

func (s *NatsStream) Send(msg *ss.Message) error {
	data, _ := proto.Marshal(msg)
	s.PublishToStream(data)
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
