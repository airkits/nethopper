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
	"strings"
	"sync/atomic"
	"time"

	"github.com/gonethopper/nethopper/base/queue"
)

const (
	// ModuleNamedIDs module id define, system reserved 1-63
	ModuleNamedIDs = iota
	// MIDMain main goruntinue
	MIDMain
	// MIDMonitor server monitor module
	MIDMonitor
	// MIDLog log module
	MIDLog
	// MIDTCP tcp module
	MIDTCP
	// MIDKCP kcp module
	MIDKCP
	// MIDQUIC quic module
	MIDQUIC
	// MIDWSServer ws server
	MIDWSServer
	// MIDGRPCServer grpc server
	MIDGRPCServer
	// MIDHTTP http module
	MIDHTTP
	// MIDLogic logic module
	MIDLogic
	// MIDRedis redis module
	MIDRedis
	// MIDTCPClient tcp client module
	MIDTCPClient
	// MIDKCPClient kcp client module
	MIDKCPClient
	// MIDQUICClient quic client module
	MIDQUICClient
	// MIDHTTPClient http client module
	MIDHTTPClient
	// MIDGRPCClient grpc client module
	MIDGRPCClient
	// MIDWSClient ws client
	MIDWSClient
	// MIDDB common db module
	MIDDB

	// MIDUserCustom User custom define named modules from 64-128
	MIDUserCustom = 64
	// MIDNamedMax named modules max ID
	MIDNamedMax = 128
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

	//Handlers set moudle handlers
	Handlers() map[string]interface{}
	//ReflectHandlers set moudle reflect handlers
	ReflectHandlers() map[string]interface{}

	// RegisterHandler register function before run
	RegisterHandler(id interface{}, f interface{})
	// RegisterReflectHandler register reflect function before run
	RegisterReflectHandler(id interface{}, f interface{})

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
	Setup(conf IConfig) (Module, error)
	//Reload reload config
	Reload(conf IConfig) error
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

	//GetReflectHandler get reflect handler
	GetReflectHandler(id interface{}) interface{}

	// Processor process callobject
	Processor(obj *CallObject) error

	// IdleTimesReset reset idle times
	IdleTimesReset()

	// IdleTimes get idle times
	IdleTimes() uint32

	// IdleTimesAdd add idle times
	IdleTimesAdd()
}

// RunSimpleFrame wrapper simple run function
func RunSimpleFrame(s Module, packageSize int) {
	for i := 0; i < packageSize; i++ {
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
	Info("Module [%s] starting", s.Name())
	for {
		s.OnRun(time.Since(start))

		if ctxDone, exitFlag = s.CanExit(ctxDone); exitFlag {
			return
		}

		start = time.Now()
		if s.MQ().Length() == 0 {
			t := time.Duration(s.IdleTimes()) * time.Nanosecond
			time.Sleep(t)
			s.IdleTimesAdd()

		}
		runtime.Gosched()
	}
}

// ModuleName get the module name
func ModuleName(s Module) string {
	t := reflect.TypeOf(s)
	path := t.Elem().PkgPath()
	pos := strings.LastIndex(path, "/")
	if pos >= 0 {
		prefix := []byte(path)[pos+1 : len(path)]
		rs := string(prefix)
		return rs
	}
	return "unknown module"
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
	funcs      map[interface{}]interface{} //handlers
	rfuncs     map[interface{}]interface{} //reflect handlers
	processers IWorkerPool
	idleTimes  uint32
}

//Handlers set moudle handlers
func (s *BaseContext) Handlers() map[string]interface{} {
	return nil
}

//ReflectHandlers set moudle reflect handlers
func (s *BaseContext) ReflectHandlers() map[string]interface{} {
	return nil
}

// RegisterHandler register function before run
func (s *BaseContext) RegisterHandler(id interface{}, f interface{}) {

	// switch f.(type) {
	// case func(Module, *CallObject, string) (string, error):
	// default:
	// 	panic(fmt.Sprintf("function id %v: definition of function is invalid,%v", id, reflect.TypeOf(f)))
	// }

	if _, ok := s.funcs[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	s.funcs[id] = f
}

// RegisterReflectHandler register reflect function before run
func (s *BaseContext) RegisterReflectHandler(id interface{}, f interface{}) {

	// switch f.(type) {
	// case func(Module, *CallObject, string) (string, error):
	// default:
	// 	panic(fmt.Sprintf("function id %v: definition of function is invalid,%v", id, reflect.TypeOf(f)))
	// }

	if _, ok := s.rfuncs[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	s.rfuncs[id] = f
}

// GetHandler get call handler
func (s *BaseContext) GetHandler(id interface{}) interface{} {
	return s.funcs[id]
}

// GetReflectHandler get call reflect handler
func (s *BaseContext) GetReflectHandler(id interface{}) interface{} {
	return s.rfuncs[id]
}

// IdleTimesReset reset idle times
func (s *BaseContext) IdleTimesReset() {
	atomic.StoreUint32(&s.idleTimes, 500)
}

// IdleTimes get idle times
func (s *BaseContext) IdleTimes() uint32 {
	return atomic.LoadUint32(&s.idleTimes)
}

// IdleTimesAdd add idle times
func (s *BaseContext) IdleTimesAdd() {
	t := s.IdleTimes()
	if t >= 20000000 { //2s
		return
	}
	atomic.AddUint32(&s.idleTimes, 100)
}

// MakeContext init base module queue and create context
func (s *BaseContext) MakeContext(p Module, queueSize int32) {
	s.parent = p
	s.q = queue.NewChanQueue(queueSize)
	s.funcs = make(map[interface{}]interface{})
	s.rfuncs = make(map[interface{}]interface{})
	if p == nil {
		s.ctx, s.cancel = context.WithCancel(context.Background())
	} else {
		s.ctx, s.cancel = context.WithCancel(p.Context())
		p.ChildAdd()
	}

}

// Processor process callobject
func (s *BaseContext) Processor(obj *CallObject) error {
	Debug("[%s] cmd [%s] process", s.Name(), obj.Cmd)
	var err error
	if s.processers == nil {
		err = errors.New("no processor pool")
	} else {
		err = s.processers.Submit(obj)
	}
	if err != nil {
		obj.ChanRet <- RetObject{
			Ret: nil,
			Result: Result{
				Code: -1,
				Err:  err,
			},
		}
	}
	return err
}

// Call async send message to module
func (s *BaseContext) Call(option int32, obj *CallObject) error {
	s.IdleTimesReset()
	if err := s.q.AsyncPush(obj); err != nil {
		Error(err.Error())
	}
	return nil
}

// CreateWorkerPool create processor pool
func (s *BaseContext) CreateWorkerPool(m Module, cap uint32, expired time.Duration, isNonBlocking bool) (err error) {
	if s.processers, err = NewFixedWorkerPool(m, cap, expired); err != nil {
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

// ChildAdd child module created and tell parent module, ref count +1
func (s *BaseContext) ChildAdd() {
	atomic.AddInt32(&s.childRef, 1)
}

// ChildDone child module exit and tell parent module, ref count -1
func (s *BaseContext) ChildDone() {
	atomic.AddInt32(&s.childRef, -1)
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

// TryExit check child ref count , if ref count == 0 then return true, if parent not nil, and will fire parent.ChildDone()
func (s *BaseContext) TryExit() bool {

	count := atomic.LoadInt32(&s.childRef)
	if count > 0 {
		return false
	}
	if s.parent != nil {
		s.parent.ChildDone()
	}
	return true
}

// CanExit if receive ctx.Done() and all child exit and queue is empty ,then return true
func (s *BaseContext) CanExit(doneFlag bool) (bool, bool) {
	if doneFlag {
		if s.q.Length() == 0 && s.TryExit() {
			return doneFlag, true
		}
	}
	select {
	case <-s.ctx.Done():
		doneFlag = true
		if s.q.Length() == 0 && s.TryExit() {
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
func (s *BaseContext) ReadConfig(conf IConfig) error {
	return nil
}

//Reload reload config
func (s *BaseContext) Reload(conf IConfig) error {
	return nil
}

// Stop goruntine
func (s *BaseContext) Stop() error {

	return nil
}

//to override end

// RegisterModule register module name to create function mapping
func RegisterModule(name string, createFunc func() (Module, error)) error {
	if IsModuleRegistered(name) {
		return fmt.Errorf("Already register Module %s", name)
	}
	relModules[name] = createFunc
	return nil
}

//IsModuleRegistered check module is registered
func IsModuleRegistered(name string) bool {
	if _, ok := relModules[name]; ok {
		return true
	}
	return false
}

// CreateModule create module by name
func CreateModule(name string) (Module, error) {
	if f, ok := relModules[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("You need register Module %s first", name)
}

// GetModuleByID get module instance by id
func GetModuleByID(MID int32) (Module, error) {
	se, ok := App.Modules.Load(MID)
	if ok {
		return se.(Module), nil
	}
	return nil, fmt.Errorf("cant get module ID %d", MID)
}

// NewNamedModule create named module
func NewNamedModule(MID int32, name string, createFunc func() (Module, error), parent Module, conf IConfig) (Module, error) {
	if !IsModuleRegistered(name) {
		if err := RegisterModule(name, createFunc); err != nil {
			panic(err)
		}
	}
	return createModuleByID(MID, name, parent, conf)
}

func cmdRegister(s Module) {
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
func createModuleByID(MID int32, name string, parent Module, conf IConfig) (Module, error) {
	se, err := CreateModule(name)
	if err != nil {
		return nil, err
	}
	se.MakeContext(nil, int32(conf.GetQueueSize()))
	se.SetName(ModuleName(se))
	cmdRegister(se)
	se.Setup(conf)
	se.SetID(MID)
	App.Modules.Store(MID, se)
	if MID == MIDLog {
		GLoggerModule = se
	}
	GOWithContext(ModuleRun, se)
	return se, nil
}

// NewModule create anonymous module
func NewModule(name string, parent Module, conf IConfig) (Module, error) {
	//Inc AnonymousMID count = count +1
	MID := atomic.AddInt32(&AnonymousMID, 1)
	return createModuleByID(MID, name, parent, conf)
}

// Call get info from modules
// same option will run in same processor
func Call(destMID int32, cmd string, option int32, args ...interface{}) (interface{}, Result) {
	var obj = NewCallObject(cmd, option, args...)
	m, err := GetModuleByID(destMID)
	if err != nil {
		return nil, Result{Code: 0, Err: err}
	}
	if err = m.Call(option, obj); err != nil {
		return nil, Result{Code: 0, Err: err}
	}
	result := <-obj.ChanRet
	return result.Ret, result.Result
}
