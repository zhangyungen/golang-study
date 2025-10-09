package main

import (
	"fmt"
	"github.com/petermattis/goid"
	"strconv"
	"time"
)

var (
	stop = make(chan struct{}) // 告诉 goroutine 停止
	done = make(chan struct{}) // 告诉我们 goroutine 退出了
)

func main() {
	id := goid.Get() // 直接获取当前 goroutine 的 Id
	fmt.Printf("Current goroutine Id:    %d\n", id)
	go func() {
		defer close(done)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for i := 10; i < 1000; i++ {
			fmt.Println("go routine print" + strconv.Itoa(i))
		}
		//for {
		//	select {
		//	case <-ticker.C:
		//		fmt.Println("tick")
		//	case <-stop:
		//		fmt.Println("stop")
		//		return
		//	}
		//}
	}()
	// 其它...
	time.Sleep(5 * time.Second)
	close(stop) // 指示 goroutine 停止
	<-done      // and wait for it to exit
	println("goroutine done")
}
