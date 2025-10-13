package ginutil

import (
	"context"
	"github.com/gin-gonic/gin"
	"zyj.com/golang-study/define"
	"zyj.com/golang-study/util/strutil"
)

func GetRequestToken(c *gin.Context) string {
	return c.GetHeader("Authorization")
}

func GetTraceID(c *gin.Context) string {
	if c != nil {
		traceId := c.GetHeader(define.HEADER_TRACE_ID_KEY)
		if traceId == "" {
			traceId = c.GetString(define.HEADER_TRACE_ID_KEY)
			if traceId == "" {
				traceId = strutil.GenerateUUIDV4()
				c.Set(define.HEADER_TRACE_ID_KEY, traceId)
			}
		}
		return traceId
	}
	return ""
}
func AppendContextDebugMsg(c *gin.Context, msg string) {
	if c != nil {
		oriMsg := c.GetString(define.LOG_DEBUG_MSG)
		if oriMsg != "" {
			oriMsg += ", "
		}
		c.Set(define.LOG_DEBUG_MSG, oriMsg+msg)
	}
}

func GetCheeseID(c *gin.Context) string {
	return c.GetHeader("Cheese-ID")
}
func GetUserID(c *gin.Context) string {
	return c.GetHeader("User-ID")
}

func GetHeaderTs(c *gin.Context) string {
	return c.GetHeader("ts")
}

// GetTraceIDFromContext 从标准 context 中获取 trace_id
func GetTraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	if traceID, ok := ctx.Value(define.HEADER_TRACE_ID_KEY).(string); ok {
		return traceID
	}

	return ""
}

// ContextWithTraceID 从 gin.Context 获取 trace_id 并创建包含 trace_id 的 context.Context
// 优先使用 gin.Context 的请求 context，如果已包含 trace_id 则直接返回
// 否则创建新的 context 并添加 trace_id
func ContextWithTraceID(c *gin.Context) context.Context {
	if c == nil {
		return context.Background()
	}

	// 获取请求的 context
	reqCtx := c.Request.Context()

	// 检查是否已经有 trace_id（可能由中间件设置）
	if existingTraceID := GetTraceIDFromContext(reqCtx); existingTraceID != "" {
		return reqCtx
	}

	// 从 gin.Context 获取 trace_id
	traceID := GetTraceID(c)

	// 如果没有 trace_id，返回原始请求 context
	if traceID == "" {
		return reqCtx
	}

	// 创建包含 trace_id 的新 context
	return context.WithValue(reqCtx, define.HEADER_TRACE_ID_KEY, traceID)
}
