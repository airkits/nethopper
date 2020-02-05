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
// * @Date: 2019-06-06 16:57:24
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-06 16:57:24

package log_test

import (
	"runtime"
	"testing"

	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/server"
)

const Step = 1000000

// BenchmarkFormatLog format test
func BenchmarkFormatLog(t *testing.B) {

	msg := "format log test"
	for i := 0; i < Step; i++ {
		_ = server.FormatLog(server.INFO, msg)
	}

}

func BenchmarkFormatLogWithParams(t *testing.B) {

	msg := "format %d log test"
	for i := 0; i < Step; i++ {
		_ = server.FormatLog(server.INFO, msg, i)
	}
}

func BenchmarkWriteLog(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	m := map[string]interface{}{
		"filename":    "logs/server.log",
		"level":       4,
		"maxSize":     50,
		"maxLines":    100000,
		"hourEnabled": false,
		"dailyEnable": true,
		"queueSize":   1000,
	}

	logger, err := log.NewFileLogger(m)
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < Step; i++ {
		logger.Debug("helloword true filename:testserver.log hourEnabled:true level:7 maxLines:100000")
	}
	logger.Close()

}
