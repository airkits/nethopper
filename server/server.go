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
	"sync"
	"sync/atomic"
)

// log variable start

var logger Service

// LogLevel log level setting
var LogLevel int32

// log variable end

//SetLogLevel set log level to app global var
func SetLogLevel(level int32) {
	LogLevel = atomic.LoadInt32(&level)
}

// WG global goruntine wait group
var WG sync.WaitGroup

// service variable start

// AnonymousServiceID Anonymous Service Counter
var AnonymousServiceID int32 = ServiceIDNamedMax

// refServices relate name to create service function
var refServices = make(map[string]func() (Service, error))

// service variable end

// App server instance
var App = &Server{
	GoCount: 0,
}

// Server server entity, only one instance
type Server struct {
	// GoCount total goruntine count
	GoCount  int32
	Services sync.Map
}
