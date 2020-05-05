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
// * @Date: 2019-06-11 07:48:41
// * @Last Modified by:   ankye
// * @Last Modified time: 2019-06-11 07:48:41

package utils

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// FileIsExist if file exist ,return true
func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// GetWorkDirectory get current exec file directory
func GetWorkDirectory() (string, error) {

	// in test case or go run xxx, work directory is temp directory
	// workDir := os.Getenv("WORK_DIR")
	// if len(workDir) > 0 {
	// 	return filepath.Abs(workDir)
	// }

	// workDir, err := os.Getwd()
	// if err != nil {
	// 	return "", err
	// }

	// file, err := exec.LookPath(os.Args[0])
	// if err != nil {
	// 	return "", err
	// }
	// path, err := filepath.Abs(file)
	// if err != nil {
	// 	return "", err
	// }
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
		i := strings.LastIndex(path, "/")
		if i < 0 {
			return "", errors.New("work directory invalid")
		}
		return string(path[0 : i+1]), nil
	}
	return path + "/", nil

}

// GetAbsDirectory get file directory abs path
func GetAbsDirectory(filename string) string {
	filename = filepath.FromSlash(path.Clean(filename))
	if runtime.GOOS == "windows" {
		filename = strings.Replace(filename, "\\", "/", -1)
	}
	i := strings.LastIndex(filename, "/")
	if i < 0 {
		return ""
	}
	filename = string(filename[0 : i+1])
	if filepath.IsAbs(filename) {
		return filename
	}
	workDir, err := GetWorkDirectory()
	if err != nil {
		return ""
	}
	return filepath.Join(workDir, filename)
}

// GetAbsFilePath get file abs path
func GetAbsFilePath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	workDir, err := GetWorkDirectory()
	if err != nil {
		return ""
	}
	return filepath.Join(workDir, filename)
}

// FileLines get file lines
func FileLines(filename string) (int32, error) {
	fd, err := os.Open(filename)
	//fd, err := mmap.Open(filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	maxbuf := 32768
	buf := make([]byte, maxbuf) // 32k
	var count int32
	lineSep := []byte{'\n'}
	offset := int64(0)
	for {
		c, err := fd.Read(buf)
		//c, err := fd.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return count, nil
		}
		offset += int64(c)
		count += int32(bytes.Count(buf[:c], lineSep))
		if err == io.EOF {
			break
		}
	}
	return count, nil
}
