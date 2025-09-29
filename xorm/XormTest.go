package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"zyj.com/golang-study/xorm/database"
	"zyj.com/golang-study/xorm/model"
	"zyj.com/golang-study/xorm/service"
)

func main() {
	//"user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	database.Init("mysql", "root:zj123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	//database.GetEngine().Sync2(new(model.User))

	service.UserServiceInstance.CreateUser(&model.User{Name: "zyj", Email: "zyj@163.com"})

	log.Printf("after insert")
}
