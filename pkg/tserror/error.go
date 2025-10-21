package tserror

import (
	"github.com/pkg/errors"
)

type RespCode int
type HttpCode int

const skipStack = 2

// 成功码
const CodeSuccess RespCode = 0

// 通用错误码
const (
	ClientError          RespCode = 40000 // 参数验证失败
	CodeParamInvalid     RespCode = 40001 // 参数验证失败
	CodeParamBind        RespCode = 40002 // 参数绑定失败
	CodeUnauthorized     RespCode = 40003 // 未授权
	CodePermissionDenied RespCode = 40004 // 权限不足
)

// 业务模块错误码（用户模块示例）
const (
	CodeUserNotFound RespCode = 60001 // 用户不存在
	CodeUserExist    RespCode = 60002 // 用户已存在
)

// 系统错误码
const (
	CodeServerInternalError  RespCode = 50001 // 服务器内部错误
	CodeDatabaseError        RespCode = 50002 // 数据库错误
	CodeExternalServiceError RespCode = 50003 // 外部服务错误
)

// 错误码与消息映射
var errorMessages = map[RespCode]string{
	CodeSuccess:             "成功",
	CodeParamInvalid:        "参数验证失败",
	CodeUserNotFound:        "用户不存在",
	CodeServerInternalError: "系统内部错误",
}

// GetErrorMessage 获取错误消息
func GetErrorMessage(code RespCode) string {
	if msg, exists := errorMessages[code]; exists {
		return msg
	}
	return "未知错误"
}

const (
	HTTP_NOT_FOUND         HttpCode = 404
	HTTP_TOO_MANY_REQUESTS HttpCode = 429
)

// BusinessError 业务逻辑错误
type BusinessError struct {
	Code    RespCode `json:"code"`    // 业务错误码
	Message string   `json:"message"` // 用户友好错误信息
}

func (e *BusinessError) Error() string {
	return e.Message
}

// NewBusinessError 创建业务错误
func NewBusinessError(message string) error {
	return errors.Wrap(&BusinessError{
		Code:    ClientError,
		Message: message,
	}, message)
}

func NewBusinessErrorCode(code RespCode, message string) error {
	return errors.Wrap(&BusinessError{
		Code:    code,
		Message: message,
	}, message)
}

// SystemError 系统内部错误
type SystemError struct {
	Code     RespCode `json:"code"`
	Message  string   `json:"message"`
	Detail   string   `json:"detail"` // 内部错误详情（仅日志记录）
	Original error    `json:"-"`      // 原始错误，不序列化
}

func (e *SystemError) Error() string {
	if e.Original != nil {
		return e.Message + ": " + e.Original.Error()
	}
	return e.Message
}

// NewSystemError 创建系统错误
func NewSystemError(message string, original error) error {
	detail := ""
	if original != nil {
		detail = original.Error()
	}
	err := &SystemError{
		Code:     CodeServerInternalError,
		Message:  message,
		Detail:   detail,
		Original: original,
	}

	return errors.Wrap(err, message)
}

func NewSystemErrorCode(code RespCode, message string, original error) error {
	detail := ""
	if original != nil {
		detail = original.Error()
	}

	err := &SystemError{
		Code:     code,
		Message:  message,
		Detail:   detail,
		Original: original,
	}
	return errors.Wrap(err, message)

}
