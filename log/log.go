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
	SetLevel(level int) error
	Emergency(v ...interface{}) error
	Alert(v ...interface{}) error
	Critical(v ...interface{}) error
	Error(v ...interface{}) error
	Warning(v ...interface{}) error
	Notice(v ...interface{}) error
	Info(v ...interface{}) error
	Debug(v ...interface{}) error
	QuitChan() <-chan struct{}
	Close() error
}

// Log Levels Define rfc5424  https://tools.ietf.org/html/rfc5424
// 0       Emergency: system is unusable
// 1       Alert: action must be taken immediately
// 2       Critical: critical conditions
// 3       Error: error conditions
// 4       Warning: warning conditions
// 5       Notice: normal but significant condition
// 6       Informational: informational messages
// 7       Debug: debug-level messages
const (
	EMEGENCY = iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

// LogLevelPrefix level format to string
var LogLevelPrefix = [DEBUG + 1]string{" [EMEGENCY] ", " [ALERT] ", " [CRITICAL] ", " [ERROR] ", " [WARNING] ", " [NOTICE] ", " [INFO] ", " [DEBUG] "}

// FormatLog format log and return string
// if len(v) > 1 ,format = v[0]
func FormatLog(level int, v ...interface{}) string {
	if level < EMEGENCY || level > DEBUG {
		level = EMEGENCY
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

	// msg := utils.TimeYMDHIS() + LogLevelPrefix[level]
	// if len(v) > 0 {
	// 	msg = msg + fmt.Sprintf(format, v...)
	// } else {
	// 	msg = msg + format
	// }
	// return msg + "\n"

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
