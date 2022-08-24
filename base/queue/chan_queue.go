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
// * @Date: 2019-12-18 10:47:24
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-18 10:47:24

package queue

import (
	"sync/atomic"
	"time"
)

// ChanQueue use chan queue
type ChanQueue struct {
	innerChan  chan interface{}
	capacity   int32
	size       int32
	timer      *time.Timer
	closedChan chan struct{}
}

// NewChanQueue create chan queue
func NewChanQueue(capacity int32) Queue {
	return &ChanQueue{
		innerChan:  make(chan interface{}, capacity),
		capacity:   capacity,
		size:       0,
		timer:      time.NewTimer(time.Second),
		closedChan: make(chan struct{}),
	}
}

// Pop sync block pop
func (q *ChanQueue) Pop() (val interface{}, err error) {

	v, ok := <-q.innerChan
	if ok {
		atomic.AddInt32(&q.size, -1)
		return v, nil
	}
	return nil, ErrQueueIsClosed

}
func (q *ChanQueue) AutoPop() ([]interface{}, error) {
	var batch []interface{}
	var val interface{}
	var err error
	for i := 0; i < 128; i++ {
		val, err = q.AsyncPop()
		if err != nil {
			break
		}
		batch = append(batch, val)
	}
	if len(batch) == 0 {
		val, err = q.Pop()
		if err != nil {
			return nil, err
		}
		batch = append(batch, val)
	}
	return batch, nil
}

// AsyncPop async pop
func (q *ChanQueue) AsyncPop() (val interface{}, err error) {

	select {
	case v, ok := <-q.innerChan:
		if ok {
			atomic.AddInt32(&q.size, -1)
			return v, nil
		}
		return nil, ErrQueueIsClosed
	default:
		return nil, ErrQueueEmpty
	}

}

// Push sync push data
func (q *ChanQueue) Push(x interface{}) error {

	if q.IsClosed() {
		return ErrQueueIsClosed
	}
	q.innerChan <- x
	atomic.AddInt32(&q.size, 1)
	return nil
}

// AsyncPush async push data
func (q *ChanQueue) AsyncPush(x interface{}) (err error) {

	if q.IsClosed() {
		err = ErrQueueIsClosed
		return err
	}

	for i := 0; i < 3; i++ {
		err = q.tryPush(x)
		if err != ErrQueueFull {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}
func (q *ChanQueue) tryPush(x interface{}) error {
	select {
	case q.innerChan <- x:
		atomic.AddInt32(&q.size, 1)
		return nil
	default:
		return ErrQueueFull
	}
}

// Length get chan queue length
func (q *ChanQueue) Length() int32 {
	return q.size
}

// Capacity get queue capacity
func (q *ChanQueue) Capacity() int32 {
	return q.capacity
}

// IsFull queue is full return true
func (q *ChanQueue) IsFull() bool {
	return len(q.innerChan) == cap(q.innerChan)
}

// Close 不需要关闭innerChan,交给GC回收,多写的时候直接关闭innerChan会出问题
func (q *ChanQueue) Close() error {
	if q.IsClosed() {
		return ErrQueueIsClosed
	}
	close(q.closedChan)

	return nil
}
func (q *ChanQueue) ForceClose() error {

	close(q.innerChan)
	return nil
}

// IsClosed if chan is close,return true
func (q *ChanQueue) IsClosed() bool {
	select {
	case <-q.closedChan:
		return true
	default:
	}
	return false

}

func (q *ChanQueue) getChan(timeout time.Duration) (<-chan interface{}, <-chan error) {
	timeoutChan := make(chan error, 1)
	resultChan := make(chan interface{}, 1)
	go func() {
		if timeout < 0 {
			item := <-q.innerChan
			atomic.AddInt32(&q.size, -1)
			resultChan <- item
		} else {
			select {
			case item := <-q.innerChan:
				atomic.AddInt32(&q.size, -1)
				resultChan <- item
			case <-time.After(timeout):
				timeoutChan <- ErrQueueTimeout
			}
		}
	}()
	return resultChan, timeoutChan
}
