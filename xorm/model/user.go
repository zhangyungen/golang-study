// model/user.go
package model

import "time"

type User struct {
	Id        int64     `xorm:"'id' pk autoincr" json:"id"`
	Name      string    `xorm:"'name' varchar(100) notnull" json:"name"`
	Email     string    `xorm:"'email' varchar(100) notnull unique" json:"email"`
	Age       int       `xorm:"'age' int" json:"age"`
	Pwd       string    `xorm:"'pwd' varchar(255)" json:"pwd"`
	Status    int       `xorm:"'status' int default 1" json:"status"`
	CreatedAt time.Time `xorm:"'created_at' created" json:"createdAt"`
	UpdatedAt time.Time `xorm:"'updated_at' updated" json:"updatedAt"`
	Version   int       `xorm:"'version' version" json:"version"`
	//注意标记数据库字段必须用单引号包起来
	//notnull 也要标记，否者会出现插入nil情况，
	Deleted int8 `xorm:"'deleted' tinyint notnull default 0" json:"deleted"`
}

func (User) TableName() string {
	return "users"
}
