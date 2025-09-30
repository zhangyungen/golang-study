// service/user_service.go
package biz

import (
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/service"
)

// UserService 用户Service
type UserCmdBiz struct {
	*BaseCmdBiz
	*service.UserService
}

// 全局UserCmdBizIns实例
var UserCmdBizIns = &UserCmdBiz{BaseCmdBizIns, service.UserServiceIns}

// CreateUser 创建用户
func (biz *UserCmdBiz) CreateUser(user *param.UserCreate) error {
	return biz.ExecuteTx(func() error {
		return service.UserServiceIns.CreateUser(param.ConvertToModel(user))
	})
}
