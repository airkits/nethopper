package mediator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/airkits/nethopper/base"
	"github.com/airkits/nethopper/config"
	"github.com/airkits/nethopper/log"
	"github.com/airkits/nethopper/utils"
)

var instance *Mediator
var once sync.Once

// M mediator instance
func M() *Mediator {
	once.Do(func() {
		instance = NewMediator()
	})
	return instance
}

func GracefulExit() {
	M().Exit()
	M().Wait()
}
func GetAppCtx() *base.AppContext {
	return M().AppCtx
}

// NewModule create module
func NewModule(mid uint8, createFunc func() (IModule, error), conf config.IConfig, dep []uint8) (IModule, error) {
	if M().HasModule(mid) {
		panic(fmt.Sprintf("already exist module %d", mid))
	}
	data := NewMData(mid, createFunc, conf, dep)
	return M().CreateModule(data)
}

func GetModuleByID(mid uint8) IModule {
	return M().GetModuleByID(mid)
}

// AsyncCall async get data from modules,return call object
// same option value will run in same processor
func AsyncCall(destMID uint8, cmdID int32, option int32, args ...interface{}) (*base.CallObject, error) {
	m := M().GetModuleByID(destMID)
	if m == nil {
		return nil, fmt.Errorf("get module failed module [%d] cmd[%d]", destMID, cmdID)
	}
	obj := base.NewCallObject(m, cmdID, option, args...)
	obj.SetTrace(destMID)
	if err := m.Call(option, obj); err != nil {
		return nil, err
	}
	return obj, nil

}

// AsyncNotify async send data from modules,return call object
// same option value will run in same processor
func AsyncNotify(destMID uint8, cmdID int32, option int32, args ...interface{}) *base.CallObject {
	m := M().GetModuleByID(destMID)
	if m == nil {
		return nil
	}
	obj := base.NewNotifyObject(m, cmdID, option, args...)
	obj.SetTrace(destMID)
	if err := m.Call(option, obj); err != nil {
		return nil
	}
	return obj

}

// Call sync get data from modules
// same option value will run in same processor
func Call(destMID uint8, cmdID int32, option int32, args ...interface{}) *base.Ret {
	obj, err := AsyncCall(destMID, cmdID, option, args...)
	if err != nil {
		result := base.NewRet(base.ErrCodeModule, err, nil)
		result.SetTrace(destMID)
		return result
	}
	result := <-obj.ChanRet
	result.SetTrace(destMID)
	return result
}

// NewWorkerPool create Processor pool
func NewWorkerPool(cap uint32, queueSize uint32, expired time.Duration) (IWorkerPool, error) {
	if cap == 0 {
		return nil, ErrInvalidcapacity
	}

	// create Processor pool
	p := &WorkerPool{
		capacity:        cap,
		expiredDuration: expired,
		workers:         make([]*Processor, 0, cap),
	}

	p.Setup(queueSize)

	go p.ExpiredCleaning()

	return p, nil
}

// NewFixedWorkerPool create fixed Processor pool
func NewFixedWorkerPool(cap uint32, queueSize uint32, expired time.Duration) (IWorkerPool, error) {
	if cap == 0 {
		return nil, ErrInvalidcapacity
	}
	capacity, power := utils.PowerCalc(int32(cap))
	// create FixedProcessor pool
	p := &FixedWorkerPool{
		capacity:        uint32(capacity),
		expiredDuration: expired,
		workers:         make([]*Processor, capacity, capacity),
		power:           power,
	}

	p.Setup(queueSize)
	go p.ExpiredCleaning()

	return p, nil
}

// ModuleName get the module name
func ModuleName(s IModule) string {
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

func ExecuteHandler(s IModule, obj *base.CallObject) *base.Ret {
	var result *base.Ret
	f := s.GetHandler(obj.CmdID)
	if f != nil {
		switch f.(type) {
		case func(interface{}) *base.Ret:
			result = f.(func(interface{}) *base.Ret)(s)
		case func(interface{}, interface{}) *base.Ret:
			result = f.(func(interface{}, interface{}) *base.Ret)(s, obj.Args[0])
		case func(interface{}, interface{}, interface{}) *base.Ret:
			result = f.(func(interface{}, interface{}, interface{}) *base.Ret)(s, obj.Args[0], obj.Args[1])
		case func(interface{}, interface{}, interface{}, interface{}) *base.Ret:
			result = f.(func(interface{}, interface{}, interface{}, interface{}) *base.Ret)(s, obj.Args[0], obj.Args[1], obj.Args[2])
		default:
			panic(fmt.Sprintf("function cmd %v: definition of function is invalid,%v", obj.CmdID, reflect.TypeOf(f)))
		}

	} else {
		f = s.GetReflectHandler(obj.CmdID)
		if f == nil {
			err := fmt.Errorf("module[%s],handler id %v: function not registered", s.Name(), obj.CmdID)
			panic(err)
		} else {
			args := []interface{}{s}
			args = append(args, obj.Args...)
			values := base.CallFunction(f, args...)
			if values == nil {
				err := errors.New("unsupport handler,need return (interface{},Result) or ([]interface{},Result)")
				panic(err)
			} else {
				l := len(values)
				if l == 1 {
					result = values[0].Interface().(*base.Ret)
				} else {
					err := errors.New("unsupport params length")
					result = base.NewRet(base.ErrCodeModule, err, nil)
					panic(err)
				}
			}
		}
	}
	return result
}

// RunSimpleFrame wrapper simple run function
func RunSimpleFrame(s IModule) {
	m, err := s.MQ().Pop()
	if err != nil {
		return
	}
	obj := m.(*base.CallObject)
	// if err := s.DoWorker(obj); err != nil {
	// 	log.Error("%s error %s", s.Name(), err.Error())
	// }

	if !s.HasWorkerPool() {
		//err = errors.New("no processor pool")
		result := s.Execute(obj)
		if !obj.Notify {
			obj.ChanRet <- result
		}
		return
	}
	err = s.WorkerPoolSubmit(obj)

	if err != nil && !obj.Notify {

		obj.ChanRet <- base.NewRet(base.ErrCodeWorker, err, nil)
	}
}

// ModuleRun wrapper module goruntine and in an orderly way to exit
func ModuleRun(s IModule) {
	ctxDone := false
	exitFlag := false
	start := time.Now()
	log.Info("Module [%s] starting", s.Name())
	for {
		s.OnRun(time.Since(start))

		if ctxDone, exitFlag = s.CanExit(ctxDone); exitFlag {
			fmt.Printf("module exit %s\n", s.Name())
			return
		}

		//start = time.Now()
		//if s.MQ().Length() == 0 {
		// t := time.Duration(s.IdleTimes()) * time.Nanosecond
		// time.Sleep(t)
		// s.IdleTimesAdd()

		//}
		//runtime.Gosched()
	}
}
