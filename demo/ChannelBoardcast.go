package main

import (
	"fmt"
	"sync"
)

type signal struct{}

func spawnGroup(count int, groupSignal chan struct{}) <-chan signal {
	c := make(chan signal) //用于让main goroutine阻塞的channel

	var wg sync.WaitGroup
	wg.Add(count)

	//创建goroutine
	for i := 1; i <= count; i++ {
		go func(index int) {
			<-groupSignal //等待main goroutine通知执行
			//处理业务逻辑
			//...
			fmt.Println(index, "子goroutine任务执行完成")

			wg.Done()
		}(i)
	}

	go func() {
		// 等待所有子goroutine执行完成
		wg.Wait()
		c <- signal{} // 通知main goroutine
	}()

	return c
}

func main() {

	groupSignal := make(chan struct{})
	c := spawnGroup(5, groupSignal)

	//执行业务逻辑
	//...

	fmt.Println("给子goroutine发送执行信号")
	// main goroutine任务执行完成，通知执行相关逻辑
	close(groupSignal) //通知刚创建的所有 routine

	<-c // 等待子goroutine执行完成

	fmt.Println("所有任务都执行完成")
}
