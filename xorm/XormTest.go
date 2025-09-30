package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"zyj.com/golang-study/util/str"
	"zyj.com/golang-study/xorm/biz"
	"zyj.com/golang-study/xorm/database"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/param"
)

func main() {
	//初始化数据库
	err := database.Init("mysql", "root:zj123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	//database.GetEngine().Sync2(new(model.User))

	//业务代码开始
	err = biz.UserCmdBizIns.CreateUser(&param.UserCreate{Name: "zyj2fdsa", Email: "zyj000@163kkkk.com"})
	if err != nil {
		log.Println("error", err)
	}

	err = biz.UserCmdBizIns.UpdateUser(&model.User{Id: 4, Name: "4name"})
	if err != nil {
		log.Println(err)
	}

	in, err := biz.UserQueryBizIns.LogIn(&param.UserLogin{Email: "zyj@163fff.com", Pwd: "123456"})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("登录状态 %v", in)
	}

	pages, err := biz.UserQueryBizIns.PageListUser(&param.PageParam{Page: 1, PageSize: 10})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("分页查询结果", str.ObjToJson(pages))
	}
	//业务代码结束

	log.Printf("after biz")
}
