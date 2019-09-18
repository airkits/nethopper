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

package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/gonethopper/nethopper/server"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// HTTPServiceCreate  service create function
func HTTPServiceCreate() (server.Service, error) {
	return &HTTPService{}, nil
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

// HTTPService struct to define service
type HTTPService struct {
	server.BaseContext
	Address string
	router  *mux.Router
}

// UserData service custom option, can you store you data and you must keep goruntine safe
func (s *HTTPService) UserData() int32 {
	return 0
}

// Setup init custom service and pass config map to service
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *HTTPService) Setup(m map[string]interface{}) (server.Service, error) {
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
func (s *HTTPService) web() {
	if err := http.ListenAndServe(s.Address, s.router); err != nil {
		panic(err)
	}

}

// config map
// address default :80
func (s *HTTPService) readConfig(m map[string]interface{}) error {

	address, err := server.ParseValue(m, "address", ":11080")
	if err != nil {
		return err
	}
	s.Address = address.(string)

	return nil
}

//Reload reload config
func (s *HTTPService) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ServiceRun to call this function
func (s *HTTPService) OnRun(dt time.Duration) {
	for i := 0; i < 128; i++ {
		m, err := s.MQ().AsyncPop()
		if err != nil {
			break
		}

		message := m.(*server.Message)
		msgType := message.MsgType
		switch msgType {
		case server.MTRequest:
			{
				s.processRequest(message)
				break
			}
		case server.MTResponse:
			{
				s.processResponse(message)
				break
			}
		}
	}
}
func (s *HTTPService) processRequest(req *server.Message) {
	server.Info("%s receive one request message from mq,cmd = %s", s.Name(), req.Cmd)

}
func (s *HTTPService) processResponse(resp *server.Message) {
	server.Info("%s receive one response message from mq,cmd = %s", s.Name(), resp.Cmd)

	sess := server.GetSession(resp.SessionID)
	if sess != nil {
		sess.Response = resp
		sess.NotifyDone()
	}
}

// Stop goruntine
func (s *HTTPService) Stop() error {
	return nil
}

// PushMessage async send message to service
func (s *HTTPService) PushMessage(option int32, msg *server.Message) error {
	if err := s.MQ().AsyncPush(msg); err != nil {
		server.Error(err.Error())
	}
	return nil
}

// PushBytes async send string or bytes to queue
func (s *HTTPService) PushBytes(option int32, buf []byte) error {
	return nil
}
