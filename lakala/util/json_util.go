package util

import (
	"encoding/json"
	"reflect"
)

// JsonUtils JSON工具类
type JsonUtils struct{}

// ToJSONString 对象转JSON字符串
func (j *JsonUtils) ToJSONString(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(data)
}

// ToJSONStringIndent 对象转格式化的JSON字符串
func (j *JsonUtils) ToJSONStringIndent(obj interface{}) string {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}
	return string(data)
}

// Parse JSON字符串解析为对象
func (j *JsonUtils) Parse(text string, obj interface{}) error {
	return json.Unmarshal([]byte(text), obj)
}

// ParseToMap JSON字符串解析为map
func (j *JsonUtils) ParseToMap(text string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(text), &result)
	return result, err
}

// ParseToList JSON字符串解析为对象列表
func (j *JsonUtils) ParseToList(text string, obj interface{}) error {
	// 使用反射获取slice的类型
	sliceType := reflect.SliceOf(reflect.TypeOf(obj))
	sliceValue := reflect.New(sliceType).Interface()

	err := json.Unmarshal([]byte(text), &sliceValue)
	if err != nil {
		return err
	}

	// 将结果赋值给传入的指针
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(sliceValue).Elem())
	return nil
}

// 全局工具实例
var jsonUtils = &JsonUtils{}

// ToJSONString 全局函数：对象转JSON字符串
func ToJSONString(obj interface{}) string {
	return jsonUtils.ToJSONString(obj)
}

// ToJSONStringIndent 全局函数：对象转格式化的JSON字符串
func ToJSONStringIndent(obj interface{}) string {
	return jsonUtils.ToJSONStringIndent(obj)
}

// Parse 全局函数：JSON字符串解析为对象
func Parse(text string, obj interface{}) error {
	return jsonUtils.Parse(text, obj)
}

// ParseToMap 全局函数：JSON字符串解析为map
func ParseToMap(text string) (map[string]interface{}, error) {
	return jsonUtils.ParseToMap(text)
}

// ParseToList 全局函数：JSON字符串解析为对象列表
func ParseToList(text string, slicePtr interface{}) error {
	return json.Unmarshal([]byte(text), slicePtr)
}
