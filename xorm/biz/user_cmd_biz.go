package biz

import (
	"errors"
	"time"
	"zyj.com/golang-study/util/obj"
	"zyj.com/golang-study/xorm/base"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/service"
)

// UserCmdBiz 用户CMDService
type UserCmdBiz struct {
	*base.BaseCmdBiz

	userService         *service.UserService
	userLoginLogService *service.UserLoginLogService
}

// 全局UserCmdBizIns实例
var UserCmdBizIns = &UserCmdBiz{base.BaseCmdBizIns,
	service.UserServiceIns, service.UserLoginLogServiceIns}

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

// LogIn 登录
func (biz *UserCmdBiz) LogIn(param *param.UserLogin) (bool, error) {
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
