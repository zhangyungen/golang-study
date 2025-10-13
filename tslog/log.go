package tslog

import (
	"context"
	"os"
	"zyj.com/golang-study/define"
	"zyj.com/golang-study/util/ginutil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// 打印所有级别的日志
	lowPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging = zapcore.Lock(os.Stdout)

	consoleEncoder = zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core = zapcore.NewTee(
		// 打印在控制台
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
)

func Info(msg string, fields ...zapcore.Field) {
	Logger.Info(msg, fields...)
}
func Warn(msg string, fields ...zapcore.Field) {
	Logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zapcore.Field) {
	Logger.Error(msg, fields...)
}

// getTraceIDFromContext 从不同类型的 context 中提取 trace_id
func getTraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// 尝试从 context.Value 中获取 trace_id（如果有中间件存储了）
	if traceID, ok := ctx.Value(define.HEADER_TRACE_ID_KEY).(string); ok && traceID != "" {
		return traceID
	}

	// 尝试从 Gin context 中获取（如果 context 是从 Gin 传递过来的）
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginutil.GetTraceID(ginCtx)
	}

	return ""
}

// InfoCtx 带 context 的 Info 日志，自动添加 trace_id
func InfoCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	traceID := getTraceIDFromContext(ctx)
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	Logger.Info(msg, fields...)
}

// WarnCtx 带 context 的 Warn 日志，自动添加 trace_id
func WarnCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	traceID := getTraceIDFromContext(ctx)
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	Logger.Warn(msg, fields...)
}

// ErrorCtx 带 context 的 Error 日志，自动添加 trace_id
func ErrorCtx(ctx context.Context, msg string, fields ...zapcore.Field) {
	traceID := getTraceIDFromContext(ctx)
	if traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}
	Logger.Error(msg, fields...)
}
