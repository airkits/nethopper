package mediator_test

import (
	"sync"
	"testing"
	"time"

	"github.com/airkits/nethopper/base/queue"
	"github.com/airkits/nethopper/mediator"
	"github.com/airkits/nethopper/utils"
)

func init() {
	// appContext := base.NewAppContext()
	// conf := log.Config{
	// 	Filename:     "logs/server_log.log",
	// 	Level:        7,
	// 	MaxLines:     1000,
	// 	MaxSize:      50,
	// 	HourEnabled:  true,
	// 	DailyEnabled: true,
	// 	QueueSize:    1000,
	// }
	// log.InitLogger(appContext, &conf)

}

type MCall struct {
}

func (m *MCall) Execute(obj *mediator.CallObject) *mediator.RetObject {
	utils.GenUID()
	return mediator.NewRetObject(1, nil, obj.Args[0])

}
func BenchmarkWorkerPool(b *testing.B) {
	b.ResetTimer()
	m := &MCall{}
	wp, err := mediator.NewFixedWorkerPool(2048, 256, 30*time.Second)
	if err != nil {
		b.Error(err)
	}
	obj := mediator.NewCallObject(m, 1, int32(1), 1)
	wp.Submit(obj)

	for j := 0; j < b.N; j++ {

		obj := mediator.NewCallObject(m, 1, int32(j), j)
		wp.Submit(obj)

		result := <-obj.ChanRet
		if result.Data != j {
			b.Error("get result failed")
		}

	}
}
func BenchmarkWorkerSend(b *testing.B) {
	b.ResetTimer()
	q := queue.NewChanQueue(1000)
	m := &MCall{}
	wp, err := mediator.NewFixedWorkerPool(128, 128, 30*time.Second)
	if err != nil {
		b.Error(err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		for j := 0; j < b.N; j++ {
			q.Push(j)
		}
		wg.Done()
	}()
	f := func() {
		for j := 0; j < b.N; j++ {
			v, _ := q.Pop()
			obj := mediator.NewCallObject(m, 1, utils.RandomInt32(0, 1000), v)
			wp.Submit(obj)
			result := <-obj.ChanRet
			if result.Data != v {
				b.Error("get result failed")
			}
		}
		wg.Done()
	}
	go f()

	wg.Wait()
}

func BenchmarkWorker(b *testing.B) {
	b.ResetTimer()
	q := queue.NewChanQueue(1000)
	p := queue.NewChanQueue(1000)
	m := &MCall{}
	wp, err := mediator.NewFixedWorkerPool(128, 128, 30*time.Second)
	if err != nil {
		b.Error(err)
	}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		for j := 0; j < b.N; j++ {
			q.Push(j)
		}
		wg.Done()
	}()
	go func() {
		for j := 0; j < b.N; j++ {
			v, _ := q.Pop()
			if v.(int) != j {
				b.Error("push error %ld", j)
			}
			obj := mediator.NewCallObject(m, 1, utils.RandomInt32(0, 1000), j)
			wp.Submit(obj)
			result := <-obj.ChanRet
			p.Push(result.Data)
		}
		wg.Done()
	}()
	go func() {
		var sum = 0
		var s = 0
		for j := 0; j < b.N; j++ {
			s += j
			v, _ := p.Pop()
			sum += v.(int)
		}
		if s != sum {
			b.Error("push error %ld %ld", sum, s)
		}
		wg.Done()
	}()
	wg.Wait()
}
