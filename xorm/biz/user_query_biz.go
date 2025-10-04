package biz

import (
	"errors"
	"time"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/result"
	"zyj.com/golang-study/xorm/service"
)

// UserQueryBiz 用户QueryBizService
type UserQueryBiz struct {
	userService         *service.UserService
	userLoginLogService *service.UserLoginLogService
}

// 全局UserQueryBizIns实例
var UserQueryBizIns = &UserQueryBiz{userService: service.UserServiceIns,
	userLoginLogService: service.UserLoginLogServiceIns}

func (biz *UserQueryBiz) PageUser(param *param.PageParam) (result.PageVO[model.User], error) {
	return biz.userService.Page(param)
}

func (biz *UserQueryBiz) ListUserByIds(ids []int64) ([]model.User, error) {
	return biz.userService.ListByIds(ids)
}

// LogIn 登录
func (biz *UserQueryBiz) LogIn(param *param.UserLogin) (bool, error) {
	if user, err := biz.userService.GetUserByEmail(param.Email); err != nil {
		return false, err
	} else if user != nil {
		if param.Pwd != user.Pwd {
			return false, errors.New("用户名或密码错误")
		} else {
			err := biz.userLoginLogService.Create(&model.UserLoginLog{
				LoginIp:     "192.168.1.1",
				LoginTime:   time.Now(),
				UserId:      user.Id,
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
				DeletedTime: time.Now(),
			})
			if err != nil {
				return false, err
			}
			return true, nil
		}
	} else {
		return false, errors.New("用户不存在")
	}
}
