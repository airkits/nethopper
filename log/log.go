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
// * @Date: 2019-06-05 23:08:35
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-05 23:08:35

package log

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/gonethopper/nethopper/utils"
)

//Log Interface
type Log interface {
	ParseConfig(v map[string]interface{}) error
	InitLogger() error
	SetLevel(level int32) error
	GetLevel() int32

	// Fatal system is unusable
	Fatal(v ...interface{}) error
	// Error error conditions
	Error(v ...interface{}) error
	// Warning warning conditions
	Warning(v ...interface{}) error
	// Info informational messages
	Info(v ...interface{}) error
	// Debug debug-level messages
	Debug(v ...interface{}) error
	// WriteLog write log to file, return immediately if not meet the conditions
	WriteLog(msg []byte, count int32) error

	CanLog(msgSize int32, count int32) bool

	Close() error
}

// Log Levels Define
// 0       Fatal: system is unusable
// 1       Error: error conditions
// 2       Warning: warning conditions
// 3       Info: informational messages
// 4       Debug: debug-level messages
const (
	FATAL = iota
	ERROR
	WARNING
	INFO
	DEBUG
)

// LogLevelPrefix level format to string
var LogLevelPrefix = [DEBUG + 1]string{" [FATAL] ", " [ERROR] ", " [WARNING] ", " [INFO] ", " [DEBUG] "}

// FormatLog format log and return string
// if len(v) > 1 ,format = v[0]
func FormatLog(level int32, v ...interface{}) string {
	if level < FATAL || level > DEBUG {
		level = FATAL
	}
	var buf bytes.Buffer
	buf.WriteString(utils.TimeYMDHIS())
	buf.WriteString(LogLevelPrefix[level])
	format := v[0].(string)

	if len(v) > 1 {
		buf.WriteString(fmt.Sprintf(format, v[1:]...))
	} else {
		buf.WriteString(format)
	}
	buf.WriteString("\n")
	return buf.String()
}

// ParseValue read config from map,if not exist return default value
func ParseValue(m map[string]interface{}, key string, opt interface{}) (interface{}, error) {
	value, ok := m[key]
	if !ok {
		return opt, nil
	}
	if reflect.TypeOf(value) != reflect.TypeOf(opt) {
		return nil, fmt.Errorf("config %s type failed", key)
	}
	return value, nil
}
