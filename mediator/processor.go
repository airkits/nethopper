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
	"time"

	"github.com/airkits/nethopper/base"
	"github.com/airkits/nethopper/base/queue"
)

// NewProcessor create new processor
func NewProcessor(wp IWorkerPool, queueSize uint32) *Processor {
	return &Processor{
		wp:      wp,
		q:       queue.NewChanQueue(int32(queueSize)),
		timeout: time.Now(),
	}
}

// Processor process job
type Processor struct {
	wp IWorkerPool
	//CallObject chan
	q queue.Queue
	//timeout set to tigger timeout event
	timeout time.Time
}

// Process goruntine process pre call
func Process(obj *base.CallObject) *base.Ret {
	//defer TraceCost(s.Name() + ":" + obj.Cmd)()
	var result *base.Ret

	defer func() {
		if r := recover(); r != nil {
			result = base.NewRet(base.ErrCodeProcessor, r.(error), nil)
			if obj.Type == base.CallObejctNone {
				obj.ChanRet <- result
			}

		}
	}()

	result = obj.Caller.Execute(obj)
	if obj.Type == base.CallObejctNone {
		obj.ChanRet <- result
	}
	return result
}

// f := s.(IModule).GetHandler(obj.CmdID)
// if f != nil {

// 	switch f.(type) {
// 	case func(interface{}) *Ret:
// 		result = f.(func(interface{}) *Ret)(s)
// 	case func(interface{}, interface{}) *Ret:
// 		result = f.(func(interface{}, interface{}) *Ret)(s, obj.Args[0])
// 	case func(interface{}, interface{}, interface{}) *Ret:
// 		result = f.(func(interface{}, interface{}, interface{}) *Ret)(s, obj.Args[0], obj.Args[1])
// 	case func(interface{}, interface{}, interface{}, interface{}) *Ret:
// 		result = f.(func(interface{}, interface{}, interface{}, interface{}) *Ret)(s, obj.Args[0], obj.Args[1], obj.Args[2])
// 	default:
// 		panic(fmt.Sprintf("function cmd %v: definition of function is invalid,%v", obj.CmdID, reflect.TypeOf(f)))
// 	}

// } else {
// 	f = s.(IModule).GetReflectHandler(obj.CmdID)
// 	if f == nil {
// 		err := fmt.Errorf("module[%s],handler id %v: function not registered", s.Name(), obj.CmdID)
// 		panic(err)
// 	} else {
// 		args := []interface{}{s}
// 		args = append(args, obj.Args...)
// 		values := base.CallFunction(f, args...)
// 		if values == nil {
// 			err := errors.New("unsupport handler,need return (interface{},Result) or ([]interface{},Result)")
// 			panic(err)
// 		} else {
// 			l := len(values)
// 			if l == 1 {
// 				result = values[0].Interface().(*Ret)
// 			} else {
// 				err := errors.New("unsupport params length")
// 				result = NewRet(-1, err, nil)
// 				panic(err)
// 			}
// 		}
// 	}
// }

// Run Processor goruntine
func (w *Processor) Run() {
	//	atomic.AddUint32(&w.owner.WorkerCount, 1)
	w.wp.AddRef()
	go func() {
		for {
			obj, err := w.q.Pop()
			if err == nil && obj == nil {
				//atomic.AddUint32(&w.owner.WorkerCount, ^uint32(-(-1)-1))
				w.wp.DecRef()
				w.wp.CachePut(w)
				break
			}
			if err == nil {
				Process(obj.(*base.CallObject))

			}
			if w.q.Length() == 0 {
				if ok := w.wp.RecycleProcessor(w); !ok {
					break
				}
			}
		}
	}()
}

// Submit task to processor
func (w *Processor) Submit(obj *base.CallObject) error {
	if err := w.q.AsyncPush(obj); err != nil {
		return err
	}

	return nil
}
