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
	"sync"
	"sync/atomic"
	"time"

	"github.com/gonethopper/queue"
)

const (
	// DefaultTimeout Processor default timeout
	DefaultTimeout = 10
	// CLOSED status
	CLOSED = 1
)

var (
	//ErrProcessorPoolClosed Processor pool is release
	ErrProcessorPoolClosed = errors.New("Processor pool already closed")
	//ErrProcessorPoolBusy Processor pool is busy
	ErrProcessorPoolBusy = errors.New("Processor pool is busy,please try again")
	//ErrInvalidcapacity set invalid capacity
	ErrInvalidcapacity = errors.New("invalid capacity")
)

// NewProcessorPool create Processor pool
func NewProcessorPool(s Service, cap uint32, expired time.Duration, isNonBlocking bool) (*ProcessorPool, error) {
	if cap == 0 {
		return nil, ErrInvalidcapacity
	}

	// create Processor pool
	p := &ProcessorPool{
		capacity:        cap,
		expiredDuration: expired,
		workers:         make([]*Processor, 0, cap),
		s:               s,
		name:            s.Name(),
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
func NewProcessor(owner *ProcessorPool, queueSize uint32) *Processor {
	return &Processor{
		owner:   owner,
		q:       queue.NewChanQueue(int32(queueSize)),
		timeout: time.Now(),
	}
}

// Processor process job
type Processor struct {
	owner *ProcessorPool
	//CallObject chan
	q queue.Queue
	//timeout set to tigger timeout event
	timeout time.Time
}

// Process goruntine process pre call
func Process(s Service, obj *CallObject) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	var ret = RetObject{
		Ret: nil,
		Err: nil,
	}
	f := s.(Service).GetHandler(obj.Cmd)
	if f == nil {
		err = Error("handler id %v: function not registered", obj.Cmd)
		panic(err)
	} else {
		args := []interface{}{s, obj}
		args = append(args, obj.Args...)
		values := CallUserFunc(f, args...)
		if values == nil {
			err = Error("unsupport handler,need return (interface{},error) or ([]interface{},error)")
			panic(err)
		} else {
			l := len(values)
			if l == 2 {
				ret.Ret = values[0].Interface()
				if values[1].Interface() != nil {
					err = values[1].Interface().(error)
					ret.Err = err
				}
			} else {
				panic(err)
			}
		}
	}
	obj.ChanRet <- ret
}

// Run Processor goruntine
func (w *Processor) Run() {
	atomic.AddUint32(&w.owner.workerCount, 1)
	go func() {
		for {
			obj, err := w.q.Pop()
			if err == nil && obj == nil {
				atomic.AddUint32(&w.owner.workerCount, ^uint32(-(-1)-1))
				w.owner.cache.Put(w)
				break
			}
			if err == nil {
				Process(w.owner.s, obj.(*CallObject))
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

// ProcessorPool Processor pool limit goruntine max count
type ProcessorPool struct {
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
	// // 空闲的Processor队列
	workers []*Processor
	s       Service
}

// RecycleProcessor return back Processor to pool
func (p *ProcessorPool) RecycleProcessor(w *Processor) bool {
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
func (p *ProcessorPool) Count() uint32 {
	return atomic.LoadUint32(&p.workerCount)
}

// Cap get the capacity
func (p *ProcessorPool) Cap() uint32 {
	return atomic.LoadUint32(&p.capacity)
}

// GetFree get the free count
func (p *ProcessorPool) GetFree() uint32 {
	return atomic.LoadUint32(&p.capacity) - atomic.LoadUint32(&p.workerCount)
}

// Submit add obj to Processor
func (p *ProcessorPool) Submit(obj *CallObject) error {
	if atomic.LoadUint32(&p.isClosed) == CLOSED {
		return ErrProcessorPoolClosed
	}
	if w := p.GetProcessor(); w != nil {
		if err := w.Submit(obj); err != nil {
			return err
		}
	} else {
		return ErrProcessorPoolBusy
	}
	return nil
}

// Resize change Processor pool capacity
func (p *ProcessorPool) Resize(cap uint32) error {
	if cap == 0 {
		return ErrInvalidcapacity
	} else if cap != p.capacity {
		atomic.StoreUint32(&p.capacity, cap)
		freeCount := int(atomic.LoadUint32(&p.workerCount)) - int(cap)
		for i := 0; i < freeCount; i++ {
			p.GetProcessor().q.AsyncPush(nil)
		}
	}
	return nil
}

// Release ProcessorPool remove all Processors
func (p *ProcessorPool) Release() {
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
func (p *ProcessorPool) ExpiredCleaning() {
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

// GetProcessor get one Processor from pool
func (p *ProcessorPool) GetProcessor() *Processor {
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
