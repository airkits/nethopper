package mediator

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/airkits/nethopper/base"
	"github.com/airkits/nethopper/log"
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

// IWorkerPool process pool interface
type IWorkerPool interface {
	//Setup
	Setup(queueSize uint32)
	// Submit add obj to Processor
	Submit(obj *base.CallObject) error

	//AddRef current goruntine count +1
	AddRef()

	//DecRef current goruntine count -1
	DecRef()

	//CachePut put processor return to pool
	CachePut(w *Processor)

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

	// workers list
	workers []*Processor

	// fixed pool flag
	fixed bool
}

func (p *WorkerPool) Setup(queueSize uint32) {
	//bind signal and lock
	p.cond = sync.NewCond(&p.lock)
	p.cache = sync.Pool{
		New: func() interface{} {
			return NewProcessor(p, queueSize)
		},
	}
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

// AddRef current goruntine count +1
func (p *WorkerPool) AddRef() {
	atomic.AddUint32(&p.workerCount, 1)
}

// DecRef current goruntine count -1
func (p *WorkerPool) DecRef() {
	atomic.AddUint32(&p.workerCount, ^uint32(-(-1)-1))
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

	p.lock.Lock()
	defer p.lock.Unlock()
	var w *Processor
	// 首先看running是否到达容量限制和是否存在空闲Processor
	idles := p.workers
	if p.workerCount < p.capacity && len(idles) == 0 {
		if cacheWorker := p.cache.Get(); cacheWorker != nil {
			log.Info("get Processor from cache")
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

	return w
}

// Submit add obj to Processor
func (p *WorkerPool) Submit(obj *base.CallObject) error {
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
	// workers list
	workers []*Processor
	// fixed pool flag
	fixed bool
}

func (p *FixedWorkerPool) Setup(queueSize uint32) {
	//bind signal and lock
	p.cond = sync.NewCond(&p.lock)
	p.cache = sync.Pool{
		New: func() interface{} {
			return NewProcessor(p, queueSize)
		},
	}
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

// AddRef current goruntine count +1
func (p *FixedWorkerPool) AddRef() {
	atomic.AddUint32(&p.workerCount, 1)
}

// DecRef current goruntine count -1
func (p *FixedWorkerPool) DecRef() {
	atomic.AddUint32(&p.workerCount, ^uint32(-(-1)-1))
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
	// p.lock.Lock()
	// defer p.lock.Unlock()
	// 首先看running是否到达容量限制和是否存在空闲Processor
	workers := p.workers

	hash := opt & ((1 << p.power) - 1)
	if hash >= p.capacity {
		panic("hash function calc error")
	}

	w = workers[hash]
	//log.Info("[Mediator] worker user id:[%d] processor", hash)
	if w == nil {
		if cacheWorker := p.cache.Get(); cacheWorker != nil {
			w = cacheWorker.(*Processor)
			w.Run()
			p.lock.Lock()
			workers[hash] = w
			p.lock.Unlock()
		}
	}

	return w
}

// Submit add obj to Processor
func (p *FixedWorkerPool) Submit(obj *base.CallObject) error {
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
