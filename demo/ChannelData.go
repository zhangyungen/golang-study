package main

import (
	"fmt"
)

func main() {
	//c := make(chan int)
	//out := make(chan bool)
	//go func() {
	//	go func() {
	//		time.Sleep(time.Second)
	//		c <- 2
	//	}()
	//	select {
	//	case <-time.After(time.Second * 1):
	//		fmt.Println("超时了")
	//	case v := <-c:
	//		fmt.Println("c中读取出的数据：", v)
	//	}
	//	out <- true
	//}()
	//<-out

	chanTest()
}

func chanTest() {
	var ch = make(chan int, 1)
	for i := 0; i < 10; i++ {
		select {
		case ch <- i:
			fmt.Println("send to chan ", i)
		case v := <-ch:
			fmt.Println("read from chan ", v)
		}
	}
}
