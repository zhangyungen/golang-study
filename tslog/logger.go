// logger/logger.go
package tslog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
)

// 全局变量
var (
	globalLogger *zap.SugaredLogger
	once         sync.Once
)

// Init 初始化日志（简单版本）
func Init(logPath string, debug bool) error {
	var err error
	once.Do(func() {
		globalLogger, err = createLogger(logPath, debug)
	})
	return err
}

// InitWithConfig 初始化日志（带配置版本）
func InitWithConfig(config *Config) error {
	var err error
	once.Do(func() {
		globalLogger, err = createLoggerWithConfig(config)
	})
	return err
}

// Config 日志配置
type Config struct {
	LogPath    string // 日志文件路径，空则只输出到控制台
	Level      string // 日志级别: debug, info, warn, error
	MaxSize    int    // 单个日志文件最大大小(MB)
	MaxBackups int    // 保留的旧日志文件最大个数
	MaxAge     int    // 保留旧日志文件的最大天数
	Compress   bool   // 是否压缩旧日志文件
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		MaxSize:    100,  // 100MB
		MaxBackups: 10,   // 保留10个文件
		MaxAge:     30,   // 保留30天
		Compress:   true, // 压缩旧文件
	}
}

// createLogger 创建日志器（简单版本）
func createLogger(logPath string, debug bool) (*zap.SugaredLogger, error) {
	config := DefaultConfig()
	config.LogPath = logPath
	if debug {
		config.Level = "debug"
	}
	return createLoggerWithConfig(config)
}

// createLoggerWithConfig 创建日志器（配置版本）
func createLoggerWithConfig(config *Config) (*zap.SugaredLogger, error) {
	// 设置日志级别
	level := getZapLevel(config.Level)

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 大写级别编码
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径调用者
	}

	// 创建编码器
	var encoder zapcore.Encoder
	if config.Level == "debug" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // 开发环境用控制台格式
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // 生产环境用JSON格式
	}

	// 创建多个写入器
	cores := []zapcore.Core{}

	// 控制台输出（始终启用）
	consoleCore := zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level)
	cores = append(cores, consoleCore)

	// 文件输出（如果配置了日志路径）
	if config.LogPath != "" {
		fileCore, err := createFileCore(encoder, level, config)
		if err != nil {
			return nil, err
		}
		cores = append(cores, fileCore)
	}

	// 创建核心
	core := zapcore.NewTee(cores...)

	// 创建日志器
	logger := zap.New(core,
		zap.AddCaller(),      // 添加调用者信息
		zap.AddCallerSkip(1), // 跳过一层调用栈
	)

	return logger.Sugar(), nil
}

// createFileCore 创建文件日志核心
func createFileCore(encoder zapcore.Encoder, level zapcore.LevelEnabler, config *Config) (zapcore.Core, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(filepath.Dir(config.LogPath), 0755); err != nil {
		return nil, err
	}

	// 创建 lumberjack 日志滚动器
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.LogPath,
		MaxSize:    config.MaxSize,    // 单个文件最大大小(MB)
		MaxBackups: config.MaxBackups, // 保留的旧文件最大个数
		MaxAge:     config.MaxAge,     // 保留旧文件的最大天数
		Compress:   config.Compress,   // 是否压缩旧文件
		LocalTime:  true,              // 使用本地时间
	}

	// 创建文件写入器
	fileWriteSyncer := zapcore.AddSync(lumberjackLogger)

	return zapcore.NewCore(encoder, fileWriteSyncer, level), nil
}

// getZapLevel 将字符串级别转换为 zapcore.Level
func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// 全局日志函数
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

//func Info(args ...interface{}) {
//	globalLogger.Info(args...)
//}
//
//func Warn(args ...interface{}) {
//	globalLogger.Warn(args...)
//}
//
//func Error(args ...interface{}) {
//	globalLogger.Error(args...)
//}

func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	globalLogger.Panic(args...)
}

// 格式化日志函数
func Debugf(format string, args ...interface{}) {
	globalLogger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	globalLogger.Panicf(format, args...)
}

// 带字段的日志函数
func Debugw(msg string, keysAndValues ...interface{}) {
	globalLogger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	globalLogger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	globalLogger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	globalLogger.Errorw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	globalLogger.Fatalw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	globalLogger.Panicw(msg, keysAndValues...)
}

// 创建带有字段的日志器
func With(args ...interface{}) *zap.SugaredLogger {
	return globalLogger.With(args...)
}

// Sync 刷新日志缓冲区
func Sync() error {
	return globalLogger.Sync()
}

// GetLogger 获取原始日志器（高级用法）
func GetLogger() *zap.SugaredLogger {
	return globalLogger
}

// 便捷初始化函数
func InitDevelopment() error {
	return Init("", true)
}

func InitProduction(logPath string) error {
	return Init(logPath, false)
}
