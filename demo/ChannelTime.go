package main

import (
	"log"
	"time"
)

func main() {
	//ch := make(chan struct{})
	//
	////执行任务
	//go runJob(ch)
	//
	//select {
	//case <-ch:
	//	fmt.Println("任务完成!")
	//case <-time.After(2 * time.Second): //2秒超时
	//	fmt.Println("超时!")
	//}

	ch := make(chan struct{})

	//执行任务
	go runJobBeat(ch)

	select {}
}

func runJob(c chan<- struct{}) {
	//执行业务逻辑
	//...
	time.Sleep(3 * time.Second) // 模拟处理业务话费时间，给3秒，会输出超时，改为1秒，就会输出任务完成
	//执行完成后给信号
	c <- struct{}{}
}

func runJobBeat(c chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-c:
			//处理业务逻辑
		case <-ticker.C:
			log.Println("心跳一次")
		}
	}
}
