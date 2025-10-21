package api

import (
	"fmt"
	"github.com/pkg/errors"
	"zyj.com/golang-study/pkg/tserror"
)

type Result struct {
	Code    tserror.RespCode `json:"code"`              // 状态码
	Message string           `json:"message,omitempty"` // 响应消息
	Data    interface{}      `json:"data,omitempty"`    // 响应数据
	Stack   string           `json:"stack,omitempty"`
}

func NewResult(data interface{}, err error) *Result {
	result := &Result{
		Code: tserror.CodeSuccess,
		Data: data,
	}
	if err != nil {
		causeErr := errors.Cause(err)
		var e *tserror.BizError
		var internalError *tserror.SystemError
		if errors.As(causeErr, &e) {
			result.Code = e.Code
			result.Message = e.Message
			result.Stack = fmt.Sprintf("%+v", err)
		} else if errors.As(causeErr, &internalError) {
			result.Code = internalError.Code
			result.Message = internalError.Detail
		} else {
			//logs.Error("result err", zap.Error(err))
			result.Code = tserror.CodeServerInternalError
			result.Message = "未知异常"
			result.Stack = fmt.Sprintf("%+v", err)
		}

	}
	return result
}
