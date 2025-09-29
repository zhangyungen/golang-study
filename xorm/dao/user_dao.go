// dao/user_dao.go
package dao

import (
	"errors"
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/common"
	"zyj.com/golang-study/xorm/model"
)

// UserDAO 用户DAO
type UserDAO struct {
	*BaseDAO[model.User, int64]
}

// 全局用户DAO实例
var UserDaoInstance = &UserDAO{&BaseDAO[model.User, int64]{}}

// Create 创建用户
func (d *UserDAO) Create(session *xorm.Session, user *model.User) error {
	//_, err := session.Insert(user)
	//return err
	return d.BaseDAO.Create(session, user)
}

// GetByID 根据ID获取用户
func (d *UserDAO) GetByID(session *xorm.Session, id int64) (*model.User, error) {
	//var user model.User
	//has, err := session.ID(id).Get(&user)
	//if err != nil {
	//	return nil, err
	//}
	//if !has {
	//	return nil, errors.New("user not found")
	//}
	//return &user, nil
	return d.BaseDAO.GetByID(session, id)
}

// GetByEmail 根据邮箱获取用户
func (d *UserDAO) GetByEmail(session *xorm.Session, email string) (*model.User, error) {
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
func (d *UserDAO) Update(session *xorm.Session, user *model.User) error {
	//_, err := session.ID(user.ID).Update(user)
	//return err
	return d.BaseDAO.Update(session, user)
}

// Delete 删除用户
func (d *UserDAO) Delete(session *xorm.Session, id int64) error {
	//user := &model.User{ID: id}
	//_, err := session.ID(id).Delete(user)
	//return err
	return d.BaseDAO.Delete(session, id)
}

// List 用户列表
func (d *UserDAO) PageList(session *xorm.Session, param *common.PageParam) ([]*model.User, error) {
	//var users []*model.User
	//err := session.Limit(param.PageSize, param.PageSize*(param.Page-1)).Find(&users)
	//return users, err

	return d.BaseDAO.PageList(session, param)
}

// ExistByEmail 检查邮箱是否存在
func (d *UserDAO) ExistByEmail(session *xorm.Session, email string) (bool, error) {
	return session.Where("email = ?", email).Exist(&model.User{})
}
