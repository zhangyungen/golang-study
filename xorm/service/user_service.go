// service/user_service.go
package service

import (
	"errors"
	"zyj.com/golang-study/xorm/dao"
	"zyj.com/golang-study/xorm/model"
)

// UserService 用户Service
type UserService struct {
	*BaseService
	userDAO *dao.UserDAO
}

// 全局用户Service实例
var UserServiceInstance = &UserService{}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *model.User) error {
	//session := s.BaseService.getSession()

	// 数据验证
	if user.Name == "" {
		return errors.New("name is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}

	// 检查邮箱是否已存在
	exist, err := s.userDAO.ExistByEmail(user.Email)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("email already exists")
	}

	return s.userDAO.Create(user)
}

// GetUser 获取用户
func (s *UserService) GetUser(id int64) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.userDAO.GetByID(id)
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.userDAO.GetByEmail(email)
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(user *model.User) error {
	if user.ID <= 0 {
		return errors.New("invalid user id")
	}

	// 检查用户是否存在
	existing, err := s.userDAO.GetByID(user.ID)
	if err != nil {
		return err
	}

	// 如果邮箱有变更，检查新邮箱是否被其他用户使用
	if user.Email != existing.Email {
		exist, err := s.userDAO.ExistByEmail(user.Email)
		if err != nil {
			return err
		}
		if exist {
			return errors.New("email already used by another user")
		}
	}

	return s.userDAO.Update(user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int64) error {
	if id <= 0 {
		return errors.New("invalid user id")
	}
	return s.userDAO.Delete(id)
}

// ListUsers 用户列表
func (s *UserService) ListUsers(page, pageSize int) ([]*model.User, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	return s.userDAO.List(page, pageSize)
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
