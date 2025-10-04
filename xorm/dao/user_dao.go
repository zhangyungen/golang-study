// dao/user_dao.go
package dao

import (
	"errors"
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/model"
)

// UserDAO 用户DAO
type UserDAO struct {
	*BaseDAO[model.User, int64]
}

// 全局用户DAO实例
var UserDaoIns = &UserDAO{&BaseDAO[model.User, int64]{}}

// CreateUser 创建用户
func (ud *UserDAO) CreateUser(session *xorm.Session, user *model.User) error {
	// 检查邮箱是否已存在
	exist, err := ud.ExistByEmail(session, user.Email)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("email already exists")
	}
	return ud.BaseDAO.Insert(session, user)
}

// GetByEmail 根据邮箱获取用户
func (ud *UserDAO) GetByEmail(session *xorm.Session, email string) (*model.User, error) {
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
func (ud *UserDAO) Update(session *xorm.Session, user *model.User) error {
	//检查用户是否存在
	existing, err := ud.GetByID(session, user.Id)
	if err != nil {
		return err
	}
	// 如果邮箱有变更，检查新邮箱是否被其他用户使用
	if user.Email != existing.Email {
		exist, err := ud.ExistByEmailAndNotId(session, user)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("email already used by another user")
		}
	}
	return ud.BaseDAO.UpdateById(session, user.Id, user)
}

// ExistByEmail 检查邮箱是否存在
func (ud *UserDAO) ExistByEmail(session *xorm.Session, email string) (bool, error) {
	return session.Where("email = ?", email).Exist(&model.User{})
}

// ExistByEmail 检查邮箱是否存在
func (ud *UserDAO) ExistByEmailAndNotId(session *xorm.Session, user *model.User) (bool, error) {
	return session.Where("email = ? and id != ?", user.Email, user.Id).Exist(&model.User{})
}

//
//// DeleteById 删除用户
//func (d *UserDAO) DeleteById(session *xorm.Session, id int64) error {
//	//user := &model.User{Id: id}
//	//_, err := session.Id(id).DeleteById(user)
//	//return err
//	return d.BaseDAO.DeleteById(session, id)
//}

// List 用户列表
//func (d *UserDAO) PageList(session *xorm.Session, param *common.PageParam) ([]*model.User, error) {
//	//var users []*model.User
//	//err := session.Limit(param.PageSize, param.PageSize*(param.Page-1)).Find(&users)
//	//return users, err
//
//	return d.BaseDAO.PageList(session, param)
//}
