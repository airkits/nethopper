package conv

import (
	"math/big"
	"reflect"
	"strconv"

	"github.com/gonethopper/nethopper/server"
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
func Str2Bool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		server.Info("Str2Bool convert error %s", err.Error())
		return false
	}
	return v
}

//Str2Int string convert to Int
func Str2Int(s string) int {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		server.Info("Str2Int convert error %s", err.Error())
		return 0
	}
	return int(v)
}

//Str2Int8 string convert to Int8
func Str2Int8(s string) int8 {
	v, e := strconv.ParseInt(s, 10, 8)
	if e != nil {
		server.Info("Str2Int8 convert error %s", e.Error())
		return int8(0)
	}
	return int8(v)
}

//Str2Int16 string convert to Int16
func Str2Int16(s string) int16 {
	v, e := strconv.ParseInt(s, 10, 16)
	if e != nil {
		server.Info("Str2Int16 convert error %s", e.Error())
		return int16(0)
	}
	return int16(v)
}

//Str2Int32 string convert to Int32
func Str2Int32(s string) int32 {
	v, e := strconv.ParseInt(s, 10, 32)
	if e != nil {
		server.Info("Str2Int32 convert error %s", e.Error())
		return int32(0)
	}
	return int32(v)
}

//Str2Int64 string convert to Int64
func Str2Int64(s string) int64 {
	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		bigInt := &big.Int{}
		val, ok := bigInt.SetString(s, 10)
		if !ok {
			server.Info("Str2Int64 convert error %s", e.Error())
			return int64(0)
		}
		return val.Int64()
	}
	return int64(v)
}

//Str2Uint string convert to Uint
func Str2Uint(s string) uint {
	v, e := strconv.ParseUint(s, 10, 64)
	if e != nil {
		server.Info("Str2Uint convert error %s", e.Error())
		return uint(0)
	}
	return uint(v)
}

//Str2Uint8 string convert to Uint8
func Str2Uint8(s string) uint8 {
	v, e := strconv.ParseUint(s, 10, 8)
	if e != nil {
		server.Info("Str2Uint8 convert error %s", e.Error())
		return uint8(0)
	}
	return uint8(v)
}

//Str2Uint16 string convert to Uint16
func Str2Uint16(s string) uint16 {
	v, e := strconv.ParseUint(s, 10, 16)
	if e != nil {
		server.Info("Str2Uint16 convert error %s", e.Error())
		return uint16(0)
	}
	return uint16(v)
}

//Str2Uint32 string convert to Uint32
func Str2Uint32(s string) uint32 {
	v, e := strconv.ParseUint(s, 10, 32)
	if e != nil {
		server.Info("Str2Uint32 convert error %s", e.Error())
		return uint32(0)
	}
	return uint32(v)
}

//Str2Uint64 string convert to Uint64
func Str2Uint64(s string) uint64 {
	v, e := strconv.ParseUint(s, 10, 64)
	if e != nil {
		bigInt := &big.Int{}
		val, ok := bigInt.SetString(s, 10)
		if !ok {
			server.Info("Str2Uint64 convert error %s", e.Error())
			return uint64(0)
		}
		return val.Uint64()
	}
	return uint64(v)
}

//Str2Float32 string convert to Float32
func Str2Float32(s string) float32 {
	v, e := strconv.ParseFloat(s, 32)
	if e != nil {
		server.Info("Str2Float32 convert error %s", e.Error())
		return float32(0)
	}
	return float32(v)
}

//Str2Float64 string convert to Float64
func Str2Float64(s string) float64 {
	v, e := strconv.ParseFloat(s, 3642)
	if e != nil {
		server.Info("Str2Float64 convert error %s", e.Error())
		return float64(0)
	}
	return float64(v)
}
