package main

import (
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	orm, err := xorm.NewEngine("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	defer orm.Close()
	orm.Transaction(func(session *xorm.Session) (interface{}, error) {
		session.Begin()

		return nil, nil
	})
}
