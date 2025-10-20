package utils

import (
	"bytes"
	"encoding/json"
)

// ToJSONString 编码为JSON字符串
func ToJSONString(v interface{}) (string, error) {
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return "", err
	}
	// 去掉Encode追加的换行符
	result := buf.String()
	if len(result) > 0 && result[len(result)-1] == '\n' {
		result = result[:len(result)-1]
	}
	return result, nil
}

// ParseJSON 解析JSON字符串
func ParseJSON[T any](text string) (T, error) {
	var target T
	decoder := json.NewDecoder(bytes.NewBufferString(text))
	if err := decoder.Decode(&target); err != nil {
		var zero T
		return zero, err
	}
	return target, nil
}

// ParseJSONList 解析JSON数组
func ParseJSONList[T any](text string) ([]T, error) {
	var target []T
	decoder := json.NewDecoder(bytes.NewBufferString(text))
	if err := decoder.Decode(&target); err != nil {
		return nil, err
	}
	return target, nil
}
