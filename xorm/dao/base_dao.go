// dao/base_dao.go
package dao

import (
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/database"
)

// BaseDAO 基础DAO
type BaseDAO struct{}

// 全局基础DAO实例
var BaseDaoInstance = &BaseDAO{}

// GetDBSession 获取数据库会话
func (d *BaseDAO) GetDBSession() *xorm.Session {
	return database.NewSession()
}

// WithTransaction 执行事务
func (d *BaseDAO) WithTransaction(fn func(*xorm.Session) error) error {
	session := d.GetDBSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	if err := fn(session); err != nil {
		_ = session.Rollback()
		return err
	}

	return session.Commit()
}
