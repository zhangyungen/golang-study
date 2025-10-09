package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

func main() {
	go func() {
		// pprof 服务器，将暴露在 6060 端口
		if err := http.ListenAndServe(":6060", nil); err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			fmt.Printf("当前 goroutine 数量: %d\n", runtime.NumGoroutine())
			time.Sleep(1 * time.Second)
		}
	}()
	// 一个模拟的 goroutine 泄漏：不断创建永不退出的 goroutine
	for {
		go func() {
			select {} // 新创建的 goroutine 永远阻塞在此
		}()
		time.Sleep(2 * time.Second) // 控制创建速度，便于观察
	}

	var wg sync.WaitGroup
	ch := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ch // 这个 goroutine 会一直阻塞在此，因为 ch 永远不会被关闭或发送数据
	}()
	wg.Wait() // 主 goroutine}
}
