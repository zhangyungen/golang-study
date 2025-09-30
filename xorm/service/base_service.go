// service/base_service.go
package service

import (
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/dao"
	"zyj.com/golang-study/xorm/database"
)

// BaseService 基础Service
type BaseService[T any, K any] struct {
	baseDAO *dao.BaseDAO[T, K]
}

// 全局基础Service实例
//var BaseServiceInstance = &BaseService{ &dao.BaseDAO{}}

// GetDBSession 获取数据库会话
func (b BaseService[T, K]) getDBSession() *xorm.Session {
	return database.GetDBSession()
}

func (b BaseService[T, K]) closeDBSession(session *xorm.Session) {
	database.CloseSession(session)
}
