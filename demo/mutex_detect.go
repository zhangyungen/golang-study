package main

import (
	"github.com/sasha-s/go-deadlock"
	"time"
)

var mu1, mu2 deadlock.Mutex

func workerA() {
	mu1.Lock()
	time.Sleep(10 * time.Millisecond)
	mu2.Lock() // 可能在此等待 mu2
	defer mu2.Unlock()
	defer mu1.Unlock()
}
func workerB() {
	mu2.Lock()
	time.Sleep(10 * time.Millisecond)
	mu1.Lock() // 可能在此等待 mu1 (与workerA顺序相反)
	defer mu1.Unlock()
	defer mu2.Unlock()

}

func main() {
	deadlock.Opts.DeadlockTimeout = 3 * time.Second // 设置锁超时时间
	go workerA()
	go workerB()
	time.Sleep(20 * time.Second) // 等待足够长时间以便观察
}
