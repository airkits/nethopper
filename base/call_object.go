package base

import (
	"sync"
	"time"
)

const (
	// CallObejctType call object type
	CallObejctType      = iota
	CallObejctNone      = 1 //普通模式
	CallObejctNotify    = 2 //通知模式，无响应
	CallObejctTransport = 3 //透传模式，handler可以收到callobject
)

type ICaller interface {
	Execute(obj *CallObject) *Ret
}

var GCallObjectPool = sync.Pool{
	New: func() interface{} {
		return &CallObject{
			Timer: time.NewTimer(TimeoutChanTime),
		}
	},
}

// CallObject call struct
type CallObject struct {
	Caller  ICaller
	CmdID   int32
	Option  int32
	Type    int8
	Args    []interface{}
	Trace   []uint8
	ChanRet chan *Ret
	Timer   *time.Timer
}

func (c *CallObject) Init(t int8, caller ICaller, cmdID int32, opt int32, args ...interface{}) *CallObject {
	c.Caller = caller
	c.CmdID = cmdID
	c.Option = opt
	c.Args = args
	c.Type = t
	if c.Type != CallObejctNotify {
		c.ChanRet = make(chan *Ret, 1)
	}
	c.Timer.Reset(TimeoutChanTime)
	c.Trace = make([]uint8, 0, 3)
	return c
}
func (c *CallObject) SetTrace(mid ...uint8) {
	c.Trace = append(c.Trace, mid...)
}

func (c *CallObject) Reset() *CallObject {
	c.Args = nil
	c.Caller = nil
	c.CmdID = 0
	c.Option = 0
	c.Trace = nil
	c.Type = 0
	if c.ChanRet != nil {
		close(c.ChanRet)
		c.ChanRet = nil
	}
	GCallObjectPool.Put(c)
	return c
}

type Callback func(ret *Ret)

// Ret call return object
type Ret struct {
	Data  interface{}
	Code  int32
	Err   error
	Trace []uint8
}

// NewCallObject create call object
func NewCallObject(caller ICaller, cmdID int32, opt int32, args ...interface{}) *CallObject {
	obj := GCallObjectPool.Get()
	return obj.(*CallObject).Init(CallObejctNone, caller, cmdID, opt, args...)
}

// NewNotifyObject create notify object
func NewNotifyObject(caller ICaller, cmdID int32, opt int32, args ...interface{}) *CallObject {
	obj := GCallObjectPool.Get()
	return obj.(*CallObject).Init(CallObejctNotify, caller, cmdID, opt, args...)
}

// NewTransportObject create transport object
func NewTransportObject(caller ICaller, cmdID int32, opt int32, args ...interface{}) *CallObject {
	obj := GCallObjectPool.Get()
	return obj.(*CallObject).Init(CallObejctTransport, caller, cmdID, opt, args...)
}

// NewRet create ret object
func NewRet(code int32, err error, data interface{}) *Ret {
	obj := &Ret{}
	return obj.Init(code, err, data)
}

func (c *Ret) Init(code int32, err error, data interface{}) *Ret {
	c.Code = code
	c.Data = data
	c.Err = err
	c.Trace = make([]uint8, 0, 3)
	return c
}
func (c *Ret) SetTrace(mid ...uint8) {
	c.Trace = append(c.Trace, mid...)
}

func (c *Ret) Reset() *Ret {
	c.Code = 0
	c.Data = nil
	c.Err = nil
	c.Trace = nil
	return c
}
