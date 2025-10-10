// service/user_service.go
package service

import (
	"errors"
	"zyj.com/golang-study/xorm/base"
	"zyj.com/golang-study/xorm/dao"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/result"
)

// UserService 用户Service
type UserService struct {
	*base.BaseService[model.User, int64]
	userDAO *dao.UserDAO
}

// 全局用户Service实例
var UserServiceIns = &UserService{&base.BaseService[model.User, int64]{},
	dao.UserDaoIns}

// CreateUser 创建用户
func (us *UserService) CreateUser(user *model.User) error {
	// 数据验证
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	session := us.GetDBSession()
	defer us.ReturnDBSession(session)
	//todo check and model biz
	return us.userDAO.CreateUser(session, user)
}

// GetUserById 获取用户
func (us *UserService) GetUserById(id int64) (*model.User, error) {
	//todo check and model biz
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	return us.GetByID(id)
}

// GetUserByEmail 根据邮箱获取用户
func (us *UserService) GetUserByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	session := us.GetDBSession()
	defer us.ReturnDBSession(session)
	return us.userDAO.GetByEmail(session, email)
}

// UpdateUser 更新用户
func (us *UserService) UpdateUser(user *model.User) error {
	//todo check and model biz
	if user.Id <= 0 {
		return errors.New("invalid user id")
	}
	session := us.GetDBSession()
	defer us.ReturnDBSession(session)
	return us.userDAO.UpdateUser(session, user)
}

// DeleteUser 删除用户
func (us *UserService) DeleteUserById(id int64) error {
	if id <= 0 {
		return errors.New("invalid user id")
	}
	return us.DeleteById(id, &model.User{})
}

// ListUsers 用户列表
func (us *UserService) PageListUser(param *param.PageParam) (result.PageVO[model.User], error) {
	//session := us.getDBSession()
	//defer us.closeDBSession(session)
	//todo check and model biz
	return us.Page(param)
}

// ValidateUser 验证用户数据
func (us *UserService) ValidateUser(email, name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if email == "" {
		return errors.New("email is required")
	}
	// 可以添加更复杂的验证逻辑
	return nil
}
