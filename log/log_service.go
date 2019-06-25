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
	"fmt"

	"github.com/gonethopper/nethopper/server"
)

// LogService struct implements the interface Service
type LogService struct {
	server.BaseContext
	logger  server.Log
	console server.Log
	//for stat
	buf     bytes.Buffer
	count   int32
	msgSize int32
}

// LogServiceCreate log service create function
func LogServiceCreate() (server.Service, error) {
	return &LogService{}, nil
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
func (s *LogService) Setup(m map[string]interface{}) (server.Service, error) {

	logger, err := NewFileLogger(m)
	if err != nil {
		return nil, err
	}
	s.logger = logger
	console, err := NewConsoleLogger(m)
	if err != nil {
		return nil, err
	}
	s.console = console
	return s, nil
}

// Reload reload config from map
func (s *LogService) Reload(m map[string]interface{}) error {
	level, err := server.ParseValue(m, "level", 7)
	if err != nil {
		return err
	}
	return s.logger.SetLevel(int32(level.(int)))

}

// Run create goruntine and run
func (s *LogService) Run() {
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
func (s *LogService) Stop() error {
	s.MQ().Close()
	return s.logger.Close()
}

// SendMessage async push message to queue
func (s *LogService) SendMessage(option int32, msg *server.Message) error {
	return fmt.Errorf("TODO LogService SendMessage")
}

// UserData service custom option, can you store you self value
func (s *LogService) UserData() int32 {
	return s.logger.GetLevel()
}

// SendBytes async push string or bytes to queue, with option
func (s *LogService) SendBytes(option int32, buf []byte) error {
	if err := s.MQ().Push(buf); err != nil {
		return err
	}
	return nil
}
