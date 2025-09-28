package main

import (
	"fmt"
	"time"
)

// 通过限制 limit 的缓存数量，决定并发时有多少协程在并行运行
var limit = make(chan struct{}, 3)

/*
*
Time: 2024-07-13T11:44:01.9047944+08:00, Goroutine: 0 exec i : 0, v: 1
Time: 2024-07-13T11:44:01.9047944+08:00, Goroutine: 11 exec i : 11, v: 12
Time: 2024-07-13T11:44:01.9047944+08:00, Goroutine: 4 exec i : 4, v: 5
Time: 2024-07-13T11:44:02.9200384+08:00, Goroutine: 8 exec i : 8, v: 9
Time: 2024-07-13T11:44:02.9200384+08:00, Goroutine: 9 exec i : 9, v: 10
Time: 2024-07-13T11:44:02.9201398+08:00, Goroutine: 10 exec i : 10, v: 11
Time: 2024-07-13T11:44:03.9207458+08:00, Goroutine: 5 exec i : 5, v: 6
Time: 2024-07-13T11:44:03.9207458+08:00, Goroutine: 1 exec i : 1, v: 2
Time: 2024-07-13T11:44:03.9207458+08:00, Goroutine: 2 exec i : 2, v: 3
Time: 2024-07-13T11:44:04.9240989+08:00, Goroutine: 6 exec i : 6, v: 7
Time: 2024-07-13T11:44:04.9243366+08:00, Goroutine: 3 exec i : 3, v: 4
Time: 2024-07-13T11:44:04.9243366+08:00, Goroutine: 7 exec i : 7, v: 8
Time: 2024-07-13T11:44:07+08:00, 主线程退出！
*/
func main() {
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	for i, v := range tasks {
		// 每个task开启一个协程
		go func(i, v int) {
			// 通过chan控制并发
			limit <- struct{}{}
			defer func() {
				<-limit
			}()
			// 具体的任务执行
			fmt.Printf("Time: %v, Goroutine: %v exec i : %d, v: %v\n", time.Now().Format(time.RFC3339Nano), i, i, v)
			time.Sleep(time.Second)

		}(i, v)
	}

	time.Sleep(6 * time.Second)
	fmt.Printf("Time: %v, 主线程退出！", time.Now().Format(time.RFC3339))
}
