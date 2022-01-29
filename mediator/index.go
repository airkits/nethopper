package mediator

import (
	"fmt"
	"sync"

	"github.com/airkits/nethopper/config"
)

var instance *Mediator
var once sync.Once

//M mediator instance
func M() *Mediator {
	once.Do(func() {
		instance = NewMediator()
	})
	return instance
}

// NewModule create module
func NewModule(mid uint8, createFunc func() (IModule, error), conf config.IConfig, dep []uint8) (IModule, error) {
	if M().HasModule(mid) {
		panic(fmt.Sprintf("already exist module %d", mid))
	}
	data := NewMData(mid, createFunc, conf, dep)
	return M().CreateModule(data)
}

// AsyncCall async get data from modules,return call object
// same option value will run in same processor
func AsyncCall(destMID uint8, cmd string, option int32, args ...interface{}) (*CallObject, error) {
	m := M().GetModuleByID(destMID)
	if m != nil {
		return nil, fmt.Errorf("get module failed")
	}
	obj := NewCallObject(cmd, option, args...)
	if err := m.Call(option, obj); err != nil {
		return nil, err
	}
	return obj, nil

}

// Call sync get data from modules
// same option value will run in same processor
func Call(destMID uint8, cmd string, option int32, args ...interface{}) (interface{}, Ret) {
	obj, err := AsyncCall(destMID, cmd, option, args...)
	if err != nil {
		return nil, Ret{Code: -1, Err: err}
	}
	result := <-obj.ChanRet
	return result.Data, result.Ret
}
