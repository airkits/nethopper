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
	"fmt"
	"time"

	"github.com/gonethopper/nethopper/cache/redis"
	"github.com/gonethopper/nethopper/codec"
	"github.com/gonethopper/nethopper/examples/simple_server/common"
	"github.com/gonethopper/nethopper/server"
)

// RedisService struct to define service
type RedisService struct {
	server.BaseContext
	rdb *redis.RedisCache
}

// RedisServiceCreate  service create function
func RedisServiceCreate() (server.Service, error) {
	return &RedisService{}, nil
}

// UserData service custom option, can you store you data and you must keep goruntine safe
func (s *RedisService) UserData() int32 {
	return 0
}

// Setup init custom service and pass config map to service
// config
// m := map[string]interface{}{
//  "queueSize":1000,
// }
func (s *RedisService) Setup(m map[string]interface{}) (server.Service, error) {

	cache, err := redis.NewRedisCache(m)
	if err != nil {
		return nil, err
	}
	s.rdb = cache

	return s, nil
}

//Reload reload config
func (s *RedisService) Reload(m map[string]interface{}) error {
	return nil
}

// OnRun goruntine run and call OnRun , always use ServiceRun to call this function
func (s *RedisService) OnRun(dt time.Duration) {
	for i := 0; i < 128; i++ {
		m, err := s.MQ().AsyncPop()
		if err != nil {
			break
		}
		message := m.(*server.Message)
		server.Info("%s receive one request message from mq,cmd = %s", s.Name(), message.Cmd)

		s.ProcessMessage(message)
	}
}

// ProcessMessage receive message from mq and process message
func (s *RedisService) ProcessMessage(message *server.Message) {
	cmd := message.Cmd
	if cmd == "login" {
		var v = make(map[string]interface{})
		server.Info("%s", string(message.Payload))
		if err := codec.JSONCodec.Unmarshal(message.Payload, &v, nil); err != nil {
			server.Info(err)
			return
		}
		password, err := s.rdb.GetString(s.Context(), fmt.Sprintf("uid_%d", v["uid"]))
		if err != nil {
			server.Info(err.Error())
			message.ErrCode = common.ErrorCodeRedisKeyNotExist
		} else {
			message.ErrCode = server.ErrorCodeOK
			message.Payload = []byte(password)
		}
		message.SrcID = s.ID()
		message.DestID = message.PopSeqID()
		message.MsgType = server.MTResponse

		server.SendMessage(message.DestID, 0, message)

	}
	// msgType := message.MsgType
	// switch msgType {
	// case server.MTRequest:
	// 	{
	// 		server.Info("receive message %s", message.Cmd)
	// 		message.SrcID = s.ID()

	// 		server.SendMessage(message.DestID, 0, message)
	// 		break
	// 	}
	// }
}

// Stop goruntine
func (s *RedisService) Stop() error {
	return nil
}

// PushMessage async send message to service
func (s *RedisService) PushMessage(option int32, msg *server.Message) error {
	if err := s.MQ().AsyncPush(msg); err != nil {
		server.Error(err.Error())
	}
	return nil

}

// PushBytes async send string or bytes to queue
func (s *RedisService) PushBytes(option int32, buf []byte) error {
	return nil
}