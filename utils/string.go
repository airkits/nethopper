package utils

import (
	"math/big"
	"reflect"
	"strconv"
)

//Struct2Map struct 2 map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//Str2Bool string convert to bool
func Str2Bool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

//Str2Int string convert to Int
func Str2Int(s string) (int, error) {
	v, e := strconv.ParseInt(s, 10, 32)
	return int(v), e
}

//Str2Int8 string convert to Int8
func Str2Int8(s string) (int8, error) {
	v, e := strconv.ParseInt(s, 10, 8)
	return int8(v), e
}

//Str2Int16 string convert to Int16
func Str2Int16(s string) (int16, error) {
	v, e := strconv.ParseInt(s, 10, 16)
	return int16(v), e
}

//Str2Int32 string convert to Int32
func Str2Int32(s string) (int32, error) {
	v, e := strconv.ParseInt(s, 10, 32)
	return int32(v), e
}

//Str2Int64 string convert to Int64
func Str2Int64(s string) (int64, error) {
	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		bigInt := &big.Int{}
		val, ok := bigInt.SetString(s, 10)
		if !ok {
			return v, e
		}
		return val.Int64(), nil
	}
	return int64(v), e
}

//Str2Uint string convert to Uint
func Str2Uint(s string) (uint, error) {
	v, e := strconv.ParseUint(s, 10, 64)
	return uint(v), e
}

//Str2Uint8 string convert to Uint8
func Str2Uint8(s string) (uint8, error) {
	v, e := strconv.ParseUint(s, 10, 8)
	return uint8(v), e
}

//Str2Uint16 string convert to Uint16
func Str2Uint16(s string) (uint16, error) {
	v, e := strconv.ParseUint(s, 10, 16)
	return uint16(v), e
}

//Str2Uint32 string convert to Uint32
func Str2Uint32(s string) (uint32, error) {
	v, e := strconv.ParseUint(s, 10, 32)
	return uint32(v), e
}

//Str2Uint64 string convert to Uint64
func Str2Uint64(s string) (uint64, error) {
	v, e := strconv.ParseUint(s, 10, 64)
	if e != nil {
		bigInt := &big.Int{}
		val, ok := bigInt.SetString(s, 10)
		if !ok {
			return v, e
		}
		return val.Uint64(), nil
	}
	return uint64(v), e
}

//Str2Float32 string convert to Float32
func Str2Float32(s string) (float32, error) {
	v, e := strconv.ParseFloat(s, 32)
	return float32(v), e
}

//Str2Float64 string convert to Float64
func Str2Float64(s string) (float64, error) {
	v, e := strconv.ParseFloat(s, 3642)
	return float64(v), e
}
