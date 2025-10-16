package lakala

import "fmt"

// NotificationResult 拉卡拉通知响应实体
// 响应失败,拉卡拉会重复通知;详情请参考拉卡拉开放平台官方文档
type NotificationResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Ok 创建成功响应
func Ok() *NotificationResult {
	return &NotificationResult{
		Code:    "SUCCESS",
		Message: "执行成功",
	}
}

// Fail 创建失败响应
func Fail() *NotificationResult {
	return &NotificationResult{
		Code:    "FAIL",
		Message: "失败",
	}
}

// Custom 创建自定义响应
func Custom(code, message string) *NotificationResult {
	return &NotificationResult{
		Code:    code,
		Message: message,
	}
}

// GetCode 获取响应码
func (n *NotificationResult) GetCode() string {
	return n.Code
}

// SetCode 设置响应码
func (n *NotificationResult) SetCode(code string) {
	n.Code = code
}

// GetMessage 获取响应消息
func (n *NotificationResult) GetMessage() string {
	return n.Message
}

// SetMessage 设置响应消息
func (n *NotificationResult) SetMessage(message string) {
	n.Message = message
}

// IsSuccess 判断是否成功
func (n *NotificationResult) IsSuccess() bool {
	return n.Code == "SUCCESS"
}

// ToJSON 转换为JSON字符串
func (n *NotificationResult) ToJSON() string {
	return fmt.Sprintf(`{"code":"%s","message":"%s"}`, n.Code, n.Message)
}
