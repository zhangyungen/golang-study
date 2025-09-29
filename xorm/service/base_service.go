// service/base_service.go
package service

import (
	"xorm.io/xorm"
	"zyj.com/golang-study/xorm/dao"
)

// BaseService 基础Service
type BaseService struct {
	baseDAO *dao.BaseDAO
}

// 全局基础Service实例
var BaseServiceInstance = &BaseService{}

func (BaseService) getSession() *xorm.Session {
	return dao.BaseDaoInstance.GetDBSession()

}
