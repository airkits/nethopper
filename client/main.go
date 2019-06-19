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
// * @Date: 2019-06-06 13:21:31
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-06 13:21:31

package main

import (
	"context"
	"fmt"

	. "github.com/gonethopper/nethopper/server"
)

type S interface {
	MakeContext(parent S) S
	GetContext() context.Context
	ObserveChild(ch chan struct{})
	Run()
	IsClosed(ch chan struct{}) bool
	Close()
	OnExit()
}

type Base struct {
	ctx           context.Context
	cancel        context.CancelFunc
	childrenChans []chan struct{}
	quitChan      chan struct{}
	Name          string
}

func (a *Base) IsClosed(ch chan struct{}) bool {
	_, isClose := <-ch
	return isClose
}
func (a *Base) GetContext() context.Context {
	return a.ctx
}
func (a *Base) ObserveChild(ch chan struct{}) {
	a.childrenChans = append(a.childrenChans, ch)
}
func (a *Base) Close() {
	a.cancel()
}
func (a *Base) MakeContext(p S) S {
	a.childrenChans = make([]chan struct{}, 0)
	a.quitChan = make(chan struct{})
	if p == nil {
		a.ctx, a.cancel = context.WithCancel(context.Background())
	} else {
		a.ctx, a.cancel = context.WithCancel(p.GetContext())
		p.ObserveChild(a.quitChan)
	}
	return a
}
func (a *Base) OnExit() {
	for _, v := range a.childrenChans {
		if a.IsClosed(v) {
		}
	}
	fmt.Println(a.Name + " Exit")
	close(a.quitChan)
}

func (a *Base) Run() {
	for {
		select {
		case <-a.ctx.Done():
			a.OnExit()
			return
		default:
		}
	}
}

type Root struct {
	Base
}

func NewRoot(p S) S {
	return (&Root{Base{Name: "Root"}}).MakeContext(p)
}

type A struct {
	Base
}

func NewA(p S) S {
	return (&A{Base{Name: "A"}}).MakeContext(p)
}

type B struct {
	Base
}

func NewB(p S) S {
	return (&B{Base{Name: "B"}}).MakeContext(p)
}

type C struct {
	Base
}

func NewC(p S) S {
	return (&C{Base{Name: "C"}}).MakeContext(p)
}

type AA struct {
	Base
}

func NewAA(p S) S {
	return (&AA{Base{Name: "AA"}}).MakeContext(p)
}

type AB struct {
	Base
}

func NewAB(p S) S {
	return (&AB{Base{Name: "AB"}}).MakeContext(p)
}

type ABC struct {
	Base
}

func NewABC(p S) S {
	return (&ABC{Base{Name: "ABC"}}).MakeContext(p)
}
func main() {

	root := NewRoot(nil)
	a := NewA(root)
	b := NewB(root)
	c := NewC(root)
	aa := NewAA(a)
	ab := NewAB(a)
	abc := NewABC(ab)
	GO(root.Run)
	GO(a.Run)
	GO(b.Run)
	GO(c.Run)
	GO(aa.Run)
	GO(ab.Run)
	GO(abc.Run)
	root.Close()

	WG.Wait()
}
