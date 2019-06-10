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
	"sync"
	"testing"

	"github.com/gonethopper/nethopper/log"
)

const Step = 1000000

//BenchmarkFormatLog format test
func BenchmarkFormatLog(t *testing.B) {
	msg := "format log test"
	for i := 0; i < Step; i++ {
		_ = log.FormatLog(log.INFO, msg)
	}

}

func BenchmarkFormatLogWithParams(t *testing.B) {
	msg := "format %d log test"
	for i := 0; i < Step; i++ {
		_ = log.FormatLog(log.INFO, msg, i)
	}
}

func BenchmarkWriteFileLog(t *testing.B) {
	m := map[string]interface{}{
		"filename":    "test/server.log",
		"level":       7,
		"maxSize":     50,
		"maxLines":    int(Step / 10),
		"hourEnabled": true,
		"dailyEnable": true,
	}
	logger, err := log.NewFileLogger(m)
	if err != nil {
		t.Error(err)
	}
	defer logger.Close()
	var wg sync.WaitGroup
	wg.Add(5)
	for j := 0; j < 5; j++ {
		go func() {
			for i := 0; i < Step/5; i++ {
				logger.Debug("helloword %d true filename:test/server.log hourEnabled:true level:7 maxLines:100000", i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
