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
// * @Date: 2019-06-12 15:53:22
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-12 15:53:22

package server

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// GBytesPool pre-create []byte pool
var GBytesPool *BytesPool

// GMessagePool pre-create message pool
var GMessagePool *MessagePool

// GSessionPool pre-create session
var GSessionPool *SessionPool

// log variable start

// GLoggerService global log service
var GLoggerService Service

// WG global goruntine wait group
var WG sync.WaitGroup

// service variable start

// AnonymousServiceID Anonymous Service Counter
var AnonymousServiceID int32 = ServiceIDNamedMax

// relServices relate name to create service function
var relServices = make(map[string]func() (Service, error))

// service variable end

// App server instance
var App *Server

func init() {
	GBytesPool = NewBytesPool()
	GMessagePool = NewMessagePool()
	GSessionPool = NewSessionPool()
	App = &Server{
		GoCount: 0,
	}

	fmt.Println("Nethopper Framework init")
}

// GracefulExit server exit by call root context close
func GracefulExit() {
	WG.Wait()
	GLoggerService.Close()
	// wait root context done
	for {
		if _, exitFlag := GLoggerService.CanExit(true); exitFlag {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// Server server entity, only one instance
type Server struct {
	// GoCount total goruntine count
	GoCount  int32
	Services sync.Map
}

// UpdateGoCount update goruntine use count ,+/- is all ok
func (s *Server) UpdateGoCount(value int32) {
	atomic.AddInt32(&s.GoCount, value)
}
