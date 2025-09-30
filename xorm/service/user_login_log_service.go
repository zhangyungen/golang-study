// service/user_service.go
package service

import (
	"zyj.com/golang-study/xorm/dao"
	"zyj.com/golang-study/xorm/model"
)

// UserService 用户Service
type UserLoginLogService struct {
	*BaseService[model.UserLoginLog, int64]
	userLoginLogDAO *dao.UserLoginLogDAO
}

// 全局用户Service实例
var UserLoginLogServiceIns = &UserLoginLogService{&BaseService[model.UserLoginLog, int64]{}, dao.UserLoginLogDaoIns}
