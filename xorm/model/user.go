// model/user.go
package model

import "time"

type User struct {
	ID        int64     `xorm:"'id' pk autoincr" json:"id"`
	Name      string    `xorm:"'name' varchar(100) notnull" json:"name"`
	Email     string    `xorm:"'email' varchar(100) notnull unique" json:"email"`
	Age       int       `xorm:"'age' int" json:"age"`
	Status    int       `xorm:"'status' int default 1" json:"status"`
	CreatedAt time.Time `xorm:"'created_at' created" json:"created_at"`
	UpdatedAt time.Time `xorm:"'updated_at' updated" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
