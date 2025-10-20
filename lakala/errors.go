package laopsdk

import "fmt"

// ErrorDescriptor 用于描述SDK错误代码及默认文案
// 中文注释遵循原Java枚举的含义
type ErrorDescriptor struct {
	Code    string
	Message string
}

// SDKError 对应Java中的 SDKException
type SDKError struct {
	Code    string
	Message string
	Err     error
}

func (e *SDKError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *SDKError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func newSDKError(desc ErrorDescriptor, extra string, err error) *SDKError {
	msg := desc.Message
	if extra != "" {
		msg = fmt.Sprintf("%s[%s]", msg, extra)
	}
	return &SDKError{Code: desc.Code, Message: msg, Err: err}
}

// 预置错误码定义，保持与Java版本一致
var (
	ErrPostError    = ErrorDescriptor{Code: "SDK000001", Message: "网络连接异常"}
	ErrResponseNil  = ErrorDescriptor{Code: "SDK000002", Message: "返回数据为空"}
	ErrBadRequest   = ErrorDescriptor{Code: "SDK000003", Message: "请求异常"}
	ErrSM4InitFail  = ErrorDescriptor{Code: "SDK000004", Message: "未初始化SM4"}
	ErrAppIDNotInit = ErrorDescriptor{Code: "SDK000005", Message: "SDK中APPID未初始化"}
	ErrSDKNotInit   = ErrorDescriptor{Code: "SDK000006", Message: "SDK未初始化"}
	ErrKeystoreInit = ErrorDescriptor{Code: "SDK000007", Message: "初始化文件异常"}
	ErrFileRead     = ErrorDescriptor{Code: "SDK000008", Message: "文件读取失败"}
	ErrCheckFail    = ErrorDescriptor{Code: "SDK000009", Message: "字段校验异常"}
	ErrUnknown      = ErrorDescriptor{Code: "SDK999999", Message: "未知异常"}
)

func newError(desc ErrorDescriptor, err error) error {
	if desc.Code == "" {
		desc = ErrUnknown
	}
	return newSDKError(desc, "", err)
}

func newErrorWithExtra(desc ErrorDescriptor, extra string, err error) error {
	if desc.Code == "" {
		desc = ErrUnknown
	}
	return newSDKError(desc, extra, err)
}
