// common/database/xorm.go
package database

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/petermattis/goid"
	"log"
	"sync"
	"time"
	"xorm.io/xorm"
)

var (
	engine     *xorm.Engine
	once       sync.Once
	sessionMap = make(map[int64]*xorm.Session, 10)
)

// Init 初始化数据库连接
func Init(driver, dsn string) error {
	var err error

	once.Do(func() {
		engine, err = xorm.NewEngine(driver, dsn)
		if err != nil {
			log.Println("Failed to connect to database:", err)
			panic(err)
			return
		}

		// 配置连接池
		engine.SetMaxOpenConns(100)
		engine.SetMaxIdleConns(10)
		engine.SetConnMaxLifetime(time.Hour)

		// 显示SQL日志（开发环境）
		engine.ShowSQL(true)

		// 测试连接
		if err = engine.Ping(); err != nil {
			log.Println("Failed to ping database:", err)
			panic(err)
		}
		log.Println("Database connected successfully")
	})

	return err
}

// GetEngine 获取数据库引擎（单例）
func GetEngine() *xorm.Engine {
	if engine == nil {
		panic("database engine not initialized, please call Init first")
	}
	return engine
}

// Close 关闭数据库连接
func CloseEngine() error {
	if engine != nil {
		return engine.Close()
	}
	return nil
}

// NewSession 创建新的数据库会话
func NewSession() *xorm.Session {
	id := goid.Get() // 直接获取当前 goroutine 的 ID
	sessionMap[id] = GetEngine().NewSession()
	return sessionMap[id]
}

// CloseSession 关闭数据库回话
func CloseSession(session *xorm.Session) error {
	id := goid.Get() // 直接获取当前 goroutine 的 ID
	if sessionMap[id] != session {
		return errors.New("关闭session 失败，不是同一个session，请忽使用 goroutine ")
	}
	session.Commit()
	err := session.Close()
	delete(sessionMap, id)
	return err
}
