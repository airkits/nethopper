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
// * @Date: 2019-06-14 14:40:51
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-14 14:40:51

package server

import (
	"github.com/gonethopper/nethopper/log"
)

//WriteLog send log to queue
func WriteLog(level int32, v ...interface{}) error {
	// UserData return logger level
	if level > GLoggerService.UserData() {
		return nil
	}
	msg := log.FormatLog(level, v...)
	if err := GLoggerService.SendBytes(level, []byte(msg)); err != nil {
		return err
	}
	return nil
}

// Emergency system is unusable
func Emergency(v ...interface{}) error {
	return WriteLog(log.EMEGENCY, v...)
}

// Alert action must be taken immediately
func Alert(v ...interface{}) error {
	return WriteLog(log.ALERT, v...)
}

// Critical critical conditions
func Critical(v ...interface{}) error {
	return WriteLog(log.CRITICAL, v...)
}

// Error error conditions
func Error(v ...interface{}) error {
	return WriteLog(log.ERROR, v...)
}

// Warning warning conditions
func Warning(v ...interface{}) error {
	return WriteLog(log.WARNING, v...)
}

// Notice normal but significant condition
func Notice(v ...interface{}) error {
	return WriteLog(log.NOTICE, v...)
}

// Info informational messages
func Info(v ...interface{}) error {
	return WriteLog(log.INFO, v...)
}

// Debug debug-level messages
func Debug(v ...interface{}) error {
	return WriteLog(log.DEBUG, v...)
}
