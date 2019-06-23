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
// * @Date: 2019-06-06 14:49:24
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-06 14:49:24

package log_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/utils"
)

func TestFormatLog(t *testing.T) {
	//test msg without params
	msg := "format log test"
	format := "%s%s%s\n"
	var level int32
	for level = log.FATAL; level < log.DEBUG; level++ {
		expect := fmt.Sprintf(format, utils.TimeYMDHIS(), log.LogLevelPrefix[level], msg)
		result := log.FormatLog(level, msg)
		if expect != result {
			t.Errorf("\nexpect :%s,\nresult :%s", expect, result)
		}
	}
	//test msg with params
	msg = "format %s log test %d"
	for level = log.FATAL; level < log.DEBUG; level++ {
		expect := fmt.Sprintf(format, utils.TimeYMDHIS(), log.LogLevelPrefix[level], fmt.Sprintf(msg, strconv.Itoa(int(level)), level))
		result := log.FormatLog(level, msg, strconv.Itoa(int(level)), level)
		if expect != result {
			t.Errorf("\nexpect :%s,\nresult :%s", expect, result)
		}
	}

}

func TestNewFileLogger(t *testing.T) {
	m := map[string]interface{}{
		"filename":    "logs/server_log.log",
		"level":       7,
		"maxSize":     50,
		"maxLines":    1000,
		"hourEnabled": false,
		"dailyEnable": true,
		"queueSize":   1000,
	}
	logger, err := log.NewFileLogger(m)
	if err != nil {
		t.Error(err)
	}

	logger.Debug("helloword %d", 1234)

	logger.Close()

}
