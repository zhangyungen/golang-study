package main

import (
	"fmt"
	"zyj.com/golang-study/util/data"
)

type IDable interface {
	GetID() int
}

// 用户结构体
type User struct {
	ID   int
	Name string
	Age  int
}

func (u User) GetID() int {
	return u.ID
}

// 订单结构体
type Order struct {
	ID     int
	UserID int
	Amount float64
}

func (o Order) GetID() int {
	return o.ID
}

func main() {
	// 示例数据
	users := []User{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
		{ID: 3, Name: "Charlie", Age: 35},
	}

	orders := []Order{
		{ID: 1, UserID: 1, Amount: 100.0},
		{ID: 2, UserID: 2, Amount: 200.0},
		{ID: 4, UserID: 4, Amount: 400.0},
	}

	// 使用方法一

	// 使用方法二
	fmt.Println("\n=== 方法二：灵活泛型 ===")
	result2 := data.MergeSlicesGeneric(users, orders,
		func(u User) int { return u.ID },
		func(o Order) int { return o.UserID },
		data.LeftJoin,
	)
	for _, r := range result2 {
		fmt.Printf("Result: %+v\n", r)
	}

}
