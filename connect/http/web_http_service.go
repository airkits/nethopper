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
	"strings"
	"time"

	"github.com/gonethopper/nethopper/server"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// WebHTTPServiceCreate  service create function
func WebHTTPServiceCreate() (server.Service, error) {
	return &WebHTTPService{}, nil
}

// SessionHTTPMiddleware define http middleware to create session id
func SessionHTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Do stuff here
		// log.Println(r.RequestURI)
		authStr := req.Header.Get("Authorization")
		if authStr == "" {
			arr := strings.Split(req.RemoteAddr, ":")
			host := arr[0]
			port := arr[1]
			sess := server.CreateSession(server.ServiceIDHTTP, host, port)
			server.Info("new connection from:%v port:%v sessionid:%s", host, port, sess.SessionID)
			context.Set(req, "token", sess.SessionID)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, req)
	})
}

// WebHTTPService struct to define service
type WebHTTPService struct {
	server.BaseContext
	Address string
	router  *mux.Router
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

	router := mux.NewRouter()
	s.router = router
	RegisterAPI(router)
	server.Info("http listening on:  %s", s.Address)
	router.Use(SessionHTTPMiddleware)

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
