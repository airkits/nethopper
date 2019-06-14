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
)

var logger Service

// LogLevel log level setting
var LogLevel int32

//SetLogLevel set log level to app global var
func SetLogLevel(level int32) {
	LogLevel = atomic.LoadInt32(&level)
}

// AnonymousServiceID Anonymous Service Counter
var AnonymousServiceID int32 = ServiceIDNamedMax

// WG global goruntine wait group
var WG sync.WaitGroup

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

// GetServiceByID get service instance by id
func (s *Server) GetServiceByID(serviceID int32) (Service, error) {
	se, ok := s.Services.Load(serviceID)
	if ok {
		return se.(Service), nil
	}
	return nil, fmt.Errorf("cant get service ID")
}

// RegisterNamedService register named service
func (s *Server) RegisterNamedService(serviceID int32, se Service) (Service, error) {
	return s.registerServiceByID(serviceID, se)
}
func (s *Server) registerServiceByID(serviceID int32, se Service) (Service, error) {
	se.SetID(serviceID)
	s.Services.Store(serviceID, se)
	if serviceID == ServiceIDLog {
		logger = se
	}
	se.Start()
	return se, nil
}

// RegisterService register service
func (s *Server) RegisterService(se Service) (Service, error) {
	//Inc AnonymousServiceID count = count +1
	serviceID := atomic.AddInt32(&AnonymousServiceID, 1)
	return s.registerServiceByID(serviceID, se)
}

// RemoveService unregister service
func (s *Server) RemoveService(serviceID int32) error {
	se, err := s.GetServiceByID(serviceID)
	if err != nil {
		return err
	}
	s.Services.Delete(serviceID)
	se.Stop()

	return nil
}

//RemoveAllServices traversing services
func (s *Server) RemoveAllServices() {
	s.Services.Range(func(key interface{}, v interface{}) bool {
		s.RemoveService(key.(int32))
		return true
	})
}
