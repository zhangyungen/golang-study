package main

import (
	"fmt"
	"zyj.com/golang-study/util/datautil"
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
	Name   string
	Age    int
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
		{ID: 4, Name: "Charlie4", Age: 3544},
	}

	orders := []Order{
		{ID: 1, UserID: 1, Amount: 100.0},
		{ID: 2, UserID: 3, Amount: 200.0},
		{ID: 4, UserID: 3, Amount: 400.0},
	}

	// 使用方法一

	// 使用方法二
	fmt.Println("\n=== 方法二：灵活泛型 ===")
	result2 := datautil.MergeSlicesGeneric(users, orders,
		func(u User) int { return u.ID },
		func(o Order) int { return o.UserID },
		datautil.LeftJoin,
	)
	for _, r := range result2 {
		fmt.Printf("Result: %+v\n", r)
	}
	result3 := datautil.LeftJoinData[Order, User, Order](orders, users,
		func(l Order) int { return l.UserID },
		func(r User) int { return r.ID },
		func(l Order, r User) Order {
			return Order{
				ID:     l.ID,
				UserID: l.UserID,
				Name:   r.Name,
				Age:    r.Age,
				Amount: l.Amount,
			}
		},
	)
	for _, r := range result3 {
		fmt.Printf("Result3: %+v\n", r)
	}

}
