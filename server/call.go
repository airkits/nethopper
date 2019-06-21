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
// * @Date: 2019-06-12 15:53:22
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-12 15:53:22

package server

import (
	"reflect"
)

// CallUserFunc simply to dynamically call a function or a method on an object
// Calls the callback given by the first parameter and passes the remaining parameters as arguments.
// Zero or more parameters to be passed to the callback.
// Returns the return value of the callback.
func CallUserFunc(f interface{}, v ...interface{}) []reflect.Value {
	valueFunc := reflect.ValueOf(f)
	paramsList := []reflect.Value{}
	if len(v) > 0 {
		for i := 0; i < len(v); i++ {
			paramsList = append(paramsList, reflect.ValueOf(v[i]))
		}
	}
	return valueFunc.Call(paramsList)

}

// CallUserMethod simply to dynamically call a method on an object
// Calls the instance given by the first parameter and method name as the second parameter
// and passes the remaining parameters as arguments.
// Zero or more parameters to be passed to the method.
// Returns the return value of the method.
func CallUserMethod(instance interface{}, method string, v ...interface{}) []reflect.Value {
	valueS := reflect.ValueOf(instance)
	m := valueS.MethodByName(method)
	paramsList := []reflect.Value{}
	if len(v) > 0 {
		for i := 0; i < len(v); i++ {
			paramsList = append(paramsList, reflect.ValueOf(v[i]))
		}
	}
	return m.Call(paramsList)
}

// GO wapper exec goruntine and stat count
func GO(v ...interface{}) {
	f := v[0]
	WG.Add(1)
	go func() {
		CallUserFunc(f, v[1:]...)
		WG.Done()
	}()
}
