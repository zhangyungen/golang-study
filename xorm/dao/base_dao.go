// dao/base_dao.go
package dao

import (
	"errors"
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/common"
)

// BaseDAO 基础DAO
type BaseDAO[T any, K any] struct{}

// 全局基础DAO实例
//var BaseDaoInstance = &BaseDAO{}

// GetByID 根据ID获取用户
func (d *BaseDAO[T, K]) GetByID(session *xorm.Session, id K) (*T, error) {
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
func (d *BaseDAO[T, K]) UpdateById(session *xorm.Session, id K, entity *T) error {
	_, err := session.ID(id).Update(entity)
	return err
}

// DeleteById 删除用户
func (d *BaseDAO[T, K]) DeleteById(session *xorm.Session, id int64, entity *T) error {
	_, err := session.ID(id).Delete(entity)
	return err
}

func (d *BaseDAO[T, K]) Insert(session *xorm.Session, entity *T) error {
	_, err := session.Insert(entity)
	return err
}

func (d *BaseDAO[T, K]) PageList(session *xorm.Session, param *common.PageParam) ([]*T, error) {
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.PageSize <= 0 {
		param.PageSize = 10
	}
	var entitys []*T
	err := session.Limit(param.PageSize, param.PageSize*(param.Page-1)).Find(&entitys)
	return entitys, err
}

// Count 统计数量
func (d *BaseDAO[T, k]) Count(session *xorm.Session, entity *T) (int64, error) {
	return session.Count(entity)
}
