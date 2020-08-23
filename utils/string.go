package utils

import (
	"fmt"
	"reflect"
	"strings"
)

//StrSub 字符串截取
func StrSub(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)
	if length == 0 {
		return str
	}

	if start < 0 {
		start = 0
	} else if start >= length {
		return ""
	}

	if start == end {
		return ""
	} else if end <= start || end > length {
		end = length
	}

	return string(rs[start:end])
}

//Implode arr to string
func Implode(list interface{}, seq string) string {
	listValue := reflect.Indirect(reflect.ValueOf(list))
	if listValue.Kind() != reflect.Slice {
		return ""
	}
	count := listValue.Len()
	listStr := make([]string, 0, count)
	for i := 0; i < count; i++ {
		v := listValue.Index(i)
		if str, err := getValue(v); err == nil {
			listStr = append(listStr, str)
		}
	}
	return strings.Join(listStr, seq)
}

func getValue(value reflect.Value) (res string, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		res, err = getValue(value.Elem())
	default:
		res = fmt.Sprint(value.Interface())
	}
	return
}
