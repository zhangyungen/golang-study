// dao/user_dao.go
package dao

import (
	"zyj.com/golang-study/xorm/base"
	"zyj.com/golang-study/xorm/model"
)

// UserDAO 用户DAO
type UserLoginLogDAO struct {
	*base.BaseDAO[model.UserLoginLog, int64]
}

// 全局用户DAO实例
var UserLoginLogDaoIns = &UserLoginLogDAO{&base.BaseDAO[model.UserLoginLog, int64]{}}
