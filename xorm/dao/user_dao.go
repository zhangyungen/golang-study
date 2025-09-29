// dao/user_dao.go
package dao

import (
	"errors"
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/model"
)

// UserDAO 用户DAO
type UserDAO struct {
	*BaseDAO
}

// 全局用户DAO实例
var UserDaoInstance = &UserDAO{}

// Create 创建用户
func (d *UserDAO) Create(user *model.User) error {
	session := d.GetDBSession()
	defer session.Close()

	_, err := session.Insert(user)
	return err
}

// GetByID 根据ID获取用户
func (d *UserDAO) GetByID(id int64) (*model.User, error) {
	session := d.GetDBSession()
	defer session.Close()

	var user model.User
	has, err := session.ID(id).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (d *UserDAO) GetByEmail(email string) (*model.User, error) {
	session := d.GetDBSession()
	defer session.Close()

	var user model.User
	has, err := session.Where("email = ?", email).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// Update 更新用户
func (d *UserDAO) Update(user *model.User) error {
	session := d.GetDBSession()
	defer session.Close()

	_, err := session.ID(user.ID).Update(user)
	return err
}

// Delete 删除用户
func (d *UserDAO) Delete(id int64) error {
	session := d.GetDBSession()
	defer session.Close()

	user := &model.User{ID: id}
	_, err := session.ID(id).Delete(user)
	return err
}

// List 用户列表
func (d *UserDAO) List(page, pageSize int) ([]*model.User, error) {
	session := d.GetDBSession()
	defer session.Close()

	var users []*model.User
	err := session.Limit(pageSize, (page-1)*pageSize).Find(&users)
	return users, err
}

// Count 统计用户数量
func (d *UserDAO) Count() (int64, error) {
	session := d.GetDBSession()
	defer session.Close()

	return session.Count(&model.User{})
}

// ExistByEmail 检查邮箱是否存在
func (d *UserDAO) ExistByEmail(email string) (bool, error) {
	session := d.GetDBSession()
	defer session.Close()

	return session.Where("email = ?", email).Exist(&model.User{})
}

// CreateWithSession 使用外部会话创建用户（用于事务）
func (d *UserDAO) CreateWithSession(session *xorm.Session, user *model.User) error {
	_, err := session.Insert(user)
	return err
}
