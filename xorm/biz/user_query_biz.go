// service/user_service.go
package biz

import (
	"errors"
	"zyj.com/golang-study/xorm/dao"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/service"
)

// UserService 用户Service
type UserQueryBiz struct {
	*service.BaseService[model.User, int64]
	userDAO *dao.UserDAO
}

// 全局用户Service实例
var UserQueryBizIns = &UserQueryBiz{&service.BaseService[model.User, int64]{}, dao.UserDaoInstance}

// LogIn 登录
func (biz *UserQueryBiz) LogIn(param *param.UserLogin) (bool, error) {
	if user, err := service.UserServiceIns.GetUserByEmail(param.Email); err != nil {
		return false, err
	} else if user != nil {
		if param.Pwd != user.Pwd {
			return false, errors.New("用户名或密码错误")
		} else {
			return true, nil
		}
	} else {
		return false, errors.New("用户不存在")
	}
}
