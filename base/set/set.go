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
// * @Date: 2019-12-20 19:40:23
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-20 19:40:23

package set

type Set interface {
	Add(e interface{}) bool
	Remove(e interface{})
	Clear()
	Contains(e interface{}) bool
	Len() int
	Same(other Set) bool
	Elements() []interface{}
	String() string
}

// 判断集合 one 是否是集合 other 的超集
func IsSuperset(one Set, other Set) bool {
	if one == nil || other == nil {
		return false
	}
	oneLen := one.Len()
	otherLen := other.Len()
	if oneLen == 0 || oneLen == otherLen {
		return false
	}
	if oneLen > 0 && otherLen == 0 {
		return true
	}
	for _, v := range other.Elements() {
		if !one.Contains(v) {
			return false
		}
	}
	return true
}

// 生成集合 one 和集合 other 的并集
func Union(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	unionedSet := NewSimpleSet()
	for _, v := range one.Elements() {
		unionedSet.Add(v)
	}
	if other.Len() == 0 {
		return unionedSet
	}
	for _, v := range other.Elements() {
		unionedSet.Add(v)
	}
	return unionedSet
}

// 生成集合 one 和集合 other 的交集
func Intersect(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	intersectedSet := NewSimpleSet()
	if other.Len() == 0 {
		return intersectedSet
	}
	if one.Len() < other.Len() {
		for _, v := range one.Elements() {
			if other.Contains(v) {
				intersectedSet.Add(v)
			}
		}
	} else {
		for _, v := range other.Elements() {
			if one.Contains(v) {
				intersectedSet.Add(v)
			}
		}
	}
	return intersectedSet
}

// 生成集合 one 对集合 other 的差集
func Difference(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	differencedSet := NewSimpleSet()
	if other.Len() == 0 {
		for _, v := range one.Elements() {
			differencedSet.Add(v)
		}
		return differencedSet
	}
	for _, v := range one.Elements() {
		if !other.Contains(v) {
			differencedSet.Add(v)
		}
	}
	return differencedSet
}

// 生成集合 one 和集合 other 的对称差集
func SymmetricDifference(one Set, other Set) Set {
	if one == nil || other == nil {
		return nil
	}
	diffA := Difference(one, other)
	if other.Len() == 0 {
		return diffA
	}
	diffB := Difference(other, one)
	return Union(diffA, diffB)
}

func NewSimpleSet() Set {
	return NewHashSet()
}

func IsSet(value interface{}) bool {
	if _, ok := value.(Set); ok {
		return true
	}
	return false
}
