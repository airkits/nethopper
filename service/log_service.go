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
	"fmt"

	"github.com/gonethopper/nethopper/log"
	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/queue"
)

// LogService struct implements the interface Service
type LogService struct {
	logger log.Log
	q      queue.Queue
	id     int
}

// Create instance
func (s *LogService) Create(serviceID int, m map[string]interface{}) (server.Service, error) {
	s = &LogService{
		id: serviceID,
	}
	s.q = queue.NewChanQueue(m["queueSize"].(int))
	// m := map[string]interface{}{
	// 	"filename":    "server.log",
	// 	"level":       7,
	// 	"maxSize":     50,
	// 	"maxLines":    1000,
	// 	"hourEnabled": false,
	// 	"dailyEnable": true,
	//  "queueSize":1000,
	// }
	logger, err := log.NewFileLogger(m)
	if err != nil {
		return nil, err
	}
	s.logger = logger
	server.App.SetLogLevel(logger.GetLevel())

	return s, nil
}

//ID service ID
func (s *LogService) ID() int {
	return s.id
}

// Start create goruntine and run
func (s *LogService) Start(m map[string]interface{}) error {
	server.GO(s.logger.RunLogger)
	return nil
}

// Stop goruntine
func (s *LogService) Stop() error {
	return s.logger.Close()
}

// Send async send message to other goruntine
func (s *LogService) Send(msg *server.Message) error {
	return fmt.Errorf("TODO LogServer Send")
}

// SendBuffer async send buffer to other goruntine
func (s *LogService) SendBuffer(buf []byte) error {
	return s.logger.AsyncWrite(buf)
}
