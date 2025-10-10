package biz

import (
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/result"
	"zyj.com/golang-study/xorm/service"
)

// UserQueryBiz 用户QueryBizService
type UserQueryBiz struct {
	userService *service.UserService
}

// 全局UserQueryBizIns实例
var UserQueryBizIns = &UserQueryBiz{userService: service.UserServiceIns}

func (biz *UserQueryBiz) PageUser(param *param.PageParam) (result.PageVO[model.User], error) {
	return biz.userService.Page(param)
}

func (biz *UserQueryBiz) ListUserByIds(ids []int64) ([]model.User, error) {
	return biz.userService.ListByIds(ids)
}
