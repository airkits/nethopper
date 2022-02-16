// MIT License

// Copyright (c) 2019 gonethopper

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * @Author: ankye
// * @Date: 2019-06-14 14:15:06
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-14 14:15:06

package mediator

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/airkits/nethopper/base/queue"
	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
)

// IModule interface define
type IModule interface {
	//BaseContext start
	// ID module id
	ID() int32
	//SetID set module ID
	SetID(v int32)
	// Name module name
	Name() string
	//SetName set module name
	SetName(v string)

	//Handlers set moudle handlers
	Handlers() map[int32]interface{}
	//ReflectHandlers set moudle reflect handlers
	ReflectHandlers() map[int32]interface{}

	// RegisterHandler register function before run
	RegisterHandler(id int32, f interface{})
	// RegisterReflectHandler register reflect function before run
	RegisterReflectHandler(id int32, f interface{})

	// MakeContext init base module queue and create context
	MakeContext(queueSize int32)
	// Context get module context
	Context() context.Context
	// Close call context cancel ,self and all child module will receive context.Done()
	Close()
	// Queue return module queue
	MQ() queue.Queue
	// CanExit if receive ctx.Done() and child ref = 0 and queue is empty ,then return true
	CanExit(doneflag bool) (bool, bool)

	//BaseContext end

	// UserData module custom option, can you store you data and you must keep goruntine safe
	UserData() int32

	HasWorkerPool() bool

	WorkerPoolSubmit(obj *CallObject) error
	// Setup init custom module and pass config map to module
	Setup(conf config.IConfig) (IModule, error)
	//Reload reload config
	Reload(conf config.IConfig) error
	// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
	OnRun(dt time.Duration)
	// Stop goruntine
	Stop() error
	// Call async send callobject to module
	Call(option int32, obj *CallObject) error
	// Execute callobject
	Execute(obj *CallObject) *RetObject

	// PushBytes async send string or bytes to queue
	PushBytes(option int32, buf []byte) error
	//GetHandler get call handler
	GetHandler(id int32) interface{}

	//GetReflectHandler get reflect handler
	GetReflectHandler(id int32) interface{}

	// DoWorker dispatch callobject to worker processor
	DoWorker(obj *CallObject) error

	//IdleTimesReset reset idle times
	// IdleTimesReset()

	// // IdleTimes get idle times
	// IdleTimes() uint32

	// // IdleTimesAdd add idle times
	// IdleTimesAdd()
}

//BaseContext use context to close all module and using the bubbling method to exit
type BaseContext struct {
	ctx        context.Context
	cancel     context.CancelFunc
	parent     IModule
	childRef   int32
	q          queue.Queue
	name       string
	id         int32
	funcs      map[int32]interface{} //handlers
	rfuncs     map[int32]interface{} //reflect handlers
	workerPool IWorkerPool
	// idleTimes  uint32
}

//Handlers set moudle handlers
func (s *BaseContext) Handlers() map[int32]interface{} {
	return nil
}

//ReflectHandlers set moudle reflect handlers
func (s *BaseContext) ReflectHandlers() map[int32]interface{} {
	return nil
}

// RegisterHandler register function before run
func (s *BaseContext) RegisterHandler(id int32, f interface{}) {

	switch f.(type) {
	case func(interface{}) *RetObject:
	case func(interface{}, interface{}) *RetObject:
	case func(interface{}, interface{}, interface{}) *RetObject:
	case func(interface{}, interface{}, interface{}, interface{}) *RetObject:
	default:
		panic(fmt.Sprintf("function id %v: definition of function is invalid,%v", id, reflect.TypeOf(f)))
	}

	if _, ok := s.funcs[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	s.funcs[id] = f
}

// RegisterReflectHandler register reflect function before run
func (s *BaseContext) RegisterReflectHandler(id int32, f interface{}) {

	if _, ok := s.rfuncs[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	s.rfuncs[id] = f
}

// GetHandler get call handler
func (s *BaseContext) GetHandler(id int32) interface{} {
	return s.funcs[id]
}

// GetReflectHandler get call reflect handler
func (s *BaseContext) GetReflectHandler(id int32) interface{} {
	return s.rfuncs[id]
}

// Execute callobject
func (s *BaseContext) Execute(obj *CallObject) *RetObject {
	return NewRetObject(-1, fmt.Errorf("must override execute"), nil)
}

// IdleTimesReset reset idle times
// func (s *BaseContext) IdleTimesReset() {
// 	atomic.StoreUint32(&s.idleTimes, 500)
// }

// IdleTimes get idle times
// func (s *BaseContext) IdleTimes() uint32 {
// 	return atomic.LoadUint32(&s.idleTimes)
// }

// // IdleTimesAdd add idle times
// func (s *BaseContext) IdleTimesAdd() {
// 	t := s.IdleTimes()
// 	if t >= 20000000 { //2s
// 		return
// 	}
// 	atomic.AddUint32(&s.idleTimes, 100)
// }

// MakeContext init base module queue and create context
func (s *BaseContext) MakeContext(queueSize int32) {
	//s.parent = p
	s.q = queue.NewChanQueue(queueSize)
	s.funcs = make(map[int32]interface{})
	s.rfuncs = make(map[int32]interface{})
	// if p == nil {
	// 	s.ctx, s.cancel = context.WithCancel(context.Background())
	// } else {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	// 	p.ChildAdd()
	// }

}

func (s *BaseContext) HasWorkerPool() bool {
	return s.workerPool != nil
}
func (s *BaseContext) WorkerPoolSubmit(obj *CallObject) error {
	return s.workerPool.Submit(obj)
}

// DoWorker process callobject
func (s *BaseContext) DoWorker(obj *CallObject) error {
	//Debug("[%s] cmd [%s] process", s.Name(), obj.Cmd)
	var err error
	if s.workerPool == nil {
		//err = errors.New("no processor pool")
		result := s.Execute(obj)
		result.SetTrace(uint8(s.ID()))
		obj.ChanRet <- result
		return nil
	} else {
		err = s.workerPool.Submit(obj)
	}
	if err != nil {
		result := NewRetObject(-1, err, nil)
		result.SetTrace(uint8(s.ID()))
		obj.ChanRet <- result
	}
	return err
}

// Call async send message to module
func (s *BaseContext) Call(option int32, obj *CallObject) error {
	//	s.IdleTimesReset()
	if err := s.q.AsyncPush(obj); err != nil {
		log.Error(err.Error())
	}

	return nil
}

// CreateWorkerPool create processor pool
func (s *BaseContext) CreateWorkerPool(cap uint32, queueSize uint32, expired time.Duration, isNonBlocking bool) (err error) {
	if s.workerPool, err = NewFixedWorkerPool(cap, queueSize, expired); err != nil {
		return err
	}
	return nil
}

// MQ return module queue
func (s *BaseContext) MQ() queue.Queue {
	return s.q
}

// Context get module context
func (s *BaseContext) Context() context.Context {
	return s.ctx
}

// Close call context cancel ,self and all child module will receive context.Done()
func (s *BaseContext) Close() {
	s.cancel()
}

//ID module ID
func (s *BaseContext) ID() int32 {
	return s.id
}

//SetID set module id
func (s *BaseContext) SetID(v int32) {
	s.id = v
}

//Name module name
func (s *BaseContext) Name() string {
	return s.name
}

//SetName set module name
func (s *BaseContext) SetName(v string) {
	s.name = v
}

// CanExit if receive ctx.Done() and all child exit and queue is empty ,then return true
func (s *BaseContext) CanExit(doneFlag bool) (bool, bool) {
	if doneFlag {
		if s.q.Length() == 0 {
			return doneFlag, true
		}
	}
	select {
	case <-s.ctx.Done():
		doneFlag = true
		if s.q.Length() == 0 {
			return doneFlag, true
		}
	default:
	}
	return doneFlag, false
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *BaseContext) OnRun(dt time.Duration) {
	fmt.Printf("module %s do Nothing \n", s.Name())

}

// to override start

//PushBytes push buffer
func (s *BaseContext) PushBytes(option int32, buf []byte) error {
	return nil
}

// UserData module custom option, can you store you data and you must keep goruntine safe
func (s *BaseContext) UserData() int32 {
	return 0
}

// ReadConfig config map
// m := map[string]interface{}{
// }
func (s *BaseContext) ReadConfig(conf config.IConfig) error {
	return nil
}

//Reload reload config
func (s *BaseContext) Reload(conf config.IConfig) error {
	return nil
}

// Stop goruntine
func (s *BaseContext) Stop() error {

	return nil
}

//to override end

func cmdRegister(s IModule) {
	cmds := s.Handlers()
	if cmds != nil {
		for k, v := range cmds {
			s.RegisterHandler(k, v)
		}
	}
	cmds = s.ReflectHandlers()
	if cmds != nil {
		for k, v := range cmds {
			s.RegisterReflectHandler(k, v)
		}
	}
}
