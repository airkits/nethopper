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
// * @Date: 2019-12-18 10:47:08
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-18 10:47:08

package queue

//Queue interface define
type Queue interface {
	// Push 阻塞写队列
	Push(data interface{}) error
	// AsyncPush 异步非阻塞写队列
	AsyncPush(data interface{}) error
	// Pop 阻塞读队列
	Pop() (interface{}, error)
	// AsyncPop 异步读队列
	AsyncPop() (interface{}, error)
	// Capacity 队列大小
	Capacity() int32
	// Length 队列占用数量
	Length() int32
	// Close 关闭队列
	Close() error
	// IsClosed 队列是否已经关闭 关闭返回true,否则返回false
	IsClosed() bool
}
