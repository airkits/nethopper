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

package http

import (
	"net/http"
	"time"

	"github.com/gonethopper/nethopper/server"
	"github.com/julienschmidt/httprouter"
)

// WebHTTPServiceCreate  service create function
func WebHTTPServiceCreate() (server.Service, error) {
	return &WebHTTPService{}, nil
}

// WebHTTPService struct to define service
type WebHTTPService struct {
	server.BaseContext
	Address string
	router  *httprouter.Router
}

// UserData service custom option, can you store you data and you must keep goruntine safe
func (s *WebHTTPService) UserData() int32 {
	return 0
}

// Setup init custom service and pass config map to service
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *WebHTTPService) Setup(m map[string]interface{}) (server.Service, error) {
	if err := s.readConfig(m); err != nil {
		panic(err)
	}

	router := httprouter.New()
	s.router = router
	RegisterAPI(router)
	server.Info("http listening on:  %s", s.Address)

	server.GO(s.web)

	return s, nil
}
func (s *WebHTTPService) web() {
	if err := http.ListenAndServe(s.Address, s.router); err != nil {
		panic(err)
	}

}

// config map
// address default :80
func (s *WebHTTPService) readConfig(m map[string]interface{}) error {

	address, err := server.ParseValue(m, "address", ":11080")
	if err != nil {
		return err
	}
	s.Address = address.(string)

	return nil
}

//Reload reload config
func (s *WebHTTPService) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ServiceRun to call this function
func (s *WebHTTPService) OnRun(dt time.Duration) {

}

// Stop goruntine
func (s *WebHTTPService) Stop() error {
	return nil
}

// PushMessage async send message to service
func (s *WebHTTPService) PushMessage(option int32, msg *server.Message) error {
	return nil
}

// PushBytes async send string or bytes to queue
func (s *WebHTTPService) PushBytes(option int32, buf []byte) error {
	return nil
}
