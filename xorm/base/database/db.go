// common/database/db.go
package database

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/petermattis/goid"
	"log"
	"reflect"
	"sync"
	"time"
	"xorm.io/xorm"
)

var (
	engine        *xorm.Engine
	once          sync.Once
	sessionMap    = make(map[int64]*xorm.Session, 1)
	idKeyTableMap = make(map[interface{}]string, 1)
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
		//group, err := xorm.NewEngineGroup("postgres", nil)
		//group.Master()
		//group.Slave()
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
func getEngine() *xorm.Engine {
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

// getDBSession 创建新的数据库会话
func GetDBSession() *xorm.Session {
	id := goid.Get() // 直接获取当前 goroutine 的 Id
	if sessionMap[id] != nil {
		return sessionMap[id]
	}
	sessionMap[id] = getEngine().NewSession()
	return sessionMap[id]
}

// ReturnSession 关闭数据库回话
func ReturnSession(session *xorm.Session) error {
	id := goid.Get() // 直接获取当前 goroutine 的 Id
	if session == nil {
		return errors.New("argument session is nil ,don't close ")
	}
	if sessionMap[id] != session {
		return errors.New("close session failed，close session and goroutine session is not same，don use goroutine for you session code ")
	}
	if !session.IsInTx() {
		err := session.Close()
		delete(sessionMap, id)
		return err
	} else {
		log.Println("session is in transaction,don't close")
	}
	return nil
}

func WithTransaction(fn func() (err error)) (err error) {
	session := GetDBSession()
	defer func(session *xorm.Session) {
		err := ReturnSession(session)
		if err != nil {
			log.Println("close session failed:", err)
		}
	}(session)
	if !session.IsInTx() {
		if err := session.Begin(); err != nil {
			return err
		}
	}
	if err := fn(); err != nil {
		_ = session.Rollback()
		return err
	}
	return session.Commit()
}

func GetPrimaryKey[T any]() (string, error) {
	var t T
	r := reflect.TypeOf(t)
	typeString := r.PkgPath() + "." + r.Name()
	if idKeyTableMap[typeString] != "" {
		return idKeyTableMap[typeString], nil
	}
	engine := getEngine()
	//tableName := engine.TableName(model, true)
	table, err := engine.TableInfo(t)
	if err != nil {
		return "", err
	}

	for _, col := range table.Columns() {
		if col.IsPrimaryKey {
			idKeyTableMap[typeString] = col.Name
			return col.Name, nil
		}
	}
	return "", errors.New("no primary key found")
}
