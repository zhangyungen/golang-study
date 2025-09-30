package model

import (
	"time"
)

type UserLoginLog struct {
	Id          int64     `xorm:"not null pk comment('主键id') BIGINT"`
	LoginIp     string    `xorm:"comment('上一次登入ip') VARCHAR(63)"`
	LoginTime   time.Time `xorm:"comment('上一次登录时间') DATETIME"`
	UserId      int64     `xorm:"not null comment('用户id') INT"`
	CreatedTime time.Time `xorm:"DATETIME"`
	UpdatedTime time.Time `xorm:"DATETIME"`
	DeletedTime time.Time `xorm:"DATETIME"`
}

func (m *UserLoginLog) TableName() string {
	return "user_login_log"
}
