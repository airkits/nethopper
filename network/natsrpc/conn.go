package natsrpc

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/mq"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/proto/ss"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/types/known/anypb"
)

// ErrQueueIsClosed queue is closed
var ErrQueueIsClosed = errors.New("Queue is Closed")

var ErrQueueEmpty = errors.New("Queue is empty")
var ErrQueueFull = errors.New("Queue is full")

// ConnSet grpc conn set
type ConnSet map[*nats.Conn]struct{}

// INatsStream define rpc stream interface
type INatsStream interface {
	Send(*ss.Message) error
	Recv() (*ss.Message, error)
	// Context returns the context for this stream.
	Context() nats.JetStreamContext
	Close()
}

// Conn grpc conn define
type Conn struct {
	sync.Mutex
	nc             *nats.Conn
	stream         nats.JetStreamContext
	sendChan       chan *ss.Message
	recvChan       chan *ss.Message
	maxMessageSize uint32
	closeFlag      bool

	funcs map[string](func(*ss.Message) *ss.Message) //handlers
}

// NewConn create websocket conn
func NewConn(conn *nats.Conn, rwQueueSize int, maxMessageSize uint32) network.IConn {
	natsConn := &Conn{}
	natsConn.nc = conn
	natsConn.funcs = make(map[string](func(*ss.Message) *ss.Message))
	js, err := conn.JetStream(nats.PublishAsyncMaxPending(4096),
		nats.PublishAsyncErrHandler(func(stream nats.JetStream, msg *nats.Msg, err error) {
			// todo jetstream error handling
			fmt.Println(err.Error())
		}),
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	natsConn.stream = js
	natsConn.sendChan = make(chan *ss.Message, rwQueueSize)
	natsConn.recvChan = make(chan *ss.Message, rwQueueSize)
	natsConn.maxMessageSize = maxMessageSize

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.PrintStack(false)
			}
		}()
		for m := range natsConn.sendChan {
			if m == nil {
				break
			}
			subject := natsConn.GetSubject(m.MsgType, m.DestType, m.DestID, m.SrcType, m.SrcID)
			if m.MsgType == mq.MTRequestAny {
				msg1, err1 := natsConn.Request(subject, m)
				if err1 != nil {
					natsConn.recvChan <- msg1
					continue
				}
				natsConn.recvChan <- msg1
			} else if m.MsgType == mq.MTResponseAny {
				msg4, err4 := natsConn.Reply(m)
				if err4 != nil {
					natsConn.recvChan <- msg4
					continue
				}
			} else if m.MsgType == mq.MTRequestPush || m.MsgType == mq.MTResponsePush {
				msg2, err2 := natsConn.publishToNats(subject, m)
				if err2 != nil {
					natsConn.recvChan <- msg2
					continue
				}
			} else {
				msg3, err3 := natsConn.publishToStream(subject, m)
				if err3 != nil {
					natsConn.recvChan <- msg3
					continue
				}
			}

		}
		//	conn.Close()
		natsConn.Lock()
		natsConn.closeFlag = true
		natsConn.Unlock()
	}()

	return natsConn
}
func (c *Conn) DirectErrorMsg(m *ss.Message, err error) *ss.Message {
	msgType := m.MsgType
	if msgType == mq.MTRequest {
		msgType = mq.MTResponse
	} else if msgType == mq.MTRequestAny {
		msgType = mq.MTResponseAny
	} else if msgType == mq.MTRequestPush {
		msgType = mq.MTResponsePush
	}
	msg := &ss.Message{
		ID:       m.ID,
		UID:      m.UID,
		MsgID:    m.MsgID,
		MsgType:  msgType,
		Seq:      m.Seq,
		SrcType:  m.SrcType,
		SrcID:    m.SrcID,
		DestType: m.SrcType,
		DestID:   m.SrcID,
		Time:     0,
		Reply:    err.Error(),
		Options:  map[string][]byte{},
		Body:     &anypb.Any{},
	}
	return msg
}
func (c *Conn) GetStreamName(msgType, srcType, srcID uint32) string {
	if msgType == mq.MTBroadcast {
		return fmt.Sprintf("gjst%ds%d", srcType, srcID)
	} else if msgType == mq.MTRequestAny {
		return fmt.Sprintf("anyt%ds%d", srcType, srcID)
	} else if msgType == mq.MTRequestPush {
		return fmt.Sprintf("pusht%ds%d", srcType, srcID)
	}
	return fmt.Sprintf("jst%ds%d", srcType, srcID)
}
func (c *Conn) GetSubject(msgType, destType, destID, srcType, srcID uint32) string {
	if msgType == mq.MTBroadcast {
		return fmt.Sprintf("gjst%ds%d.t%ds%d", destType, destID, srcType, srcID)
	} else if msgType == mq.MTRequestAny {
		return fmt.Sprintf("anyt%ds%d.t%ds%d", destType, destID, srcType, srcID)
	} else if msgType == mq.MTRequestPush {
		return fmt.Sprintf("pusht%ds%d.t%ds%d", destType, destID, srcType, srcID)
	}
	return fmt.Sprintf("jst%ds%d.t%ds%d", destType, destID, srcType, srcID)

}

func (c *Conn) RegisterService(srcType, srcID uint32) error {

	if err := c.RegisterStream(mq.MTBroadcast, srcType, srcID); err != nil {
		return err
	}
	if err := c.RegisterStream(mq.MTRequest, srcType, srcID); err != nil {
		return err
	}
	if err := c.RegisterSubject(mq.MTRequestAny, srcType, srcID); err != nil {
		return err
	}
	if err := c.RegisterSubject(mq.MTRequestPush, srcType, srcID); err != nil {
		return err
	}
	return nil
}
func (c *Conn) RegisterStream(msgType, srcType, srcID uint32) error {
	name := c.GetStreamName(msgType, srcType, srcID)
	subject := fmt.Sprintf("%s.*", name)
	if err := c.createStream(name, []string{subject}); err != nil {
		return err
	}
	c.SubscribeToStream(name, subject)
	return nil
}
func (c *Conn) RegisterSubject(msgType, srcType, srcID uint32) error {
	name := c.GetStreamName(msgType, srcType, srcID)
	subject := fmt.Sprintf("%s.*", name)
	if msgType == mq.MTRequestPush {
		c.SubscribeToNats(name, subject)
	} else if msgType == mq.MTRequestAny {
		c.SubscribeToReply(name, subject)
	}

	return nil
}

/*
*
最大年龄	流中任何消息的最长期限，以微秒为单位
最大字节数	当合并后的流大小超过此旧消息时，将删除Stream的大小
最大消息大小	流将接受的最大消息
最大消息	流中可能有多少条消息，如果流超过此大小，则最早的消息将被删除
MaxConsumers	可以为给定的流定义多少个消费者，
名称	流的名称，不能包含空格，制表符或.
NoAck	禁用确认流接收的消息
复制品	每个邮件要保留多少个副本（截至2020年1月尚未实现）
保留	如何考虑保留邮件
LimitsPolicy（默认）
InterestPolicy或WorkQueuePolicy丢弃	当流达到其限制时，
DiscardNew拒绝新消息，而
DiscardOld（默认）删除旧消息
存储	该类型的存储后端，file并memory
科目	要使用的主题列表，支持通配符
重复项	跟踪重复消息的窗口
*
*/
func (c *Conn) createStream(name string, subjects []string) error {

	info, err := c.stream.StreamInfo(name)
	conf := &nats.StreamConfig{
		Name:         name,
		Subjects:     subjects,
		MaxConsumers: 1,
		MaxMsgs:      -1, // unlimitted
		MaxBytes:     -1, // stream size unlimitted
		MaxAge:       7 * 24 * time.Hour,
		MaxMsgSize:   640000,
		Duplicates:   1 * time.Hour,
	}
	fmt.Print(info)
	if err != nil {
		info, err = c.stream.AddStream(conf, nats.PublishAsyncMaxPending(100000))
	} else {
		info, err = c.stream.UpdateStream(conf, nats.PublishAsyncMaxPending(100000))
	}
	fmt.Print(info)
	if err != nil {
		return err
	}

	return nil
}

// RegisterReply register function before run
// func (c *Conn) RegisterReply(subject string, f func(*ss.Message) *ss.Message) {

// 	if _, ok := c.funcs[subject]; ok {
// 		panic(fmt.Sprintf("function id %v: already registered", subject))
// 	}
// 	c.funcs[subject] = f
// 	c.reply(subject, f)
// }

func (c *Conn) Request(subject string, msg *ss.Message) (*ss.Message, error) {

	data, err := proto.Marshal(msg)
	if err != nil {
		return c.DirectErrorMsg(msg, err), err
	}
	result, err := c.nc.Request(subject, data, time.Second*30)
	if err != nil {
		return c.DirectErrorMsg(msg, err), err
	}
	outMsg := &ss.Message{}
	err = proto.Unmarshal(result.Data, outMsg)
	if err != nil {
		return c.DirectErrorMsg(msg, err), err
	}
	return outMsg, nil
}
func (c *Conn) Reply(m *ss.Message) (*ss.Message, error) {
	data, err := proto.Marshal(m)
	if err != nil {
		return c.DirectErrorMsg(m, err), err
	}
	err = c.nc.Publish(m.Reply, data)
	if err != nil {
		return c.DirectErrorMsg(m, err), err
	}
	return nil, nil
}
func (c *Conn) SubscribeToReply(name, subject string) {
	fmt.Printf("Subscribing to %s", subject)
	result, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		//	fmt.Printf("Msg recieved")
		//	fmt.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		ss := &ss.Message{}
		proto.Unmarshal(msg.Data, ss)
		ss.Reply = msg.Reply
		c.recvChan <- ss
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}
func (c *Conn) SubscribeToNats(name, subject string) {
	fmt.Printf("Subscribing to %s", subject)
	result, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		//	fmt.Printf("Msg recieved")
		//	fmt.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		ss := &ss.Message{}
		proto.Unmarshal(msg.Data, ss)
		c.recvChan <- ss
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}
func (c *Conn) SubscribeToStream(name, subject string) {
	fmt.Printf("Subscribing to %s", subject)
	result, err := c.stream.Subscribe(subject, func(msg *nats.Msg) {
		//	fmt.Printf("Msg recieved")
		msg.Ack()
		//	fmt.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		ss := &ss.Message{}
		proto.Unmarshal(msg.Data, ss)
		c.recvChan <- ss
	}, nats.Durable(name), nats.ManualAck())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

func (c *Conn) publishToNats(subject string, msg *ss.Message) (*ss.Message, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return c.DirectErrorMsg(msg, err), err
	}

	err2 := c.nc.Publish(subject, data)
	if err2 != nil {
		fmt.Println(err2.Error())
		return c.DirectErrorMsg(msg, err2), err2
	}
	//	fmt.Printf("\nsend reqid = %d,seq=%d \n", result.Sequence, msg.Seq)
	return nil, nil
}
func (c *Conn) publishToStream(subject string, msg *ss.Message) (*ss.Message, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return c.DirectErrorMsg(msg, err), err
	}

	_, err2 := c.stream.PublishAsync(subject, data)
	if err2 != nil {
		fmt.Println(err2.Error())
		return c.DirectErrorMsg(msg, err2), err2
	}
	//	fmt.Printf("\nsend reqid = %d,seq=%d \n", result.Sequence, msg.Seq)
	return nil, nil
}
func (c *Conn) doDestroy() {

	if !c.closeFlag {
		close(c.sendChan)
		c.closeFlag = true
	}
}

// Destroy grpc conn destory
func (c *Conn) Destroy() {
	c.Lock()
	defer c.Unlock()

	c.doDestroy()
}

// Close grpc conn close
func (c *Conn) Close() {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return
	}

	c.doWrite(nil)
	c.closeFlag = true

	if err := c.nc.FlushTimeout(5 * time.Second); err != nil {
		log.Error("error while flushing jetstream during close operation")
	}
	c.nc.Close()
}

func (c *Conn) doWrite(b *ss.Message) error {
	// if len(c.sendChan) == cap(c.sendChan) {
	// 	log.Error("close conn: channel full")
	// 	//c.doDestroy()
	// 	return ErrQueueFull
	// }

	c.sendChan <- b

	return nil
}

// LocalAddr get local addr
func (c *Conn) LocalAddr() net.Addr {
	log.Error("[LocalAddr] invoke LocalAddr() unsupport")
	return nil
}

// RemoteAddr get remote addr
func (c *Conn) RemoteAddr() net.Addr {

	return c.nc.Opts.Dialer.LocalAddr

}

// ReadMessage goroutine not safe
func (c *Conn) ReadMessage() (interface{}, error) {

	v, ok := <-c.recvChan
	if ok {
		return v, nil
	}
	return nil, ErrQueueIsClosed
	// select {
	// case v, ok := <-c.recvChan:
	// 	if ok {
	// 		return v, nil
	// 	}
	// 	return nil, ErrQueueIsClosed
	// default:
	// 	return nil, ErrQueueEmpty
	// }
}

// WriteMessage args must not be modified by the others goroutines
func (c *Conn) WriteMessage(args ...interface{}) error {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return nil
	}

	for i := 0; i < len(args); i++ {
		err := c.doWrite(args[i].(*ss.Message))
		if err != nil {
			return err
		}
	}
	return nil
}
