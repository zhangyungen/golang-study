package tslog

import (
	"context"
	"errors"
	"gorm.io/gorm/utils"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"zyj.com/golang-study/define"
	"zyj.com/golang-study/util/ginutil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type ContextFn func(ctx context.Context) []zapcore.Field

// defaultContextFn 默认的 context 字段提取器，提取 trace_id
func defaultContextFn(ctx context.Context) []zapcore.Field {
	if ctx == nil {
		return nil
	}

	var fields []zapcore.Field

	// 尝试从 context.Value 中获取 trace_id
	if traceID, ok := ctx.Value(define.HEADER_TRACE_ID_KEY).(string); ok && traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	} else if ginCtx, ok := ctx.(*gin.Context); ok {
		// 尝试从 Gin context 中获取
		if traceID := ginutil.GetTraceID(ginCtx); traceID != "" {
			fields = append(fields, zap.String("trace_id", traceID))
		}
	}

	return fields
}

type XLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormlogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
	Context                   ContextFn
}

func New(zapLogger *zap.Logger) XLogger {
	return XLogger{
		ZapLogger:                 zapLogger,
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
		Context:                   defaultContextFn,
	}
}

func (l XLogger) SetAsDefault() {
	gormlogger.Default = l
}

func (l XLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return XLogger{
		ZapLogger:                 l.ZapLogger,
		SlowThreshold:             l.SlowThreshold,
		LogLevel:                  level,
		SkipCallerLookup:          l.SkipCallerLookup,
		IgnoreRecordNotFoundError: l.IgnoreRecordNotFoundError,
		Context:                   l.Context,
	}
}

func (l XLogger) Info(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Info {
		return
	}
	l.logger(ctx).Sugar().Debugf(str, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l XLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Warn {
		return
	}
	l.logger(ctx).Sugar().Warnf(str, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l XLogger) Error(ctx context.Context, str string, args ...interface{}) {
	if l.LogLevel < gormlogger.Error {
		return
	}
	l.logger(ctx).Sugar().Errorf(str, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l XLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	logger := l.logger(ctx)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		logger.Error("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		logger.Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case l.LogLevel >= gormlogger.Info:
		sql, rows := fc()
		logger.Debug("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}

var (
	gormPackage    = filepath.Join("gorm.io", "gorm")
	zapgormPackage = filepath.Join("moul.io", "zapgorm2")
)

func (l XLogger) logger(ctx context.Context) *zap.Logger {
	logger := l.ZapLogger
	if l.Context != nil {
		fields := l.Context(ctx)
		logger = logger.With(fields...)
	}

	if l.SkipCallerLookup {
		return logger
	}
	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, "gorm"):
		//case strings.Contains(file, zapgormPackage):
		default:
			return logger.WithOptions(zap.AddCallerSkip(i - 1))
		}
	}
	return logger
}
