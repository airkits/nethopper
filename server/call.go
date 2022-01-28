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
	"fmt"
	"reflect"
	"runtime"
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
	App.ModifyGoCount(1)
	go func() {
		CallUserFunc(f, v[1:]...)
		App.ModifyGoCount(-1)
		WG.Done()
	}()
}

// GOWithContext wapper exec goruntine and use context to manager goruntine
func GOWithContext(v ...interface{}) {
	f := v[0]
	App.ModifyGoCount(1)
	go func() {
		CallUserFunc(f, v[1:]...)
		App.ModifyGoCount(-1)
	}()
}

// // Future async call function
// func Future(f func() (interface{}, error)) func() (interface{}, error) {
// 	var result interface{}
// 	var err error

// 	c := make(chan struct{}, 1)
// 	go func() {
// 		defer close(c)
// 		result, err = f()
// 	}()

// 	return func() (interface{}, error) {
// 		<-c
// 		return result, err
// 	}
// }

// RunFuncName 获取正在运行的函数名
func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// RunModuleFuncName 获取正在运行的module函数名
func RunModuleFuncName(s Module) string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return fmt.Sprintf("[%s] %s", s.Name(), f.Name())
}

// CallObject call struct
type CallObject struct {
	Cmd     string
	Option  int32
	Args    []interface{}
	ChanRet chan RetObject
}

type Callback func(interface{}, Ret)

//Ret define code and error
type Ret struct {
	Code int32
	Err  error
}

// RetObject call return object
type RetObject struct {
	Data interface{}
	Ret  Ret
}

// NewCallObject create call object
func NewCallObject(cmd string, opt int32, args ...interface{}) *CallObject {
	return &CallObject{
		Cmd:     cmd,
		Option:  opt,
		Args:    args,
		ChanRet: make(chan RetObject, 1),
	}
}

// AsyncCall async get data from modules,return call object
// same option value will run in same processor
func AsyncCall(destMID int32, cmd string, option int32, callback Callback, args ...interface{}) error {
	obj, err := processCall(destMID, cmd, option, args...)
	if err != nil {
		return err
	}
	go func(obj *CallObject) {
		result := <-obj.ChanRet
		callback(result.Data, result.Ret)
	}(obj)
	return err
}

func processCall(destMID int32, cmd string, option int32, args ...interface{}) (*CallObject, error) {
	m, err := GetModuleByID(destMID)
	if err != nil {
		return nil, err
	}
	var obj = NewCallObject(cmd, option, args...)
	if err = m.Call(option, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// Call sync get data from modules
// same option value will run in same processor
func Call(destMID int32, cmd string, option int32, args ...interface{}) (interface{}, Ret) {
	obj, err := processCall(destMID, cmd, option, args...)
	if err != nil {
		return nil, Ret{Code: -1, Err: err}
	}
	result := <-obj.ChanRet
	return result.Data, result.Ret
}
