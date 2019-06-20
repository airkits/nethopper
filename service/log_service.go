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

package service

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gonethopper/nethopper/log"
	. "github.com/gonethopper/nethopper/server"
)

// LogService struct implements the interface Service
type LogService struct {
	BaseService
	logger log.Log
	//q      queue.Queue
	id int32
	//for stat
	buf     bytes.Buffer
	count   int32
	msgSize int32
}

// LogServiceCreate log service create function
func LogServiceCreate() (Service, error) {
	return &LogService{}, nil
}

// Setup init and setup config
// Log config
// m := map[string]interface{}{
// 	"filename":    "server.log",
// 	"level":       7,
// 	"maxSize":     50,
// 	"maxLines":    1000,
// 	"hourEnabled": false,
// 	"dailyEnable": true,
//  "queueSize":1000,
// }
func (s *LogService) Setup(m map[string]interface{}) (Service, error) {
	queueSize, ok := m["queueSize"]
	if !ok {
		return nil, errors.New("params queueSize needed")
	}
	s.MakeContext(nil, int32(queueSize.(int)))

	logger, err := log.NewFileLogger(m)
	if err != nil {
		return nil, err
	}
	s.logger = logger
	SetLogLevel(logger.GetLevel())

	return s, nil
}

//ID service ID
func (s *LogService) ID() int32 {
	return s.id
}

//SetID set service id
func (s *LogService) SetID(v int32) {
	s.id = v
}

// Run create goruntine and run
func (s *LogService) Run(v ...interface{}) {

	for i := 0; i < 128; i++ {

		if v, err := s.Queue().AsyncPop(); err == nil {
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
		s.buf.Reset()
	} else {
		// ensure queue is empty
		if s.Queue().IsClosed() && s.Queue().Length() == 0 {
			s.logger.Close()
		}
	}

}

// Stop goruntine
func (s *LogService) Stop() error {
	return s.logger.Close()
}

// Send async send message to other goruntine
func (s *LogService) Send(msg *Message) error {
	return fmt.Errorf("TODO LogServer Send")
}

// SendBytes async send buffer to other goruntine
func (s *LogService) SendBytes(buf []byte) error {
	if err := s.Queue().Push(buf); err != nil {
		return err
	}
	return nil
}
