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

	"github.com/gonethopper/nethopper/examples/usercenter/model"
	"github.com/gonethopper/nethopper/server"
	"github.com/gonethopper/nethopper/utils/conv"
)

func getUserInfoKey(uid uint64) string {
	return fmt.Sprintf("userinfo_%d", uid)
}

// GetUserInfoHander 获取用户信息
func GetUserInfoHander(s *Module, obj *server.CallObject, uid uint64) (*model.User, error) {
	defer server.TraceCost("GetUserInfoHander")()
	key := getUserInfoKey(uid)
	results, err := s.rdb.HMGet(s.Context(), key, "uid", "appid", "openid", "uuid", "avatar", "name", "gender", "channel", "gold", "coin", "status")
	if err == nil {

		user := &model.User{
			UID:     conv.Str2Uint64(results["uid"]),
			AppID:   results["appid"],
			OpenID:  results["openid"],
			UUID:    results["uuid"],
			Name:    results["name"],
			Channel: results["channel"],
			Avatar:  results["avatar"],
			Gender:  conv.Str2Int(results["gender"]),
			Gold:    conv.Str2Uint64(results["gold"]),
			Coin:    conv.Str2Uint64(results["coin"]),
			Status:  conv.Str2Int(results["status"]),
		}
		return user, err
	}
	return nil, err
}

// UpdateUserInfoHandler update user info
func UpdateUserInfoHandler(s *Module, obj *server.CallObject, u *model.User) (bool, error) {

	var key = getUserInfoKey(u.UID)
	params := map[interface{}]interface{}{
		"uid":      u.UID,
		"appid":    u.AppID,
		"openid":   u.OpenID,
		"uuid":     u.UUID,
		"name":     u.Name,
		"channel":  u.Channel,
		"avatar":   u.Avatar,
		"password": u.Password,
		"phone":    u.Phone,
		"gender":   u.Gender,
		"age":      u.Age,
		"gold":     u.Gold,
		"coin":     u.Coin,
		"status":   u.Status,
		"loginAt":  u.LoginAt,
		"loginIP":  u.LoginIP,
		"createAt": u.CreateAt,
	}
	if err := s.rdb.HMSet(s.Context(), key, params); err != nil {
		return false, err
	}
	return true, nil
}
