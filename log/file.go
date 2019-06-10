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
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gonethopper/nethopper/utils"
	"golang.org/x/exp/mmap"
)

//NewFileLogger create FileLog instance
func NewFileLogger(m map[string]interface{}) (Log, error) {
	logger := &FileLog{}
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
	sync.RWMutex
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

	fmt.Println(m)
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
	filename := l.prefix
	if len(timestr) > 0 {
		filename += "_" + timestr
	}
	if num > 0 {
		filename += "_" + strconv.Itoa(num)
	}
	return filename + l.suffix
}

//nextNumTest test the file actually exists in the filesystem,return the next file num
//if test failed return -1
func (l *FileLog) nextNumTest(timestr string) int {
	MaxNum := 1000
	for i := 1; i < MaxNum; i++ {
		filename := l.genFilename(timestr, i)
		_, err := os.Stat(filename)
		if err != nil { //not exist
			return i
		}
	}
	return -1
}

//filecutTest check time/maxsize/maxlines
func (l *FileLog) filecutTest() bool {
	timestr := l.genCurrentTime()
	if strings.Compare(l.currentTime, timestr) != 0 {
		fmt.Println("compare")
		return true
	}
	if l.currentSize >= l.maxSize || l.currentLines >= l.maxLines {
		return true
	}
	return false
}

//InitLogger init logger
func (l *FileLog) InitLogger() error {
	return l.createNewFile()
}

//Close close file logger
func (l *FileLog) Close() error {

	defer l.currentWriter.Close()
	if err := l.flush(); err != nil {
		return err
	}
	return nil
}

func (l *FileLog) flush() error {
	return l.currentWriter.Sync()
}

//lines
func (l *FileLog) lines() (int, error) {
	//fd, err := os.Open(l.fileName)
	fd, err := mmap.Open(l.fileName)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	maxbuf := 32768
	buf := make([]byte, maxbuf) // 32k
	count := 0
	lineSep := []byte{'\n'}
	offset := int64(0)
	for {
		//c, err := fd.Read(buf)
		c, err := fd.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return count, nil
		}
		offset += int64(c)
		count += bytes.Count(buf[:c], lineSep)
		if err == io.EOF {
			break
		}
	}
	return count, nil
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
func (l *FileLog) isExist(path string) bool {
	_, err := os.Stat(l.fileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

//createNewFile if file exist,then check current lines and filesize
func (l *FileLog) createNewFile() error {
	l.currentTime = l.genCurrentTime()
	l.fileName = l.genFilename(l.currentTime, 0)
	flag := l.isExist(l.fileName)
	if flag { //exist file
		lines, err := l.lines()
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

//WriteLog write message to file, return immediately if not meet the conditions
func (l *FileLog) WriteLog(level int, format string, v ...interface{}) error {

	if level > l.level {
		return nil
	}

	msg := FormatLog(level, format, v...)
	l.Lock()
	defer l.Unlock()

	if l.filecutTest() {
		l.moveFile()
		l.createNewFile()
	}
	_, err := l.currentWriter.Write([]byte(msg))
	if err == nil {
		l.currentLines++
		l.currentSize += len(msg)
	}

	return nil
}

//Emergency system is unusable
func (l *FileLog) Emergency(format string, v ...interface{}) error {
	return l.WriteLog(EMEGENCY, format, v...)
}

//Alert action must be taken immediately
func (l *FileLog) Alert(format string, v ...interface{}) error {
	return l.WriteLog(ALERT, format, v...)
}

//Critical critical conditions
func (l *FileLog) Critical(format string, v ...interface{}) error {
	return l.WriteLog(CRITICAL, format, v...)
}

//Error error conditions
func (l *FileLog) Error(format string, v ...interface{}) error {
	return l.WriteLog(ERROR, format, v...)
}

//Warning warning conditions
func (l *FileLog) Warning(format string, v ...interface{}) error {
	return l.WriteLog(WARNING, format, v...)
}

//Notice normal but significant condition
func (l *FileLog) Notice(format string, v ...interface{}) error {
	return l.WriteLog(NOTICE, format, v...)
}

//Info informational messages
func (l *FileLog) Info(format string, v ...interface{}) error {
	return l.WriteLog(INFO, format, v...)
}

//Debug debug-level messages
func (l *FileLog) Debug(format string, v ...interface{}) error {
	return l.WriteLog(DEBUG, format, v...)
}
