// dao/base_dao.go
package base

import (
	"errors"
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/base/database"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/result"
)

// BaseDAO 基础DAO
type BaseDAO[T any, K any] struct {
}

// 全局基础DAO实例
//var BaseDaoInstance = &BaseDAO{}

// GetByID 根据ID获取用户
func (bd *BaseDAO[T, K]) GetByID(session *xorm.Session, id K) (*T, error) {
	var entity T
	has, err := session.ID(id).Get(&entity)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("entity not found")
	}
	return &entity, nil
}

// Update 更新用户
func (bd *BaseDAO[T, K]) UpdateById(session *xorm.Session, id K, entity *T) error {
	_, err := session.ID(id).Update(entity)
	return err
}

// DeleteById 删除用户
func (bd *BaseDAO[T, K]) DeleteById(session *xorm.Session, id int64, entity *T) error {
	_, err := session.ID(id).Delete(entity)
	return err
}

func (bd *BaseDAO[T, K]) Insert(session *xorm.Session, entity *T) error {
	_, err := session.Insert(entity)
	return err
}

func (bd *BaseDAO[T, K]) PageList(session *xorm.Session, param *param.PageParam) (result.PageVO[T], error) {
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.PageSize <= 0 {
		param.PageSize = 10
	}
	var entities []T
	err := session.Limit(param.PageSize, param.PageSize*(param.Page-1)).Find(&entities)
	var t T
	count, err := session.Count(t)
	if err != nil {
		return result.PageVO[T]{}, err
	}
	return result.Convert2PageVO[T](param, count, entities), nil
}

// Count 统计数量
func (bd *BaseDAO[T, k]) Count(session *xorm.Session, entity *T) (int64, error) {
	return session.Count(entity)
}

func (bd *BaseDAO[T, K]) ListByIds(session *xorm.Session, ids []K) ([]T, error) {
	var entities []T
	key, err := database.GetPrimaryKey[T]()
	err = session.In(key, ids).Find(&entities)
	return entities, err
}
