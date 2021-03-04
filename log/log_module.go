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
// * @Date: 2019-06-14 14:40:51
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-14 14:40:51

package log

import (
	"bytes"
	"time"

	"github.com/airkits/nethopper/server"
)

// LogModule struct implements the interface Module
type LogModule struct {
	server.BaseContext
	logger  server.Log
	console server.Log
	//for stat
	buf     bytes.Buffer
	count   int32
	msgSize int32
}

// LogModuleCreate log module create function
func LogModuleCreate() (server.Module, error) {
	return &LogModule{}, nil
}

// Setup init and setup config
// Log config
// m := map[string]interface{}{
// 	"filename":    "server.log",
// 	"level":       4,
// 	"maxSize":     50,
// 	"maxLines":    1000,
// 	"hourEnabled": false,
// 	"dailyEnable": true,
//  "queueSize":1000,
// }
func (s *LogModule) Setup(conf server.IConfig) (server.Module, error) {

	logger, err := NewFileLogger(conf.(*Config))
	if err != nil {
		return nil, err
	}
	s.logger = logger
	console, err := NewConsoleLogger(conf.(*Config))
	if err != nil {
		return nil, err
	}
	s.console = console
	return s, nil
}

// Reload reload config from map
func (s *LogModule) Reload(conf server.IConfig) error {

	return s.logger.SetLevel(conf.(*Config).Level)

}

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *LogModule) OnRun(dt time.Duration) {
	s.msgSize = 0
	s.count = 0
	for i := 0; i < 128; i++ {

		if v, err := s.MQ().AsyncPop(); err == nil {
			if n, e := s.buf.Write(v.([]byte)); e == nil {
				s.msgSize += int32(n)
				s.count++
				if !s.logger.CanLog(s.msgSize, s.count) {
					break
				}
			}
		}
	}

	if s.buf.Len() > 0 {
		s.logger.WriteLog(s.buf.Bytes(), s.count)
		s.console.WriteLog(s.buf.Bytes(), s.count)
		s.buf.Reset()
	} else {
		// ensure queue is empty
		if s.MQ().IsClosed() && s.MQ().Length() == 0 {
			s.logger.Close()
		}
	}

}

// Stop goruntine
func (s *LogModule) Stop() error {
	s.MQ().Close()
	return s.logger.Close()
}

// Call async push message to queue
// func (s *LogModule) Call(option int32, obj *server.CallObject) error {
// 	return fmt.Errorf("TODO LogModule Call")
// }

// UserData module custom option, can you store you self value
func (s *LogModule) UserData() int32 {
	return s.logger.GetLevel()
}

// PushBytes async push string or bytes to queue, with option
func (s *LogModule) PushBytes(option int32, buf []byte) error {
	if err := s.MQ().Push(buf); err != nil {
		return err
	}
	return nil
}
