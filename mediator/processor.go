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
// * @Date: 2019-12-06 08:28:07
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-06 08:28:07

package mediator

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/airkits/nethopper/base"
	"github.com/airkits/nethopper/base/queue"
)

// NewProcessor create new processor
func NewProcessor(owner IWorkerPool, queueSize uint32) *Processor {
	return &Processor{
		owner:   owner,
		q:       queue.NewChanQueue(int32(queueSize)),
		timeout: time.Now(),
	}
}

// Processor process job
type Processor struct {
	owner IWorkerPool
	//CallObject chan
	q queue.Queue
	//timeout set to tigger timeout event
	timeout time.Time
}

// Process goruntine process pre call
func Process(s IModule, obj *CallObject) (result Ret) {
	//defer TraceCost(s.Name() + ":" + obj.Cmd)()
	var ret = RetObject{
		Data: nil,
		Ret:  Ret{Err: nil, Code: 0},
	}

	defer func() {
		if r := recover(); r != nil {
			result.Err = r.(error)
			result.Code = -1
		}
	}()

	f := s.(IModule).GetHandler(obj.Cmd)
	if f != nil {
		var data interface{}
		var result Ret
		switch f.(type) {
		case func(interface{}) (interface{}, Ret):
			data, result = f.(func(interface{}) (interface{}, Ret))(s)
		case func(interface{}, interface{}) (interface{}, Ret):
			data, result = f.(func(interface{}, interface{}) (interface{}, Ret))(s, obj.Args[0])
		case func(interface{}, interface{}, interface{}) (interface{}, Ret):
			data, result = f.(func(interface{}, interface{}, interface{}) (interface{}, Ret))(s, obj.Args[0], obj.Args[1])
		case func(interface{}, interface{}, interface{}, interface{}) (interface{}, Ret):
			data, result = f.(func(interface{}, interface{}, interface{}, interface{}) (interface{}, Ret))(s, obj.Args[0], obj.Args[1], obj.Args[2])
		default:
			panic(fmt.Sprintf("function cmd %v: definition of function is invalid,%v", obj.Cmd, reflect.TypeOf(f)))
		}
		ret.Data = data
		ret.Ret = result
	} else {
		f = s.(IModule).GetReflectHandler(obj.Cmd)
		if f == nil {
			err := fmt.Errorf("module[%s],handler id %v: function not registered", s.Name(), obj.Cmd)
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
				if l == 2 {
					ret.Data = values[0].Interface()
					if values[1].Interface() != nil {
						result = values[1].Interface().(Ret)
						ret.Ret.Code = result.Code
						ret.Ret.Err = result.Err
					}
				} else {
					err := errors.New("unsupport params length")
					panic(err)
				}
			}
		}
	}

	obj.ChanRet <- ret
	return result
}

// Run Processor goruntine
func (w *Processor) Run() {
	//	atomic.AddUint32(&w.owner.WorkerCount, 1)
	w.owner.WorkerCountInc()
	go func() {
		for {
			obj, err := w.q.Pop()
			if err == nil && obj == nil {
				//atomic.AddUint32(&w.owner.WorkerCount, ^uint32(-(-1)-1))
				w.owner.WorkerCountDec()
				w.owner.CachePut(w)
				break
			}
			if err == nil {
				if result := Process(w.owner.Owner(), obj.(*CallObject)); result.Err != nil {
					obj.(*CallObject).ChanRet <- RetObject{Data: nil, Ret: result}
				}
			}
			if w.q.Length() == 0 {
				if ok := w.owner.RecycleProcessor(w); !ok {
					break
				}
			}
		}
	}()
}

//Submit task to processor
func (w *Processor) Submit(obj *CallObject) error {
	if err := w.q.AsyncPush(obj); err != nil {
		return err
	}
	return nil
}
