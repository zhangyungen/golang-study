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
func (bs BaseService[T, K]) getDBSession() *xorm.Session {
	return database.GetDBSession()
}

func (bs BaseService[T, K]) closeDBSession(session *xorm.Session) {
	database.ReturnSession(session)
}

// GetByID 根据ID获取实体
func (bs *BaseService[T, K]) GetByID(id K) (*T, error) {
	return bs.baseDAO.GetByID(bs.getDBSession(), id)
}

// Update 实体
func (bs *BaseService[T, K]) UpdateById(id K, entity *T) error {
	return bs.baseDAO.UpdateById(bs.getDBSession(), id, entity)
}

// DeleteById 删除用户
func (bs *BaseService[T, K]) DeleteById(id int64, entity *T) error {
	return bs.baseDAO.DeleteById(bs.getDBSession(), id, entity)
}

func (bs *BaseService[T, K]) Create(entity *T) error {
	return bs.baseDAO.Insert(bs.getDBSession(), entity)
}

func (bs *BaseService[T, K]) PageList(param *param.PageParam) ([]T, error) {
	return bs.baseDAO.PageList(bs.getDBSession(), param)
}

// Count 统计数量
func (bs *BaseService[T, K]) Count(entity *T) (int64, error) {
	return bs.baseDAO.Count(bs.getDBSession(), entity)
}

func (bs *BaseService[T, K]) ExecuteTx(fn func() error) error {
	return database.WithTransaction(fn)
}
