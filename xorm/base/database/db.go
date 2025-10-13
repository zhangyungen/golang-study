package database

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/petermattis/goid"
	"github.com/puzpuzpuz/xsync/v3"
	"log"
	"reflect"
	"sync"
	"time"
	"xorm.io/xorm"
	"zyj.com/golang-study/util/objutil"
)

var (
	engine *xorm.Engine
	once   sync.Once
	//sessionMap    = make(map[int64]*xorm.Session, 1)
	//idKeyTableMap = make(map[string]string, 1)
	//sessionMap    sync.Map
	//idKeyTableMap sync.Map

	//idKeyTableMap = cmap.New[string]()
	//sessionMap    = cmap.New[*xorm.Session]()

	sessionMap    = xsync.NewMapOf[int64, *xorm.Session]()
	idKeyTableMap = xsync.NewMapOf[string, string]()
)

// Init 初始化数据库连接
func Init(driver, dsn string) error {
	var err error

	once.Do(func() {
		engine, err = xorm.NewEngine(driver, dsn)
		if err != nil {
			log.Println("Failed to connect to database:", err)
			panic(err)
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
	//value, ok := sessionMap.Load(id)
	value, ok := sessionMap.Load(id)
	if ok {
		//return value.(*xorm.Session)
		return value
	}
	session := getEngine().NewSession()
	sessionMap.Store(id, session)
	return session
}

// ReturnSession 关闭数据库回话
func ReturnSession(session *xorm.Session) {
	id := goid.Get() // 直接获取当前 goroutine 的 Id
	holdSession, ok := sessionMap.Load(id)
	if !ok {
		log.Println("argument session is not ok ,don't close ")
	}
	if holdSession == nil {
		log.Println("argument session is nil ,don't close ")
	}
	if holdSession != session {
		log.Println("close session failed，close session and goroutine session is not same，don use goroutine for you session code ")
	}
	if !session.IsInTx() {
		err := session.Close()
		sessionMap.Delete(id)
		if err != nil {
			log.Println("close session failed:", err)
		}
	} else {
		log.Println("session is in transaction,don't close")
	}
}

func WithTransaction(fn func() (err error)) (err error) {
	session := GetDBSession()
	defer ReturnSession(session)
	if !session.IsInTx() {
		if err := session.Begin(); err != nil {
			return err
		}
	}
	if err := fn(); err != nil {
		err2 := session.Rollback()
		return errors.Join(err, err2)
	}
	return session.Commit()
}

func WithTransactionSession(fn func(session *xorm.Session) (err error)) (err error) {
	session := GetDBSession()
	defer ReturnSession(session)
	if !session.IsInTx() {
		if err := session.Begin(); err != nil {
			return err
		}
	}
	if err := fn(session); err != nil {
		err2 := session.Rollback()
		return errors.Join(err, err2)
	}
	return session.Commit()
}

func GetPrimaryKey[T any]() (string, error) {
	var t T
	r := reflect.TypeOf(t)
	typeString := r.PkgPath() + "." + r.Name()
	value, b := idKeyTableMap.Load(typeString)
	if b {
		//return value.(string), nil
		return value, nil
	}
	engine := getEngine()
	//tableName := engine.TableName(model, true)
	table, err := engine.TableInfo(t)
	if err != nil {
		return "", err
	}

	for _, col := range table.Columns() {
		if col.IsPrimaryKey {
			idKeyTableMap.Store(typeString, col.Name)
			return col.Name, nil
		}
	}
	return "", errors.New("no primary key found")
}

func QueryRowsBySql[T any](session *xorm.Session, sql string) ([]T, error) {
	exec, err := session.QueryInterface(sql)
	if err != nil {
		return nil, err
	}
	count := len(exec)
	var rows = make([]T, 0, count)
	for _, v := range exec {
		t := objutil.MapToObjByStr[T](v)
		rows = append(rows, *t)
	}
	return rows, nil
}

func QueryRowBySql[T any](session *xorm.Session, sql string) (*T, error) {
	exec, err := session.QueryInterface(sql)
	var t T
	if err != nil {
		return &t, err
	}
	if len(exec) > 1 {
		return &t, errors.New("multi rows found")
	}
	if len(exec) == 0 {
		return &t, errors.New("not rows found")
	}
	for _, v := range exec {
		t = *objutil.MapToObjByStr[T](v)
	}
	return &t, nil
}

func ExecuteTxSession(fn func(session *xorm.Session) error) error {
	session := GetDBSession()
	defer ReturnSession(session)
	return fn(session)
}
