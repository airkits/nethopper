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
	// ModuleNamedIDs module id define, system reserved 1-63
	ModuleNamedIDs = iota
	// ModuleIDMain main goruntinue
	ModuleIDMain
	// ModuleIDMonitor server monitor module
	ModuleIDMonitor
	// ModuleIDLog log module
	ModuleIDLog
	// ModuleIDTCP tcp module
	ModuleIDTCP
	// ModuleIDKCP kcp module
	ModuleIDKCP
	// ModuleIDHTTP http module
	ModuleIDHTTP
	// ModuleIDLogic logic module
	ModuleIDLogic
	// ModuleIDRedis redis module
	ModuleIDRedis
	// ModuleIDTCPClient tcp client module
	ModuleIDTCPClient
	// ModuleIDKCPClient kcp client module
	ModuleIDKCPClient
	// ModuleIDHTTPClient http client module
	ModuleIDHTTPClient
	// ModuleIDDB common db module
	ModuleIDDB
	// ModuleIDUserCustom User custom define named modules from 64-128
	ModuleIDUserCustom = 64
	// ModuleIDNamedMax named modules max ID
	ModuleIDNamedMax = 128
)

// Module interface define
type Module interface {
	//BaseContext start
	// ID module id
	ID() int32
	//SetID set module ID
	SetID(v int32)
	// Name module name
	Name() string
	//SetName set module name
	SetName(v string)

	// MakeContext init base module queue and create context
	MakeContext(p Module, queueSize int32)
	// Context get module context
	Context() context.Context
	// ChildAdd after child module created and tell parent module, ref count +1
	ChildAdd()
	// ChildDone child module exit and tell parent module, ref count -1
	ChildDone()
	// Close call context cancel ,self and all child module will receive context.Done()
	Close()
	// Queue return module queue
	MQ() queue.Queue
	// CanExit if receive ctx.Done() and child ref = 0 and queue is empty ,then return true
	CanExit(doneflag bool) (bool, bool)
	// TryExit check child ref count , if ref count == 0 then return true, if parent not nil, fire parent.ChildDone()
	TryExit() bool
	//BaseContext end

	// UserData module custom option, can you store you data and you must keep goruntine safe
	UserData() int32
	// Setup init custom module and pass config map to module
	Setup(m map[string]interface{}) (Module, error)
	//Reload reload config
	Reload(m map[string]interface{}) error
	// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
	OnRun(dt time.Duration)
	// Stop goruntine
	Stop() error
	// Call async send callobject to module
	Call(option int32, obj *CallObject) error
	// PushBytes async send string or bytes to queue
	PushBytes(option int32, buf []byte) error
	//GetHandler get call handler
	GetHandler(id interface{}) interface{}
	// Processor process callobject
	Processor(obj *CallObject) error
}

// RunSimpleFrame wrapper simple run function
func RunSimpleFrame(s Module) {
	for i := 0; i < 128; i++ {
		m, err := s.MQ().AsyncPop()
		if err != nil {
			break
		}
		obj := m.(*CallObject)

		if err := s.Processor(obj); err != nil {
			Error("%s error %s", s.Name(), err.Error())
			break
		}
	}
}

// ModuleRun wrapper module goruntine and in an orderly way to exit
func ModuleRun(s Module) {
	ctxDone := false
	exitFlag := false
	start := time.Now()
	Info("Module %s start ", s.Name())
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

// ModuleName get the module name
func ModuleName(s Module) string {
	t := reflect.TypeOf(s)
	return t.Elem().Name()
}

//BaseContext use context to close all module and using the bubbling method to exit
type BaseContext struct {
	ctx        context.Context
	cancel     context.CancelFunc
	parent     Module
	childRef   int32
	q          queue.Queue
	name       string
	id         int32
	functions  map[interface{}]interface{}
	processers IWorkerPool
}

// RegisterHandler register function before run
func (a *BaseContext) RegisterHandler(id interface{}, f interface{}) {

	// switch f.(type) {
	// case func(Module, *CallObject, string) (string, error):
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

// MakeContext init base module queue and create context
func (a *BaseContext) MakeContext(p Module, queueSize int32) {
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

// Call async send message to module
func (a *BaseContext) Call(option int32, obj *CallObject) error {
	if err := a.q.AsyncPush(obj); err != nil {
		Error(err.Error())
	}
	return nil
}

// CreateWorkerPool create processor pool
func (a *BaseContext) CreateWorkerPool(s Module, cap uint32, expired time.Duration, isNonBlocking bool) (err error) {
	if a.processers, err = NewFixedWorkerPool(s, cap, expired); err != nil {
		return err
	}
	return nil
}

// MQ return module queue
func (a *BaseContext) MQ() queue.Queue {
	return a.q
}

// Context get module context
func (a *BaseContext) Context() context.Context {
	return a.ctx
}

// ChildAdd child module created and tell parent module, ref count +1
func (a *BaseContext) ChildAdd() {
	atomic.AddInt32(&a.childRef, 1)
}

// ChildDone child module exit and tell parent module, ref count -1
func (a *BaseContext) ChildDone() {
	atomic.AddInt32(&a.childRef, -1)
}

// Close call context cancel ,self and all child module will receive context.Done()
func (a *BaseContext) Close() {
	a.cancel()
}

//ID module ID
func (a *BaseContext) ID() int32 {
	return a.id
}

//SetID set module id
func (a *BaseContext) SetID(v int32) {
	a.id = v
}

//Name module name
func (a *BaseContext) Name() string {
	return a.name
}

//SetName set module name
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

// OnRun module run
func (a *BaseContext) OnRun(dt time.Duration) {
	fmt.Printf("module %s do Nothing \n", a.Name())

}

// RegisterModule register module name to create function mapping
func RegisterModule(name string, createFunc func() (Module, error)) error {
	if _, ok := relModules[name]; ok {
		return fmt.Errorf("Already register Module %s", name)
	}
	relModules[name] = createFunc
	return nil
}

// CreateModule create module by name
func CreateModule(name string) (Module, error) {
	if f, ok := relModules[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("You need register Module %s first", name)
}

// GetModuleByID get module instance by id
func GetModuleByID(moduleID int32) (Module, error) {
	se, ok := App.Modules.Load(moduleID)
	if ok {
		return se.(Module), nil
	}
	return nil, fmt.Errorf("cant get module ID")
}

// NewNamedModule create named module
func NewNamedModule(moduleID int32, name string, parent Module, m map[string]interface{}) (Module, error) {
	return createModuleByID(moduleID, name, parent, m)
}
func createModuleByID(moduleID int32, name string, parent Module, m map[string]interface{}) (Module, error) {
	se, err := CreateModule(name)
	if err != nil {
		return nil, err
	}
	queueSize, ok := m["queueSize"]
	if !ok {
		return nil, errors.New("params queueSize needed")
	}
	se.MakeContext(nil, int32(queueSize.(int)))
	se.SetName(ModuleName(se))
	se.Setup(m)
	se.SetID(moduleID)
	App.Modules.Store(moduleID, se)
	if moduleID == ModuleIDLog {
		GLoggerModule = se
	}
	GOWithContext(ModuleRun, se)
	return se, nil
}

// NewModule create anonymous module
func NewModule(name string, parent Module, m map[string]interface{}) (Module, error) {
	//Inc AnonymousModuleID count = count +1
	moduleID := atomic.AddInt32(&AnonymousModuleID, 1)
	return createModuleByID(moduleID, name, parent, m)
}

// Call get info from modules
func Call(destModuleID int32, cmd string, option int32, args ...interface{}) (interface{}, error) {
	var obj = NewCallObject(cmd, option, args...)
	s, err := GetModuleByID(destModuleID)
	if err != nil {
		return nil, err
	}
	if err = s.Call(option, obj); err != nil {
		return nil, err
	}
	result := <-obj.ChanRet
	return result.Ret, result.Err
}
