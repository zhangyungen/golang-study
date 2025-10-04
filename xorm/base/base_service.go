// service/base_service.go
package base

import (
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/base/database"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/result"
)

// BaseService 基础Service
type BaseService[T any, K any] struct {
	baseDAO *BaseDAO[T, K]
}

// 全局基础Service实例
//var BaseServiceInstance = &BaseService{ &dao.BaseDAO{}}

// GetDBSession 获取数据库会话
func (bs BaseService[T, K]) GetDBSession() *xorm.Session {
	return database.GetDBSession()
}

func (bs BaseService[T, K]) ReturnDBSession(session *xorm.Session) {
	database.ReturnSession(session)
}

// GetByID 根据ID获取实体
func (bs *BaseService[T, K]) GetByID(id K) (*T, error) {
	return bs.baseDAO.GetByID(bs.GetDBSession(), id)
}

// Update 实体
func (bs *BaseService[T, K]) UpdateById(id K, entity *T) error {
	return bs.baseDAO.UpdateById(bs.GetDBSession(), id, entity)
}

// DeleteById 删除用户
func (bs *BaseService[T, K]) DeleteById(id int64, entity *T) error {
	return bs.baseDAO.DeleteById(bs.GetDBSession(), id, entity)
}

func (bs *BaseService[T, K]) Create(entity *T) error {
	return bs.baseDAO.Insert(bs.GetDBSession(), entity)
}

func (bs *BaseService[T, K]) Page(param *param.PageParam) (result.PageVO[T], error) {
	return bs.baseDAO.Page(bs.GetDBSession(), param)
}

func (bs *BaseService[T, K]) ListByIds(ids []K) ([]T, error) {
	return bs.baseDAO.ListByIds(bs.GetDBSession(), ids)
}

// Count 统计数量
func (bs *BaseService[T, K]) Count(entity *T) (int64, error) {
	return bs.baseDAO.Count(bs.GetDBSession(), entity)
}

func (bs *BaseService[T, K]) ExecuteTx(fn func() error) error {
	return database.WithTransaction(fn)
}
