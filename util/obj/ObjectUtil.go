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
		log.Println("ObjToJsonStr error ", err)
		return ""
	}
	return string(meta)
}

// JsonStrToObj 字符串转对象
func JsonStrToObj[T any](str string) *T {
	obj := new(T)
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		log.Println("JsonStrToObj error ", err)
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
		log.Println("ObjToObj error ", err)
	}
	return &t
}

//func ObjToObj(s interface{}, t interface{}) {
//	mapstructure.DecodeMetadata(s, t, &mapstructure.Metadata{Keys: ``})
//	err := copier.Copy(&t, s)
//	if err != nil {
//		log.Println("ObjToObj", err)
//	}
//	return &t
//}

func CopyToObj(s interface{}, t interface{}) {
	err := copier.Copy(t, s)
	if err != nil {
		log.Println("CopyToObj error ", err)
	}
}

//func CopierObj(s interface{}, t interface{}) interface{} {
//	if !isPointer(t) || !isPointer(s) {
//		panic("argument  must be pointer")
//	}
//	err := copier.Copy(t, s)
//	if err != nil {
//		log.Println("CopierObj", err)
//	}
//	return t
//}

func ObjToMap(s interface{}) map[interface{}]interface{} {
	var myMap map[interface{}]interface{}
	err := mapstructure.Decode(s, &myMap)
	if err != nil {
		log.Println("ObjToMap error ", err)
	}
	return myMap
}

func MapToObj[T any](param map[interface{}]interface{}) *T {
	var obj T
	err := mapstructure.Decode(param, &obj)
	if err != nil {
		log.Println("MapObject error ", err)
	}
	return &obj
}
func MapToObjByStr[T any](param map[string]interface{}) *T {
	var obj T
	err := mapstructure.Decode(param, &obj)
	if err != nil {
		log.Println("MapObject error ", err)
	}
	return &obj
}

func isPointer(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Ptr
}
