// service/user_service.go
package service

import (
	"errors"
	"zyj.com/golang-study/xorm/dao"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/result"
)

// UserService 用户Service
type UserService struct {
	*BaseService[model.User, int64]
	userDAO *dao.UserDAO
}

// 全局用户Service实例
var UserServiceIns = &UserService{&BaseService[model.User, int64]{}, dao.UserDaoIns}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *model.User) error {
	// 数据验证
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	session := s.getDBSession()
	defer s.closeDBSession(session)
	return s.userDAO.CreateUser(session, user)
}

// GetUser 获取用户
func (s *UserService) GetUser(id int64) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	session := s.getDBSession()
	defer s.closeDBSession(session)
	return s.userDAO.GetByID(session, id)
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	session := s.getDBSession()
	defer s.closeDBSession(session)
	return s.userDAO.GetByEmail(session, email)
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *model.User) error {
	if user.Id <= 0 {
		return errors.New("invalid user id")
	}
	session := s.getDBSession()
	defer s.closeDBSession(session)
	return s.userDAO.Update(session, user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUserById(id int64) error {
	if id <= 0 {
		return errors.New("invalid user id")
	}
	session := s.getDBSession()
	defer s.closeDBSession(session)
	return s.userDAO.DeleteById(session, id, &model.User{})
}

// ListUsers 用户列表
func (s *UserService) PageList(param *param.PageParam) (result.PageVO[model.User], error) {
	session := s.getDBSession()
	defer s.closeDBSession(session)

	list, err := s.userDAO.PageList(session, param)
	if err != nil {
		return result.PageVO[model.User]{}, err
	}
	count, err := session.Count(&model.User{})
	if err != nil {
		return result.PageVO[model.User]{}, err
	}
	return result.Convert2PageVO(param, count, list), nil
}

// ValidateUser 验证用户数据
func (s *UserService) ValidateUser(email, name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	if email == "" {
		return errors.New("email is required")
	}
	// 可以添加更复杂的验证逻辑
	return nil
}
