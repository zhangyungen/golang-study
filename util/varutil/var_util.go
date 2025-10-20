package varutil

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
	"log"
	"reflect"
)

// JsonStr 对象转JSON字符串
func JsonStr(v interface{}) string {
	meta, err := json.Marshal(v)
	if err != nil {
		log.Println("JsonStr error ", err)
		return ""
	}
	return string(meta)
}

func PrettyJsonStr(v interface{}) string {

	marshal, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("json序列化失败: %v", err)
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, marshal, "", "\t")
	if err != nil {
		log.Println("JSON parse error: ", err)
		return ""
	}
	return string(prettyJSON.Bytes())
}

// JsonStrToStruct 字符串转对象
func JsonStrToStruct[T any](str string) *T {
	obj := new(T)
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		log.Println("JsonStrToStruct error ", err)
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
func ConvertTo[T any](s interface{}) *T {
	var t T
	err := copier.Copy(&t, s)
	if err != nil {
		log.Println("ConvertTo error ", err)
	}
	return &t
}

//func ConvertTo(s interface{}, t interface{}) {
//	mapstructure.DecodeMetadata(s, t, &mapstructure.Metadata{Keys: ``})
//	err := copier.Copy(&t, s)
//	if err != nil {
//		log.Println("ConvertTo", err)
//	}
//	return &t
//}

func CopyTo(s interface{}, t interface{}) {
	err := copier.Copy(t, s)
	if err != nil {
		log.Println("CopyTo error ", err)
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

func StructToMap(s interface{}) map[interface{}]interface{} {
	var myMap map[interface{}]interface{}
	err := mapstructure.Decode(s, &myMap)
	if err != nil {
		log.Println("StructToMap error ", err)
	}
	return myMap
}

func MapToStruct[T any](param map[interface{}]interface{}) *T {
	var obj T
	err := mapstructure.Decode(param, &obj)
	if err != nil {
		log.Println("MapObject error ", err)
	}
	return &obj
}
func MapToStructByStr[T any](param map[string]interface{}) *T {
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
