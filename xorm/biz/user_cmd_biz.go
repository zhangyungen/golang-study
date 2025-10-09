package biz

import (
	"zyj.com/golang-study/util/obj"
	"zyj.com/golang-study/xorm/base"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/service"
)

// UserCmdBiz 用户CMDService
type UserCmdBiz struct {
	*base.BaseCmdBiz

	userService *service.UserService
}

// 全局UserCmdBizIns实例
var UserCmdBizIns = &UserCmdBiz{base.BaseCmdBizIns,
	service.UserServiceIns}

// CreateUser 创建用户
func (biz *UserCmdBiz) CreateUser(user *param.UserCreate) error {
	return biz.ExecuteTx(func() error {
		toObj := obj.ObjToObj[model.User](user)
		return service.UserServiceIns.CreateUser(toObj)
	})
}

func (biz *UserCmdBiz) UpdateUser(user *param.UserUpdate) error {
	userModel := obj.ObjToObj[model.User](user)
	return service.UserServiceIns.UpdateUser(userModel)
}
