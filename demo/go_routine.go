package main

import (
	"fmt"
	"github.com/petermattis/goid"
)

func main() {
	id := goid.Get() // 直接获取当前 goroutine 的 ID
	fmt.Printf("Current goroutine ID: %d\n", id)

}
