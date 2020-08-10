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
func getOpenID2UIDKey(openID string) string {
	return fmt.Sprintf("openid2uid_%s", openID)
}

//GetUIDByOpenID 通过openID映射为uid，没有就创建一个
func GetUIDByOpenID(s *Module, obj *server.CallObject, openID string) (uint64, error) {
	defer server.TraceCost("ConvertOpenID2UID")()
	key := getOpenID2UIDKey(openID)
	uid, err := s.rdb.GetUint64(s.Context(), key)
	return uid, err
}

// SetUIDByOpenID set openid to uid mapping
func SetUIDByOpenID(s *Module, obj *server.CallObject, openID string, uid uint64) (bool, error) {
	var key = getOpenID2UIDKey(openID)
	err := s.rdb.Set(s.Context(), key, uid, 0)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(s *Module, obj *server.CallObject, uid uint64) (*model.User, error) {
	defer server.TraceCost("GetUserInfo")()
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

// UpdateUserInfo update user info
func UpdateUserInfo(s *Module, obj *server.CallObject, u *model.User) (bool, error) {
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
