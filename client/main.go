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
	Create(ctx context.Context)
	Close()
	GetContext() context.Context
	Run()
}
type Root struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *Root) Create(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}
func (a *Root) Close() {
	fmt.Println("Root")
	a.cancel()
}
func (a *Root) GetContext() context.Context {
	return a.ctx
}
func (a *Root) Run() {
	for {
		select {
		case <-a.ctx.Done():
			fmt.Println("Root程序结束")
			return
		default:
		}
	}
}

type A struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *A) Create(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}
func (a *A) Close() {
	fmt.Println("A")
	a.cancel()
}
func (a *A) GetContext() context.Context {
	return a.ctx
}
func (a *A) Run() {
	for {
		select {
		case <-a.ctx.Done():
			fmt.Println("A程序结束")
			return
		default:
		}
	}
}

type B struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *B) Create(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}
func (a *B) Close() {
	fmt.Println("B")
	a.cancel()
}
func (a *B) GetContext() context.Context {
	return a.ctx
}
func (a *B) Run() {
	for {
		select {
		case <-a.ctx.Done():
			fmt.Println("B程序结束")
			return
		default:
		}
	}
}

type C struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *C) Create(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}
func (a *C) Close() {
	fmt.Println("C")
	a.cancel()
}
func (a *C) GetContext() context.Context {
	return a.ctx
}
func (a *C) Run() {
	for {
		select {
		case <-a.ctx.Done():
			fmt.Println("C程序结束")
			return
		default:
		}
	}
}

type AA struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *AA) Create(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}
func (a *AA) Close() {
	fmt.Println("AA")
	a.cancel()
}
func (a *AA) GetContext() context.Context {
	return a.ctx
}
func (a *AA) Run() {
	for {
		select {
		case <-a.ctx.Done():
			fmt.Println("AA程序结束")
			return
		default:
		}
	}
}

type AB struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *AB) Create(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}
func (a *AB) Close() {
	fmt.Println("AB")
	a.cancel()
}
func (a *AB) GetContext() context.Context {
	return a.ctx
}
func (a *AB) Run() {
	for {
		select {
		case <-a.ctx.Done():
			fmt.Println("AB程序结束")
			return
		default:
		}
	}
}

type ABC struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (a *ABC) Create(ctx context.Context) {
	a.ctx, a.cancel = context.WithCancel(ctx)
}
func (a *ABC) Close() {
	fmt.Println("ABC")
	a.cancel()
}
func (a *ABC) GetContext() context.Context {
	return a.ctx
}
func (a *ABC) Run() {
	for {
		select {
		case <-a.ctx.Done():
			fmt.Println("ABC程序结束")
			return
		default:
		}
	}
}
func main() {

	root := &Root{}
	root.Create(context.Background())

	a := &A{}
	a.Create(root.GetContext())
	b := &B{}
	b.Create(root.GetContext())
	c := &C{}
	c.Create(root.GetContext())
	aa := &AA{}
	aa.Create(a.GetContext())
	ab := &AB{}
	ab.Create(a.GetContext())
	abc := &ABC{}
	abc.Create(ab.GetContext())
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
