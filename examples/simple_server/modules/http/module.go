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

// ModuleCreate  module create function
func ModuleCreate() (server.Module, error) {
	return &Module{}, nil
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
			sess := server.CreateSession(server.ModuleIDHTTP, host, port)
			server.Info("new connection from:%v port:%v sessionid:%s", host, port, sess.SessionID)
			context.Set(req, "token", sess.SessionID)
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, req)
	})
}

// Module struct to define module
type Module struct {
	server.BaseContext
	Address string
	router  *mux.Router
}

// UserData module custom option, can you store you data and you must keep goruntine safe
// func (s *Module) UserData() int32 {
// 	return 0
// }

// Setup init custom module and pass config map to module
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *Module) Setup(m map[string]interface{}) (server.Module, error) {
	if err := s.ReadConfig(m); err != nil {
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
func (s *Module) web() {
	if err := http.ListenAndServe(s.Address, s.router); err != nil {
		panic(err)
	}

}

// ReadConfig config map
// address default :80
func (s *Module) ReadConfig(m map[string]interface{}) error {
	if err := server.ParseConfigValue(m, "address", ":11080", &s.Address); err != nil {
		return err
	}
	return nil
}

//Reload reload config
// func (s *Module) Reload(m map[string]interface{}) error {
// 	return nil
// }

// OnRun goruntine run and call OnRun , always use ModuleRun to call this function
func (s *Module) OnRun(dt time.Duration) {
}

// Stop goruntine
func (s *Module) Stop() error {
	return nil
}

// // Call async send message to module
// func (s *Module) Call(option int32, obj *server.CallObject) error {
// 	return nil
// }

// PushBytes async send string or bytes to queue
// func (s *Module) PushBytes(option int32, buf []byte) error {
// 	return nil
// }
