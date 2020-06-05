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
// * @Date: 2020-01-09 11:01:34
// * @Last Modified by:   ankye
// * @Last Modified time: 2020-01-09 11:01:34

package logic

import (
	"github.com/gonethopper/nethopper/server"
)

// UUIDHandler get one uniq id
// @Summary UUIDHandler
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param channel query int32 1 "channel"
// @Success 200 {uint64} uint64 成功后返回值
// @Router /call/UUIDHandler [put]
func UUIDHandler(s *Module, obj *server.CallObject, channel int32) (uint64, error) {
	defer server.TraceCost("UUIDHandler")()
	//	opt, err := strconv.Atoi(uid)
	return s.GenerateUUID()
}

// UUIDsHandler get one uniq id
// @Summary UUIDHandler
// @Tags LogicModule
// @version 1.0
// @Accept  plain
// @Produce plain
// @Param channel query int64 1 "channel"
// @Success 200 {uint64} uint64 成功后返回值
// @Router /call/UUIDsHandler [put]
func UUIDsHandler(s *Module, obj *server.CallObject, channel int32, num int32) ([]uint64, error) {
	defer server.TraceCost("UUIDsHandler")()
	//	opt, err := strconv.Atoi(uid)
	uids := make([]uint64, num)
	var err error
	var uid uint64
	for i := 0; i < int(num); i++ {
		if uid, err = s.GenerateUUID(); err == nil {
			uids[i] = uid
		} else {
			return []uint64{}, err
		}
	}
	return uids, nil
}