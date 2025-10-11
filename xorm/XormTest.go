package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"xorm.io/xorm"
	"zyj.com/golang-study/util/obj"
	"zyj.com/golang-study/xorm/base/database"
	"zyj.com/golang-study/xorm/biz"
	"zyj.com/golang-study/xorm/dao/sql"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
	"zyj.com/golang-study/xorm/service"
)

func main() {

	//初始化数据库
	err := database.Init("mysql", "root:zj123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	sqlStr := sql.GetUserLoginSql(param.UserCreate{Name: "192.168.1.3"})
	log.Println(sqlStr)
	var userLoginLogs []model.UserLoginLog
	err = service.UserServiceIns.ExecuteTxSession(func(session *xorm.Session) error {
		userLoginLogs, err = database.QueryRowsBySql[model.UserLoginLog](session, sqlStr)
		return err
	})
	if err != nil {
		return
	}
	log.Println("userLoginLogs:", userLoginLogs)

	//database.GetEngine().Sync2(new(model.User))
	//业务代码开始
	user, err := biz.UserCmdBizIns.CreateUser(&param.UserCreate{Name: "zyj2fdsa", Email: "zyj0009@163kkkk.com"})

	if err != nil {
		log.Println("error", err)
	} else {
		log.Println("create user", obj.ObjToJsonStr(user))
	}

	err = biz.UserCmdBizIns.UpdateUserById(&param.UserUpdate{Id: 1, Name: "zyj2fdsa", Email: "zyj000@163kkkk.com"})

	err = biz.UserCmdBizIns.UpdateUserById(&param.UserUpdate{Name: "updatetest"})

	if err != nil {
		log.Println(err)
	}

	in, err := biz.UserCmdBizIns.LogIn(&param.UserLogin{Email: "zyj@163fff.com", Pwd: "123456"})

	if err != nil {
		log.Println(err)
	} else {
		log.Println("登录状态", in)
	}

	pages, err := biz.UserQueryBizIns.PageUser(&param.PageParam{Page: 1, PageSize: 10})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("分页查询结果", obj.ObjToJsonStr(pages))
	}

	entities, err := biz.UserQueryBizIns.ListUserByIds([]int64{1, 2, 3})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("列表查询结果", obj.ObjToJsonStr(entities))
	}

	//业务代码结束
	log.Println("after biz")
	err = database.CloseEngine()
	if err != nil {
		log.Println(err)
		return
	}
}
