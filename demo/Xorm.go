package main

import (
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var eg *xorm.EngineGroup

func main() {
	var err error
	master, err := xorm.NewEngine("postgres", "postgres://postgres:root@localhost:5432/test?sslmode=disable")
	if err != nil {
		return
	}

	slave1, err := xorm.NewEngine("postgres", "postgres://postgres:root@localhost:5432/test1?sslmode=disable")
	if err != nil {
		return
	}

	slave2, err := xorm.NewEngine("postgres", "postgres://postgres:root@localhost:5432/test2?sslmode=disable")
	if err != nil {
		return
	}

	slaves := []*xorm.Engine{slave1, slave2}
	eg, err = xorm.NewEngineGroup(master, slaves)
}
