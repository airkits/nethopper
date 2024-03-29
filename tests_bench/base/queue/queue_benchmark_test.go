package queue_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/airkits/nethopper/base/queue"
)

func bench_sync(q queue.Queue) {
	var wg sync.WaitGroup
	num := 1000000
	pushNum := 5
	wg.Add(pushNum + 1)
	for i := 0; i < pushNum; i++ {
		go func(num int) {
			for j := 0; j < num; j++ {
				q.Push(j)
			}
			wg.Done()
		}(num)
	}
	go func() {
		idx := 0
		for idx < num*pushNum {
			_, err := q.Pop()
			if err == nil {
				idx += 1
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

func bench_syncb(q queue.Queue) {
	var wg sync.WaitGroup
	num := 1000000
	pushNum := 5
	wg.Add(pushNum + 1)
	for i := 0; i < pushNum; i++ {
		go func(num int) {
			for j := 0; j < num; j++ {
				q.Push(j)
			}
			wg.Done()
		}(num)
	}
	go func() {
		idx := 0
		for idx < num*pushNum {
			result, err := q.AutoPop()
			if err == nil {
				idx += len(result)
			}
		}
		wg.Done()
	}()
	wg.Wait()
}

func Benchmark_ChanQueue(b *testing.B) {
	q := queue.NewChanQueue(1000)
	//bench_sync(q)
	bench_syncb(q)
}

func BenchmarkPush(b *testing.B) {
	b.ResetTimer()
	q := queue.NewChanQueue(1000)
	p := queue.NewChanQueue(1000)
	o := queue.NewChanQueue(1000)
	var wg sync.WaitGroup
	wg.Add(4)
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
			p.Push(v)
		}
		wg.Done()
	}()
	go func() {
		for j := 0; j < b.N; j++ {
			v, _ := p.Pop()
			if v.(int) != j {
				b.Error("push error %ld", j)
			}
			o.Push(v)
		}
		wg.Done()
	}()
	go func() {
		for j := 0; j < b.N; j++ {
			v, _ := o.Pop()
			if v.(int) != j {
				b.Error("push error %ld", j)
			}

		}
		wg.Done()
	}()
	wg.Wait()
}
func genRandomList(size int) []int {
	list := make([]int, size)
	for i, _ := range list {
		v := rand.Intn(10)
		list[i] = v
	}
	return list
}

func test_queue(q queue.Queue, size int) bool {
	list := genRandomList(size)
	fmt.Printf("test_queue, %d\n", size)

	go func() {
		for _, v := range list {
			// fmt.Printf("put, %d\n", v)
			q.Push(v)
		}

	}()

	for _, v := range list {
		v2, _ := q.Pop()
		if v != v2 {
			fmt.Printf("vail, %d != %d", v, v2)
			return false
		}
	}

	return true
}

func Test_ChanQueue(t *testing.T) {
	q := queue.NewChanQueue(1000)
	r := test_queue(q, 10000)
	if r == false {
		t.Error("chan queue error")
	}
}
