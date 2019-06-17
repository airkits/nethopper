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
// * @Date: 2019-06-14 14:15:06
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-14 14:15:06

package server

import (
	"fmt"
	"sync/atomic"
)

const (
	// ServiceNamedIDs service id define, system reserved 1-63
	ServiceNamedIDs = iota
	// ServiceIDMain main goruntinue
	ServiceIDMain
	// ServiceIDMonitor server monitor service
	ServiceIDMonitor
	//ServiceIDLog log service
	ServiceIDLog
	//ServiceIDUserCustom User custom define named services from 64-128
	ServiceIDUserCustom = 64
	//ServiceIDNamedMax named services max ID
	ServiceIDNamedMax = 128
)

// Service interface define
type Service interface {

	// Setup
	Setup(m map[string]interface{}) (Service, error)
	// ID service id
	ID() int32
	//SetID set service ID
	SetID(v int32)
	// Start create goruntine and run
	Start() error
	// Stop goruntine
	Stop() error
	// Send async send message to other goruntine
	Send(msg *Message) error
	SendBytes(buf []byte) error
}

// RegisterService register service name to create function mapping
func RegisterService(name string, createFunc func() (Service, error)) error {
	if _, ok := refServices[name]; ok {
		return fmt.Errorf("Already register Service %s", name)
	}
	refServices[name] = createFunc
	return nil
}

// CreateService create service by name
func CreateService(name string) (Service, error) {
	if f, ok := refServices[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("You need register Service %s first", name)
}

// GetServiceByID get service instance by id
func GetServiceByID(serviceID int32) (Service, error) {
	se, ok := App.Services.Load(serviceID)
	if ok {
		return se.(Service), nil
	}
	return nil, fmt.Errorf("cant get service ID")
}

// NewNamedService create named service
func NewNamedService(serviceID int32, name string, m map[string]interface{}) (Service, error) {
	return createServiceByID(serviceID, name, m)
}
func createServiceByID(serviceID int32, name string, m map[string]interface{}) (Service, error) {
	se, err := CreateService(name)
	if err != nil {
		return nil, err
	}
	se.Setup(m)
	se.SetID(serviceID)
	App.Services.Store(serviceID, se)
	if serviceID == ServiceIDLog {
		logger = se
	}
	se.Start()
	return se, nil
}

// NewService create anonymous service
func NewService(name string, m map[string]interface{}) (Service, error) {
	//Inc AnonymousServiceID count = count +1
	serviceID := atomic.AddInt32(&AnonymousServiceID, 1)
	return createServiceByID(serviceID, name, m)
}

// DeleteService unregister service
func DeleteService(serviceID int32) error {
	se, err := GetServiceByID(serviceID)
	if err != nil {
		return err
	}
	App.Services.Delete(serviceID)
	se.Stop()

	return nil
}

// DeleteAllServices traversing services
func DeleteAllServices() {
	App.Services.Range(func(key interface{}, v interface{}) bool {
		DeleteService(key.(int32))
		return true
	})
}

// SendMessage send message to services
func SendMessage(serviceID int32, msg *Message) error {
	s, err := GetServiceByID(serviceID)
	if err != nil {
		return err
	}
	return s.Send(msg)
}
