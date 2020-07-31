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
// * @Date: 2020-01-09 11:01:43
// * @Last Modified by:   ankye
// * @Last Modified time: 2020-01-09 11:01:43

package redis

import (
	"fmt"

	"github.com/gonethopper/nethopper/server"
)

// GetUserInfoHander 获取用户信息
// func GetUserInfoHander(s *Module, obj *server.CallObject, uid string) (*model.User, error) {
// 	defer server.TraceCost("GetUserInfoHander")()
// 	key := fmt.Sprintf("uid_%s", uid)
// 	results, err ;= s.rdb.HMGet(s.Context(),key,"uid", "appid", "openid", "uuid", "avatar", "name", "gender", "channel", "gold", "coin", "status")

// 	return password, err

// }

// UpdateUserInfoHandler update user info
func UpdateUserInfoHandler(s *Module, obj *server.CallObject, uid string, pwd string) (bool, error) {

	var key = fmt.Sprintf("uid_%s", uid)
	err := s.rdb.Set(s.Context(), key, pwd, 0)
	if err != nil {
		return false, err
	}
	return true, nil
}
