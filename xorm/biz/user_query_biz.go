// service/user_service.go
package biz

import (
	"errors"
	"time"
	"zyj.com/golang-study/xorm/dao"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/service"
)

// UserService 用户Service
type UserQueryBiz struct {
	*service.BaseService[model.User, int64]
	userLoginLogService *service.UserLoginLogService
	userDAO             *dao.UserDAO
	userLoginLogDAO     *dao.UserLoginLogDAO
}

// 全局UserQueryBizIns实例
var UserQueryBizIns = &UserQueryBiz{BaseService: &service.BaseService[model.User, int64]{},
	userDAO: dao.UserDaoIns, userLoginLogDAO: dao.UserLoginLogDaoIns,
	userLoginLogService: service.UserLoginLogServiceIns}

// LogIn 登录
func (biz *UserQueryBiz) LogIn(param *param.UserLogin) (bool, error) {
	if user, err := service.UserServiceIns.GetUserByEmail(param.Email); err != nil {
		return false, err
	} else if user != nil {
		if param.Pwd != user.Pwd {
			return false, errors.New("用户名或密码错误")
		} else {
			biz.userLoginLogService.Create(&model.UserLoginLog{
				LoginIp:     "192.168.1.1",
				LoginTime:   time.Time{},
				UserId:      user.Id,
				CreatedTime: time.Time{},
				UpdatedTime: time.Time{},
				DeletedTime: time.Time{},
			})
			return true, nil
		}
	} else {
		return false, errors.New("用户不存在")
	}
}
