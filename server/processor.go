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

package server

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gonethopper/nethopper/base/queue"
)

const (
	// DefaultTimeout Processor default timeout
	DefaultTimeout = 10
	// CLOSED status
	CLOSED = 1
)

var (
	//ErrWorkerPoolClosed Processor pool is release
	ErrWorkerPoolClosed = errors.New("Processor pool already closed")
	//ErrWorkerPoolBusy Processor pool is busy
	ErrWorkerPoolBusy = errors.New("Processor pool is busy,please try again")
	//ErrInvalidcapacity set invalid capacity
	ErrInvalidcapacity = errors.New("invalid capacity")
	//ErrorTodo todo error
	ErrorTodo = errors.New("todo,not implementation")
)

// NewWorkerPool create Processor pool
func NewWorkerPool(owner Module, cap uint32, expired time.Duration) (*WorkerPool, error) {
	if cap == 0 {
		return nil, ErrInvalidcapacity
	}

	// create Processor pool
	p := &WorkerPool{
		capacity:        cap,
		expiredDuration: expired,
		workers:         make([]*Processor, 0, cap),
		owner:           owner,
		name:            owner.Name(),
	}

	//bind signal and lock
	p.cond = sync.NewCond(&p.lock)
	p.cache = sync.Pool{
		New: func() interface{} {
			return NewProcessor(p, 128)
		},
	}
	go p.ExpiredCleaning()

	return p, nil
}

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
func Process(s Module, obj *CallObject) (result Ret) {
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

	f := s.(Module).GetHandler(obj.Cmd)
	if f != nil {
		switch f.(type) {
		case func(interface{}) (interface{}, Ret):
			data, result := f.(func(interface{}) (interface{}, Ret))(s)
			ret.Data = data
			ret.Ret = result
		case func(interface{}, interface{}) (interface{}, Ret):
			data, result := f.(func(interface{}, interface{}) (interface{}, Ret))(s, obj.Args[0])
			ret.Data = data
			ret.Ret = result
		case func(interface{}, interface{}, interface{}) (interface{}, Ret):
			data, result := f.(func(interface{}, interface{}, interface{}) (interface{}, Ret))(s, obj.Args[0], obj.Args[1])
			ret.Data = data
			ret.Ret = result
		case func(interface{}, interface{}, interface{}, interface{}) (interface{}, Ret):
			data, result := f.(func(interface{}, interface{}, interface{}, interface{}) (interface{}, Ret))(s, obj.Args[0], obj.Args[1], obj.Args[2])
			ret.Data = data
			ret.Ret = result
		default:
			panic(fmt.Sprintf("function cmd %v: definition of function is invalid,%v", obj.Cmd, reflect.TypeOf(f)))
		}
	} else {
		f = s.(Module).GetReflectHandler(obj.Cmd)
		if f == nil {
			err := fmt.Errorf("module[%s],handler id %v: function not registered", s.Name(), obj.Cmd)
			panic(err)
		} else {
			args := []interface{}{s}
			args = append(args, obj.Args...)
			values := CallUserFunc(f, args...)
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

// IWorkerPool process pool interface
type IWorkerPool interface {

	// Owner get the module who own the processor pool
	Owner() Module

	//WorkerCountInc current goruntine count +1
	WorkerCountInc()

	//WorkerCountDec current goruntine count -1
	WorkerCountDec()

	//CachePut put processor return to pool
	CachePut(w *Processor)

	// Submit add obj to Processor
	Submit(obj *CallObject) error

	//RecycleProcessor return back Processor to pool
	RecycleProcessor(w *Processor) bool

	// Release WorkerPool remove all Processors
	Release()
}

// WorkerPool Processor pool limit goruntine max count
// -- dynamic processor pool
type WorkerPool struct {
	// capacity max goruntine count
	capacity uint32

	//workerCount current goruntine Processors
	workerCount uint32

	//Cache sync pool to store Processor
	cache sync.Pool

	// 当关闭该Pool支持通知所有Processor退出运行以防goroutine泄露
	isClosed uint32

	// expiredDuration set timeout for Processor
	expiredDuration time.Duration

	// 互斥锁
	lock sync.Mutex

	// 信号量
	cond *sync.Cond

	// 确保关闭操作只执行一次
	once sync.Once
	name string
	// workers list
	workers []*Processor
	owner   Module
	// fixed pool flag
	fixed bool
}

// Name get the processor pool name
func (p *WorkerPool) Name() string {
	return p.name
}

// RecycleProcessor return back Processor to pool
func (p *WorkerPool) RecycleProcessor(w *Processor) bool {
	if atomic.LoadUint32(&p.isClosed) == CLOSED {
		return false
	}
	w.timeout = time.Now()
	p.lock.Lock()
	p.workers = append(p.workers, w)
	p.cond.Signal()
	p.lock.Unlock()
	return true
}

// Count get current running Processors count
func (p *WorkerPool) Count() uint32 {
	return atomic.LoadUint32(&p.workerCount)
}

//WorkerCountInc current goruntine count +1
func (p *WorkerPool) WorkerCountInc() {
	atomic.AddUint32(&p.workerCount, 1)
}

//WorkerCountDec current goruntine count -1
func (p *WorkerPool) WorkerCountDec() {
	atomic.AddUint32(&p.workerCount, ^uint32(-(-1)-1))
}

// Owner get processor pool owner
func (p *WorkerPool) Owner() Module {
	return p.owner
}

// CachePut return processor back to cache
func (p *WorkerPool) CachePut(w *Processor) {
	p.cache.Put(w)
}

// Cap get the capacity
func (p *WorkerPool) Cap() uint32 {
	return atomic.LoadUint32(&p.capacity)
}

// GetFree get the free count
func (p *WorkerPool) GetFree() uint32 {
	return atomic.LoadUint32(&p.capacity) - atomic.LoadUint32(&p.workerCount)
}

// Resize change Processor pool capacity
func (p *WorkerPool) Resize(cap uint32) error {
	if cap == 0 {
		return ErrInvalidcapacity
	} else if cap != p.capacity {
		atomic.StoreUint32(&p.capacity, cap)
		freeCount := int(atomic.LoadUint32(&p.workerCount)) - int(cap)
		for i := 0; i < freeCount; i++ {
			p.getProcessor().q.AsyncPush(nil)
		}
	}
	return nil
}

// Release WorkerPool remove all Processors
func (p *WorkerPool) Release() {
	p.once.Do(func() {
		atomic.StoreUint32(&p.isClosed, 1)
		p.lock.Lock()
		idleWorkers := p.workers
		for i, v := range idleWorkers {
			v.q.AsyncPush(nil)
			idleWorkers[i] = nil
		}
		p.workers = nil
		p.lock.Unlock()
	})
}

// ExpiredCleaning clean expired Processors
func (p *WorkerPool) ExpiredCleaning() {
	for {
		if atomic.LoadUint32(&p.isClosed) == CLOSED {
			break
		}
		time.Sleep(p.expiredDuration)
		now := time.Now()
		p.lock.Lock()
		idleWorkers := p.workers
		var temp []*Processor
		for i, v := range idleWorkers {
			if now.Sub(v.timeout) > p.expiredDuration {
				v.q.AsyncPush(nil)
				idleWorkers[i] = nil
			} else {
				temp = append(temp, v)
			}
		}
		p.workers = temp
		p.lock.Unlock()
	}
}

// getProcessor get one Processor from pool
func (p *WorkerPool) getProcessor() *Processor {
	var w *Processor
	p.lock.Lock()
	// 首先看running是否到达容量限制和是否存在空闲Processor
	idles := p.workers
	if p.workerCount < p.capacity && len(idles) == 0 {
		if cacheWorker := p.cache.Get(); cacheWorker != nil {
			Info("get Processor from cache")
			w = cacheWorker.(*Processor)
			w.Run()
		}
	} else if p.workerCount < p.capacity && len(idles) != 0 {
		w = idles[0]
		p.workers = idles[1:]

	} else if p.workerCount >= p.capacity {
		p.cond.Wait()
		w = idles[0]
		p.workers = idles[1:]

	}
	p.lock.Unlock()
	return w
}

// Submit add obj to Processor
func (p *WorkerPool) Submit(obj *CallObject) error {
	if atomic.LoadUint32(&p.isClosed) == CLOSED {
		return ErrWorkerPoolClosed
	}
	if w := p.getProcessor(); w != nil {
		if err := w.Submit(obj); err != nil {
			return err
		}
	} else {
		return ErrWorkerPoolBusy
	}
	return nil
}

///////////////////

// NewFixedWorkerPool create fixed Processor pool
func NewFixedWorkerPool(owner Module, cap uint32, expired time.Duration) (IWorkerPool, error) {
	if cap == 0 {
		return nil, ErrInvalidcapacity
	}
	capacity, power := PowerCalc(int32(cap))
	// create FixedProcessor pool
	p := &FixedWorkerPool{
		capacity:        uint32(capacity),
		expiredDuration: expired,
		workers:         make([]*Processor, capacity, capacity),
		owner:           owner,
		power:           power,
		name:            owner.Name(),
	}

	//bind signal and lock
	p.cond = sync.NewCond(&p.lock)
	p.cache = sync.Pool{
		New: func() interface{} {
			return NewProcessor(p, 128)
		},
	}
	go p.ExpiredCleaning()

	return p, nil
}

// FixedWorkerPool fixed hash processor pool
type FixedWorkerPool struct {
	// capacity max goruntine count
	capacity uint32
	//power the capacity power
	power uint8
	//workerCount current goruntine Processors
	workerCount uint32

	//Cache sync pool to store Processor
	cache sync.Pool

	// 当关闭该Pool支持通知所有Processor退出运行以防goroutine泄露
	isClosed uint32

	// expiredDuration set timeout for Processor
	expiredDuration time.Duration

	// 互斥锁
	lock sync.Mutex

	// 信号量
	cond *sync.Cond

	// 确保关闭操作只执行一次
	once sync.Once
	name string
	// workers list
	workers []*Processor
	owner   Module
	// fixed pool flag
	fixed bool
}

// Name get the processor pool name
func (p *FixedWorkerPool) Name() string {
	return p.name
}

// RecycleProcessor return back Processor to pool
func (p *FixedWorkerPool) RecycleProcessor(w *Processor) bool {
	if atomic.LoadUint32(&p.isClosed) == CLOSED {
		return false
	}
	w.timeout = time.Now()
	return true
}

// Count get current running Processors count
func (p *FixedWorkerPool) Count() uint32 {
	return atomic.LoadUint32(&p.workerCount)
}

//WorkerCountInc current goruntine count +1
func (p *FixedWorkerPool) WorkerCountInc() {
	atomic.AddUint32(&p.workerCount, 1)
}

//WorkerCountDec current goruntine count -1
func (p *FixedWorkerPool) WorkerCountDec() {
	atomic.AddUint32(&p.workerCount, ^uint32(-(-1)-1))
}

// Owner get processor pool owner
func (p *FixedWorkerPool) Owner() Module {
	return p.owner
}

// CachePut return processor back to cache
func (p *FixedWorkerPool) CachePut(w *Processor) {
	p.cache.Put(w)
}

// Cap get the capacity
func (p *FixedWorkerPool) Cap() uint32 {
	return atomic.LoadUint32(&p.capacity)
}

// GetFree get the free count
func (p *FixedWorkerPool) GetFree() uint32 {
	return atomic.LoadUint32(&p.capacity) - atomic.LoadUint32(&p.workerCount)
}

// Resize change Processor pool capacity
func (p *FixedWorkerPool) Resize(cap uint32) error {
	if cap == 0 {
		return ErrInvalidcapacity
	} else if cap != p.capacity {
		// atomic.StoreUint32(&p.capacity, cap)
		// freeCount := int(atomic.LoadUint32(&p.workerCount)) - int(cap)
		// for i := 0; i < freeCount; i++ {
		// 	p.GetProcessor().q.AsyncPush(nil)
		// }
		//todo
		return ErrorTodo
	}
	return nil
}

// Release WorkerPool remove all Processors
func (p *FixedWorkerPool) Release() {
	p.once.Do(func() {
		atomic.StoreUint32(&p.isClosed, 1)
		p.lock.Lock()
		workers := p.workers
		for i, v := range workers {
			if v != nil {
				v.q.AsyncPush(nil)
				p.workers[i] = nil
			}
		}
		p.workers = nil
		p.lock.Unlock()
	})
}

// ExpiredCleaning clean expired Processors
func (p *FixedWorkerPool) ExpiredCleaning() {
	for {
		if atomic.LoadUint32(&p.isClosed) == CLOSED {
			break
		}
		time.Sleep(p.expiredDuration)
		now := time.Now()
		p.lock.Lock()
		workers := p.workers

		for i, v := range workers {
			if v != nil {
				if now.Sub(v.timeout) > p.expiredDuration {
					v.q.AsyncPush(nil)
					p.workers[i] = nil
				}
			}
		}
		p.lock.Unlock()
	}
}

// getProcessor get one Processor from pool
func (p *FixedWorkerPool) getProcessor(opt uint32) *Processor {
	var w *Processor
	p.lock.Lock()
	// 首先看running是否到达容量限制和是否存在空闲Processor
	workers := p.workers

	hash := opt & ((1 << p.power) - 1)
	if hash >= p.capacity {
		panic("hash function calc error")
	}
	w = workers[hash]

	if w == nil {
		if cacheWorker := p.cache.Get(); cacheWorker != nil {
			w = cacheWorker.(*Processor)
			w.Run()
			workers[hash] = w
		}
	}

	p.lock.Unlock()
	return w
}

// Submit add obj to Processor
func (p *FixedWorkerPool) Submit(obj *CallObject) error {
	if atomic.LoadUint32(&p.isClosed) == CLOSED {
		return ErrWorkerPoolClosed
	}
	if w := p.getProcessor(uint32(obj.Option)); w != nil {
		if err := w.Submit(obj); err != nil {
			return err
		}
	} else {
		return ErrWorkerPoolBusy
	}
	return nil
}
