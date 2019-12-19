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
// * @Date: 2019-06-24 09:49:04
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-24 09:49:04

package server

import (
	"fmt"
	"reflect"
)

// ParseConfigValue read config from map,if not exist return default value,support string,int,bool
func ParseConfigValue(m map[string]interface{}, key string, opt interface{}, result interface{}) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("Invalid type %s", reflect.TypeOf(result))
	}

	value, ok := m[key]
	if !ok {
		value = opt
	}
	if reflect.TypeOf(value) != reflect.TypeOf(opt) {
		return fmt.Errorf("config %s type failed", key)
	}
	rv.Elem().Set(reflect.ValueOf(value).Convert(rv.Elem().Type()))
	return nil
}
