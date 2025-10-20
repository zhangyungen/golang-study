package enums

import "strings"

// FunctionCode 对接Java枚举
type FunctionCode struct {
	Code string
	Name string
}

// URL 返回请求路径
func (f FunctionCode) URL() string {
	return strings.ReplaceAll(f.Code, ".", "/")
}
