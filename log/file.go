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
// * @Date: 2019-06-06 13:49:42
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-06 13:49:42

package log

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gonethopper/nethopper/utils"
	"github.com/gonethopper/queue"
)

//NewFileLogger create FileLog instance
func NewFileLogger(m map[string]interface{}) (Log, error) {
	logger := &FileLog{
		closedChan: make(chan struct{}),
	}
	if err := logger.ParseConfig(m); err != nil {
		return nil, err
	}
	if err := logger.InitLogger(); err != nil {
		return nil, err
	}
	return logger, nil
}

//FileLog implements Log interface
//write log to file,if file reached limit,rename file match format filename
//support filesize limit / time frequency / lines limit
//filename read from config like server.log and the real filename like server_20190101-01.log or server_20190101.log
//hourEnabled if enabled , filename like server_20190101-01.log
//dailyEnabled if enabled, filename like server_20190101.log
//else filename like server.log
//filename format = filename_ymd-h_num.suffix
type FileLog struct {
	//set level and  atomic incr CurrentSize and CurrentLines
	//write log one by one
	sync.Mutex
	level         int
	fileName      string //real filename
	currentTime   string //gen date ymd / ymd-h
	suffix        string //filename suffix,like .log .txt
	prefix        string //filename prefix, like server
	maxSize       int    //filesize limit
	currentSize   int
	maxLines      int //lines limit
	currentLines  int
	currentNum    int  //current file nums
	hourEnabled   bool //time frequency
	dailyEnabled  bool
	currentWriter *os.File //current File Writer
	q             queue.Queue
	closedChan    chan struct{}
}

//InitLogger init logger
func (l *FileLog) InitLogger() error {
	l.q = queue.NewChanQueue(1024)

	go l.startLogger()
	return l.createNewFile()
}

//QuitChan write all message from queue and tigger closed notify
func (l *FileLog) QuitChan() <-chan struct{} {
	return l.closedChan
}
func (l *FileLog) startLogger() {
	var buf bytes.Buffer
	var count int
	var msgSize int
	for {
		count = 0
		msgSize = 0
		for i := 0; i < 128; i++ {
			if v, err := l.q.AsyncGet(); err == nil {
				if n, e := buf.WriteString(v.(string)); e == nil {
					msgSize += n
					count++
					if (msgSize+l.currentSize) >= l.maxSize || (count+l.currentLines) >= l.maxLines {
						break
					}
				}
			}
		}

		if buf.Len() > 0 {
			l.writeLog(buf.String(), count)
			buf.Reset()
		} else {
			if l.q.IsClosed() && l.q.Length() == 0 {
				l.flush()
				l.currentWriter.Close()
				close(l.closedChan)
				return
			}
		}
	}
}

//writeLog write message to file, return immediately if not meet the conditions
func (l *FileLog) writeLog(msg string, count int) error {

	// l.Lock()
	// defer l.Unlock()
	if l.fileCutTest() {
		l.moveFile()
		l.createNewFile()
	}
	_, err := l.currentWriter.Write([]byte(msg))
	if err == nil {
		l.currentLines += count
		l.currentSize += len(msg)
	}

	return nil
}

//SetLevel update log level
func (l *FileLog) SetLevel(level int) error {
	if level < EMEGENCY || level > DEBUG {
		return fmt.Errorf("log level:[%d] invalid", level)
	}
	l.Lock()
	defer l.Unlock()
	l.level = level
	return nil
}

//ParseConfig read config from map[string]interface{}
// config key
// filename default server.log
// level default 7
// maxSize default 1024
// maxLines default 100000
// hourEnabled default false
// dailyEnabled default true
func (l *FileLog) ParseConfig(m map[string]interface{}) error {

	filename, err := ParseValue(m, "filename", "server.log")
	if err != nil {
		return err
	}
	l.suffix = filepath.Ext(filename.(string))
	l.prefix = strings.TrimSuffix(filename.(string), l.suffix)
	level, err := ParseValue(m, "level", 7)
	if err != nil {
		return err
	}
	l.level = level.(int)

	maxSize, err := ParseValue(m, "maxSize", 1024)
	if err != nil {
		return err
	}
	l.maxSize = maxSize.(int) * 1024 * 1024
	maxLines, err := ParseValue(m, "maxLines", 100000)
	if err != nil {
		return err
	}
	l.maxLines = maxLines.(int)
	hourEnabled, err := ParseValue(m, "hourEnabled", false)
	if err != nil {
		return err
	}
	l.hourEnabled = hourEnabled.(bool)
	dailyEnabled, err := ParseValue(m, "dailyEnabled", true)
	if err != nil {
		return err
	}
	l.dailyEnabled = dailyEnabled.(bool)

	return nil
}
func (l *FileLog) genCurrentTime() string {
	currentTime := ""
	if l.hourEnabled {
		currentTime = utils.TimeYMDH()
	} else if l.dailyEnabled {
		currentTime = utils.TimeYMD()
	}
	return currentTime
}

//genFilename filename format = filename_ymd-h_num.suffix
//if num == 0, then format = filename_ymd-h.suffix
//else if hourEnabled == false, then format = filename_ymd.suffix
//else if  dailyEnabled == false, then format = filename.suffix
func (l *FileLog) genFilename(timestr string, num int) string {
	var buf bytes.Buffer
	buf.WriteString(l.prefix)
	if len(timestr) > 0 {
		buf.WriteString("_")
		buf.WriteString(timestr)
	}
	if num > 0 {
		buf.WriteString("_")
		buf.WriteString(strconv.Itoa(num))
	}
	buf.WriteString(l.suffix)
	return buf.String()
	// filename := l.prefix
	// if len(timestr) > 0 {
	// 	filename += "_" + timestr
	// }
	// if num > 0 {
	// 	filename += "_" + strconv.Itoa(num)
	// }
	// return filename + l.suffix
}

//nextNumTest test the file actually exists in the filesystem,return the next file num
//if test failed return -1
func (l *FileLog) nextNumTest(timestr string) int {
	MaxNum := 1000
	for i := 1; i < MaxNum; i++ {
		filename := l.genFilename(timestr, i)
		if !utils.FileIsExist(filename) {
			return i
		}
	}
	return -1
}

//fileCutTest check time/maxsize/maxlines
func (l *FileLog) fileCutTest() bool {

	timestr := l.genCurrentTime()
	if strings.Compare(l.currentTime, timestr) != 0 {
		return true
	}
	if l.currentSize >= l.maxSize || l.currentLines >= l.maxLines {
		return true
	}
	return false
}

//Close close file logger
func (l *FileLog) Close() error {
	return l.q.Close()
}

func (l *FileLog) flush() error {
	return l.currentWriter.Sync()
}

func (l *FileLog) moveFile() error {
	num := l.nextNumTest(l.currentTime)
	if num < 0 {
		return fmt.Errorf("max file reached")
	}
	l.currentWriter.Sync()
	l.currentWriter.Close()
	l.currentWriter = nil
	filename := l.genFilename(l.currentTime, num)
	err := os.Rename(l.fileName, filename)
	if err != nil {
		return fmt.Errorf("file rename Error %s", l.fileName)
	}
	return nil
}

//createNewFile if file exist,then check current lines and filesize
func (l *FileLog) createNewFile() error {
	l.currentTime = l.genCurrentTime()
	l.fileName = l.genFilename(l.currentTime, 0)
	flag := utils.FileIsExist(l.fileName)
	if flag { //exist file
		lines, err := utils.FileLines(l.fileName)
		if err != nil {
			return err
		}
		l.currentLines = lines
	} else {
		l.currentLines = 0
	}
	fd, err := os.OpenFile(l.fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	l.currentWriter = fd
	if flag {
		stat, err := fd.Stat()
		if err != nil {
			return err
		}
		l.currentSize = int(stat.Size())
	} else {
		l.currentSize = 0
	}
	return nil
}

//PushLog push log to queue
func (l *FileLog) PushLog(level int, v ...interface{}) error {
	if level > l.level {
		return nil
	}
	msg := FormatLog(level, v...)
	if err := l.q.Put(msg); err != nil {
		return err
	}
	return nil
}

//Emergency system is unusable
func (l *FileLog) Emergency(v ...interface{}) error {
	return l.PushLog(EMEGENCY, v...)
}

//Alert action must be taken immediately
func (l *FileLog) Alert(v ...interface{}) error {
	return l.PushLog(ALERT, v...)
}

//Critical critical conditions
func (l *FileLog) Critical(v ...interface{}) error {
	return l.PushLog(CRITICAL, v...)
}

//Error error conditions
func (l *FileLog) Error(v ...interface{}) error {
	return l.PushLog(ERROR, v...)
}

//Warning warning conditions
func (l *FileLog) Warning(v ...interface{}) error {
	return l.PushLog(WARNING, v...)
}

//Notice normal but significant condition
func (l *FileLog) Notice(v ...interface{}) error {
	return l.PushLog(NOTICE, v...)
}

//Info informational messages
func (l *FileLog) Info(v ...interface{}) error {
	return l.PushLog(INFO, v...)
}

//Debug debug-level messages
func (l *FileLog) Debug(v ...interface{}) error {
	return l.PushLog(DEBUG, v...)
}
