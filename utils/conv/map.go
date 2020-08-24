package conv

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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

//JSON2Map json 转换为 map
func JSON2Map(s string) (map[string]interface{}, error) {

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil, err
	}
	return result, nil
}

//Map2Struct map转换为struct
func Map2Struct(data map[string]interface{}, out interface{}) error {
	if out == nil {
		return errors.New("Nil pointer")
	}

	valAtPtr := reflect.ValueOf(out)
	if valAtPtr.Type().Kind() == reflect.Ptr {
		ele := valAtPtr.Elem()
		valAtPtr = reflect.ValueOf(ele.Addr().Interface()).Elem()
	}

	if valAtPtr.Kind() != reflect.Struct {
		return errors.New("Output is not struct")
	}

	tagMp := make(map[string]reflect.Value)
	fieldCount := valAtPtr.NumField()
	for i := 0; i < fieldCount; i++ {
		fldValue := valAtPtr.Field(i)
		fldType := valAtPtr.Type().Field(i)
		tag := fldType.Tag.Get("ms")
		tagMp[tag] = fldValue
	}
	for key, val := range data {
		fld, ok := tagMp[key]
		if !ok {
			continue
		}
		switch fld.Kind() {
		case reflect.Int:
			fld.SetInt(int64(val.(int)))
		case reflect.String:
			fld.SetString(val.(string))
		case reflect.Bool:
			fld.SetBool(val.(bool))
		case reflect.Ptr:
			newEle := reflect.New(fld.Type().Elem())
			err := Map2Struct(val.(map[string]interface{}), newEle.Interface())
			if err != nil {
				return err
			}
			fld.Set(newEle)
		case reflect.Struct:
			newEle := reflect.New(fld.Type())
			err := Map2Struct(val.(map[string]interface{}), newEle.Interface())
			if err != nil {
				return err
			}
			fld.Set(newEle.Elem())
		default:
			return fmt.Errorf("Invalid type for field %v", fld.Type().Name())
		}
	}
	return nil
}
