package obj

import (
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"log"
	"reflect"
)

// ObjToJsonStr 对象转JSON字符串
func ObjToJsonStr(obj interface{}) string {
	meta, err := json.Marshal(obj)
	if err != nil {
		log.Printf("ObjToJsonStr", err)
		return ""
	}
	return string(meta)
}

// JsonStrToObj 字符串转对象
func JsonStrToObj[T any](str string) *T {
	obj := new(T)
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		log.Printf("JsonStrToObj", err)
		return obj
	}
	return obj
}

//var metadata mapstructure.Metadata

/**
 * 对象转对象
 * @param s 源对象
 * @param t 目标对象
 */
func ObjToObj[T any](s interface{}) *T {
	var t T
	err := copier.Copy(&t, s)
	if err != nil {
		log.Printf("ObjToObj", err)
	}
	return &t
}

//func CopierObj(s interface{}, t interface{}) interface{} {
//	if !isPointer(t) || !isPointer(s) {
//		panic("argument  must be pointer")
//	}
//	err := copier.Copy(t, s)
//	if err != nil {
//		log.Printf("CopierObj", err)
//	}
//	return t
//}

func ObjToMap(s interface{}) map[interface{}]interface{} {
	var myMap map[interface{}]interface{}
	err := mapstructure.Decode(s, &myMap)
	if err != nil {
		log.Printf("ObjToMap", err)
	}
	return myMap
}

func MapToObj[T any](param map[interface{}]interface{}) *T {
	var obj T
	err := mapstructure.Decode(param, &obj)
	if err != nil {
		log.Printf("MapObject", err)
	}
	return &obj
}

func isPointer(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Ptr
}
