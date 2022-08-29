package natsrpc

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/airkits/nethopper/base"
	"github.com/airkits/nethopper/codec/json"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/mq"
	"github.com/airkits/nethopper/network"
	"github.com/airkits/nethopper/utils"
	"github.com/airkits/proto/s2s"
	"github.com/airkits/proto/ss"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

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
	nc        *nats.Conn
	stream    nats.JetStreamContext
	sendChan  chan *ss.Message
	recvChan  chan *ss.Message
	closeFlag bool
	sendCount int64
	services  sync.Map
	Conf      *NatsConfig
}

// NewConn create websocket conn
func NewConn(conn *nats.Conn, conf *NatsConfig) network.IConn {
	natsConn := &Conn{}
	natsConn.nc = conn
	natsConn.Conf = conf
	natsConn.sendCount = 0
	js, err := conn.JetStream(nats.PublishAsyncMaxPending(int(conf.AsyncMaxPending)),
		nats.PublishAsyncErrHandler(func(stream nats.JetStream, msg *nats.Msg, err error) {
			// todo jetstream error handling
			fmt.Println(err.Error())
			outMsg := &ss.Message{}
			_ = proto.Unmarshal(msg.Data, outMsg)
			if err != nil {
				natsConn.recvChan <- natsConn.DirectErrorMsg(outMsg, err)
			}
		}),
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	natsConn.stream = js
	natsConn.sendChan = make(chan *ss.Message, conf.SocketQueueSize)
	natsConn.recvChan = make(chan *ss.Message, conf.SocketQueueSize)

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
	m.MsgID = uint32(s2s.MessageCmd_ERROR)
	body := &s2s.ErrorResp{
		Result: &s2s.Result{
			Code: base.ErrCodeRouter,
			Msg:  err.Error(),
		},
		Time: utils.LocalMilliscond(),
	}
	any, _ := anypb.New(body)
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
		Time:     m.Time,
		Reply:    m.Reply,
		Options:  map[string][]byte{"e": []byte(err.Error())},
		Body:     any,
	}

	return msg
}
func (c *Conn) GetStreamName(msgType, srcType, srcID uint32) string {
	if msgType == mq.MTBroadcast {
		return fmt.Sprintf("gjetst%ds%d", srcType, srcID)
	} else if msgType == mq.MTRequestAny {
		return fmt.Sprintf("replyt%ds%d", srcType, srcID)
	} else if msgType == mq.MTRequestPush {
		return fmt.Sprintf("pusht%ds%d", srcType, srcID)
	}
	return fmt.Sprintf("jetst%ds%d", srcType, srcID)
}
func (c *Conn) GetSubject(msgType, destType, destID, srcType, srcID uint32) string {
	if msgType == mq.MTBroadcast {
		return fmt.Sprintf("gjetst%ds%d.t%ds%d", destType, destID, srcType, srcID)
	} else if msgType == mq.MTRequestAny || msgType == mq.MTResponseAny {
		return fmt.Sprintf("replyt%ds%d.t%ds%d", destType, destID, srcType, srcID)
	} else if msgType == mq.MTRequestPush || msgType == mq.MTResponsePush {
		return fmt.Sprintf("pusht%ds%d.t%ds%d", destType, destID, srcType, srcID)
	}
	return fmt.Sprintf("jetst%ds%d.t%ds%d", destType, destID, srcType, srcID)

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
	os, err := c.getObjectStore()
	if err != nil {
		return err
	}
	for i := 0; i < len(c.Conf.Services); i++ {
		c.LoadServiceInfo(os, &c.Conf.Services[i])
	}
	go func() {

		// Create key watcher.
		wopts := []nats.WatchOpt{}
		watcher, err := os.Watch(wopts...)
		if err != nil {
			fmt.Printf("ERROR: nats.KeyValue.WatchAll failed, err: %v", err)
		}
		for {
			select {
			case kve := <-watcher.Updates():
				if kve != nil {
					fmt.Printf("RECV: key: %v", kve)
				}
			case <-time.After(base.TimeoutChanTime):
				continue
			}
		}
	}()
	return nil
}
func (c *Conn) RegisterStream(msgType, srcType, srcID uint32) error {
	name := c.GetStreamName(msgType, srcType, srcID)
	subject := fmt.Sprintf("%s.*", name)
	maxConsumers := 1
	if msgType == mq.MTBroadcast {
		maxConsumers = 1024
	}
	if err := c.createStream(name, []string{subject}, maxConsumers); err != nil {
		log.Error("[NatsRPC] Create or Update stream error %s", err.Error())
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
func (c *Conn) getObjectStore() (nats.ObjectStore, error) {
	cfg := &nats.ObjectStoreConfig{Bucket: NatsServiceKey}
	return c.stream.CreateObjectStore(cfg)
}
func (c *Conn) LoadServiceInfo(os nats.ObjectStore, localInfo *ServiceGroup) error {

	result, err := os.GetString(localInfo.Key)
	if err != nil {
		infoByte, err1 := json.Marshal(localInfo)
		if err1 != nil {
			return err1
		}
		os.PutString(localInfo.Key, string(infoByte))
		c.services.Store(localInfo.Type, localInfo)
		return err
	}
	remoteInfo := &ServiceGroup{}
	err = json.Unmarshal([]byte(result), remoteInfo)
	if err != nil {
		return err
	}
	if localInfo.Version > remoteInfo.Version {
		infoByte, err1 := json.Marshal(localInfo)
		if err1 != nil {
			return err1
		}
		os.PutString(localInfo.Key, string(infoByte))
		c.services.Store(localInfo.Type, localInfo)
	} else {
		c.services.Store(localInfo.Type, remoteInfo)
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
func (c *Conn) createStream(name string, subjects []string, maxConsumers int) error {

	_, err := c.stream.StreamInfo(name)
	conf := &nats.StreamConfig{
		Name:         name,
		Subjects:     subjects,
		MaxConsumers: maxConsumers,
		MaxMsgs:      1000000, // unlimitted
		MaxBytes:     -1,      // stream size unlimitted
		MaxAge:       1 * 24 * time.Hour,
		Duplicates:   1 * time.Hour,
	}

	if err != nil {
		_, err = c.stream.AddStream(conf, nats.PublishAsyncMaxPending(100000))
	} else {
		_, err = c.stream.UpdateStream(conf, nats.PublishAsyncMaxPending(100000))
	}

	if err != nil {
		return err
	}

	return nil
}

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
	log.Info("[NatsRPC] SubscribeToReply %s to %s", name, subject)
	result, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
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
	log.Info("[NatsRPC] SubscribeToNats %s to %s", name, subject)
	result, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
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
	log.Info("[NatsRPC] SubscribeToStream %s to %s", name, subject)
	sub, err := c.stream.Subscribe(subject, func(msg *nats.Msg) {
		//	log.Info("recv msg from stream %s %v ", subject, msg)

		msg.Ack()
		//	fmt.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		ss := &ss.Message{}
		proto.Unmarshal(msg.Data, ss)
		c.recvChan <- ss
	}, nats.Durable(name), nats.ManualAck())
	if err != nil {
		fmt.Println(err.Error())
	}
	msgLimit, byteLimit, _ := sub.PendingLimits()
	log.Info("[NatsRPC] subscribe stream success,msgLimit:%d byteLimit:%d", msgLimit, byteLimit)
	sub.SetPendingLimits(-1, -1)
	fmt.Println(sub)
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
	//log.Info("publish msg to stream %s %v ", subject, msg)
	data, err := proto.Marshal(msg)
	if err != nil {
		return c.DirectErrorMsg(msg, err), err
	}

	c.sendCount += 1
	_, err2 := c.stream.PublishAsync(subject, data)
	if err2 != nil {
		fmt.Println(err2.Error())
		return c.DirectErrorMsg(msg, err2), err2
	}
	if c.sendCount > 5120 {
		select {
		case <-c.stream.PublishAsyncComplete():
			c.sendCount = 0
		case <-time.After(50 * time.Millisecond):
			fmt.Println("publish async Did not resolve in time")
		}
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
	c.sendChan <- b
	return nil
	// select {
	// case c.sendChan <- b:
	// 	return nil
	// case <-time.After(10 * time.Second):
	// 	return base.ErrReadChanTimeout
	// }
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
	return nil, base.ErrQueueIsClosed
	// select {
	// case v, ok := <-c.recvChan:
	// 	if ok {
	// 		return v, nil
	// 	}
	// 	return nil, base.ErrQueueIsClosed
	// case <-time.After(10 * time.Second):
	// 	return nil, base.ErrReadChanTimeout
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
