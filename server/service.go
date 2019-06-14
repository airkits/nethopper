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
	// ID service id
	ID() int
	//Create instance
	Create(serviceID int, m map[string]interface{}) (Service, error)
	// Start create goruntine and run
	Start(m map[string]interface{}) error
	// Stop goruntine
	Stop() error
	// Send async send message to other goruntine
	Send(msg *Message) error
	SendBuffer(buf []byte) error
}
