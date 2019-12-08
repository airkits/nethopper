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

package server

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/gonethopper/queue"
)

const (
	// ServiceNamedIDs service id define, system reserved 1-63
	ServiceNamedIDs = iota
	// ServiceIDMain main goruntinue
	ServiceIDMain
	// ServiceIDMonitor server monitor service
	ServiceIDMonitor
	// ServiceIDLog log service
	ServiceIDLog
	// ServiceIDTCP tcp service
	ServiceIDTCP
	// ServiceIDKCP kcp service
	ServiceIDKCP
	// ServiceIDHTTP http service
	ServiceIDHTTP
	// ServiceIDLogic logic service
	ServiceIDLogic
	// ServiceIDRedis redis service
	ServiceIDRedis
	// ServiceIDTCPClient tcp client service
	ServiceIDTCPClient
	// ServiceIDKCPClient kcp client service
	ServiceIDKCPClient
	// ServiceIDHTTPClient http client service
	ServiceIDHTTPClient
	// ServiceIDDB common db service
	ServiceIDDB
	// ServiceIDUserCustom User custom define named services from 64-128
	ServiceIDUserCustom = 64
	// ServiceIDNamedMax named services max ID
	ServiceIDNamedMax = 128
)

// Service interface define
type Service interface {
	//BaseContext start
	// ID service id
	ID() int32
	//SetID set service ID
	SetID(v int32)
	// Name service name
	Name() string
	//SetName set service name
	SetName(v string)

	// MakeContext init base service queue and create context
	MakeContext(p Service, queueSize int32)
	// Context get service context
	Context() context.Context
	// ChildAdd after child service created and tell parent service, ref count +1
	ChildAdd()
	// ChildDone child service exit and tell parent service, ref count -1
	ChildDone()
	// Close call context cancel ,self and all child service will receive context.Done()
	Close()
	// Queue return service queue
	MQ() queue.Queue
	// CanExit if receive ctx.Done() and child ref = 0 and queue is empty ,then return true
	CanExit(doneflag bool) (bool, bool)
	// TryExit check child ref count , if ref count == 0 then return true, if parent not nil, fire parent.ChildDone()
	TryExit() bool
	//BaseContext end

	// UserData service custom option, can you store you data and you must keep goruntine safe
	UserData() int32
	// Setup init custom service and pass config map to service
	Setup(m map[string]interface{}) (Service, error)
	//Reload reload config
	Reload(m map[string]interface{}) error
	// OnRun goruntine run and call OnRun , always use ServiceRun to call this function
	OnRun(dt time.Duration)
	// Stop goruntine
	Stop() error
	// Call async send callobject to service
	Call(option int32, obj *CallObject) error
	// PushBytes async send string or bytes to queue
	PushBytes(option int32, buf []byte) error
	//GetHandler get call handler
	GetHandler(id interface{}) interface{}
}

// ServiceRun wrapper service goruntine and in an orderly way to exit
func ServiceRun(s Service) {
	ctxDone := false
	exitFlag := false
	start := time.Now()
	Info("Service %s start ", s.Name())
	for {
		s.OnRun(time.Since(start))
		if ctxDone, exitFlag = s.CanExit(ctxDone); exitFlag {
			return
		}
		start = time.Now()
		if s.MQ().Length() == 0 {
			time.Sleep(time.Millisecond)
		}
		runtime.Gosched()
	}
}

// ServiceName get the service name
func ServiceName(s Service) string {
	t := reflect.TypeOf(s)
	return t.Elem().Name()
}

//BaseContext use context to close all service and using the bubbling method to exit
type BaseContext struct {
	ctx        context.Context
	cancel     context.CancelFunc
	parent     Service
	childRef   int32
	q          queue.Queue
	name       string
	id         int32
	functions  map[interface{}]interface{}
	processers IProcessorPool
}

// RegisterHandler register function before run
func (a *BaseContext) RegisterHandler(id interface{}, f interface{}) {

	// switch f.(type) {
	// case func(Service, *CallObject, string) (string, error):
	// default:
	// 	panic(fmt.Sprintf("function id %v: definition of function is invalid,%v", id, reflect.TypeOf(f)))
	// }

	if _, ok := a.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	a.functions[id] = f
}

// GetHandler get call handler
func (a *BaseContext) GetHandler(id interface{}) interface{} {
	return a.functions[id]
}

// MakeContext init base service queue and create context
func (a *BaseContext) MakeContext(p Service, queueSize int32) {
	a.parent = p
	a.q = queue.NewChanQueue(queueSize)
	a.functions = make(map[interface{}]interface{})
	if p == nil {
		a.ctx, a.cancel = context.WithCancel(context.Background())
	} else {
		a.ctx, a.cancel = context.WithCancel(p.Context())
		p.ChildAdd()
	}

}

// Processor process callobject
func (a *BaseContext) Processor(obj *CallObject) error {
	Debug("%s start do Processor,cmd = %s", a.Name(), obj.Cmd)
	var err error
	if a.processers == nil {
		err = errors.New("no processor pool")
	} else {
		err = a.processers.Submit(obj)
	}
	if err != nil {
		obj.ChanRet <- RetObject{
			Ret: nil,
			Err: err,
		}
	}
	return err
}

// CreateProcessorPool create processor pool
func (a *BaseContext) CreateProcessorPool(s Service, cap uint32, expired time.Duration, isNonBlocking bool) (err error) {
	if a.processers, err = NewFixedProcessorPool(s, cap, expired); err != nil {
		return err
	}
	return nil
}

// MQ return service queue
func (a *BaseContext) MQ() queue.Queue {
	return a.q
}

// Context get service context
func (a *BaseContext) Context() context.Context {
	return a.ctx
}

// ChildAdd child service created and tell parent service, ref count +1
func (a *BaseContext) ChildAdd() {
	atomic.AddInt32(&a.childRef, 1)
}

// ChildDone child service exit and tell parent service, ref count -1
func (a *BaseContext) ChildDone() {
	atomic.AddInt32(&a.childRef, -1)
}

// Close call context cancel ,self and all child service will receive context.Done()
func (a *BaseContext) Close() {
	a.cancel()
}

//ID service ID
func (a *BaseContext) ID() int32 {
	return a.id
}

//SetID set service id
func (a *BaseContext) SetID(v int32) {
	a.id = v
}

//Name service name
func (a *BaseContext) Name() string {
	return a.name
}

//SetName set service name
func (a *BaseContext) SetName(v string) {
	a.name = v
}

// TryExit check child ref count , if ref count == 0 then return true, if parent not nil, and will fire parent.ChildDone()
func (a *BaseContext) TryExit() bool {

	count := atomic.LoadInt32(&a.childRef)
	if count > 0 {
		return false
	}
	if a.parent != nil {
		a.parent.ChildDone()
	}
	return true
}

// CanExit if receive ctx.Done() and all child exit and queue is empty ,then return true
func (a *BaseContext) CanExit(doneFlag bool) (bool, bool) {
	if doneFlag {
		if a.q.Length() == 0 && a.TryExit() {
			return doneFlag, true
		}
	}
	select {
	case <-a.ctx.Done():
		doneFlag = true
		if a.q.Length() == 0 && a.TryExit() {
			return doneFlag, true
		}
	default:
	}
	return doneFlag, false
}

// OnRun service run
func (a *BaseContext) OnRun(dt time.Duration) {
	fmt.Printf("service %s do Nothing \n", a.Name())

}

// RegisterService register service name to create function mapping
func RegisterService(name string, createFunc func() (Service, error)) error {
	if _, ok := relServices[name]; ok {
		return fmt.Errorf("Already register Service %s", name)
	}
	relServices[name] = createFunc
	return nil
}

// CreateService create service by name
func CreateService(name string) (Service, error) {
	if f, ok := relServices[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("You need register Service %s first", name)
}

// GetServiceByID get service instance by id
func GetServiceByID(serviceID int32) (Service, error) {
	se, ok := App.Services.Load(serviceID)
	if ok {
		return se.(Service), nil
	}
	return nil, fmt.Errorf("cant get service ID")
}

// NewNamedService create named service
func NewNamedService(serviceID int32, name string, parent Service, m map[string]interface{}) (Service, error) {
	return createServiceByID(serviceID, name, parent, m)
}
func createServiceByID(serviceID int32, name string, parent Service, m map[string]interface{}) (Service, error) {
	se, err := CreateService(name)
	if err != nil {
		return nil, err
	}
	queueSize, ok := m["queueSize"]
	if !ok {
		return nil, errors.New("params queueSize needed")
	}
	se.MakeContext(nil, int32(queueSize.(int)))
	se.SetName(ServiceName(se))
	se.Setup(m)
	se.SetID(serviceID)
	App.Services.Store(serviceID, se)
	if serviceID == ServiceIDLog {
		GLoggerService = se
	}
	GOWithContext(ServiceRun, se)
	return se, nil
}

// NewService create anonymous service
func NewService(name string, parent Service, m map[string]interface{}) (Service, error) {
	//Inc AnonymousServiceID count = count +1
	serviceID := atomic.AddInt32(&AnonymousServiceID, 1)
	return createServiceByID(serviceID, name, parent, m)
}

// Call get info from services
func Call(destServiceID int32, cmd string, option int32, args ...interface{}) (interface{}, error) {
	var obj = NewCallObject(cmd, option, args...)
	s, err := GetServiceByID(destServiceID)
	if err != nil {
		return nil, err
	}
	if err = s.Call(option, obj); err != nil {
		return nil, err
	}
	result := <-obj.ChanRet
	return result.Ret, result.Err
}
