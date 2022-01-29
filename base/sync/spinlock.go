// Package sync provides synchronization primitive implementations for spinlocks
// and semaphore.
package sync

import (
	"runtime"
	"sync/atomic"
)

var (
	// TODO: replace with real yield function when context-switching is implemented.
	yieldFn func()
)

const maxBackoff = 64

// noCopy may be embedded into structs which must not be copied
// after the first use.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

// Spinlock implements a lock where each task trying to acquire it busy-waits
// till the lock becomes available.
type Spinlock struct {
	state  uint32
	noCopy noCopy
}

// Acquire blocks until the lock can be acquired by the currently active task.
// Any attempt to re-acquire a lock already held by the current task will cause
// a deadlock.
//func (l *Spinlock) Acquire() {
//	// archAcquireSpinlock(&l.state, 1)
//}

// TryToAcquire attempts to acquire the lock and returns true if the lock could
// be acquired or false otherwise.
func (sl *Spinlock) TryToAcquire() bool {
	return atomic.SwapUint32(&sl.state, 1) == 0
}

// Release relinquishes a held lock allowing other tasks to acquire it. Calling
// Release while the lock is free has no effect.
func (sl *Spinlock) Release() {
	atomic.StoreUint32(&sl.state, 0)
}

func (sl *Spinlock) Lock() {
	backoff := 1
	for !atomic.CompareAndSwapUint32(&sl.state, 0, 1) {
		// Leverage the exponential backoff algorithm, see https://en.wikipedia.org/wiki/Exponential_backoff.
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < maxBackoff {
			backoff <<= 1
		}
	}
}

func (sl *Spinlock) Unlock() {
	atomic.StoreUint32(&sl.state, 0)
}

// NewSpinLock instantiates a spin-lock.
func NewSpinLock() *Spinlock {
	return new(Spinlock)
}

// archAcquireSpinlock is an arch-specific implementation for acquiring the lock.
// func archAcquireSpinlock(state *uint32, attemptsBeforeYielding uint32)
