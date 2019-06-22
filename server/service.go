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
	"sync/atomic"

	"github.com/gonethopper/queue"
)

const (
	// ServiceNamedIDs service id define, system reserved 1-63
	ServiceNamedIDs = iota
	// ServiceIDMain main goruntinue
	ServiceIDMain
	// ServiceIDMonitor server monitor service
	ServiceIDMonitor
	//ServiceIDLog log service
	ServiceIDLog
	//ServiceIDUserCustom User custom define named services from 64-128
	ServiceIDUserCustom = 64
	//ServiceIDNamedMax named services max ID
	ServiceIDNamedMax = 128
)

// Service interface define
type Service interface {
	//BaseService start
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
	Queue() queue.Queue
	// CanExit if receive ctx.Done() and child ref = 0 and queue is empty ,then return true
	CanExit(doneflag bool) (bool, bool)
	// TryExit check child ref count , if ref count == 0 then return true, if parent not nil, fire parent.ChildDone()
	TryExit() bool
	//BaseService end

	// UserData service custom option, can you store you data and you must keep goruntine safe
	UserData() int32
	// Setup init custom service and pass config map to service
	Setup(m map[string]interface{}) (Service, error)
	//Reload reload config
	Reload(m map[string]interface{}) error
	// Run create goruntine and run, always use ServiceRun to call this function
	Run()
	// Stop goruntine
	Stop() error
	// SendMessage async send message to service
	SendMessage(option int32, msg *Message) error
	// SendBytes async send string or bytes to queue
	SendBytes(option int32, buf []byte) error
}

// ServiceRun wrapper service goruntine and in an orderly way to exit
func ServiceRun(s Service) {
	ctxDone := false
	exitFlag := false
	for {
		s.Run()
		if ctxDone, exitFlag = s.CanExit(ctxDone); exitFlag {
			return
		}
	}
}

// ServiceName get the service name
func ServiceName(s Service) string {
	t := reflect.TypeOf(s)
	return t.Elem().Name()
}

//BaseService use context to close all service and using the bubbling method to exit
type BaseService struct {
	ctx      context.Context
	cancel   context.CancelFunc
	parent   Service
	childRef int32
	q        queue.Queue
	name     string
	id       int32
}

// MakeContext init base service queue and create context
func (a *BaseService) MakeContext(p Service, queueSize int32) {
	a.parent = p
	a.q = queue.NewChanQueue(queueSize)
	if p == nil {
		a.ctx, a.cancel = context.WithCancel(context.Background())
	} else {
		a.ctx, a.cancel = context.WithCancel(p.Context())
		p.ChildAdd()
	}
}

// Queue return service queue
func (a *BaseService) Queue() queue.Queue {
	return a.q
}

// Context get service context
func (a *BaseService) Context() context.Context {
	return a.ctx
}

// ChildAdd child service created and tell parent service, ref count +1
func (a *BaseService) ChildAdd() {
	atomic.AddInt32(&a.childRef, 1)
}

// ChildDone child service exit and tell parent service, ref count -1
func (a *BaseService) ChildDone() {
	atomic.AddInt32(&a.childRef, -1)
}

// Close call context cancel ,self and all child service will receive context.Done()
func (a *BaseService) Close() {
	a.cancel()
}

//ID service ID
func (a *BaseService) ID() int32 {
	return a.id
}

//SetID set service id
func (a *BaseService) SetID(v int32) {
	a.id = v
}

//Name service name
func (a *BaseService) Name() string {
	return a.name
}

//SetName set service name
func (a *BaseService) SetName(v string) {
	a.name = v
}

// TryExit check child ref count , if ref count == 0 then return true, if parent not nil, and will fire parent.ChildDone()
func (a *BaseService) TryExit() bool {

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
func (a *BaseService) CanExit(doneFlag bool) (bool, bool) {
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

// Run service run
func (a *BaseService) Run() {
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

// DeleteService unregister service
func DeleteService(serviceID int32) error {
	se, err := GetServiceByID(serviceID)
	if err != nil {
		return err
	}
	App.Services.Delete(serviceID)
	se.Stop()

	return nil
}

// DeleteAllServices traversing services
func DeleteAllServices() {
	App.Services.Range(func(key interface{}, v interface{}) bool {
		DeleteService(key.(int32))
		return true
	})
}

// SendMessage send message to services
func SendMessage(serviceID int32, option int32, msg *Message) error {
	s, err := GetServiceByID(serviceID)
	if err != nil {
		return err
	}
	return s.SendMessage(option, msg)
}
