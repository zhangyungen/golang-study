package param

import (
	"github.com/jinzhu/copier"
	"log"
	"zyj.com/golang-study/xorm/model"
)

type UserCreate struct {
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
	Age   int    `json:"age" form:"age"`
}

type UserUpdate struct {
	Id     int64  ` json:"id"`
	Name   string ` json:"name"`
	Email  string ` json:"email"`
	Age    int    ` json:"age"`
	Pwd    string `json:"pwd"`
	Status int    ` json:"status"`
}

type UserLogin struct {
	Email string `json:"email" form:"email"`
	Pwd   string `json:"pwd" form:"pwd"`
}

func ConvertToModel(user *UserCreate) *model.User {
	userReturn := &model.User{}
	err := copier.Copy(userReturn, &user) // 核心拷贝操作
	if err != nil {
		log.Println("copier error ", err)
		return userReturn
	}
	return userReturn
}
