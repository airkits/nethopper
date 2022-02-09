package utils

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
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

//Str2Bool string convert to bool
func Str2Bool(s string) bool {
	if len(s) <= 0 {
		return false
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		fmt.Printf("Str2Bool convert error %s", err.Error())
		return false
	}
	return v
}

//Str2Int string convert to Int
func Str2Int(s string) int {
	if len(s) <= 0 {
		return int(0)
	}
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		fmt.Printf("Str2Int convert error %s", err.Error())
		return 0
	}
	return int(v)
}

//Str2Int8 string convert to Int8
func Str2Int8(s string) int8 {
	if len(s) <= 0 {
		return int8(0)
	}
	v, e := strconv.ParseInt(s, 10, 8)
	if e != nil {
		fmt.Printf("Str2Int8 convert error %s", e.Error())
		return int8(0)
	}
	return int8(v)
}

//Str2Int16 string convert to Int16
func Str2Int16(s string) int16 {
	if len(s) <= 0 {
		return int16(0)
	}
	v, e := strconv.ParseInt(s, 10, 16)
	if e != nil {
		fmt.Printf("Str2Int16 convert error %s", e.Error())
		return int16(0)
	}
	return int16(v)
}

//Str2Int32 string convert to Int32
func Str2Int32(s string) int32 {
	if len(s) <= 0 {
		return int32(0)
	}
	v, e := strconv.ParseInt(s, 10, 32)
	if e != nil {
		fmt.Printf("Str2Int32 convert error %s", e.Error())
		return int32(0)
	}
	return int32(v)
}

//Str2Int64 string convert to Int64
func Str2Int64(s string) int64 {
	if len(s) <= 0 {
		return int64(0)
	}
	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		bigInt := &big.Int{}
		val, ok := bigInt.SetString(s, 10)
		if !ok {
			fmt.Printf("Str2Int64 convert error %s", e.Error())
			return int64(0)
		}
		return val.Int64()
	}
	return int64(v)
}

//Str2Uint string convert to Uint
func Str2Uint(s string) uint {
	if len(s) <= 0 {
		return uint(0)
	}
	v, e := strconv.ParseUint(s, 10, 64)
	if e != nil {
		fmt.Printf("Str2Uint convert error %s", e.Error())
		return uint(0)
	}
	return uint(v)
}

//Str2Uint8 string convert to Uint8
func Str2Uint8(s string) uint8 {
	if len(s) <= 0 {
		return uint8(0)
	}
	v, e := strconv.ParseUint(s, 10, 8)
	if e != nil {
		fmt.Printf("Str2Uint8 convert error %s", e.Error())
		return uint8(0)
	}
	return uint8(v)
}

//Str2Uint16 string convert to Uint16
func Str2Uint16(s string) uint16 {
	if len(s) <= 0 {
		return uint16(0)
	}
	v, e := strconv.ParseUint(s, 10, 16)
	if e != nil {
		fmt.Printf("Str2Uint16 convert error %s", e.Error())
		return uint16(0)
	}
	return uint16(v)
}

//Str2Uint32 string convert to Uint32
func Str2Uint32(s string) uint32 {
	if len(s) <= 0 {
		return uint32(0)
	}
	v, e := strconv.ParseUint(s, 10, 32)
	if e != nil {
		fmt.Printf("Str2Uint32 convert error %s", e.Error())
		return uint32(0)
	}
	return uint32(v)
}

//Str2Uint64 string convert to Uint64
func Str2Uint64(s string) uint64 {
	if len(s) <= 0 {
		return uint64(0)
	}
	v, e := strconv.ParseUint(s, 10, 64)
	if e != nil {
		bigInt := &big.Int{}
		val, ok := bigInt.SetString(s, 10)
		if !ok {
			fmt.Printf("Str2Uint64 convert error %s", e.Error())
			return uint64(0)
		}
		return val.Uint64()
	}
	return uint64(v)
}

//Str2Float32 string convert to Float32
func Str2Float32(s string) float32 {
	if len(s) <= 0 {
		return float32(0)
	}
	v, e := strconv.ParseFloat(s, 32)
	if e != nil {
		fmt.Printf("Str2Float32 convert error %s", e.Error())
		return float32(0)
	}
	return float32(v)
}

//Str2Float64 string convert to Float64
func Str2Float64(s string) float64 {
	if len(s) <= 0 {
		return float64(0)
	}
	v, e := strconv.ParseFloat(s, 3642)
	if e != nil {
		fmt.Printf("Str2Float64 convert error %s", e.Error())
		return float64(0)
	}
	return float64(v)
}
