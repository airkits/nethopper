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
)

var logger Service

// WG global goruntine wait group
var WG sync.WaitGroup

// App server instance
var App = &Server{
	GoCount:            0,
	ServiceCount:       0,
	ServiceAnonymousID: ServiceIDNamedMax,
	Services:           make(map[int]Service),
}

// Server server entity, only one instance
type Server struct {
	// GoCount total goruntine count
	sync.Mutex
	LogLevel           int
	GoCount            int
	ServiceCount       int
	ServiceAnonymousID int
	Services           map[int]Service
}

//SetLogLevel set log level to app global var
func (s *Server) SetLogLevel(level int) {
	s.Lock()
	defer s.Unlock()
	s.LogLevel = level
}

// GetServiceByID get service instance by id
func (s *Server) GetServiceByID(serviceID int) (Service, error) {
	s.Lock()
	defer s.Unlock()
	se, ok := s.Services[serviceID]
	if ok {
		return se, nil
	}
	return nil, fmt.Errorf("cant get service ID")
}

// CreateNamedService register named service
func (s *Server) CreateNamedService(serviceID int, service Service, m map[string]interface{}) (Service, error) {
	s.Lock()
	defer s.Unlock()
	se, err := service.Create(serviceID, m)
	if err != nil {
		return nil, err
	}
	s.Services[serviceID] = se
	s.ServiceCount++
	return se, nil
}

// CreateService register service
func (s *Server) CreateService(service Service, m map[string]interface{}) (Service, error) {
	s.Lock()
	defer s.Unlock()

	se, err := service.Create(s.ServiceAnonymousID, m)
	if err != nil {
		return nil, err
	}
	s.Services[s.ServiceAnonymousID] = se
	s.ServiceAnonymousID++
	s.ServiceCount++
	return se, nil
}

// RemoveService unregister service
func (s *Server) RemoveService(serviceID int) error {
	s.Lock()
	defer s.Unlock()
	se, ok := s.Services[serviceID]
	if !ok {
		return fmt.Errorf("cant exist serviceID %d", serviceID)
	}
	se.Stop()
	delete(s.Services, serviceID)
	s.ServiceCount--
	return nil
}
