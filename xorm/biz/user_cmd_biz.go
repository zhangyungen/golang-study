// service/user_service.go
package biz

import (
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/service"
)

// UserService 用户Service
type UserCmdBiz struct {
	*service.UserService
}

// 全局用户Service实例
var UserCmdBizIns = &UserCmdBiz{service.UserServiceIns}

// CreateUser 创建用户
func (biz *UserCmdBiz) CreateUser(user *param.UserCreate) error {
	return service.UserServiceIns.CreateUser(param.ConvertToModel(user))
}
