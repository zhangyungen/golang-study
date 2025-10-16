package lakala

import "fmt"

// SDKExceptionEnums SDK异常枚举
type SDKExceptionEnums struct {
	Code    string
	Message string
}

func (e SDKExceptionEnums) GetCode() string {
	return e.Code
}

func (e SDKExceptionEnums) GetMessage() string {
	return e.Message
}

func (e *SDKExceptionEnums) SetMessage(message string) {
	e.Message = message
}

// 定义异常枚举
var (
	POST_ERROR                = SDKExceptionEnums{"SDK000001", "网络连接异常"}
	RES_IS_NULL               = SDKExceptionEnums{"SDK000002", "返回数据为空"}
	BAD_REQ                   = SDKExceptionEnums{"SDK000003", "请求异常"}
	SM4_INIT_FAIL             = SDKExceptionEnums{"SDK000004", "未初始化SM4"}
	SDK_APPID_NOT_INIT        = SDKExceptionEnums{"SDK000005", "SDK中APPID未初始化"}
	SDK_NOT_INIT              = SDKExceptionEnums{"SDK000006", "SDK未初始化"}
	INITIALIZE_KEYSTORE_ERROR = SDKExceptionEnums{"SDK000007", "初始化文件异常"}
	FILE_READ_FAIL_EXCEPTION  = SDKExceptionEnums{"SDK000008", "文件读取失败"}
	CHECK_FAIL                = SDKExceptionEnums{"SDK000009", "字段校验异常"}
	ERROR                     = SDKExceptionEnums{"SDK999999", "未知异常"}
)

// SDKException SDK异常
type SDKException struct {
	Code    string
	Message string
	Cause   error
}

func (e *SDKException) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("SDKException[%s]: %s, Cause: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("SDKException[%s]: %s", e.Code, e.Message)
}

func NewSDKException(message string, cause error) *SDKException {
	return &SDKException{Message: message, Cause: cause}
}

func NewSDKExceptionWithCode(code, message string) *SDKException {
	return &SDKException{Code: code, Message: message}
}

func NewSDKExceptionFromEnums(err SDKExceptionEnums) *SDKException {
	return &SDKException{Code: err.Code, Message: err.Message}
}

func NewSDKExceptionFromEnumsWithInfo(err SDKExceptionEnums, errInfo string) *SDKException {
	return &SDKException{
		Code:    err.Code,
		Message: fmt.Sprintf("%s[%s]", err.Message, errInfo),
	}
}
