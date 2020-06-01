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
// snowflake 算法参考 https://github.com/gonet2/snowflake/blob/master/service.go

// * @Author: ankye
// * @Date: 2019-06-24 11:07:19
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-24 11:07:19

package logic

import (
	"sync"
	"time"

	"github.com/gonethopper/nethopper/examples/model/common"
	"github.com/gonethopper/nethopper/examples/snowflake/global"
	"github.com/gonethopper/nethopper/server"
)

const (
	BACKOFF    = 100  // max backoff delay millisecond
	CONCURRENT = 128  // max concurrent connections to etcd
	UUID_QUEUE = 1024 // uuid process queue
)

const (
	TS_MASK         = 0x1FFFFFFFFFF // 41bit
	SN_MASK         = 0xFFF         // 12bit
	MACHINE_ID_MASK = 0x3FF         // 10bit
)

// Module struct to define module
type Module struct {
	server.BaseContext
	pkroot     string
	uuidkey    string
	machine_id uint64 // 10-bit machine id
	chProc     chan chan uint64
	muNext     sync.Mutex
}

// get timestamp
func ts() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// ModuleCreate  module create function
func ModuleCreate() (server.Module, error) {
	return &Module{}, nil
}

// UserData module custom option, can you store you data and you must keep goruntine safe
// func (s *Module) UserData() int32 {
// 	return 0
// }

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *Module) Setup(conf server.IConfig) (server.Module, error) {
	s.RegisterHandler(common.CallIDGenUIDCmd, UUIDHandler)
	s.CreateWorkerPool(s, 128, 10*time.Second, true)
	cfg := global.GetInstance().GetConfig()
	s.chProc = make(chan chan uint64, UUID_QUEUE)
	// shifted machine id
	s.machine_id = (uint64(cfg.SID) & MACHINE_ID_MASK) << 12
	go s.uuidTask()

	return s, nil
}

// GenerateUUID an unique uuid
func (s *Module) GenerateUUID() (uint64, error) {
	req := make(chan uint64, 1)
	s.chProc <- req
	return <-req, nil
}

// uuidTask generator
func (s *Module) uuidTask() {
	var sn uint64    // 12-bit serial no
	var lastTs int64 // last timestamp
	for {
		ret := <-s.chProc
		// get a correct serial number
		t := ts()
		if t < lastTs { // clock shift backward
			server.Warning("clock shift happened, waiting until the clock moving to the next millisecond.")
			t = s.waitMillisecond(lastTs)
		}

		if lastTs == t { // same millisecond
			sn = (sn + 1) & SN_MASK
			if sn == 0 { // serial number overflows, wait until next ms
				t = s.waitMillisecond(lastTs)
			}
		} else { // new millsecond, reset serial number to 0
			sn = 0
		}
		// remember last timestamp
		lastTs = t

		// generate uuid, format:
		//
		// 0		0.................0		0..............0	0........0
		// 1-bit	41bit timestamp			10bit machine-id	12bit sn
		var uuid uint64
		uuid |= (uint64(t) & TS_MASK) << 22
		uuid |= s.machine_id
		uuid |= sn
		ret <- uuid
	}
}

// waitMillisecond will wait untill last_ts
func (s *Module) waitMillisecond(lastTs int64) int64 {
	t := ts()
	for t < lastTs {
		time.Sleep(time.Duration(lastTs-t) * time.Millisecond)
		t = ts()
	}
	return t
}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
	server.RunSimpleFrame(s, 128)
}

// Stop goruntine
func (s *Module) Stop() error {
	return nil
}

// Call async send message to module
// func (s *Module) Call(option int32, obj *server.CallObject) error {
// 	if err := s.MQ().AsyncPush(obj); err != nil {
// 		server.Error(err.Error())
// 	}
// 	return nil
// }

// PushBytes async send string or bytes to queue
// func (s *Module) PushBytes(option int32, buf []byte) error {
// 	return nil
// }
