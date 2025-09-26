// example/main.go
package main

import (
	"errors"
	"fmt"
	"time"
	"zyj.com/golang-study/util/logger"
)

func main() {
	// 示例1: 最简单的初始化（开发环境）
	//exampleSimple()

	// 示例2: 生产环境配置
	exampleProduction()

	// 示例3: 自定义配置
	//exampleCustom()

	// 示例4: 各种日志用法
	exampleUsage()
}

func exampleSimple() {
	fmt.Println("=== 示例1: 开发环境 ===")

	// 初始化开发环境日志（只输出到控制台）
	err := logger.InitDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("开发环境日志初始化成功")
	logger.Debug("这是一条调试信息")
	logger.Warn("这是一条警告信息")

	// 模拟业务逻辑
	processOrder("order_123")
}

func exampleProduction() {
	fmt.Println("\n=== 示例2: 生产环境 ===")

	// 初始化生产环境日志（输出到文件和控制台）
	err := logger.InitProduction("logs/app.log")
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("生产环境日志初始化成功")
	logger.Debug("这条调试信息在生产环境不会显示") // 因为默认级别是info

	// 业务逻辑
	processUser("user_456")
}

func exampleCustom() {
	fmt.Println("\n=== 示例3: 自定义配置 ===")

	config := &logger.Config{
		LogPath:    "logs/custom.log",
		Level:      "debug", // 设置为debug级别
		MaxSize:    50,      // 50MB
		MaxBackups: 5,       // 保留5个文件
		MaxAge:     7,       // 保留7天
		Compress:   true,
	}

	err := logger.InitWithConfig(config)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("自定义配置日志初始化成功")
	logger.Debug("调试信息现在可以显示了")
}

func exampleUsage() {
	fmt.Println("\n=== 示例4: 各种日志用法 ===")

	// 确保日志已初始化
	if err := logger.InitDevelopment(); err != nil {
		panic(err)
	}
	defer logger.Sync()

	// 1. 基本日志
	logger.Info("程序启动")
	logger.Infof("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	// 2. 带字段的日志（结构化日志）
	logger.Infow("用户登录",
		"user_id", 12345,
		"username", "john_doe",
		"ip", "192.168.1.100",
		"timestamp", time.Now().Unix(),
	)

	// 3. 错误处理
	if err := simulateError(); err != nil {
		logger.Errorw("处理用户请求失败",
			"error", err,
			"operation", "user_login",
			"user_id", 12345,
		)
	}

	// 4. 创建带固定字段的日志器
	requestLogger := logger.With("request_id", "req_789", "service", "auth")
	requestLogger.Infow("开始处理请求")
	requestLogger.Infow("请求处理完成", "duration_ms", 150)

	// 5. 性能监控
	start := time.Now()
	time.Sleep(100 * time.Millisecond) // 模拟耗时操作
	duration := time.Since(start)
	logger.Infow("操作完成", "duration", duration.String(), "duration_ms", duration.Milliseconds())
}

// 模拟业务函数
func processOrder(orderID string) {
	logger.Infow("开始处理订单",
		"order_id", orderID,
		"step", "validation",
	)

	// 模拟处理逻辑
	time.Sleep(50 * time.Millisecond)

	logger.Infow("订单处理完成",
		"order_id", orderID,
		"status", "success",
		"amount", 99.99,
	)
}

func processUser(userID string) {
	logger.Infow("处理用户信息",
		"user_id", userID,
		"action", "profile_update",
	)

	// 模拟业务逻辑
	if userID == "user_456" {
		logger.Warnw("用户状态异常",
			"user_id", userID,
			"reason", "account_locked",
		)
	}
}

func simulateError() error {
	return errors.New("数据库连接超时")
}
