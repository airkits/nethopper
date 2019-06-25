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
// * @Date: 2019-06-24 13:35:23
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-24 13:35:23

package log

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/gonethopper/nethopper/server"
)

// NewConsoleLogger create FileLog instance
func NewConsoleLogger(m map[string]interface{}) (server.Log, error) {
	logger := &ConsoleLog{}
	if err := logger.ParseConfig(m); err != nil {
		return nil, err
	}
	if err := logger.InitLogger(); err != nil {
		return nil, err
	}
	return logger, nil
}

// ConsoleLog writes messages to terminal.
type ConsoleLog struct {
	//set level and  atomic incr CurrentSize and CurrentLines
	level  int32
	writer io.Writer
}

// ParseConfig read config from map[string]interface{}
// config key map
// level default 7
func (c *ConsoleLog) ParseConfig(m map[string]interface{}) error {
	level, err := server.ParseValue(m, "level", 7)
	if err != nil {
		return err
	}
	c.level = int32(level.(int))

	return nil
}

// InitLogger init logger
func (c *ConsoleLog) InitLogger() error {
	c.writer = os.Stdout
	return nil
}

// SetLevel atomic set level
func (c *ConsoleLog) SetLevel(level int32) error {
	if level < server.FATAL || level > server.DEBUG {
		return fmt.Errorf("log level:[%d] invalid", level)
	}
	atomic.StoreInt32(&c.level, level)
	return nil
}

// GetLevel atomic get level
func (c *ConsoleLog) GetLevel() int32 {
	level := atomic.LoadInt32(&c.level)
	return level
}

// WriteLog write log to file, return immediately if not meet the conditions
func (c *ConsoleLog) WriteLog(msg []byte, count int32) error {
	c.writer.Write(msg)
	return nil
}

// CanLog check log status
func (c *ConsoleLog) CanLog(msgSize int32, count int32) bool {
	return true
}

// Close and flush
func (c *ConsoleLog) Close() error {
	return nil
}

// PushLog push log to queue
func (c *ConsoleLog) PushLog(level int32, v ...interface{}) error {
	if level > c.level {
		return nil
	}
	msg := server.FormatLog(level, v...)
	return c.WriteLog([]byte(msg), 1)

}

// Fatal system is unusable
func (c *ConsoleLog) Fatal(v ...interface{}) error {
	return c.PushLog(server.FATAL, v...)
}

// Error error conditions
func (c *ConsoleLog) Error(v ...interface{}) error {
	return c.PushLog(server.ERROR, v...)
}

// Warning warning conditions
func (c *ConsoleLog) Warning(v ...interface{}) error {
	return c.PushLog(server.WARNING, v...)
}

// Info informational messages
func (c *ConsoleLog) Info(v ...interface{}) error {
	return c.PushLog(server.INFO, v...)
}

// Debug debug-level messages
func (c *ConsoleLog) Debug(v ...interface{}) error {
	return c.PushLog(server.DEBUG, v...)
}
