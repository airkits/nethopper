package utils

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// TVRBits tvr bits
	TVRBits = 8
	// TVRSize tvr size
	TVRSize = 1 << TVRBits
	// TVRMask tvr mask
	TVRMask = TVRSize - 1

	// TVNBits tvn bits
	TVNBits = 6
	// TVNSize tvn size
	TVNSize = 1 << TVNBits
	// TVNMask tvn mask
	TVNMask = TVNSize - 1
)
const (
	//TimerTypeOnce timer run once and remove after run
	TimerTypeOnce = 0
	//TimerTypeLoop timer run forever
	TimerTypeLoop = 1
)

func offset(n uint64) uint64 {
	return TVRSize + n*TVNSize
}

func index(v, n uint64) uint64 {
	return ((v >> (TVRBits + n*TVNBits)) & TVNMask)
}

// TimerID define timerid
type TimerID uint64

// Timer definition
type Timer struct {
	interval uint64
	expire   uint64
	index    int32
	id       TimerID
	_type    int32
	f        func(interface{})
	// for test
	Str string
	I   int
}

// TimerManager to manger timer
type TimerManager struct {
	//采用hashmap 主要考虑到需要快速update 带来的影响是需要同时触发的定时器不是先进先出的
	tv        [TVRSize + 4*TVNSize]map[TimerID]*Timer
	checktime uint64
	/* 定时器精度 ，确定了定时范围 (0-2^32)*tick.精度越高,定时越不准。
	* 5us以下的定时可能不准确,不支持 0间隔定时器 interval=0 会直接添加失败
	 */
	tick    time.Duration
	idIndex map[TimerID]int32
	mutex   sync.Mutex
}

var (
	id uint64
)

// NewTimerID create TimerID
func NewTimerID() TimerID {
	atomic.AddUint64(&id, 1)
	return TimerID(id)
}

// NewTimer create timer by time type
func NewTimer(timerType int32) (*Timer, TimerID) {
	_id := NewTimerID()
	timer := &Timer{id: _id, _type: timerType}
	return timer, _id
}

//NewTimerManager use timermanager to manger timers
func NewTimerManager(tick time.Duration) *TimerManager {
	tm := &TimerManager{tick: tick}
	tm.mutex.Lock()
	tm.checktime = uint64(time.Now().UnixNano() / int64(tick))
	for i, _ := range tm.tv {
		tm.tv[i] = make(map[TimerID]*Timer)
	}
	tm.idIndex = make(map[TimerID]int32)
	tm.mutex.Unlock()
	return tm
}

// Stop timer
func (t *Timer) Stop(tm *TimerManager) {
	if t.index != -1 {
		tm.RemoveTimerInLock(t)
		t.index = -1
	}
}

// Start timer with timerManager
func (t *Timer) Start(interval uint64, f func(interface{}), tm *TimerManager) error {
	if interval == 0 {
		err := fmt.Errorf("invalid interval time! %d", interval)
		return err
	}
	t.Stop(tm)
	t.interval = interval
	t.f = f
	t.expire = interval + uint64(time.Now().UnixNano()/int64(tm.tick))
	tm.AddTimerInLock(t)
	return nil
}

// Update timer update with timerManager
func (t *Timer) Update(interval uint64, f func(interface{}), tm *TimerManager) {
	t.Start(interval, f, tm)
}

// AddTimerInLock add timer to TimerManager,goruntine safe
func (t *TimerManager) AddTimerInLock(timer *Timer) {
	t.mutex.Lock()
	t.addTimer(timer)
	t.mutex.Unlock()
}

// addTimer add timer to TimerManager
func (t *TimerManager) addTimer(timer *Timer) {
	expires := timer.expire
	idx := expires - t.checktime
	if idx < TVRSize {
		timer.index = int32(expires & TVRMask)
	} else if idx < 1<<(TVRBits+TVNBits) {
		timer.index = int32(offset(0) + index(expires, 0))
	} else if idx < 1<<(TVRBits+2*TVNBits) {
		timer.index = int32(offset(1) + index(expires, 1))
	} else if idx < 1<<(TVRBits+3*TVNBits) {
		timer.index = int32(offset(2) + index(expires, 2))
	} else if int64(idx) < 0 {
		timer.index = int32(t.checktime & TVRMask)
	} else {
		if idx > 0xffffffff {
			idx = 0xffffffff
			expires = idx + t.checktime
		}
		timer.index = int32(offset(3) + index(expires, 3))
	}
	t.idIndex[timer.id] = timer.index
	timermap := t.tv[timer.index]
	timermap[timer.id] = timer
}

func (t *TimerManager) removeTimer(timer *Timer) {
	timermap := t.tv[timer.index]
	delete(timermap, timer.id)
	delete(t.idIndex, timer.id)
}

// RemoveTimerInLock remove timer ,goruntine safe
func (t *TimerManager) RemoveTimerInLock(timer *Timer) {
	t.mutex.Lock()
	t.removeTimer(timer)
	t.mutex.Unlock()
}

func (t *TimerManager) detectTimer() {
	now := uint64(time.Now().UnixNano() / int64(t.tick))
	for tnow := now; t.checktime <= tnow; t.checktime++ {
		index1 := t.checktime & TVRMask
		if index1 == 0 &&
			t.cascade(offset(0), index(t.checktime, 0)) == 0 &&
			t.cascade(offset(1), index(t.checktime, 1)) == 0 &&
			t.cascade(offset(2), index(t.checktime, 2)) == 0 {
			t.cascade(offset(3), index(t.checktime, 3))
		}
		timermap := t.tv[index1]
		for k, v := range timermap {
			v.f(v.id)
			delete(timermap, k)
			if v._type == TimerTypeLoop {
				v.expire = v.interval + now
				t.addTimer(v)
			}
		}
	}
}

func (t *TimerManager) cascade(offset, index uint64) uint64 {
	timermap := t.tv[offset+index]
	for k, v := range timermap {
		t.addTimer(v)
		delete(timermap, k)
	}
	return index
}

// DetectTimerInLock detect timer and clean timeout
func (t *TimerManager) DetectTimerInLock() {
	t.mutex.Lock()
	t.detectTimer()
	t.mutex.Unlock()
}

//GetTimerByID get timer by id
func (t *TimerManager) GetTimerByID(ID TimerID) *Timer {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if v, ok := t.idIndex[ID]; ok {
		timermap := t.tv[v]
		if v1, ok := timermap[ID]; ok {
			return v1
		}
	}
	return nil

}
