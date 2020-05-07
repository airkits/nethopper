package utils

// reference: https://github.com/mohae/deepcopy
import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
)

func deepCopy(dst, src reflect.Value) {
	switch src.Kind() {
	case reflect.Interface:
		value := src.Elem()
		if !value.IsValid() {
			return
		}
		newValue := reflect.New(value.Type()).Elem()
		deepCopy(newValue, value)
		dst.Set(newValue)
	case reflect.Ptr:
		value := src.Elem()
		if !value.IsValid() {
			return
		}
		dst.Set(reflect.New(value.Type()))
		deepCopy(dst.Elem(), value)
	case reflect.Map:
		dst.Set(reflect.MakeMap(src.Type()))
		keys := src.MapKeys()
		for _, key := range keys {
			value := src.MapIndex(key)
			newValue := reflect.New(value.Type()).Elem()
			deepCopy(newValue, value)
			dst.SetMapIndex(key, newValue)
		}
	case reflect.Slice:
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			deepCopy(dst.Index(i), src.Index(i))
		}
	case reflect.Struct:
		typeSrc := src.Type()
		for i := 0; i < src.NumField(); i++ {
			value := src.Field(i)
			tag := typeSrc.Field(i).Tag
			if value.CanSet() && tag.Get("deepcopy") != "-" {
				deepCopy(dst.Field(i), value)
			}
		}
	default:
		dst.Set(src)
	}
}

func DeepCoderCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func DeepCopy(dst, src interface{}) {
	typeDst := reflect.TypeOf(dst)
	typeSrc := reflect.TypeOf(src)
	if typeDst != typeSrc {
		panic("DeepCopy: " + typeDst.String() + " != " + typeSrc.String())
	}
	if typeSrc.Kind() != reflect.Ptr {
		panic("DeepCopy: pass arguments by address")
	}

	valueDst := reflect.ValueOf(dst).Elem()
	valueSrc := reflect.ValueOf(src).Elem()
	if !valueDst.IsValid() || !valueSrc.IsValid() {
		panic("DeepCopy: invalid arguments")
	}

	deepCopy(valueDst, valueSrc)
}

func DeepClone(v interface{}) interface{} {
	dst := reflect.New(reflect.TypeOf(v)).Elem()
	deepCopy(dst, reflect.ValueOf(v))
	return dst.Interface()
}

func CopyFloat32Map(source map[int32]float32) map[int32]float32 {
	target := make(map[int32]float32)
	for key, value := range source {
		target[key] = value
	}
	return target
}

func CopyFloat64Map(source map[int32]float64) map[int32]float64 {
	target := make(map[int32]float64)
	for key, value := range source {
		target[key] = value
	}
	return target
}

func CopyInt32Map(source map[int32]int32) map[int32]int32 {
	target := make(map[int32]int32)
	for key, value := range source {
		target[key] = value
	}
	return target
}

func CopyInt64Map(source map[int32]int64) map[int32]int64 {
	target := make(map[int32]int64)
	for key, value := range source {
		target[key] = value
	}
	return target
}

func MargeMap(a map[int32]int32, b map[int32]int32) {
	for key, value := range b {
		a[key] = value
	}
}

func CopyMapInt32(a map[int32]int32) map[int32]int32 {
	results := make(map[int32]int32)
	for key, value := range a {
		results[key] = value
	}
	return results
}

func CopyJSON(marshaler interface{}, unMarshaler interface{}) error {
	data, error := json.Marshal(marshaler)
	if error != nil {
		return error
	}
	return json.Unmarshal(data, unMarshaler)
}
