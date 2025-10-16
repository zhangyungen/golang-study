package main

import "fmt"

type Coder interface {
	code()
}

type Gopher struct {
	name string
}

func (g Gopher) code() {
	fmt.Printf("%s is coding\n", g.name)
}

func main() {
	// nil 为引用类型零值 接口是引用类型

	var c Coder
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)

	//申明 Gopher 为零值的指针变量
	var g *Gopher
	var f Gopher

	fmt.Println("Gopher nil", g == nil)

	fmt.Println("&Gopher nil", &g == nil) // 取指针的指针无法编译  false

	//fmt.Println("*Gopher nil", *g == nil) // 取指针的指针无法编译
	fmt.Println("*Gopher nil", f) // 取指针的指针无法编译

	c = g
	fmt.Println(c == nil)
	fmt.Printf("c: %T, %v\n", c, c)
}
