// service/base_service.go
package service

import (
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/dao"
	"zyj.com/golang-study/xorm/database"
	"zyj.com/golang-study/xorm/param"
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

// GetByID 根据ID获取实体
func (b *BaseService[T, K]) GetByID(id K) (*T, error) {
	return b.baseDAO.GetByID(b.getDBSession(), id)
}

// Update 实体
func (b *BaseService[T, K]) UpdateById(id K, entity *T) error {
	return b.baseDAO.UpdateById(b.getDBSession(), id, entity)
}

// DeleteById 删除用户
func (b *BaseService[T, K]) DeleteById(id int64, entity *T) error {
	return b.baseDAO.DeleteById(b.getDBSession(), id, entity)
}

func (b *BaseService[T, K]) Create(entity *T) error {
	return b.baseDAO.Insert(b.getDBSession(), entity)
}

func (b *BaseService[T, K]) PageList(param *param.PageParam) ([]*T, error) {
	return b.baseDAO.PageList(b.getDBSession(), param)
}

// Count 统计数量
func (b *BaseService[T, K]) Count(entity *T) (int64, error) {
	return b.baseDAO.Count(b.getDBSession(), entity)
}

func (biz *BaseService[T, K]) ExecuteTx(fn func() error) error {
	return database.WithTransaction(fn)
}
