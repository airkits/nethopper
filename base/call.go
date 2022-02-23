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

package base

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
)

// CallFunction simply to dynamically call a function or a method on an object
// Calls the callback given by the first parameter and passes the remaining parameters as arguments.
// Zero or more parameters to be passed to the callback.
// Returns the return value of the callback.
func CallFunction(f interface{}, v ...interface{}) []reflect.Value {
	valueFunc := reflect.ValueOf(f)
	paramsList := []reflect.Value{}
	if len(v) > 0 {
		for i := 0; i < len(v); i++ {
			paramsList = append(paramsList, reflect.ValueOf(v[i]))
		}
	}
	return valueFunc.Call(paramsList)

}

// CallMethod simply to dynamically call a method on an object
// Calls the instance given by the first parameter and method name as the second parameter
// and passes the remaining parameters as arguments.
// Zero or more parameters to be passed to the method.
// Returns the return value of the method.
func CallMethod(instance interface{}, method string, v ...interface{}) []reflect.Value {
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

// GOFunctionWithWG wapper exec goruntine,waitgroup and ref count
func GOFunctionWithWG(wg *sync.WaitGroup, ref IRef, v ...interface{}) {
	f := v[0]
	wg.Add(1)
	if ref != nil {
		ref.AddRef()
	}
	go func(wg *sync.WaitGroup, ref IRef, v ...interface{}) {
		defer wg.Done()
		if ref != nil {
			defer ref.DecRef()
		}
		CallFunction(f, v[1:]...)

	}(wg, ref, v...)
}

// GOFunction wapper exec goruntine and ref count
func GOFunction(ref IRef, v ...interface{}) {
	f := v[0]
	if ref != nil {
		ref.AddRef()
	}
	go func(ref IRef, v ...interface{}) {
		if ref != nil {
			defer ref.DecRef()
		}
		CallFunction(f, v[1:]...)
	}(ref, v...)
}

// GOMethodWithWG wapper exec goruntine,waitgroup and ref count
func GOMethodWithWG(wg *sync.WaitGroup, ref IRef, instance interface{}, method string, v ...interface{}) {
	wg.Add(1)
	if ref != nil {
		ref.AddRef()
	}
	go func(wg *sync.WaitGroup, ref IRef, instance interface{}, method string, v ...interface{}) {
		defer wg.Done()
		if ref != nil {
			defer ref.DecRef()
		}
		CallMethod(instance, method, v...)
	}(wg, ref, instance, method, v...)
}

// GOMethod wapper exec class method goruntine and ref count
func GOMethod(ref IRef, instance interface{}, method string, v ...interface{}) {
	if ref != nil {
		ref.AddRef()
	}
	go func(ref IRef, instance interface{}, method string, v ...interface{}) {
		if ref != nil {
			defer ref.DecRef()
		}
		CallMethod(instance, method, v...)
	}(ref, instance, method, v...)
}

// GO wapper exec goruntine and ref count
func GO(v ...interface{}) {
	f := v[0]
	go func(v ...interface{}) {
		CallFunction(f, v[1:]...)
	}(v...)
}

// Future async call function
func Future(f func() (interface{}, error)) func() (interface{}, error) {
	var result interface{}
	var err error

	c := make(chan struct{}, 1)
	go func() {
		defer close(c)
		result, err = f()
	}()

	return func() (interface{}, error) {
		<-c
		return result, err
	}
}

// GetMethodName 获取正在运行的函数名
func GetFunctionName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// GetClassMethodName 获取正在运行的class函数名
func GetClassMethodName(s IClass) string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return fmt.Sprintf("[%s] %s", s.Name(), f.Name())
}
