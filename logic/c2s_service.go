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
// * @Date: 2019-06-24 11:07:19
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-24 11:07:19

package logic

import "github.com/gonethopper/nethopper/server"

// C2SService struct to define service
type C2SService struct {
	server.BaseContext
}

// C2SServiceCreate  service create function
func C2SServiceCreate() (server.Service, error) {
	return &C2SService{}, nil
}

// UserData service custom option, can you store you data and you must keep goruntine safe
func (s *C2SService) UserData() int32 {
	return 0
}

// Setup init custom service and pass config map to service
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *C2SService) Setup(m map[string]interface{}) (server.Service, error) {
	return s, nil
}

//Reload reload config
func (s *C2SService) Reload(m map[string]interface{}) error {
	return nil
}

// Run create goruntine and run, always use ServiceRun to call this function
func (s *C2SService) Run() {

}

// Stop goruntine
func (s *C2SService) Stop() error {
	return nil
}

// SendMessage async send message to service
func (s *C2SService) SendMessage(option int32, msg *server.Message) error {
	return nil
}

// SendBytes async send string or bytes to queue
func (s *C2SService) SendBytes(option int32, buf []byte) error {
	return nil
}
