package natsrpc

import (
	"fmt"
	"net"
	"sync"

	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/proto/ss"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

//ConnSet grpc conn set
type ConnSet map[*nats.Conn]struct{}

//INatsStream define rpc stream interface
type INatsStream interface {
	Send(*ss.Message) error
	Recv() (*ss.Message, error)
	// Context returns the context for this stream.
	Context() nats.JetStreamContext
	Close()
}

//Conn grpc conn define
type Conn struct {
	sync.Mutex
	nc             *nats.Conn
	stream         nats.JetStreamContext
	writeChan      chan *ss.Message
	readChan       chan *ss.Message
	maxMessageSize uint32
	closeFlag      bool
	streamInfo     *nats.StreamInfo
}

//NewConn create websocket conn
func NewConn(conn *nats.Conn, rwQueueSize int, maxMessageSize uint32) network.IConn {
	natsConn := new(Conn)
	natsConn.nc = conn

	natsConn.writeChan = make(chan *ss.Message, rwQueueSize)
	natsConn.readChan = make(chan *ss.Message, rwQueueSize)
	natsConn.maxMessageSize = maxMessageSize
	natsConn.CreateStream("query", []string{"query.*"})
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.PrintStack(false)
			}
		}()
		for b := range natsConn.writeChan {
			if b == nil {
				break
			}

			//err := stream.Send(b)
			// if err != nil {
			// 	break
			// }
		}

		//	conn.Close()
		natsConn.Lock()
		natsConn.closeFlag = true
		natsConn.Unlock()
	}()

	return natsConn
}

func (s *Conn) CreateStream(name string, subjects []string) error {

	if s.streamInfo != nil {
		err := s.stream.DeleteStream(name)
		if err != nil {
			return err
		}
	}

	info, err := s.stream.AddStream(&nats.StreamConfig{
		Name:     name,
		Subjects: subjects,
	})
	if err != nil {
		return err
	}
	s.streamInfo = info
	return nil
}

func (c *Conn) SubscribeToStream(subject string) {
	fmt.Printf("Subscribing to query.serialized")
	result, err := c.stream.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("Msg recieved")
		msg.Ack()
		fmt.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		ss := &ss.Message{}
		proto.Unmarshal(msg.Data, ss)
		c.readChan <- ss
	}, nats.Durable("monitor"), nats.ManualAck())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}
func (c *Conn) PublishToStream(subject string, msg *ss.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.stream.Publish(subject, data)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
func (c *Conn) doDestroy() {
	// c.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	// c.conn.Close()

	if !c.closeFlag {
		close(c.writeChan)
		c.closeFlag = true
	}
}

//Destroy grpc conn destory
func (c *Conn) Destroy() {
	c.Lock()
	defer c.Unlock()

	c.doDestroy()
}

//Close grpc conn close
func (c *Conn) Close() {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return
	}

	c.doWrite(nil)
	c.closeFlag = true
}

func (c *Conn) doWrite(b *ss.Message) {
	if len(c.writeChan) == cap(c.writeChan) {
		log.Error("close conn: channel full")
		c.doDestroy()
		return
	}

	c.writeChan <- b
}

//LocalAddr get local addr
func (c *Conn) LocalAddr() net.Addr {
	log.Error("[LocalAddr] invoke LocalAddr() unsupport")
	return nil
}

//RemoteAddr get remote addr
func (c *Conn) RemoteAddr() net.Addr {
	// pr, ok := peer.FromContext(c.stream.Context())
	// if !ok {
	// 	log.Error("[RemoteAddr] invoke FromContext() failed")
	// 	return nil
	// }
	// if pr.Addr == net.Addr(nil) {
	// 	log.Error("[RemoteAddr] peer.Addr is nil")
	// 	return nil
	// }

	//return pr.Addr
	return nil
}

//ReadMessage goroutine not safe
func (c *Conn) ReadMessage() (interface{}, error) {
	return c.stream.Recv()
}

//WriteMessage args must not be modified by the others goroutines
func (c *Conn) WriteMessage(args ...interface{}) error {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return nil
	}

	for i := 0; i < len(args); i++ {
		c.doWrite(args[i].(*ss.Message))
	}
	return nil
}
