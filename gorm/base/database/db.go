package database

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/puzpuzpuz/xsync/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"reflect"
	"sync"
	"zyj.com/golang-study/tslog"
)

var (
	dbEngine *gorm.DB
	once     sync.Once
	//sessionMap    = make(map[int64]*xorm.Session, 1)
	//idKeyTableMap = make(map[string]string, 1)
	//sessionMap    sync.Map
	//idKeyTableMap sync.Map

	//idKeyTableMap = cmap.New[string]()
	//sessionMap    = cmap.New[*xorm.Session]()

	sessionMap     = xsync.NewMapOf[int64, *gorm.DB]()
	isInTx         = xsync.NewMapOf[int64, bool]()
	idKeyTableMap  = xsync.NewMapOf[string, string]()
	idKeyTableMap2 = new(sync.Map)
)

// Init 初始化数据库连接
func Init(driver, dsn string) error {
	var err error

	once.Do(func() {

		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{Logger: tslog.New(tslog.Logger)})

		dbEngine = db
		if err != nil {
			log.Fatal("failed to connect database: %v", err)
		}
	})

	return err
}

// GetEngine 获取数据库引擎（单例）
func getDBEngine() *gorm.DB {
	if dbEngine == nil {
		panic("database engine not initialized, please call Init first")
	}
	return dbEngine
}

// Close 关闭数据库连接
func CloseEngine() error {
	if dbEngine != nil {
		dbEngine = nil
	}
	return nil
}

// todo  待测试
func WithTransaction(ctx context.Context, fn func(txCtx context.Context) (err error)) (err error) {
	return getDBEngine().Transaction(func(tx *gorm.DB) error {
		return fn(ctx)
	})
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

	schema, err := schema.Parse(t, idKeyTableMap2, getDBEngine().NamingStrategy)
	if err != nil {
		return "", err
	}
	// 查找主键字段
	for _, field := range schema.Fields {
		if field.PrimaryKey {
			// 返回主键在数据库中的列名
			return field.DBName, nil
		}
	}
	return "", fmt.Errorf("primary key not found for model %T", t)

}

func QueryRowsBySql[T any](sql string) ([]T, error) {
	var rows = make([]T, 0)
	result := gorm.WithResult()
	err := gorm.G[T](getDBEngine(), result).Exec(context.Background(), sql, rows)

	if err != nil {
		return nil, err
	}
	return rows, nil
}

func QueryRowBySql[T any](sql string) (*T, error) {
	var t T
	result := gorm.WithResult()
	err := gorm.G[T](getDBEngine(), result).Exec(context.Background(), sql, t)
	//.Error()
	if err != nil {
		return &t, err
	}
	return &t, nil
}
