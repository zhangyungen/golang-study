package base

//todo change  xorm to gorm
import (
	"zyj.com/golang-study/xorm/base/database"
)

// UserService 用户Service
type BaseCmdBiz struct {
}

// 全局用户Service实例
var BaseCmdBizIns = &BaseCmdBiz{}

func (biz *BaseCmdBiz) ExecuteTx(fn func() error) error {
	return database.WithTransaction(fn)
}
