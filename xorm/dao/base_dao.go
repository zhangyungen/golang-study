// dao/base_dao.go
package dao

import (
	"errors"
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/common"
	"zyj.com/golang-study/xorm/model"
)

// BaseDAO 基础DAO
type BaseDAO[T any, K any] struct{}

// 全局基础DAO实例
//var BaseDaoInstance = &BaseDAO{}

// GetByID 根据ID获取用户
func (d *BaseDAO[T, K]) GetByID(session *xorm.Session, id K) (*T, error) {
	var user T
	has, err := session.ID(id).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// Update 更新用户
func (d *BaseDAO[T, K]) Update(session *xorm.Session, user *T) error {
	_, err := session.Update(user)
	return err
}

// Delete 删除用户
func (d *BaseDAO[T, K]) Delete(session *xorm.Session, id int64) error {
	user := &model.User{ID: id}
	_, err := session.ID(id).Delete(user)
	return err
}

func (d *BaseDAO[T, K]) Create(session *xorm.Session, user *T) error {
	_, err := session.Insert(user)
	return err
}

func (d *BaseDAO[T, K]) PageList(session *xorm.Session, param *common.PageParam) ([]*T, error) {
	var users []*T
	err := session.Limit(param.PageSize, param.PageSize*(param.Page-1)).Find(&users)
	return users, err
}

// Count 统计数量
func (d *BaseDAO[T, k]) Count(session *xorm.Session) (int64, error) {
	return session.Count(&model.User{})
}
