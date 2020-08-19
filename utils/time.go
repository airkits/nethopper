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
// * @Date: 2019-12-18 10:46:52
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-12-18 10:46:52

package utils

import (
	"fmt"
	"time"
)

// TimeYMDHIS get current time
// return format yearmouthday hour:minute:second
func TimeYMDHIS() string {
	return time.Now().Format("20060102 15:04:05")
}

// TimeYMDH get current time
// return format yearmouthday-hour
func TimeYMDH() string {
	return time.Now().Format("20060102-15")
}

// TimeYMD get current time
// return format yearmouthday
func TimeYMD() string {
	return time.Now().Format("20060102")
}

//LocalMilliscond 当前毫秒
func LocalMilliscond() int64 {
	return time.Now().UnixNano() / 1e6
}

//LocalTimestamp 当前时间戳秒
func LocalTimestamp() int64 {
	return time.Now().Unix()
}

//GetTodayTime 获取当天开始时间
func GetTodayTime() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

//GetTodayHourTime 获取当天整点时间
func GetTodayHourTime(hour int) time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), hour, 0, 0, 0, t.Location())
}

//GetEmptyTime 获取空时间结构
func GetEmptyTime() time.Time {
	return time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
}

//GetTomorrowTime 获取明天的开始时间
func GetTomorrowTime() time.Time {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return tm1.Add(24 * time.Hour)
}

//Today 获取当前年月日
func Today() string {

	year, month, day := time.Now().Date()
	str := fmt.Sprintf("%d-%d-%d 00:00:00", year, month, day)
	return str
}

//HourNow 现在小时字符串
func HourNow() string {
	year, month, day := time.Now().Date()
	str := fmt.Sprintf("%d-%d-%d %d:00:00", year, month, day, time.Now().Hour())
	return str
}

//Tomorrow 获取明天年月日
func Tomorrow() string {
	tomorrow := time.Now().Add(24 * time.Hour)
	year, month, day := tomorrow.Date()
	str := fmt.Sprintf("%d-%d-%d 00:00:00", year, month, day)
	return str
}

//TodayTimestamp 获取今天最早的时间戳毫秒
func TodayTimestamp() int64 {
	now := time.Now()
	year, month, day := now.Date()
	todaytime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return todaytime.UnixNano() / 1e6
}

//TimeFormat 时间戳格式化
func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

//HourTimestamp 获取整点时间戳
func HourTimestamp() int64 {
	now := time.Now()
	timestamp := now.Unix() - int64(now.Second()) - int64((60 * now.Minute()))

	return timestamp
}

//Hour 获取整点时间戳
func Hour(t time.Time) int64 {
	timestamp := t.Unix() - int64(t.Second()) - int64((60 * t.Minute()))
	return timestamp
}
